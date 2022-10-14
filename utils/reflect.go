package utils

import (
	"fmt"
	"reflect"
)

//ReflectGetVTByOrigin 获取反射的type和value，保留指针
func ReflectGetVTByOrigin(model interface{}) (t reflect.Type, v reflect.Value) {
	return reflect.TypeOf(model), reflect.ValueOf(model)
}

//ReflectGetVT 获取反射的type和value，忽略指针
func ReflectGetVT(model interface{}) (t reflect.Type, v reflect.Value) {
	t = reflect.TypeOf(model)
	v = reflect.ValueOf(model)

	if v.Kind() == reflect.Ptr {
		v = reflect.ValueOf(model).Elem()

		if v.IsValid() {
			t = v.Type()
		}
	}

	return
}

//ReflectSetValue 反射赋值
func ReflectSetValue(model interface{}, key string, value interface{}) {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	v.FieldByName(key).Set(reflect.ValueOf(value))
}

//CopyFields 复制结构体
func CopyFields(dst interface{}, src interface{}, fields ...string) (err error) {
	at := reflect.TypeOf(dst)
	av := reflect.ValueOf(dst)
	bt := reflect.TypeOf(src)
	bv := reflect.ValueOf(src)

	// 简单判断下
	if at.Kind() != reflect.Ptr {
		err = fmt.Errorf("a must be a struct pointer")
		return
	}
	av = reflect.ValueOf(av.Interface())

	// 要复制哪些字段
	_fields := make([]string, 0)
	if len(fields) > 0 {
		_fields = fields
	} else {
		for i := 0; i < bv.NumField(); i++ {
			_fields = append(_fields, bt.Field(i).Name)
		}
	}

	if len(_fields) == 0 {
		fmt.Println("no fields to copy")
		return
	}

	// 复制
	for i := 0; i < len(_fields); i++ {
		name := _fields[i]
		f := av.Elem().FieldByName(name)
		bValue := bv.FieldByName(name)

		// a中有同名的字段并且类型一致才复制
		if f.IsValid() && f.Kind() == bValue.Kind() {
			f.Set(bValue)
		} else {
			fmt.Printf("no such field or different kind, fieldName: %s\n", name)
		}
	}

	return
}
