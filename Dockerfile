# FROM golang:1.16

# WORKDIR /go/src/app

# RUN go mod init main
# RUN go mod tidy
# RUN go install app
# EXPOSE 
FROM golang

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/LeCHAUVEaLaCape/Projet-Forum

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

RUN go mod init
# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["go-sample-app"]