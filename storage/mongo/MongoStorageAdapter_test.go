package mongo

import (
	"reflect"
	"testing"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/okobsamoht/talisman/errs"
	"github.com/okobsamoht/talisman/types"
	"github.com/okobsamoht/talisman/utils"
)

func Test_ClassExists(t *testing.T) {
	adapter := getAdapter()
	var name string
	var exist bool
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user2").insertOne(types.M{"_id": "01"})
	name = "user"
	exist = adapter.ClassExists(name)
	if exist == false {
		t.Error("expect:", true, "result:", exist)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user2").insertOne(types.M{"_id": "01"})
	name = "user3"
	exist = adapter.ClassExists(name)
	if exist == true {
		t.Error("expect:", false, "result:", exist)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user2").insertOne(types.M{"_id": "01"})
	name = "user3"
	exist = adapter.ClassExists(name)
	if exist == true {
		t.Error("expect:", false, "result:", exist)
	}
	adapter.adaptiveCollection("user3").insertOne(types.M{"_id": "01"})
	name = "user3"
	exist = adapter.ClassExists(name)
	if exist == false {
		t.Error("expect:", true, "result:", exist)
	}
	adapter.DeleteAllClasses()
}

func Test_SetClassLevelPermissions(t *testing.T) {
	adapter := getAdapter()
	var className string
	var clps types.M
	var err error
	var result []types.M
	var expect types.M
	/*****************************************************/
	className = "user"
	clps = nil
	adapter.CreateClass(className, nil)
	err = adapter.SetClassLevelPermissions(className, clps)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	result, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
		"_metadata": types.M{
			"class_permissions": types.M{},
		},
	}
	if err != nil || result == nil || reflect.DeepEqual(expect, result[0]) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	clps = types.M{
		"find":   types.M{"*": true},
		"get":    types.M{"*": true},
		"create": types.M{"*": true},
		"update": types.M{"*": true},
	}
	adapter.CreateClass(className, nil)
	err = adapter.SetClassLevelPermissions(className, clps)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	result, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
		"_metadata": types.M{
			"class_permissions": types.M{
				"find":   types.M{"*": true},
				"get":    types.M{"*": true},
				"create": types.M{"*": true},
				"update": types.M{"*": true},
			},
		},
	}
	if err != nil || result == nil || reflect.DeepEqual(expect, result[0]) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
}

