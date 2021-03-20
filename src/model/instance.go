/**
 * @Author: lzw5399
 * @Date: 2021/3/9 14:18
 * @Desc:
 */
package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type ProcessInstance struct {
	AuditableBase
	Title               string         `gorm:"type:text" json:"title" form:"title"`                                                  // 工单标题
	Priority            int            `gorm:"type:smallint" json:"priority" form:"priority"`                                        // 工单优先级 1，正常 2，紧急 3，非常紧急
	ProcessDefinitionId int            `gorm:"type:integer" json:"processDefinitionId" form:"processDefinitionId"`                   // 流程ID
	ClassifyId          int            `gorm:"type:integer" json:"classifyId" form:"classifyId"`                                     // 分类ID
	IsEnd               bool           `gorm:"default:false" json:"isEnd" form:"isEnd"`                                              // 是否结束
	IsDenied            bool           `gorm:"default:false" json:"isDenied" form:"isDenied"`                                        // 是否被拒绝
	State               datatypes.JSON `gorm:"type:jsonb" json:"state" form:"state"`                                                 // 状态信息
	RelatedPerson       pq.Int64Array  `gorm:"type:integer[]; default:array[]::integer[]" json:"relatedPerson" form:"relatedPerson"` // 工单所有处理人
	UrgeCount           int            `gorm:"type:integer; default:0" json:"urgeCount" form:"urgeCount"`                            // 催办次数
	UrgeLastTime        time.Time      `gorm:"type:timestamp" json:"urgeLastTime" form:"urgeLastTime"`                  // 上一次催促时间
	TenantId            int            `gorm:"index" json:"tenantId" form:"tenantId"`                                                // 租户id
}
