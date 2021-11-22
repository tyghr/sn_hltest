FROM golang:1.17 as build
WORKDIR /opt/chat
COPY . .
RUN CGO_ENABLED=0 go build -o /opt/chat/snbalancer cmd/chat/balancer/main.go

FROM alpine:latest as release
COPY --from=build /opt/chat/snbalancer /bin/snbalancer
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Europe/Moscow
ENV ZONEINFO=/zoneinfo.zip
ENTRYPOINT ["/bin/snbalancer"]
