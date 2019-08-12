package utils

import "reflect"

type StructField struct {
	TagName   string
	FieldName string
	FieldType reflect.Type
}

//根据model获取名称
func StructGetLineName(model interface{}) string {
	t, _ := ReflectGetVT(model)
	return ConvertCamel2Line(t.Name())
}

//根据model获取所有子项
func StructGetFields(model interface{}, tag string) (cols []StructField) {
	t, _ := ReflectGetVT(model)
	n := t.NumField()

	for i := 0; i < n; i++ {
		f := t.Field(i)

		sc := StructField{
			TagName:   f.Tag.Get(tag),
			FieldName: f.Name,
			FieldType: f.Type,
		}

		cols = append(cols, sc)
	}

	return
}
