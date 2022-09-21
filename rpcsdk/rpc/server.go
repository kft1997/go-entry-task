package rpc

import (
	"bufio"
	"io"
	"log"
	"net"
	"reflect"
	"time"
)

var funcs map[string]reflect.Value

func init() {
	funcs = make(map[string]reflect.Value)
}

func Register(rpcName string, f interface{}) {
	if _, ok := funcs[rpcName]; ok {
		return
	}
	// map中没有值，则将映射添加进map,便于调用
	fVal := reflect.ValueOf(f)
	funcs[rpcName] = fVal
}

func Serve(lis net.Listener) {
	for {
		// 拿到连接
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("accept err:%v", err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(defaultConfig().ConnTimeOut) * time.Second))
		data, err := unpack(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("read err:%v\n", err)
			}
			break
		}
		rsp := call(data)
		rsp = packet(rsp)
		_, err = conn.Write(rsp)
		if err != nil {
			log.Printf("write error:%v\n", err)
			break
		}
	}
	conn.Close()
	return
}

func call(data []byte) (rsp []byte) {
	var err error
	msg, _ := decode(data)
	f, ok := funcs[msg.Name]
	if !ok {
		log.Printf("invalid func name:%s\n", msg.Name)
	}
	inArgs := make([]reflect.Value, 1, 1)
	inArgs[0] = reflect.ValueOf(msg.Arg)
	// 反射调用方法，传入参数
	out := f.Call(inArgs)
	msg.Arg = out[0].Interface()
	if out[1].Interface() != nil {
		msg.Err = out[1].Interface().(error).Error()
	}
	rsp, err = encode(msg)
	if err != nil {
		log.Printf("encode error:%s\n", err)
	}
	return rsp
}
