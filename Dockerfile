FROM golang:alpine

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /app

# Copy the current directory contents into the container at /app
COPY ./src /app

RUN ["go", "get", "-u", "github.com/go-telegram-bot-api/telegram-bot-api"]
RUN ["go", "get", "github.com/lib/pq"]

# Make port 80 available to the world outside this container
#EXPOSE 80

# Run go build when the container launches
RUN ["go", "build"]

#CMD ["./app"]