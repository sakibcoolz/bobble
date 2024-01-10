FROM golang:alpine AS builder
RUN apk add --no-cache git gcc libc-dev
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o zsync_service ./
RUN ls -lrth

# final stage
FROM alpine:latest
LABEL Name=zsync_service Version=0.0.1
RUN apk --no-cache add ca-certificates curl
COPY --from=builder /app/zsync_service zsync_service
EXPOSE 9000:9020
ENTRYPOINT ["/zsync_service"]