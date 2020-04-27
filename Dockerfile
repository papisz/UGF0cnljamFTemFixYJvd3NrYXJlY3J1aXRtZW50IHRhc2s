FROM golang:1.14-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/weather

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test ./...

# Build the Go app
RUN go build -o ./out/weather cmd/weather/weather.go

# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates


COPY --from=build_base /tmp/weather/out/weather /app/weather

# This container exposes port 5555 to the outside world
EXPOSE 5555

# Run the binary program produced by `go install`
CMD ["/app/weather"]