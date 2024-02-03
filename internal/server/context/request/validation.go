package request

import (
	"fmt"
	"strings"
	"time"
	tagsList "web/storage/tags-list"
)

type MultiError []error

func (m MultiError) Error() string {
	var errStrings []string
	for _, err := range m {
		errStrings = append(errStrings, err.Error())
	}
	return strings.Join(errStrings, "; ")
}

func (t *TaskRequest) ValidateRequest(allTagsList *tagsList.TagsList) error {
	var errors MultiError

	err := t.validateText()
	if err != nil {
		errors = append(errors, err)
	}

	err = t.ValidateTags(allTagsList)
	if err != nil {
		errors = append(errors, err)
	}

	err = t.ValidateDue()
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (t *TaskRequest) validateText() error {
	if len(t.Text) <= 0 {
		return fmt.Errorf("expect non-empty text")
	}
	return nil
}

func (t *TaskRequest) ValidateTags(allTagsList *tagsList.TagsList) error {
	if len(t.Tags) <= 0 {
		return fmt.Errorf("expect non-empty tags")
	}

	for _, tag := range t.Tags {
		if _, ok := (*allTagsList)[tag]; !ok {
			return fmt.Errorf("tag '%s' not found", tag)
		}
	}

	uniqTag := make(map[string]int)

	for _, tag := range t.Tags {
		uniqTag[tag]++
		if uniqTag[tag] > 1 {
			return fmt.Errorf("duplicate tag '%s'", tag)
		}
	}

	return nil
}

func (t *TaskRequest) ValidateDue() error {
	// TODO check due date format
	//const dateFormat = "2006-01-02T15:04:05Z"
	//
	//dueDate, err := time.Parse(dateFormat, due.String())
	//fmt.Println(dueDate)
	//if err != nil {
	//	return nil, fmt.Errorf("invalid date format")
	//}
	date, err := time.Parse(time.RFC3339, t.Due)
	if err != nil {
		return fmt.Errorf("expect due date in RFC3339 format, given: '%v'", t.Due)
	}

	if date.IsZero() {
		return fmt.Errorf("expect non-zero due date, given: %v", t.Due)
	}

	if date.Before(time.Now()) {
		return fmt.Errorf("expect due date in the future")
	}

	return nil
}
