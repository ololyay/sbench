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

# sbench -n 5000 -t 50 -u http://example.com
Starting test http://example.com
Requests: 5000
Threads: 50
Timeout: 10 sec

18% done...
35% done...
52% done...
69% done...
86% done...


174.17 requests/sec
287 ms mean response time
Percentage of requests processed within a certain time:
25%: 118 ms
50%: 123 ms
75%: 138 ms
90%: 1250 ms
95%: 1275 ms
98%: 1296 ms
5000 requests total.
0 requests failed.

```

---
## Examples:

Make 10 POST requests with data:
```
sbench -u http://example.com -n 10 --method=POST --content-type="application/x-www-form-urlencoded" --body="key=value"
```

Make 50 GET requests in 2 parallel threads:
```
sbench -u http://example.com -n 50 -t 2
```
