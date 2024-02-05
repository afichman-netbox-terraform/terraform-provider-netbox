package netbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetboxManufacturerDataSource_basic(t *testing.T) {
	testSlug := "manufacturer_ds_basic"
	testName := testAccGetTestName(testSlug)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
			
resource "netbox_manufacturer" "test" {
  name = "%[1]s"
}
data "netbox_manufacturer" "test" {
  depends_on = [netbox_manufacturer.test]
  filter {
	  name = "%[1]s"
  }
}`, testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.netbox_manufacturer.test", "id", "netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_manufacturer.test", "slug", "netbox_manufacturer.test", "slug"),
				),
			},
		},
	})
}
