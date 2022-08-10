FROM golang:1.18-alpine as build
WORKDIR /app
COPY ["go.mod", "go.sum" , "./"]
RUN go mod download
COPY . .
RUN go build src/main.go

FROM alpine
WORKDIR /app
RUN apk add chromium
COPY --from=build /app/main ./main
CMD ["./main"]