package clapper

import "strings"

// SanitizerFn is a function that sanitizes a slice of strings.
type SanitizerFn = func([]string) []string

// ArgumentSanitizer sanitizes a slice of strings `.With()` given sanitizer functions applied on `.Get()`.
type ArgumentSanitizer struct {
	sanitizers []SanitizerFn
	args       []string
}

// NewDefaultArgumentSanitizer returns a new ArgumentSanitizer without any actual sanitizers.
func NewArgumentSanitizer(args []string) *ArgumentSanitizer {
	return &ArgumentSanitizer{
		sanitizers: make([]SanitizerFn, 0),
		args:       args,
	}
}

// With adds a sanitizer function to the ArgumentSanitizer.
func (s *ArgumentSanitizer) With(fn SanitizerFn) *ArgumentSanitizer {
	s.sanitizers = append(s.sanitizers, fn)
	return s
}

// ExplodeShortsSanitizer retruns a sanitizer with all sanitizers enabled.
func NewDefaultArgumentSanitizer(args []string) *ArgumentSanitizer {
	return NewArgumentSanitizer(args).
		With(SanitizerSkipLeadingValues).
		With(SanitizeSplitAssignmets).
		With(SanitizeExplodeShorts)
}

// Get returns the sanitized arguments after applying all sanitizers.
func (s *ArgumentSanitizer) Get() []string {
	for _, fn := range s.sanitizers {
		s.args = fn(s.args)
	}
	return s.args
}

// SanitizeSplitAssignmets splits an argument into its key and value if present (iE --foo=bar -> --foo bar).
func SanitizeSplitAssignmets(args []string) []string {
	result := make([]string, 0)
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		for _, part := range parts {
			result = append(result, part)
		}
	}
	return result
}

// SanitizerSkipLeadingValues removes prefixed values that can not be assigned to any argument.
func SanitizerSkipLeadingValues(args []string) []string {
	for index, arg := range args {
		argType := NewArgType(arg)

		if argType != ArgTypeValue {
			return args[index:]
		}
	}
	return make([]string, 0)
}

// SanitizeExplodeShorts splits combined short flags into separate arguments (iE -abc -> -a -b -c).
func SanitizeExplodeShorts(args []string) []string {
	sanitized := make([]string, 0)
	for _, arg := range args {
		argType := NewArgType(arg)

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
