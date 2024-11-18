terraform {
  required_providers {
    discord = {
      source = "registry.terraform.io/justarecord/discord"
    }
  }
}

provider "discord" {}

data "discord_guild" "merchtracker_dev" {
  name = "MerchTracker DEV"
}

data "discord_channel" "logs" {
  guild_id = data.discord_guild.merchtracker_dev.id
  name     = "logs"
}

resource "discord_channel" "testing" {
  guild_id = data.discord_guild.merchtracker_dev.id
  name     = "testing"
}

resource "discord_role" "tf_testing" {
  guild_id    = data.discord_guild.merchtracker_dev.id
  name        = "tf_testing"
  mentionable = false
  hoist       = true
  color       = "#FFFFFF"
  permissions = [
    "READ_MESSAGE_HISTORY",
  ]
}

# import {
#   to = discord_permissions.testing
#   id = "${data.discord_guild.merchtracker_dev.id}/${data.discord_channel.logs.id}/role/${resource.discord_role.testing.id}"
# }

resource "discord_role" "testing" {
  name        = "testing"
  guild_id    = data.discord_guild.merchtracker_dev.id
  color       = "#1ABC9C"
  hoist       = true
  mentionable = false
  permissions = [
    "VIEW_CHANNEL",
    "EMBED_LINKS",
    "ATTACH_FILES",
  ]
}

resource "discord_permissions" "tf_testing" {
  guild_id   = data.discord_guild.merchtracker_dev.id
  channel_id = resource.discord_channel.testing.id
  id         = resource.discord_role.tf_testing.id
  type       = "role"
  allow = [
    "VIEW_CHANNEL",
  ]
  deny = [
  ]
}

resource "discord_permissions" "testing" {
  guild_id   = data.discord_guild.merchtracker_dev.id
  channel_id = resource.discord_channel.testing.id
  id         = resource.discord_role.testing.id
  type       = "role"
  allow = [
    "VIEW_CHANNEL",
  ]
  deny = [
  ]
}

data "discord_guild" "merchtracker" {
  name = "MerchTracker"
}

data "discord_member" "member" {
  guild_id = data.discord_guild.merchtracker_dev.id
  username     = "swiftsudo"
}

resource "discord_role_members" "testing" {
  guild_id = data.discord_guild.merchtracker_dev.id
  role_id  = resource.discord_role.testing.id
  members = [
    "swiftsudo",
  ]
}

output "role_members" {
  value = resource.discord_role_members.testing
}

output "member" {
  value = data.discord_member.member
}

data "discord_channel" "category" {
  guild_id = data.discord_guild.merchtracker_dev.id
  name     = "artists: us"
}

output "category_children" {
  # Number of children in the category
  value = length(data.discord_channel.category.children)
}
