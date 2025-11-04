.PHONY: format lint tidy vulncheck secrets secscan precommit

format:
	go fmt ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

vulncheck:
	govulncheck ./...

secrets:
	gitleaks detect --source . --no-git --report-path gitleaks.report.json || true

secscan:
	gosec ./...

precommit: format lint tidy vulncheck secrets secscan
