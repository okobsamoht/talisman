package rest

import (
	"sort"
	"strings"

	"github.com/okobsamoht/tomato/cloud"
	"github.com/okobsamoht/tomato/config"
	"github.com/okobsamoht/tomato/errs"
	"github.com/okobsamoht/tomato/files"
	"github.com/okobsamoht/tomato/orm"
	"github.com/okobsamoht/tomato/types"
	"github.com/okobsamoht/tomato/utils"
)

// Query 处理查询请求的结构体
type Query struct {
	auth              *Auth
	className         string
	Where             types.M
	restOptions       types.M
	findOptions       types.M
	response          types.M
	doCount           bool
	include           [][]string
	keys              []string
	redirectKey       string
	redirectClassName string
	clientSDK         map[string]string
}

var alwaysSelectedKeys = []string{"objectId", "createdAt", "updatedAt"}

// NewQuery 组装查询对象
func NewQuery(
	auth *Auth,
	className string,
	where types.M,
	options types.M,
	clientSDK map[string]string,
) (*Query, error) {
	if auth == nil {
		auth = Nobody()
	}
	if where == nil {
		where = types.M{}
	}
	query := &Query{
		auth:              auth,
		className:         className,
		Where:             where,
		restOptions:       options,
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         clientSDK,
	}

	if auth.IsMaster == false {
		// 当前权限为 Master 时，findOptions 中不存在 acl 这个 key
		if auth.User != nil {
			query.findOptions["acl"] = []string{utils.S(auth.User["objectId"])}
		} else {
			query.findOptions["acl"] = nil
		}
		if className == "_Session" {
			if query.findOptions["acl"] == nil {
				return nil, errs.E(errs.InvalidSessionToken, "This session token is invalid.")
			}
			user := types.M{
				"user": types.M{
					"__type":    "Pointer",
					"className": "_User",
					"objectId":  auth.User["objectId"],
				},
			}
			query.Where = types.M{
				"$and": types.S{where, user},
			}
		}
	}

	keys := []string{}
	if k, ok := options["keys"]; ok {
		if s, ok := k.(string); ok {
			keys = strings.Split(s, ",")
			for i, key := range keys {
				keys[i] = strings.TrimSpace(key)
			}
		}
		// 在 includePath 中 会使用 restOptions["keys"] ，所以需要设置过滤后的数据
		options["keys"] = strings.Join(keys, ",")
	}

	// 当 keys 包含 n 级时，在 include 中自动加入 n-1 级
	if len(keys) > 0 {
		includeKeys := []string{}
		for _, key := range keys {
			if len(strings.Split(key, ".")) > 1 {
				key = key[:strings.LastIndex(key, ".")]
				includeKeys = append(includeKeys, key)
			}
		}
		keysForInclude := strings.Join(includeKeys, ",")

		if keysForInclude != "" {
			if include := utils.S(options["include"]); include == "" {
				options["include"] = keysForInclude
			} else {
				options["include"] = include + "," + keysForInclude
			}
		}
	}

	for k, v := range options {
		switch k {
		case "keys":
			if len(keys) > 0 {
				query.keys = append(keys, alwaysSelectedKeys...)
			}
		case "count":
			query.doCount = true
		case "skip":
			query.findOptions["skip"] = v
		case "limit":
			query.findOptions["limit"] = v
		case "order":
			if s, ok := v.(string); ok {
				fields := strings.Split(s, ",")
				for i, field := range fields {
					fields[i] = strings.TrimSpace(field)
				}
				// sortMap := map[string]int{}
				// for _, v := range fields {
				// 	if strings.HasPrefix(v, "-") {
				// 		sortMap[v[1:]] = -1
				// 	} else {
				// 		sortMap[v] = 1
				// 	}
				// }
				// query.findOptions["sort"] = sortMap
				query.findOptions["sort"] = fields
			}
		case "include":
			if s, ok := v.(string); ok { // v = "user.session,name.friend"
				paths := strings.Split(s, ",") // paths = ["user.session","name.friend"]
				for i, path := range paths {
					paths[i] = strings.TrimSpace(path)
				}
				pathSet := map[string]bool{}
				for _, path := range paths {
					parts := strings.Split(path, ".") // parts = ["user","session"]
					for lenght := 1; lenght <= len(parts); lenght++ {
						pathSet[strings.Join(parts[0:lenght], ".")] = true
					} // pathSet = {"user":true,"user.session":true}
				} // pathSet = {"user":true,"user.session":true,"name":true,"name.friend":true}
				pathArray := []string{}
				for k := range pathSet {
					pathArray = append(pathArray, k)
				}
				sort.Strings(pathArray) // pathArray = ["name","name.friend","user","user.session"]
				for _, set := range pathArray {
					query.include = append(query.include, strings.Split(set, "."))
				} // query.include = [["name"],["name","friend"],["user"],["user","seeeion"]]
			}
		case "redirectClassNameForKey":
			if s, ok := v.(string); ok {
				query.redirectKey = s
				query.redirectClassName = ""
			}
		default:
			return nil, errs.E(errs.InvalidJSON, "bad option: "+k)
		}
	}

	return query, nil
}

