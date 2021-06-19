VERSION := $(shell git describe --tags --always)

docker:
	go build
	docker build -t earthquakesan/freeipa-group-sync:${VERSION} .
	go clean

docker-test-run:
	docker run -it --rm \
	-v $(shell pwd)/groups.yaml:/data/groups.yaml \
	-v /etc/ipa/ca.crt:/etc/ipa/ca.crt \
	--name freeipa-group-sync \
	earthquakesan/freeipa-group-sync:${VERSION}
