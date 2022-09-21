package config

type RedisConfig struct {
	Addr        string //地址
	PoolConn    int    //redis连接池数
	IdleTimeout int    //redis空闲超时
}

type DBConfig struct {
	Host        string //mysql域名
	MaxConn     int    //mysql最大连接数
	IdleConn    int    //mysql闲置连接数
	IdleTimeout int    //mysql空闲超时
}

func DefaultDBConfig() DBConfig {
	return DBConfig{
		Host:     "root:123456@tcp(127.0.0.1:3307)/entry_task",
		MaxConn:  100,
		IdleConn: 100,
	}
}

func DefaultRedisConfig() RedisConfig {
	return RedisConfig{
		Addr:        "127.0.0.1:6379",
		PoolConn:    1000,
		IdleTimeout: 300,
	}
}
