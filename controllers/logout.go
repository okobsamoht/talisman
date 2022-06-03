package controllers

import (
	"github.com/okobsamoht/tomato/rest"
	"github.com/okobsamoht/tomato/types"
	"github.com/okobsamoht/tomato/utils"
)

// LogoutController 处理 /logout 接口的请求
type LogoutController struct {
	ClassesController
}

// HandleLogOut 处理用户退出请求
// @router / [post]
func (l *LogoutController) HandleLogOut() {
	if l.Info != nil && l.Info.SessionToken != "" {
		where := types.M{
			"sessionToken": l.Info.SessionToken,
		}
		records, err := rest.Find(rest.Master(), "_Session", where, types.M{}, l.Info.ClientSDK)

		if err != nil {
			l.HandleError(err, 0)
			return
		}
		if utils.HasResults(records) {
			results := utils.A(records["results"])
			obj := utils.M(results[0])
			err := rest.Delete(rest.Master(), "_Session", utils.S(obj["objectId"]))
			if err != nil {
				l.HandleError(err, 0)
				return
			}
		}
	}
	l.Data["json"] = types.M{}
	l.ServeJSON()
}

// Get ...
// @router / [get]
func (l *LogoutController) Get() {
	l.ClassesController.Get()
}

// Delete ...
// @router / [delete]
func (l *LogoutController) Delete() {
	l.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (l *LogoutController) Put() {
	l.ClassesController.Put()
}
