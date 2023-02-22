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

Lookup an asteroid by its SPKID or Asteroid Name
```go
neo, err := client.Lookup.Find("2023 BU")
```

Get approaching asteroids for the next 7 days
```go
feed, err := client.Feed.Fetch(time.Now(), time.Now().AddDate(0, 0, 7))
```

Get today's approaching asteroids
```go
feed, err := client.Feed.Today()
```

Get **detailed** information about approaching asteroids from "2021-06-01" to "2021-06-02" (YYYY-mm-dd)
```go
start, err := time.Parse("2006-01-02", "2021-06-01")
if err != nil {
    t.Errorf("Error parsing start date: %s", err)
}

end, err := time.Parse("2006-01-02", "2021-06-02")
if err != nil {
    t.Errorf("Error parsing end date: %s", err)
}

fs := NewFeedService(client, &FeedOptions{Detailed: true})
feed, err := fs.Fetch(start, end)
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
    start, err := time.Parse("2006-01-02", "2021-06-01")
    if err != nil {
        t.Errorf("Error parsing start date: %s", err)
    }

    end, err := time.Parse("2006-01-02", "2021-06-02")
    if err != nil {
        t.Errorf("Error parsing end date: %s", err)
    }

    feed, err := client.Feed.Fetch()
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

