package clapper

import (
	"reflect"
	"strings"
)

type ArgsState int

const (
	ArgsStart ArgsState = iota
	ArgsShort
	ArgsLong
)

type ArgsParser struct {
	state    ArgsState
	lastKey  string
	Params   map[string][]string
	Trailing []string
}

func NewArgsParser() *ArgsParser {
	return &ArgsParser{
		state:    ArgsStart,
		lastKey:  "",
		Params:   make(map[string][]string),
		Trailing: make([]string, 0),
	}
}

func (pa *ArgsParser) With(modifier func(args *ArgsParser)) *ArgsParser {
	modifier(pa)
	return pa
}

func (pa *ArgsParser) ValuesEqualWith(other *ArgsParser) bool {
	return reflect.DeepEqual(pa.Params, other.Params) &&
		reflect.DeepEqual(pa.Trailing, other.Trailing)
}

func (pa *ArgsParser) initKey(key string) {
	_, ok := pa.Params[key]
	if !ok {
		pa.Params[key] = make([]string, 0)
	}
}

func (pa *ArgsParser) addValue(key string, value string) *ArgsParser {
	pa.initKey(key)
	pa.Params[key] = append(pa.Params[key], value)
	pa.Trailing = append(pa.Trailing, value)
	return pa
}

func (pa *ArgsParser) clearTrailBuffer() *ArgsParser {
	pa.Trailing = make([]string, 0)
	return pa
}

func (pa *ArgsParser) Add(arg string) error {
	if len(arg) == 0 {
		return ErrEmptyArgument
	}

	if arg[0] != '-' && pa.state == ArgsStart {
		pa.Trailing = append(pa.Trailing, arg)
		return nil
	}

	if arg[0] == '-' {
		pa.clearTrailBuffer()
	}

	if (pa.state == ArgsShort || pa.state == ArgsLong) && arg[0] != '-' {
		if pa.lastKey == "" {
			pa.Trailing = append(pa.Trailing, arg)
			return nil
		}

		pa.addValue(pa.lastKey, arg)
	}

	if strings.HasPrefix(arg, "--") {
		pa.state = ArgsLong
		arg = arg[2:]
		pa.lastKey = arg
		pa.initKey(pa.lastKey)
	} else if strings.HasPrefix(arg, "-") {
		pa.state = ArgsShort
		arg = arg[1:]
		for _, char := range arg {
			pa.lastKey = string(char)
			pa.initKey(pa.lastKey)
		}
	}

	return nil
}

func (pa *ArgsParser) PopTrailing(took []string) *ArgsParser {
	if took == nil {
		return pa
	}
	for _, item := range took {
		if len(pa.Trailing) == 0 {
			return pa
		}
		if pa.Trailing[0] == item {
			pa.Trailing = pa.Trailing[1:]
		}
	}
	return pa
}

func (pa *ArgsParser) Parse(args []string) (*ArgsParser, error) {
	for _, arg := range args {
		// --foo=2 becomes --foo 2
		// Maybe we have to watch for "foo bar"
		splitted := strings.SplitN(arg, "=", 2)
		for _, item := range splitted {
			err := pa.Add(item)
			if err != nil {
				return nil, err
			}
		}
	}
	return pa, nil
}
