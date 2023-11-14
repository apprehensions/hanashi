package main

import (
	"log"

	"github.com/diamondburned/arikawa/v3/gateway"
)

func (s *State) onReady(r *gateway.ReadyEvent) {
	log.Printf("catched %T", r)

	s.addGuilds(r.UserSettings.GuildFolders)
}
