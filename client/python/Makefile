#
#   Copyright (c) 2021, 2121, ducesoft and/or its affiliates. All rights reserved.
#   DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#

.PHONY: pub

.PHONY: pub
pub:
	$(eval Pyun:=$(or $(PYUN),root))
	$(eval Pypd:=$(or $(PYPD),Abc12345))
	$(eval Pypi:=$(or $(PYPI),http://10.12.0.78:13141/))
	@pip3 install poetry==1.2.0 -i https://mirrors.aliyun.com/pypi/simple/
	@poetry config repositories.mesh $(Pypi)root/dev/+simple/
	@poetry config http-basic.mesh $(Pyun) $(Pypd)
	@poetry build
	#@poetry publish -u $(Pyun) -p $(Pypd) -r mesh
	pip3 install devpi-client -i https://mirrors.aliyun.com/pypi/simple/
	devpi use -l $(Pypi)
	devpi login $(Pyun) --password=$(Pypd)
	devpi use root/dev/+simple/
	devpi upload --from-dir dist
	#devpi push example==mesh-0.0.1 root/dev

.PHONY: venv
venv:
	python3 -m venv .venv && source .venv/bin/activate && pip3 install -r requirements.txt

.PHONY: token
token:
	poetry config pypi-token.pypi pypi-AgEIcHlwaS5vcmcCJDczMjBjNzU1LTAxYTYtNDgyNy1hOTAxLTgyYmZjYjhlZDU5YQACKlszLCJjNzE5OGFlNi0xNWRlLTQyZmEtYjNhYS05ZmQ1MzNmMWQwZTIiXQAABiAO7SSEFhjcWGlpCC5H4lTh0M8K9LlssDwsojjBmkK1-w

.PHONY: deploy
deploy:
	poetry publish --build
