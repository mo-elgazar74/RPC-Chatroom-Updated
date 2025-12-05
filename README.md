# 🚀 Real-Time Broadcast Chatroom (Go Concurrency)

This assignment extends the previous **Go RPC chatroom** by redesigning
it into a **real-time concurrent broadcast system** using Go's
**goroutines**, **channels**, and **mutex synchronization**.

The server manages multiple TCP clients simultaneously, broadcasting
messages instantly to all connected users (except the sender), and
announcing when users join or leave the chatroom.

------------------------------------------------------------------------

## Features

### ✅ Concurrent Real-Time Server

-   Uses **goroutines** for concurrent client handling
-   Broadcasts all messages instantly to every other client
-   Sends join/leave notifications:
    -   `User [ID] joined`
    -   `User [ID] left`
-   Each client has:
    -   A dedicated **reader goroutine**
    -   A dedicated **writer goroutine**
    -   A buffered outbound **channel**
-   Shared client list protected by **sync.Mutex**

------------------------------------------------------------------------

### ✅ Interactive TCP Client

-   Runs locally using:

        go run client.go

-   Sends realtime messages to the server

-   Displays all messages broadcasted by other users

-   Does **not** echo its own messages back

-   Cleanly exits with `exit`

------------------------------------------------------------------------

## Folder Structure

    RPC-Chatroom-Updated/
    │
    ├── server/
    │   └── server.go
    │
    └── client/
        └── client.go

------------------------------------------------------------------------

## How to Run (Server + Multiple Clients)

### 1) Run the Server

Open a terminal:

``` bash
cd server
go run server.go
```

The server listens on:

    localhost:1234

------------------------------------------------------------------------

### 2) Run Clients

Open **one or more terminals**:

``` bash
cd client
go run client.go
```

Example interaction:

    Enter your name: Mo
    Welcome Mo!
    User [NO] joined
    [NO]: hi
    [Mo]: hello everyone

Another client:

    Enter your name: NO
    [Mo]: hello everyone
    how are you?

------------------------------------------------------------------------

## Example Behaviors

### ▶ When a new client joins:

    User [Alice] joined

### ▶ When a client sends a message:

    [Alice]: hello everyone!

### ▶ When a client leaves:

    User [Alice] left

------------------------------------------------------------------------

## Concurrency Design

### Server-side concurrency

  Component            Mechanism
  -------------------- ------------------------
  Client readers       goroutine per client
  Client writers       goroutine per client
  Broadcast messages   channels
  Client list          map protected by Mutex

### Client-side concurrency

  Component                  Mechanism
  -------------------------- ----------------------------------
  Server → client messages   goroutine listening continuously
  User input → server        main thread

------------------------------------------------------------------------

## Summary

This project demonstrates:

-   Real-time communication
-   Concurrency with goroutines
-   Synchronization using mutex
-   Channel-based messaging
-   Multi-client handling over TCP
-   Event-driven broadcasting
