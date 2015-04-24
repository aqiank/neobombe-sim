package twitter

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

var settings = struct {
	ConsumerKey string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	AccessToken string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
	Track string `json:"track"`
}{}

var twitter *anaconda.TwitterApi

func load() {
	var file *os.File
	var err error

	if file, err = os.Open("twitter_settings.json"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = json.NewDecoder(file).Decode(&settings); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func connect() (bool, error) {
	anaconda.SetConsumerKey(settings.ConsumerKey)
	anaconda.SetConsumerSecret(settings.ConsumerSecret)
	twitter = anaconda.NewTwitterApi(settings.AccessToken, settings.AccessTokenSecret)

	ok, err := twitter.VerifyCredentials()
	if err != nil {
		return false, fmt.Errorf("connect: %v", err)
	}

	return ok, nil
}

func listen(track string, msgChan, decChan chan string) {
	vals := url.Values{}
	vals.Set("track", track)
	stream := twitter.PublicStreamFilter(vals)
	fmt.Println("twitter: tracking", track)

	for {
		data := <-stream.C
		switch t := data.(type) {
		case anaconda.Tweet:
			msgChan <-t.Text
			<-decChan
			flush(stream)
		}
	}
}

func flush(stream anaconda.Stream) {
	c := make(chan bool, 1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		c <- true
	}()

	go func() {
		for range stream.C {}
	}()

	select {
	case <-c:
		return
	}
}

func Run(cs map[string]interface{}) {
	msgChan := cs["message"].(chan string)
	decChan := cs["decrypted"].(chan string)
	sigChan := cs["signal"].(chan os.Signal)

	load()

	fmt.Println("twitter starting..")
	if ok, err := connect(); err != nil {
		fmt.Println("twitter:", err)
		os.Exit(1)
	} else if ok {
		fmt.Println("twitter: successfully connected")
		go listen(settings.Track, msgChan, decChan)
	} else {
		fmt.Println("twitter: failed to verify credential")
		os.Exit(1)
	}

	sigChan <- <-sigChan
	fmt.Println("twitter exiting..")
}

func Track() string {
	return settings.Track
}
