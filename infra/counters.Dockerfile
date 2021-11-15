FROM golang:1.17
WORKDIR /opt/sn
COPY . .
RUN CGO_ENABLED=0 go build -o sncounters cmd/counters/main.go
RUN cp /opt/sn/sncounters /bin/sncounters
ENTRYPOINT ["/bin/sncounters"]