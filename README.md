# Clapper 👏

Yes there is clap, but this is clapper (pun intended).

A more complete command line parser with optionals, tagged-defaults and auto-help.

## Differences to clap

A lot. Like pointer receivers (optionals), defaults, auto-help, etc. I know that clap originally deals with defaults by assigning them inside the struct on init. It did not like it. However, I liked clap as a command line parser but was missing some of these features and the pre-defined order so I always struggled with it.

## Maturity

**NOTE** with version 1.0.0 this package must be included like
```golang
import "github.com/Dirk007/clapper"
```

the `pkg/clapper`-Path has been removed.

Expectations could be assured from the [argument-parser](./pkg/clapper/arg_parser_test.go)- and [interpreter](./pkg/clapper/clapper_test.go) tests which are hopefully exhautive enough.

If you find a bug feel free to add a test and file a PR.

## Usage

Define your clapper-tag like this: 
```golang
type Foo struct {
    // evaluates --some-defaulted-value - or -s. If not given, default applies.
    SomeDefaultedValue int     `clapper:"short=s,long,default=42,help='Set this to have some value here instead of 42'"`
    // evaluates only --optional-value as `short` is not given. If not given, the value is nil.
    OptionalValue      *string `clapper:"long"`
    // this values has to be given or the `Parse()` will fail.
    MandatoryValue     string  `clapper:"long`
    // bool value will be `false` if not given, or `true` if -F is given.
    FlagValueOptional  bool    `clapper:"short`
    // unevaluated value as there is no `clapper`-Tag.
    IgnoredValue       bool
}

// Parse all os.Args and try to set then to `foo`
var foo Foo
trailing, err := clapper.Parse(&foo)
```

the inner order of the tag does not matter at all. The given example gives command line options `-s` and `--some-value`, defaulting to value `42` if not given. `OptionalValue` is optional. 
`FlagValueOptional` will become `true` if `-F` is given, `false` if not.


## Changes

### 1.1.0

Added `command` tag to support a struct embedded command from the command line.
Usage:

```golang
type Foo struct {
    Help    bool   `clapper:"long"`
    Command string `clapper:"command,help=show|hide"`
}
```

Given an invocation like `someProgram foo bar`, `foo` will be assigned to `Command` in the Foo-struct.

`command`-tags can have any supported type. This means if it has a slice type, all given trailing inputs will be assigned to the `command`-tag-value. If you have a single type, only the first trailing input will be assigned and the rest still remains as `trailing` from the `Parse()` invocation.

If a required command is not given, the resulting error will respect any `help` inside the tag and enrich the error message accordingly.


## Guarantees and rules

- If no arguments than the target have been given to `Parse()`, the `os.Arguemnts` will be taken.
- All value properties except `bool` are **mandatory** unless a `default` ist given.
- All pointer properties are **optional**. `default` applies.
- A bool proerty not given will remain untouched.
- At least short or long name must be provided.
- - If both are given, the long-name provided value has higher priority. `--some 1 -s 2` -> 1.
- If a flag is provided several times with different values and the property is not a slice, only the first value will be taken. `--foo 1 --foo 2` -> `foo=1`
- If a flag is provided several times and the property is a slice, the values are appended in the given order. `--foo 1 --foo 2` -> `foo=[1,2]`
- Slice properties can also be filled like `--foo a b c -d`. -> `foo=[a,b,c]`
- Short flags `-s -a -d` can be combined as `-sad` and will be interpreted as `-s -a -d`.
- If combined short-flags are provided with a value `-sad 123`, the value will be bound to the last short-flag. `d=123`
- Unknown but given flags are silently discared.
- Short flags are always `-[char]` - one dash, one character
- Long flags have to be always `--[some-string]` two dahes, longer than 1 char.
- Boolean properties can be set to `true` if the flag is given by command line. `-f`.
- - Values for bools are not accepted. Defined means `true`.
- Command line input like `--foo=bar` or `--foo bar` are interpreted as the same.
- If the last command line parameters are assigned to a slice `--foo a b c` then all these parameters will be appended to the slice. There are no trailing parameters then.
- - In opposite if there is a slice `--foo a b c --bar baz 1 2 3`, then the trailing parameters `1 2 3` will be returned from the `Parse`.
- Only one `command`-tag can be defined. If it is defined more than once, the `Parse()` will fail.
- If a `command`-tag is defined, the input is mandatory.
- A single `command`-tag target will be set with the first trailing argument.
- A slice type `command`-tag will be set to all trailing arguments.

## Tag-Options

### short
Given only as `short`, assumes that the first lowercased-letter of the property is taken as the short parameter. For exmaple `Something` becomes `-s`.

```
someprogram -s "hello world"
```

 To override the default assumption, you can write `short=X` where `X` would be the overriden short parameter. 

 ```
 someprogram -X "hello world"
 ```
 
### long
Given as `long`, assumes that the Kebab-Case-converted name of the property is the paramter name. For example `FooBar` becomes `--foo-bar`

```
someprogram --foo-bar "Hello World"
```

Also here, the name assumption can be overriden by paramter.
```long=not-so-foo```

```
someprogram --not-so-foo "Hello World"
```

### default
The default value taken if the parameter is not provided via command line. Set this to have all mandatory flags to be optional with this default. 

If `default` is defined for a pointer-value, the default will also be applied then.

### help
Clapper has a auto-help feature and this optional tag-option can be set to let your users have some extra idea of the meaning of your flag.

Exmaple:
```golang
type Foo struct {
    WeirdParameter int `clappper:"short,help='Set to the value you whish to have on your bank account.'"`
}

help := clapper.HelpDefault()
fmt.Println(help)
```

Call `clapper.Help()` to print out the avaiable options with their help text. 

FYI cllapper `HelpDefault()` will use the originally designed help format and prints out the parameters only.
You may want to have a 
```golang
type Foo struct {
    // ...
    ShowHelp bool `clapper:"long=help"`
}
```
in order to check that and display the help. But it is up to you.

## command

Up from version 1.1.0 clapper supports a `command`-tag which will be filled with the trailing arguments given. Only one field with `command` can be specified.

If the struct field type is a slice, all trailing arguments will be assigned to it. If the type is a single value property, only the first trailing argument will be assigned and the rest will still be returned as trailing arguments from the `Parse()` function as before.

If a command is specified, it becomes mandatory and the `Parse()` will fail if no command was found (i.e. no trailing arguments found).
If such a required command is not given, the resulting error will respect any optionally given `help` inside the tag and enrich the error message accordingly..

Example:
```golang
type Foo struct {
    Help    bool   `clapper:"long"`
    Command string `clapper:"command,help=show|hide"`
}
```

## Trailing?

Clapper works different from clap and does not include `trailing` as a struct property. Trailing parameters are returned from the `Parse()` command. It is up to you to do whatever you like with them.
Please also see the `command`-tag above.
