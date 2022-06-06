package main

import (
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// UpperCase letterBytes for symbols
const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Create a cache with a default expiration time of 5 minutes, and which purges expired items every 10 minutes
var c = cache.New(5*time.Minute, 10*time.Minute)

// Generate Random String with given Lenght
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Random Integer between given range
func randomNumber(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// Random Float between given range
func randomFloat(min float64, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

// Handler for default route
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	// Only Entertain Default route for the time being
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Only Entertain GET method for the time being
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// Handling request here
	type Tick struct {
		Time   time.Time `json:"time"`
		Symbol string    `json:"symbol"`
		Open   float64   `json:"open"`
		High   float64   `json:"high"`
		Low    float64   `json:"low"`
		Close  float64   `json:"close"`
		Volume int       `json:"volume"`
	}

	// Generating 10 random Ticks and storing it in the Cache
	for i := 0; i < 10; i++ {
		tick := Tick{time.Now(), RandStringBytes(4), 100.00, 100.00, 100.00, 100.00, 10000}

		// converting index to string to store it as key
		c.Set(strconv.Itoa(i), tick, cache.DefaultExpiration)

	}

	// Ticker for every 100ms which will publish a stock tick
	ticker := time.NewTicker(100 * time.Millisecond)

	for _ = range ticker.C {
		var randomIndex = randomNumber(0, 9)

		// Get the string associated with the key "i" from the cache
		foo, found := c.Get(strconv.Itoa(randomIndex))
		if found {
			// type assert it to a compatible type before accessing its fields
			foo := foo.(Tick)
			foo.Time = time.Now()
			foo.Volume = foo.Volume + randomNumber(0, 1000)
			foo.Close = randomFloat(foo.Close*0.9, foo.Close*1.1)

			if foo.Close > foo.High {
				foo.High = foo.Close
			} else if foo.Close < foo.Low {
				foo.Low = foo.Close
			}

			// updating the tick.
			c.Set(strconv.Itoa(randomIndex), foo, cache.DefaultExpiration)

			// Publishing Tick
			jsonByte, _ := json.Marshal(foo)
			// Writing it to Browser
			fmt.Fprintf(w, string(jsonByte))
			// Printing it to the console
			fmt.Print(string(jsonByte), "\n")
		}
	}

}

func main() {

	// Default route Handler
	http.HandleFunc("/", defaultHandler)

	// Server initialisation at port 9000
	fmt.Printf("Starting server at port 9000\n")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}

}
