package client

import (
	"MyEntryTask/utils"
	"errors"
	"sync"
	"time"
)

const (
	kNetwork = "tcp"
)

type idleConn struct {
	cli  *RPCClient
	time time.Time
}

type CliPoolConf struct {
	Size   int
	MaxCap int
}

type CliPool struct {
	mu      sync.Mutex
	conns   chan *idleConn
	factory func() *RPCClient
	timeout time.Duration
	maxCap  int
	size    int
	cloesd  bool
}

var (
	ErrClosed        = errors.New("pool is closed")
	ErrPoolNotEnough = errors.New("pool is not enough")
)

func NewClientPool(cpConf *CliPoolConf) (*CliPool, error) {
	cp := &CliPool{
		conns:   make(chan *idleConn, cpConf.Size),
		factory: NewRpcClient,
		timeout: 10 * time.Minute,
		maxCap:  cpConf.MaxCap,
		size:    cpConf.Size,
		mu:      sync.Mutex{},
	}

	for i := 0; i < cp.size; i++ {
		cli := cp.factory()
		// 连接池连接
		err := cli.Dial(make(chan interface{}, 1), kNetwork, RpcHost)
		if err != nil {
			utils.Logs.Warn("pool cli dial err: [%v]\n", err)
			continue
		}
		cp.conns <- &idleConn{cli, time.Now()}
	}

	return cp, nil
}

func (cp *CliPool) Get() (*RPCClient, error) {
	for {
		select {
		case conn := <-cp.conns:
			if conn == nil {
				return nil, ErrClosed
			}
			// 连接是否超时 超时获取新的
			if conn.time.Add(cp.timeout).Before(time.Now()) {
				cp.Release(conn)
				newConn := cp.factory()
				return newConn, nil
			}

			return conn.cli, nil
		default:
			// 连接池数量不够用
			cp.mu.Lock()
			if cp.size >= cp.maxCap {
				cp.mu.Unlock()
				utils.Logs.Warn("pool is not enough")
				return nil, ErrPoolNotEnough
			}
			newConn := cp.factory()
			cp.size++
			cp.mu.Unlock()
			return newConn, nil
		}
	}
}

func (cp *CliPool) Put(cli *RPCClient) error {
	if cli == nil {
		return errors.New("connection is nil. reject")
	}
	if cp.cloesd {
		utils.Logs.Warn("pool is closed")
		return ErrClosed
	}
	cp.conns <- &idleConn{cli, time.Now()}
	return nil
}

// Close 关闭连接池
func (cp *CliPool) Close() {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.factory = nil
	cp.size = 0
	cp.maxCap = 0
	cp.cloesd = true

	for conn := range cp.conns {
		cp.Release(conn)
	}
	close(cp.conns)
}

// Release 释放rpc client 连接
func (cp *CliPool) Release(conn *idleConn) {
	conn.cli.Close()
}
