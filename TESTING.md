# Testing Guide for Terraform Provider for Venafi

This document describes the tests that are run to confirm the core functionality of the Terraform Provider for Venafi.

## Quick Start: Running Tests

### Step 1: Clone and Setup

```bash
git clone https://github.com/Venafi/terraform-provider-venafi.git
cd terraform-provider-venafi
go mod download
```

### Step 2: Run Unit Tests (No External Services Required)

```bash
# Run all unit tests
go test -v ./venafi

# Expected output:
# === RUN   TestProvider
# --- PASS: TestProvider (0.00s)
# === RUN   TestProvider_impl
# --- PASS: TestProvider_impl (0.00s)
# === RUN   TestNormalizedZones
# --- PASS: TestNormalizedZones (0.00s)
# === RUN   TestSetTLSConfig
# --- PASS: TestSetTLSConfig (0.01s)
# PASS
# ok      github.com/Venafi/terraform-provider-venafi/venafi
```

### Step 3: Run Dev Mode Tests (Requires TF_ACC=1)

```bash
TF_ACC=1 go test -v ./venafi -run "^TestDev" -timeout 10m

# Expected output shows certificate generation and validation
```

### Step 4: Run Integration Tests (Requires External Services)

See detailed instructions in the "Integration Tests" section below.

---

## Test Evidence

### Unit Tests Evidence

The following tests run successfully without any external dependencies:

```
=== RUN   TestProvider
--- PASS: TestProvider (0.00s)
=== RUN   TestProvider_impl
--- PASS: TestProvider_impl (0.00s)
=== RUN   TestNormalizedZones
--- PASS: TestNormalizedZones (0.00s)
=== RUN   TestSetTLSConfig
--- PASS: TestSetTLSConfig (0.01s)
PASS
ok      github.com/Venafi/terraform-provider-venafi/venafi    0.012s
```

### Available Tests Summary

| Category | Test Count | Build Tag Required | External Service Required |
|----------|------------|-------------------|--------------------------|
| Unit Tests | 4 | None | No |
| Dev Mode Tests | 5 | None | No (uses `TF_ACC=1`) |
| TPP Integration Tests | 26 | `tpp` | Yes (CyberArk Certificate Manager, Self-Hosted) |
| VaaS Integration Tests | 12 | `vaas` | Yes (CyberArk Certificate Manager, SaaS) |

**Total: 47 tests**

---

## Test Categories

The provider includes several categories of tests, each serving a specific purpose:

### 1. Unit Tests (No Build Tags Required)

These tests run without any external services and verify core functionality:

| Test Name | Location | Purpose |
|-----------|----------|---------|
| `TestProvider` | `venafi/provider_test.go` | Validates the provider schema and internal configuration |
| `TestProvider_impl` | `venafi/provider_test.go` | Verifies provider implementation |
| `TestNormalizedZones` | `venafi/provider_test.go` | Tests zone normalization for TPP paths |
| `TestSetTLSConfig` | `venafi/provider_test.go` | Validates TLS configuration with PKCS#12 certificates |

### 2. Dev Mode Tests (No External Services Required)

These tests use the provider's development mode to verify certificate functionality without connecting to CyberArk Certificate Manager:

| Test Name | Location | Purpose |
|-----------|----------|---------|
| `TestDevSignedCert` | `venafi/resource_venafi_certificate_test.go` | Tests RSA certificate generation in dev mode |
| `TestDevSignedCertECDSA` | `venafi/resource_venafi_certificate_test.go` | Tests ECDSA certificate generation in dev mode |
| `TestDevSignedCertWithCSRPem` | `venafi/resource_venafi_certificate_test.go` | Tests certificate enrollment with user-provided CSR |
| `TestDevSignedCertBackwardCompatibility` | `venafi/resource_venafi_certificate_test.go` | Verifies backward compatibility for local CSR generation |
| `TestDevCSRPemValidation` | `venafi/resource_venafi_certificate_test.go` | Tests CSR validation |

