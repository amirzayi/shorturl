FROM golang:1.24-alpine3.21 as builder
LABEL authors="amirzayi"
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 go build -o app -ldflags '-s -w' .

FROM alpine
COPY --from=builder /app/app .
ENV HTTP_PORT=8000
EXPOSE 8000
ENTRYPOINT ["./app"]