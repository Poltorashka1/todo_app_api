package request

type Request interface {
	Request() bool
}

// Todo Due optional?

// TaskRequest http request struct
type TaskRequest struct {
	Text string   `json:"text" validate:"required, max=100"`
	Tags []string `json:"tags" validate:"required"`
	Due  string   `json:"due" validate:"required"`
}

type TagsRequest struct {
	Tags []string `json:"tags" validate:"required"`
}

func (t *TagsRequest) Request() bool {
	return true
}

func (t *TaskRequest) Request() bool {
	return true
}

type TagRequest struct {
	Name string `json:"name" validate:"required, max=100"`
}

func (t *TagRequest) Request() bool {
	return true
}
