package gatherers

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/razuos/covid19-bot/db"
	"github.com/razuos/covid19-bot/models"
	"github.com/razuos/covid19-bot/utils"
	"github.com/ungerik/go-rss"
)

func newChannel(url string) (rss.Channel, error) {
	resp, err := rss.Read(url, false)
	if err != nil {
		log.Fatal(err.Error())
	}
	channel, err := rss.Regular(resp)
	return *channel, err
}

func getEntriesFromChannel(channel rss.Channel) (results []rss.Item) {
	return channel.Item
}

func getSameDayEntries(entries []rss.Item) (results []rss.Item) {
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

func hasAnyKeywords(str string) bool {
	keywords := []string{"covid", "coronavirus", "corona", "casos", "atualizacao", "pandemia"}
	for _, keyword := range keywords {
		if strings.Contains(str, keyword) {
			return true
		}
	}
	return false
}

func filterEntriesByRelevancy(entries []rss.Item) (results []rss.Item) {

	for _, item := range entries {
		title := utils.RemoveSpecialChars(item.Title)
		body := utils.RemoveSpecialChars(item.Description)

		if hasAnyKeywords(title) || hasAnyKeywords(body) {
			log.Println("Found Relevant Item: ", item.GUID)
			results = append(results, item)
		}
	}
	return results
}

func addDotToTitle(str string) string {
	if string(str[len(str)-1:]) != "." {
		return str + string(".")
	}
	return str
}

func getFilteredEntries(channel rss.Channel) []rss.Item {
	getEntriesFromChannel(channel)
	filtered := getSameDayEntries(channel.Item)
	filtered = filterEntriesByRelevancy(filtered)
	return filtered
}

func genUUID(url string) (result string) {
	result = ""
	str1 := utils.Reverse(url)
	for _, char := range str1 {
		if string(char) != "=" {
			result = result + string(char)
		} else {
			return utils.Reverse(result)
		}
	}
	return
}

func isEntryInDB(db *gorm.DB, uuid string) bool {
	find := models.RSSData{}
	if err := db.First(&find, "uuid = ?", uuid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		log.Fatalf("Error in isEntryInDB: %s", err.Error())
	}
	return true
}

// RSSWorker gathers informartion from a RSS feed and sends to the twitter worker.
func RSSWorker(url string, wg *sync.WaitGroup, status chan<- string, clock *time.Ticker, done <-chan bool) {
	db := db.Connect()
	log.Printf("Starting RSS Gathering Worker on %s", url)

	for {
		select {
		case <-clock.C:
			log.Println("Updating RSS channel.")
			channel, _ := newChannel(url)
			entries := getFilteredEntries(channel)
			for _, item := range entries {
				uuid := genUUID(item.GUID)
				if !isEntryInDB(db, uuid) {
					title := addDotToTitle(item.Title)
					time, _ := item.PubDate.Parse()
					entryTime := time.Format("2006-01-02")
					log.Printf("New tweet from UUID: %s\n", uuid)
					text := fmt.Sprintf("%s\r\rFonte: %s", title, item.GUID)
					status <- text
					db.Create(&models.RSSData{UUID: uuid, Title: title, Source: item.GUID, PubDate: entryTime})
				} else {
					log.Printf("Ignoring UUID: %s", uuid)
				}
			}
		case <-done:
			wg.Done()
			db.Close()
			return
		}
	}
}
