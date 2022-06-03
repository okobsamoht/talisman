package rest

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/okobsamoht/tomato/config"
	"github.com/okobsamoht/tomato/errs"
	"github.com/okobsamoht/tomato/orm"
	"github.com/okobsamoht/tomato/types"
	"github.com/okobsamoht/tomato/utils"
)

func Test_HandleLoginAttempt(t *testing.T) {
	// TODO
}

func Test_notLocked(t *testing.T) {
	var username string
	var object, schema types.M
	var accountLockout *AccountLockout
	var err, expectErr error
	var expiresAtStr string
	/*****************************************************************/
	config.TConfig.AccountLockoutThreshold = 3
	config.TConfig.AccountLockoutDuration = 5
	expiresAtStr = utils.TimetoString(time.Now().UTC().Add(time.Duration(config.TConfig.AccountLockoutDuration) * time.Minute))
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "01",
		"username": username,
		"_account_lockout_expires_at": types.M{
			"__type": "Date",
			"iso":    expiresAtStr,
		},
		"_failed_login_count": 3,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	accountLockout = NewAccountLockout(username)
	err = accountLockout.notLocked()
	expectErr = errs.E(errs.ObjectNotFound, "Your account is locked due to multiple failed login attempts. Please try again after "+
		strconv.Itoa(config.TConfig.AccountLockoutDuration)+" minute(s)")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/*****************************************************************/
	config.TConfig.AccountLockoutThreshold = 3
	config.TConfig.AccountLockoutDuration = 5
	expiresAtStr = utils.TimetoString(time.Now().UTC().Add(time.Duration(config.TConfig.AccountLockoutDuration) * time.Minute))
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "01",
		"username": username,
		"_account_lockout_expires_at": types.M{
			"__type": "Date",
			"iso":    expiresAtStr,
		},
		"_failed_login_count": 1,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	accountLockout = NewAccountLockout(username)
	err = accountLockout.notLocked()
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/*****************************************************************/
	config.TConfig.AccountLockoutThreshold = 3
	config.TConfig.AccountLockoutDuration = 5
	expiresAtStr = utils.TimetoString(time.Now().UTC().Add(-time.Duration(config.TConfig.AccountLockoutDuration) * time.Minute))
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "01",
		"username": username,
		"_account_lockout_expires_at": types.M{
			"__type": "Date",
			"iso":    expiresAtStr,
		},
		"_failed_login_count": 3,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	accountLockout = NewAccountLockout(username)
	err = accountLockout.notLocked()
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_setFailedLoginCount(t *testing.T) {
	var username string
	var object, schema types.M
	var accountLockout *AccountLockout
	var err error
	var results, expect []types.M
	/*****************************************************************/
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "01",
		"username": username,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	accountLockout = NewAccountLockout(username)
	err = accountLockout.setFailedLoginCount(0)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = orm.Adapter.Find("_User", schema, types.M{}, types.M{})
	expect = []types.M{
		types.M{
			"objectId":            "01",
			"username":            username,
			"_failed_login_count": 0,
		},
	}
	if reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_handleFailedLoginAttempt(t *testing.T) {
	var username string
	var object, schema types.M
	var accountLockout *AccountLockout
	var err error
	var results, expect []types.M
	/*****************************************************************/
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "01",
		"username": username,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	accountLockout = NewAccountLockout(username)
	err = accountLockout.handleFailedLoginAttempt()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = orm.Adapter.Find("_User", schema, types.M{}, types.M{})
	expect = []types.M{
		types.M{
			"objectId":            "01",
			"username":            username,
			"_failed_login_count": 1,
		},
	}
	if reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
	/*****************************************************************/
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":            "01",
		"username":            username,
		"_failed_login_count": 2,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	config.TConfig.AccountLockoutThreshold = 3
	config.TConfig.AccountLockoutDuration = 5
	accountLockout = NewAccountLockout(username)
	err = accountLockout.handleFailedLoginAttempt()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = orm.Adapter.Find("_User", schema, types.M{}, types.M{})
	expect = []types.M{
		types.M{
			"objectId":            "01",
			"username":            username,
			"_failed_login_count": 3,
		},
	}
	if _, ok := results[0]["_account_lockout_expires_at"]; ok == false {
		t.Error("expect:", "_account_lockout_expires_at", "result:", "")
	}
	delete(results[0], "_account_lockout_expires_at")
	if reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_initFailedLoginCount(t *testing.T) {
	var username string
	var object, schema types.M
	var accountLockout *AccountLockout
	var err error
	var results, expect []types.M
	/*****************************************************************/
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "01",
		"username": username,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	accountLockout = NewAccountLockout(username)
	err = accountLockout.initFailedLoginCount()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = orm.Adapter.Find("_User", schema, types.M{}, types.M{})
	expect = []types.M{
		types.M{
			"objectId":            "01",
			"username":            username,
			"_failed_login_count": 0,
		},
	}
	if reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_incrementFailedLoginCount(t *testing.T) {
	var username string
	var object, schema types.M
	var accountLockout *AccountLockout
	var err error
	var results, expect []types.M
	/*****************************************************************/
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":            "01",
		"username":            username,
		"_failed_login_count": 0,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	accountLockout = NewAccountLockout(username)
	err = accountLockout.incrementFailedLoginCount()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = orm.Adapter.Find("_User", schema, types.M{}, types.M{})
	expect = []types.M{
		types.M{
			"objectId":            "01",
			"username":            username,
			"_failed_login_count": 1,
		},
	}
	if reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_setLockoutExpiration(t *testing.T) {
	var username string
	var object, schema types.M
	var accountLockout *AccountLockout
	var err error
	var results, expect []types.M
	/*****************************************************************/
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":            "01",
		"username":            username,
		"_failed_login_count": 1,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	config.TConfig.AccountLockoutThreshold = 3
	config.TConfig.AccountLockoutDuration = 5
	accountLockout = NewAccountLockout(username)
	err = accountLockout.setLockoutExpiration()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = orm.Adapter.Find("_User", schema, types.M{}, types.M{})
	expect = []types.M{
		types.M{
			"objectId":            "01",
			"username":            username,
			"_failed_login_count": 1,
		},
	}
	if reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
	/*****************************************************************/
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":            "01",
		"username":            username,
		"_failed_login_count": 3,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	config.TConfig.AccountLockoutThreshold = 3
	config.TConfig.AccountLockoutDuration = 5
	expiresAtStr := utils.TimetoString(time.Now().UTC().Add(time.Duration(config.TConfig.AccountLockoutDuration) * time.Minute))
	expiresAt, _ := utils.StringtoTime(expiresAtStr)
	accountLockout = NewAccountLockout(username)
	err = accountLockout.setLockoutExpiration()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = orm.Adapter.Find("_User", schema, types.M{}, types.M{})
	expect = []types.M{
		types.M{
			"objectId":                    "01",
			"username":                    username,
			"_failed_login_count":         3,
			"_account_lockout_expires_at": expiresAt.Local(),
		},
	}
	if reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_isFailedLoginCountSet(t *testing.T) {
	var username string
	var object, schema types.M
	var accountLockout *AccountLockout
	var isSet bool
	var err error
	/*****************************************************************/
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "01",
		"username": username,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	accountLockout = NewAccountLockout(username)
	isSet, err = accountLockout.isFailedLoginCountSet()
	if err != nil || isSet != false {
		t.Error("expect:", false, "result:", isSet, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/*****************************************************************/
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":            "01",
		"username":            username,
		"_failed_login_count": 3,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	accountLockout = NewAccountLockout(username)
	isSet, err = accountLockout.isFailedLoginCountSet()
	if err != nil || isSet != true {
		t.Error("expect:", true, "result:", isSet, err)
	}
	orm.TomatoDBController.DeleteEverything()
}
