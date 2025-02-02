- Build a system that can efficiently analyze and process large amounts of log data
- Process each line in the log file simultaneously to perform log analysis (e.g. error log count, log detection of specific pattern, etc.).
- Analysis results must be output to the console or saved as a new file.
- Design the system considering optimized processing performance and memory usage

1. Error log number and ratio analysis:
   Calculate the number of logs for error levels (“ERROR”, “WARNING”, “DEBUG”) in the log file, calculate the ratio for each, and output.
2. Log analysis by HTTP status code:
   Calculate and output the number of logs for each HTTP status code in the log file.
3. Response time analysis:
   Extract the response time from the log file and calculate and output the average response time, minimum response time, and maximum response time.
   Calculate and output the percentage of requests with response times exceeding the specified threshold (500ms).
4. Parse the request URI:
   Find the most requested domains in the log file and print them, up to 5th place in ranking.
5. Analysis by time period
   Analyze and output request frequency by time zone.
6. (Additional items to do if you have time) Build with cli
7. Branch Rule: Use the format feature/name for branch names.

### Prerequire

Go version 1.22 or more

### Run program

```
    go run main.go
```

Or

```
    make dev
```

### test performance on macbook m1

version 1 _feature/log-insight_

- exec time: 191.739417ms

version 2 _feature/v2-with-concurrency_

- faster then v1 86.41%
- exec time: 26.060583ms
- fix loadtime zone to global
- decode json with https://github.com/json-iterator/go
- add go concurrency
- go version 1.22
