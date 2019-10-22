package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"quiet_hn/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}

		doneCh := make(chan struct{})
		storyList := StoryList{
			Items:  make([]Item, len(ids)),
			Status: make([]int, len(ids)),
			doneCh: doneCh,
		}

		idCh := broadcastIds(doneCh, ids)
		itemCh := getItems(doneCh, idCh, client)
		go addToSlice(itemCh, storyList)
		<-doneCh

		unprocessedlist := storyList.GetList()
		var stories []Item
		for _, story := range unprocessedlist {
			if isStoryLink(story) {
				stories = append(stories, story)
				if len(stories) >= 30 {
					break
				}
			}
		}

		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	}
}

type idPayload struct {
	id  int
	pos int
}

// stage 1
// work generator
func broadcastIds(done <-chan struct{}, ids []int) <-chan idPayload {
	out := make(chan idPayload)
	go func() {
		defer close(out)
		for pos, id := range ids {
			select {
			case out <- idPayload{id, pos}:
				fmt.Println(pos)
			case <-done:
				fmt.Println("closing worker in broadcasterId")
				return
			}
		}
	}()
	return out
}

type itemPayload struct {
	item Item
	pos  int
}

// stage 2
// idPayload -> itemPayload
// fan out to worker pool
func getItems(done <-chan struct{}, idCh <-chan idPayload, client hn.Client) <-chan itemPayload {
	var wg sync.WaitGroup
	out := make(chan itemPayload)
	// define the function that performs the heavy lifting
	getter := func(done <-chan struct{}, idCh <-chan idPayload, id int) {
		defer wg.Done()
		for idPayload := range idCh {
			// do work on the idPayload before entering select statement
			item, err := client.GetItem(idPayload.id)
			if err != nil {
				log.Printf("failed to process id %d, err: %v", idPayload.id, err)
			}
			newItem := parseHNItem(item)
			// enter select statement
			select {
			case out <- itemPayload{newItem, idPayload.pos}:
			case <-done:
				fmt.Printf("closing worker #%d in getItems\n", id)
				return
			}
		}
	}

	var goroutine_count = 15
	wg.Add(goroutine_count)
	for i := 0; i < goroutine_count; i++ {
		fmt.Printf("opening worker %d\n", i)
		go getter(done, idCh, i)
	}

	go func() {
		wg.Wait()
		fmt.Println("closing itemPayload channel in getItems")
		close(out)
	}()

	return out
}

// stage 3
// fan-in results from worker pool
func addToSlice(itemCh <-chan itemPayload, list StoryList) {
	for item := range itemCh {
		list.Add(item)
	}
}
