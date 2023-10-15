package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	redis "github.com/swapbotAA/relayer/dao/redis"
	"github.com/swapbotAA/relayer/ethereum"
	"github.com/swapbotAA/relayer/service/limitorder/bo"
	"github.com/swapbotAA/relayer/service/limitorder/order"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var (
	OrderChannels map[string]chan bool
)

// Adds a new order to the manager
func AddOrder(bo *bo.LimitOrderBO) (string, error) {
	order := new(order.Order)
	order.UserID = bo.User

	if err := order.Init(bo); err != nil {
		return "", fmt.Errorf("could not initialize order: %w", err)
	}

	// Approve Order
	//_, err := order.ApproveAmount(client)
	//if err != nil {
	//	return nil, fmt.Errorf("could not approve order: %w", err)
	//}

	// Add order to database
	newID := generateIncrementalID()
	order.ID = strconv.FormatInt(newID, 10)
	addOrder(order)

	go checkPriceAndSwap(order)
	return order.ID, nil
}

func addOrder(order *order.Order) {
	// 序列化结构数据为JSON
	jsonData, err := json.Marshal(order)
	if err != nil {
		panic(err)
	}
	redis.AddOrder(order.ID, order.UserID, jsonData)
}

func generateIncrementalID() int64 {
	// 获取当前时间戳（纳秒级别）
	timestamp := time.Now().UnixNano()

	// 生成一个随机数
	rand.Seed(time.Now().UnixNano()) // 使用当前时间作为随机种子
	randomValue := rand.Intn(1000)   // 生成一个0到999之间的随机数

	// 将时间戳和随机数组合在一起，形成唯一ID
	uniqueID := timestamp*1000 + int64(randomValue)

	return uniqueID
}

// Removes an order from the manager
//func (m *Manager) RemoveOrder(id uint) error {
//	// Get stop channel for the order
//	m.channelsMutex.Lock()
//	stopChan, ok := m.OrderChannels[id]
//	m.channelsMutex.Unlock()
//
//	if !ok {
//		return ErrOrderNotFound
//	}
//
//	// Signal the manager to stop checking prices
//	stopChan <- true
//
//	// Delete channel from map
//	m.channelsMutex.Lock()
//	delete(m.OrderChannels, id)
//	m.channelsMutex.Unlock()
//
//	return m.DB.DeleteOrder(id)
//}

// Compares a stream o prices and swaps the tokens if the price is right
func checkPriceAndSwap(o *order.Order) {
	client := ethereum.Config.Client

	// May need to make stop channel buffered
	done, stopChan := make(chan bool), make(chan bool)

	// Saves the channel for stopping the order
	OrderChannels[o.ID] = stopChan

	priceChan, errChan := o.GetPriceStream()
	log.Printf("Started watching order #%d\n", o.ID)
	for {
		select {
		case currentPrice := <-priceChan:
			log.Printf("Current price for swap: %s\n", currentPrice.Text('f', 18))
			receipt, err := o.CompAndSwap(client, currentPrice)

			if err != nil {
				log.Println(err)
				//if err := redis.DeleteOrder(o.ID); err != nil {
				//	log.Printf("Could not delete order from database: %s", err)
				//}
				return
			}

			if receipt != nil {
				done <- true
				log.Printf("Swap transaction %s <-> %s completed: %s\n", o.TokenOut, o.TokenIn, receipt.TxHash.Hex())

				// Delete channel from map
				delete(OrderChannels, o.ID)

				if err := redis.DeleteOrder(o.ID, o.UserID); err != nil {
					log.Printf("Could not delete order from database: %s", err)
				}

				return
			}
		case err := <-errChan:
			log.Println(err)
			if !errors.Is(err, order.ErrCanContinue) {
				return
			}
		case <-stopChan:
			done <- true
			return
		}
	}
}
