//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup


func producer(stream Stream, tweets chan *Tweet) {
	defer wg.Done()
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			break
		}
		tweets <- tweet
	}
	close(tweets)

}

func consumer(tweets chan *Tweet) {
	defer wg.Done()
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	wg.Add(2)

	tweetsChan := make(chan *Tweet)
	start := time.Now()
	stream := GetMockStream()
	go producer(stream, tweetsChan)
	// Producer
	go consumer(tweetsChan)
	fmt.Printf("Process took %s\n", time.Since(start))
	wg.Wait()

}
