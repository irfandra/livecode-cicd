FROM golang:alpine as build
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o simple-golang

FROM alpine
WORKDIR /app
COPY --from=build /app/simple-golang /app
ENTRYPOINT [ "/app/simple-golang" ]