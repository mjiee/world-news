# build backend
FROM golang:1.25.4-alpine3.22 as go-builder

WORKDIR /build

COPY . .

# RUN go env -w GOPROXY=https://goproxy.io,direct
RUN mkdir -p /build/frontend/dist && \
    go mod tidy && \
    go build -tags web -o world-news main.go

# runtime stage
FROM alpine:3.22

WORKDIR /app
COPY --from=go-builder /build/world-news /app/world-news

CMD ["./world-news"]
