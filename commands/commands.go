package commands

import (
	"errors"
	"fmt"

	config "github.com/joao-alho/gator/internal"
)

type State struct {
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Insufficient number of arguments.")
	}
	username := cmd.Args[0]
	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("User has been set to -> %s\n", username)
	return nil
}

type Commands struct {
	Cmds map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) error {
	c.Cmds[name] = f
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	f, ok := c.Cmds[cmd.Name]
	if !ok {
		return errors.New("Invalid command")
	}
	if err := f(s, cmd); err != nil {
		return err
	}
	return nil
}
