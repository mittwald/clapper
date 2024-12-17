package clapper

import (
	"slices"
	"strings"
)

type ArgType int

const (
	ArgTypeShort ArgType = iota
	ArgTypeLong
	ArgTypeValue
)

func NewArgType(arg string) ArgType {
	if !strings.HasPrefix(arg, "-") {
		return ArgTypeValue
	}

	if strings.HasPrefix(arg, "--") {
		return ArgTypeLong
	}

	return ArgTypeShort
}

func (t ArgType) Value(from string) string {
	switch t {
	case ArgTypeShort:
		return from[1:]
	case ArgTypeLong:
		return from[2:]
	default:
		return from
	}
}

type ArgValue struct {
	Type     ArgType
	Value    string
	Consumed bool
}

type ArgParserExt struct {
	Args []ArgValue
}

// splitAssignmets splits an argument into its key and value if present.
func splitAssignmets(args []string) []string {
	result := make([]string, 0)
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		for _, part := range parts {
			result = append(result, part)
		}
	}
	return result
}

// sanitizeArgs removes prefixed values and expands combined short flags as well as splits "="-assignments".
func sanitizeArgs(args []string) []string {
	args = splitAssignmets(args)
	sanitized := make([]string, 0)
	first := true
	for _, arg := range args {
		argType := NewArgType(arg)

		// Skip any leading values as they can not be assigned to any argument.
		if argType == ArgTypeValue && first {
			continue
		}
		first = false

		// Expand combined short flags (iE -abc becomes -a -b -c).
		if argType == ArgTypeShort && len(arg) > 2 {
			for _, c := range arg[1:] {
				sanitized = append(sanitized, "-"+string(c))
			}
			continue
		}

		sanitized = append(sanitized, arg)
	}
	return sanitized
}

func NewArgParserExt(args []string) *ArgParserExt {
	sanitized := sanitizeArgs(args)
	ext := &ArgParserExt{
		Args: make([]ArgValue, 0, len(sanitized)),
	}
	for _, arg := range sanitized {
		argType := NewArgType(arg)
		value := argType.Value(arg)
		ext.Args = append(ext.Args, ArgValue{
			Type:     argType,
			Value:    value,
			Consumed: false,
		})
	}
	return ext
}

// findAll finds all arguments matching the given key and type.
// if the key occurs multiple times, all values are returned (iE -a 1 -a 2 -> a=[1,2]).
func (ext *ArgParserExt) findAll(key string, argType ArgType) (args []*ArgValue, ok bool) {
	args = make([]*ArgValue, 0)
	ok = false

	for indexArg, arg := range ext.Args {
		if arg.Type == argType && arg.Value == key {
			ok = true
			// Get all values for this key if there are any.
			for index := indexArg + 1; index < len(ext.Args) && ext.Args[index].Type == ArgTypeValue; index++ {
				args = append(args, &ext.Args[index])
				indexArg++
			}
		}
	}

	return args, ok
}

func (ext *ArgParserExt) Get(key string, argType ArgType) (values []string, ok bool) {
	values = make([]string, 0)
	args, ok := ext.findAll(key, argType)
	if !ok {
		return nil, false
	}

	for _, arg := range args {
		values = append(values, arg.Value)
	}

	return values, true
}

func (ext *ArgParserExt) Consume(key string, argType ArgType, n int) *ArgParserExt {
	args, ok := ext.findAll(key, argType)
	if !ok {
		return ext
	}
	for i, arg := range args {
		if i >= n {
			break
		}
		arg.Consumed = true
	}
	return ext
}

func (ext *ArgParserExt) GetTrailing() []string {
	result := make([]string, 0)

	for i := len(ext.Args) - 1; i >= 0; i-- {
		if ext.Args[i].Type != ArgTypeValue {
			break
		}
		if ext.Args[i].Consumed {
			break
		}
		result = append(result, ext.Args[i].Value)
	}

	slices.Reverse(result)
	return result
}

func (ext *ArgParserExt) ConsumeTrailing(n int) *ArgParserExt {
	trailing := ext.GetTrailing()
	if n > len(trailing) {
		n = len(trailing)
	}
	for i := range n {
		ext.Args[len(ext.Args)-len(trailing)+i].Consumed = true
	}
	return ext
}