// Execute 执行查询请求，返回的数据包含 results count 两个字段
func (q *Query) Execute(executeOptions ...types.M) (types.M, error) {

	err := q.BuildRestWhere()
	if err != nil {
		return nil, err
	}
	err = q.runFind(executeOptions...)
	if err != nil {
		return nil, err
	}
	err = q.runCount()
	if err != nil {
		return nil, err
	}
	err = q.handleInclude()
	if err != nil {
		return nil, err
	}
	err = q.runAfterFindTrigger()
	if err != nil {
		return nil, err
	}
	return q.response, nil
}

// BuildRestWhere 展开查询参数，组装设置项
func (q *Query) BuildRestWhere() error {
	err := q.getUserAndRoleACL()
	if err != nil {
		return err
	}
	err = q.redirectClassNameForKey()
	if err != nil {
		return err
	}
	err = q.validateClientClassCreation()
	if err != nil {
		return err
	}
	err = q.replaceSelect()
	if err != nil {
		return err
	}
	err = q.replaceDontSelect()
	if err != nil {
		return err
	}
	err = q.replaceInQuery()
	if err != nil {
		return err
	}
	err = q.replaceNotInQuery()
	if err != nil {
		return err
	}
	q.replaceEquality()
	return nil
}

// getUserAndRoleACL 获取当前用户角色信息，以及用户 id，添加到设置项 acl 中
func (q *Query) getUserAndRoleACL() error {
	if q.auth.IsMaster || q.auth.User == nil {
		return nil
	}
	acl := []string{utils.S(q.auth.User["objectId"])}
	roles := q.auth.GetUserRoles()
	acl = append(acl, roles...)
	q.findOptions["acl"] = acl
	return nil
}

// redirectClassNameForKey 修改 className 为 redirectKey 字段对应的相关类型
func (q *Query) redirectClassNameForKey() error {
	if q.redirectKey == "" {
		return nil
	}

	newClassName := orm.TomatoDBController.RedirectClassNameForKey(q.className, q.redirectKey)
	q.className = newClassName
	q.redirectClassName = newClassName

	return nil
}

// validateClientClassCreation 验证当前请求是否能创建类
func (q *Query) validateClientClassCreation() error {
	// 检测配置项是否允许
	if config.TConfig.AllowClientClassCreation {
		return nil
	}
	if q.auth.IsMaster {
		return nil
	}
	// 允许操作系统表
	for _, v := range orm.SystemClasses {
		if v == q.className {
			return nil
		}
	}
	// 允许操作已存在的表
	schema := orm.TomatoDBController.LoadSchema(nil)
	hasClass := schema.HasClass(q.className)
	if hasClass {
		return nil
	}

	// 无法操作不存在的表
	return errs.E(errs.OperationForbidden, "This user is not allowed to access non-existent class: "+q.className)
}

