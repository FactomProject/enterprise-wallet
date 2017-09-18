FROM golang:1.8.3

RUN pwd

# Get git
RUN apt-get update \
    && apt-get -y install curl git \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Get glide
RUN go get github.com/Masterminds/glide

# Where enterprise-wallet sources will live
WORKDIR $GOPATH/src/github.com/FactomProject/enterprise-wallet

# Get the dependencies
COPY glide.yaml glide.lock ./

# Install dependencies
RUN glide install -v

# Populate the rest of the source
COPY . .

ARG GOOS=linux

RUN echo "1"
RUN cat /go/src/github.com/FactomProject/enterprise-wallet/vendor/github.com/FactomProject/factomd/common/constants/constants.go

RUN cat /go/src/github.com/FactomProject/enterprise-wallet/vendor/github.com/FactomProject/factomd/common/constants/checkpoints.go

# Build and install enterprise-wallet
RUN go install -v

ENTRYPOINT ["/go/bin/enterprise-wallet"]

EXPOSE 8091
