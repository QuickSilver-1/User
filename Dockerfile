FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk update && apk add --no-cache git
WORKDIR /user
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /user/cmd/user
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user/main .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder user/. .
WORKDIR /user/cmd/user
COPY --from=builder user/cmd/user .
EXPOSE 8080
CMD ["user/main"]