// replaceSelect 执行 $select 中的查询语句，把结果放入 $in 中，替换掉 $select
// 替换前的格式如下：
// {
//     "hometown":{
//         "$select":{
//             "query":{
//                 "className":"Team",
//                 "where":{
//                     "winPct":{
//                         "$gt":0.5
//                     }
//                 }
//             },
//             "key":"city"
//         }
//     }
// }
// 转换后格式如下
// {
//     "hometown":{
//         "$in":["abc","cba"]
//     }
// }
func (q *Query) replaceSelect() error {
	selectObject := findObjectWithKey(q.Where, "$select")
	if selectObject == nil {
		return nil
	}
	// 必须包含两个 key ： query key
	selectValue := utils.M(selectObject["$select"])
	if selectValue == nil || len(selectValue) != 2 {
		return errs.E(errs.InvalidQuery, "improper usage of $select")
	}
	queryValue := utils.M(selectValue["query"])
	key := utils.S(selectValue["key"])
	if queryValue == nil || key == "" {
		return errs.E(errs.InvalidQuery, "improper usage of $select")
	}
	className := utils.S(queryValue["className"])
	where := utils.M(queryValue["where"])
	// iOS SDK 中不设置 where 时，没有 where 字段，所以此处不检测 where
	if className == "" {
		return errs.E(errs.InvalidQuery, "improper usage of $select")
	}
	if where == nil {
		where = types.M{}
	}

	// where 与 className 之外的字段默认为 Options
	delete(queryValue, "where")
	delete(queryValue, "className")
	additionalOptions := queryValue

	query, err := NewQuery(q.auth, className, where, additionalOptions, q.clientSDK)
	if err != nil {
		return err
	}
	response, err := query.Execute()
	if err != nil {
		return err
	}
	// 组装查询到的对象
	values := []types.M{}
	if utils.HasResults(response) == true {
		for _, v := range utils.A(response["results"]) {
			result := utils.M(v)
			if result == nil {
				continue
			}
			values = append(values, result)
		}
	}
	// 替换 $select 为 $in
	transformSelect(selectObject, key, values)
	// 继续搜索替换
	return q.replaceSelect()
}

// replaceDontSelect 执行 $dontSelect 中的查询语句，把结果放入 $nin 中，替换掉 $select
// 数据结构与 replaceSelect 类似
func (q *Query) replaceDontSelect() error {
	dontSelectObject := findObjectWithKey(q.Where, "$dontSelect")
	if dontSelectObject == nil {
		return nil
	}
	// 必须包含两个 key ： query key
	dontSelectValue := utils.M(dontSelectObject["$dontSelect"])
	if dontSelectValue == nil || len(dontSelectValue) != 2 {
		return errs.E(errs.InvalidQuery, "improper usage of $dontSelect")
	}
	queryValue := utils.M(dontSelectValue["query"])
	key := utils.S(dontSelectValue["key"])
	if queryValue == nil || key == "" {
		return errs.E(errs.InvalidQuery, "improper usage of $dontSelect")
	}
	className := utils.S(queryValue["className"])
	where := utils.M(queryValue["where"])
	// iOS SDK 中不设置 where 时，没有 where 字段，所以此处不检测 where
	if className == "" {
		return errs.E(errs.InvalidQuery, "improper usage of $dontSelect")
	}
	if where == nil {
		where = types.M{}
	}

	// where 与 className 之外的字段默认为 Options
	delete(queryValue, "where")
	delete(queryValue, "className")
	additionalOptions := queryValue

	query, err := NewQuery(q.auth, className, where, additionalOptions, q.clientSDK)
	if err != nil {
		return err
	}
	response, err := query.Execute()
	if err != nil {
		return err
	}
	// 组装查询到的对象
	values := []types.M{}
	if utils.HasResults(response) == true {
		for _, v := range utils.A(response["results"]) {
			result := utils.M(v)
			if result == nil {
				continue
			}
			values = append(values, result)
		}
	}
	// 替换 $dontSelect 为 $nin
	transformDontSelect(dontSelectObject, key, values)
	// 继续搜索替换
	return q.replaceDontSelect()
}

