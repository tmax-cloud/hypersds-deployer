.PHONY: build container clean
UBUNTU_IMAGE = ubuntu:18.04
CRI = docker
PROVISIONER_IMAGE = hypersds-provisioner
PROVISIONER_TAG = test

build:
	mkdir -p build
	go build -o build/ ./...

container:
ifeq ($(REGISTRY),)
	$(CRI) build -t $(PROVISIONER_IMAGE):$(PROVISIONER_TAG) . --build-arg BASE_IMAGE=$(UBUNTU_IMAGE)
else
	$(CRI) build -t $(PROVISIONER_IMAGE):$(PROVISIONER_TAG) . --build-arg BASE_IMAGE=$(REGISTRY)/$(UBUNTU_IMAGE)
	$(CRI) tag $(PROVISIONER_IMAGE):$(PROVISIONER_TAG) $(REGISTRY)/$(PROVISIONER_IMAGE):$(PROVISIONER_TAG)
	$(CRI) push $(REGISTRY)/$(PROVISIONER_IMAGE):$(PROVISIONER_TAG)
endif

clean:
	rm -rf build
