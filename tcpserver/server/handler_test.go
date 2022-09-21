package server

import (
	"go-entry-task/common"
	"go-entry-task/protocol"
	"testing"
)

func init() {

}

func TestLoginHandler(t *testing.T) {
	req := protocol.LoginReq{Username: "user_0", Passwd: common.GetMD5("pass_0")}
	_, err := server.LoginHandler(req)
	if err != nil {
		t.Error(err)
	}
}

func TestQueryHandler(t *testing.T) {
	loginReq := protocol.LoginReq{Username: "user_0", Passwd: common.GetMD5("pass_0")}
	loginRsp, err := server.LoginHandler(loginReq)
	if err != nil {
		t.Error(err)
	}
	req := protocol.QueryReq{Username: "user_0", Token: loginRsp.Token}
	_, err = server.QueryHandler(req)
	if err != nil {
		t.Error(err)
	}
}

func TestEditNickHandler(t *testing.T) {
	loginReq := protocol.LoginReq{"user_0", common.GetMD5("pass_0")}
	loginRsp, err := server.LoginHandler(loginReq)
	if err != nil {
		t.Error(err)
	}
	req := protocol.EditPeq{protocol.QueryReq{"user_0", loginRsp.Token}, "nick_0"}
	_, err = server.EditNickHandler(req)
	if err != nil {
		t.Error(err)
	}
}

func TestEditPicHandler(t *testing.T) {
	loginreq := protocol.LoginReq{"user_0", common.GetMD5("pass_0")}
	loginrsp, err := server.LoginHandler(loginreq)
	if err != nil {
		t.Error(err)
	}
	req := protocol.EditPeq{protocol.QueryReq{"user_0", loginrsp.Token}, ""}
	_, err = server.EditPicHandler(req)
	if err != nil {
		t.Error(err)
	}
}