// replaceInQuery 执行 $inQuery 中的查询语句，把结果放入 $in 中，替换掉 $inQuery
// 替换前的格式：
// {
//     "post":{
//         "$inQuery":{
//             "where":{
//                 "image":{
//                     "$exists":true
//                 }
//             },
//             "className":"Post"
//         }
//     }
// }
// 替换后的格式
// {
//     "post":{
//         "$in":[
// 			{
// 				"__type":    "Pointer",
// 				"className": "className",
// 				"objectId":  "objectId",
// 			},
// 			{...}
// 		]
//     }
// }
func (q *Query) replaceInQuery() error {
	inQueryObject := findObjectWithKey(q.Where, "$inQuery")
	if inQueryObject == nil {
		return nil
	}
	// 必须包含两个 key ： where className
	inQueryValue := utils.M(inQueryObject["$inQuery"])
	if inQueryValue == nil {
		return errs.E(errs.InvalidQuery, "improper usage of $inQuery")
	}
	where := utils.M(inQueryValue["where"])
	className := utils.S(inQueryValue["className"])
	if where == nil || className == "" {
		return errs.E(errs.InvalidQuery, "improper usage of $inQuery")
	}

	// where 与 className 之外的字段默认为 Options
	delete(inQueryValue, "where")
	delete(inQueryValue, "className")
	additionalOptions := inQueryValue

	query, err := NewQuery(q.auth, className, where, additionalOptions, q.clientSDK)
	if err != nil {
		return err
	}
	response, err := query.Execute()
	if err != nil {
		return err
	}
	// 组装查询到的对象
	values := []types.M{}
	if utils.HasResults(response) == true {
		for _, v := range utils.A(response["results"]) {
			result := utils.M(v)
			if result == nil {
				continue
			}
			values = append(values, result)
		}
	}
	// 替换 $inQuery 为 $in
	transformInQuery(inQueryObject, query.className, values)
	// 继续搜索替换
	return q.replaceInQuery()
}

// replaceNotInQuery 执行 $notInQuery 中的查询语句，把结果放入 $nin 中，替换掉 $notInQuery
// 数据格式与 replaceInQuery 类似
func (q *Query) replaceNotInQuery() error {
	notInQueryObject := findObjectWithKey(q.Where, "$notInQuery")
	if notInQueryObject == nil {
		return nil
	}
	// 必须包含两个 key ： where className
	notInQueryValue := utils.M(notInQueryObject["$notInQuery"])
	if notInQueryValue == nil {
		return errs.E(errs.InvalidQuery, "improper usage of $notInQuery")
	}
	where := utils.M(notInQueryValue["where"])
	className := utils.S(notInQueryValue["className"])
	if where == nil || className == "" {
		return errs.E(errs.InvalidQuery, "improper usage of $notInQuery")
	}

	// where 与 className 之外的字段默认为 Options
	delete(notInQueryValue, "where")
	delete(notInQueryValue, "className")
	additionalOptions := notInQueryValue

	query, err := NewQuery(q.auth, className, where, additionalOptions, q.clientSDK)
	if err != nil {
		return err
	}
	response, err := query.Execute()
	if err != nil {
		return err
	}
	// 组装查询到的对象
	values := []types.M{}
	if utils.HasResults(response) == true {
		for _, v := range utils.A(response["results"]) {
			result := utils.M(v)
			if result == nil {
				continue
			}
			values = append(values, result)
		}
	}
	// 替换 $notInQuery 为 $nin
	transformNotInQuery(notInQueryObject, query.className, values)
	// 继续搜索替换
	return q.replaceNotInQuery()
}

func (q *Query) replaceEquality() {
	for key := range q.Where {
		q.Where[key] = replaceEqualityConstraint(q.Where[key])
	}
}

