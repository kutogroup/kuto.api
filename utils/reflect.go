package utils

import "reflect"

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
