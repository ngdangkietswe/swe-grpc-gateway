FROM golang:1.23-alpine as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine
WORKDIR /app/
COPY --from=builder /app/app .
EXPOSE 7777
ENTRYPOINT ["./app"]