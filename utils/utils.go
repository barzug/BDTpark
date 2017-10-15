package utils

import "reflect"

func CheckEmpty(obj interface{}) bool {
	reflectObject := reflect.ValueOf(obj).Elem()
	for i := 0; i < reflectObject.NumField(); i++ {
		field := reflectObject.Field(i)
		fieldValue := field.Interface()
		if fieldValue == "" {
			return true
		}
	}
	return false
}

func AdditionObject(base, additional interface{}) {
	baseObj := reflect.ValueOf(base).Elem()
	additionalObj := reflect.ValueOf(additional).Elem()

	for i := 0; i < baseObj.NumField(); i++ {
		if baseObj.Field(i).String() == "" {
			baseObj.Field(i).SetString(additionalObj.Field(i).String())
		}
	}
}
