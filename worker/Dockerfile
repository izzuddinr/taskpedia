#build stage worker
FROM golang:alpine AS worker-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o taskpedia-worker ./main.go

WORKDIR /app/dist

RUN cp /app/taskpedia-worker .

#final stage worker
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY . .
COPY --from=worker-builder /app/dist/taskpedia-worker /

ENTRYPOINT ["/taskpedia-worker"]
