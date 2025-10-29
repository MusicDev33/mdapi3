# mdapi3

This repo's more of an experiment with AI coding agents. Essentially, I just made an empty repo with some Go code in it and then told Claude code to translate parts of MDAPI2 (written in TypeScript) into Go and put them in this repo. This README is also AI-generated too, though I'll probably chop it down soon. The repo itself is just my backend for all my services I like to write myself.

## What I Learned
AI in 2025 is damn impressive. I actually can't believe it was able to open my MDAPI2 repo and translate even just meaningfully small amounts of code into Go. It's not super great at the translation process though; I had to explicitly tell it not to just copy the structure of the old code base. It wasn't going to make the 'internal/' directory and was literally copying the directory structure from the TS version. Still some work to do there, but not really that big of a problem.

I also deliberately chose Go. I've written enough both on my own time and at work to know the ins and outs fairly well, so I was able to keep it on the rails. I think this is going to be the make or break moment for a lot of AI devs (which is going to be a thing soon if it isn't already): if you don't understand what you're doing on a deep level, your code's going to look wack and if you ever need help debugging it, people are going to spend more time parsing through a non-idiomatic code base than really anything else.

The PR was pretty big, and while that itself isn't a huge deal, because it was so easily generated, it's easy to see the PR itself as cheap. I only glanced over it to make sure nothing egregious made its way in, but these AI-generated PRs are probably going to be a gold mine of subtle logic errors. But maybe I'm just fearmongering here, I don't know. Just my predictions.

Last, AI is so damn verbose. It's annoying when engineers make comments that say "looping through list and doing x". Like yes, I can see that. I know when you're adding middleware. If I read the code, it tells me what it's doing most of the time. AI feels the need to tell me every minute thing it's doing and that's kind of stupid. I also asked it on a separate occasion to dockerize something for me and it made me an entire Docker/Docker Compose setup, which I didn't really ask for. It's way too verbose, and I suspect a large part of our jobs will soon be taken up by hacking away at a lot of the stuff AI generates. The flip side of that coin is the AI being more conservative with the code it writes and then we spend more time prompting the AI, and I think that sounds less desirable.

## Features

- **Multi-Engine AI Chat**: Support for OpenAI GPT-4, Anthropic Claude, and DeepSeek
- **User Authentication**: Secure bcrypt password hashing with whitelist-based registration
- **Conversation Management**: Create, retrieve, and delete chat conversations
- **Token Management**: Automatic context window management to prevent token overflow
- **MongoDB Storage**: Persistent storage for users, conversations, and chat messages
- **RESTful API**: Clean HTTP API with JSON responses

## Prerequisites

- Go 1.21 or higher
- MongoDB (running locally or remotely)
- API keys for at least one AI provider:
  - OpenAI API key (for GPT-4)
  - Anthropic API key (for Claude)
  - DeepSeek API key (for DeepSeek)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/MusicDev33/mdapi3.git
cd mdapi3
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o mdapi3
```

## Configuration

Create a `config.yaml` file in the project root with the following structure:

```yaml
port: 8080
mongoPort: 27017
dbName: mdapi3
akAnthropic: "your-anthropic-api-key"
akOpenAI: "your-openai-api-key"
akDeepSeek: "your-deepseek-api-key"
whitelistCors: "*"
whitelistUsers:
  - "user1"
  - "user2"
  - "admin"
```

### Configuration Options

| Field | Description | Required |
|-------|-------------|----------|
| `port` | Port for the HTTP server | Yes |
| `mongoPort` | MongoDB port | Yes |
| `dbName` | MongoDB database name | Yes |
| `akAnthropic` | Anthropic API key | Optional* |
| `akOpenAI` | OpenAI API key | Optional* |
| `akDeepSeek` | DeepSeek API key | Optional* |
| `whitelistCors` | CORS whitelist | Yes |
| `whitelistUsers` | List of usernames allowed to register | Yes |

*At least one AI provider API key is required.

## Running the Application

1. Start MongoDB:
```bash
# If using Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Or if MongoDB is installed locally
mongod --dbpath /path/to/data
```

2. Run the application:
```bash
./mdapi3
```

Or run directly without building:
```bash
go run main.go
```

The server will start on the configured port (default: 8080).

## API Endpoints

### Authentication

#### Register New User
```http
POST /login/create
Content-Type: application/json

{
  "username": "user1",
  "password": "secure-password"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "username": "user1",
    "_id": "507f1f77bcf86cd799439011"
  }
}
```

#### Login
```http
POST /auth
Content-Type: application/json

{
  "username": "user1",
  "password": "secure-password"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "username": "user1",
    "_id": "507f1f77bcf86cd799439011"
  }
}
```

#### Check Username Availability
```http
GET /verify/:username
```

**Response:**
```json
{
  "success": true  // true if username is available and whitelisted
}
```

### Conversations

#### Create Chat / Send Message
```http
POST /code
Content-Type: application/json

