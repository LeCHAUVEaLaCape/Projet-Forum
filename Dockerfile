FROM golang:1.15.1-alpine

RUN mkdir /serveurforumGroupe6
ADD . /serveurforumGroupe6
WORKDIR /serveurforumGroupe6

# enable commands to get github repositories
RUN apk add git
# use a gcc compiler to get sqlite repo
RUN apk add build-base

# get all repo needed
RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/satori/go.uuid
RUN go get golang.org/x/crypto/bcrypt
RUN go get golang.org/x/oauth2
RUN go get golang.org/x/oauth2/facebook
RUN go get golang.org/x/oauth2/google

RUN go build -o main .

CMD ["/serveurforumGroupe6/main"]

EXPOSE 8000