package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/appleboy/go-fcm"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var deviceTokens []string
var count int

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("API_KEY")

	// Create a FCM client to send the message.
	client, err := fcm.NewClient(apiKey)
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hero Zawrudo!")
	})

	r.POST("/send/all", func(c *gin.Context) {
		for _, deviceToken := range deviceTokens {
			sendMessage(client, deviceToken)
		}
	})

	r.POST("/send", func(c *gin.Context) {
		type fcmData struct {
			Data struct {
				Count string `json:"count"`
				Via   string `json:"via"`
			} `json:"data"`
			Notification struct {
				Body  string `json:"body"`
				Title string `json:"title"`
			} `json:"notification"`
			Token string `json:"token"`
		}
		data := new(fcmData)
		err := c.BindJSON(data)
		if err != nil {
			log.Fatalln(err)
		}
		deviceToken := data.Token
		deviceTokens = append(deviceTokens, deviceToken)
		log.Println(data.Token)
		c.JSON(http.StatusOK, gin.H{
			"tokens": deviceTokens,
			"data":   data,
		})
		sendMessage(client, deviceToken)
	})

	http.ListenAndServe(":3100", r)
}

func noop(x interface{}) {}

func sendMessage(client *fcm.Client, token string) {
	count++
	// Time to live is 1 hour
	// https://firebase.google.com/docs/cloud-messaging/concept-options#ttl
	// time_to_live in https://firebase.google.com/docs/cloud-messaging/http-server-ref#downstream
	var TTL uint = 1 * 60 * 60
	// Create the message to be sent.
	msg := &fcm.Message{
		To:       token,
		Priority: "high",
		Data: map[string]interface{}{
			"foo": "bar",
		},
		Notification: &fcm.Notification{
			Title: "Hello World " + strconv.Itoa(count),
			Body:  "Hello? World?" + strconv.Itoa(count),
		},
		TimeToLive: &TTL,
	}
	// Send the message and receive the response without retries.
	response, err := client.Send(msg)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", response)
}
