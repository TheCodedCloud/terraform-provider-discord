package channel

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	datasourceMetadataName = "channel"
	datasourceMetadataType = "data source"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ChannelDataSource{}
	_ datasource.DataSourceWithConfigure = &ChannelDataSource{}
)

const (
	resourceMetadataName = "channel"
	resourceMetadataType = "resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ChannelResource{}
	_ resource.ResourceWithConfigure   = &ChannelResource{}
	_ resource.ResourceWithImportState = &ChannelResource{}
)
