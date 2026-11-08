.PHONY: lint architecture test vet build verify

lint:
	python3 lab/ferramentas/crystalline-lint/run.py --fail-on warning .

architecture:
	go test ./lab/bancada/arquitetura/0129/boundaries ./lab/bancada/arquitetura/0131/contracts ./lab/bancada/arquitetura/0145

test:
	go test ./...

vet:
	go vet ./...

build:
	go build -o /tmp/bateia ./04_wiring

verify:
	python3 lab/ferramentas/verify.py
