package orm

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func Test_CollectionExists(t *testing.T) {
	// TODO
}

func Test_PurgeCollection(t *testing.T) {
	// LoadSchema
	// TODO
}

func Test_Find(t *testing.T) {
	// LoadSchema
	// reduceRelationKeys
	// reduceInRelation
	// addPointerPermissions
	// addReadACL
	// validateQuery
	// filterSensitiveData
	// TODO
}

func Test_Destroy(t *testing.T) {
	// LoadSchema
	// addPointerPermissions
	// addWriteACL
	// validateQuery
	// TODO
}

func Test_Update(t *testing.T) {
	// LoadSchema
	// handleRelationUpdates
	// addPointerPermissions
	// addWriteACL
	// validateQuery
	// sanitizeDatabaseResult
	// TODO
}

func Test_Create(t *testing.T) {
	// validateClassName
	// LoadSchema
	// handleRelationUpdates
	// TODO
}

func Test_validateClassName(t *testing.T) {
	// TODO
}

func Test_handleRelationUpdates(t *testing.T) {
	// addRelation
	// removeRelation
	// TODO
}

func Test_addRelation(t *testing.T) {
	// TODO
}

func Test_removeRelation(t *testing.T) {
	// TODO
}

func Test_ValidateObject(t *testing.T) {
	// LoadSchema
	// canAddField
	// TODO
}

func Test_LoadSchema(t *testing.T) {
	// TODO
}

func Test_DeleteEverything(t *testing.T) {
	// TODO
}

func Test_RedirectClassNameForKey(t *testing.T) {
	// LoadSchema
	// TODO
}

func Test_canAddField(t *testing.T) {
	// TODO
}

func Test_reduceRelationKeys(t *testing.T) {
	// relatedIds
	// addInObjectIdsIds
	// TODO
}

func Test_relatedIds(t *testing.T) {
	// TODO
}

func Test_addInObjectIdsIds(t *testing.T) {
	// TODO
}

func Test_addNotInObjectIdsIds(t *testing.T) {
	// TODO
}

func Test_reduceInRelation(t *testing.T) {
	// owningIds
	// addNotInObjectIdsIds
	// addInObjectIdsIds
	// TODO
}

func Test_owningIds(t *testing.T) {
	// TODO
}

func Test_DeleteSchema(t *testing.T) {
	// LoadSchema
	// CollectionExists
	// TODO
}

func Test_addPointerPermissions(t *testing.T) {
	// TODO
}

//////////////////////////////////////////////////////

func Test_sanitizeDatabaseResult(t *testing.T) {
	// TODO
}

func Test_keysForQuery(t *testing.T) {
	// TODO
}

func Test_joinTableName(t *testing.T) {
	// TODO
}

func Test_filterSensitiveData(t *testing.T) {
	// TODO
}

func Test_addWriteACL(t *testing.T) {
	// TODO
}

func Test_addReadACL(t *testing.T) {
	// TODO
}

func Test_validateQuery(t *testing.T) {
	// TODO
}

func Test_transformObjectACL(t *testing.T) {
	var object types.M
	var result types.M
	var expect types.M
	/*************************************************/
	object = nil
	result = transformObjectACL(object)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	object = types.M{}
	result = transformObjectACL(object)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	object = types.M{"ACL": "hello"}
	result = transformObjectACL(object)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	object = types.M{
		"ACL": types.M{
			"userid": types.M{
				"read":  true,
				"write": true,
			},
			"role:xxx": types.M{
				"read":  true,
				"write": true,
			},
			"*": types.M{
				"read": true,
			},
		},
	}
	result = transformObjectACL(object)
	expect = types.M{
		"_rperm": types.S{"userid", "role:xxx", "*"},
		"_wperm": types.S{"userid", "role:xxx"},
	}
	if utils.CompareArray(expect["_rperm"], result["_rperm"]) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	if utils.CompareArray(expect["_wperm"], result["_wperm"]) == false {
		t.Error("expect:", expect, "get result:", result)
	}
}

