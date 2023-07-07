package queue

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/lgdzz/vingo-utils/vingo"
	"reflect"
	"time"
)

var Redis RedisQueue

type RedisQueueConfig struct {
	Debug             *bool // 调试模式，为true时日志在控制台输出，否则记录到日志文件，默认为true
	AutoBootTime      *int  // 监控器异常自动重启时间间隔，默认3秒
	SortedSetRestTime *int  // 有序集合中没有消息时休息等待时间，默认2秒
	RetryWaitTime     *int  // 消费失败重试等待时间，默认5秒
}

type RedisQueue struct {
	Config RedisQueueConfig
}

// 初始化服务（只需要执行1次）
func InitRedisQueue(config RedisQueueConfig) {
	if config.Debug != nil {
		Redis.Config.Debug = config.Debug
	} else {
		Redis.Config.Debug = vingo.BoolPointer(true)
	}

	if config.AutoBootTime != nil {
		Redis.Config.AutoBootTime = config.AutoBootTime
	} else {
		Redis.Config.AutoBootTime = vingo.IntPointer(3)
	}

	if config.SortedSetRestTime != nil {
		Redis.Config.SortedSetRestTime = config.SortedSetRestTime
	} else {
		Redis.Config.SortedSetRestTime = vingo.IntPointer(2)
	}

	if config.RetryWaitTime != nil {
		Redis.Config.RetryWaitTime = config.RetryWaitTime
	} else {
		Redis.Config.RetryWaitTime = vingo.IntPointer(5)
	}
}

// 将消息转换为字符串类型
func (s *RedisQueue) toString(value any) string {
	var kind = reflect.TypeOf(value).Kind()
	switch kind {
	case reflect.Struct:
		return vingo.JsonToString(value)
	case reflect.String:
		return value.(string)
	default:
		panic("RedisQueue.Push未知消息数据类型")
	}
}

func (s *RedisQueue) getTopic(topic string) string {
	return fmt.Sprintf("%v.queue", topic)
}

func (s *RedisQueue) getDelayTopic(topic string) string {
	return fmt.Sprintf("%v.queue.delay", topic)
}

// 推送实时任务
// topic-消息队列主题
// value-消息内容，可选类型[struct|string]
func (s *RedisQueue) Push(topic string, value any) bool {
	r, err := vingo.Redis.RPush(s.getTopic(topic), s.toString(value)).Result()
	if err != nil {
		panic(err.Error())
	}
	return r > 0
}

// 推送延迟任务
// topic-消息队列主题
// value-消息内容，可选类型[struct|string]
// delayed-延迟时间，单位：秒，如：60秒后执行，则传入60
func (s *RedisQueue) PushDelay(topic string, value any, delayed int64) bool {
	var nowTime = time.Now()
	var score = float64(nowTime.Add(time.Duration(delayed) * time.Second).Unix())
	r, err := vingo.Redis.ZAdd(s.getDelayTopic(topic), redis.Z{Member: s.toString(value), Score: score}).Result()
	if err != nil {
		panic(err.Error())
	}
	return r >= 0
}

// 开始监听队列信息
func (s *RedisQueue) StartMonitor(topic string, handler Handler) {
	go s.monitorGuard(topic, handler)
}

// 队列监听守卫
func (s *RedisQueue) monitorGuard(topic string, handler Handler) {
	defer func() {
		if err := recover(); err != nil {
			// 等待3秒后重启监听器
			time.Sleep(time.Second * time.Duration(*s.Config.AutoBootTime))
			if *s.Config.Debug {
				fmt.Println(fmt.Sprintf("[消息队列]监听器异常，进行重启."))
			} else {
				vingo.LogError(fmt.Sprintf("[消息队列]监听器异常，进行重启."))
			}
			s.monitorGuard(topic, handler)
		}
	}()
	s.monitor(topic, handler)
}

// 队列监听
func (s *RedisQueue) monitor(topic string, handler Handler) {
	topicQueue := s.getTopic(topic)
	for {
		r, err := vingo.Redis.BLPop(0, topicQueue).Result()
		if err != nil {
			panic(err.Error())
		}
		value := r[1]
		func(topic string, value string) {
			defer func() {
				if err := recover(); err != nil {
					// 如果消息处理异常，则将任务推送到延迟队列，在指定时间后再次消费
					s.PushDelay(topic, value, int64(*s.Config.RetryWaitTime))
				}
			}()
			// 执行消息处理
			handler.HandleMessage(&value)
		}(topic, value)
	}
}

// 开始监听队列信息(延迟)
func (s *RedisQueue) StartMonitorDelay(topic string) {
	go s.monitorGuardDelay(topic)
}

// 队列监听守卫(延迟)
func (s *RedisQueue) monitorGuardDelay(topic string) {
	defer func() {
		if err := recover(); err != nil {
			// 等待3秒后重启监听器
			time.Sleep(time.Second * time.Duration(*s.Config.AutoBootTime))
			if *s.Config.Debug {
				fmt.Println(fmt.Sprintf("[消息队列]监听器delay异常，进行重启."))
			} else {
				vingo.LogError(fmt.Sprintf("[消息队列]监听器delay异常，进行重启."))
			}
			s.monitorGuardDelay(topic)
		}
	}()
	s.monitorDelay(topic)
}

// 队列监听(延迟)
func (s *RedisQueue) monitorDelay(topic string) {
	topicDelay := s.getDelayTopic(topic)
	for {
		r, err := vingo.Redis.ZRangeWithScores(topicDelay, 0, 0).Result()
		if err != nil {
			panic(err.Error())
		}
		if len(r) > 0 {
			// 将分数转换为时间戳
			// 获取第一个记录的分数和成员
			score := r[0].Score
			member := r[0].Member.(string)
			// 将分数转换为时间戳
			expiryTime := time.Unix(int64(score), 0)
			// 当前时间
			now := time.Now()
			// 判断是否已经到期
			if expiryTime.After(now) {
				// 计算剩余时间
				remainingTime := expiryTime.Sub(now)
				// 暂停等待剩余时间
				if remainingTime.Seconds() < float64(*s.Config.SortedSetRestTime) {
					// 剩余时间小于休息时间，则按剩余时间暂停
					time.Sleep(remainingTime)
				} else {
					// 否则直接用休息时间暂停
					time.Sleep(time.Second * time.Duration(*s.Config.SortedSetRestTime))
				}
			} else {
				// 将任务加入到实时队列
				s.Push(topic, member)
				// 删除记录有序集合中的记录
				vingo.Redis.ZRem(topicDelay, member)
			}
		} else {
			// 有序集合中没有消息时休息等待
			time.Sleep(time.Second * time.Duration(*s.Config.SortedSetRestTime))
		}
	}
}

type Handler interface {
	HandleMessage(message *string)
}
