FROM golang:1.11-alpine3.7
WORKDIR /usr/local/go/src/github.com/Samollo/maain/
ADD . . 
RUN go build -o engine .
RUN ls
ENTRYPOINT ["/usr/local/go/src/github.com/Samollo/maain/engine"]