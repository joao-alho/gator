package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joao-alho/gator/commands"
	config "github.com/joao-alho/gator/internal"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	state := commands.State{Cfg: &cfg}

	cmds := commands.Commands{
		Cmds: make(map[string]func(*commands.State, commands.Command) error),
	}
	cmds.Register("login", commands.HandlerLogin)

	args := os.Args
	if len(args) <= 2 {
		log.Fatal("Too few arguments!")
	}

	login_cmd := commands.Command{
		Name: args[1],
		Args: args[2:],
	}
	cmds.Run(&state, login_cmd)

	cfg, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Read config: %v\n", cfg)
}