func Test_CreateClass(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var result types.M
	var err error
	var expect types.M
	var results []types.M
	/*****************************************************/
	className = "user"
	schema = nil
	result, err = adapter.CreateClass(className, schema)
	expect = types.M{
		"className": className,
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"classLevelPermissions": types.M{
			"find":     types.M{"*": true},
			"get":      types.M{"*": true},
			"create":   types.M{"*": true},
			"update":   types.M{"*": true},
			"delete":   types.M{"*": true},
			"addField": types.M{"*": true},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	} else {
		results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
		expect = types.M{
			"_id":       className,
			"objectId":  "string",
			"updatedAt": "string",
			"createdAt": "string",
		}
		if results == nil || len(results) != 1 {
			t.Error("expect:", expect, "result:", results, err)
		} else {
			if reflect.DeepEqual(expect, results[0]) == false {
				t.Error("expect:", expect, "result:", results[0], err)
			}
		}
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{},
	}
	result, err = adapter.CreateClass(className, schema)
	expect = types.M{
		"className": className,
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"classLevelPermissions": types.M{
			"find":     types.M{"*": true},
			"get":      types.M{"*": true},
			"create":   types.M{"*": true},
			"update":   types.M{"*": true},
			"delete":   types.M{"*": true},
			"addField": types.M{"*": true},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	} else {
		results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
		expect = types.M{
			"_id":       className,
			"objectId":  "string",
			"updatedAt": "string",
			"createdAt": "string",
		}
		if results == nil || len(results) != 1 {
			t.Error("expect:", expect, "result:", results, err)
		} else {
			if reflect.DeepEqual(expect, results[0]) == false {
				t.Error("expect:", expect, "result:", results[0], err)
			}
		}
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"k1": types.M{
				"type":        "Pointer",
				"targetClass": "user",
			},
			"k2": types.M{
				"type":        "Relation",
				"targetClass": "user",
			},
			"k3": types.M{
				"type": "String",
			},
		},
	}
	result, err = adapter.CreateClass(className, schema)
	expect = types.M{
		"className": className,
		"fields": types.M{
			"k1": types.M{
				"type":        "Pointer",
				"targetClass": "user",
			},
			"k2": types.M{
				"type":        "Relation",
				"targetClass": "user",
			},
			"k3":        types.M{"type": "String"},
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"classLevelPermissions": types.M{
			"find":     types.M{"*": true},
			"get":      types.M{"*": true},
			"create":   types.M{"*": true},
			"update":   types.M{"*": true},
			"delete":   types.M{"*": true},
			"addField": types.M{"*": true},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	} else {
		results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
		expect = types.M{
			"_id":       className,
			"k1":        "*user",
			"k2":        "relation<user>",
			"k3":        "string",
			"objectId":  "string",
			"updatedAt": "string",
			"createdAt": "string",
		}
		if results == nil || len(results) != 1 {
			t.Error("expect:", expect, "result:", results, err)
		} else {
			if reflect.DeepEqual(expect, results[0]) == false {
				t.Error("expect:", expect, "result:", results[0], err)
			}
		}
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"k1": types.M{
				"type":        "Pointer",
				"targetClass": "user",
			},
			"k2": types.M{
				"type":        "Relation",
				"targetClass": "user",
			},
			"k3": types.M{
				"type": "String",
			},
		},
		"classLevelPermissions": types.M{
			"find":   types.M{"*": true},
			"get":    types.M{"*": true},
			"create": types.M{"*": true},
			"update": types.M{"*": true},
		},
	}
	result, err = adapter.CreateClass(className, schema)
	expect = types.M{
		"className": className,
		"fields": types.M{
			"k1": types.M{
				"type":        "Pointer",
				"targetClass": "user",
			},
			"k2": types.M{
				"type":        "Relation",
				"targetClass": "user",
			},
			"k3":        types.M{"type": "String"},
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"classLevelPermissions": types.M{
			"find":     types.M{"*": true},
			"get":      types.M{"*": true},
			"create":   types.M{"*": true},
			"update":   types.M{"*": true},
			"delete":   types.M{"*": true},
			"addField": types.M{"*": true},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	} else {
		results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
		expect = types.M{
			"_id":       className,
			"k1":        "*user",
			"k2":        "relation<user>",
			"k3":        "string",
			"objectId":  "string",
			"updatedAt": "string",
			"createdAt": "string",
			"_metadata": types.M{
				"class_permissions": types.M{
					"find":   types.M{"*": true},
					"get":    types.M{"*": true},
					"create": types.M{"*": true},
					"update": types.M{"*": true},
				},
			},
		}
		if results == nil || len(results) != 1 {
			t.Error("expect:", expect, "result:", results, err)
		} else {
			if reflect.DeepEqual(expect, results[0]) == false {
				t.Error("expect:", expect, "result:", results[0], err)
			}
		}
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{},
	}
	result, err = adapter.CreateClass(className, schema)
	result, err = adapter.CreateClass(className, schema)
	expectErr := errs.E(errs.DuplicateValue, "Class already exists.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	adapter.DeleteAllClasses()
}

func Test_AddFieldIfNotExists(t *testing.T) {
	// 测试用例与 MongoSchemaCollection.addFieldIfNotExists 相同
}

func Test_DeleteClass(t *testing.T) {
	adapter := getAdapter()
	var className string
	var result types.M
	var err error
	var expect types.M
	/*****************************************************/
	className = "user"
	result, err = adapter.DeleteClass(className)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/*****************************************************/
	className = "user"
	adapter.adaptiveCollection(className).insertOne(types.M{"_id": "1024"})
	adapter.CreateClass(className, nil)
	result, err = adapter.DeleteClass(className)
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	} else {
		results, err := adapter.rawFind(className, types.M{})
		if results != nil && len(results) > 0 {
			t.Error("expect:", 0, "result:", results, err)
		}
		results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
		if results != nil && len(results) > 0 {
			t.Error("expect:", 0, "result:", results, err)
		}
	}
	adapter.DeleteAllClasses()
}

func Test_DeleteAllClasses(t *testing.T) {
	adapter := getAdapter()
	var err error
	var names []string
	/*****************************************************/
	err = adapter.DeleteAllClasses()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	names = adapter.getCollectionNames()
	if names != nil && len(names) != 0 {
		t.Error("expect:", 0, "result:", names)
	}
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user1").insertOne(types.M{"_id": "01"})
	err = adapter.DeleteAllClasses()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	names = adapter.getCollectionNames()
	if names != nil && len(names) != 0 {
		t.Error("expect:", 0, "result:", names)
	}
}

func Test_DeleteFields(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var fieldNames []string
	var err error
	var results []types.M
	var expect types.M
	/*****************************************************/
	className = "user"
	schema = nil
	fieldNames = nil
	err = adapter.DeleteFields(className, schema, fieldNames)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	/*****************************************************/
	className = "user"
	schema = nil
	fieldNames = []string{}
	err = adapter.DeleteFields(className, schema, fieldNames)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	fieldNames = []string{"key"}
	adapter.adaptiveCollection(className).insertOne(types.M{"_id": "01", "key": "hello"})
	adapter.CreateClass(className, schema)
	err = adapter.DeleteFields(className, schema, fieldNames)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = adapter.rawFind(className, types.M{"_id": "01"})
	expect = types.M{
		"_id": "01",
	}
	if err != nil || results == nil || len(results) != 1 || reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if err != nil || results == nil || len(results) != 1 || reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "Pointer"},
		},
	}
	fieldNames = []string{"key"}
	adapter.adaptiveCollection(className).insertOne(types.M{"_id": "01", "_p_key": "hello"})
	adapter.CreateClass(className, schema)
	err = adapter.DeleteFields(className, schema, fieldNames)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = adapter.rawFind(className, types.M{"_id": "01"})
	expect = types.M{
		"_id": "01",
	}
	if err != nil || results == nil || len(results) != 1 || reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if err != nil || results == nil || len(results) != 1 || reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	adapter.DeleteAllClasses()
}

func Test_CreateObject(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var object types.M
	var err error
	var result []types.M
	var expect types.M
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	tmpTime, _ := utils.StringtoTime(tmpTimeStr)
	/*****************************************************/
	className = "user"
	schema = nil
	object = types.M{
		"objectId":  "1024",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
	}
	err = adapter.CreateObject(className, schema, object)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	result, err = adapter.rawFind(className, types.M{"_id": "1024"})
	expect = types.M{
		"_id":         "1024",
		"_updated_at": tmpTime.Local(),
		"_created_at": tmpTime.Local(),
	}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(expect, result[0]) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = nil
	object = types.M{
		"objectId":  "1024",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key": types.M{
			"__type":    "Pointer",
			"className": "abc",
			"objectId":  "123",
		},
	}
	err = adapter.CreateObject(className, schema, object)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	result, err = adapter.rawFind(className, types.M{"_id": "1024"})
	expect = types.M{
		"_id":         "1024",
		"_updated_at": tmpTime.Local(),
		"_created_at": tmpTime.Local(),
		"_p_key":      "abc$123",
	}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(expect, result[0]) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{
				"type": "Pointer",
			},
		},
	}
	object = types.M{
		"objectId":  "1024",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key": types.M{
			"__type":    "Other",
			"className": "abc",
			"objectId":  "123",
		},
	}
	err = adapter.CreateObject(className, schema, object)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	result, err = adapter.rawFind(className, types.M{"_id": "1024"})
	expect = types.M{
		"_id":         "1024",
		"_updated_at": tmpTime.Local(),
		"_created_at": tmpTime.Local(),
		"_p_key": types.M{
			"__type":    "Other",
			"className": "abc",
			"objectId":  "123",
		},
	}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(expect, result[0]) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = nil
	object = types.M{
		"objectId":  "1024",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
	}
	err = adapter.CreateObject(className, schema, object)
	err = adapter.CreateObject(className, schema, object)
	expectErr := errs.E(errs.DuplicateValue, "A duplicate value for a field with unique values was provided")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	result, err = adapter.rawFind(className, types.M{"_id": "1024"})
	expect = types.M{
		"_id":         "1024",
		"_updated_at": tmpTime.Local(),
		"_created_at": tmpTime.Local(),
	}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(expect, result[0]) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
}

