package guild

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

const (
	datasourceMetadataName = "guild"
	datasourceMetadataType = "data source"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GuildDataSource{}
	_ datasource.DataSourceWithConfigure = &GuildDataSource{}
)
