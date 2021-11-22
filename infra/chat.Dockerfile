FROM golang:1.17 as build
WORKDIR /opt/chat
COPY . .
RUN CGO_ENABLED=0 go build -o /opt/chat/snchat cmd/chat/main.go

FROM alpine:latest as release
COPY --from=build /opt/chat/snchat /bin/snchat
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Europe/Moscow
ENV ZONEINFO=/zoneinfo.zip
ENTRYPOINT ["/bin/snchat"]
