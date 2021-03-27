/**
 * @Author: lzw5399
 * @Date: 2021/3/19 17:10
 * @Desc: 工单的流转相关方法
 */
package engine

import (
	"time"

	"workflow/src/global/constant"
	"workflow/src/model"
	"workflow/src/model/dto"
	"workflow/src/model/request"
	"workflow/src/util"
)

// 一般流转处理，兼顾了会签的判断
func (engine *ProcessEngine) CommonProcessing(newStates dto.StateArray) error {
	// 如果是拒绝的流程直接跳转
	if engine.linkEdge.FlowProperties == "0" {
		return engine.Circulation(newStates)
	}

	// TODO 暂不支持并行网关，这边先判断0
	state := engine.ProcessInstance.State[0]
	// 不是会签直接跳
	if !state.IsCounterSign {
		return engine.Circulation(newStates)
	}

	// 判断当前用户是否会签的最后一个人
	isLastPerson, err := engine.IsCounterSignLastProcessor()
	if err != nil {
		return err
	}
	if isLastPerson {
		return engine.Circulation(newStates)
	}

	// 不是会签的最后一个人, 则更新
	toUpdate := map[string]interface{}{
		"state":          engine.ProcessInstance.State,
		"update_time":    time.Now().Local(),
		"update_by":      engine.currentUserId,
		"related_person": engine.ProcessInstance.RelatedPerson,
		"variables":      engine.ProcessInstance.Variables,
	}

	err = engine.tx.Model(&engine.ProcessInstance).
		Updates(toUpdate).
		Error

	return err
}

// processInstance流转处理
func (engine *ProcessEngine) Circulation(newStates dto.StateArray) error {
	toUpdate := map[string]interface{}{
		"state":          newStates,
		"related_person": engine.ProcessInstance.RelatedPerson,
		"is_end":         false,
		"update_time":    time.Now().Local(),
		"update_by":      engine.currentUserId,
		"variables":      engine.ProcessInstance.Variables,
	}

	// 如果是跳转到结束节点，则需要修改节点状态
	if engine.targetNode.Clazz == constant.End {
		toUpdate["is_end"] = true
	}

	err := engine.tx.
		Model(&engine.ProcessInstance).
		Updates(toUpdate).
		Error

	return err
}

// 否决
func (engine *ProcessEngine) Deny(r *request.DenyInstanceRequest) error {
	// 获取最新的相关者RelatedPerson
	engine.UpdateRelatedPerson()

	// 更新instance字段
	toUpdate := map[string]interface{}{
		"related_person": engine.ProcessInstance.RelatedPerson,
		"is_denied":      true,
		"update_time":    time.Now().Local(),
		"update_by":      engine.currentUserId,
	}

	err := engine.tx.
		Model(&engine.ProcessInstance).
		Updates(toUpdate).
		Error

	// 获取上一条的流转历史的CreateTime来计算CostDuration
	var lastCirculation model.CirculationHistory
	err = engine.tx.
		Where("process_instance_id = ?", engine.ProcessInstance.Id).
		Order("create_time desc").
		Select("create_time").
		First(&lastCirculation).
		Error
	if err != nil {
		return err
	}
	duration := util.FmtDuration(time.Since(lastCirculation.CreateTime))

	// 创建新的一条流转历史
	// todo 这里先判断[0]
	state := engine.ProcessInstance.State[0]
	cirHistory := model.CirculationHistory{
		AuditableBase: model.AuditableBase{
			CreateBy: engine.currentUserId,
			UpdateBy: engine.currentUserId,
		},
		Title:             engine.ProcessInstance.Title,
		ProcessInstanceId: engine.ProcessInstance.Id,
		SourceState:       state.Label,
		SourceId:          state.Id,
		TargetId:          "",
		Circulation:       "否决",
		ProcessorId:       engine.currentUserId,
		CostDuration:      duration,
		Remarks:           r.Remarks,
	}

	err = engine.tx.
		Model(&model.CirculationHistory{}).
		Create(&cirHistory).
		Error

	if err != nil {
		return err
	}

	return err
}

// 更新relatedPerson
func (engine *ProcessEngine) UpdateRelatedPerson() {
	// 获取最新的相关者RelatedPerson
	exist := false
	for _, person := range engine.ProcessInstance.RelatedPerson {
		if uint(person) == engine.currentUserId {
			exist = true
			break
		}
	}
	if !exist {
		engine.ProcessInstance.RelatedPerson = append(engine.ProcessInstance.RelatedPerson, int64(engine.currentUserId))
	}
}