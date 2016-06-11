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
  timeout = kingpin.Flag("timeout", "Timeout for requests (sec)").Default("10").Int()
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
  fmt.Printf("Starting test %s\n", *url)
  fmt.Printf("Requests: %v\n", *size)
  fmt.Printf("Threads: %v\n", *threads_count)
  fmt.Printf("Timeout: %v sec\n\n", *timeout)
  for i := 0; i < *threads_count; i++ {
    wg.Add(1)
    go RequestsThread(*method, *url)
  }
  wg.Wait()
  ShowResults()
}

func RequestsThread(method string, host string) {
  var thread_response_times []int

  for ; isNeedRequest(); {
    MakeRequest(method, host, &thread_response_times)
  }

  mutex.Lock()
  response_times = append(response_times, thread_response_times...)
  mutex.Unlock()

  wg.Done()
}

func isNeedRequest() bool {
  mutex.Lock()
  defer mutex.Unlock()
  if (requests_count < *size) {
    requests_count++
    return true
  } else {
    return false
  }
}

func MakeRequest(method string, host string, thread_response_times *[]int) {
  start_time := time.Now()
  req, _ := http.NewRequest(method, host, strings.NewReader(*http_body))
  if *http_content_type != "" {
    req.Header.Add("Content-Type", *http_content_type)
  }
  client := &http.Client{
    Timeout: time.Duration(*timeout) * time.Second,
  }
  resp, _ := client.Do(req)
  if resp != nil {
    defer resp.Body.Close()
    *thread_response_times = append(*thread_response_times, int(time.Since(start_time).Nanoseconds() / 1000000))
  }
}

func ShowResults() {
  if *size == failedRequestsCount() {
    fmt.Printf("\nAll requests failed!\n")
    os.Exit(1)
    return
  }
  fmt.Printf("\n\n%.2f requests/sec\n", meanRequestsPerSec())
  fmt.Printf("%.0f ms mean response time\n", meanResponseTime())
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
  return float32(*threads_count) * 1000.0 * float32(len(response_times)) / float32(sumTime())
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
    if requests_count < *size {
      percent := int32((float32(requests_count) / float32(*size)) * 100.0)
      fmt.Printf("%d%% done...\n", percent)
    } else {
      return
    }
  }
}

func failedRequestsCount() int {
  return *size - len(response_times)
}
