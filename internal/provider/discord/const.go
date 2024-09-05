package discord

import "github.com/bwmarrin/discordgo"

const (
	UNIMPLEMENTED = "UNIMPLEMENTED"
)

// Guild constants
var (
	MessageNotifications = map[discordgo.MessageNotifications]string{
		discordgo.MessageNotificationsAllMessages:  "ALL_MESSAGES",
		discordgo.MessageNotificationsOnlyMentions: "ONLY_MENTIONS",
	}

	ExplicitContentFilters = map[discordgo.ExplicitContentFilterLevel]string{
		discordgo.ExplicitContentFilterDisabled:            "DISABLED",
		discordgo.ExplicitContentFilterMembersWithoutRoles: "MEMBERS_WITHOUT_ROLES",
		discordgo.ExplicitContentFilterAllMembers:          "ALL_MEMBERS",
	}

	MFALevels = map[discordgo.MfaLevel]string{
		discordgo.MfaLevelNone:     "NONE",
		discordgo.MfaLevelElevated: "ELEVATED",
	}

	VerificationLevels = map[discordgo.VerificationLevel]string{
		discordgo.VerificationLevelNone:     "NONE",
		discordgo.VerificationLevelLow:      "LOW",
		discordgo.VerificationLevelMedium:   "MEDIUM",
		discordgo.VerificationLevelHigh:     "HIGH",
		discordgo.VerificationLevelVeryHigh: "VERY_HIGH",
	}

	GuildNSFWLevels = map[discordgo.GuildNSFWLevel]string{
		discordgo.GuildNSFWLevelDefault:       "DEFAULT",
		discordgo.GuildNSFWLevelExplicit:      "EXPLICIT",
		discordgo.GuildNSFWLevelSafe:          "SAFE",
		discordgo.GuildNSFWLevelAgeRestricted: "AGE_RESTRICTED",
	}

	PremiumTiers = map[discordgo.PremiumTier]string{
		discordgo.PremiumTierNone: "NONE",
		discordgo.PremiumTier1:    "TIER_1",
		discordgo.PremiumTier2:    "TIER_2",
		discordgo.PremiumTier3:    "TIER_3",
	}

	SystemChannelFlags = map[string]discordgo.SystemChannelFlag{
		"SUPPRESS_JOIN_NOTIFICATIONS":                              discordgo.SystemChannelFlagsSuppressJoinNotifications,
		"SUPPRESS_PREMIUM_SUBSCRIPTIONS":                           discordgo.SystemChannelFlagsSuppressPremium,
		"SUPPRESS_GUILD_REMINDER_NOTIFICATIONS":                    discordgo.SystemChannelFlagsSuppressGuildReminderNotifications,
		"SUPPRESS_JOIN_NOTIFICATION_REPLIES":                       discordgo.SystemChannelFlagsSuppressJoinNotificationReplies,
		"SUPPRESS_ROLE_SUBSCRIPTION_PURCHASE_NOTIFICATIONS":        1 << 4,
		"SUPPRESS_ROLE_SUBSCRIPTION_PURCHASE_NOTIFICATION_REPLIES": 1 << 5,
	}

	// Unused
	GuildFeatures = map[string]string{
		string(discordgo.GuildFeatureAnimatedBanner):                "ANIMATED_BANNER",
		string(discordgo.GuildFeatureAnimatedIcon):                  "ANIMATED_ICON",
		"APPLICATION_COMMAND_PERMISSIONS_V2":                        "APPLICATION_COMMAND_PERMISSIONS_V2",
		string(discordgo.GuildFeatureAutoModeration):                "AUTO_MODERATION",
		string(discordgo.GuildFeatureBanner):                        "BANNER",
		string(discordgo.GuildFeatureCommunity):                     "COMMUNITY", // mutable
		"CREATOR_MONETIZABLE_PROVISIONAL":                           "CREATOR_MONETIZABLE_PROVISIONAL",
		"CREATOR_STORE_PAGE":                                        "CREATOR_STORE_PAGE",
		"DEVELOPER_SUPPORT_SERVER":                                  "DEVELOPER_SUPPORT_SERVER",
		string(discordgo.GuildFeatureDiscoverable):                  "DISCOVERABLE", // mutable
		string(discordgo.GuildFeatureFeaturable):                    "FEATURABLE",
		"INVITES_DISABLED":                                          "INVITES_DISABLED", // mutable
		string(discordgo.GuildFeatureInviteSplash):                  "INVITE_SPLASH",
		string(discordgo.GuildFeatureMemberVerificationGateEnabled): "MEMBER_VERIFICATION_GATE_ENABLED",
		// string(discordgo.GuildFeatureMonetizationEnabled):           "MONETIZATION_ENABLED",
		string(discordgo.GuildFeatureMoreStickers):   "MORE_STICKERS",
		string(discordgo.GuildFeatureNews):           "NEWS",
		string(discordgo.GuildFeaturePartnered):      "PARTNERED",
		string(discordgo.GuildFeaturePreviewEnabled): "PREVIEW_ENABLED",
		// string(discordgo.GuildFeaturePrivateThreads):        "PRIVATE_THREADS",
		"RAID_ALERTS_DISABLED":                              "RAID_ALERTS_DISABLED", // mutable
		string(discordgo.GuildFeatureRoleIcons):             "ROLE_ICONS",
		"ROLE_SUBSCRIPTIONS_AVAILABLE_FOR_PURCHASE":         "ROLE_SUBSCRIPTIONS_AVAILABLE_FOR_PURCHASE",
		"ROLE_SUBSCRIPTIONS_ENABLED":                        "ROLE_SUBSCRIPTIONS_ENABLED",
		string(discordgo.GuildFeatureTicketedEventsEnabled): "TICKETED_EVENTS_ENABLED",
		string(discordgo.GuildFeatureVanityURL):             "VANITY_URL",
		string(discordgo.GuildFeatureVerified):              "VERIFIED",
		string(discordgo.GuildFeatureVipRegions):            "VIP_REGIONS",
		string(discordgo.GuildFeatureWelcomeScreenEnabled):  "WELCOME_SCREEN_ENABLED",
	}
)

