package main

import (
  "fmt"
  "net/http"
  "time"
  "sync"
  "sort"
  "os"
  "strings"
  "gopkg.in/alecthomas/kingpin.v2"
)

var (
  version string
  response_times []int
  requests_count int
  wg sync.WaitGroup
  mutex sync.Mutex
  url = kingpin.Flag("url", "URL to benchmark").Required().Short('u').String()
  method = kingpin.Flag("method", "HTTP method").Default("GET").Short('m').String()
  size = kingpin.Flag("number", "Number of requests to perform").Default("100").Short('n').Int()
  threads_count = kingpin.Flag("threads", "Number of threads").Default("1").Short('t').Int()
  http_content_type = kingpin.Flag("content-type", "HTTP content type").Default("").String()
  http_body = kingpin.Flag("body", "Body of HTTP request").Default("").String()
)


func main() {
  kingpin.Version(version)
  kingpin.Parse()
  if *size < *threads_count {
    fmt.Printf("Requests count should be greater than threads count.\n")
    os.Exit(1)
  }
  requests_count = 0
  go showProgress()
  fmt.Printf("Starting %v threads to make %v %v requests to %s\n", *threads_count, *size, *method, *url)
  for i := 0; i < *threads_count; i++ {
    go RequestsThread(*method, *url)
  }
  time.Sleep(time.Second)
  wg.Wait()
  ShowResults()
}

func RequestsThread(method string, host string) {
  wg.Add(1)
  for requests_count < *size {
    requests_count++
    MakeRequest(method, host)
  }
  wg.Done()
}

func MakeRequest(method string, host string) {
  start_time := time.Now()
  req, _ := http.NewRequest(method, host, strings.NewReader(*http_body))
  if *http_content_type != "" {
    req.Header.Add("Content-Type", *http_content_type)
  }
  client := &http.Client{}
  resp, _ := client.Do(req)
  mutex.Lock()
  if resp != nil {
    defer resp.Body.Close()
    response_times = append(response_times, int(time.Since(start_time).Nanoseconds() / 1000000))
  }
  mutex.Unlock()
}

func ShowResults() {
  if *size == failedRequestsCount() {
    fmt.Printf("\nAll requests failed!\n")
    os.Exit(1)
    return
  }
  fmt.Printf("%g requests/sec\n", meanRequestsPerSec())
  fmt.Printf("%g ms mean response time\n", meanResponseTime())
  sort.Ints(response_times)
  percentiles := []float32{0.25, 0.5, 0.75, 0.90, 0.95, 0.98}
  fmt.Printf("Percentage of requests processed within a certain time:\n")
  for i := 0; i < len(percentiles); i++ {
    percentile_index := int(float32(len(response_times)) * percentiles[i])
    fmt.Printf("%v%%: %v ms\n", int(percentiles[i] * 100), response_times[percentile_index])
  }
  fmt.Printf("%d requests total.\n", *size)
  fmt.Printf("%d requests failed.\n\n", failedRequestsCount())
}

func meanResponseTime() float32 {
  return float32(sumTime()) / float32(len(response_times))
}

func meanRequestsPerSec() float32 {
  return 1000.0 / meanResponseTime()
}

func sumTime() int32 {
  total := 0
  for i := 0; i < len(response_times); i++ {
    total += response_times[i]
  }
  return int32(total)
}

func showProgress() {
  for {
    time.Sleep(5 * time.Second)
    if failedRequestsCount() + len(response_times) < *size {
      percent := int32((float32(failedRequestsCount() + len(response_times)) / float32(*size)) * 100.0)
      fmt.Printf("%d%% done...\n", percent)
    } else {
      return
    }
  }
}

func failedRequestsCount() int {
  return *size - len(response_times)
}
