package role_members

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	datasourceMetadataName = "role_members"
	datasourceMetadataType = "data source"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &RoleMembersDataSource{}
	_ datasource.DataSourceWithConfigure = &RoleMembersDataSource{}
)

const (
	resourceMetadataName = "role_members"
	resourceMetadataType = "resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &RoleMembersResource{}
	_ resource.ResourceWithConfigure   = &RoleMembersResource{}
	_ resource.ResourceWithImportState = &RoleMembersResource{}
)
