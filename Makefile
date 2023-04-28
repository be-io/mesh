#
#   Copyright (c) 2021, 2121, ducesoft and/or its affiliates. All rights reserved.
#   DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#

equals = $(and $(findstring x$(1),x$(2)), $(findstring x$(2),x$(1)))
contains = $(call equals, $(findstring x$(2), x$(1)), x$(2))
set = $(if $(call contains,x$(2),x$(3)), $(eval $(1):=$(4)))

.PHONY: clean build

.PHONY: setenv
setenv:
	$(eval TmpDir := $(shell mktemp -d))
	$(eval ProtoPath := $(TmpDir)/src:$(ProtoPath):$(TmpDir)/src/types/protos:$(TmpDir)/src/types/pb:.)
	$(eval GOOS:=$(shell go env GOOS))
	$(eval Branch:=$(or $(shell git branch --show-current),$(CI_BUILD_REF_NAME)))
	$(eval CommitID:=$(or $(shell git rev-parse --short HEAD),$(CI_COMMIT_SHORT_SHA)))
	$(eval Tag:=$(or $(shell git describe --tags --abbrev=0),$(CI_COMMIT_TAG)))
	#$(if $(Tag),$(shell echo $(Tag) | sed -e "s/[master\/|release\/|dev\/|feature\/|alpha\/|beta\/]//g"),$(Stage))
	$(eval Version:=$(shell echo $(or $(Branch),1.0.0) | sed -e "s/\//-/g"))
	$(eval MinaVersion:=$(shell echo $(Version) | sed -e "s/release-//g"))
	$(call set,Stage,$(Branch),dev,dev)
	$(call set,Stage,$(Branch),feature,dev)
	$(call set,Stage,$(Branch),beta,dev)
	$(call set,Stage,$(Branch),release,release)
	$(call set,Stage,$(Branch),master,dev)
	$(eval ImageRegistryUsername:=$(or $(HARBOR_USERNAME),admin))
	$(eval ImageRegistryPassword:=$(or $(HARBOR_PASSWORD),Abc12345))
	$(eval ImageRegistryUrl:=$(or $(REGISTRY_URL),10.12.0.78:5000))
	$(eval ImageID:=$(ImageRegistryUrl)/middleware-$(Stage)/gaia-mesh:$(Version)-$(CommitID))
	$(call set,ImageID,$(Stage),release,$(ImageRegistryUrl)/middleware-$(Version)/gaia-mesh:$(MinaVersion)-beta-$(CommitID))
	$(eval Flags:="-X github.com/be-io/mesh/client/golang/prsim.Version=$(MinaVersion) -X github.com/be-io/mesh/client/golang/prsim.CommitID=$(CommitID)")

.PHONY: dockerhub
dockerhub: setenv
	git fetch
	git pull origin $(Branch) --rebase
	docker login -u $(ImageRegistryUsername) -p $(ImageRegistryPassword) $(ImageRegistryUrl)

.PHONY: babel
babel:
	babel -i client/java -f java -o client/cpp -t cpp
	babel -i client/java -f java -o client/golang -t golang
	babel -i client/java -f java -o client/python -t python
	babel -i client/java -f java -o client/javascript -t javascript

.PHONY: clean
clean: setenv
	rm -rf bin
	cd client/java && mvn clean

.PHONY: test
test: setenv clean
	#go test ./...
	#cd client/java && mvn clean test
	#cd client/python
	#cd client/cpp
	#docker compose up .

.PHONY: image
image:setenv
	buildctl --addr tcp://10.99.253.223:1234 build \
	--frontend dockerfile.v0 \
	--local context=. \
	--local dockerfile=./build/oci/golang/ \
	--output type=image,registry.insecure=true,name=$(ImageID),push=true \
	--allow security.insecure \
	--opt filename=Dockerfile \
	--opt platform=linux/amd64 \
	--opt build-arg:flags=$(Flags)

.PHONY: exe
exe:setenv
	buildctl --addr tcp://10.99.253.223:1234 build \
	--frontend dockerfile.v0 \
	--local context=. \
	--local dockerfile=./build/oci/golang/ \
	--output type=local,dest=. \
	--allow security.insecure \
	--opt filename=Gorun.Dockerfile \
	--opt platform=linux/amd64,linux/arm64 \
	--opt build-arg:flags=$(Flags)

.PHONY: check
check:
	@./depcheck.sh && \
		(echo "Installing proto libraries to versions in go.mod." ; \
		go install github.com/golang/protobuf/protoc-gen-go ; \
		go install github.com/gogo/protobuf/protoc-gen-gogofaster)

.PHONY: regenerate
regenerate: setenv clean check
	@protoc \
		--proto_path=/usr/local/include \
		--proto_path=/usr/include \
		--proto_path=$(ProtoPath) \
		--gogofaster_out=plugins=grpc,Mapi.proto=$(pwd)/protos/api:pb \
		pb.proto
	@rm -rf $(TmpDir)
	@echo Done.

.PHONY: ci-golang
ci-golang:
	buildctl --addr tcp://10.99.253.223:1234 build \
	--frontend dockerfile.v0 --local context=. \
	--local dockerfile=./build/oci/golang \
	--output type=image,registry.insecure=true,name=10.12.0.78:5000/cosmos/ci:golang-1.19,push=true \
	--allow security.insecure \
	--opt filename=BUILD.Dockerfile \
	--opt platform=linux/amd64,linux/arm64
