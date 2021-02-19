.PHONY: build container clean

build:
	go build -o build/hypersds-provisioner

container:
	docker build .
	
clean:
	rm -rf build
