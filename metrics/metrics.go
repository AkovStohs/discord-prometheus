// metrics/metrics.go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DiscordUserMessages = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "discord_user_channel_messages_total",
			Help: "Total number of messages per user per guild per channel.",
		},
		[]string{"guild", "channel", "user"},
	)
	DiscordReactionAddTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "discord_reaction_add_total",
			Help: "Total number of reaction_add events per guild and emoji.",
		},
		[]string{"guild", "emoji"},
	)
	DiscordReactionRemoveTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "discord_reaction_remove_total",
			Help: "Total number of reaction_remove events per guild and emoji.",
		},
		[]string{"guild", "emoji"},
	)

	// New join/leave metrics
	GuildMemberJoins = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "discord_guild_member_joins_total",
			Help: "Total number of members who have joined each guild.",
		},
		[]string{"guild"},
	)
	GuildMemberLeaves = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "discord_guild_member_leaves_total",
			Help: "Total number of members who have left each guild.",
		},
		[]string{"guild"},
	)

	// Presence change counter
	PresenceChanges = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "discord_presence_changes_total",
			Help: "Total number of presence update events, partitioned by guild and new status.",
		},
		[]string{"guild", "user", "status"},
	)

	// Optional gauges for current state
	GuildMembersOnline = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "discord_guild_members_online",
			Help: "Current number of members online per guild.",
		},
		[]string{"guild"},
	)
)

func init() {
	prometheus.MustRegister(
		DiscordUserMessages,
		DiscordReactionAddTotal,
		DiscordReactionRemoveTotal,
		GuildMemberJoins,
		GuildMemberLeaves,
		PresenceChanges,
		GuildMembersOnline,
	)
}
