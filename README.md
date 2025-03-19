# Parse Practice
Generated logs with GPT and went through parsing into structs, then using an interface to print the different struct types.

## To Run
**Required:** Go

Clone repo to local system

By default the program parses "./log/error.log"

To run with default path open a terminal at the root of the repo and use:
```bash
go run . 
```
The program can also accept filename arguments to specify which log to parse:
```bash
go run . debug.log
go run . error.log
go run . info.log
```

The program can parse the following logs.

* Debug
```
2025-03-18T15:06:12Z INFO HTTP request received: method=GET, path=/api/items, status=200, duration=150ms
```
* Error
```
2025-03-18T15:05:01Z ERROR Database connection failed: timeout error (db=users, retry=3)
```
* Info
```
2025-03-18T15:04:05Z INFO Server started on port 8080
```
