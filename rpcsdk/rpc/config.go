package rpc

type rpcConfig struct {
	PoolInitConn uint //连接池初始连接数
	PoolMaxConn  uint //连接池最大连接数
	PoolTimeOut  int  //连接池中连接最大限制时间
	ConnTimeOut  int  //tcp连接最大阻塞时间
}

func defaultConfig() rpcConfig {
	return rpcConfig{
		PoolInitConn: 200,
		PoolMaxConn:  500,
		PoolTimeOut:  240,
		ConnTimeOut:  300,
	}
}
