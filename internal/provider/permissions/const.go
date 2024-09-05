package permissions

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	datasourceMetadataName = "permissions"
	datasourceMetadataType = "data source"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &PermissionsDataSource{}
	_ datasource.DataSourceWithConfigure = &PermissionsDataSource{}
)

const (
	resourceMetadataName = "permissions"
	resourceMetadataType = "resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &PermissionsResource{}
	_ resource.ResourceWithConfigure   = &PermissionsResource{}
	_ resource.ResourceWithImportState = &PermissionsResource{}
)
