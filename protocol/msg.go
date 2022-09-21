package protocol

import (
	"encoding/gob"
)

func init() {
	gob.Register(LoginReq{})
	gob.Register(LoginRsp{})
	gob.Register(QueryReq{})
	gob.Register(EditPeq{})
	gob.Register(QueryRsp{})
	gob.Register(StatusSuccess)
	gob.Register(StatusFail)
}

type Status int32

const (
	StatusUnknown Status = iota
	StatusSuccess
	StatusFail
)

//用户登录消息题
type LoginReq struct {
	Username string
	Passwd   string
}

type LoginRsp struct {
	Token string
	Info  QueryRsp
}

//修改请求消息体，修改昵称、图片使用
type EditPeq struct {
	Query   QueryReq
	Content string
}

//查询用户信息消息体
type QueryReq struct {
	Username string
	Token    string
}

type QueryRsp struct {
	Nickname string
	PicUrl   string
}
