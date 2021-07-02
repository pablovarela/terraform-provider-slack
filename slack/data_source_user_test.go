package slack

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSlackUserDataSource_basic(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.slack_user.test"

	var providers []*schema.Provider

	t.Run("search non-existent user by name", func(t *testing.T) {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories(&providers),
			Steps: []resource.TestStep{
				{
					Config:      testAccCheckSlackUserDataSourceConfigNonExistentByName,
					ExpectError: regexp.MustCompile(`your query returned no results`),
				},
			},
		})
	})

	t.Run("search non-existent user by email", func(t *testing.T) {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories(&providers),
			Steps: []resource.TestStep{
				{
					Config:      testAccCheckSlackUserDataSourceConfigNonExistentByEmail,
					ExpectError: regexp.MustCompile(`users_not_found`),
				},
			},
		})
	})

	t.Run("search without setting any field", func(t *testing.T) {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories(&providers),
			Steps: []resource.TestStep{
				{
					Config:      testAccCheckSlackUserDataSourceConfigMissingFields,
					ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
				},
			},
		})
	})

	t.Run("search by name and email", func(t *testing.T) {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories(&providers),
			Steps: []resource.TestStep{
				{
					Config:      testAccCheckSlackUserDataSourceConfigExistentByNameAndEmail,
					ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
				},
			},
		})
	})

	t.Run("search by name", func(t *testing.T) {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories(&providers),
			Steps: []resource.TestStep{
				{
					Config: testAccCheckSlackUserDataSourceConfigExistentByName,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSlackUserDataSourceID(dataSourceName),
						resource.TestCheckResourceAttr(dataSourceName, "name", testUser00.name),
						resource.TestCheckResourceAttr(dataSourceName, "id", testUser00.id),
						resource.TestCheckResourceAttr(dataSourceName, "email", testUser00.email),
					),
				},
			},
		})
	})

	t.Run("search by email", func(t *testing.T) {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories(&providers),
			Steps: []resource.TestStep{
				{
					Config: testAccCheckSlackUserDataSourceConfigExistentByEmail,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSlackUserDataSourceID(dataSourceName),
						resource.TestCheckResourceAttr(dataSourceName, "name", testUser00.name),
						resource.TestCheckResourceAttr(dataSourceName, "id", testUser00.id),
						resource.TestCheckResourceAttr(dataSourceName, "email", testUser00.email),
					),
				},
			},
		})
	})
}

func testAccCheckSlackUserDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find slack conversation data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("slack conversation data source id not set")
		}
		return nil
	}
}

const (
	testAccCheckSlackUserDataSourceConfigNonExistentByName = `
data slack_user test {
 name = "non-existent"
}
`

	testAccCheckSlackUserDataSourceConfigNonExistentByEmail = `
data slack_user test {
 email = "non-existent"
}
`

	testAccCheckSlackUserDataSourceConfigMissingFields = `
data slack_user test {
}
`
)

var (
	testAccCheckSlackUserDataSourceConfigExistentByName = fmt.Sprintf(`
data slack_user test {
 name = "%s"
}
`, testUser00.name)

	testAccCheckSlackUserDataSourceConfigExistentByEmail = fmt.Sprintf(`
data slack_user test {
 email = "%s"
}
`, testUser00.email)

	testAccCheckSlackUserDataSourceConfigExistentByNameAndEmail = fmt.Sprintf(`
data slack_user test {
 name = "%s"
 email = "%s"
}
`, testUser00.name, testUser00.email)
)
