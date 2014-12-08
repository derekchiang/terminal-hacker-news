package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
	"os"
)

func printErr(err error) {
	println("Something went wrong.  " +
		"Please report the following error to https://github.com/derekchiang/terminal-hacker-news/issues\n")
	fmt.Println(err)
	os.Exit(1)
}

const (
	baseURL = "https://hacker-news.firebaseio.com/"
	version = "v0"
	fullURL = baseURL + version
)

type topStories []int

type story struct {
	By    string
	Id    int
	Kids  []int
	Score int
	Time  int
	Title string
	Url   string
}

func httpGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		printErr(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		printErr(err)
	}

	return body
}

func getTopStories() topStories {
	body := httpGet(fullURL + "/topstories.json")
	var stories topStories
	err := json.Unmarshal(body, &stories)

	if err != nil {
		printErr(err)
	}
	return stories
}

func getStory(id int) story {
	body := httpGet(fullURL + fmt.Sprintf("/item/%v.json", id))
	var ret story
	err := json.Unmarshal(body, &ret)
	if err != nil {
		printErr(err)
	}
	return ret
}

func topStoriesCommand(c *cli.Context) {
	for idx, story_id := range getTopStories() {
		story := getStory(story_id)
		fmt.Printf("%v. %v\n", idx, story.Title)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "hn"
	app.Usage = "Read Hacker News like a real hacker"

	app.Action = topStoriesCommand

	app.Commands = []cli.Command{
		{
			Name:   "top",
			Usage:  "Show the top stories",
			Action: topStoriesCommand,
		},
	}

	app.Run(os.Args)
}
