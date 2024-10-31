package clapper

import (
	"fmt"
	"reflect"
	"strings"
)

type FormatterFn = func(item *HelpItem, formatting *HelpFormatting) string

type HelpFormatting struct {
	InvokationMax int
	DefaultMax    int
}

func DefaultHelpFormatting() *HelpFormatting {
	return &HelpFormatting{}
}

func (h *HelpFormatting) Update(item *HelpItem) *HelpFormatting {
	if h.InvokationMax < len(item.Invokation) {
		h.InvokationMax = len(item.Invokation)
	}
	if item.Default != nil {
		if h.DefaultMax < len(*item.Default) {
			h.DefaultMax = len(*item.Default)
		}
	}
	return h
}

type HelpItem struct {
	Invokation string
	Default    *string
	Help       *string
}

func (h *HelpItem) Display(formatting HelpFormatting) string {
	result := h.Invokation
	result += strings.Repeat(" ", formatting.InvokationMax-len(result))
	def := ""
	if h.Default != nil {
		def = *h.Default
	}
	def += strings.Repeat(" ", formatting.DefaultMax-len(def))
	result += " " + def
	if h.Help != nil {
		result += " - " + *h.Help
	}
	return result
}

func HelpItemFromTags(tags map[TagType]Tag) *HelpItem {
	invoke := ""
	var def *string
	var help *string
	if shortTag, ok := tags[TagShort]; ok {
		name := strings.ToLower(shortTag.Name)
		if shortTag.HasValue() {
			name = shortTag.Value
		}
		invoke = "-" + name
	}
	if longTag, ok := tags[TagLong]; ok {
		if invoke != "" {
			invoke += ", "
		}
		name := longTag.Name
		if longTag.HasValue() {
			name = longTag.Value
		}
		invoke += "--" + name
	}
	defaultTag, ok := tags[TagDefault]
	if ok {
		def = ptr(fmt.Sprintf("(default: %s)", defaultTag.Value))
	}
	helpTag, ok := tags[TagHelp]
	if ok {
		help = &helpTag.Value
	}

	return &HelpItem{
		Invokation: invoke,
		Default:    def,
		Help:       help,
	}
}

func DefaultHelpFormatter(item *HelpItem, formatting *HelpFormatting) string {
	return item.Display(*formatting)
}

func HelpDefault[T any](target *T) (string, error) {
	return Help(target, DefaultHelpFormatter)
}

func Help[T any](target *T, formatter FormatterFn) (string, error) {
	t := reflect.TypeOf(*target)
	if t.Kind() != reflect.Struct {
		return "", ErrNoStruct
	}

	parsedTags, err := parseStructTags(t)
	if err != nil {
		return "", err
	}

	formatting := DefaultHelpFormatting()
	helpItems := make([]*HelpItem, 0, len(parsedTags))
	for _, tags := range parsedTags {
		item := HelpItemFromTags(tags)
		helpItems = append(helpItems, item)
		formatting.Update(item)
	}

	help := ""
	for _, h := range helpItems {
		help += formatter(h, formatting) + "\n"
	}

	return help, nil
}
