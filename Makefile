.PHONY: build container clean

build:
	go build -o build/hypersds-provisioner

container:
	docker build -t hypersds-provisioner:canary .
	
clean:
	rm -rf build
