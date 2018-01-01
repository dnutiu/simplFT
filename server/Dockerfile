# Golang is our base images.
FROM golang:1.7

# Make a directory called simplFT
RUN mkdir -p /go/src/github.com/metonimie/simplFT/server

# Copy the current dir contents into simplFT
ADD . /go/src/github.com/metonimie/simplFT/server

# Set the working directory to simplFT
WORKDIR /go/src/github.com/metonimie/simplFT/server

# Make port 80 available to the world outside this container
EXPOSE 8080

# Install dependencies
RUN go get "github.com/zyxar/image2ascii/ascii"
RUN go get "github.com/spf13/viper"

# Build the application
RUN go build ./simplFTP.go

# Run simplFT when the container launches
CMD ["./simplFTP", "-config-name", "docker-config"]
