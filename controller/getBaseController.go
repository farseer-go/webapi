package controller

import "reflect"

func getBaseController(controllerVal reflect.Value) BaseController {
	controllerVal = reflect.Indirect(controllerVal)
	for i := 0; i < controllerVal.NumField(); i++ {
		fieldVal := controllerVal.Field(i)
		if fieldVal.Type().String() == "controller.BaseController" {
			return fieldVal.Interface().(BaseController)
		}
	}
	return BaseController{}
}
