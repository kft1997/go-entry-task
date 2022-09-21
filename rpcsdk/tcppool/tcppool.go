package tcppool

import (
	"errors"
	"net"
	"sync"
	"time"
)

var ErrPoolClosed = errors.New("资源池已经关闭。")

type IConn struct {
	conn net.Conn
	t    time.Time
}

type TcpPool struct {
	m       sync.Mutex
	res     chan IConn
	factory func() (net.Conn, error)
	closed  bool
	maxNum  uint
	curNUm  uint
	timeout int
}

//创建一个资源池
func New(fn func() (net.Conn, error), curconn, maxconn uint, expire int) (*TcpPool, error) {
	if curconn > maxconn {
		return nil, errors.New("参数设置错误。")
	}
	p := &TcpPool{
		factory: fn,
		res:     make(chan IConn, maxconn),
		maxNum:  maxconn,
		curNUm:  curconn,
		timeout: expire,
	}
	var i uint
	for i = 0; i < curconn; i++ {
		c, err := fn()
		if err != nil {
			return nil, err
		}
		p.res <- IConn{c, time.Now()}
	}
	return p, nil
}

//从资源池里获取一个资源
func (p *TcpPool) Get() (net.Conn, error) {
	for {
		select {
		//当前有空余连接
		case wrapConn, ok := <-p.res:
			if !ok {
				return nil, ErrPoolClosed
			}
			//判断超时
			if p.timeout > 0 {
				expire := time.Second * time.Duration(p.timeout)
				if wrapConn.t.Add(expire).Before(time.Now()) {
					wrapConn.conn.Close()
					p.m.Lock()
					p.curNUm--
					p.m.Unlock()
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			p.m.Lock()
			if p.curNUm < p.maxNum {
				wrapConn, err := p.factory()
				if err != nil {
					p.m.Unlock()
					return nil, err
				}
				//log.Println("cur conn num:",p.curConn)
				p.curNUm++
				p.m.Unlock()
				return wrapConn, nil
			}
			p.m.Unlock()
			wrapConn, ok := <-p.res
			if !ok {
				return nil, ErrPoolClosed
			}
			if p.timeout > 0 {
				expire := time.Second * time.Duration(p.timeout)
				if wrapConn.t.Add(expire).Before(time.Now()) {
					wrapConn.conn.Close()
					p.m.Lock()
					p.curNUm--
					p.m.Unlock()
					continue
				}
			}
			//log.Println("Get:共享资源")
			return wrapConn.conn, nil
		}
	}
}

//关闭资源池，释放资源
func (p *TcpPool) Close() {
	p.m.Lock()
	defer p.m.Unlock()
	if p.closed {
		return
	}

	p.closed = true

	//关闭通道，不让写入了
	close(p.res) //关闭通道里的资源
	for r := range p.res {
		r.conn.Close()
	}
}

func (p *TcpPool) Put(c net.Conn) {
	if p.closed { //池子已经关闭
		c.Close()
		return
	}
	p.res <- IConn{c, time.Now()}
}
