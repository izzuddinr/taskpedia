#build stage rest
FROM golang:alpine AS rest-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o taskpedia-rest ./main.go

WORKDIR /app/dist

RUN cp /app/taskpedia-rest .

#final stage rest
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY . .
COPY --from=rest-builder /app/dist/taskpedia-rest /

ENTRYPOINT ["/taskpedia-rest"]