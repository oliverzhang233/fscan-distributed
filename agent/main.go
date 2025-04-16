package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/streadway/amqp"
)

// 消息队列配置
const (
	amqpURL = "amqp://guest:guest@localhost:5672/"
	taskQueue = "scan_tasks"
	resultQueue = "scan_results"
)

// 初始化消息队列连接
func initRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	return conn, ch, nil
}

// 执行 Fscan 扫描
func executeScan(task map[string]interface{}) (string, error) {
	host, ok := task["host"].(string)
	if!ok {
		return "", fmt.Errorf("invalid host in task")
	}
	cmd := exec.Command("fscan", "-h", host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// 任务处理回调
func handleTask(d amqp.Delivery, ch *amqp.Channel) {
	var task map[string]interface{}
	err := json.Unmarshal(d.Body, &task)
	if err != nil {
		log.Printf("Failed to unmarshal task: %v", err)
		return
	}
	log.Printf("Received task: %v", task)
	result, err := executeScan(task)
	if err != nil {
		log.Printf("Failed to execute scan: %v", err)
		return
	}
	resultData := map[string]interface{}{
		"host":   task["host"],
		"result": result,
	}
	body, err := json.Marshal(resultData)
	if err != nil {
		log.Printf("Failed to marshal result: %v", err)
		return
	}
	err = ch.Publish(
		"",           // exchange
		resultQueue,  // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Printf("Failed to publish result: %v", err)
		return
	}
	log.Printf("Result sent successfully")
}

func main() {
	conn, ch, err := initRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		taskQueue, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Waiting for tasks...")
	for d := range msgs {
		handleTask(d, ch)
	}
}    
