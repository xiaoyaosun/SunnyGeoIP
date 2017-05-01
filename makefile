
PACKAGE = geoip
GOSRCS  = $(wildcard src/$(PACKAGE)/*.go)

run: bin/$(PACKAGE)
	bin/$(PACKAGE) 

build: bin/$(PACKAGE)
	@echo "$(PACKAGE) built ok"

test: $(GOSRCS)
	export GOPATH=`pwd` && go test  -c -o bin/$(PACKAGE)test $(GOSRCS) && bin/$(PACKAGE)test

lint: bin/gometalinter
	export GOPATH=`pwd` && bin/gometalinter src/$(PACKAGE)/...

bin/gometalinter:
	export GOPATH=`pwd` && go get -u github.com/alecthomas/gometalinter
	export GOPATH=`pwd` && bin/gometalinter --install --update

clean:
	rm -rf pkg bin

bin/$(PACKAGE): $(GOSRCS)
	@export GOPATH=`pwd` && go fmt $(PACKAGE)
	export GOPATH=`pwd` && go build -o bin/$(PACKAGE) $(GOSRCS)

srcinstall:
	cd src && glide install

docker:
	docker build -t $(PACKAGE) .

docker-run:
	docker run -d -p 8087:8087 --name $(PACKAGE) $(PACKAGE)
