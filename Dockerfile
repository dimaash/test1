FROM golang:1.14.3-alpine AS build_base

RUN apk add --no-cache git

# Set the Current tmp Working Directory inside the container
WORKDIR /tmp/go-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
# COPY go.sum .

RUN go mod download

COPY . .

RUN pwd
RUN ls -al 

# Build our app
RUN go build -o ./out/go-app .
RUN cp /tmp/go-app/users.json /tmp/go-app/out/users.json

# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates

WORKDIR /app

COPY --from=build_base /tmp/go-app/out/go-app /app/go-app
COPY --from=build_base /tmp/go-app/out/users.json /app/users.json


CMD ["/app/go-app"]


