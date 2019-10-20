package main

import (
	"github.com/sterlingdeng/cyoa/server"
	"github.com/sterlingdeng/cyoa/story"
	"log"
)

const filename = "./gopher.json"
const PORT = 8080

func main() {
	repo, err := story.NewStory(filename)
	if err != nil {
		log.Fatal(err)
	}
	svr := server.NewServer(PORT, repo)
	svr.Start()
}
