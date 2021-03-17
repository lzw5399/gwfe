/**
 * @Author: lzw5399
 * @Date: 2021/1/16 22:58
 * @Desc: 流程实例服务
 */
package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"

	"workflow/src/global"
	"workflow/src/global/shared"
	"workflow/src/model"
	"workflow/src/model/request"
	"workflow/src/model/response"
	"workflow/src/service/engine"
	"workflow/src/util"
)

type InstanceService interface {
	CreateProcessInstance(*request.ProcessInstanceRequest, uint) (uint, error)
	Get(uint) (*model.ProcessInstance, error)
	List(*request.PagingRequest) (*response.PagingResponse, error)
}

type instanceService struct {
}

func NewInstanceService() *instanceService {
	return &instanceService{}
}

// 创建实例
func (i *instanceService) CreateProcessInstance(r *request.ProcessInstanceRequest, currentUserId uint) (uint, error) {
	var (
		variableValue     []map[string]interface{} // 变量值
		err               error
		processDefinition model.ProcessDefinition // 流程模板
		processInstance   = r.ToProcessInstance(currentUserId)
		processEngine     *engine.ProcessEngine  // 流程定义引擎
		instanceEngine    *engine.InstanceEngine // 流程实例引擎
		condExprStatus    bool
		sourceEdges       []map[string]interface{}
		targetEdges       []map[string]interface{}
	)

	// 查询对应的流程模板
	err = global.BankDb.
		Where("id = ?", processInstance.ProcessDefinitionId).
		Find(&processDefinition).
		Error
	if err != nil {
		return 0, err
	}

	// 实例化流程引擎
	processEngine, err = engine.NewProcessEngine(processDefinition)
	if err != nil {
		return 0, err
	}

	instanceEngine, err = engine.NewInstanceEngine(processDefinition)
	if err != nil {
		return 0, err
	}

	// 将初始状态赋值给当前的流程实例
	if state, err := instanceEngine.GetInstanceInitialState(); err != nil {
		return 0, err
	} else {
		processInstance.State = util.MarshalToBytes(state)
	}

	// 把对应的流程模板的structure单独反序列化出来处理
	nodeValue, err := processEngine.GetNode(variableValue[0]["id"].(string))
	if err != nil {
		return 0, err
	}

	// 获取变量值
	err = json.Unmarshal(processInstance.State, &variableValue)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	// 填充process instance的当前状态state字段

	// TODO 省略了processInstance.State针对变量的预处理

	switch nodeValue["clazz"] {
	// 排他网关
	case "exclusiveGateway":
		var sourceEdges []map[string]interface{}
		sourceEdges, err = processEngine.GetEdge(nodeValue["id"].(string), "source")
		if err != nil {
			return 0, err
		}
	breakTag:
		for _, edge := range sourceEdges {
			edgeCondExpr := make([]map[string]interface{}, 0)
			err = json.Unmarshal([]byte(edge["conditionExpression"].(string)), &edgeCondExpr)
			if err != nil {
				return 0, err
			}
			for _, condExpr := range edgeCondExpr {
				// 条件判断
				condExprStatus, err = processEngine.ConditionalJudgment(condExpr)
				if err != nil {
					return 0, err
				}
				if condExprStatus {
					// 进行节点跳转
					nodeValue, err = processEngine.GetNode(edge["target"].(string))
					if err != nil {
						return 0, err
					}

					if nodeValue["clazz"] == "userTask" || nodeValue["clazz"] == "receiveTask" {
						if nodeValue["assignValue"] == nil || nodeValue["assignType"] == "" {
							err = errors.New("处理人不能为空")
							return 0, err
						}
					}
					variableValue[0]["id"] = nodeValue["id"].(string)
					variableValue[0]["label"] = nodeValue["label"]
					variableValue[0]["processor"] = nodeValue["assignValue"]
					variableValue[0]["process_method"] = nodeValue["assignType"]
					break breakTag
				}
			}
		}
		if !condExprStatus {
			return 0, errors.New("所有流转均不符合条件，请确认。")
		}
	// 并行网关
	case "parallelGateway":
		// 入口，判断
		sourceEdges, err = processEngine.GetEdge(nodeValue["id"].(string), "source")
		if err != nil {
			return 0, fmt.Errorf("查询流转信息失败，%v", err.Error())
		}

		targetEdges, err = processEngine.GetEdge(nodeValue["id"].(string), "target")
		if err != nil {
			return 0, fmt.Errorf("查询流转信息失败，%v", err.Error())
		}

		if len(sourceEdges) > 0 {
			nodeValue, err = processEngine.GetNode(sourceEdges[0]["target"].(string))
			if err != nil {
				return 0, err
			}
		} else {
			return 0, errors.New("并行网关流程不正确")
		}

		if len(sourceEdges) > 1 && len(targetEdges) == 1 {
			// 入口
			variableValue = []map[string]interface{}{}
			for _, edge := range sourceEdges {
				targetStateValue, err := processEngine.GetNode(edge["target"].(string))
				if err != nil {
					return 0, err
				}
				variableValue = append(variableValue, map[string]interface{}{
					"id":             edge["target"].(string),
					"label":          targetStateValue["label"],
					"processor":      targetStateValue["assignValue"],
					"process_method": targetStateValue["assignType"],
				})
			}
		} else {
			return 0, errors.New("并行网关流程配置不正确")
		}
	}

	// 变量的预处理
	err = preprocessVariables(variableValue, currentUserId)
	if err != nil {
		log.Error(err)
		return 0, errors.New("获取处理人变量值失败")
	}

	// 将最新的“变量/状态信息”赋值给processInstance
	processInstance.State, err = json.Marshal(variableValue)
	if err != nil {
		return 0, err
	}

	// processInstance某些字段更新
	relatedPerson, _ := json.Marshal([]uint{currentUserId})
	processInstance.RelatedPerson = relatedPerson

	// 开启事务
	err = global.BankDb.Transaction(func(tx *gorm.DB) error {

		// 创建
		err = tx.Create(&processInstance).Error
		if err != nil {
			return fmt.Errorf("创建工单失败，%v", err.Error())
		}

		// todo 省略了【创建工单模版关联数据】

		// todo 省略了【获取当前用户信息】

		// 创建历史记录
		var stateList []interface{}
		err = json.Unmarshal(processInstance.State, &stateList)
		if err != nil {
			return fmt.Errorf("json序列化失败，%s", err.Error())
		}

		err = tx.Create(&model.CirculationHistory{
			Title:             processInstance.Title,
			ProcessInstanceId: processInstance.Id,
			//State:             r.SourceState, // todo带填充
			//Source:            r.Source, // todo带填充
			Target:      stateList[0].(map[string]interface{})["id"].(string),
			Circulation: "新建",
			Processor:   "", // todo 上面两个也要
			ProcessorId: currentUserId,
		}).Error
		if err != nil {
			return fmt.Errorf("新建历史记录失败，%v", err.Error())
		}

		// 更新process_definition表的提交数量统计
		err = tx.Model(&model.ProcessDefinition{}).
			Where("id = ?", processInstance.ProcessDefinitionId).
			Update("submit_count", processDefinition.SubmitCount+1).Error
		if err != nil {
			return fmt.Errorf("更新流程提交数量统计失败，%v", err.Error())
		}

		return nil
	})

	// todo 暂时省略了发送通知

	// todo 暂时省略了执行脚本任务

	return 0, err
}

