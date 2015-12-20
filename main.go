package main

import (
	"encoding/json"
	"github.com/antoinefinkelstein/intercom-twitter-follow/Godeps/_workspace/src/github.com/ChimeraCoder/anaconda"
	"github.com/antoinefinkelstein/intercom-twitter-follow/Godeps/_workspace/src/gopkg.in/intercom/intercom-go.v1"
	redix "github.com/antoinefinkelstein/intercom-twitter-follow/Godeps/_workspace/src/menteslibres.net/gosexy/redis"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	redis           *redix.Client
	twitterAPI      = anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
	intercomAPI     = intercom.NewClient(os.Getenv("INTERCOM_APP_ID"), os.Getenv("INTERCOM_API_KEY"))
	redisConnection = &redisInfos{path: os.Getenv("REDIS_URL")}
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Success")
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var f interface{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&f)
	if err != nil {
		return
	}

	log.Println(f)

	response := f.(map[string]interface{})
	if response["type"] != "user" {
		return
	}

	go enqueueUserID(response["id"].(string))

	io.WriteString(w, "Webhook processed")
}

func enqueueUserID(id string) {
	log.Println("Received id " + id)

	redis.ZAdd("queue:users", time.Now().Unix(), id)
	return
}

func main() {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))

	if redisConnection.path == "" {
		redisConnection.path = "redis://127.0.0.1:6379"
	}
	redisConnection.parse()

	redis = redix.New()
	err := redis.Connect(redisConnection.host, redisConnection.port)
	if err != nil {
		log.Fatalf("Connect failed: %s\n", err.Error())
		return
	}
	// if redisConnection.auth != "" {
	_, err = redis.Auth(redisConnection.auth)
	if err != nil {
		log.Fatalf("Connect failed: %s\n", err.Error())
		return
	}
	// }

	go startWorkers()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/webhook", webhookHandler)
	http.ListenAndServe(":8000", nil)
}
