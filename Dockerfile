FROM golang:1.16.4-buster

RUN mkdir /serveurforumGroupe6

ADD . /serveurforumGroupe6
WORKDIR /serveurforumGroupe6
RUN go mod init main
RUN go mod tidy
RUN go build -o main .
CMD ["/serveurforumGroupe6/main"]
EXPOSE 8000