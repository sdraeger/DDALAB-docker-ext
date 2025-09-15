.PHONY: build install clean push dev

EXTENSION_NAME = sdraeger/ddalab-manager
EXTENSION_TAG = latest

build:
	docker build -t $(EXTENSION_NAME):$(EXTENSION_TAG) .

install: build
	docker extension install $(EXTENSION_NAME):$(EXTENSION_TAG)

update: build
	docker extension update $(EXTENSION_NAME):$(EXTENSION_TAG)

uninstall:
	docker extension rm $(EXTENSION_NAME):$(EXTENSION_TAG)

clean:
	rm -f backend/ddalab-manager

push: build
	docker push $(EXTENSION_NAME):$(EXTENSION_TAG)

dev:
	docker extension dev debug $(EXTENSION_NAME):$(EXTENSION_TAG)

validate:
	docker extension validate $(EXTENSION_NAME):$(EXTENSION_TAG)

logs:
	docker extension dev logs $(EXTENSION_NAME):$(EXTENSION_TAG)