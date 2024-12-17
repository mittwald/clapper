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

func UsageHelp(tags ParsedTags) (string, bool) {
	for _, tagItems := range tags {
		if !tagItems.HasTagType(TagCommand) {
			continue
		}

		helpTag, ok := tagItems[TagHelp]
		if !ok {
			// Without help tag, we can not tell anything about the command's usage.
			break
		}

		return fmt.Sprintf("Available commands: %s", helpTag.Value), true
	}

	return "", false
}

// HelpItemFromTags creates a HelpItem from the given tags or retruns nil if the tags represent an informational tag line only.
func HelpItemFromTags(tags TagMap) *HelpItem {
	invoke := ""
	var def *string
	var help *string
	if !tags.HasInputTag() {
		return nil
	}
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

	help := ""
	if usageHelp, ok := UsageHelp(parsedTags); ok {
		help = usageHelp + "\n"
	}

	formatting := DefaultHelpFormatting()
	helpItems := make([]*HelpItem, 0, len(parsedTags))
	for _, tags := range parsedTags {
		if item := HelpItemFromTags(tags); item != nil {
			helpItems = append(helpItems, item)
			formatting.Update(item)
		}
	}

	for _, h := range helpItems {
		help += formatter(h, formatting) + "\n"
	}

	return help, nil
}
