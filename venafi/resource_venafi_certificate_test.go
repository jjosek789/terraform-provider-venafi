package venafi

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	devConfig = `
provider "venafi" {
	alias = "dev"
	dev_mode = true
}
resource "venafi_certificate" "dev_certificate" {
	provider = "venafi.dev"
	common_name = "%s"
	%s
	san_dns = [
		"%s"
	]
	san_ip = [
		"10.1.1.1",
		"192.168.0.1"
	]
	san_email = [
		"dev@venafi.com",
		"dev2@venafi.com"
	]
}
output "certificate" {
	value = "${venafi_certificate.dev_certificate.certificate}"
}
output "private_key" {
	value = "${venafi_certificate.dev_certificate.private_key_pem}"
	sensitive = true
}`
)

func TestDevSignedCert(t *testing.T) {
	t.Log("Testing Dev RSA certificate")
	data := testData{}
	data.cn = "dev-random.venafi.example.com"
	data.dns_ns = "dev-web01-random.example.com"
	data.key_algo = rsa2048
	config := fmt.Sprintf(devConfig, data.cn, data.key_algo, data.dns_ns)
	t.Logf("Testing dev certificate with config:\n %s", config)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					err := checkStandardCert(t, &data, s)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	})
}

func TestDevSignedCertECDSA(t *testing.T) {
	t.Log("Testing Dev ECDSA certificate")
	data := testData{}
	data.cn = "dev-random.venafi.example.com"
	data.dns_ns = "dev-web01-random.example.com"
	data.key_algo = ecdsa521
	config := fmt.Sprintf(devConfig, data.cn, data.key_algo, data.dns_ns)
	t.Logf("Testing dev certificate with config:\n %s", config)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					err := checkStandardCert(t, &data, s)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	})
}

var devConfigWithCSRFile = `
provider "venafi" {
alias = "dev"
dev_mode = true
}
resource "venafi_certificate" "dev_certificate_csr_file" {
provider = "venafi.dev"
common_name = "test.venafi.example.com"
csr_origin = "file"
csr_file = "%s"
}
output "certificate" {
value = "${venafi_certificate.dev_certificate_csr_file.certificate}"
}
`

func TestDevSignedCertWithCSRFile(t *testing.T) {
	t.Log("Testing Dev certificate with user-provided CSR file")
	csrPath := "../test_files/test-csr.pem"
	config := fmt.Sprintf(devConfigWithCSRFile, csrPath)
	t.Logf("Testing dev certificate with CSR file config:\n %s", config)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					// Check that certificate was created
					gotUntyped := s.RootModule().Resources["venafi_certificate.dev_certificate_csr_file"]
					if gotUntyped == nil {
						return fmt.Errorf("resource not found in state")
					}

					got := gotUntyped.Primary
					if got == nil {
						return fmt.Errorf("primary instance not found")
					}

					// Verify certificate is present
					cert := got.Attributes["certificate"]
					if cert == "" {
						return fmt.Errorf("certificate attribute is empty")
					}

					// Verify chain is present
					chain := got.Attributes["chain"]
					if chain == "" {
						return fmt.Errorf("chain attribute is empty")
					}

					// Verify private_key_pem is NOT present (managed externally)
					privateKey := got.Attributes["private_key_pem"]
					if privateKey != "" {
						return fmt.Errorf("private_key_pem should not be stored for user-provided CSR, but got: %s", privateKey)
					}

					t.Logf("Certificate with user-provided CSR successfully created")
					return nil
				},
			},
		},
	})
}