### 3. TPP Integration Tests (Requires `tpp` Build Tag)

These tests require a connection to CyberArk Certificate Manager, Self-Hosted:

**Certificate Tests** (`venafi/resource_venafi_certificate_tpp_test.go`):

| Test Name | Purpose |
|-----------|---------|
| `TestTPPSignedCert` | Tests basic RSA certificate enrollment |
| `TestTPPSignedCertWithNickname` | Tests certificate enrollment with custom nickname |
| `TestTPPECDSASignedCert` | Tests ECDSA certificate enrollment |
| `TestTPPSignedCertUpdate` | Tests certificate renewal functionality |
| `TestTPPTokenSignedCert` | Tests certificate enrollment using access token |
| `TestTPPTokenECDSASignedCert` | Tests ECDSA enrollment with access token |
| `TestTPPTokenSignedCertUpdateRenew` | Tests renewal with access token |
| `TestTPPSignedCertCustomFields` | Tests custom fields support |
| `TestTPPTokenSignedCertValidDays` | Tests valid_days parameter |
| `TestTPPTokenSignedCertUpdateSetGreaterExpWindow` | Tests expiration window updates |
| `TestTPPTppCsrService` | Tests server-generated CSR (service mode) |
| `TestTPPSansCsrService` | Tests SANs with service-generated CSR |
| `TestTPPImportCertificate` | Tests certificate import functionality |
| `TestTPPImportCertificateWithNickname` | Tests import with nickname |
| `TestTPPImportCertificateWithCustomFields` | Tests import with custom fields |
| `TestTPPImportCertificateECDSA` | Tests ECDSA certificate import |
| `TestTPPValidateWrongImportEntries` | Tests import validation error handling |
| `TestTPPManyCerts` | Stress tests with multiple certificates |
| `TestTPPManyCertsCsrService` | Stress tests with service-generated CSR |

**Policy Tests** (`venafi/resource_venafi_policy_tpp_test.go`):

| Test Name | Purpose |
|-----------|---------|
| `TestTPPCreateEmptyPolicy` | Tests creating an empty policy folder |
| `TestTPPCreatePolicy` | Tests creating a policy with specifications |
| `TestTPPImportPolicy` | Tests importing existing policy |

**SSH Certificate Tests** (`venafi/resource_venafi_ssh_certificate_test.go`):

| Test Name | Purpose |
|-----------|---------|
| `TestTPPSshCert` | Tests SSH certificate with service-generated key |
| `TestTPPSshCertNewAttrPrincipals` | Tests SSH cert with new principals attribute |
| `TestTPPSshCertLocalPublicKey` | Tests SSH cert with locally generated key |
| `TestTPPSshCertLocalPublicKeyNewAttrPrincipals` | Tests SSH cert local key with principals |

**SSH Config Tests** (`venafi/resource_venafi_ssh_config_test.go`):

| Test Name | Purpose |
|-----------|---------|
| `TestTPPSshConfig` | Tests SSH configuration retrieval |

### 4. VaaS Integration Tests (Requires `vaas` Build Tag)

These tests require a connection to CyberArk Certificate Manager, SaaS:

**Certificate Tests** (`venafi/resource_venafi_certificate_vaas_test.go`):

| Test Name | Purpose |
|-----------|---------|
| `TestVAASSignedCert` | Tests basic certificate enrollment |
| `TestVAASSignedCertWithDN` | Tests certificate with Distinguished Name |
| `TestVAASSignedCertWithUnacceptableDN` | Tests DN validation errors |
| `TestVAASSignedCertWithDNServiceGeneratedCSR` | Tests DN with service CSR |
| `TestVAASSignedCertWithUnacceptableDNServiceGeneratedCSR` | Tests DN validation with service CSR |
| `TestVAASSignedCertUpdateRenew` | Tests certificate renewal |
| `TestVAASSignedCertUpdateSetGreaterExpWindow` | Tests expiration window updates |
| `TestVAASImportCertificate` | Tests certificate import |
| `TestVAASCsrService` | Tests service-generated CSR |

