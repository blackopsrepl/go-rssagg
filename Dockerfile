# Stage 1: BUILD

FROM golang:1.24.2-alpine AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go mod vendor

# -o specifies output name, CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -C cmd/go-rssagg -a -installsuffix cgo -o go-rssagg .


# Stage 2: RUNTIME IMAGE

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/cmd/go-rssagg/go-rssagg .

ENV DB_URL=$DB_URL
ENV PORT=$PORT

CMD ["./go-rssagg"]