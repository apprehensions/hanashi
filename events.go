package main

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/ningen/v3"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (s *State) onReady(r *gateway.ReadyEvent) {
	s.addGuilds(r.UserSettings.GuildFolders)
}

func (s *State) onMessageCreate(m *gateway.MessageCreateEvent) {
	mention := s.MessageMentions(&m.Message)

	if mention == 0 || mention == ningen.MessageMentions {
		return
	}

	s.guildNode.Walk(func(n, _ *tview.TreeNode) bool {
		ref, ok := n.GetReference().(discord.GuildID)
		if !ok || ref != m.Message.GuildID {
			return true
		}

		s.QueueUpdateDraw(func() {
			n.SetColor(tcell.ColorRed)
		})

		return false
	})
}

func (s *State) onGuildCreate(gc *gateway.GuildCreateEvent) {
	s.addGuild(s.guildNode, &gc.Guild)
}

func (s *State) onGuildDelete(gd *gateway.GuildDeleteEvent) {
	s.guildNode.Walk(func(n, p *tview.TreeNode) bool {
		ref, ok := n.GetReference().(discord.GuildID)
		if !ok || ref != gd.ID {
			return true
		}

		s.QueueUpdateDraw(func() {
			p.RemoveChild(n)
		})

		return false
	})
}

func (s *State) onGuildUpdate(gu *gateway.GuildUpdateEvent) {
	s.guildNode.Walk(func(n, _ *tview.TreeNode) bool {
		ref, ok := n.GetReference().(discord.GuildID)
		if !ok || ref != gu.ID {
			return true
		}

		if n.GetText() == gu.Name {
			return false
		}

		s.QueueUpdateDraw(func() {
			n.SetText(gu.Name)
		})

		return false
	})
}
