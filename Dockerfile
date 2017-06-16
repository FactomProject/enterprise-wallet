FROM golang:1.8.3-alpine

# Get git
RUN apk add --no-cache curl git

# Get glide
RUN go get github.com/Masterminds/glide

# Where enterprise-wallet sources will live
WORKDIR $GOPATH/src/github.com/FactomProject/enterprise-wallet

# Populate the source
COPY . .

# Install dependencies
RUN glide install -v

ARG GOOS=linux

# Build and install enterprise-wallet
RUN go install

ENTRYPOINT ["/go/bin/enterprise-wallet"]

EXPOSE 8091
