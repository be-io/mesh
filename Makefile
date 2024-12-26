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
	$(eval ImageID:=bfia/mesh:$(Version)-$(CommitID))
	$(eval Flags:="-X github.com/opendatav/mesh/client/golang/prsim.Version=$(MinaVersion) -X github.com/opendatav/mesh/client/golang/prsim.CommitID=$(CommitID)")

.PHONY: image
image:setenv
	buildctl --addr tcp://10.99.253.223:1234 build \
	--frontend dockerfile.v0 \
	--local context=. \
	--local dockerfile=. \
	--output type=image,registry.insecure=true,name=$(ImageID),push=true \
	--allow security.insecure \
	--opt filename=Dockerfile \
	--opt platform=linux/amd64 \
	--opt build-arg:flags=$(Flags)
