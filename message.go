package main

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/ningen/v3"
	"github.com/rivo/tview"
)

func (s *State) addMessage(m *discord.Message) {
	if m.Type != discord.DefaultMessage {
		return
	}

	fmt.Fprintf(s.messageView, "%s (%s): %s\n", m.Author.DisplayName, m.Author.Username, m.Content)
}

func (s *State) drawMessages(c discord.ChannelID) {
	s.messageView.Clear()

	ms, err := s.Messages(c, uint(120))
	if err != nil {
		log.Println(err)
		return
	}

	for i := len(ms)-1; i >= 0; i-- {
		s.addMessage(&ms[i])
	}
}

func (s *State) onMessageCreate(m *gateway.MessageCreateEvent) {
	mgid := m.Message.GuildID
	log.Println("catched MessageCreateEvent in", mgid)

	if s.selectedChannel == m.Message.ChannelID {
		s.addMessage(&m.Message)
		return
	}

	s.guildNode.Walk(func(n, _ *tview.TreeNode) bool {
		ref, ok := n.GetReference().(discord.GuildID)
		if !ok || ref != mgid {
			return true
		}

		ui := ningen.ChannelUnread

		if s.MessageMentions(&m.Message) == ningen.MessageMentions {
			ui = ningen.ChannelMentioned
		}

		g, err := s.Cabinet.Guild(mgid)
		if err != nil {
			log.Fatal(err)
		}

		s.QueueUpdateDraw(func() {
			n.SetText(formatGuildName(g.Name, ui))
		})

		return false
	})
}