// Role constants
var (
	RoleFlags = map[string]discordgo.RoleFlags{
		"IN_PROMPT": discordgo.RoleFlagInPrompt,
	}
)

// Permissions constants
var (
	OverwriteTypes = map[discordgo.PermissionOverwriteType]string{
		discordgo.PermissionOverwriteTypeRole:   "ROLE",
		discordgo.PermissionOverwriteTypeMember: "MEMBER",
	}

	// Full Map of all permissions
	Permissions = map[string]int64{
		"CREATE_INSTANT_INVITE":               discordgo.PermissionCreateInstantInvite,
		"KICK_MEMBERS":                        discordgo.PermissionKickMembers,
		"BAN_MEMBERS":                         discordgo.PermissionBanMembers,
		"ADMINISTRATOR":                       discordgo.PermissionAdministrator,
		"MANAGE_CHANNELS":                     discordgo.PermissionManageChannels,
		"MANAGE_GUILD":                        discordgo.PermissionManageServer,
		"ADD_REACTIONS":                       discordgo.PermissionAddReactions,
		"VIEW_AUDIT_LOG":                      discordgo.PermissionViewAuditLogs,
		"PRIORITY_SPEAKER":                    discordgo.PermissionVoicePrioritySpeaker,
		"STREAM":                              discordgo.PermissionVoiceStreamVideo,
		"VIEW_CHANNEL":                        discordgo.PermissionViewChannel,
		"SEND_MESSAGES":                       discordgo.PermissionSendMessages,
		"SEND_TTS_MESSAGES":                   discordgo.PermissionSendTTSMessages,
		"MANAGE_MESSAGES":                     discordgo.PermissionManageMessages,
		"EMBED_LINKS":                         discordgo.PermissionEmbedLinks,
		"ATTACH_FILES":                        discordgo.PermissionAttachFiles,
		"READ_MESSAGE_HISTORY":                discordgo.PermissionReadMessageHistory,
		"MENTION_EVERYONE":                    discordgo.PermissionMentionEveryone,
		"USE_EXTERNAL_EMOJIS":                 discordgo.PermissionUseExternalEmojis,
		"VIEW_GUILD_INSIGHTS":                 discordgo.PermissionViewGuildInsights,
		"CONNECT":                             discordgo.PermissionVoiceConnect,
		"SPEAK":                               discordgo.PermissionVoiceSpeak,
		"MUTE_MEMBERS":                        discordgo.PermissionVoiceMuteMembers,
		"DEAFEN_MEMBERS":                      discordgo.PermissionVoiceDeafenMembers,
		"MOVE_MEMBERS":                        discordgo.PermissionVoiceMoveMembers,
		"USE_VAD":                             discordgo.PermissionVoiceUseVAD,
		"CHANGE_NICKNAME":                     discordgo.PermissionChangeNickname,
		"MANAGE_NICKNAMES":                    discordgo.PermissionManageNicknames,
		"MANAGE_ROLES":                        discordgo.PermissionManageRoles,
		"MANAGE_WEBHOOKS":                     discordgo.PermissionManageWebhooks,
		"MANAGE_GUILD_EXPRESSIONS":            1 << 30,
		"USE_APPLICATION_COMMANDS":            1 << 31,
		"REQUEST_TO_SPEAK":                    discordgo.PermissionVoiceRequestToSpeak,
		"MANAGE_EVENTS":                       discordgo.PermissionManageEvents,
		"MANAGE_THREADS":                      discordgo.PermissionManageThreads,
		"CREATE_PUBLIC_THREADS":               discordgo.PermissionCreatePublicThreads,
		"CREATE_PRIVATE_THREADS":              discordgo.PermissionCreatePrivateThreads,
		"USE_EXTERNAL_STICKERS":               discordgo.PermissionUseExternalStickers,
		"SEND_MESSAGES_IN_THREADS":            discordgo.PermissionSendMessagesInThreads,
		"USE_EMBEDDED_ACTIVITIES":             discordgo.PermissionUseActivities,
		"MODERATE_MEMBERS":                    discordgo.PermissionModerateMembers,
		"VIEW_CREATOR_MONETIZATION_ANALYTICS": 1 << 41,
		"USE_SOUNDBOARD":                      1 << 42,
		"CREATE_GUILD_EXPRESSIONS":            1 << 43,
		"CREATE_EVENTS":                       1 << 44,
		"USE_EXTERNAL_SOUNDS":                 1 << 45,
		"SEND_VOICE_MESSAGES":                 1 << 46,
		"SEND_POLLS":                          1 << 49,
		"USE_EXTERNAL_APPS":                   1 << 50,
	}
)

