PACKAGE_OS := linux darwin windows
PACKAGE_ARCH := amd64 386
PRODUCT_NAME := elastic-watcher

build:
	go build .

package: clean
	go get github.com/mitchellh/gox && \
	mkdir -p build && \
	gox -os="$(PACKAGE_OS)" -arch="$(PACKAGE_ARCH)" -output="build/$(PRODUCT_NAME)_{{.OS}}_{{.Arch}}/$(PRODUCT_NAME)" && \
	mkdir -p dist && \
	ls -1 build | xargs -I% tar zcf "dist/%.tar.gz" -C "build/%" $(PRODUCT_NAME)

clean:
	rm -rf build dist