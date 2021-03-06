/**
 * @Author: lzw5399
 * @Date: 2021/3/20 21:59
 * @Desc: 用户角色控制器(外部系统的)
 */
package controller

import (
	"github.com/labstack/echo/v4"

	"workflow/src/global/response"
	"workflow/src/model/request"
	"workflow/src/service"
	"workflow/src/util"
)

// @Tags role-users
// @Summary 批量同步(创建或更新)角色用户映射关系
// @Accept  json
// @Produce json
// @param request body request.BatchSyncUserRoleRequest true "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/role-users/_batch [POST]
func BatchSyncRoleUsers(c echo.Context) error {
	var (
		r   request.BatchSyncUserRoleRequest
		err error
	)

	if err = c.Bind(&r); err != nil {
		return response.BadRequest(c)
	}

	tenantId := util.GetCurrentTenantId(c)
	err = service.BatchSyncRoleUsers(&r, tenantId)
	if err != nil {
		return response.Failed(c, err)
	}

	return response.OkWithMessage(c, "更新成功")
}

// @Tags role-users
// @Summary 同步(创建或更新)角色用户映射关系
// @Accept  json
// @Produce json
// @param request body request.SyncRoleUsersRequest true "request"
// @param WF-TENANT-CODE header string true "WF-TENANT-CODE"
// @param WF-CURRENT-USER header string true "WF-CURRENT-USER"
// @Success 200 {object} response.HttpResponse
// @Router /api/wf/role-users [POST]
//func SyncRoleUsers(c echo.Context) error {
//	var (
//		r   request.SyncRoleUsersRequest
//		err error
//	)
//
//	if err = c.Bind(&r); err != nil {
//		return response.BadRequest(c)
//	}
//
//	tenantId := util.GetCurrentTenantId(c)
//	err = service.SyncRoleUsers(&r, tenantId)
//	if err != nil {
//		return response.Failed(c, err)
//	}
//
//	return response.OkWithMessage(c, "更新成功")
//}
