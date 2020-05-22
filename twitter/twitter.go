package twitter

import (
	"log"
	"os"
	"strconv"
	"time"

	twitter "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func NewTwitter() (*twitter.Client, twitter.User) {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessSecret := os.Getenv("ACCESS_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	if consumerKey == "" || consumerSecret == "" || accessSecret == "" || accessToken == "" {
		log.Fatal("Fail to get env vars")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)
	user, _, err := client.Accounts.VerifyCredentials(nil)
	if err != nil {
		log.Fatalf("Error authenticating user %s: %s", user.Name, err.Error())
	}
	return client, *user
}

func UpdateStatus(client *twitter.Client, status string) {
	lat, err := strconv.ParseFloat(os.Getenv("CITY_LAT"), 64)
	long, err := strconv.ParseFloat(os.Getenv("CITY_LONG"), 64)

	var params twitter.StatusUpdateParams

	if err == nil {
		params = twitter.StatusUpdateParams{
			Lat:  &lat,
			Long: &long,
		}
	}

	client.Statuses.Update(status, &params)
}

func TwitterWorker(status <-chan string, clock *time.Ticker) {

	quit := make(chan struct{})
	client, user := NewTwitter()
	log.Printf("Starting Twitter Worker for user %s.", user.Name)
	for {
		select {
		case <-clock.C:
			text := <-status
			log.Println("Tweeting.")
			UpdateStatus(client, text)
		case <-quit:
			clock.Stop()
			return
		}
	}
}