func Test_GetClass(t *testing.T) {
	adapter := getAdapter()
	var className string
	var result types.M
	var err error
	var expect types.M
	/*****************************************************/
	className = "user"
	result, err = adapter.GetClass(className)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	adapter.CreateClass(className, nil)
	result, err = adapter.GetClass(className)
	expect = types.M{
		"className": className,
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"classLevelPermissions": types.M{
			"find":     types.M{"*": true},
			"get":      types.M{"*": true},
			"create":   types.M{"*": true},
			"update":   types.M{"*": true},
			"delete":   types.M{"*": true},
			"addField": types.M{"*": true},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
}

func Test_GetAllClasses(t *testing.T) {
	adapter := getAdapter()
	var className string
	var result []types.M
	var err error
	var expect []types.M
	/*****************************************************/
	result, err = adapter.GetAllClasses()
	expect = []types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	adapter.CreateClass(className, nil)
	className = "user1"
	adapter.CreateClass(className, nil)
	result, err = adapter.GetAllClasses()
	expect = []types.M{
		types.M{
			"className": "user",
			"fields": types.M{
				"objectId":  types.M{"type": "String"},
				"updatedAt": types.M{"type": "Date"},
				"createdAt": types.M{"type": "Date"},
				"ACL":       types.M{"type": "ACL"},
			},
			"classLevelPermissions": types.M{
				"find":     types.M{"*": true},
				"get":      types.M{"*": true},
				"create":   types.M{"*": true},
				"update":   types.M{"*": true},
				"delete":   types.M{"*": true},
				"addField": types.M{"*": true},
			},
		},
		types.M{
			"className": "user1",
			"fields": types.M{
				"objectId":  types.M{"type": "String"},
				"updatedAt": types.M{"type": "Date"},
				"createdAt": types.M{"type": "Date"},
				"ACL":       types.M{"type": "ACL"},
			},
			"classLevelPermissions": types.M{
				"find":     types.M{"*": true},
				"get":      types.M{"*": true},
				"create":   types.M{"*": true},
				"update":   types.M{"*": true},
				"delete":   types.M{"*": true},
				"addField": types.M{"*": true},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
}

func Test_getCollectionNames(t *testing.T) {
	adapter := getAdapter()
	var names []string
	/*****************************************************/
	names = adapter.getCollectionNames()
	if names != nil && len(names) > 0 {
		t.Error("expect:", 0, "result:", len(names))
	}
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user1").insertOne(types.M{"_id": "01"})
	names = adapter.getCollectionNames()
	if names == nil || len(names) != 2 {
		t.Error("expect:", 2, "result:", len(names))
	} else {
		expect := []string{"talismanuser", "talismanuser1"}
		if reflect.DeepEqual(expect, names) == false {
			t.Error("expect:", expect, "result:", names)
		}
	}
	adapter.adaptiveCollection("user").drop()
	adapter.adaptiveCollection("user1").drop()
}

func Test_DeleteObjectsByQuery(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var query types.M
	var err error
	var object types.M
	var results []types.M
	var expect interface{}
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	/*****************************************************/
	className = "user"
	schema = nil
	object = types.M{
		"objectId":  "01",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       3,
	}
	adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "02",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       1,
	}
	adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "03",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       1,
	}
	adapter.CreateObject(className, schema, object)
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{
		"key": 1,
	}
	err = adapter.DeleteObjectsByQuery(className, schema, query)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = adapter.rawFind(className, types.M{"key": 1})
	expect = []types.M{}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{
		"key": 1,
	}
	err = adapter.DeleteObjectsByQuery(className, schema, query)
	expect = errs.E(errs.ObjectNotFound, "Object not found.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	results, err = adapter.rawFind(className, types.M{"key": 1})
	expect = []types.M{}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}

	adapter.DeleteAllClasses()
}

func Test_UpdateObjectsByQuery(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var query types.M
	var update types.M
	var err error
	var object types.M
	var results []types.M
	var expect []types.M
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	tmpTime, _ := utils.StringtoTime(tmpTimeStr)
	/*****************************************************/
	className = "user"
	schema = nil
	object = types.M{
		"objectId":  "01",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       3,
	}
	adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "02",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       1,
	}
	adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "03",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       1,
	}
	adapter.CreateObject(className, schema, object)
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{
		"key": 1,
	}
	update = types.M{
		"key": 30,
	}
	err = adapter.UpdateObjectsByQuery(className, schema, query, update)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = adapter.rawFind(className, types.M{"key": 30})
	expect = []types.M{
		types.M{
			"_id":         "02",
			"_updated_at": tmpTime.Local(),
			"_created_at": tmpTime.Local(),
			"key":         30,
		},
		types.M{
			"_id":         "03",
			"_updated_at": tmpTime.Local(),
			"_created_at": tmpTime.Local(),
			"key":         30,
		},
	}
	if err != nil || results == nil || len(results) == 0 || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{
		"objectId": "04",
	}
	update = types.M{
		"key": 30,
	}
	err = adapter.UpdateObjectsByQuery(className, schema, query, update)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = adapter.rawFind(className, types.M{"_id": "04"})
	expect = []types.M{}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}

	adapter.DeleteAllClasses()
}

