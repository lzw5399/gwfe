/**
 * @Author: lzw5399
 * @Date: 2021/3/20 15:44
 * @Desc:
 */
package initialize

import (
	"log"

	"github.com/patrickmn/go-cache"

	"workflow/src/global"
	"workflow/src/model"
)

// 初始化租户缓存
func init() {
	c := cache.New(-1, -1)
	global.BankCache = c

	var tenants []model.Tenant
	err := global.BankDb.
		Model(&model.Tenant{}).
		Find(&tenants).
		Error
	if err != nil {
		log.Fatalf("初始化租户失败, err:%s", err.Error())
	}

	c.SetDefault("tenants", tenants)
	log.Printf("-------租户列表缓存成功，当前租户个数:%d--------\n", len(tenants))
}
