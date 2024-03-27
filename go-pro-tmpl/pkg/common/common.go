package common

import (
	"encoding/json"
	"errors"
)

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

type DataList struct {
	Item  interface{} `json:"item"`
	Total uint        `json:"total"`
}

func BuildListResponse(item interface{}, total uint) Response {
	return Response{
		Status: 200,
		Data: DataList{
			Item:  item,
			Total: total,
		},
		Msg: "ok",
	}
}

func ErrorResponse(err error) Response {
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return Response{
			Status: 400,
			Msg:    "json type not match!",
			Error:  err.Error(),
		}
	}

	return Response{
		Status: 400,
		Msg:    "args error!",
		Error:  err.Error(),
	}
}
