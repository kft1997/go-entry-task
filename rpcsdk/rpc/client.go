package rpc

import (
	"bufio"
	"errors"
	"go-entry-task/rpcsdk/tcppool"
	"log"
	"net"
	"reflect"
	"time"
)

var pool *tcppool.TcpPool

func Init(addr string) error {
	conf := defaultConfig()
	//暂不支持服务发现，需要手动将服务器ip输入
	fn := func() (net.Conn, error) {
		return net.Dial("tcp", addr)
	}
	var err error
	pool, err = tcppool.New(fn, conf.PoolInitConn, conf.PoolMaxConn, conf.PoolTimeOut)
	return err
}

func RpcCall(name string, req interface{}, rsp interface{}) error {
	msg := rpcData{name, req, ""}
	data, err := encode(msg)
	if err != nil {
		return err
	}
	//打包消息
	data = packet(data)
	conn, err := pool.Get()
	if err != nil {
		log.Printf("get conn error:%s\n", err)
		return err
	}
	//写数据
	_, err = conn.Write(data)
	if err != nil {
		log.Printf("send data err:%s\n", err)
	}
	conn.SetReadDeadline(time.Now().Add(time.Duration(defaultConfig().ConnTimeOut) * time.Second))
	//读数据
	reader := bufio.NewReader(conn)
	data, err = unpack(reader)
	if err != nil {
		log.Printf("read data err:%s\n", err)
		conn.Close()
		return err
	}
	pool.Put(conn)
	msg, err = decode(data)
	if err != nil {
		log.Printf("decode error:%s\n", err)
		return err
	}
	reflect.ValueOf(rsp).Elem().Set(reflect.ValueOf(msg.Arg))
	if msg.Err != "" {
		return errors.New(msg.Err)
	}
	return nil
}
