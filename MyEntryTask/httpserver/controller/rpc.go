package controller

import (
	"MyEntryTask/tcpserver/rpc/client"
	"MyEntryTask/utils"
	"context"
	"errors"
	"google.golang.org/protobuf/proto"
	"syscall"
)

var (
	pool *client.CliPool
)

func PoolInit() (err error) {
	poolConf := &client.CliPoolConf{
		Size:   2000,
		MaxCap: 2500,
	}

	pool, err = client.NewClientPool(poolConf)
	if err != nil {
		return err
	}

	return nil
}

func CallRPC(ctx context.Context, svcMethod string, args, reply proto.Message) (code uint32, err error) {
	// 从rpc连接池获取一个rpc client
	cli, err := pool.Get()
	if err != nil {
		utils.Logs.Warn("pool get err :[%v]\n", err)
		return utils.CodeInternalErr, err
	}

	code, err = cli.CallRpc(ctx, svcMethod, args, reply)
	if err != nil {
		// 如果broken pipe
		if errors.Is(err, syscall.EPIPE) {
			// 客户端关闭
			cli.Close()
			cli = client.NewRpcClient()
			// 重试
			code, err = cli.CallRpc(ctx, svcMethod, args, reply)
			if err != nil {
				cli.Close()
				utils.Logs.Warn("retry call rpc err :[%v]\n", err)
				return utils.CodeTCPRpcServiceErr, err
			}
			err = pool.Put(cli)
			return code, nil
		}
		cli.Close()
		return utils.CodeTCPRpcServiceErr, err
	}
	err = pool.Put(cli)
	if err != nil {
		utils.Logs.Warn("pool put err :[%v]\n", err)
		return utils.CodeInternalErr, err
	}
	return
}
