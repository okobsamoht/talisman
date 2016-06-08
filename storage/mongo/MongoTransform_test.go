package mongo

import (
	"reflect"
	"testing"
	"time"

	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func Test_transformKey(t *testing.T) {
	tf := NewTransform()
	var fieldName string
	var schema types.M
	var result string
	/*************************************************/
	fieldName = "objectId"
	result = tf.transformKey("", fieldName, schema)
	if result != "_id" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	fieldName = "createdAt"
	result = tf.transformKey("", fieldName, schema)
	if result != "_created_at" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	fieldName = "updatedAt"
	result = tf.transformKey("", fieldName, schema)
	if result != "_updated_at" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	fieldName = "sessionToken"
	result = tf.transformKey("", fieldName, schema)
	if result != "_session_token" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	schema = nil
	fieldName = "user"
	result = tf.transformKey("", fieldName, schema)
	if result != "user" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	schema = types.M{
		"fields": "type is string",
	}
	fieldName = "user"
	result = tf.transformKey("", fieldName, schema)
	if result != "user" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	schema = types.M{
		"fields": types.M{
			"user": "type is string",
		},
	}
	fieldName = "user"
	result = tf.transformKey("", fieldName, schema)
	if result != "user" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	schema = types.M{
		"fields": types.M{
			"user": types.M{
				"__type": "String",
			},
		},
	}
	fieldName = "user"
	result = tf.transformKey("", fieldName, schema)
	if result != "user" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	schema = types.M{
		"fields": types.M{
			"user": types.M{
				"__type": "Pointer",
			},
		},
	}
	fieldName = "user"
	result = tf.transformKey("", fieldName, schema)
	if result != "_p_user" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
}

func Test_transformKeyValueForUpdate(t *testing.T) {
	// transformTopLevelAtom
	// transformInteriorValue
	// transformUpdateOperator
	// TODO
}

func Test_transformQueryKeyValue(t *testing.T) {
	// transformWhere
	// transformConstraint
	// transformTopLevelAtom
	// TODO
}

func Test_transformConstraint(t *testing.T) {
	// transformInteriorAtom
	// transformTopLevelAtom
	// TODO
}

func Test_transformTopLevelAtom(t *testing.T) {
	// dateCoder
	// bytesCoder
	// geoPointCoder
	// TODO
}

func Test_transformUpdateOperator(t *testing.T) {
	// transformInteriorAtom
	// TODO
}

func Test_parseObjectToMongoObjectForCreate(t *testing.T) {
	// transformAuthData
	// transformACL
	// parseObjectKeyValueToMongoObjectKeyValue
	// TODO
}

func Test_parseObjectKeyValueToMongoObjectKeyValue(t *testing.T) {
	// transformTopLevelAtom
	// transformInteriorValue
	// transformUpdateOperator
	// TODO
}

func Test_transformAuthData(t *testing.T) {
	// TODO
}

func Test_transformACL(t *testing.T) {
	// TODO
}

func Test_transformWhere(t *testing.T) {
	// transformQueryKeyValue
	// TODO
}

func Test_transformUpdate(t *testing.T) {
	// transformAuthData
	// transformACL
	// transformKeyValueForUpdate
	// TODO
}

func Test_nestedMongoObjectToNestedParseObject(t *testing.T) {
	// TODO
}

func Test_mongoObjectToParseObject(t *testing.T) {
	// nestedMongoObjectToNestedParseObject
	// dateCoder
	// bytesCoder
	// untransformACL
	// TODO
}

func Test_untransformACL(t *testing.T) {
	// TODO
}

func Test_transformInteriorAtom(t *testing.T) {
	// dateCoder
	// bytesCoder
	// TODO
}

func Test_transformInteriorValue(t *testing.T) {
	// transformInteriorAtom
	// TODO
}

func Test_dateCoder(t *testing.T) {
	// TODO
}

func Test_bytesCoder(t *testing.T) {
	// TODO
}

func Test_geoPointCoder(t *testing.T) {
	gpc := geoPointCoder{}
	var databaseObject interface{}
	var jsonObject types.M
	var ok bool
	var expect interface{}
	/*************************************************/
	databaseObject = "Incorrect type"
	jsonObject = gpc.databaseToJSON(databaseObject)
	expect = types.M{
		"__type":    "GeoPoint",
		"longitude": 0,
		"latitude":  0,
	}
	if reflect.DeepEqual(jsonObject, expect) == false {
		t.Error("expect:", expect, "get jsonObject:", jsonObject)
	}
	/*************************************************/
	databaseObject = types.S{20, 20, 20}
	jsonObject = gpc.databaseToJSON(databaseObject)
	expect = types.M{
		"__type":    "GeoPoint",
		"longitude": 0,
		"latitude":  0,
	}
	if reflect.DeepEqual(jsonObject, expect) == false {
		t.Error("expect:", expect, "get jsonObject:", jsonObject)
	}
	/*************************************************/
	databaseObject = types.S{20, 20}
	jsonObject = gpc.databaseToJSON(databaseObject)
	expect = types.M{
		"__type":    "GeoPoint",
		"longitude": 20,
		"latitude":  20,
	}
	if reflect.DeepEqual(jsonObject, expect) == false {
		t.Error("expect:", expect, "get jsonObject:", jsonObject)
	}
	/*************************************************/
	databaseObject = "Incorrect type"
	ok = gpc.isValidDatabaseObject(databaseObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	databaseObject = types.S{20, 20, 20}
	ok = gpc.isValidDatabaseObject(databaseObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	databaseObject = types.S{"20", "20"}
	ok = gpc.isValidDatabaseObject(databaseObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	databaseObject = types.S{20, 20}
	ok = gpc.isValidDatabaseObject(databaseObject)
	if !ok {
		t.Error("expect:", "true", "get:", ok)
	}
	// /*************************************************/
	// jsonObject = nil
	// databaseObject, _ = gpc.jsonToDatabase(jsonObject)
	// if databaseObject != nil {
	// 	t.Error("expect:", "nil", "get:", databaseObject)
	// }
	// /*************************************************/
	// jsonObject = types.M{
	// 	"__type": "File",
	// 	"name":   "pic.jpg",
	// }
	// databaseObject, _ = gpc.jsonToDatabase(jsonObject)
	// if reflect.DeepEqual("pic.jpg", databaseObject) == false {
	// 	t.Error("expect:", "pic.jpg", "get:", databaseObject)
	// }
	// /*************************************************/
	// jsonObject = nil
	// ok = gpc.isValidJSON(jsonObject)
	// if ok {
	// 	t.Error("expect:", "false", "get:", ok)
	// }
	// /*************************************************/
	// jsonObject = types.M{}
	// ok = gpc.isValidJSON(jsonObject)
	// if ok {
	// 	t.Error("expect:", "false", "get:", ok)
	// }
	// /*************************************************/
	// jsonObject = types.M{"__type": "Date"}
	// ok = gpc.isValidJSON(jsonObject)
	// if ok {
	// 	t.Error("expect:", "false", "get:", ok)
	// }
	// /*************************************************/
	// jsonObject = types.M{"__type": "File"}
	// ok = gpc.isValidJSON(jsonObject)
	// if ok {
	// 	t.Error("expect:", "false", "get:", ok)
	// }
	// /*************************************************/
	// jsonObject = types.M{"__type": "File", "name": 1024}
	// ok = gpc.isValidJSON(jsonObject)
	// if ok {
	// 	t.Error("expect:", "false", "get:", ok)
	// }
	// /*************************************************/
	// jsonObject = types.M{"__type": "File", "name": "pic.jpg"}
	// ok = gpc.isValidJSON(jsonObject)
	// if !ok {
	// 	t.Error("expect:", "true", "get:", ok)
	// }
	// TODO
}

func Test_fileCoder(t *testing.T) {
	fc := fileCoder{}
	var databaseObject interface{}
	var jsonObject types.M
	var ok bool
	var expect interface{}
	/*************************************************/
	databaseObject = "pic.jpg"
	jsonObject = fc.databaseToJSON(databaseObject)
	expect = types.M{
		"__type": "File",
		"name":   "pic.jpg",
	}
	if reflect.DeepEqual(jsonObject, expect) == false {
		t.Error("expect:", expect, "get jsonObject:", jsonObject)
	}
	/*************************************************/
	databaseObject = 1024
	ok = fc.isValidDatabaseObject(databaseObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	databaseObject = "pic.jpg"
	ok = fc.isValidDatabaseObject(databaseObject)
	if !ok {
		t.Error("expect:", "true", "get:", ok)
	}
	/*************************************************/
	jsonObject = nil
	databaseObject, _ = fc.jsonToDatabase(jsonObject)
	if databaseObject != nil {
		t.Error("expect:", "nil", "get:", databaseObject)
	}
	/*************************************************/
	jsonObject = types.M{
		"__type": "File",
		"name":   "pic.jpg",
	}
	databaseObject, _ = fc.jsonToDatabase(jsonObject)
	if reflect.DeepEqual("pic.jpg", databaseObject) == false {
		t.Error("expect:", "pic.jpg", "get:", databaseObject)
	}
	/*************************************************/
	jsonObject = nil
	ok = fc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{}
	ok = fc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "Date"}
	ok = fc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "File"}
	ok = fc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "File", "name": 1024}
	ok = fc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "File", "name": "pic.jpg"}
	ok = fc.isValidJSON(jsonObject)
	if !ok {
		t.Error("expect:", "true", "get:", ok)
	}
}

func Test_valueAsDate(t *testing.T) {
	var value interface{}
	var date time.Time
	var ok bool
	/*************************************************/
	value = 1024
	date, ok = valueAsDate(value)
	if ok {
		t.Error("value:", value, "date:", date, "expect: false", "get:", ok)
	}
	/*************************************************/
	value = "Incorrect string time"
	date, ok = valueAsDate(value)
	if ok {
		t.Error("value:", value, "date:", date, "expect: false", "get:", ok)
	}
	/*************************************************/
	value = "2006-01-02T15:04:05.000Z"
	date, ok = valueAsDate(value)
	if !ok || utils.TimetoString(date) != "2006-01-02T15:04:05.000Z" {
		t.Error("value:", value, "date:", date, "expect: true 2006-01-02T15:04:05.000Z", "get:", ok, utils.TimetoString(date))
	}
	/*************************************************/
	value = time.Now().UTC()
	date, ok = valueAsDate(value)
	if !ok {
		t.Error("value:", value, "date:", date, "expect: true", "get:", ok)
	}
}
