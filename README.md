# High Throughput Pub/Sub Logger (RabbitMQ + Graylog)

The Go Lang program generates random Apache-Commons logs every 400ms into a RabbitMQ direct exchange queue and subscriber recieves the logs and sends them to Graylog for monitoring

(Note: This repo experiments with the concepts of RabbitMQ, specifically the direct exhange)

### _**Disclaimer:**_ Works only on MacOS

## Pre-requisites:

* Requires MacOS with Homebrew installed (https://brew.sh/)
* Install Docker (https://docs.docker.com/docker-for-mac/install/)
* Using homebrew install `flog` (https://github.com/mingrammer/flog)


## Steps to Run the project:

* Install and setup Go 
    - https://golang.org/doc/install
    - Go Modules: https://github.com/golang/go/wiki/Modules
* `go mod download` to download the requirements
* `go mod vendor` to download to local reposistory
* Spin-up a RabbitMQ docker container from the _./RabbitMQ_ directory (`docker-compose up -d`)
* Spin-up a Graylog Instance from the _./graylog_ directory (`docker-compose up -d`)

#### Running the subscriber: 
* Change to the publisher directory&nbsp;&nbsp; `cd ./subscriber`
* Start the publisher:&nbsp;&nbsp; `go run subscribe.go`

#### Running the publisher: 
* Change to the publisher directory&nbsp;&nbsp; `cd ./publisher`
* Start the publisher:&nbsp;&nbsp; `go run publish.go`

``` diff
! Alternatively, use MakeFile to run the targets
+ Example: make -f MakeFile install
```

#### _TODO:_

* Dockerize the application
* Make it platform/OS independent





