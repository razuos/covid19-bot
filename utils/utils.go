package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/dghubble/oauth1"
	twauth "github.com/dghubble/oauth1/twitter"
)

var config oauth1.Config

func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func CheckEnv(variables []string) {
	for _, variable := range variables {
		if os.Getenv(variable) == "" {
			log.Fatal("Missing required environment variable: ", variable)
		}
	}
}

func StartAuthorizationFlow() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	if consumerKey == "" || consumerSecret == "" {
		log.Fatal("Env vars missing.")
	}

	config = oauth1.Config{
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
