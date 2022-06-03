package rest

import (
	"reflect"
	"testing"

	"github.com/okobsamoht/talisman/cloud"
	"github.com/okobsamoht/talisman/errs"
	"github.com/okobsamoht/talisman/types"
	"github.com/okobsamoht/talisman/utils"
)

func Test_maybeRunTrigger(t *testing.T) {
	var result types.M
	var err error
	var expect types.M
	var expectErr error
	/****************************************************************************************/
	cloud.BeforeSave("user", func(request cloud.TriggerRequest, response cloud.Response) {
		object := request.Object
		if username := utils.S(object["username"]); username != "" {
			object["username"] = username + "_talisman"
			response.Success(nil)
		} else {
			response.Error(1, "need a username")
		}
	})
	_, err = maybeRunTrigger(cloud.TypeBeforeSave, Master(), types.M{"className": "user"}, nil)
	expectErr = errs.E(1, "need a username")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	result, err = maybeRunTrigger(cloud.TypeBeforeSave, Master(), types.M{"className": "user", "username": "joe"}, nil)
	expect = types.M{
		"object": types.M{
			"className": "user",
			"username":  "joe_talisman",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	cloud.UnregisterAll()
}
