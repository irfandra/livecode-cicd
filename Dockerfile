FROM golang:alpine as build
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o book-go

FROM alpine
WORKDIR /app
COPY --from=build /app/book-go /app
ENTRYPOINT [ "/app/book-go" ]