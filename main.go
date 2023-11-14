package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/ningen/v3"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type State struct {
	*ningen.State
	*tview.Application

	messageView *tview.TextView
	guildNode   *tview.TreeNode

	selectedChannel discord.ChannelID
}

func main() {
	var token string
	flag.StringVar(&token, "token", "", "")
	flag.Parse()

	if token == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	state := &State{
		State:       ningen.FromState(state.New(token)),
		Application: tview.NewApplication(),

		messageView: tview.NewTextView(),
		guildNode:   tview.NewTreeNode("Guilds"),
	}

	state.messageView.SetBorder(true)
	state.messageView.SetBackgroundColor(tcell.ColorDefault)
	state.messageView.ScrollToEnd()
	state.messageView.SetChangedFunc(func() {
		state.Draw()
	})

	guildTree := tview.NewTreeView()
	guildTree.SetBorder(true)
	guildTree.SetTopLevel(1)
	guildTree.SetBackgroundColor(tcell.ColorDefault)
	guildTree.SetRoot(state.guildNode)
	state.messageView.SetChangedFunc(func() {
		state.Draw()
	})
	guildTree.SetSelectedFunc(func(n *tview.TreeNode) {
		state.selectedChannel = 0

		if len(n.GetChildren()) != 0 {
			n.SetExpanded(!n.IsExpanded())
			return
		}

		switch ref := n.GetReference().(type) {
		case discord.GuildID:
			log.Println("Got GuildID Selected", ref)
		case discord.ChannelID:
			log.Println("Got ChannelID Selected", ref)
			state.drawMessages(ref)
			state.selectedChannel = ref
		case nil:
			log.Println("Got Unknown Selected")
		}
	})

	logView := tview.NewTextView()
	logView.SetBorder(true)
	logView.ScrollToEnd()
	logView.SetBackgroundColor(tcell.ColorDefault)
	logView.SetBorderColor(tcell.ColorRed)
	logView.SetChangedFunc(func() {
		state.Draw()
	})
	log.SetOutput(logView)

	lFlex := tview.NewFlex()
	lFlex.SetDirection(tview.FlexRow)
	lFlex.AddItem(guildTree, 0, 4, true)
	lFlex.AddItem(logView, 0, 1, false)

	rFlex := tview.NewFlex()
	rFlex.AddItem(state.messageView, 0, 1, false)

	mainFlex := tview.NewFlex()
	mainFlex.AddItem(lFlex, 0, 1, false)
	mainFlex.AddItem(rFlex, 0, 3, false)

	// README: there is no guild folder update event
	state.AddHandler(state.onReady)
	state.AddHandler(state.onReadUpdate)
	state.AddHandler(state.onMessageCreate)
	state.AddHandler(state.onGuildCreate)
	state.AddHandler(state.onGuildDelete)
	state.AddHandler(state.onGuildUpdate)

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
