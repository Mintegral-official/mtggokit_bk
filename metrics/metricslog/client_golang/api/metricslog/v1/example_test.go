// Copyright 2019 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package v1_test provides examples making requests to Prometheus using the
// Golang client.
package v1_test

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/Mintegral-official/mtggokit/metrics/metricslog/client_golang/api"
    v1 "github.com/Mintegral-official/mtggokit/metrics/metricslog/client_golang/api/metricslog/v1"
)

func ExampleAPI_Query() {
    client, err := api.NewClient(api.Config{
        Address: "http://demo.robustperception.io:9090",
    })
    if err != nil {
        fmt.Printf("Error creating client: %v\n", err)
        os.Exit(1)
    }

    api := v1.NewAPI(client)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    result, warnings, err := api.Query(ctx, "up", time.Now())
    if err != nil {
        fmt.Printf("Error querying Prometheus: %v\n", err)
        os.Exit(1)
    }
    if len(warnings) > 0 {
        fmt.Printf("Warnings: %v\n", warnings)
    }
    fmt.Printf("Result:\n%v\n", result)
}

func ExampleAPI_QueryRange() {
    client, err := api.NewClient(api.Config{
        Address: "http://demo.robustperception.io:9090",
    })
    if err != nil {
        fmt.Printf("Error creating client: %v\n", err)
        os.Exit(1)
    }

    api := v1.NewAPI(client)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    r := v1.Range{
        Start: time.Now().Add(-time.Hour),
        End:   time.Now(),
        Step:  time.Minute,
    }
    result, warnings, err := api.QueryRange(ctx, "rate(metricslog_tsdb_head_samples_appended_total[5m])", r)
    if err != nil {
        fmt.Printf("Error querying Prometheus: %v\n", err)
        os.Exit(1)
    }
    if len(warnings) > 0 {
        fmt.Printf("Warnings: %v\n", warnings)
    }
    fmt.Printf("Result:\n%v\n", result)
}

func ExampleAPI_Series() {
    client, err := api.NewClient(api.Config{
        Address: "http://demo.robustperception.io:9090",
    })
    if err != nil {
        fmt.Printf("Error creating client: %v\n", err)
        os.Exit(1)
    }

    api := v1.NewAPI(client)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    lbls, warnings, err := api.Series(ctx, []string{
        "{__name__=~\"scrape_.+\",job=\"node\"}",
        "{__name__=~\"scrape_.+\",job=\"metricslog\"}",
    }, time.Now().Add(-time.Hour), time.Now())
    if err != nil {
        fmt.Printf("Error querying Prometheus: %v\n", err)
        os.Exit(1)
    }
    if len(warnings) > 0 {
        fmt.Printf("Warnings: %v\n", warnings)
    }
    fmt.Println("Result:")
    for _, lbl := range lbls {
        fmt.Println(lbl)
    }
}
