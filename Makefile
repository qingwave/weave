PKG=qinng.io/weave
GOPATH=$(shell pwd)/.gopath
TESTPKG=$(shell go list ./...)

base: clean
		@mkdir -p $(GOPATH)/src/$(shell dirname $(PKG))
		@mkdir -p bin
		@ln -s $(shell pwd) $(GOPATH)/src/$(PKG)
		@ln -s $(shell pwd)/bin $(GOPATH)/bin

		GO111MODULE=on go install -v  $(GOPATH)/src/$(PKG)/cmd/weave/main.go
		chmod -R 755 $(GOPATH)

test:
	go test -ldflags -s -v --cover $(TESTPKG)

TAG = latest
REPOSITORY = cr.d.xiaomi.net/containercloud/weave:$(TAG)
OCEAN_REPOSITORY = push.docker.pt.xiaomi.com/base/weave:$(TAG)
dockerbuild: test
	docker build -t $(REPOSITORY) .
	docker tag $(REPOSITORY) $(OCEAN_REPOSITORY)
dockerpush: dockerbuild
	docker push $(REPOSITORY)

clean:
		@rm -rf $(GOPATH)/bin $(GOPATH)/src
		@rm -rf bin/