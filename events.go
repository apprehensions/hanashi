package main

import (
	"github.com/diamondburned/arikawa/v3/gateway"
)

func (s *State) onReady(r *gateway.ReadyEvent) {
	s.addGuilds(r.UserSettings.GuildFolders)
}
