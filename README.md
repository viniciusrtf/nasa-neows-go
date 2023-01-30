# NASA NeoWs API Go SDK
This is a Golang client for the [NASA NeoWs API](https://api.nasa.gov/#asteroids-neows). The API is maintained by the [SpaceRocks](https://github.com/SpaceRocks/) team (David Greenfield, Arezu Sarvestani, Jason English and Peter Baunach).

I am not affiliated with the aforementioned groups and this is not an official initiative from them. This is also a work in progress (WIP) and it is not production-ready.

## Usage

```bash
go get https://github.com/viniciusrtf/nasa-neows-go
```

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
    startDate := time.Now()
    endDate := time.Now().AddDate(0, 0 ,1)
    feed, err := client.Feed.Fetch(startDate, endDate)
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

