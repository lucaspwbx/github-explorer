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

func GetTrendingRepos(language, since string) []Repo {
	var repos []Repo
	url := fmt.Sprintf("https://github-trending-api.now.sh/repositories?language=%s&since=%s", language, since)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	err = json.Unmarshal([]byte(body), &repos)
	if err != nil {
		log.Fatal(err)
	}
	return repos
}
