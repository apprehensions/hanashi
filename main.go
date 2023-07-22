package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/ningen/v3"
	"github.com/rivo/tview"
)

type State struct {
	*ningen.State
	*tview.Application

	guildNode *tview.TreeNode
}

func main() {
	var token string
	flag.StringVar(&token, "token", "", "")
	flag.Parse()

	state := &State{
		State:       ningen.FromState(state.New(token)),
		Application: tview.NewApplication(),

		guildNode: tview.NewTreeNode("Guilds"),
	}

	guildTree := tview.NewTreeView()
	guildTree.SetBorder(true)
	guildTree.SetTopLevel(1)
	guildTree.SetRoot(state.guildNode)

	mainFlex := tview.NewFlex()
	mainFlex.AddItem(guildTree, 0, 1, false)

	state.AddHandler(state.onReady)

	state.SetRoot(mainFlex, true)
	state.EnableMouse(true)

	if err := state.Open(context.TODO()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := state.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	state.Close()
}
