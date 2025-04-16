package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

// 任务下发处理函数
func sendTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var task map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	conn, ch, err := initRabbitMQ()
	if err != nil {
		http.Error(w, "Failed to connect to RabbitMQ", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	defer ch.Close()

	body, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Failed to marshal task", http.StatusInternalServerError)
		return
	}
	err = ch.Publish(
		"",           // exchange
		taskQueue,    // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		http.Error(w, "Failed to publish task", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Task sent successfully")
}

// 结果收集处理函数
func receiveResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var result map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// 这里可以将结果存储到数据库
	log.Printf("Received result: %v", result)
	fmt.Fprintf(w, "Result received successfully")
}

func main() {
	http.HandleFunc("/send_task", sendTaskHandler)
	http.HandleFunc("/receive_result", receiveResultHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}    
