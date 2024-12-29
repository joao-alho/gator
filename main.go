package main

import (
	"database/sql"
	"log"
	"os"

	config "github.com/joao-alho/gator/internal/config"
	database "github.com/joao-alho/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", cfg.DBURL)
	}
	defer db.Close()
	dbQueries := database.New(db)

	current_state := state{cfg: &cfg, db: dbQueries}
	cmds := commands{
		registered_commands: make(map[string]func(*state, command) error),
	}
	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerGetUsers)
	cmds.Register("agg", handlerAgg)
	cmds.Register("addfeed", handlerAddFeed)

	args := os.Args
	if len(args) <= 1 {
		log.Fatal("Too few arguments!")
	}
	cmd := command{
		Name: args[1],
		Args: args[2:],
	}
	if err := cmds.Run(&current_state, cmd); err != nil {
		log.Fatal(err)
	}
}
