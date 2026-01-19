
build: cmd/sample-go-app/main.go
	go build ./...

clean::
	rm sample-go-app

security: trivy

trivy:
	trivy fs .

gofmt:
	gofmt -w .
GOSEC_VERSION ?= latest
GOSV_VERSION ?= latest
GITLEAKS_VERSION ?= latest
# Install/Update Security Tools (Runs before scan targets)
tools::
	go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION)
	go install github.com/google/osv-scanner/v2/cmd/osv-scanner@$(GOSV_VERSION)
	go install github.com/gitleaks/gitleaks/v8@$(GITLEAKS_VERSION)

# Static Analysis Security Testing (SAST)
scan-sast: tools
	@echo "--- Running Gosec (SAST) ---"
	gosec -fmt=json -out=gosec-report.json ./...
	@echo "Gosec report generated: gosec-report.json"

# Dependency Vulnerability Scanning
scan-deps: tools
	@echo "--- Running OSV-Scanner (Dependency Scan) ---"
	osv-scanner --format=json --output=osv-report.json ./...
	@echo "OSV-Scanner report generated: osv-report.json"

# Secrets Scanning
scan-secrets: tools
	@echo "--- Running GitLeaks (Secrets Scan) ---"
	gitleaks detect --verbose --report-path=gitleaks-report.json
	@echo "GitLeaks report generated: gitleaks-report.json"

# Combined Security Scan (Runs all checks)
scan: scan-sast scan-deps scan-secrets

# Clean up generated reports
clean::
	rm -f gosec-report.json osv-report.json gitleaks-report.json
