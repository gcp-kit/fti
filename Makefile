.PHONY: test
test:
	go test ./...

.PHONY: gen_samples
gen_samples:
	go run ./cmd/fti -c sample/config.yaml

.PHONY: emulator
emulator:
	npm run emulator