// runFind 从数据库查找数据，并处理返回结果
func (q *Query) runFind(executeOptions ...types.M) error {
	var options types.M
	if len(executeOptions) > 0 {
		options = executeOptions[0]
	} else {
		options = types.M{}
	}

	if q.findOptions["limit"] != nil {
		if l, ok := q.findOptions["limit"].(float64); ok {
			if l == 0 {
				q.response["results"] = types.S{}
				return nil
			}
		} else if l, ok := q.findOptions["limit"].(int); ok {
			if l == 0 {
				q.response["results"] = types.S{}
				return nil
			}
		}
	}

	findOptions := types.M{}
	for k, v := range q.findOptions {
		findOptions[k] = v
	}

	if len(q.keys) > 0 {
		keys := []string{}
		for _, k := range q.keys {
			keys = append(keys, strings.Split(k, ".")[0])
		}
		findOptions["keys"] = keys
	}
	if v, ok := options["op"].(string); ok && v != "" {
		findOptions["op"] = v
	}
	response, err := orm.TomatoDBController.Find(q.className, q.Where, findOptions)
	if err != nil {
		return err
	}
	// 从 _User 表中删除敏感字段
	if q.className == "_User" {
		for _, v := range response {
			if user := utils.M(v); user != nil {
				cleanResultOfSensitiveUserInfo(user, q.auth)
				cleanResultAuthData(user)
			}
		}
	}

	// 展开文件类型
	files.ExpandFilesInObject(response)

	if q.redirectClassName != "" {
		for _, v := range response {
			if r := utils.M(v); r != nil {
				r["className"] = q.redirectClassName
			}
		}
	}

	q.response["results"] = response
	return nil
}

// runCount 查询符合条件的结果数量
func (q *Query) runCount() error {
	if q.doCount == false {
		return nil
	}
	q.findOptions["count"] = true
	delete(q.findOptions, "skip")
	delete(q.findOptions, "limit")
	// 当需要取 count 时，数据库返回结果的第一个即为 count
	result, err := orm.TomatoDBController.Find(q.className, q.Where, q.findOptions)
	if err != nil {
		return err
	}
	if result == nil || len(result) == 0 {
		q.response["count"] = 0
	} else {
		q.response["count"] = result[0]
	}
	return nil
}

// handleInclude 展开 include 对应的内容
func (q *Query) handleInclude() error {
	if len(q.include) == 0 {
		return nil
	}
	// includePath 中会直接更新 q.response
	err := includePath(q.auth, q.response, q.include[0], q.restOptions)
	if err != nil {
		return err
	}

	if len(q.include) > 0 {
		q.include = q.include[1:]
		return q.handleInclude()
	}

	return nil
}

// runAfterFindTrigger 运行查询后的回调
func (q *Query) runAfterFindTrigger() error {
	if q.response == nil {
		return nil
	}
	results := utils.A(q.response["results"])
	hasAfterFindHook := cloud.TriggerExists(cloud.TypeAfterFind, q.className)
	if hasAfterFindHook == false {
		return nil
	}
	results, err := maybeRunAfterFindTrigger(cloud.TypeAfterFind, q.className, results, q.auth)
	if err != nil {
		return err
	}
	if results == nil {
		results = types.S{}
	}
	q.response["results"] = results

	return nil
}

