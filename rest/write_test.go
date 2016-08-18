package rest

import (
	"reflect"
	"testing"
	"time"

	"github.com/lfq7413/tomato/cache"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func Test_NewWrite(t *testing.T) {
	var auth *Auth
	var className string
	var query types.M
	var data types.M
	var originalData types.M
	var clientSDK map[string]string
	var result *Write
	var err error
	var expect *Write
	var expectErr error
	/***************************************************************/
	auth = nil
	className = "user"
	query = nil
	data = types.M{
		"objectId": "1001",
	}
	originalData = nil
	clientSDK = nil
	_, err = NewWrite(auth, className, query, data, originalData, clientSDK)
	expectErr = errs.E(errs.InvalidKeyName, "objectId is an invalid field name.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	auth = nil
	className = "user"
	query = nil
	data = types.M{
		"key": "hello",
	}
	originalData = nil
	clientSDK = nil
	result, err = NewWrite(auth, className, query, data, originalData, clientSDK)
	expect = &Write{
		auth:                       Nobody(),
		className:                  "user",
		query:                      nil,
		data:                       types.M{"key": "hello"},
		originalData:               nil,
		storage:                    types.M{},
		RunOptions:                 types.M{},
		response:                   nil,
		updatedAt:                  utils.TimetoString(time.Now().UTC()),
		responseShouldHaveUsername: false,
		clientSDK:                  nil,
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/***************************************************************/
	auth = nil
	className = "user"
	query = types.M{
		"objectId": "1001",
	}
	data = types.M{
		"key": "hello",
	}
	originalData = types.M{
		"key": "hi",
	}
	clientSDK = nil
	result, err = NewWrite(auth, className, query, data, originalData, clientSDK)
	expect = &Write{
		auth:                       Nobody(),
		className:                  "user",
		query:                      types.M{"objectId": "1001"},
		data:                       types.M{"key": "hello"},
		originalData:               types.M{"key": "hi"},
		storage:                    types.M{},
		RunOptions:                 types.M{},
		response:                   nil,
		updatedAt:                  utils.TimetoString(time.Now().UTC()),
		responseShouldHaveUsername: false,
		clientSDK:                  nil,
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_Execute_Write(t *testing.T) {
	// TODO
}

func Test_getUserAndRoleACL_Write(t *testing.T) {
	var schema types.M
	var object types.M
	var className string
	var w *Write
	var auth *Auth
	var query types.M
	var data types.M
	var originalData types.M
	var expect []string
	/***************************************************************/
	auth = Master()
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(auth, "user", query, data, originalData, nil)
	w.getUserAndRoleACL()
	if _, ok := w.RunOptions["acl"]; ok {
		t.Error("findOptions[acl] exist")
	}
	/***************************************************************/
	auth = Nobody()
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(auth, "user", query, data, originalData, nil)
	w.getUserAndRoleACL()
	expect = []string{"*"}
	if reflect.DeepEqual(expect, w.RunOptions["acl"]) == false {
		t.Error("expect:", expect, "result:", w.RunOptions["acl"])
	}
	/***************************************************************/
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
	w, _ = NewWrite(auth, "user", query, data, originalData, nil)
	w.getUserAndRoleACL()
	expect = []string{"*", "9001", "role:role1001", "role:role1002"}
	if reflect.DeepEqual(expect, w.RunOptions["acl"]) == false {
		t.Error("expect:", expect, "result:", w.RunOptions["acl"])
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_validateClientClassCreation_Write(t *testing.T) {
	// 测试用例与 query.validateClientClassCreation 相同
}

func Test_validateSchema(t *testing.T) {
	// 测试用例与 DBController.ValidateObject 相同
}

func Test_handleInstallation(t *testing.T) {
	// TODO
}

func Test_handleSession(t *testing.T) {
	// Execute
	// TODO
}

func Test_validateAuthData(t *testing.T) {
	var className string
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var result error
	var expect error
	/***************************************************************/
	initEnv()
	className = "user"
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.validateAuthData()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	query = nil
	data = types.M{
		"key": "hello",
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.validateAuthData()
	expect = errs.E(errs.UsernameMissing, "bad or missing username")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	query = nil
	data = types.M{
		"username": "joe",
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.validateAuthData()
	expect = errs.E(errs.PasswordMissing, "password is required")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	query = nil
	data = types.M{
		"username": "joe",
		"password": "123",
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.validateAuthData()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.validateAuthData()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"key": "1001",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.validateAuthData()
	expect = errs.E(errs.UnsupportedService, "This authentication method is unsupported.")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_runBeforeTrigger(t *testing.T) {
	// TODO
}

func Test_setRequiredFieldsIfNeeded(t *testing.T) {
	// TODO
}

func Test_transformUser(t *testing.T) {
	// TODO
}

func Test_expandFilesForExistingObjects(t *testing.T) {
	// TODO
}

func Test_runDatabaseOperation(t *testing.T) {
	// TODO
}

func Test_createSessionTokenIfNeeded(t *testing.T) {
	// createSessionToken
	// TODO
}

func Test_handleFollowup(t *testing.T) {
	// createSessionToken
	// TODO
}

func Test_runAfterTrigger(t *testing.T) {
	// TODO
}

func Test_cleanUserAuthData(t *testing.T) {
	// TODO
}

/////////////////////////////////////////////////////////////

func Test_handleAuthData(t *testing.T) {
	var className string
	var schema types.M
	var object types.M
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var result error
	var expect error
	var response types.M
	var location string
	/***************************************************************/
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"other": types.M{
				"id": "1001",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = errs.E(errs.UnsupportedService, "This authentication method is unsupported.")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/***************************************************************/
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	object = types.M{
		"objectId": "102",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = errs.E(errs.AccountAlreadyLinked, "this auth is already used")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id": "1002",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	if reflect.DeepEqual("facebook", w.storage["authProvider"]) == false {
		t.Error("expect:", "facebook", "result:", w.storage["authProvider"])
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	response = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	if reflect.DeepEqual(utils.M(response), w.response["response"]) == false {
		t.Error("expect:", response, "result:", w.response["response"])
	}
	location = "http://www.g.cn/users/101"
	if reflect.DeepEqual(location, w.response["location"]) == false {
		t.Error("expect:", location, "result:", w.response["location"])
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "aaa",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "abc",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	response = types.M{
		"objectId": "101",
		"authData": map[string]interface{}{
			"facebook": types.M{
				"id":    "1001",
				"token": "abc",
			},
		},
	}
	if reflect.DeepEqual(utils.M(response), w.response["response"]) == false {
		t.Error("expect:", response, "result:", w.response["response"])
	}
	location = "http://www.g.cn/users/101"
	if reflect.DeepEqual(location, w.response["location"]) == false {
		t.Error("expect:", location, "result:", w.response["location"])
	}
	r, _ := orm.TomatoDBController.Find(className, types.M{"objectId": "101"}, types.M{})
	response = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "abc",
			},
		},
	}
	if reflect.DeepEqual(response, r[0]) == false {
		t.Error("expect:", response, "result:", r[0])
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "aaa",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = types.M{"objectId": "101"}
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "aaa",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "aaa",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = types.M{"objectId": "102"}
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "aaa",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = errs.E(errs.AccountAlreadyLinked, "this auth is already used")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_handleAuthDataValidation(t *testing.T) {
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var authData types.M
	var result error
	var expect error
	/***************************************************************/
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	authData = types.M{
		"facebook": types.M{
			"id": "1001",
		},
	}
	result = w.handleAuthDataValidation(authData)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/***************************************************************/
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	authData = types.M{
		"other": types.M{
			"id": "1001",
		},
	}
	result = w.handleAuthDataValidation(authData)
	expect = errs.E(errs.UnsupportedService, "This authentication method is unsupported.")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_findUsersWithAuthData(t *testing.T) {
	var schema types.M
	var object types.M
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var className string
	var authData types.M
	var result types.S
	var err error
	var expect types.S
	/***************************************************************/
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	query = nil
	data = types.M{}
	originalData = nil
	className = "user"
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	authData = types.M{
		"facebook": types.M{
			"id": "1002",
		},
	}
	result, err = w.findUsersWithAuthData(authData)
	expect = types.S{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	query = nil
	data = types.M{}
	originalData = nil
	className = "_User"
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	authData = types.M{
		"facebook": types.M{
			"id": "1001",
		},
	}
	result, err = w.findUsersWithAuthData(authData)
	expect = types.S{
		types.M{
			"objectId": "101",
			"authData": types.M{
				"facebook": types.M{
					"id": "1001",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	object = types.M{
		"objectId": "102",
		"authData": types.M{
			"twitter": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	query = nil
	data = types.M{}
	originalData = nil
	className = "_User"
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	authData = types.M{
		"facebook": types.M{
			"id": "1001",
		},
		"twitter": types.M{
			"id": "1001",
		},
	}
	result, err = w.findUsersWithAuthData(authData)
	expect = types.S{
		types.M{
			"objectId": "101",
			"authData": types.M{
				"facebook": types.M{
					"id": "1001",
				},
			},
		},
		types.M{
			"objectId": "102",
			"authData": types.M{
				"twitter": types.M{
					"id": "1001",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_createSessionToken(t *testing.T) {
	// Execute
	// TODO
}

func Test_location(t *testing.T) {
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var className string
	var result string
	var expect string
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	query = nil
	data = types.M{}
	originalData = nil
	className = "post"
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	w.data["objectId"] = "1001"
	result = w.location()
	expect = "http://www.g.cn/classes/post/1001"
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	query = nil
	data = types.M{}
	originalData = nil
	className = "_User"
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	w.data["objectId"] = "1001"
	result = w.location()
	expect = "http://www.g.cn/users/1001"
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_objectID(t *testing.T) {
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var result interface{}
	var expect interface{}
	/***************************************************************/
	query = nil
	data = types.M{"key": "hello"}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	w.data["objectId"] = "1001"
	result = w.objectID()
	expect = "1001"
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/***************************************************************/
	query = types.M{"objectId": "1001"}
	data = types.M{"key": "hello"}
	originalData = types.M{"key": "hi"}
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	result = w.objectID()
	expect = "1001"
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_sanitizedData(t *testing.T) {
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var result types.M
	var expect types.M
	/***************************************************************/
	query = nil
	data = types.M{
		"key":              "hello",
		"_auth_data":       "facebook",
		"_hashed_password": "123456",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	result = w.sanitizedData()
	expect = types.M{
		"key": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_updateResponseWithData(t *testing.T) {
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var response, updateData types.M
	var result types.M
	var expect types.M
	/***************************************************************/
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	response = types.M{
		"key":  "hello",
		"key1": 10,
	}
	updateData = types.M{
		"key1": types.M{
			"__op": "Increment",
		},
		"key2": "world",
		"key3": types.M{
			"__op": "Delete",
		},
	}
	result = w.updateResponseWithData(response, updateData)
	expect = types.M{
		"key":  "hello",
		"key1": 10,
		"key2": "world",
		"key3": types.M{
			"__op": "Delete",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}
