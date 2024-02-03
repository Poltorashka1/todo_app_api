package storage

import (
	"net/http"
	"strings"
)

// Tasks Models
type Task struct {
	Id   int      `json:"id"`
	Text string   `json:"text"`
	Tags []string `json:"tags"`
	Due  string   `json:"due"`
}

func NewTask(id int, text, tags, due string) *Task {
	tag := strings.Split(tags, "; ")

	return &Task{
		Id:   id,
		Text: text,
		Tags: tag,
		Due:  due,
	}
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

func (t *Tasks) Header() http.Header {
	return http.Header{}

}

// Tags Models
type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewTag(id int, name string) *Tag {
	return &Tag{
		Id:   id,
		Name: name,
	}
}

type Tags struct {
	Tags []Tag `json:"tags"`
}
