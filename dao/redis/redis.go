package config

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/swapbotAA/relayer/settings"
	"log"
)

const (
	ORDER_PREFIX = "ORDER_"
)

var client *redis.Client

func Init(config *settings.RedisConfig) error {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password, // no password set
		DB:       0,               // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(pong, err)
		return err
	}
	return nil
}

func GetIncrId() int64 {
	// Redis键的名称
	key := "counter"

	// 使用INCR命令递增计数器
	newID, err := client.Incr(key).Result()
	if err != nil {
		log.Fatalf("Failed to increment counter: %v", err)
		return 0
	}

	return newID
}

func AddOrder(orderId string, userId string, jsonData []byte) {
	//// 序列化结构数据为JSON
	//jsonData, err := json.Marshal(order)
	//if err != nil {
	//	panic(err)
	//}

	// 保存JSON字符串到RedisSet key为userID
	err := client.SAdd(ORDER_PREFIX+userId, orderId).Err()
	if err != nil {
		panic(err)
	}

	// 保存JSON字符串到Redis key为orderId
	err = client.Set(ORDER_PREFIX+orderId, jsonData, 0).Err()
	if err != nil {
		panic(err)
	}
}

func DeleteOrder(orderId string, userId string) error {
	_, err := client.SRem(ORDER_PREFIX+userId, orderId).Result()
	if err != nil {
		return err
	}

	_, err = client.Del(ORDER_PREFIX + orderId).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetOrder(orderId string) string {
	// 从Redis获取JSON字符串
	jsonData, err := client.Get(ORDER_PREFIX + orderId).Result()
	if err != nil {
		panic(err)
	}
	return jsonData

	//// 反序列化JSON为结构
	//var order *order.Order
	//err = json.Unmarshal([]byte(jsonData), &order)
	//if err != nil {
	//	panic(err)
	//}
	//return order
}
