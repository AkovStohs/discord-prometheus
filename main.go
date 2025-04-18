package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}

var root = cobra.Command{
	Use:   "discord-prometheus",
	Short: "Discord metrics exporter",
	PersistentPreRun: func(c *cobra.Command, _ []string) {
		pflags := c.Flags()
		debug, err := pflags.GetBool(flagDebug)
		if err != nil {
			panic(err.Error())
		}
		if debug {
			logLevel.SetLevel(zapcore.DebugLevel)
		} else {
			logLevel.SetLevel(zapcore.InfoLevel)
		}
	},
}

func init() {
	pflags := root.PersistentFlags()
	pflags.Bool(flagDebug, false, "Enable debug log")
	pflags.StringVar(&metricsAddr, "metrics-addr", ":9090", "Metrics server address")
	root.AddCommand(&live)
}

const (
	flagDebug = "debug"
)

var log *zap.Logger
var logLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
var metricsAddr string
var guildID string

func init() {
	config := zap.Config{
		Level:       logLevel,
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:    "L",
			MessageKey:  "M",
			LineEnding:  zapcore.DefaultLineEnding,
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02T15:04:05Z0700"))
			},
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     true,
		DisableStacktrace: true,
	}
	var err error
	log, err = config.Build()
	if err != nil {
		panic(err.Error())
	}
}

var discord *discordgo.Session

func readDiscordToken() (discordToken string) {
	discordTokenPath := os.Getenv("DISCORD_TOKEN_FILE")
	if discordTokenPath != "" {
		discordTokenBuf, err := os.ReadFile(discordTokenPath)
		if err != nil {
			panic(err.Error())
		}
		discordToken = strings.TrimSpace(string(discordTokenBuf))
	} else {
		discordToken = strings.TrimSpace(os.Getenv("DISCORD_TOKEN"))
	}
	if discordToken == "" {
		log.Fatal("No Discord token found")
	}
	if !strings.ContainsRune(discordToken, ' ') {
		discordToken = "Bot " + discordToken
	}
	return
}

func initDiscord() {
	apiToken := readDiscordToken()
	guildID = strings.TrimSpace(os.Getenv("DISCORD_GUILD"))
	if guildID == "" {
		log.Fatal("No Discord Guild found")
	}
	var err error
	discord, err = discordgo.New(apiToken)
	if err != nil {
		log.Fatal("Error creating Discord token", zap.Error(err))
	}
}
