# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy root go.mod and go.sum (from context root)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project to build the specific agent
COPY . .

# Build the agent
RUN go build -o /app/agent ./agents/agent-analyst/main.go

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/agent .

EXPOSE 50051
CMD ["./agent"]
