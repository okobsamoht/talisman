package controllers

import (
	"github.com/okobsamoht/talisman/errs"
	"github.com/okobsamoht/talisman/rest"
	"github.com/okobsamoht/talisman/types"
	"github.com/okobsamoht/talisman/utils"
)

// SessionsController 处理 /sessions 接口的请求
type SessionsController struct {
	ClassesController
}

// HandleFind 处理查找 session 请求
// @router / [get]
func (s *SessionsController) HandleFind() {
	s.ClassName = "_Session"
	s.ClassesController.HandleFind()
}

// HandleGet 处理获取指定 session 请求
// @router /:objectId [get]
func (s *SessionsController) HandleGet() {
	s.ClassName = "_Session"
	s.ObjectID = s.Ctx.Input.Param(":objectId")
	s.ClassesController.HandleGet()
}

// HandleCreate 处理 session 创建请求
// @router / [post]
func (s *SessionsController) HandleCreate() {
	s.ClassName = "_Session"
	s.ClassesController.HandleCreate()
}

// HandleUpdate 处理更新指定 session 请求
// @router /:objectId [put]
func (s *SessionsController) HandleUpdate() {
	s.ClassName = "_Session"
	s.ObjectID = s.Ctx.Input.Param(":objectId")
	s.ClassesController.HandleUpdate()
}

// HandleDelete 处理删除指定 session 请求
// @router /:objectId [delete]
func (s *SessionsController) HandleDelete() {
	objectID := s.Ctx.Input.Param(":objectId")
	if objectID == "me" {
		s.ClassesController.Delete()
		return
	}
	s.ClassName = "_Session"
	s.ObjectID = objectID
	s.ClassesController.HandleDelete()
}

// HandleGetMe 处理当前请求 session
// @router /me [get]
func (s *SessionsController) HandleGetMe() {
	if s.Info == nil || s.Info.SessionToken == "" {
		s.HandleError(errs.E(errs.InvalidSessionToken, "Session token required."), 0)
		return
	}
	where := types.M{
		"sessionToken": s.Info.SessionToken,
	}
	response, err := rest.Find(rest.Master(), "_Session", where, types.M{}, s.Info.ClientSDK)
	if err != nil {
		s.HandleError(err, 0)
		return
	}
	if utils.HasResults(response) == false {
		s.HandleError(errs.E(errs.InvalidSessionToken, "Session token not found."), 0)
		return
	}
	results := utils.A(response["results"])
	s.Data["json"] = results[0]
	s.ServeJSON()
}

// HandleUpdateMe 仅用于更新 installationId
// @router /me [put]
func (s *SessionsController) HandleUpdateMe() {
	if s.Info == nil || s.Info.SessionToken == "" {
		s.HandleError(errs.E(errs.InvalidSessionToken, "Session token required."), 0)
		return
	}
	if s.Info.InstallationID == "" {
		s.Data["json"] = types.M{}
		s.ServeJSON()
		return
	}
	where := types.M{
		"sessionToken": s.Info.SessionToken,
	}
	response, err := rest.Find(rest.Master(), "_Session", where, types.M{}, s.Info.ClientSDK)
	if err != nil {
		s.HandleError(err, 0)
		return
	}
	if utils.HasResults(response) == false {
		s.HandleError(errs.E(errs.InvalidSessionToken, "Session token not found."), 0)
		return
	}
	results := utils.A(response["results"])
	session := utils.M(results[0])
	update := types.M{"installationId": s.Info.InstallationID}
	result, err := rest.Update(rest.Master(), "_Session", utils.S(session["objectId"]), update, nil)
	if err != nil {
		s.HandleError(err, 0)
		return
	}
	s.Data["json"] = result["response"]
	s.ServeJSON()
}

// Put ...
// @router / [put]
func (s *SessionsController) Put() {
	s.ClassesController.Put()
}

// Delete ...
// @router / [delete]
func (s *SessionsController) Delete() {
	s.ClassesController.Delete()
}
