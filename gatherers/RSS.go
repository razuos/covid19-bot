package gatherers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"unicode"

	"github.com/covid19-bot/covid19-bot/utils"
	"github.com/ungerik/go-rss"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func NewChannel(url string) (rss.Channel, error) {
	resp, err := rss.Read(url, false)
	if err != nil {
		log.Fatal(err.Error())
	}
	channel, err := rss.Regular(resp)
	return *channel, err
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func processString(str string) string {
	transformer := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	result := strings.ToLower(str)
	result, _, _ = transform.String(transformer, result)

	return result
}

func GetEntriesFromChannel(channel rss.Channel) (results []rss.Item) {
	return channel.Item
}

func GetSameDayEntries(entries []rss.Item) (results []rss.Item) {
	log.Println("Getting RSS entries from same day")
	localtime := time.Now().Local().Format("2006-01-02")
	log.Printf("Local time: %s", localtime)

	for i, item := range entries {
		time, err := item.PubDate.Parse()
		if err != nil {
			log.Println("Failed to get pub time for entry #", i)
		}
		entryTime := time.Format("2006-01-02")
		// log.Println("Entry time: ", entryTime)
		// log.Println("Entry ID: ", item.GUID)
		if localtime == entryTime {
			results = append(results, item)
		}
	}

	log.Printf("Got %d entries from the same day.", len(results))

	return results
}

func HasAnyKeywords(str string) bool {
	keywords := []string{"covid", "coronavirus", "corona", "casos", "atualizacao", "pandemia"}
	for _, keyword := range keywords {
		if strings.Contains(str, keyword) {
			return true
		}
	}
	return false
}

func FilterEntriesByRelevancy(entries []rss.Item) (results []rss.Item) {

	for _, item := range entries {
		title := processString(item.Title)
		body := processString(item.Description)

		// log.Println("Item title: ", title)
		// log.Println("Item body: ", body)

		if HasAnyKeywords(title) || HasAnyKeywords(body) {
			log.Println("Found Relevant Item: ", item.GUID)
			results = append(results, item)
		}
	}
	return results
}

func GetFilteredEntries(channel rss.Channel) []rss.Item {
	GetEntriesFromChannel(channel)
	filtered := GetSameDayEntries(channel.Item)
	filtered = FilterEntriesByRelevancy(filtered)
	return filtered
}

func RSSGathererWorker(status chan<- string, clock *time.Ticker) {
	log.Printf("Starting RSS Gathering Worker.")
	quit := make(chan struct{})
	var knownEntries []string
	for {
		select {
		case <-clock.C:
			log.Println("Updating RSS channel.")
			log.Println("Known entries: ", knownEntries)
			channel, _ := NewChannel("http://www.riogrande.rs.gov.br/corona/feed")
			entries := GetFilteredEntries(channel)
			log.Println("Updating RSS channel.")
			for _, item := range entries {
				if !utils.IsStringInSlice(item.GUID, knownEntries) {
					log.Println("Adding new tweet to queue:", item.GUID)
					knownEntries = append(knownEntries, item.GUID)
					text := fmt.Sprintf("%s\r\rFonte: %s", item.Title, item.GUID)
					status <- text
				} else {
					log.Println("Item already tweetted about:", item.GUID)
				}
			}

		case <-quit:
			clock.Stop()
			return
		}
	}
}
