package main

import (
	"os"
	"fmt"
	"context"
	"flag"

	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/ningen/v3"
)

type State struct {
	*ningen.State
}

func main() {
	var token string
	flag.StringVar(&token, "token", "", "")
	flag.Parse()

	state := &State{
		State: ningen.FromState(state.New(token)),
	}
	
	if err := state.Open(context.TODO()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	state.Close()
}
