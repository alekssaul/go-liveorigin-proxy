FROM golang:latest as builder
WORKDIR /go/src/github.com/alekssaul/go-liveorigin-proxy/origin
COPY . .
RUN mkdir -p /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/origin .


FROM alpine:latest
RUN apk update ;  apk add --no-cache ca-certificates ; update-ca-certificates ; mkdir /app
WORKDIR /app
COPY --from=builder /app .
CMD /app/origin
EXPOSE 8080