# build backend
FROM golang:1.26.1-alpine3.23 as go-builder

WORKDIR /build

COPY . .

# RUN go env -w GOPROXY=https://goproxy.io,direct
RUN mkdir -p /build/frontend/dist && \
    touch /build/frontend/dist/index.html && \
    go mod tidy && \
    go build -tags web -o world-news main.go

# runtime stage
FROM alpine:3.23

WORKDIR /app
COPY --from=go-builder /build/world-news /app/world-news

CMD ["./world-news"]
