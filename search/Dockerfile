#build stage search
FROM golang:alpine AS search-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o taskpedia-search ./main.go

WORKDIR /app/dist

RUN cp /app/taskpedia-search .

#final stage search
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY . .
COPY --from=search-builder /app/dist/taskpedia-search /

ENTRYPOINT ["/taskpedia-search"]