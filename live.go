package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/AkovStohs/discord-prometheus/metrics"
	"github.com/bwmarrin/discordgo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var live = cobra.Command{
	Use:   "live",
	Short: "Continually export live stats",
	Long:  "Exports live event statistics from the channels the Discord Bot is in.",
	Args:  cobra.NoArgs,
	PreRun: func(_ *cobra.Command, _ []string) {
		initDiscord()
		http.Handle("/metrics", promhttp.Handler())
		go func() {
			log.Info("Metrics endpoint listening", zap.String("addr", metricsAddr))
			if err := http.ListenAndServe(metricsAddr, nil); err != nil {
				log.Fatal("Metrics server failed", zap.Error(err))
			}
		}()
	},

	Run: runLive,
}

func runLive(c *cobra.Command, _ []string) {
	defer func() { _ = discord.Close() }()
	log.Info("Starting Discord-Prometheus live exporter")
	defer log.Info("Stopping")
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		log := log.With(
			zap.String("guild_id", m.GuildID),
			zap.String("channel_id", m.ChannelID),
			zap.String("message_id", m.Message.ID))
		// timestamp := messageTimestamp(m.Message.ID)
		metrics.DiscordUserMessages.WithLabelValues(m.GuildID, m.ChannelID, m.Author.String()).Inc()
		log.Debug("MessageCreate")
	})
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		log := log.With(
			zap.String("guild_id", m.GuildID),
			zap.String("channel_id", m.ChannelID),
			zap.String("message_id", m.MessageID),
			zap.String("emoji", m.Emoji.Name))
		metrics.DiscordReactionAddTotal.WithLabelValues(m.GuildID, m.Emoji.Name).Inc()
		log.Debug("MessageReactionAdd")
	})
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
		log := log.With(
			zap.String("guild_id", m.GuildID),
			zap.String("channel_id", m.ChannelID),
			zap.String("message_id", m.MessageID),
			zap.String("emoji", m.Emoji.Name))
		// timestamp := messageTimestamp(m.MessageID)
		metrics.DiscordReactionRemoveTotal.WithLabelValues(m.GuildID, m.Emoji.Name).Inc()
		log.Debug("MessageReactionRemove")
	})
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		log := log.With(
			zap.String("guild_id", m.GuildID),
		)
		metrics.GuildMemberJoins.WithLabelValues(m.GuildID).Inc()
		g, err := s.Guild(m.GuildID)
		if err != nil {
			log.Error("Failed to get guild", zap.Error(err))
			return
		}
		metrics.GuildMembersOnline.
			WithLabelValues(m.GuildID).
			Set(float64(g.MemberCount))
		log.Debug("UserJoined")
	})
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
		log := log.With(
			zap.String("guild_id", m.GuildID),
		)
		metrics.GuildMemberLeaves.WithLabelValues(m.GuildID).Inc()
		g, err := s.Guild(m.GuildID)
		if err != nil {
			log.Error("Failed to get guild", zap.Error(err))
			return
		}
		metrics.GuildMembersOnline.WithLabelValues(m.GuildID).Set(float64(g.MemberCount))
		log.Debug("UserLeft")
	})
	discord.AddHandler(func(s *discordgo.Session, p *discordgo.PresenceUpdate) {
		log := log.With(
			zap.String("guild_id", p.GuildID),
			zap.String("user_id", p.User.ID),
			zap.String("status", string(p.Status)),
			zap.String("activity", p.Activities[0].Name),
		)

		metrics.PresenceChanges.WithLabelValues(p.GuildID, p.User.ID, string(p.Status)).Inc()
		g, err := s.Guild(p.GuildID)
		if err != nil {
			log.Error("Failed to get guild", zap.Error(err))
			return
		}
		metrics.GuildMembersOnline.WithLabelValues(p.GuildID).Set(float64(g.MemberCount))
		log.Debug("PresenceUpdate")
	})
	discord.State.TrackPresences = true
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions
	if err := discord.Open(); err != nil {
		log.Fatal("Failed to connect to Discord", zap.Error(err))
	}

	guild, err := discord.GuildWithCounts(guildID)
	if err != nil {
		log.Fatal("Failed to connect to fetch guild", zap.Error(err))
	}
	_ = discord.UpdateWatchStatus(0, "metrics grow")
	metrics.GuildMembersOnline.WithLabelValues(guild.ID).Set(float64(guild.ApproximateMemberCount))
	log.Info("Connected to Discord")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
