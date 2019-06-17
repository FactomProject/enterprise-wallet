FROM golang:1.9

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

# Build and install enterprise-wallet
RUN go install

ENTRYPOINT ["/go/bin/enterprise-wallet"]

EXPOSE 8091
