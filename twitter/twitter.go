package twitter

import (
	"log"
	"sync"
	"time"

	twitter "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/razuos/covid19-bot/config"
)

func newTwitter() (*twitter.Client, twitter.User) {
	config.Check([]string{"consumerkey", "consumersecret", "accesstoken", "accesstokensecret"})

	consumerKey := config.AppConfig.Twitter.ConsumerKey
	consumerSecret := config.AppConfig.Twitter.ConsumerSecret
	accessToken := config.AppConfig.Twitter.AccessToken
	accessSecret := config.AppConfig.Twitter.AccessTokenSecret

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

func updateStatus(client *twitter.Client, status string) {
	var params twitter.StatusUpdateParams

	if config.Check([]string{"lat", "long"}) {
		params = twitter.StatusUpdateParams{
			Lat:  &config.AppConfig.Twitter.Lat,
			Long: &config.AppConfig.Twitter.Long,
		}
	}

	client.Statuses.Update(status, &params)
}

// Worker is a worker function to listen for statuses and post them.
func Worker(wg *sync.WaitGroup, status <-chan string, clock *time.Ticker, done <-chan bool) {
	defer wg.Done()

	client, user := newTwitter()
	log.Printf("Starting Twitter Worker for user %s.", user.Name)
	for {
		select {
		case <-clock.C:
			log.Println("Waiting for status to tweet.")
			text := <-status
			log.Println("Tweeting.")
			updateStatus(client, text)
		case <-done:
			wg.Done()
			return
		}
	}
}
