package rest

import (
	"reflect"
	"testing"
	"time"

	"github.com/okobsamoht/talisman/cache"
	"github.com/okobsamoht/talisman/errs"
	"github.com/okobsamoht/talisman/orm"
	"github.com/okobsamoht/talisman/types"
	"github.com/okobsamoht/talisman/utils"
)

func Test_GetAuthForSessionToken(t *testing.T) {
	var schema types.M
	var object types.M
	var className string
	var sessionToken string
	var installationID string
	var result *Auth
	var err error
	var expect *Auth
	var expectErr error
	/********************************************************/
	cache.InitCache()
	initEnv()
	sessionToken = "abc"
	installationID = "111"
	_, err = GetAuthForSessionToken(sessionToken, installationID)
	expectErr = errs.E(errs.InvalidSessionToken, "invalid session token")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TalismanDBController.DeleteEverything()
	/********************************************************/
	cache.InitCache()
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"password": "123",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Session"
	schema = types.M{
		"fields": types.M{
			"user":         types.M{"type": "Pointer", "targetClass": "_User"},
			"sessionToken": types.M{"type": "String"},
			"expiresAt":    types.M{"type": "Date"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"user": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "1001",
		},
		"sessionToken": "abc1001",
		"expiresAt":    utils.TimetoString(time.Now().UTC()),
	}
	orm.Adapter.CreateObject(className, schema, object)
	sessionToken = "abc"
	installationID = "111"
	_, err = GetAuthForSessionToken(sessionToken, installationID)
	expectErr = errs.E(errs.InvalidSessionToken, "invalid session token")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TalismanDBController.DeleteEverything()
	/********************************************************/
	cache.InitCache()
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"password": "123",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Session"
	schema = types.M{
		"fields": types.M{
			"user":         types.M{"type": "Pointer", "targetClass": "_User"},
			"sessionToken": types.M{"type": "String"},
			"expiresAt":    types.M{"type": "Date"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"user": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "1001",
		},
		"sessionToken": "abc1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	sessionToken = "abc1001"
	installationID = "111"
	_, err = GetAuthForSessionToken(sessionToken, installationID)
	expectErr = errs.E(errs.InvalidSessionToken, "Session token is expired.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TalismanDBController.DeleteEverything()
	/********************************************************/
	cache.InitCache()
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"password": "123",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Session"
	schema = types.M{
		"fields": types.M{
			"user":         types.M{"type": "Pointer", "targetClass": "_User"},
			"sessionToken": types.M{"type": "String"},
			"expiresAt":    types.M{"type": "Date"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"user": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "1001",
		},
		"sessionToken": "abc1001",
		"expiresAt":    utils.TimetoString(time.Now().UTC()),
	}
	orm.Adapter.CreateObject(className, schema, object)
	sessionToken = "abc1001"
	installationID = "111"
	_, err = GetAuthForSessionToken(sessionToken, installationID)
	expectErr = errs.E(errs.InvalidSessionToken, "Session token is expired.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TalismanDBController.DeleteEverything()
	/********************************************************/
	cache.InitCache()
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"password": "123",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Session"
	schema = types.M{
		"fields": types.M{
			"user":         types.M{"type": "Pointer", "targetClass": "_User"},
			"sessionToken": types.M{"type": "String"},
			"expiresAt":    types.M{"type": "Date"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"user": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "1001",
		},
		"sessionToken": "abc1001",
		"expiresAt":    utils.TimetoString(time.Now().UTC().Add(5 * time.Second)),
	}
	orm.Adapter.CreateObject(className, schema, object)
	sessionToken = "abc1001"
	installationID = "111"
	result, err = GetAuthForSessionToken(sessionToken, installationID)
	expect = &Auth{
		IsMaster:       false,
		InstallationID: "111",
		User: types.M{
			"__type":       "Object",
			"className":    "_User",
			"objectId":     "1001",
			"username":     "joe",
			"sessionToken": "abc1001",
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	orm.TalismanDBController.DeleteEverything()
}

func Test_CouldUpdateUserID(t *testing.T) {
	var auth *Auth
	var result bool
	var expect bool
	/********************************************************/
	auth = &Auth{
		IsMaster: true,
	}
	result = auth.CouldUpdateUserID("1001")
	expect = true
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/********************************************************/
	auth = &Auth{
		IsMaster: false,
	}
	result = auth.CouldUpdateUserID("1001")
	expect = false
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/********************************************************/
	auth = &Auth{
		IsMaster: false,
		User:     types.M{},
	}
	result = auth.CouldUpdateUserID("1001")
	expect = false
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/********************************************************/
	auth = &Auth{
		IsMaster: false,
		User:     types.M{"objectId": "1002"},
	}
	result = auth.CouldUpdateUserID("1001")
	expect = false
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/********************************************************/
	auth = &Auth{
		IsMaster: false,
		User:     types.M{"objectId": "1001"},
	}
	result = auth.CouldUpdateUserID("1001")
	expect = true
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_GetUserRoles(t *testing.T) {
	var schema types.M
	var object types.M
	var className string
	var auth *Auth
	var result []string
	var expect []string
	/********************************************************/
	auth = Master()
	result = auth.GetUserRoles()
	expect = []string{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/********************************************************/
	auth = Nobody()
	result = auth.GetUserRoles()
	expect = []string{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/********************************************************/
	auth = &Auth{
		IsMaster: false,
		User: types.M{
			"objectId": "9001",
		},
		FetchedRoles: true,
		UserRoles:    []string{"role:role1001"},
	}
	result = auth.GetUserRoles()
	expect = []string{"role:role1001"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/********************************************************/
	auth = &Auth{
		IsMaster: false,
		User: types.M{
			"objectId": "9001",
		},
		FetchedRoles: false,
		RolePromise:  []string{"role:role1001"},
	}
	result = auth.GetUserRoles()
	expect = []string{"role:role1001"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/********************************************************/
	cache.InitCache()
	initEnv()
	className = "_Role"
	schema = types.M{
		"fields": types.M{
			"name":  types.M{"type": "String"},
			"users": types.M{"type": "Relation", "targetClass": "_User"},
			"roles": types.M{"type": "Relation", "targetClass": "_Role"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"name":     "role1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"name":     "role1002",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Join:roles:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "5001",
		"owningId":  "1002",
		"relatedId": "1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Join:users:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "5002",
		"owningId":  "1001",
		"relatedId": "9001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	auth = &Auth{
		IsMaster: false,
		User: types.M{
			"objectId": "9001",
		},
		FetchedRoles: false,
		RolePromise:  nil,
	}
	result = auth.GetUserRoles()
	expect = []string{"role:role1001", "role:role1002"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TalismanDBController.DeleteEverything()
}

func Test_loadRoles(t *testing.T) {
	var schema types.M
	var object types.M
	var className string
	var auth *Auth
	var result []string
	var expect []string
	/********************************************************/
	cache.InitCache()
	initEnv()
	className = "_Role"
	schema = types.M{
		"fields": types.M{
			"name":  types.M{"type": "String"},
			"users": types.M{"type": "Relation", "targetClass": "_User"},
			"roles": types.M{"type": "Relation", "targetClass": "_Role"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"name":     "role1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"name":     "role1002",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Join:roles:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "5001",
		"owningId":  "1002",
		"relatedId": "1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	auth = &Auth{
		IsMaster: false,
		User: types.M{
			"objectId": "9001",
		},
	}
	result = auth.loadRoles()
	expect = []string{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TalismanDBController.DeleteEverything()
	/********************************************************/
	cache.InitCache()
	initEnv()
	className = "_Role"
	schema = types.M{
		"fields": types.M{
			"name":  types.M{"type": "String"},
			"users": types.M{"type": "Relation", "targetClass": "_User"},
			"roles": types.M{"type": "Relation", "targetClass": "_Role"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"name":     "role1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"name":     "role1002",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Join:roles:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "5001",
		"owningId":  "1002",
		"relatedId": "1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Join:users:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "5002",
		"owningId":  "1001",
		"relatedId": "9001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	auth = &Auth{
		IsMaster: false,
		User: types.M{
			"objectId": "9001",
		},
	}
	result = auth.loadRoles()
	expect = []string{"role:role1001", "role:role1002"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TalismanDBController.DeleteEverything()
}

func Test_getAllRolesNamesForRoleIds(t *testing.T) {
	var schema types.M
	var object types.M
	var className string
	var roleIDs []string
	var names []string
	var queriedRoles map[string]bool
	var result []string
	var expect []string
	/********************************************************/
	initEnv()
	className = "_Role"
	schema = types.M{
		"fields": types.M{
			"name":  types.M{"type": "String"},
			"users": types.M{"type": "Relation", "targetClass": "_User"},
			"roles": types.M{"type": "Relation", "targetClass": "_Role"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	className = "_Join:roles:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	roleIDs = []string{"1001"}
	names = []string{}
	queriedRoles = map[string]bool{}
	result = Master().getAllRolesNamesForRoleIds(roleIDs, names, queriedRoles)
	expect = []string{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TalismanDBController.DeleteEverything()
	/********************************************************/
	initEnv()
	className = "_Role"
	schema = types.M{
		"fields": types.M{
			"name":  types.M{"type": "String"},
			"users": types.M{"type": "Relation", "targetClass": "_User"},
			"roles": types.M{"type": "Relation", "targetClass": "_Role"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"name":     "role1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Join:roles:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	roleIDs = []string{"1001"}
	names = []string{}
	queriedRoles = map[string]bool{}
	result = Master().getAllRolesNamesForRoleIds(roleIDs, names, queriedRoles)
	expect = []string{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TalismanDBController.DeleteEverything()
	/********************************************************/
	initEnv()
	className = "_Role"
	schema = types.M{
		"fields": types.M{
			"name":  types.M{"type": "String"},
			"users": types.M{"type": "Relation", "targetClass": "_User"},
			"roles": types.M{"type": "Relation", "targetClass": "_Role"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"name":     "role1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"name":     "role1002",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Join:roles:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "5001",
		"owningId":  "1002",
		"relatedId": "1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	roleIDs = []string{"1001"}
	names = []string{}
	queriedRoles = map[string]bool{}
	result = Master().getAllRolesNamesForRoleIds(roleIDs, names, queriedRoles)
	expect = []string{"role1002"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TalismanDBController.DeleteEverything()
	/********************************************************/
	initEnv()
	className = "_Role"
	schema = types.M{
		"fields": types.M{
			"name":  types.M{"type": "String"},
			"users": types.M{"type": "Relation", "targetClass": "_User"},
			"roles": types.M{"type": "Relation", "targetClass": "_Role"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"name":     "role1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"name":     "role1002",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1003",
		"name":     "role1003",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2001",
		"name":     "role2001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"name":     "role2002",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"name":     "role2003",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Join:roles:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "5001",
		"owningId":  "1002",
		"relatedId": "1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "5002",
		"owningId":  "1003",
		"relatedId": "1002",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "5003",
		"owningId":  "2002",
		"relatedId": "2001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "5004",
		"owningId":  "2003",
		"relatedId": "2002",
	}
	orm.Adapter.CreateObject(className, schema, object)
	roleIDs = []string{"1001", "2001"}
	names = []string{}
	queriedRoles = map[string]bool{}
	result = Master().getAllRolesNamesForRoleIds(roleIDs, names, queriedRoles)
	expect = []string{"role1002", "role2002", "role1003", "role2003"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TalismanDBController.DeleteEverything()
}
