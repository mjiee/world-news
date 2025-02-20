# build frontend 
FROM node:23.8.0-alpine3.21 as node-builder

WORKDIR /build

COPY frontend/. .

# RUN npm config set registry https://registry.npmmirror.com/
RUN npm install && \
    npm run build-web

# build backend
FROM golang:1.23.6-alpine3.21 as go-builder

WORKDIR /build

COPY . .
COPY --from=node-builder /build/dist /build/frontend/dist

# RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go mod tidy && \
    go build -tags web -o world-news main.go

# runtime stage
FROM alpine:3.21

WORKDIR /app
COPY --from=go-builder /build/world-news /app/world-news

CMD ["./world-news"]