**Policy Tests** (`venafi/resource_venafi_policy_vaas_test.go`):

| Test Name | Purpose |
|-----------|---------|
| `TestVAASCreateEmptyPolicy` | Tests creating empty application/CIT |
| `TestVAASCreatePolicy` | Tests creating policy with specifications |
| `TestVAASImportPolicy` | Tests importing existing policy |

## Running Tests

### Unit Tests (No External Services)

```bash
# Run all unit tests
go test -v ./venafi

# Run specific test
go test -v ./venafi -run TestProvider
```

### Dev Mode Tests

Dev mode tests run automatically with unit tests since they don't require external connections:

```bash
go test -v ./venafi -run "^TestDev"
```

### TPP Integration Tests

Requires environment variables to be set:

```bash
export TPP_URL="https://your-tpp-server.example.com"
export TPP_USER="your-username"
export TPP_PASSWORD="your-password"
export TPP_ZONE="YourPolicy\\Folder"
export TPP_ZONE_ECDSA="YourPolicy\\ECDSAFolder"
export TPP_ACCESS_TOKEN="your-access-token"
export TRUST_BUNDLE="/path/to/bundle.pem"

# Run TPP tests
TF_ACC=1 go test -tags=tpp -run "^TestTPP" ./venafi -v -timeout 120m
```

### VaaS Integration Tests

Requires environment variables to be set:

```bash
export CLOUD_URL="https://api.venafi.cloud"
export CLOUD_APIKEY="your-api-key"
export CLOUD_ZONE="YourApp\\YourCIT"

# Run VaaS tests
TF_ACC=1 go test -tags=vaas -run "^TestVAAS" ./venafi -v -timeout 120m
```

### Using Makefile Targets

The Makefile provides convenient targets for running tests:

| Target | Description |
|--------|-------------|
| `make test_go` | Runs unit tests with coverage |
| `make testacc` | Runs all acceptance tests |
| `make test_tpp_tf_acc` | Runs TPP acceptance tests |
| `make test_vaas_tf_acc` | Runs VaaS acceptance tests |
| `make test` | Runs full test suite (fmtcheck, linter, test_go, testacc, e2e) |
| `make test_tpp` | Runs all TPP-related tests |
| `make test_vaas` | Runs all VaaS-related tests |
| `make test_e2e` | Runs end-to-end tests with real Terraform binary |

## E2E Tests

End-to-end tests use a real Terraform binary to test the provider:

| Target | Description |
|--------|-------------|
| `make test_e2e_dev` | Tests dev mode with real Terraform |
| `make test_e2e_dev_ecdsa` | Tests dev mode ECDSA with real Terraform |
| `make test_e2e_tpp` | Tests TPP with real Terraform |
| `make test_e2e_vaas` | Tests VaaS with real Terraform |
| `make test_e2e_tpp_token` | Tests TPP token auth with real Terraform |

## How Core Functionality is Confirmed

### 1. Certificate Enrollment

The core certificate enrollment functionality is verified through:

- **Certificate Generation**: Tests verify that certificates are properly generated with correct CN, SANs, key algorithm (RSA/ECDSA)
- **Private Key Management**: Tests confirm private keys are encrypted with the provided password
- **Certificate Validation**: Each test parses the generated certificate and validates:
  - Common Name matches requested
  - DNS SANs are correct
  - Key type matches specification
  - Certificate can be paired with private key

### 2. Certificate Renewal

The renewal functionality is tested by:

- Setting `expiration_window` equal to certificate validity
- Running terraform apply twice
- Verifying serial numbers differ (new certificate issued)

### 3. CSR Origin Options

