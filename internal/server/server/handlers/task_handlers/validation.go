package handlers

import (
	"fmt"
	"strconv"
	"time"
	"web/internal/server/context/request"
	tagsList "web/storage/tags-list"
)

func validateDue(dueDate string) error {

	date, err := time.Parse(time.RFC3339, dueDate)
	if err != nil {
		return fmt.Errorf("expect due date in RFC3339 format, given: %v", dueDate)
	}
	if date.IsZero() {
		return fmt.Errorf("expect non-zero due date, given: %v", dueDate)
	}
	if date.Before(time.Now().Add(-2 * 24 * time.Hour)) {
		return fmt.Errorf("no tasks were found. 2 days after the deadline, the tasks were deleted")
	}
	return nil
}

func validateTags(tags []string, allTags *tagsList.TagsList) error {
	for _, tag := range tags {
		if tag[0] == ' ' {
			return fmt.Errorf("tags must not be empty")
		}
		if _, err := strconv.Atoi(tag); err == nil {
			return fmt.Errorf("tags must not be a digit")
		}
		if _, ok := (*allTags)[tag]; !ok {
			return fmt.Errorf("tag '%s' not found", tag)
		}
	}

	uniqTag := make(map[string]int)

	for _, tag := range tags {
		uniqTag[tag]++
		if uniqTag[tag] > 1 {
			return fmt.Errorf("duplicate tag '%s'", tag)
		}
	}

	return nil
}

// validate tag name don't use because chi do it automatically +-
func validateTagName(tagName string) error {
	if tagName[0] == ' ' {
		return fmt.Errorf("tag must not be empty")
	}
	if _, err := strconv.Atoi(tagName); err == nil {
		return fmt.Errorf("tag must not be a digit")
	}
	return nil
}

// validate id don't use because chi do it automatically +-
func validateId(idString string) (*int, error) {
	if idString[0] == ' ' {
		return nil, fmt.Errorf("id param must not be empty")
	}
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return nil, fmt.Errorf("id param must be integer") // Todo check using method by due date
	}
	if idInt < 0 {
		return nil, fmt.Errorf("id param must be positive") // Todo check using method by due date
	}
	return &idInt, nil
}

func validateTagsAndDue(tagsList []string, dueDate string, tags *tagsList.TagsList) error {
	var multiError request.MultiError
	err := validateTags(tagsList, tags)
	if err != nil {
		multiError = append(multiError, err)
	}
	err = validateDue(dueDate)
	if err != nil {
		multiError = append(multiError, err)
	}
	if len(multiError) > 0 {
		return multiError
	}
	return nil
}
