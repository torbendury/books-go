# syntax=docker/dockerfile:1
ARG DEBIAN_FRONTEND=noninteractive

# build application without CGO for linux (amd64 by default)
FROM golang:1.20.2 AS build
WORKDIR /app
COPY /go.mod /go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /books-go cmd/main.go
CMD ["/books-go"]

# run tests
FROM build AS test
RUN CGO_ENABLED=1 go test -race -v ./...

# package application into alpine (might move do distroless later but small distro is great for debugging)
FROM alpine:3.18 AS release
WORKDIR /app
COPY --from=build /books-go /books-go
EXPOSE 3000
USER nobody:nobody
ENTRYPOINT [ "/books-go" ]