// includePath 在 response 中搜索 path 路径中对应的节点，
// 查询出该节点对应的对象，然后用对象替换该节点
func includePath(auth *Auth, response types.M, path []string, restOptions types.M) error {
	if restOptions == nil {
		restOptions = types.M{}
	}
	// 查找路径对应的所有节点
	pointers := findPointers(response["results"], path)
	if len(pointers) == 0 {
		return nil
	}
	pointersHash := map[string]types.S{}
	for _, pointer := range pointers {
		// 不再区分不同 className ，添加不为空的 className
		className := utils.S(pointer["className"])
		objectID := utils.S(pointer["objectId"])
		if className != "" && objectID != "" {
			if v, ok := pointersHash[className]; ok {
				v = append(v, objectID)
				pointersHash[className] = v
			} else {
				pointersHash[className] = types.S{objectID}
			}
		}

	}

	// example1:
	// path:        []string{"user"},
	// restOptions: M{"keys": "user.id"},
	// ==>> M{"keys": "id"}
	// example2:
	// path:        []string{"user"},
	// restOptions: M{"keys": "user.id,user.name"},
	// ==>> M{"keys": "id,name"}
	// example3:
	// path:        []string{"user", "post"},
	// restOptions: M{"keys": "user.id,user.name,user.post.id"},
	// ==>> M{"keys": "id"}
	includeRestOptions := types.M{}
	if keyStr, ok := restOptions["keys"].(string); ok && keyStr != "" {
		keys := strings.Split(keyStr, ",")
		keySet := []string{}
		for _, key := range keys {
			keyPath := strings.Split(key, ".")
			if len(path) >= len(keyPath) {
				continue
			}
			m := true
			i := 0
			for i = 0; i < len(path); i++ {
				if path[i] != keyPath[i] {
					m = false
					break
				}
			}
			if m {
				keySet = append(keySet, keyPath[i])
			}
		}
		if len(keySet) > 0 {
			includeRestOptions["keys"] = strings.Join(keySet, ",")
		}
	}

	replace := types.M{}
	for clsName, ids := range pointersHash {
		// 获取所有 ids 对应的对象
		objectID := types.M{
			"$in": ids,
		}
		where := types.M{
			"objectId": objectID,
		}
		query, err := NewQuery(auth, clsName, where, includeRestOptions, nil)
		if err != nil {
			return err
		}
		includeResponse, err := query.Execute(types.M{"op": "get"})
		if err != nil {
			return err
		}
		if utils.HasResults(includeResponse) == false {
			return nil
		}

		// 组装查询到的对象
		results := utils.A(includeResponse["results"])
		for _, v := range results {
			obj := utils.M(v)
			if obj == nil {
				continue
			}
			obj["__type"] = "Object"
			obj["className"] = clsName
			if clsName == "_User" && auth.IsMaster == false {
				delete(obj, "sessionToken")
				delete(obj, "authData")
			}
			replace[utils.S(obj["objectId"])] = obj
		}
	}
	// 使用查询到的对象替换对应的节点
	replacePointers(pointers, replace)

	return nil
}

// findPointers 查询路径对应的对象列表，对象必须为 Pointer 类型
func findPointers(object interface{}, path []string) []types.M {
	if object == nil {
		return []types.M{}
	}
	// 如果是对象数组，则遍历每一个对象
	if s := utils.A(object); s != nil {
		answer := []types.M{}
		for _, v := range s {
			p := findPointers(v, path)
			answer = append(answer, p...)
		}
		return answer
	}

	// 如果不能转成 map ，则返回错误
	obj := utils.M(object)
	if obj == nil {
		return []types.M{}
	}
	// 如果当前是路径最后一个节点，判断是否为 Pointer
	if len(path) == 0 {
		if utils.S(obj["__type"]) == "Pointer" {
			return []types.M{obj}
		}
		return []types.M{}
	}
	// 取出下一个路径对应的对象，进行查找
	subobject := obj[path[0]]
	if subobject == nil {
		// 对象不存在，则不进行处理
		return []types.M{}
	}
	return findPointers(subobject, path[1:])
}

// replacePointers 把 replace 保存的对象，添加到 pointers 对应的节点中
// pointers 中保存的是指向 response 的引用，修改 pointers 中的内容，即可同时修改 response 的内容
func replacePointers(pointers []types.M, replace types.M) {
	if replace == nil {
		return
	}
	for _, pointer := range pointers {
		if pointer == nil {
			continue
		}
		objectID := utils.S(pointer["objectId"])
		if objectID == "" {
			continue
		}
		if rpl := utils.M(replace[objectID]); rpl != nil {
			// 把对象中的所有字段写入节点
			for k, v := range rpl {
				pointer[k] = v
			}
		}
	}
}

// findObjectWithKey 查找带有指定 key 的对象，root 可以是 Slice 或者 map
// 查找到一个符合条件的对象之后立即返回
func findObjectWithKey(root interface{}, key string) types.M {
	if root == nil {
		return nil
	}
	// 如果是 Slice 则遍历查找
	if s := utils.A(root); s != nil {
		for _, v := range s {
			answer := findObjectWithKey(v, key)
			if answer != nil {
				return answer
			}
		}
	}

	if m := utils.M(root); m != nil {
		// 当前 map 中存在指定的 key，表示已经找到，立即返回
		if m[key] != nil {
			return m
		}
		// 不存在指定 key 时，则遍历 map 中各对象进行查找
		for _, v := range m {
			answer := findObjectWithKey(v, key)
			if answer != nil {
				return answer
			}
		}
	}
	return nil
}

