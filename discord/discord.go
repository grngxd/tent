package discord

import (
	"os"
)

type DiscordInstalls struct {
	Discord       bool
	DiscordCanary bool
	DiscordPTB    bool
}

func Discords() (DiscordInstalls, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return DiscordInstalls{}, err
	}

	discord := false
	discordCanary := false
	discordPTB := false

	if _, err := os.Stat(dir + "\\discord"); err == nil {
		discord = true
	}

	if _, err := os.Stat(dir + "\\discordcanary"); err == nil {
		discordCanary = true
	}

	if _, err := os.Stat(dir + "\\discordptb"); err == nil {
		discordPTB = true
	}

	installs := DiscordInstalls{
		Discord:       discord,
		DiscordCanary: discordCanary,
		DiscordPTB:    discordPTB,
	}

	return installs, nil
}
