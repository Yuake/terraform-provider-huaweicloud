package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccImsImageDataSource_basic(t *testing.T) {
	imageName := "CentOS 7.4 64bit"
	dataSourceName := "data.huaweicloud_images_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImageDataSource_publicName(imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", imageName),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
				),
			},
			{
				Config: testAccImsImageDataSource_osVersion(imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
				),
			},
			{
				Config: testAccImsImageDataSource_nameRegex("^CentOS 7.4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
				),
			},
		},
	})
}

func TestAccImsImageDataSource_testQueries(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_images_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImageDataSource_base(rName),
			},
			{
				Config: testAccImsImageDataSource_queryName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "private"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
				),
			},
			{
				Config: testAccImsImageDataSource_queryTag(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
				),
			},
		},
	})
}

func testAccCheckImagesV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Image data source ID not set")
		}

		return nil
	}
}

func testAccImsImageDataSource_publicName(imageName string) string {
	return fmt.Sprintf(`
data "huaweicloud_images_image" "test" {
  name        = "%s"
  visibility  = "public"
  most_recent = true
}
`, imageName)
}

func testAccImsImageDataSource_nameRegex(regexp string) string {
	return fmt.Sprintf(`
data "huaweicloud_images_image" "test" {
  architecture = "x86"
  name_regex   = "%s"
  visibility   = "public"
  most_recent  = true
}
`, regexp)
}

func testAccImsImageDataSource_osVersion(osVersion string) string {
	return fmt.Sprintf(`
data "huaweicloud_images_image" "test" {
  architecture = "x86"
  os_version   = "%s"
  visibility   = "public"
  most_recent  = true
}
`, osVersion)
}

func testAccImsImageDataSource_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_compute_instance" "test" {
  name              = "%s"
  image_name        = "Ubuntu 18.04 server 64bit"
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_images_image" "test" {
  name        = "%s"
  instance_id = huaweicloud_compute_instance.test.id
  description = "created by Terraform AccTest"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName, rName)
}

func testAccImsImageDataSource_queryName(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
  most_recent = true
  name        = huaweicloud_images_image.test.name
}
`, testAccImsImageDataSource_base(rName))
}

func testAccImsImageDataSource_queryTag(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
  most_recent = true
  visibility  = "private"
  tag         = "foo=bar"
}
`, testAccImsImageDataSource_base(rName))
}
