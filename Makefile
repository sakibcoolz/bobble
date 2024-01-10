include compose/.env
include variables.mk

docker-image:
	make -C docker docker-image

run:
	HOSTPORT=$(HOSTPORT1) \
	OTHERS=$(OTHERS) \
	$(GORUN) main.go 

run2:
	HOSTPORT=$(HOSTPORT2) \
	OTHERS=$(OTHERS) \
	$(GORUN) main.go 

run3:
	HOSTPORT=$(HOSTPORT3) \
	OTHERS=$(OTHERS) \
	$(GORUN) main.go 

kick-start: docker-image
	docker-start


docker-start:
	make -C compose start

docker-stop:
	make -C compose stop

PHONY: start stop
