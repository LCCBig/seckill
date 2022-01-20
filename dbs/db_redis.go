package dbs

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	redisClient *redis.Client
)

/**
参数详解
http://www.lsdcloud.com/go/middleware/go-redis.html#_3-1-redis-options%E5%8F%82%E6%95%B0%E8%AF%A6%E8%A7%A3
*/
func IntiRedisClient(ctx context.Context) {
	redisClient = redis.NewClient(&redis.Options{
		//// 网络类型 tcp 或者 unix.
		//	// 默认是 tcp.
		Network: "",
		//redis地址，格式 host:port
		Addr:   viper.GetString("redis.addr"),
		Dialer: nil,
		// 新建一个redis连接的时候，会回调这个函数
		OnConnect: nil,
		//redis用户名，redis server没有设置可以为空。
		Username: "",
		//redis密码，redis server没有设置可以为空。
		Password: "",
		//redis数据库，序号从0开始，默认是0，可以不用设置
		DB: viper.GetInt("redis.db"),
		//redis操作失败最大重试次数，默认不重试。
		MaxRetries: 0,
		//最小重试时间间隔.
		//默认是 8ms ; -1 表示关闭
		MinRetryBackoff: 0,
		//最大重试时间间隔
		//默认是 512ms; -1 表示关闭.
		MaxRetryBackoff: 0,
		//redis连接超时时间.
		//默认是 5 秒.
		DialTimeout: 0,
		//socket读取超时时间
		//默认 3 秒.
		ReadTimeout: 0,
		//socket写超时时间
		WriteTimeout: 0,
		//连接池的类型。
		//FIFO 池为 true，LIFO 池为 false。
		//请注意，与 lifo 相比，fifo 具有更高的开销。
		PoolFIFO: false,
		//redis连接池的最大连接数.
		//默认连接池大小等于 cpu个数 * 10
		PoolSize: 0,
		//redis连接池最小空闲连接数.
		MinIdleConns: viper.GetInt("redis.min-idle-conns"),
		//redis连接最大的存活时间，默认不会关闭过时的连接.
		MaxConnAge: 0,
		//当你从redis连接池获取一个连接之后，连接池最多等待这个拿出去的连接多长时间。
		//默认是等待 ReadTimeout + 1 秒.
		PoolTimeout: 0,
		//redis连接池多久会关闭一个空闲连接.
		//默认是 5 分钟. -1 则表示关闭这个配置项
		IdleTimeout: 0,
		//多长时间检测一下，空闲连接
		//默认是 1 分钟. -1 表示关闭空闲连接检测
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
		Limiter:            nil,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		print(err)
	}
}

func GetRedisClinet() *redis.Client {
	return redisClient
}