Tests verify all three CSR generation methods:
- `local`: Key pair generated locally by the provider
- `service`: Key pair generated by the CyberArk platform
- `file`: User provides their own CSR (external key management)

### 4. Import Functionality

Import tests verify:

- Existing certificates can be imported into Terraform state
- Certificate attributes are correctly populated
- Private key can be retrieved and decrypted
- Custom fields are preserved

### 5. Policy Management

Policy tests verify:

- Policies can be created from JSON specification
- Policy attributes match the specification file
- Existing policies can be imported

### 6. SSH Certificates

SSH certificate tests verify:

- SSH certificates are properly generated
- Principals are correctly assigned
- CA public key is retrieved correctly
- Both service and local key generation work

## Test Utilities

The `venafi/test_util.go` file provides shared test utilities:

- `checkStandardCert`: Validates certificate attributes
- `checkStandardCertInfo`: Parses and validates certificate content
- `checkCertValidDays`: Verifies certificate validity period
- `createCertificate`: Creates test certificates using vcert library
- `randSeq`: Generates random strings for unique test data

## Coverage

To generate test coverage:

```bash
make test_go
# This outputs coverage to cov1.out and displays coverage summary
```

Or manually:

```bash
go test -v -coverprofile=coverage.out ./venafi
go tool cover -func=coverage.out
```

---

## Complete Test List

### Tests Without Build Tags (9 tests)

```
TestProvider
TestProvider_impl
TestNormalizedZones
TestSetTLSConfig
TestDevSignedCert
TestDevSignedCertECDSA
TestDevSignedCertWithCSRPem
TestDevSignedCertBackwardCompatibility
TestDevCSRPemValidation
```

### Tests With `tpp` Build Tag (35 tests total, 26 TPP-specific)

```
TestTPPSignedCertUpdate
TestTPPSignedCert
TestTPPSignedCertWithNickname
TestTPPECDSASignedCert
TestTPPTokenSignedCertUpdateRenew
TestTPPTokenSignedCert
TestTPPTokenECDSASignedCert
TestTPPSignedCertCustomFields
TestTPPTokenSignedCertValidDays
TestTPPTokenSignedCertUpdateSetGreaterExpWindow
TestTPPTppCsrService
TestTPPValidateWrongImportEntries
TestTPPImportCertificate
TestTPPImportCertificateWithNickname
TestTPPImportCertificateWithCustomFields
TestTPPImportCertificateECDSA
TestTPPManyCerts
TestTPPManyCertsCsrService
TestTPPSansCsrService
TestTPPCreateEmptyPolicy
TestTPPCreatePolicy
TestTPPImportPolicy
TestTPPSshCert
TestTPPSshCertNewAttrPrincipals
TestTPPSshCertLocalPublicKey
TestTPPSshCertLocalPublicKeyNewAttrPrincipals
TestTPPSshConfig
```

### Tests With `vaas` Build Tag (21 tests total, 12 VaaS-specific)

```
TestVAASSignedCert
TestVAASSignedCertWithDN
TestVAASSignedCertWithUnacceptableDN
TestVAASSignedCertWithDNServiceGeneratedCSR
TestVAASSignedCertWithUnacceptableDNServiceGeneratedCSR
TestVAASSignedCertUpdateRenew
TestVAASSignedCertUpdateSetGreaterExpWindow
TestVAASImportCertificate
TestVAASCsrService
TestVAASCreateEmptyPolicy
TestVAASCreatePolicy
TestVAASImportPolicy
```

---

## Verification Steps

### How to Verify Core Certificate Functionality Works

1. **Run Dev Mode Test** (tests certificate generation without external services):
   ```bash
   TF_ACC=1 go test -v ./venafi -run TestDevSignedCert -timeout 5m
   ```

2. **Expected Behavior**:
   - Provider initializes in dev mode
   - Certificate is generated with RSA 2048-bit key
   - Certificate contains correct Common Name and SANs
   - Private key is properly encrypted
   - Certificate and private key form valid keypair

