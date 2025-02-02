# Handle
## RequestTable
- When "Fetch" button is clicked, add URLs to the RequestTable
- The RequestTable is mapped 1:1 with the GUI table
- "Clear" button clears the RequestTable except "Downloading" requests
- "Start" button starts the downloading process
- "Terminate" button terminates a particular request (atomically)
- No locking required (purely UI thread)

## RequestQueue
- The request queue holds the requests to be downloaded, guarded by an R-W lock
- When "Download" button is clicked:
  - Obtain write lock
  - Move all "Downloading" requests to the queue

## Download Scheduling
- The download scheduling task periodically checks the queue for new requests (obtain the read lock)
- It holds an internal counter of concurrent workers
- It polls a number of requests from the queue such that the new count of concurrent workers do not exceed the setting
- Each eligible request is assigned to a new worker (goroutine) managed by Go runtime

## Download Worker
- The download worker is responsible to:
  - Execute the download command
  - Receive output from yt-dlp process
  - Receive termination request
  - Report download progress

## Concurrent Goroutine
- Main thread (= UI thread, ...unsure?)
- Download scheduler
- n Download workers