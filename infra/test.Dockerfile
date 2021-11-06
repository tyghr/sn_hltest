FROM golang:1.17
WORKDIR /opt/sntest
COPY . .
RUN CGO_ENABLED=0 go build -o sntest infra/bench/mysql_read.go
RUN cp /opt/sntest/sntest /bin/sntest
ENTRYPOINT ["/bin/sntest"]
