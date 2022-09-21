package server

import (
	"encoding/json"
	"fmt"
	"go-entry-task/common"
	"go-entry-task/protocol"
	"go-entry-task/rpcsdk/rpc"
	"go-entry-task/tcpserver/cache"
	"go-entry-task/tcpserver/db"
	"go-entry-task/tcpserver/util"
	"log"
	"math/rand"
	"time"
)

const (
	tokenTimeout = 90  //存入redis时token时长
	infoTimeout  = 120 //存入redis时用户信息时长
)

type Server struct {
}

func (s *Server) register() {
	rpc.Register("login", s.LoginHandler)
	rpc.Register("query", s.QueryHandler)
	rpc.Register("updateNick", s.EditNickHandler)
	//rpc.Register("logout",HandlerLogout)
	rpc.Register("uploadpic", s.EditPicHandler)
}

//用户登录
func (s *Server) LoginHandler(msg protocol.LoginReq) (protocol.LoginRsp, error) {
	var rsp protocol.LoginRsp
	//先查询redis中有没有数据
	nick, pass, url, err := cache.QueryUser(msg.Username)
	if err == nil { //redis有数据
		if pass == msg.Passwd { //密码正确
			token := generateToken(msg.Username)
			err = cache.Set(token, msg.Username, tokenTimeout)
			if err != nil {
				return rsp, err
			}
			rsp = protocol.LoginRsp{Token: token, Info: protocol.QueryRsp{Nickname: nick, PicUrl: url}}
			return rsp, nil
		} else {
			return rsp, util.ErrInvalidPassword
		}
	}
	//redis中没有数据，从mysql中查找
	nick, pass, url, err = db.Query(msg.Username)
	if err != nil {
		log.Printf("query mysql err:%v\n", err)
		return rsp, err
	}
	if pass != msg.Passwd {
		return rsp, util.ErrInvalidPassword
	}
	//插入redis
	token := generateToken(msg.Username)
	err = cache.Set(token, msg.Username, tokenTimeout)
	if err != nil {
		return rsp, err
	}
	hash := map[string]string{"url": url, "pass": pass, "nick": nick}
	text, _ := json.Marshal(hash)
	err = cache.Set(msg.Username, text, infoTimeout)
	if err != nil {
		log.Println(err)
	}
	rsp = protocol.LoginRsp{Token: token, Info: protocol.QueryRsp{Nickname: nick, PicUrl: url}}
	return rsp, nil
}

//用户查询
func (s *Server) QueryHandler(msg protocol.QueryReq) (protocol.QueryRsp, error) {
	var rsp protocol.QueryRsp
	if !cache.QueryToken(msg.Token, msg.Username) {
		return rsp, util.ErrTokenNotExist
	}
	nick, pass, url, err := cache.QueryUser(msg.Username)
	if err == nil {
		rsp = protocol.QueryRsp{Nickname: nick, PicUrl: url}
		return rsp, nil
	}
	//redis中没有数据,去mysql中查找
	nick, pass, url, err = db.Query(msg.Username)
	if err != nil {
		return rsp, err
	}
	hash := make(map[string]string)
	hash["url"] = url
	hash["pass"] = pass
	hash["nick"] = nick
	text, _ := json.Marshal(hash)
	err = cache.Set(msg.Username, text, 120)
	rsp = protocol.QueryRsp{Nickname: nick, PicUrl: url}
	return rsp, nil
}

//用户修改昵称
func (s *Server) EditNickHandler(msg protocol.EditPeq) (protocol.Status, error) {
	if !cache.QueryToken(msg.Query.Token, msg.Query.Username) {
		return protocol.StatusFail, util.ErrTokenNotExist
	}
	err := db.UpdateNic(msg.Content, msg.Query.Username)
	if err != nil {
		return protocol.StatusFail, util.ErrUpdateDB
	}
	err = cache.Del(msg.Query.Username)
	if err != nil {
		return protocol.StatusFail, util.ErrDelRedis
	}
	return protocol.StatusSuccess, nil
}

//用户修改简历
func (s *Server) EditPicHandler(msg protocol.EditPeq) (protocol.Status, error) {
	if !cache.QueryToken(msg.Query.Token, msg.Query.Username) {
		return protocol.StatusFail, util.ErrTokenNotExist
	}
	//先更新db，再删除缓存
	err := db.UpdatePic(msg.Content, msg.Query.Username)
	if err != nil {
		return protocol.StatusFail, util.ErrUpdateDB
	}
	err = cache.Del(msg.Query.Username)
	if err != nil {
		return protocol.StatusFail, util.ErrDelRedis
	}
	return protocol.StatusSuccess, nil
}

func generateToken(user string) string {
	second := time.Now().Unix()
	rand.Seed(time.Now().UnixNano())
	return common.GetMD5(fmt.Sprintf("%v%v%v", user, second, rand.Int()))
}

//用户登出
//func HandlerLogout(msg protocol.QueryReq)(int,error){
//	if !cache.QueryToken(msg.Token,msg.Username){
//		return 1,errors.New("token doesn't exist")
//	}
//	err := cache.Del(msg.Token)
//	if err!=nil{
//		return 3,errors.New("del redis error")
//	}
//	return 0,nil
//}
