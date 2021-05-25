FROM golang:1.16.4-buster

RUN mkdir /serveur-forum-BPEYRE

ADD . /serveur-forum-BPEYRE
WORKDIR /serveur-forum-BPEYRE
RUN go mod init main
RUN go mod tidy
RUN go build -o main .
CMD ["/serveur-forum-BPEYRE/main"]
EXPOSE 8000