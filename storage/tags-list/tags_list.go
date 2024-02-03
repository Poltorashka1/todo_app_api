package tagsList

import (
	"fmt"
	"log/slog"
	"web/internal/storage"
)

type TagsList map[string]bool

func NewTagsMemoryList(dataBase *storage.Storage, log *slog.Logger) *TagsList {
	const op = "tags_list.NewTagsMemoryList"
	var tagsList = TagsList{}

	tags, err := (*dataBase).GetAllTags()
	if err != nil {
		log.Error(fmt.Sprintf("%v: %v", op, err.Error()))
	}
	for _, tag := range tags.Tags {
		tagsList[tag.Name] = true
	}
	return &tagsList
}