{
  "user": "user1",
  "convId": "",  // Empty string for new conversation
  "msg": "Hello, how are you?",
  "mode": "chat",  // "chat" or "code"
  "engine": "claude"  // "claude", "chatgpt", or "deepseek"
}
```

**Response:**
```json
{
  "success": true,
  "msg": "Successfully received response.",
  "newChat": {
    "_id": "507f1f77bcf86cd799439012",
    "conversationId": "507f1f77bcf86cd799439011",
    "role": "assistant",
    "content": "Hello! I'm doing well, thank you for asking...",
    "timestamp": 1698765432000
  },
  "newConversation": {
    "_id": "507f1f77bcf86cd799439011",
    "user": "user1",
    "name": "Brave Dolphin",
    "createdAt": "2024-10-29T12:00:00Z",
    "updatedAt": "2024-10-29T12:00:00Z"
  }
}
```

#### Get User's Conversations
```http
GET /convs/:username
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "_id": "507f1f77bcf86cd799439011",
      "user": "user1",
      "name": "Brave Dolphin",
      "createdAt": "2024-10-29T12:00:00Z",
      "updatedAt": "2024-10-29T12:00:00Z"
    }
  ]
}
```

#### Get Conversation Messages
```http
GET /msgs/:convId
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "_id": "507f1f77bcf86cd799439012",
      "conversationId": "507f1f77bcf86cd799439011",
      "role": "user",
      "content": "Hello, how are you?",
      "timestamp": 1698765432000
    },
    {
      "_id": "507f1f77bcf86cd799439013",
      "conversationId": "507f1f77bcf86cd799439011",
      "role": "assistant",
      "content": "Hello! I'm doing well...",
      "timestamp": 1698765432001
    }
  ]
}
```

#### Delete Conversation
```http
DELETE /convs/:convId
```

**Response:**
```json
{
  "success": true,
  "msg": "Successfully deleted conversation!"
}
```

### Utility

#### Health Check
```http
GET /test
```

**Response:**
```json
{
  "status": "okay"
}
```

## Chat Modes

### Chat Mode
Standard conversational mode with full AI capabilities.
```json
{
  "mode": "chat"
}
```

### Code Mode
Optimized for code generation with terse responses (temperature: 0.3).
```json
{
  "mode": "code"
}
```

## AI Engine Selection

You can specify which AI engine to use for each request:

- `"chatgpt"` - OpenAI GPT-4o
- `"claude"` - Anthropic Claude 3.7 Sonnet
- `"deepseek"` - DeepSeek Chat

**Note:** DeepSeek is automatically disabled in certain security contexts and will fall back to Claude.

## Development

### Project Structure
```
mdapi3/
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # MongoDB connection
│   ├── models/          # Data models (ZUser, Conversation, Chat)
│   ├── server/          # HTTP server and routing
│   └── zokyo/          # Zokyo chat functionality
│       ├── auth.go     # Authentication routes
│       ├── chat.go     # Chat creation and management
│       ├── engine.go   # AI engine integration
│       ├── getters.go  # Data retrieval routes
│       ├── delete.go   # Deletion routes
│       ├── naming.go   # Conversation name generator
│       └── assets.go   # Error messages and utilities
├── main.go             # Application entry point
├── config.yaml         # Configuration file (create this)
├── go.mod              # Go module definition
└── README.md           # This file
```

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mdapi3 .
```

## Docker Deployment

Create a `Dockerfile` (example):
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mdapi3 .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/mdapi3 .
COPY config.yaml .
EXPOSE 8080
CMD ["./mdapi3"]
```

Build and run:
```bash
docker build -t mdapi3 .
docker run -p 8080:8080 --link mongodb:mongodb mdapi3
```

## Security Considerations

- **Password Storage**: All passwords are hashed using bcrypt with cost factor 12
- **Whitelist-Only Registration**: Only pre-approved usernames can register
- **API Key Security**: Store API keys securely and never commit them to version control
- **Engine Security**: DeepSeek is automatically disabled in certain contexts for enhanced security
- **CORS Configuration**: Configure `whitelistCors` appropriately for your deployment

## Troubleshooting

### MongoDB Connection Issues
- Ensure MongoDB is running: `mongosh` or `mongo`
- Check the port in `config.yaml` matches your MongoDB instance
- Verify MongoDB is accessible from your application

### AI API Errors
- Verify your API keys are correct and have sufficient credits
- Check API rate limits for your provider
- Ensure network connectivity to AI provider endpoints

### Build Issues
- Ensure Go 1.21+ is installed: `go version`
- Clean module cache: `go clean -modcache`
- Re-download dependencies: `go mod download`

## License

This project is part of the mdapi suite. See LICENSE for details.

## Contributing

Contributions are welcome! Please submit pull requests or open issues for bugs and feature requests.

## Acknowledgments

This is a Go port of the mdapi2 TypeScript project, recreating the Zokyo chat functionality with improved performance and type safety.
