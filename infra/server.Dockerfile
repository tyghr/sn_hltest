FROM golang:1.17
RUN mkdir -p /opt/snserver
WORKDIR /opt/snserver
COPY . .
RUN CGO_ENABLED=0 go build -o snserver cmd/server/main.go
RUN cp /opt/snserver/snserver /bin/snserver
ENTRYPOINT ["/bin/snserver"]
