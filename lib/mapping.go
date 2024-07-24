package lib

import (
	"errors"
	"reflect"
)

func MapStruct(src interface{}, dst interface{}) error {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	if srcVal.Kind() != reflect.Ptr || dstVal.Kind() != reflect.Ptr {
		return errors.New("src and dst must be pointers")
	}

	srcVal = srcVal.Elem()
	dstVal = dstVal.Elem()

	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		return errors.New("src and dst must be structs")
	}

	srcType := srcVal.Type()
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		srcFieldType := srcType.Field(i)

		dstField := dstVal.FieldByName(srcFieldType.Name)
		if dstField.IsValid() && dstField.CanSet() {
			dstField.Set(srcField)
		}
	}

	return nil
}
