package main

import (
	"context"
	"encoding/json"
	"expression-agent/internal/helpers"
	"expression-agent/internal/redis"
	"fmt"
	"github.com/spf13/viper"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	viper.SetConfigName(".env.json")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	helpers.StoreContainerId(strconv.Itoa(rand.Int()))
}

func main() {
	var miniTasksChan = make(chan redis.MiniTask, 1000)

	// Send ping signals to main server
	go func() {
		containerId, err := helpers.GetContainerId()
		fmt.Println(err)
		for {
			rdb, err := redis.GetConnection()
			if err == nil {
				rdb.Publish(context.Background(), "ping_channel", containerId)
			}

			time.Sleep(5 * time.Second)
		}
	}()

	// Check Queue for new tasks
	go func() {
		for {
			time.Sleep(5 * time.Millisecond)
			rdb, err := redis.GetConnection()
			if err == nil {
				queue := rdb.RPop(context.Background(), "mini_task_queue")

				if queue.Err() != nil {
					continue
				}
				miniTask := redis.MiniTask{}
				err = json.Unmarshal([]byte(queue.Val()), &miniTask)
				fmt.Println(queue.String(), &miniTask, queue.Val())
				if err == nil {
					miniTasksChan <- miniTask
				}
			}
		}
	}()

	//Resolve tasks and send answer to channel
	go func() {
		for {
			select {
			case miniTask, ok := <-miniTasksChan:
				if ok {
					_, _ = miniTask.ResolveMiniTask()
					redis.PublishMiniTaskAnswer(&miniTask)
					fmt.Println(miniTask)
				}
			}
		}
	}()

	// Sleep infinity
	select {}
}
