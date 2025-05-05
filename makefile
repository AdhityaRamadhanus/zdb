.PHONY: unit-test bench-bulk bench

unit-test:
	@go test -count=1 -v --cover ./... -tags="unit"

bench:
	@go test -count=1 -bench=. -benchtime=10s -v --cover ./... -tags="unit"