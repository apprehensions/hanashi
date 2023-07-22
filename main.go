package main

import (
	"os"
	"fmt"
	"context"
	"flag"

	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/ningen/v3"
	"github.com/rivo/tview"
)

type State struct {
	*ningen.State
	*tview.Application
}

func main() {
	var token string
	flag.StringVar(&token, "token", "", "")
	flag.Parse()

	state := &State{
		State: ningen.FromState(state.New(token)),
		Application: tview.NewApplication(),
	}

	mainFlex := tview.NewFlex()
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
