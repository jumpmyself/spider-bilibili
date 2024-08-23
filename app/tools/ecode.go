package tools

import "fmt"

var (
	OK  = ECode{Code: 0}
	Err = ECode{Code: 10004, Message: "获取数据失败"}
)

type ECode struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (e *ECode) String() string {
	return fmt.Sprintf("code；%d,message:%s", e.Code, e.Message)
}
