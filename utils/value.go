package utils

import "time"

type value struct{
	Code uint `json:"code"`
	Data interface{} `json:"data"`
}

func GetValue(code uint,data interface{}) *value{
	v := new(value)
	v.Code = code
	v.Data = data
	return v
}

func GetTimeMillis() int64 {
	return time.Now().UnixNano() / 1e6
}