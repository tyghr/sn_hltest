FROM golang:1.17
RUN mkdir -p /opt/chat
WORKDIR /opt/chat
COPY . .
RUN CGO_ENABLED=0 go build -o chat cmd/chat/main.go
RUN cp /opt/chat/chat /bin/chat
ENTRYPOINT ["/bin/chat"]
