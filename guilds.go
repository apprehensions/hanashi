package main

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/rivo/tview"
)

func (s *State) addGuilds(folders []gateway.GuildFolder) {
	for _, folder := range folders {
		if len(folder.GuildIDs) <= 1 {
			s.addGuild(s.guildNode, folder.GuildIDs[0])
		} else {
			s.addGuildFolder(folder)
		}
	}
}

func (s *State) formatGuild(guildID discord.GuildID) string {
	if guild, err := s.Cabinet.Guild(guildID); err == nil {
		return guild.Name
	}

	return guildID.String()
}

func (s *State) addGuild(n *tview.TreeNode, guildID discord.GuildID) {
	name := s.formatGuild(guildID)
	node := tview.NewTreeNode(name)
	node.SetReference(guildID)
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
		s.addGuild(node, guildID)
	}

	node.CollapseAll()
}
