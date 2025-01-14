package acctest

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

const SkipNoResultsTest = true

var testAccProviders map[string]*schema.Provider
var testAccClassicProviders map[string]*schema.Provider

var testAccProvider *schema.Provider
var testAccClassicProvider *schema.Provider

var credsEnvVars = []string{
	"NCLOUD_ACCESS_KEY",
	"NCLOUD_SECRET_KEY",
}

var regionEnvVar = "NCLOUD_REGION"

func init() {
	testAccProvider = getTestAccProvider(true)
	testAccClassicProvider = getTestAccProvider(false)

	testAccProviders = map[string]*schema.Provider{
		"ncloud": testAccProvider,
	}

	testAccClassicProviders = map[string]*schema.Provider{
		"ncloud": testAccClassicProvider,
	}
}

func GetTestProvider(isVpc bool) *schema.Provider {
	if isVpc {
		return testAccProvider
	}

	return testAccClassicProvider
}

func GetTestAccProviders(isVpc bool) map[string]*schema.Provider {
	if isVpc {
		return testAccProviders
	}

	return testAccClassicProviders
}

func getTestAccProvider(isVpc bool) *schema.Provider {
	p := provider.Provider()
	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		d.Set("region", testAccGetRegion())
		d.Set("support_vpc", isVpc)
		return provider.ProviderConfigure(d)
	}
	return p
}

func TestAccPreCheck(t *testing.T) {
	if v := multiEnvSearch(credsEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(credsEnvVars, ", "))
	}

	region := testAccGetRegion()
	log.Printf("[INFO] Test: Using %s as test region", region)
}

func testAccGetRegion() string {
	v := os.Getenv(regionEnvVar)
	if v == "" {
		return "KR"
	}
	return v
}

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func GetTestPrefix() string {
	rand := acctest.RandString(5)
	return fmt.Sprintf("tf%s", rand)
}

func ComposeConfig(config ...string) string {
	var str strings.Builder

	for _, conf := range config {
		str.WriteString(conf)
	}

	return str.String()
}

func TestAccCheckResourceDisappears(provider *schema.Provider, resource *schema.Resource, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceState.Primary.ID == "" {
			return fmt.Errorf("resource ID missing: %s", resourceName)
		}

		return resource.Delete(resource.Data(resourceState.Primary), provider.Meta())
	}
}

func TestAccCheckDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("source ID not set")
		}
		return nil
	}
}

func TestAccValidateOneResult(t *testing.T) {
	if err := ValidateOneResult(0); err == nil {
		t.Fatalf("0 result must throw 'no results' error")
	}
	if err := ValidateOneResult(1); err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := ValidateOneResult(2); err == nil {
		t.Fatalf("2 results must throw 'more than one found results'")
	}
}

func GetTestClusterName() string {
	rInt := acctest.RandIntRange(1, 9999)
	testClusterName := fmt.Sprintf("tf-%d-cluster", rInt)
	return testClusterName
}
