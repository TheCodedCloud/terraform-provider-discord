package webhook

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	datasourceMetadataName = "webhook"
	datasourceMetadataType = "data source"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &WebhookDataSource{}
	_ datasource.DataSourceWithConfigure = &WebhookDataSource{}
)

const (
	resourceMetadataName = "webhook"
	resourceMetadataType = "resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &WebhookResource{}
	_ resource.ResourceWithConfigure   = &WebhookResource{}
	_ resource.ResourceWithImportState = &WebhookResource{}
)
