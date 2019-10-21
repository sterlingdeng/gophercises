package story

import (
	"encoding/json"
	"errors"
	"os"
)

var NoIntroError = errors.New("intro not found")


type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type Story struct {
	Arcs map[string]Arc
}

func NewStory(filepath string) (Story, error) {
	s := Story{}
	err := s.ParseJSON(filepath)
	if err != nil {
		return Story{}, err
	}
	return s, nil
}

func (s *Story) ParseJSON(filepath string) (error) {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	storyArcs := make(map[string]Arc)
	// use json.NewDecoder because file implements io reader.
	d := json.NewDecoder(file)
	err = d.Decode(&storyArcs)
	if err != nil {
		return err
	}
	if _, exists := storyArcs["intro"]; !exists {
		return NoIntroError
	}
	s.Arcs = storyArcs
	return nil
}

