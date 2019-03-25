package main

type Repo struct {
	Name        string
	Description string
	Author      string
	Url         string
	Language    string
	Stars       int
}

type Developer struct {
	Username string
	Name     string
	Url      string
	Repos    []Repo
}
