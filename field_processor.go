package clapper

import (
	"errors"
	"reflect"
)

type StructFieldProcessor struct {
	targetType   reflect.Type
	targetValue  reflect.Value
	tags         ParsedTags
	args         *ArgParserExt
	currentIndex int
	commandHelp  string
	commandIndex *int
}

func NewStructFieldProcessor(target reflect.Type, value reflect.Value, tags ParsedTags, args *ArgParserExt) *StructFieldProcessor {
	return &StructFieldProcessor{
		targetType:   target,
		targetValue:  value,
		tags:         tags,
		args:         args,
		currentIndex: 0,
		commandHelp:  "",
		commandIndex: nil,
	}
}

func hasTagType(tags map[TagType]Tag, tagType TagType) bool {
	_, ok := tags[tagType]
	return ok
}

func (f *StructFieldProcessor) EOF() bool {
	return f.currentIndex >= f.targetType.NumField()
}

func (f *StructFieldProcessor) Next() error {
	if f.EOF() {
		return errors.New("all fields processed")
	}

	index := f.currentIndex
	f.currentIndex++

	tags, hasTags := f.tags[index]
	if !hasTags {
		return nil
	}

	if hasTagType(tags, TagCommand) {
		f.commandIndex = &index
		if hasTagType(tags, TagHelp) {
			f.commandHelp = tags[TagHelp].Value
		}
		return nil
	}

	field := f.targetType.Field(index)
	fieldValue := f.targetValue.Field(index)

	return trySetFieldConsumingArgs(field, fieldValue, tags, f.args)
}

func (f *StructFieldProcessor) HasCommand() bool {
	return f.commandIndex != nil
}

func (f *StructFieldProcessor) ProcessCommand() error {
	trailing := f.GetTrailing()
	if len(trailing) == 0 {
		return NewCommandRequiredError(f.commandHelp)
	}

	field := f.targetType.Field(*f.commandIndex)
	fieldValue := f.targetValue.Field(*f.commandIndex)

	took, err := StringReflect(field, fieldValue, trailing)
	if err != nil {
		return err
	}

	f.args.ConsumeTrailing(took)

	return nil
}

func (f *StructFieldProcessor) Finalize() error {
	if f.HasCommand() {
		if err := f.ProcessCommand(); err != nil {
			return err
		}
	}

	return nil
}

func (f *StructFieldProcessor) GetTrailing() []string {
	return f.args.GetTrailing()
}
