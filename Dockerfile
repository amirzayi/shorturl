FROM golang:1.24-alpine3.21 as builder
LABEL authors="amirzayi"
WORKDIR /app
COPY . .
RUN go mod download && go build -o app .

FROM alpine
COPY --from=builder /app/app .
ENV HTTP_PORT=8000
EXPOSE 8000
ENTRYPOINT ["./app"]