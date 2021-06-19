FROM golang:alpine AS builder

RUN apk update \
 && apk add --no-cache git \
 && mkdir -p $GOPATH/src/github.com/earthquakesan/freeipa-group-sync/ \
 && git clone https://github.com/earthquakesan/goipa.git $GOPATH/src/github.com/earthquakesan/goipa

WORKDIR $GOPATH/src/github.com/earthquakesan/freeipa-group-sync/
COPY . .
RUN go get -d -v
RUN go build -o /go/bin/freeipa-group-sync

FROM alpine:3.14

ENV IPA_HOST=ipa.example.test
ENV IPA_REALM=EXAMPLE.TEST
ENV IPA_USERNAME=admin
ENV IPA_PASSWORD=Pa33s0000rd
ENV IPA_GROUPS_YAML_PATH=/data/groups.yaml

COPY --from=builder /go/bin/freeipa-group-sync /go/bin/freeipa-group-sync

ENTRYPOINT ["/go/bin/freeipa-group-sync"]