func Test_untransformObjectACL(t *testing.T) {
	var output types.M
	var result types.M
	var expect types.M
	/*************************************************/
	output = nil
	result = untransformObjectACL(output)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	output = types.M{}
	result = untransformObjectACL(output)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	output = types.M{"_rperm": "Incorrect type"}
	result = untransformObjectACL(output)
	expect = types.M{"ACL": types.M{}}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	output = types.M{"_rperm": types.S{"userid", "role:xxx"}}
	result = untransformObjectACL(output)
	expect = types.M{
		"ACL": types.M{
			"userid": types.M{
				"read": true,
			},
			"role:xxx": types.M{
				"read": true,
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	output = types.M{"_wperm": "Incorrect type"}
	result = untransformObjectACL(output)
	expect = types.M{"ACL": types.M{}}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	output = types.M{"_wperm": types.S{"userid", "role:xxx"}}
	result = untransformObjectACL(output)
	expect = types.M{
		"ACL": types.M{
			"userid": types.M{
				"write": true,
			},
			"role:xxx": types.M{
				"write": true,
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	output = types.M{
		"_rperm": types.S{"userid", "role:xxx", "*"},
		"_wperm": types.S{"userid", "role:xxx"},
	}
	result = untransformObjectACL(output)
	expect = types.M{
		"ACL": types.M{
			"userid": types.M{
				"read":  true,
				"write": true,
			},
			"role:xxx": types.M{
				"read":  true,
				"write": true,
			},
			"*": types.M{
				"read": true,
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_transformAuthData(t *testing.T) {
	var className string
	var object types.M
	var expect types.M
	var schema types.M
	/*************************************************/
	className = "Other"
	object = types.M{"key": "value"}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{"key": "value"}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = nil
	schema = nil
	transformAuthData(className, object, schema)
	expect = nil
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{"authData": nil}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{"authData": 1024}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{
		"authData": types.M{
			"facebook": nil,
		},
	}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{
		"_auth_data_facebook": types.M{
			"__op": "Delete",
		},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{
		"authData": types.M{
			"facebook": 1024,
		},
	}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{
		"_auth_data_facebook": types.M{
			"__op": "Delete",
		},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{
		"authData": types.M{
			"facebook": types.M{},
		},
	}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{
		"_auth_data_facebook": types.M{
			"__op": "Delete",
		},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{
		"authData": types.M{
			"facebook": types.M{"id": "1024"},
			"twitter":  types.M{},
		},
	}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{
		"_auth_data_facebook": types.M{"id": "1024"},
		"_auth_data_twitter":  types.M{"__op": "Delete"},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{
		"authData": types.M{
			"facebook": types.M{"id": "1024"},
			"twitter":  types.M{},
		},
	}
	schema = types.M{
		"fields": types.M{},
	}
	transformAuthData(className, object, schema)
	expect = types.M{
		"_auth_data_facebook": types.M{"id": "1024"},
		"_auth_data_twitter":  types.M{"__op": "Delete"},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	expect = types.M{
		"fields": types.M{
			"_auth_data_facebook": types.M{"type": "Object"},
		},
	}
	if reflect.DeepEqual(expect, schema) == false {
		t.Error("expect:", expect, "result:", schema)
	}
}

func Test_flattenUpdateOperatorsForCreate(t *testing.T) {
	var object types.M
	var err error
	var expect interface{}
	/**********************************************************/
	object = nil
	err = flattenUpdateOperatorsForCreate(object)
	expect = nil
	if err != nil || object != nil {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{"key": "value"}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{"key": "value"}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{"key": "value"},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": types.M{"key": "value"},
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":   "Increment",
			"amount": 10.24,
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": 10.24,
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":   "Increment",
			"amount": 1024,
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": 1024,
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":   "Increment",
			"amount": "hello",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an number")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op": "Increment",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an number")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "Add",
			"objects": types.S{"abc", "def"},
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": types.S{"abc", "def"},
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "Add",
			"objects": "hello",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an array")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "AddUnique",
			"objects": types.S{"abc", "def"},
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": types.S{"abc", "def"},
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "AddUnique",
			"objects": "hello",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an array")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "Remove",
			"objects": types.S{"abc", "def"},
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": types.S{},
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "Remove",
			"objects": "hello",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an array")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op": "Delete",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op": "Other",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.CommandUnavailable, "The Other operator is not supported yet.")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
}
