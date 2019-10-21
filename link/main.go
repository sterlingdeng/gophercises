package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
	"strings"
)

type DocumentLinks struct {
	Links []Link
}

func (d *DocumentLinks) add(l Link) {
	d.Links = append(d.Links, l)
}

type Link struct {
	Href string
	Text string
}

func (l *Link) addText(text string) {
	l.Text = fmt.Sprintf("%s %s", l.Text, text)
}

func main() {
	exercises := []string{"ex1", "ex2", "ex3", "ex4"}

	for _, exercise := range exercises {
		fn := fmt.Sprintf("./pages/%s.html", exercise)
		f, err := os.Open(fn)
		if err != nil {
			log.Fatal(err)
		}

		node, err := html.Parse(f)
		if err != nil {
			log.Fatal(err)
		}

		list := DocumentLinks{}

		crawler(node, &list, nil)
		fmt.Printf("\nExercise: %s\n", exercise)
		for _, link := range list.Links {
			fmt.Printf("links for test file %+v\n", link)
		}
	}
}

func crawler(n *html.Node, list *DocumentLinks, link *Link) {
	if n.Type == html.ElementNode && n.Data == "a" && link == nil {
		href := getHrefValue(n.Attr)
		link = &Link{
			Href: href,
		}
	}
	if n.Type == html.TextNode && n.Data != "\n" && link != nil {
		text := strings.Trim(n.Data, "\n ")
		link.addText(text)
	}

	if n.FirstChild == nil && n.NextSibling == nil {
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		crawler(c, list, link)
		if link != nil && n.Type == html.ElementNode && n.Data == "a" {
			list.add(*link)
			link = nil
		}
	}
}

func getHrefValue(attrs []html.Attribute) string {
	for _, attr := range attrs {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}
