# gogpt
using ollama api build gpt chat

## SetUp

```bash
➜ go run cmd/chat/main.go 
Welcome to GPT chatbot! Type 'exit' to quit.

You: gRPC stream vs Unary?

Bot: gRPC is a high-performance remote procedure call (RPC) framework that allows you to define service interfaces in protocol buffers and generate client and server code in multiple programming languages. When it comes to gRPC, two primary types of requests are used: unary and streaming.

**Unary Requests**

A unary request is a single request where the client sends a message to the server, and the server returns a response. The client and server have a one-to-one conversation. Here's an example:

* Client: Send a GET request to retrieve data.
* Server: Receive the request and return the requested data.

Unary requests are suitable for simple, stateless operations like retrieving or updating data. They are efficient because they don't require any additional resources or bandwidth.

**Streaming Requests**

A streaming request is a continuous flow of messages where the client sends multiple requests to the server, and the server responds with multiple responses. The client and server have a conversation that continues until one party stops sending messages. Here's an example:

* Client: Send a sequence of GET requests to retrieve data.
* Server: Receive each request and respond with corresponding data.

Streaming requests are suitable for long-running operations like video playback, audio streaming, or large data transfer. They provide benefits like efficient resource usage and error handling because the server can handle multiple requests concurrently without creating new connections.

**Key differences**

Here's a summary of the key differences between unary and streaming gRPC requests:

| **Characteristics** | **Unary Requests** | **Streaming Requests** |
| --- | --- | --- |
| **Request-response pattern** | Single request, single response | Multiple requests, multiple responses |
| **Server handling** | Handle each request individually | Handle multiple requests concurrently |
| **Resource usage** | Create a new connection for each request | Reuse the same connection for all requests |
| **Error handling** | Typically use HTTP status codes to handle errors | Use gRPC's built-in error handling mechanisms, such as `Status` messages |

When choosing between unary and streaming gRPC requests, consider the following factors:

* The type of operation: Is it a simple stateless request or a long-running operation with multiple responses?
* Resource constraints: Will the server need to handle many concurrent requests, or can it allocate resources for each individual request?
* Error handling: How will errors be handled in your system?
You: exit
```
