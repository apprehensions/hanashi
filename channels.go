package main

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/ningen/v3"
	"github.com/diamondburned/ningen/v3/states/read"
	"github.com/rivo/tview"
)

var allowedChannelTypes = []discord.ChannelType{
	discord.GuildText,
}

func formatIndicationName(name string, ui ningen.UnreadIndication) string {
	switch ui {
	case ningen.ChannelMentioned:
		return fmt.Sprintf("[red::rb]%s[-::-]\n", name)
	case ningen.ChannelUnread:
		return fmt.Sprintf("[::b]%s[::-]\n", name)
	}

	return name
}

func formatChannelName(c *discord.Channel, ui ningen.UnreadIndication) string {
	name := c.Name

	switch c.Type {
	case discord.GuildText:
		name = "#" + c.Name
	}

	return formatIndicationName(name, ui)
}

func (s *State) addChannel(n *tview.TreeNode, c *discord.Channel) {
	ui := s.ChannelIsUnread(c.ID)
	name := formatChannelName(c, ui)
	node := tview.NewTreeNode(name)
	node.SetReference(c.ID)
	n.AddChild(node)
}

func (s *State) addGuildChannels(n *tview.TreeNode, cs *[]discord.Channel) {
	for _, c := range *cs {
		s.addChannel(n, &c)
	}

	n.CollapseAll()
}

func (s *State) onReadUpdate(u *read.UpdateEvent) {
	log.Println("catched read.UpdateEvent in", u.GuildID)

	ui := ningen.ChannelRead
	if u.Unread {
		ui = ningen.ChannelUnread
	}

	s.guildNode.Walk(func(n, _ *tview.TreeNode) bool {
		ref, ok := n.GetReference().(discord.GuildID)
		if !ok || ref != u.GuildID {
			return true
		}

		g, err := s.Cabinet.Guild(u.GuildID)
		if err != nil {
			log.Fatal(err)
		}

		s.QueueUpdateDraw(func() {
			n.SetText(formatGuildName(g.Name, ui))
		})

		n.Walk(func(n, _ *tview.TreeNode) bool {
			ref, ok := n.GetReference().(discord.ChannelID)
			if !ok || ref != u.ChannelID {
				return true
			}

			c, err := s.Cabinet.Channel(ref)
			if err != nil {
				log.Fatal(err)
			}

			s.QueueUpdateDraw(func() {
				n.SetText(formatChannelName(c, ui))
			})

			return false
		})

		return false
	})
}
