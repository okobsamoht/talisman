package rest

import (
	"github.com/okobsamoht/tomato/cache"
	"github.com/okobsamoht/tomato/cloud"
	"github.com/okobsamoht/tomato/livequery"
	"github.com/okobsamoht/tomato/orm"
	"github.com/okobsamoht/tomato/types"
	"github.com/okobsamoht/tomato/utils"
)

// Destroy 删除对象
type Destroy struct {
	auth         *Auth
	className    string
	query        types.M
	originalData types.M
}

// NewDestroy 组装 Destroy
func NewDestroy(
	auth *Auth,
	className string,
	query types.M,
	originalData types.M,
) *Destroy {
	destroy := &Destroy{
		auth:         auth,
		className:    className,
		query:        query,
		originalData: originalData,
	}
	return destroy
}

// Execute 执行删除请求
func (d *Destroy) Execute() error {
	err := d.handleSession()
	if err != nil {
		return err
	}
	err = d.runBeforeTrigger()
	if err != nil {
		return err
	}
	err = d.handleUserRoles()
	if err != nil {
		return err
	}
	err = d.runDestroy()
	if err != nil {
		return err
	}
	err = d.runAfterTrigger()
	if err != nil {
		return err
	}
	return nil
}

// handleSession 处理 _Session 表的删除操作
func (d *Destroy) handleSession() error {
	if d.className != "_Session" {
		return nil
	}
	if sessionToken := utils.S(d.originalData["sessionToken"]); sessionToken != "" {
		cache.User.Del(sessionToken)
	}

	return nil
}

// runBeforeTrigger 执行删前回调
func (d *Destroy) runBeforeTrigger() error {
	if d.originalData == nil {
		return nil
	}
	if livequery.TLiveQuery != nil {
		livequery.TLiveQuery.OnAfterDelete(d.className, d.originalData, nil)
	}

	d.originalData["className"] = d.className
	maybeRunTrigger(cloud.TypeBeforeDelete, d.auth, d.originalData, nil)

	return nil
}

// handleUserRoles 获取用户角色信息
func (d *Destroy) handleUserRoles() error {
	if d.auth.IsMaster == false {
		d.auth.GetUserRoles()
	}

	return nil
}

// runDestroy 添加 acl 字段，并执行删除对象操作
func (d *Destroy) runDestroy() error {
	options := types.M{}
	if d.auth.IsMaster == false {
		acl := []string{"*"}
		if d.auth.User != nil {
			acl = append(acl, utils.S(d.auth.User["objectId"]))
			acl = append(acl, d.auth.UserRoles...)
		}
		options["acl"] = acl
	}
	return orm.TomatoDBController.Destroy(d.className, d.query, options)
}

// runAfterTrigger 执行删后回调
func (d *Destroy) runAfterTrigger() error {
	maybeRunTrigger(cloud.TypeAfterDelete, d.auth, d.originalData, nil)
	return nil
}