// transformSelect 转换对象中的 $select
func transformSelect(selectObject types.M, key string, objects []types.M) {
	if selectObject == nil || selectObject["$select"] == nil {
		return
	}
	values := types.S{}
	for _, result := range objects {
		if result == nil || result[key] == nil {
			continue
		}
		values = append(values, result[key])
	}

	delete(selectObject, "$select")
	var in types.S
	if v := utils.A(selectObject["$in"]); v != nil {
		in = v
		in = append(in, values...)
	} else {
		in = values
	}
	selectObject["$in"] = in
}

// transformDontSelect 转换对象中的 $dontSelect
func transformDontSelect(dontSelectObject types.M, key string, objects []types.M) {
	if dontSelectObject == nil || dontSelectObject["$dontSelect"] == nil {
		return
	}
	values := types.S{}
	for _, result := range objects {
		if result == nil || result[key] == nil {
			continue
		}
		values = append(values, result[key])
	}

	delete(dontSelectObject, "$dontSelect")
	var nin types.S
	if v := utils.A(dontSelectObject["$nin"]); v != nil {
		nin = v
		nin = append(nin, values...)
	} else {
		nin = values
	}
	dontSelectObject["$nin"] = nin
}

// transformInQuery 转换对象中的 $inQuery
func transformInQuery(inQueryObject types.M, className string, results []types.M) {
	if inQueryObject == nil || inQueryObject["$inQuery"] == nil {
		return
	}
	values := types.S{}
	for _, result := range results {
		if result == nil || utils.S(result["objectId"]) == "" {
			continue
		}
		o := types.M{
			"__type":    "Pointer",
			"className": className,
			"objectId":  result["objectId"],
		}
		values = append(values, o)
	}

	delete(inQueryObject, "$inQuery")
	var in types.S
	if v := utils.A(inQueryObject["$in"]); v != nil {
		in = v
		in = append(in, values...)
	} else {
		in = values
	}
	inQueryObject["$in"] = in
}

// transformNotInQuery 转换对象中的 $notInQuery
func transformNotInQuery(notInQueryObject types.M, className string, results []types.M) {
	if notInQueryObject == nil || notInQueryObject["$notInQuery"] == nil {
		return
	}
	values := types.S{}
	for _, result := range results {
		if result == nil || utils.S(result["objectId"]) == "" {
			continue
		}
		o := types.M{
			"__type":    "Pointer",
			"className": className,
			"objectId":  result["objectId"],
		}
		values = append(values, o)
	}

	delete(notInQueryObject, "$notInQuery")
	var nin types.S
	if v := utils.A(notInQueryObject["$nin"]); v != nil {
		nin = v
		nin = append(nin, values...)
	} else {
		nin = values
	}
	notInQueryObject["$nin"] = nin
}

// cleanResultOfSensitiveUserInfo 清除用户数据中的敏感字段
func cleanResultOfSensitiveUserInfo(result types.M, auth *Auth) {
	delete(result, "password")

	if auth.IsMaster || (auth.User != nil && utils.S(auth.User["objectId"]) == utils.S(result["objectId"])) {
		return
	}

	for _, field := range config.TConfig.UserSensitiveFields {
		delete(result, field)
	}
}

// cleanResultAuthData 清理 AuthData
func cleanResultAuthData(result types.M) {
	if authData := utils.M(result["authData"]); authData != nil {
		for provider, v := range authData {
			if v == nil {
				delete(authData, provider)
			}
		}
		if len(authData) == 0 {
			delete(result, "authData")
		}
	}
}

func replaceEqualityConstraint(constraint interface{}) interface{} {
	object := utils.M(constraint)
	if object == nil {
		return constraint
	}

	equalToObject := types.M{}
	hasDirectConstraint := false
	hasOperatorConstraint := false

	for key := range object {
		if strings.Index(key, "$") != 0 {
			hasDirectConstraint = true
			equalToObject[key] = object[key]
		} else {
			hasOperatorConstraint = true
		}
	}

	if hasDirectConstraint && hasOperatorConstraint {
		object["$eq"] = equalToObject
		for key := range equalToObject {
			delete(object, key)
		}
	}

	return constraint
}
