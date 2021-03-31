FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /bin/hcledit ./cmd/hcledit

FROM alpine:3.13

COPY --from=builder /bin/hcledit /usr/bin
CMD ["hcledit"]
