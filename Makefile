SHELL := /bin/bash

.PHONY: setup dev build clean

setup:
	cd frontend && npm i
	conda create -n montaje python=3.11 -y
	conda activate montaje; pip install -r aiworker/requirements.txt

dev:
	wails dev -tags webkit2_41

build:
	wails build -tags webkit2_41

clean:
	rm -rf frontend/node_modules frontend/dist venv build