// 获取单个ProcessInstance
func (i *instanceService) Get(instanceId uint) (*model.ProcessInstance, error) {
	var instance model.ProcessInstance
	err := global.BankDb.Where("id=?", instanceId).First(&instance).Error

	return &instance, err
}

// 获取ProcessInstance列表
func (i *instanceService) List(r *request.PagingRequest) (*response.PagingResponse, error) {
	var instances []model.ProcessInstance
	db := shared.ApplyPaging(global.BankDb, r)
	err := db.Find(&instances).Error

	var totalCount int64
	global.BankDb.Model(&model.ProcessInstance{}).Count(&totalCount)

	return &response.PagingResponse{
		TotalCount:   totalCount,
		CurrentCount: int64(len(instances)),
		Data:         &instances,
	}, err
}

// 获取实例的某一个变量
//func (i *instanceService) GetVariable(r *request.GetVariableRequest) (*response.InstanceVariableResponse, error) {
//	var str string
//	err := global.BankDb.Model(&model.ProcessInstance{}).Where("id=?", r.InstanceId).Select("variables").First(&str).Error
//	if err != nil {
//		return nil, errors.New("当前instance实例不存在")
//	}
//
//	var variables []model.InstanceVariable
//	err = json.Unmarshal([]byte(str), &variables)
//	if err != nil {
//		return nil, errors.New("获取失败")
//	}
//
//	for _, v := range variables {
//		if v.Name == r.VariableName {
//			return &response.InstanceVariableResponse{
//				Name:  v.Name,
//				Type:  v.Type,
//				Value: v.Value,
//			}, nil
//		}
//	}
//
//	return nil, errors.New("变量不存在")
//}

//func (i *instanceService) ListVariables(r *request.GetVariableListRequest) (*response.PagingResponse, error) {
//	var str string
//	err := global.BankDb.Model(&model.ProcessInstance{}).Where("id=?", r.InstanceId).Select("variables").First(&str).Error
//	if err != nil {
//		return nil, errors.New("当前instance实例不存在")
//	}
//
//	var variables []model.InstanceVariable
//	err = json.Unmarshal([]byte(str), &variables)
//	if err != nil {
//		return nil, errors.New("获取失败")
//	}
//
//	//var finalVariables []model.InstanceVariable
//	//util.NewPaging(variables).Offset(r.Offset).Limit(r.Limit).Get(&finalVariables)
//
//	return &response.PagingResponse{
//		TotalCount:   int64(len(variables)),
//		CurrentCount: int64(len(variables)),
//		Data:         &variables,
//	}, err
//}

// 创建流程实例的时候预先处理变量转成实际处理人
func preprocessVariables(stateList []map[string]interface{}, creator uint) (err error) {
	//var (
	//	userInfo system.SysUser
	//	deptInfo system.Dept
	//)
	//
	//// 变量转为实际的数据
	//for _, stateItem := range stateList {
	//	if stateItem.(map[string]interface{})["process_method"] == "variable" {
	//		for processorIndex, processor := range stateItem.(map[string]interface{})["processor"].([]interface{}) {
	//			if int(processor.(float64)) == 1 {
	//				// 创建者
	//				stateItem.(map[string]interface{})["processor"].([]interface{})[processorIndex] = creator
	//			} else if int(processor.(float64)) == 2 {
	//				// 1. 查询用户信息
	//				err = orm.Eloquent.Model(&userInfo).Where("user_id = ?", creator).Find(&userInfo).Error
	//				if err != nil {
	//					return
	//				}
	//				// 2. 查询部门信息
	//				err = orm.Eloquent.Model(&deptInfo).Where("dept_id = ?", userInfo.DeptId).Find(&deptInfo).Error
	//				if err != nil {
	//					return
	//				}
	//
	//				// 3. 替换处理人信息
	//				stateItem.(map[string]interface{})["processor"].([]interface{})[processorIndex] = deptInfo.Leader
	//			}
	//		}
	//		stateItem.(map[string]interface{})["process_method"] = "person"
	//	}
	//}

	return
}
