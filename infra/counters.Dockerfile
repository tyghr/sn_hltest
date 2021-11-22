FROM golang:1.17 as build
WORKDIR /opt/sn
COPY . .
RUN CGO_ENABLED=0 go build -o /opt/sn/sncounters cmd/counters/main.go

FROM alpine:latest as release
COPY --from=build /opt/sn/sncounters /bin/sncounters
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Europe/Moscow
ENV ZONEINFO=/zoneinfo.zip
ENTRYPOINT ["/bin/sncounters"]
