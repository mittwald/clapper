package clapper

import (
	"reflect"
	"strings"
)

const (
	TagName = "clapper"
)

// ParsedTags are a map of tags for each field in a struct.
// Not each struct-field may have set a clapper tag.
type ParsedTags = map[int]map[TagType]Tag

// parseTags parses the tags for a given field (aka "one line") and returns them as a map.
func parseTags(tagItems []string, fieldName string, index int) (map[TagType]Tag, error) {
	tags := make(map[TagType]Tag, 0)
	for _, tagItem := range tagItems {
		tag, err := NewTag(tagItem, fieldName, index)
		if err != nil {
			return nil, err
		}
		tags[tag.Type] = *tag
	}
	return tags, nil
}

// parseStructTags parses a given sturct and returns all of its parsed tags.
func parseStructTags(t reflect.Type) (ParsedTags, error) {
	parsedTags := make(map[int]map[TagType]Tag, 0)
	commandTagSpecified := false
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagLine := field.Tag.Get(TagName)
		if tagLine == "" {
			continue
		}
		tagItems := strings.Split(tagLine, ",")
		tags, err := parseTags(tagItems, field.Name, i)
		if err != nil {
			return nil, NewParseError(err, i, field.Name, tagLine)
		}
		if hasTagType(tags, TagCommand) {
			if commandTagSpecified {
				return nil, ErrDuplicateCommandTag
			}
			commandTagSpecified = true
		}
		parsedTags[i] = tags
	}
	return parsedTags, nil
}
