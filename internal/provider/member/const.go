package member

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

const (
	datasourceMetadataName = "member"
	datasourceMetadataType = "data source"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &MemberDataSource{}
	_ datasource.DataSourceWithConfigure = &MemberDataSource{}
)

const (
	resourceMetadataName = "member"
	resourceMetadataType = "resource"
)

// // Ensure the implementation satisfies the expected interfaces.
// var (
// 	_ resource.Resource                = &RoleResource{}
// 	_ resource.ResourceWithConfigure   = &RoleResource{}
// 	_ resource.ResourceWithImportState = &RoleResource{}
// )
