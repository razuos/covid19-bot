package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/razuos/covid19-bot/config"
	"github.com/razuos/covid19-bot/gatherers"
	"github.com/razuos/covid19-bot/twitter"
	"github.com/razuos/covid19-bot/utils"
)

func main() {

	config.Load()

	if len(os.Args) > 1 {
		if os.Args[1] == "--getTokens" {
			utils.StartAuthorizationFlow()
			os.Exit(0)
		}
	}

	var wg sync.WaitGroup
	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Starting COVID-19 bot. Happy tweeting :)")

	go func() {
		sig := <-signals
		log.Println("Signal received:", sig)
		done <- true
	}()

	log.Printf("Tweet interval: %d, Gathering interval: %d\n", config.AppConfig.Twitter.TweetIntervalMinutes, config.AppConfig.Gatherer.UpdateIntervalMinutes)
	// tweetClock := time.NewTicker(5 * time.Second)
	// gatherClock := time.NewTicker(5 * time.Second)
	tweetClock := time.NewTicker(time.Duration(config.AppConfig.Twitter.TweetIntervalMinutes) * time.Minute)
	gatherClock := time.NewTicker(time.Duration(config.AppConfig.Gatherer.UpdateIntervalMinutes) * time.Minute)
	statusChan := make(chan string, 5)

	go twitter.Worker(&wg, statusChan, tweetClock, done)
	go gatherers.RSSWorker("https://www.riogrande.rs.gov.br/corona/feed/", &wg, statusChan, gatherClock, done)

	<-done
	tweetClock.Stop()
	gatherClock.Stop()
	wg.Wait()
	log.Println("Exiting.")

}
