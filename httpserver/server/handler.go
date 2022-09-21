package server

import (
	"encoding/json"
	"fmt"
	"go-entry-task/common"
	"go-entry-task/protocol"
	"go-entry-task/rpcsdk/rpc"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	viewFile  = "./template/view.html"
	indexFile = "./template/index.html"
	picFormat = "./picfile/%v_%v.%v"
)

type Server struct {
}

func (s *Server) ViewHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadFile(viewFile)
	fmt.Fprintf(w, string(body))
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("user")
	password := r.Form.Get("pass") //前端没有做加密，api层加密，模拟真实用户请求
	pass := common.GetMD5(password)
	req := protocol.LoginReq{Username: username, Passwd: pass}
	var rsp protocol.LoginRsp
	err := rpc.RpcCall("login", req, &rsp)
	if err != nil {
		log.Printf("rpc err: %v\n", err)
		s.replayHttpErr(w)
	} else {
		// set token
		cookie := http.Cookie{
			Name:     "token",
			Value:    rsp.Token,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)
		info := userInfo{username, rsp.Info.Nickname, rsp.Info.PicUrl}
		t, _ := template.ParseFiles(indexFile)
		t.Execute(w, info)
	}
}

func (s *Server) QueryHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("token")
	token := cookie.Value
	args := r.URL.Query()
	username := args.Get("username")
	req := protocol.QueryReq{Username: username, Token: token}
	var rsp protocol.QueryRsp
	err := rpc.RpcCall("query", req, &rsp)
	if err != nil {
		log.Printf("rpc err:%v", err)
		s.replayHttpErr(w)
	} else {
		info := userInfo{username, rsp.Nickname, rsp.PicUrl}
		s.replyHttp(w, response{0, info})
	}
}

func (s *Server) EditNickHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("token")
	token := cookie.Value
	r.ParseForm()
	username := r.Form.Get("user")
	nickname := r.Form.Get("nick")
	req := protocol.EditPeq{Query: protocol.QueryReq{Username: username, Token: token}, Content: nickname}
	var rsp protocol.Status
	err := rpc.RpcCall("updateNick", req, &rsp)
	if err != nil {
		log.Printf("rpc err:%v", err)
		s.replayHttpErr(w)
	} else {
		s.replyHttp(w, response{0, ""})
	}
}

func (s *Server) EditPicHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("token")
	token := cookie.Value
	args := r.URL.Query()
	username := args.Get("username")

	err := r.ParseMultipartForm(1024 * 1024)
	if err != nil {
		log.Printf("parse err:%v", err)
		s.replayHttpErr(w)
	}

	//删除旧文件
	queryReq := protocol.QueryReq{Username: username, Token: token}
	var queryRsp protocol.QueryRsp
	err = rpc.RpcCall("query", queryReq, &queryRsp)
	if err != nil {
		s.replyHttp(w, response{1, err.Error()})
	}
	oldFilePath := queryRsp.PicUrl
	if oldFilePath != "" && exist(oldFilePath) {
		if err = os.Remove(oldFilePath); err != nil {
			log.Printf("remove old file err:%s\n", err)
		}
	}

	//插入新文件
	attribute := args.Get("type") //文件后缀名
	image := r.MultipartForm.File["picture"]
	file, _ := image[0].Open()
	data, err := ioutil.ReadAll(file)
	picUrl := fmt.Sprintf(picFormat, time.Now().Unix(), username, attribute)
	err = ioutil.WriteFile(picUrl, data, 0644)
	if err != nil {
		log.Printf("write file err:%v", err)
		s.replayHttpErr(w)
	}
	//把文件url发送到tcpserver中
	req := protocol.EditPeq{protocol.QueryReq{username, token}, picUrl}
	var rsp protocol.Status
	err = rpc.RpcCall("uploadpic", req, &rsp)
	if err != nil {
		log.Printf("rpc err:%v", err)
		s.replayHttpErr(w)
	} else {
		s.replyHttp(w, response{0, userInfo{"", "", picUrl}})
	}
}

//
func (s *Server) replayHttpErr(w http.ResponseWriter) {
	s.replyHttp(w, errResponse) //封装一层错误
}

func (s *Server) replyHttp(w http.ResponseWriter, r response) {
	data, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
	} else {
		w.Write(data)
	}
}

func exist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
