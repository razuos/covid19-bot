package utils

import (
	"fmt"
	"log"

	"github.com/dghubble/oauth1"
	twauth "github.com/dghubble/oauth1/twitter"
	"github.com/razuos/covid19-bot/config"
)

// IsStringInSlice checks if a string is in a slice.
func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// StartAuthorizationFlow returns the api keys for twitter.
func StartAuthorizationFlow() {
	config.Check([]string{"consumerkey, consumersecret"})

	consumerKey := config.AppConfig.Twitter.ConsumerKey
	consumerSecret := config.AppConfig.Twitter.ConsumerSecret

	config := oauth1.Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		Endpoint:       twauth.AuthorizeEndpoint,
	}

	requestToken, _, err := config.RequestToken()
	if err != nil {
		log.Fatalf("Error getting request token: %s", err.Error())
	}
	fmt.Println("Request Token: ", requestToken)

	authorizationURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		log.Fatalf("Error getting authorization URL: %s", err.Error())
	}
	fmt.Println("Authorization URL: ", authorizationURL)

	fmt.Printf("Paste your PIN here: ")
	var verifier string
	_, err = fmt.Scanf("%s", &verifier)
	if err != nil {
		log.Fatalf("Error getting PIN: %s", err.Error())
	}

	accessToken, accessSecret, err := config.AccessToken(requestToken, "secret does not matter", verifier)
	if err != nil {
		log.Fatalf("Error getting access token: %s", err.Error())
	}

	authToken := oauth1.NewToken(accessToken, accessSecret)
	if err != nil {
		log.Fatalf("Error getting auth token: %s", err.Error())
	}

	fmt.Println("Authorization successful.")
	fmt.Printf("token: %s\nsecret: %s\n", authToken.Token, authToken.TokenSecret)
}
