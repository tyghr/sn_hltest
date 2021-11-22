FROM golang:1.17 as build
WORKDIR /opt/snserver
COPY . .
RUN CGO_ENABLED=0 go build -o /opt/snserver/snserver cmd/server/main.go

FROM alpine:latest as release
COPY --from=build /opt/snserver/snserver /bin/snserver
COPY --from=build /opt/snserver/html_tmpl/ /html_tmpl/
COPY --from=build /opt/snserver/migrations/ /migrations/
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Europe/Moscow
ENV ZONEINFO=/zoneinfo.zip
ENTRYPOINT ["/bin/snserver"]
