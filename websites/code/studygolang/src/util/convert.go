// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package util

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// 将url.Values（表单数据）转换为Model（struct）
func ConvertAssign(dest interface{}, form url.Values) error {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return fmt.Errorf("convertAssign(non-pointer %s)", destType)
	}
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	if destValue.Kind() != reflect.Struct {
		return fmt.Errorf("convertAssign(non-struct %s)", destType)
	}
	destType = destValue.Type()
	fieldNum := destType.NumField()
	for i := 0; i < fieldNum; i++ {
		// struct 字段的反射类型（StructField）
		fieldType := destType.Field(i)
		// 非导出字段不处理
		if fieldType.PkgPath != "" {
			continue
		}
		tag := fieldType.Tag.Get("json")
		fieldValue := destValue.Field(i)
		val := form.Get(tag)
		// 字段本身的反射类型（field type）
		fieldValType := fieldType.Type
		switch fieldValType.Kind() {
		case reflect.Int:
			if len(form[tag]) > 1 {
				// TODO:多个值如何处理？
			}
			if val == "" {
				continue
			}
			tmp, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			fieldValue.SetInt(int64(tmp))
		case reflect.String:
			if len(form[tag]) > 1 {
				// TODO:多个值如何处理？
			}
			fieldValue.SetString(val)
		case reflect.Bool:
			if len(form[tag]) > 1 {
				// TODO:多个值如何处理？
			}

			var tmp bool
			if val == "1" {
				tmp = true
			}
			fieldValue.SetBool(tmp)
		default:

		}
	}
	return nil
}

func Struct2Map(dest map[string]interface{}, src interface{}) error {
	if dest == nil {
		return fmt.Errorf("Struct2Map(dest is %v)", dest)
	}
	srcType := reflect.TypeOf(src)
	srcValue := reflect.Indirect(reflect.ValueOf(src))
	if srcValue.Kind() != reflect.Struct {
		return fmt.Errorf("Struct2Map(non-struct %s)", srcType)
	}
	srcType = srcValue.Type()
	fieldNum := srcType.NumField()
	for i := 0; i < fieldNum; i++ {
		// struct 字段的反射类型（StructField）
		fieldType := srcType.Field(i)
		// 非导出字段不处理
		if fieldType.PkgPath != "" {
			continue
		}
		tag := fieldType.Tag.Get("json")
		// 有可能包含 ctime,omitempty
		tags := strings.Split(tag, ",")

		fieldValue := srcValue.Field(i)
		dest[tags[0]] = fieldValue.Interface()
	}
	return nil
}

// model中类型提取其中的 idField(int 类型) 属性组成 slice 返回
func Models2Intslice(models interface{}, idField string) []int {
	if models == nil {
		return []int{}
	}

	// 类型检查
	modelsValue := reflect.ValueOf(models)
	if modelsValue.Kind() != reflect.Slice {
		return []int{}
	}

	var modelValue reflect.Value

	length := modelsValue.Len()
	ids := make([]int, 0, length)

	for i := 0; i < length; i++ {
		modelValue = reflect.Indirect(modelsValue.Index(i))
		if modelValue.Kind() != reflect.Struct {
			continue
		}

		val := modelValue.FieldByName(idField)
		if val.Kind() != reflect.Int {
			continue
		}

		ids = append(ids, int(val.Int()))
	}

	return ids
}
