package config

// TwitterConfig holds twitter related configs.
type TwitterConfig struct {
	ConsumerKey          string
	ConsumerSecret       string
	AccessToken          string
	AccessTokenSecret    string
	Lat                  float64
	Long                 float64
	TweetIntervalMinutes int
}

// GathererConfig hold gathering related configs.
type GathererConfig struct {
	UpdateIntervalMinutes int
}

// DBConfig holds DB related configs.
type DBConfig struct {
	Host string
	Port int
	User string
	Pass string
	Name string
}

// Config holds all the config categories.
type Config struct {
	Twitter  TwitterConfig
	Gatherer GathererConfig
	DB       DBConfig
}
