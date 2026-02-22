package common

import "github.com/bwmarrin/discordgo"

const (
	UNIMPLEMENTED = "UNIMPLEMENTED"
)

// Guild constants
var (
	MessageNotificationsType = discordgo.MessageNotificationsAllMessages // Example of type
	MessageNotifications     = map[discordgo.MessageNotifications]string{
		discordgo.MessageNotificationsAllMessages:  "ALL_MESSAGES",
		discordgo.MessageNotificationsOnlyMentions: "ONLY_MENTIONS",
	}

	ExplicitContentFiltersType = discordgo.ExplicitContentFilterDisabled // Example of type
	ExplicitContentFilters     = map[discordgo.ExplicitContentFilterLevel]string{
		discordgo.ExplicitContentFilterDisabled:            "DISABLED",
		discordgo.ExplicitContentFilterMembersWithoutRoles: "MEMBERS_WITHOUT_ROLES",
		discordgo.ExplicitContentFilterAllMembers:          "ALL_MEMBERS",
	}

	MFALevelsType = discordgo.MfaLevelNone // Example of type
	MFALevels     = map[discordgo.MfaLevel]string{
		discordgo.MfaLevelNone:     "NONE",
		discordgo.MfaLevelElevated: "ELEVATED",
	}

	VerificationLevelsType = discordgo.VerificationLevelNone // Example of type
	VerificationLevels     = map[discordgo.VerificationLevel]string{
		discordgo.VerificationLevelNone:     "NONE",
		discordgo.VerificationLevelLow:      "LOW",
		discordgo.VerificationLevelMedium:   "MEDIUM",
		discordgo.VerificationLevelHigh:     "HIGH",
		discordgo.VerificationLevelVeryHigh: "VERY_HIGH",
	}

	GuildNSFWLevelsType = discordgo.GuildNSFWLevelDefault // Example of type
	GuildNSFWLevels     = map[discordgo.GuildNSFWLevel]string{
		discordgo.GuildNSFWLevelDefault:       "DEFAULT",
		discordgo.GuildNSFWLevelExplicit:      "EXPLICIT",
		discordgo.GuildNSFWLevelSafe:          "SAFE",
		discordgo.GuildNSFWLevelAgeRestricted: "AGE_RESTRICTED",
	}

	PremiumTiersType = discordgo.PremiumTierNone // Example of type
	PremiumTiers     = map[discordgo.PremiumTier]string{
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
	PermissionsOverwriteType = discordgo.PermissionOverwriteTypeRole // Example of type
	OverwriteTypes           = map[discordgo.PermissionOverwriteType]string{
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
		"USE_CLYDE_AI":                        1 << 47,
		"SET_VOICE_CHANNEL_STATUS":            1 << 48,
		"SEND_POLLS":                          1 << 49,
		"USE_EXTERNAL_APPS":                   1 << 50,
	}
)

// Channel constants
var (
	ChannelType  = discordgo.ChannelTypeGuildText // Example of type
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

	ForumSortOrderType = discordgo.ForumSortOrderLatestActivity // Example of type
	ForumSortOrders    = map[discordgo.ForumSortOrderType]string{
		discordgo.ForumSortOrderLatestActivity: "LATEST_ACTIVITY",
		discordgo.ForumSortOrderCreationDate:   "CREATION_DATE",
	}

	ForumLayoutType = discordgo.ForumLayoutListView // Example of type
	ForumLayouts    = map[discordgo.ForumLayout]string{
		discordgo.ForumLayoutNotSet:      "NOT_SET",
		discordgo.ForumLayoutListView:    "LIST_VIEW",
		discordgo.ForumLayoutGalleryView: "GALLERY_VIEW",
	}
)

// Webhook constants
var (
	WebhookType  = discordgo.WebhookTypeIncoming // Example of type
	WebhookTypes = map[discordgo.WebhookType]string{
		discordgo.WebhookTypeIncoming:        "INCOMING",
		discordgo.WebhookTypeChannelFollower: "CHANNEL_FOLLOWER",
		3:                                    "APPLICATION",
	}
)

// Member constants
var (
	MemberFlags = map[string]discordgo.MemberFlags{
		"DID_REJOIN":                      discordgo.MemberFlagDidRejoin,
		"COMPLETED_ONBOARDING":            discordgo.MemberFlagCompletedOnboarding,
		"BYPASSES_VERIFICATION":           discordgo.MemberFlagBypassesVerification,
		"STARTED_ONBOARDING":              discordgo.MemberFlagStartedOnboarding,
		"IS_GUEST":                        1 << 4,
		"STARTED_HOME_ACTIONS":            1 << 5,
		"COMPLETED_HOME_ACTIONS":          1 << 6,
		"AUTOMOD_QUARANTINED_USER":        1 << 7,
		"DM_SETTINGS_UPSELL_ACKNOWLEDGED": 1 << 9,
	}
)

// User constants
var (
	UserPremiumType  = discordgo.UserPremiumTypeNone // Example of type
	UserPremiumTypes = map[discordgo.UserPremiumType]string{
		discordgo.UserPremiumTypeNone:         "NONE",
		discordgo.UserPremiumTypeNitroClassic: "NITRO_CLASSIC",
		discordgo.UserPremiumTypeNitro:        "NITRO",
		discordgo.UserPremiumTypeNitroBasic:   "NITRO_BASIC",
	}

	UserFlags = map[string]discordgo.UserFlags{
		"DISCORD_EMPLOYEE":            discordgo.UserFlagDiscordEmployee,
		"DISCORD_PARTNER":             discordgo.UserFlagDiscordPartner,
		"HYPESQUAD_EVENTS":            discordgo.UserFlagHypeSquadEvents,
		"BUG_HUNTER_LEVEL_1":          discordgo.UserFlagBugHunterLevel1,
		"HOUSE_BRAVERY":               discordgo.UserFlagHouseBravery,
		"HOUSE_BRILLIANCE":            discordgo.UserFlagHouseBrilliance,
		"HOUSE_BALANCE":               discordgo.UserFlagHouseBalance,
		"EARLY_SUPPORTER":             discordgo.UserFlagEarlySupporter,
		"TEAM_USER":                   discordgo.UserFlagTeamUser,
		"SYSTEM":                      discordgo.UserFlagSystem,
		"BUG_HUNTER_LEVEL_2":          discordgo.UserFlagBugHunterLevel2,
		"VERIFIED_BOT":                discordgo.UserFlagVerifiedBot,
		"VERIFIED_BOT_DEVELOPER":      discordgo.UserFlagVerifiedBotDeveloper,
		"DISCORD_CERTIFIED_MODERATOR": discordgo.UserFlagDiscordCertifiedModerator,
		"BOT_HTTP_INTERACTIONS":       discordgo.UserFlagBotHTTPInteractions,
		"ACTIVE_BOT_DEVELOPER":        discordgo.UserFlagActiveBotDeveloper,
	}
)

// Integration constants
const IntegrationURL = "https://discord.com/api/oauth2/authorize?client_id=%s&permissions=%d&scope=%s"
