package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetTrendingDevelopers(language, since string) []Developer {
	var devs []Developer
	url := fmt.Sprintf("https://github-trending-api.now.sh/developers?language=%s&since=%s", language, since)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	err = json.Unmarshal([]byte(body), &devs)
	if err != nil {
		log.Fatal(err)
	}
	return devs
}

func getRepos(language, since string) ([]Repo, error) {
	var repos []Repo
	url := fmt.Sprintf("https://github-trending-api.now.sh/repositories?language=%s&since=%s", language, since)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return []Repo{}, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return []Repo{}, err
	}
	res.Body.Close()
	err = json.Unmarshal([]byte(body), &repos)
	if err != nil {
		log.Fatal(err)
		return []Repo{}, err
	}
	return repos, nil
}

func GetTrendingRepos(language, since string, respChan chan string) error {
	var responseRepos string
	repos, err := getRepos(language, since)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		responseRepos += fmt.Sprintf("%s\n\n", repo.Name)
		responseRepos += fmt.Sprintf("* %s\n", repo.Description)
		responseRepos += fmt.Sprintf("* %s\n", repo.Author)
		responseRepos += fmt.Sprintf("* %s\n", repo.Url)
		responseRepos += fmt.Sprintf("* %d\n\n", repo.Stars)
	}
	respChan <- responseRepos
	return nil
}
