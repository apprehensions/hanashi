package main

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/ningen/v3"
	"github.com/rivo/tview"
)

var allowedChannelTypes = []discord.ChannelType{
	discord.GuildText,
}

func (s *State) addGuilds(folders []gateway.GuildFolder) {
	for _, folder := range folders {
		if len(folder.GuildIDs) >= 1 {
			s.addGuildFolder(folder)
			continue
		}

		guild, err := s.Cabinet.Guild(folder.GuildIDs[0])
		if err == nil {
			s.addGuild(s.guildNode, guild)
		}
	}
}

func (s *State) addGuild(n *tview.TreeNode, g *discord.Guild) {
	name := g.Name
	u := s.GuildIsUnread(g.ID, allowedChannelTypes)

	if u == ningen.ChannelMentioned {
		name = fmt.Sprintf("[red::rb]%s[-:-:-]\n", name)
	} else if u == ningen.ChannelUnread {
		name = fmt.Sprintf("[::b]%s[::-]\n", name)
	}

	node := tview.NewTreeNode(name)
	node.SetReference(g.ID)
	n.AddChild(node)
}

func (s *State) addGuildFolder(folder gateway.GuildFolder) {
	name := folder.Name
	if folder.Name == "" {
		name = "Guild Folder"
	}

	node := tview.NewTreeNode(name)
	node.SetReference(folder.ID)
	s.guildNode.AddChild(node)

	for _, guildID := range folder.GuildIDs {
		guild, err := s.Cabinet.Guild(guildID)
		if err != nil {
			continue
		}
		s.addGuild(node, guild)
	}

	node.CollapseAll()
}
