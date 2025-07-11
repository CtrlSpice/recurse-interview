# Recurse Center Interview Project: Database Server

Hello there!

This is a small project I built for my Recurse Center pairing interview.

## The Prompt:

"Before your interview, write a program that :
[] Runs a server that is accessible on http://localhost:4000/. 
[] When your server receives a request on http://localhost:4000/set?somekey=somevalue it should store the passed key and value in memory.
[] When it receives a request on http://localhost:4000/get?key=somekey it should return the value stored at somekey.

During your interview, you will pair on saving the data to a file. You can start with simply appending each write to the file, and work on making it more efficient if you have time."

## Interview Goals:

[] Create a write-ahead log
[] Compress the write-alead log
[] If I have time, try my hand at a really simple Bloom filter because I know about the theory, but have never implemented one in practice.

## Getting Started

### Prerequisites
- Go 1.24.2 or later

### Running the project
```bash
go run .
```

### Building the project
```bash
go build -o recurse-interview
```

### Trying it out
```bash
curl "http://localhost:4000/set?key=hello&value=there"
curl -w "\n" "http://localhost:4000/get?key=hello"
```