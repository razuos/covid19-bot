package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/covid19-bot/covid19-bot/gatherers"
	"github.com/covid19-bot/covid19-bot/twitter"
	"github.com/covid19-bot/covid19-bot/utils"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "--getTokens" {
			utils.CheckEnv([]string{"CONSUMER_KEY", "CONSUMER_SECRET"})
			utils.StartAuthorizationFlow()
			os.Exit(0)
		}
	}

	utils.CheckEnv([]string{"CONSUMER_KEY", "CONSUMER_SECRET", "ACCESS_TOKEN", "ACCESS_SECRET"})

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting COVID-19 bot. Happy tweeting :)")
		sig := <-signals
		log.Println("Signal received:", sig)
		done <- true
	}()

	tweetClock := time.NewTicker(1 * time.Minute)
	gatherClock := time.NewTicker(30 * time.Second)
	statusChan := make(chan string, 5)

	go twitter.TwitterWorker(statusChan, tweetClock)
	go gatherers.RSSGathererWorker(statusChan, gatherClock)

	<-done
	log.Println("Exiting.")

}
