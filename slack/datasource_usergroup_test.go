package slack

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/slack-go/slack"
)

func TestAccSlackUserGroupDataSource_basic(t *testing.T) {
	resourceName := "slack_usergroup.test"
	dataSourceName := "data.slack_usergroup.test"

	var providers []*schema.Provider
	t.Run("search non-existent by ID", func(t *testing.T) {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories(&providers),
			Steps: []resource.TestStep{
				{
					Config:      testAccCheckSlackUserGroupDataSourceConfigNonExistentID,
					ExpectError: regexp.MustCompile(`could not find usergroup`),
				},
			},
		})
	})

	t.Run("search non-existent by name", func(t *testing.T) {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories(&providers),
			Steps: []resource.TestStep{
				{
					Config:      testAccCheckSlackUserGroupDataSourceConfigNonExistentName,
					ExpectError: regexp.MustCompile(`could not find usergroup`),
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
					Config:      testAccCheckSlackUserGroupDataSourceConfigMissingFields,
					ExpectError: regexp.MustCompile("ExactlyOne"),
				},
			},
		})
	})

	t.Run("search by name and ID", func(t *testing.T) {
		name := acctest.RandomWithPrefix(userGroupResourceNamePrefix)
		users := []string{testUser00.id, testUser01.id}
		channel := createTestConversation(t)
		createUserGroup := testAccSlackUserGroupWithUsers(name, []string{channel.ID}, users)
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories(&providers),
			Steps: []resource.TestStep{
				{
					Config:      testAccCheckSlackUserGroupDataSourceConfigByNameAndID(createUserGroup),
					ExpectError: regexp.MustCompile("ExactlyOne"),
				},
			},
		})
	})

	t.Run("search by ID", func(t *testing.T) {
		name := acctest.RandomWithPrefix(userGroupResourceNamePrefix)
		users := []string{testUser00.id, testUser01.id}
		channel := createTestConversation(t)
		createUserGroup := testAccSlackUserGroupWithUsers(name, []string{channel.ID}, users)
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories(&providers),
			Steps: []resource.TestStep{
				{
					Config: testAccCheckSlackUserGroupDataSourceConfigByID(createUserGroup),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSlackUserGroupDataSourceID(dataSourceName),
						resource.TestCheckResourceAttrPair(dataSourceName, "usergroup_id", resourceName, "id"),
						resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
						resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
						resource.TestCheckResourceAttrPair(dataSourceName, "handle", resourceName, "handle"),
						resource.TestCheckResourceAttrPair(dataSourceName, "users", resourceName, "users"),
						resource.TestCheckResourceAttrPair(dataSourceName, "channels", resourceName, "channels"),
					),
				},
			},
		})
	})

	t.Run("search by name", func(t *testing.T) {
		name := acctest.RandomWithPrefix(userGroupResourceNamePrefix)
		users := []string{testUser00.id, testUser01.id}
		channel := createTestConversation(t)
		createUserGroup := testAccSlackUserGroupWithUsers(name, []string{channel.ID}, users)
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories(&providers),
			Steps: []resource.TestStep{
				{
					Config: testAccCheckSlackUserGroupDataSourceConfigByName(createUserGroup),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSlackUserGroupDataSourceID(dataSourceName),
						resource.TestCheckResourceAttrPair(dataSourceName, "usergroup_id", resourceName, "id"),
						resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
						resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
						resource.TestCheckResourceAttrPair(dataSourceName, "handle", resourceName, "handle"),
						resource.TestCheckResourceAttrPair(dataSourceName, "users", resourceName, "users"),
						resource.TestCheckResourceAttrPair(dataSourceName, "channels", resourceName, "channels"),
					),
				},
			},
		})
	})
}

func testAccCheckSlackUserGroupDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find slack usergroup data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("slack usergroup data source id not set")
		}
		return nil
	}
}

const (
	testAccCheckSlackUserGroupDataSourceConfigNonExistentID = `
data slack_usergroup test {
 usergroup_id = "non-existent"
}
`
	testAccCheckSlackUserGroupDataSourceConfigNonExistentName = `
data slack_usergroup test {
 name = "non-existent"
}
`

	testAccCheckSlackUserGroupDataSourceConfigMissingFields = `
data slack_usergroup test {
}
`
	testAccCheckSlackUserGroupDataSourceConfigExistentIDAndName = `
data slack_usergroup test {
  usergroup_id = slack_usergroup.test.id
  name = slack_usergroup.test.name
}
`

	testAccCheckSlackUserGroupDataSourceConfigExistentID = `
data slack_usergroup test {
  usergroup_id = slack_usergroup.test.id
}
`

	testAccCheckSlackUserGroupDataSourceConfigExistentName = `
data slack_usergroup test {
  name = slack_usergroup.test.name
}
`
)

func testAccCheckSlackUserGroupDataSourceConfigByID(group slack.UserGroup) string {
	return testAccSlackUserGroupConfig(group) + testAccCheckSlackUserGroupDataSourceConfigExistentID
}

func testAccCheckSlackUserGroupDataSourceConfigByName(group slack.UserGroup) string {
	return testAccSlackUserGroupConfig(group) + testAccCheckSlackUserGroupDataSourceConfigExistentName
}

func testAccCheckSlackUserGroupDataSourceConfigByNameAndID(group slack.UserGroup) string {
	return testAccSlackUserGroupConfig(group) + testAccCheckSlackUserGroupDataSourceConfigExistentIDAndName
}
