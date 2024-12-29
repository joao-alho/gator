package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registered_commands map[string]func(*state, command) error
}

func (c *commands) Register(name string, f func(*state, command) error) error {
	c.registered_commands[name] = f
	return nil
}

func (c *commands) Run(s *state, cmd command) error {
	f, ok := c.registered_commands[cmd.Name]
	if !ok {
		return errors.New("Invalid command")
	}
	if err := f(s, cmd); err != nil {
		return err
	}
	return nil
}
