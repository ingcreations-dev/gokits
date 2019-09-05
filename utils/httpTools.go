package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
)

func GetForm (form interface{},formType reflect.Type,ctx *gin.Context) error {
	count := formType.NumField()
	setter := reflect.ValueOf(form).Elem()
	for i:=0; i < count; i++{
		field := formType.Field(i)
		paramName := field.Tag.Get("param")
		required,err := strconv.ParseBool(field.Tag.Get("required"))
		value := ""
		value = ctx.PostForm(paramName)
		if len(value) <= 0{
			value = ctx.Query(paramName)
			if len(value) <= 0{
				value = ctx.Param(paramName)
			}
		}

		if err == nil{
			if required && len(value) <= 0{
				errors.New(fmt.Sprintf("Parameter %s not present",paramName))
			}
		}
		fieldSetter := setter.FieldByName(field.Name)
		switch fieldSetter.Kind(){
		case reflect.Int:
			iv,_ := strconv.Atoi(value)
			fieldSetter.Set(reflect.ValueOf(iv))
		case reflect.Bool:
			bv,_ := strconv.ParseBool(value)
			fieldSetter.Set(reflect.ValueOf(bv))
		default:
			fieldSetter.Set(reflect.ValueOf(value))
		}
	}
	return nil
}