// Channel constants
var (
	ChannelTypes = map[discordgo.ChannelType]string{
		discordgo.ChannelTypeGuildText:          "GUILD_TEXT",
		discordgo.ChannelTypeDM:                 "DM",
		discordgo.ChannelTypeGuildVoice:         "GUILD_VOICE",
		discordgo.ChannelTypeGroupDM:            "GROUP_DM",
		discordgo.ChannelTypeGuildCategory:      "GUILD_CATEGORY",
		discordgo.ChannelTypeGuildNews:          "GUILD_ANNOUCEMENT",
		discordgo.ChannelTypeGuildStore:         "GUILD_STORE",
		discordgo.ChannelTypeGuildNewsThread:    "ANNOUNCEMENT_THREAD",
		discordgo.ChannelTypeGuildPublicThread:  "PUBLIC_THREAD",
		discordgo.ChannelTypeGuildPrivateThread: "PRIVATE_THREAD",
		discordgo.ChannelTypeGuildStageVoice:    "GUILD_STAGE_VOICE",
		discordgo.ChannelTypeGuildDirectory:     "GUILD_DIRECTORY",
		discordgo.ChannelTypeGuildForum:         "GUILD_FORUM",
		discordgo.ChannelTypeGuildMedia:         "GUILD_MEDIA",
	}

	// Unused
	VideoQualities = map[string]int{
		"AUTO": 1,
		"FULL": 2,
	}

	ChannelFlags = map[string]discordgo.ChannelFlags{
		"PINNED":                      discordgo.ChannelFlagPinned,
		"REQUIRE_TAG":                 discordgo.ChannelFlagRequireTag,
		"HIDE_MEDIA_DOWNLOAD_OPTIONS": 1 << 15,
	}

	ForumSortOrders = map[discordgo.ForumSortOrderType]string{
		discordgo.ForumSortOrderLatestActivity: "LATEST_ACTIVITY",
		discordgo.ForumSortOrderCreationDate:   "CREATION_DATE",
	}

	ForumLayouts = map[discordgo.ForumLayout]string{
		discordgo.ForumLayoutNotSet:      "NOT_SET",
		discordgo.ForumLayoutListView:    "LIST_VIEW",
		discordgo.ForumLayoutGalleryView: "GALLERY_VIEW",
	}
)

// Webhook constants
var (
	WebhookTypes = map[discordgo.WebhookType]string{
		discordgo.WebhookTypeIncoming:        "INCOMING",
		discordgo.WebhookTypeChannelFollower: "CHANNEL_FOLLOWER",
		3:                                    "APPLICATION",
	}
)
