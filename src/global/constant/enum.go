/**
 * @Author: lzw5399
 * @Date: 2021/1/16 11:21
 * @Desc:
 */
package constant

// event的类别常量
const (
	StartEvent = iota + 1
	EndEvent
)

// process instance 的type类别
const (
	I_MyToDo   = iota + 1 // 我的待办
	I_ICreated            // 我创建的
	I_IRelated            // 和我相关的
	I_All                 // 所有
)

// process definition 的type类别
const (
	D_ICreated = iota + 1
	D_All
)

type ChainNodeStatus int

const (
	Processed   ChainNodeStatus = iota + 1 // 已处理
	CurrentNode                            // 当前节点
	Unreachable                            // 后续节点
)

const (
	VariableNumber = iota + 1
	VariableString
	VariableBool
)

// 流转历史的类型
const (
	HistoryTypeFull = iota + 1
	HistoryTypeSimple
)
