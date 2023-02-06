# NASA NeoWs API Go SDK
This is a Golang client for the [NASA NeoWs API](https://api.nasa.gov/#asteroids-neows). The API is maintained by the [SpaceRocks](https://github.com/SpaceRocks/) team (David Greenfield, Arezu Sarvestani, Jason English and Peter Baunach).

I am not affiliated with the aforementioned groups and this is not an official initiative from them. This is also a work in progress (WIP) and it is not production-ready.

## Usage

```bash
go get https://github.com/viniciusrtf/nasa-neows-go
```

Setup your client
```go
client := neows.NewClient("your-api-key")
```

Optionally, enable caching NeoWs API responses
```go
client.EnableCache()
```

Get approaching asteroids using default options (next 7 days)
```go
feed, err := client.Feed.Fetch(nil)
```

Get approaching asteroids with detailed orbital data
```go
opts := &FeedOptions{Detailed: true}
feed, err := client.Feed.Fetch(opts)
```

Get approaching asteroids from "2021-06-01" to "2021-06-02" (YYYY-mm-dd)
```go
startDate, err := time.Parse("2006-01-02", "2021-06-01")
if err != nil {
    t.Errorf("Error parsing start date: %s", err)
}

endDate, err := time.Parse("2006-01-02", "2021-06-02")
if err != nil {
    t.Errorf("Error parsing end date: %s", err)
}

opts := &FeedOptions{startDate, endDate}
feed, err := client.Feed.Fetch(opts)
```

A more complete example: Get approaching asteroids for a given date range:
```go
package main

import (
    "fmt"
    "github.com/viniciusrtf/nasa-neows-go"
)

func main() {
    // See https://api.nasa.gov/#signUp to generate an API Key
    apiKey := "your-api-key"
    client := neows.NewClient(&apiKey)
    client.EnableCache() // optional

    // Get a feed with close approaches by date range
    opts := &FeedOptions{
        StartDate: time.Now(),
        EndDate: time.Now().AddDate(0, 0, 1)
    }
    feed, err := client.Feed.Fetch(opts)
    if err != nil {
        panic(err)
    }

    if feed.ElementCount == 0 {
        fmt.Println("No close approaches for the given date range")
    }

    // Print asteroids info by date
    for date, neo := range feed.NeosByDate {
        fmt.Printf("Asteroids approaching in %s\n", date)
        for _, n := range neo {
            fmt.Printf("\tName: %s\n", n.Name)
            fmt.Printf("\tNASA JPL URL: %s\n", n.NasaJplURL)
            fmt.Printf("\tPotentialy hazardous: %t\n", n.IsPotentiallyHazardousAsteroid)
            fmt.Printf("\tMiss distance: %s km\n\n", n.CloseApproachData[0].MissDistance.Kilometers)
        }
        fmt.Printf("\n")
    }
}
```

## Related links and projects
- [NASA Open APIs](https://api.nasa.gov)
- [SpaceRocks GitHub](https://github.com/SpaceRocks/)
- [NeoWs API Swagger](http://www.neowsapp.com/swagger-ui/index.html)

