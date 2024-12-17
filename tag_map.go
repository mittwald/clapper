package clapper

type (
	// TagMap represents all tags in a single struct fields tag line.
	TagMap map[TagType]Tag
)

func NewTagMap() TagMap {
	return make(map[TagType]Tag)
}

func (t TagMap) HasTagType(tagType TagType) bool {
	_, ok := t[tagType]
	return ok
}

// HasInputTag returns true if the TagMap contains a tag that can be filled by command line arguments.
func (t TagMap) HasInputTag() bool {
	return t.HasTagType(TagShort) || t.HasTagType(TagLong)
}

// InputArgument returns the name of the command line argument.
// Long names take precedence over short names.
// If there is no input tag, it returns "<unknown>".
func (t TagMap) InputArgument() string {
	tag, ok := t[TagLong]
	if !ok {
		tag, ok = t[TagShort]
		if !ok {
			return "<unknown>"
		}
	}
	return tag.ArgumentName()
}
