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
      --timeout=10    Timeout for requests (sec)
      --version       Show application version.

# sbench -u http://example.com -n 50 -t 2
Starting test http://example.com
Requests: 50
Threads: 2
Timeout: 10 sec

58% done...


3.1080997 requests/sec
321.74 ms mean response time
Percentage of requests processed within a certain time:
25%: 274 ms
50%: 278 ms
75%: 283 ms
90%: 290 ms
95%: 305 ms
98%: 1363 ms
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