3. **Sample Test Output**:
   ```
   === RUN   TestDevSignedCert
       resource_venafi_certificate_test.go:45: Testing Dev RSA certificate
       test_util.go:110: Testing certificate with cn dev-random.venafi.example.com
       test_util.go:117: Testing certificate PEM:
        -----BEGIN CERTIFICATE-----
        MIID9jCCAt6gAwIBAgIQ...
        -----END CERTIFICATE-----
       test_util.go:255: Testing private key PEM:
        -----BEGIN PRIVATE KEY-----
        MIIEvQIBADANBgkqhkiG9w0B...
        -----END PRIVATE KEY-----
   --- PASS: TestDevSignedCert (1.53s)
   ```

### How to Verify TPP Integration Works

1. **Set up environment variables**:
   ```bash
   export TPP_URL="https://your-tpp-server.example.com"
   export TPP_ACCESS_TOKEN="your-access-token"
   export TPP_ZONE="YourPolicy\\Folder"
   export TPP_ZONE_ECDSA="YourPolicy\\ECDSAFolder"
   export TRUST_BUNDLE="/path/to/bundle.pem"
   ```

2. **Run TPP tests**:
   ```bash
   TF_ACC=1 go test -tags=tpp -run "^TestTPPSignedCert$" ./venafi -v -timeout 120m
   ```

3. **Verify in TPP UI**:
   - Log into TPP Web Admin
   - Navigate to the policy folder
   - Verify certificate object was created
   - Check certificate details match test parameters

### How to Verify VaaS Integration Works

1. **Set up environment variables**:
   ```bash
   export CLOUD_URL="https://api.venafi.cloud"
   export CLOUD_APIKEY="your-api-key"
   export CLOUD_ZONE="YourApp\\YourCIT"
   ```

2. **Run VaaS tests**:
   ```bash
   TF_ACC=1 go test -tags=vaas -run "^TestVAASSignedCert$" ./venafi -v -timeout 120m
   ```

3. **Verify in VaaS UI**:
   - Log into CyberArk Certificate Manager, SaaS
   - Navigate to the application
   - Verify certificate was issued
   - Check certificate details match test parameters

---

## Test File Locations

| File | Purpose |
|------|---------|
| `venafi/provider_test.go` | Provider unit tests |
| `venafi/resource_venafi_certificate_test.go` | Dev mode certificate tests |
| `venafi/resource_venafi_certificate_tpp_test.go` | TPP certificate integration tests |
| `venafi/resource_venafi_certificate_vaas_test.go` | VaaS certificate integration tests |
| `venafi/resource_venafi_policy_tpp_test.go` | TPP policy integration tests |
| `venafi/resource_venafi_policy_vaas_test.go` | VaaS policy integration tests |
| `venafi/resource_venafi_ssh_certificate_test.go` | TPP SSH certificate tests |
| `venafi/resource_venafi_ssh_config_test.go` | TPP SSH config tests |
| `venafi/test_util.go` | Shared test utilities |

---

## Makefile Test Targets

| Target | Command | Description |
|--------|---------|-------------|
| `test_go` | `go test -v -coverprofile=cov1.out ./venafi` | Unit tests with coverage |
| `testacc` | `TF_ACC=1 go test $(TEST) -v -timeout 120m` | All acceptance tests |
| `test_tpp_tf_acc` | `TF_ACC=1 go test -tags=tpp -run ^TestTPP ./venafi -v -timeout 120m` | TPP acceptance tests |
| `test_vaas_tf_acc` | `TF_ACC=1 go test -tags=vaas -run ^TestVAAS ./venafi -v -timeout 120m` | VaaS acceptance tests |
| `test_e2e` | Multiple terraform apply commands | End-to-end with real Terraform |
| `linter` | `golangci-lint run` | Code linting |
| `fmtcheck` | `scripts/gofmtcheck.sh` | Format checking |