func Test_FindOneAndUpdate(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var query types.M
	var update types.M
	var result types.M
	var err error
	var object types.M
	var results []types.M
	var expect interface{}
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	tmpTime, _ := utils.StringtoTime(tmpTimeStr)
	/*****************************************************/
	className = "user"
	schema = nil
	object = types.M{
		"objectId":  "01",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       3,
	}
	adapter.CreateObject(className, schema, object)
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{
		"objectId": "01",
	}
	update = types.M{
		"key": 30,
	}
	result, err = adapter.FindOneAndUpdate(className, schema, query, update)
	expect = types.M{
		"objectId":  "01",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       30,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = adapter.rawFind(className, types.M{"_id": "01"})
	expect = []types.M{
		types.M{
			"_id":         "01",
			"_updated_at": tmpTime.Local(),
			"_created_at": tmpTime.Local(),
			"key":         30,
		},
	}
	if err != nil || results == nil || len(results) == 0 || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{
		"objectId": "02",
	}
	update = types.M{
		"key": 30,
	}
	result, err = adapter.FindOneAndUpdate(className, schema, query, update)
	if err != nil || result == nil {
		t.Error("expect:", types.M{}, "result:", result, err)
	}

	adapter.DeleteAllClasses()
}

func Test_UpsertOneObject(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var query types.M
	var update types.M
	var err error
	var object types.M
	var results []types.M
	var expect []types.M
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	tmpTime, _ := utils.StringtoTime(tmpTimeStr)
	/*****************************************************/
	className = "user"
	schema = nil
	object = types.M{
		"objectId":  "01",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       3,
	}
	adapter.CreateObject(className, schema, object)
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{
		"objectId": "01",
	}
	update = types.M{
		"key": 30,
	}
	err = adapter.UpsertOneObject(className, schema, query, update)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = adapter.rawFind(className, types.M{"_id": "01"})
	expect = []types.M{
		types.M{
			"_id":         "01",
			"_updated_at": tmpTime.Local(),
			"_created_at": tmpTime.Local(),
			"key":         30,
		},
	}
	if err != nil || results == nil || len(results) == 0 || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{
		"objectId": "02",
	}
	update = types.M{
		"key": 30,
	}
	err = adapter.UpsertOneObject(className, schema, query, update)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = adapter.rawFind(className, types.M{"_id": "02"})
	expect = []types.M{
		types.M{
			"_id": "02",
			"key": 30,
		},
	}
	if err != nil || results == nil || len(results) == 0 || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}

	adapter.DeleteAllClasses()
}

