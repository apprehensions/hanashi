package main

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/ningen/v3"
	"github.com/rivo/tview"
)

func (s *State) addGuilds(folders []gateway.GuildFolder) {
	for _, folder := range folders {
		if len(folder.GuildIDs) > 1 {
			s.addGuildFolder(folder)
			continue
		}

		guild, err := s.Cabinet.Guild(folder.GuildIDs[0])
		if err == nil {
			s.addGuild(s.guildNode, guild)
		}
	}
}

func formatGuildName(name string, ui ningen.UnreadIndication) string {
	return formatIndicationName(name, ui)
}

func (s *State) addGuild(n *tview.TreeNode, g *discord.Guild) {
	ui := s.GuildIsUnread(g.ID, allowedChannelTypes)
	name := formatGuildName(g.Name, ui)

	node := tview.NewTreeNode(name)
	node.SetReference(g.ID)

	cs, err := s.Channels(g.ID, allowedChannelTypes)
	if err != nil {
		log.Println(err)
	} else {
		s.addGuildChannels(node, &cs)
	}

	n.AddChild(node)
}

func formatGuildFolder(folder gateway.GuildFolder) string {
	if folder.Name == "" {
		folder.Name = "Guild Folder"
	}

	if folder.Color == discord.NullColor {
		folder.Color = 0x7289da
	}

	return fmt.Sprintf("[%s::]%s[-::]", folder.Color, folder.Name)
}

func (s *State) addGuildFolder(folder gateway.GuildFolder) {
	name := formatGuildFolder(folder)
	node := tview.NewTreeNode(name)
	node.SetReference(folder.ID)
	s.guildNode.AddChild(node)

	for _, guildID := range folder.GuildIDs {
		guild, err := s.Cabinet.Guild(guildID)
		if err != nil {
			log.Println(err)
			continue
		}
		s.addGuild(node, guild)
	}

	node.CollapseAll()
}

func (s *State) onGuildCreate(gc *gateway.GuildCreateEvent) {
	log.Printf("catched %T", gc)

	s.addGuild(s.guildNode, &gc.Guild)
}

func (s *State) onGuildDelete(gd *gateway.GuildDeleteEvent) {
	log.Printf("catched %T", gd)

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
	log.Printf("catched %T", gu)

	s.guildNode.Walk(func(n, _ *tview.TreeNode) bool {
		ref, ok := n.GetReference().(discord.GuildID)
		if !ok || ref != gu.ID {
			return true
		}

		s.QueueUpdateDraw(func() {
			n.SetText(gu.Name)
		})

		return false
	})
}
