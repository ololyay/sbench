# sbench - Simple HTTP benchmarking tool

```
# sbench --help
usage: sbench --url=URL [<flags>]

Flags:
      --help          Show context-sensitive help (also try --help-long and --help-man).
  -u, --url=URL       URL to benchmark
  -m, --method="GET"  HTTP method
  -n, --number=100    Number of requests to perform
  -t, --threads=1     Number of threads
      --content-type=CONTENT-TYPE
                      HTTP content type
      --body=BODY     Body of HTTP request
      --version       Show application version.

# sbench -u http://example.com -n 50 -t 2
Starting 2 threads to make 50 GET requests to http://example.com
3.5345678 requests/sec
282.92 ms mean response time
Percentage of requests processed within a certain time:
25%: 277 ms
50%: 281 ms
75%: 285 ms
90%: 292 ms
95%: 307 ms
98%: 334 ms
50 requests total.
0 requests failed.
```

---
## Examples:

Make 10 POST requests with data:
```
sbench -u http://localhost:8080 -n 10 --method=POST --content-type="application/x-www-form-urlencoded" --body="key=value"
```

Make 50 GET requests in 2 parallel threads:
```
sbench -u http://example.com -n 50 -t 2
```
