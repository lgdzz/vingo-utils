package vingo

import (
	"fmt"
	"github.com/go-redis/redis"
)

// Deprecated: This function is no longer recommended for use.
// Suggested: Please use redis.Option instead.
type RedisConfig struct {
	Host         string `yaml:"host" json:"host"`
	Port         string `yaml:"port" json:"port"`
	Select       int    `yaml:"select" json:"select"`
	Password     string `yaml:"password" json:"password"`
	PoolSize     int    `yaml:"poolSize" json:"poolSize"`
	MinIdleConns int    `yaml:"minIdleConns" json:"minIdleConns"`
}

// Deprecated: This function is no longer recommended for use.
// Suggested: Please use redis.Client instead.
var Redis *redis.Client

// redis初始化
// Deprecated: This function is no longer recommended for use.
// Suggested: Please use redis.InitClient() instead.
func InitRedisService(config *RedisConfig) {

	if config.Host == "" {
		config.Host = "127.0.0.1"
	}

	if config.Port == "" {
		config.Port = "6379"
	}

	if config.PoolSize == 0 {
		config.PoolSize = 4
	}

	if config.MinIdleConns == 0 {
		config.PoolSize = 2
	}

	Redis = redis.NewClient(&redis.Options{
		//连接信息
		Network:  "tcp",                                          //网络类型，tcp or unix，默认tcp
		Addr:     fmt.Sprintf("%v:%v", config.Host, config.Port), //主机名+冒号+端口，默认localhost:6379
		Password: config.Password,                                //密码
		DB:       config.Select,                                  // redis数据库index

		//连接池容量及闲置连接数量
		PoolSize:     config.PoolSize,     // 连接池最大socket连接数，应该设置为服务器CPU核心数的两倍
		MinIdleConns: config.MinIdleConns, // 在启动阶段创建指定数量的Idle连接，一般来说，可以将其设置为PoolSize的一半
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis连接异常：%v", err.Error()))
	}
}

// Deprecated: This function is no longer recommended for use.
// Suggested: Please use redis.RedisResult() instead.
func RedisResult(cmd *redis.StringCmd) string {
	result, err := cmd.Result()
	if err == redis.Nil {
		//fmt.Println("没有这个值")
		return ""
	} else if err != nil {
		//fmt.Println(err)
		panic(err.Error())
	} else {
		return result
	}
}

// Deprecated: This function is no longer recommended for use.
// Suggested: Please use redis.RedisSaveResult() instead.
func RedisSaveResult(err error) {
	if err != nil {
		panic(err.Error())
	}
}