func Test_Find(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var query types.M
	var options types.M
	var results []types.M
	var err error
	var object types.M
	var expect []types.M
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	/*****************************************************/
	className = "user"
	schema = nil
	object = types.M{
		"objectId":  "01",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       3,
	}
	adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "02",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       1,
	}
	adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "03",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       2,
	}
	adapter.CreateObject(className, schema, object)
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{}
	options = nil
	results, err = adapter.Find(className, schema, query, options)
	expect = []types.M{
		types.M{
			"objectId":  "01",
			"updatedAt": tmpTimeStr,
			"createdAt": tmpTimeStr,
			"key":       3,
		},
		types.M{
			"objectId":  "02",
			"updatedAt": tmpTimeStr,
			"createdAt": tmpTimeStr,
			"key":       1,
		},
		types.M{
			"objectId":  "03",
			"updatedAt": tmpTimeStr,
			"createdAt": tmpTimeStr,
			"key":       2,
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{}
	options = types.M{
		"sort": []string{"key"},
	}
	results, err = adapter.Find(className, schema, query, options)
	expect = []types.M{
		types.M{
			"objectId":  "02",
			"updatedAt": tmpTimeStr,
			"createdAt": tmpTimeStr,
			"key":       1,
		},
		types.M{
			"objectId":  "03",
			"updatedAt": tmpTimeStr,
			"createdAt": tmpTimeStr,
			"key":       2,
		},
		types.M{
			"objectId":  "01",
			"updatedAt": tmpTimeStr,
			"createdAt": tmpTimeStr,
			"key":       3,
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{}
	options = types.M{
		"sort": []string{"-key"},
	}
	results, err = adapter.Find(className, schema, query, options)
	expect = []types.M{
		types.M{
			"objectId":  "01",
			"updatedAt": tmpTimeStr,
			"createdAt": tmpTimeStr,
			"key":       3,
		},
		types.M{
			"objectId":  "03",
			"updatedAt": tmpTimeStr,
			"createdAt": tmpTimeStr,
			"key":       2,
		},
		types.M{
			"objectId":  "02",
			"updatedAt": tmpTimeStr,
			"createdAt": tmpTimeStr,
			"key":       1,
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{}
	options = types.M{
		"keys": []string{"key"},
	}
	results, err = adapter.Find(className, schema, query, options)
	expect = []types.M{
		types.M{
			"objectId": "01",
			"key":      3,
		},
		types.M{
			"objectId": "02",
			"key":      1,
		},
		types.M{
			"objectId": "03",
			"key":      2,
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{"key": 3}
	options = types.M{}
	results, err = adapter.Find(className, schema, query, options)
	expect = []types.M{
		types.M{
			"objectId":  "01",
			"updatedAt": tmpTimeStr,
			"createdAt": tmpTimeStr,
			"key":       3,
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}

	adapter.DeleteAllClasses()
}

func Test_AdapterRawFind(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var query types.M
	var results []types.M
	var err error
	var object types.M
	var expect []types.M
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	tmpTime, _ := utils.StringtoTime(tmpTimeStr)
	/*****************************************************/
	className = "user"
	schema = nil
	object = types.M{
		"objectId":  "01",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       3,
	}
	adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "02",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       1,
	}
	adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "03",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       2,
	}
	adapter.CreateObject(className, schema, object)
	/*****************************************************/
	className = "user"
	query = types.M{}
	results, err = adapter.rawFind(className, query)
	expect = []types.M{
		types.M{
			"_id":         "01",
			"_updated_at": tmpTime.Local(),
			"_created_at": tmpTime.Local(),
			"key":         3,
		},
		types.M{
			"_id":         "02",
			"_updated_at": tmpTime.Local(),
			"_created_at": tmpTime.Local(),
			"key":         1,
		},
		types.M{
			"_id":         "03",
			"_updated_at": tmpTime.Local(),
			"_created_at": tmpTime.Local(),
			"key":         2,
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	adapter.DeleteAllClasses()
}

func Test_Count(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var query types.M
	var count int
	var err error
	var object types.M
	var expect int
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	/*****************************************************/
	className = "user"
	schema = nil
	object = types.M{
		"objectId":  "01",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       3,
	}
	adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "02",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       1,
	}
	adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId":  "03",
		"updatedAt": tmpTimeStr,
		"createdAt": tmpTimeStr,
		"key":       2,
	}
	adapter.CreateObject(className, schema, object)
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{}
	count, err = adapter.Count(className, schema, query)
	expect = 3
	if err != nil || count != expect {
		t.Error("expect:", expect, "result:", count, err)
	}
	/*****************************************************/
	className = "user1"
	schema = nil
	query = types.M{}
	count, err = adapter.Count(className, schema, query)
	expect = 0
	if err != nil || count != expect {
		t.Error("expect:", expect, "result:", count, err)
	}
	/*****************************************************/
	className = "user"
	schema = nil
	query = types.M{"key": 3}
	count, err = adapter.Count(className, schema, query)
	expect = 1
	if err != nil || count != expect {
		t.Error("expect:", expect, "result:", count, err)
	}

	adapter.DeleteAllClasses()
}

func Test_EnsureUniqueness(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var fieldNames []string
	var err error
	var indexes []mgo.Index
	var expect mgo.Index
	var ok bool
	/*****************************************************/
	className = "user"
	schema = nil
	fieldNames = []string{"username"}
	err = adapter.EnsureUniqueness(className, schema, fieldNames)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	indexes, err = adapter.adaptiveCollection(className).collection.Indexes()
	expect = mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}
	ok = false
	for _, i := range indexes {
		if reflect.DeepEqual(i.Key, expect.Key) &&
			i.Unique == expect.Unique &&
			i.Background == expect.Background &&
			i.Sparse == expect.Sparse {
			ok = true
			break
		}
	}
	if ok == false {
		t.Error("expect:", expect, "get result:", indexes)
	}

	adapter.DeleteAllClasses()
}

func Test_storageAdapterAllCollections(t *testing.T) {
	adapter := getAdapter()
	var result []*MongoCollection
	var expect []*MongoCollection
	/*****************************************************/
	result = storageAdapterAllCollections(adapter)
	expect = []*MongoCollection{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user1").insertOne(types.M{"_id": "01"})
	result = storageAdapterAllCollections(adapter)
	if result == nil || len(result) != 2 {
		t.Error("expect:", 2, "result:", len(result))
	} else {
		expect := []string{"talismanuser", "talismanuser1"}
		names := []string{result[0].collection.Name, result[1].collection.Name}
		if reflect.DeepEqual(expect, names) == false {
			t.Error("expect:", expect, "result:", names)
		}
	}
	adapter.adaptiveCollection("user").drop()
	adapter.adaptiveCollection("user1").drop()
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user1").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user.system.id").insertOne(types.M{"_id": "01"})
	result = storageAdapterAllCollections(adapter)
	if result == nil || len(result) != 2 {
		t.Error("expect:", 2, "result:", len(result))
	} else {
		expect := []string{"talismanuser", "talismanuser1"}
		names := []string{result[0].collection.Name, result[1].collection.Name}
		if reflect.DeepEqual(expect, names) == false {
			t.Error("expect:", expect, "result:", names)
		}
	}
	adapter.adaptiveCollection("user").drop()
	adapter.adaptiveCollection("user1").drop()
	adapter.adaptiveCollection("ser.system.id").drop()
}

func Test_convertParseSchemaToMongoSchema(t *testing.T) {
	var schema types.M
	var result types.M
	var expect types.M
	/*****************************************************/
	schema = nil
	result = convertParseSchemaToMongoSchema(schema)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	schema = types.M{}
	result = convertParseSchemaToMongoSchema(schema)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	schema = types.M{
		"className": "user",
		"fields": types.M{
			"name":             types.M{"type": "string"},
			"_rperm":           types.M{"type": "array"},
			"_wperm":           types.M{"type": "array"},
			"_hashed_password": types.M{"type": "string"},
		},
	}
	result = convertParseSchemaToMongoSchema(schema)
	expect = types.M{
		"className": "user",
		"fields": types.M{
			"name":             types.M{"type": "string"},
			"_hashed_password": types.M{"type": "string"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	schema = types.M{
		"className": "_User",
		"fields": types.M{
			"name":             types.M{"type": "string"},
			"_rperm":           types.M{"type": "array"},
			"_wperm":           types.M{"type": "array"},
			"_hashed_password": types.M{"type": "string"},
		},
	}
	result = convertParseSchemaToMongoSchema(schema)
	expect = types.M{
		"className": "_User",
		"fields": types.M{
			"name": types.M{"type": "string"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_mongoSchemaFromFieldsAndClassNameAndCLP(t *testing.T) {
	var fields types.M
	var className string
	var classLevelPermissions types.M
	var result types.M
	var expect types.M
	/*****************************************************/
	fields = nil
	className = "user"
	classLevelPermissions = nil
	result = mongoSchemaFromFieldsAndClassNameAndCLP(fields, className, classLevelPermissions)
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fields = types.M{
		"key1": types.M{
			"type":        "Pointer",
			"targetClass": "_User",
		},
		"key2": types.M{
			"type":        "Relation",
			"targetClass": "_User",
		},
		"loc": types.M{
			"type": "GeoPoint",
		},
	}
	className = "user"
	classLevelPermissions = nil
	result = mongoSchemaFromFieldsAndClassNameAndCLP(fields, className, classLevelPermissions)
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
		"key1":      "*_User",
		"key2":      "relation<_User>",
		"loc":       "geopoint",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fields = types.M{
		"key1": types.M{
			"type":        "Pointer",
			"targetClass": "_User",
		},
		"key2": types.M{
			"type":        "Relation",
			"targetClass": "_User",
		},
		"loc": types.M{
			"type": "GeoPoint",
		},
	}
	className = "user"
	classLevelPermissions = types.M{
		"find":     types.M{"*": true},
		"get":      types.M{"*": true},
		"create":   types.M{"*": true},
		"update":   types.M{"*": true},
		"delete":   types.M{"*": true},
		"addField": types.M{"*": true},
	}
	result = mongoSchemaFromFieldsAndClassNameAndCLP(fields, className, classLevelPermissions)
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
		"key1":      "*_User",
		"key2":      "relation<_User>",
		"loc":       "geopoint",
		"_metadata": types.M{
			"class_permissions": types.M{
				"find":     types.M{"*": true},
				"get":      types.M{"*": true},
				"create":   types.M{"*": true},
				"update":   types.M{"*": true},
				"delete":   types.M{"*": true},
				"addField": types.M{"*": true},
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func getAdapter() *MongoAdapter {
	return NewMongoAdapter("talisman", openDB())
}
