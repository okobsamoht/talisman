package controllers

import (
	"strconv"
	"strings"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// SchemasController 处理 /schemas 接口的请求
type SchemasController struct {
	ObjectsController
}

// Prepare 访问 /schemas 接口需要 master key
func (s *SchemasController) Prepare() {
	s.ObjectsController.Prepare()
	if s.Auth.IsMaster == false {
		s.Data["json"] = errs.ErrorMessageToMap(errs.OperationForbidden, "Need master key!")
		s.ServeJSON()
	}
}

// HandleFind 处理 schema 查找请求
// @router / [get]
func (s *SchemasController) HandleFind() {
	result, err := orm.SchemaCollection().GetAllSchemas()
	if err != nil && result == nil {
		s.Data["json"] = types.M{
			"results": types.S{},
		}
		s.ServeJSON()
		return
	}
	for i, v := range result {
		result[i] = orm.MongoSchemaToSchemaAPIResponse(v)
	}
	s.Data["json"] = types.M{
		"results": result,
	}
	s.ServeJSON()
}

// HandleGet 处理查找指定的类请求
// @router /:className [get]
func (s *SchemasController) HandleGet() {
	className := s.Ctx.Input.Param(":className")
	result, err := orm.SchemaCollection().FindSchema(className)
	if err != nil && result == nil {
		s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidClassName, "Class "+className+" does not exist.")
		s.ServeJSON()
		return
	}
	s.Data["json"] = result
	s.ServeJSON()
}

// HandleCreate 处理创建类请求，同时可匹配 / 的 POST 请求
// @router /:className [post]
func (s *SchemasController) HandleCreate() {
	className := s.Ctx.Input.Param(":className")
	var data = s.JSONBody
	if data == nil {
		s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "request body is empty")
		s.ServeJSON()
		return
	}

	bodyClassName := ""
	if data["className"] != nil && utils.String(data["className"]) != "" {
		bodyClassName = utils.String(data["className"])
	}
	if className != "" && bodyClassName != "" {
		if className != bodyClassName {
			s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidClassName, "Class name mismatch between "+bodyClassName+" and "+className+".")
			s.ServeJSON()
			return
		}
	}
	if className == "" {
		className = bodyClassName
	}
	if className == "" {
		s.Data["json"] = errs.ErrorMessageToMap(errs.MissingRequiredFieldError, "POST schemas needs a class name.")
		s.ServeJSON()
		return
	}

	schema := orm.LoadSchema(nil)
	result, err := schema.AddClassIfNotExists(className, utils.MapInterface(data["fields"]), utils.MapInterface(data["classLevelPermissions"]))
	if err != nil {
		s.Data["json"] = errs.ErrorToMap(err)
		s.ServeJSON()
		return
	}

	s.Data["json"] = orm.MongoSchemaToSchemaAPIResponse(result)
	s.ServeJSON()
}

// HandleUpdate 处理更新类请求
// @router /:className [put]
func (s *SchemasController) HandleUpdate() {
	className := s.Ctx.Input.Param(":className")
	var data = s.JSONBody
	if data == nil {
		s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "request body is empty")
		s.ServeJSON()
		return
	}

	bodyClassName := ""
	if data["className"] != nil && utils.String(data["className"]) != "" {
		bodyClassName = utils.String(data["className"])
	}
	if className != bodyClassName {
		s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidClassName, "Class name mismatch between "+bodyClassName+" and "+className+".")
		s.ServeJSON()
		return
	}

	submittedFields := types.M{}
	if data["fields"] != nil && utils.MapInterface(data["fields"]) != nil {
		submittedFields = utils.MapInterface(data["fields"])
	}

	schema := orm.LoadSchema(nil)
	result, err := schema.UpdateClass(className, submittedFields, utils.MapInterface(data["classLevelPermissions"]))
	if err != nil {
		s.Data["json"] = errs.ErrorToMap(err)
		s.ServeJSON()
		return
	}

	s.Data["json"] = result
	s.ServeJSON()
}

// HandleDelete 处理删除指定类请求
// @router /:className [delete]
func (s *SchemasController) HandleDelete() {
	className := s.Ctx.Input.Param(":className")
	if orm.ClassNameIsValid(className) == false {
		s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidClassName, orm.InvalidClassNameMessage(className))
		s.ServeJSON()
		return
	}

	exist := orm.CollectionExists(className)
	if exist == false {
		s.Data["json"] = types.M{}
		s.ServeJSON()
		return
	}

	collection := orm.AdaptiveCollection(className)
	count := collection.Count(types.M{}, types.M{})
	if count > 0 {
		s.Data["json"] = errs.ErrorMessageToMap(errs.ClassNotEmpty, "Class "+className+" is not empty, contains "+strconv.Itoa(count)+" objects, cannot drop schema.")
		s.ServeJSON()
		return
	}
	collection.Drop()

	// 从 _SCHEMA 表中删除类信息，清除相关的 _Join 表
	coll := orm.SchemaCollection()
	document, err := coll.FindAndDeleteSchema(className)
	if err != nil {
		s.Data["json"] = errs.ErrorToMap(err)
		s.ServeJSON()
		return
	}
	if document != nil {
		err = removeJoinTables(document)
		if err != nil {
			s.Data["json"] = errs.ErrorToMap(err)
			s.ServeJSON()
			return
		}
	}
	s.Data["json"] = types.M{}
	s.ServeJSON()
	return
}

// removeJoinTables 清除类中的所有关联表
// 需要查找的类型： "field":"relation<otherClass>"
// 需要删除的表明： "_Join:field:className"
func removeJoinTables(mongoSchema types.M) error {
	for field, v := range mongoSchema {
		fieldType := utils.String(v)
		if strings.HasPrefix(fieldType, "relation<") {
			collectionName := "_Join:" + field + ":" + utils.String(mongoSchema["_id"])
			err := orm.DropCollection(collectionName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Delete ...
// @router / [delete]
func (s *SchemasController) Delete() {
	s.ObjectsController.Delete()
}

// Put ...
// @router / [put]
func (s *SchemasController) Put() {
	s.ObjectsController.Put()
}
