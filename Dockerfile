FROM linuxbrew/linuxbrew

WORKDIR /go/src/github.com/josephnp732/RabbitMQ-Logger

COPY . .

RUN brew tap mingrammer/flog
RUN	brew install flog


