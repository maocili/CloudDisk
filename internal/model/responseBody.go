package model

/*
	：:TODO ResponseBody 加锁
*/
type ResponseBody struct {
	Code int         `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty" `
}
