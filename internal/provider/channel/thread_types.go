package channel

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ThreadMetadata struct {
	// Whether the thread is archived.
	Archived types.Bool `tfsdk:"archived"`

	// The thread will stop showing in the channel list after auto_archive_duration minutes of inactivity, can be set to: 60, 1440, 4320, 10080
	AutoArchiveDuration types.Int32 `tfsdk:"auto_archive_duration"`

	// Timestamp when the thread's archive status was last changed, used for calculating recent activity
	ArchiveTimestamp *time.Time `tfsdk:"archive_timestamp"`

	// Whether the thread is locked; when a thread is locked, only users with MANAGE_THREADS can unarchive it
	Locked types.Bool `tfsdk:"locked"`

	// Whether non-moderators can add other non-moderators to a thread; only available on private threads
	Invitable types.Bool `tfsdk:"invitable"`

	// Timestamp when the thread was created; only populated for threads created after 2022-01-09
	CreateTimestamp *time.Time `tfsdk:"create_timestamp"`
}

var ThreadMetadataSchema = map[string]schema.Attribute{
	"archived": schema.BoolAttribute{
		Optional: true,
		Computed: true,
	},
	"auto_archive_duration": schema.Int32Attribute{
		Optional: true,
		Computed: true,
	},
	"archive_timestamp": schema.StringAttribute{
		Computed: true,
	},
	"locked": schema.BoolAttribute{
		Optional: true,
		Computed: true,
	},
	"invitable": schema.BoolAttribute{
		Optional: true,
		Computed: true,
	},
	"create_timestamp": schema.StringAttribute{
		Computed: true,
	},
}

type ThreadMember struct {
	// ID of the thread
	ID types.String `tfsdk:"id"`

	// ID of the user
	UserID types.String `tfsdk:"user_id"`

	// Time the user last joined the thread
	JoinTimestamp *time.Time `tfsdk:"join_timestamp"`

	// Any user-thread settings, currently only used for notifications
	Flags types.String `tfsdk:"flags"`

	// Additional information about the user
	// Member *User `tfsdk:"member"`
}
