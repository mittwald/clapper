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

Expectations could be assured from the [argument-parser](./pkg/clapper/arg_parser_test.go)- and [interpreter](./pkg/clapper/clapper_test.go) tests which are nopefully exhautive enough.

If you find a bug feel free to add a test and file a PR.

## Usage

Define your clapper-tag like this: 
```golang
type Foo struct {

    // evaluates --some-defaulted-value - or -s. If not given, default applies.
    SomeDefaultedValue int     `clapper:"short,long,default=42,help='Set this to have some value here instead of 42'"`
    // evaluates only --optional-value as `short` is not given. If not given, the value is nil.
    OptionalValue      *string `clapper:"long"`
    // this values has to be given or the `Parse()` will fail.
    MandatoryValue     string  `clapper:"long`
    // bool value will be `false` if not given, or `true` if -f is given.
    FlagValueOptional  bool    `clapper:"short`
    // unevaluated value as there is no `clapper`-Tag.
    IgnoredValue       bool
}

var foo Foo
trailing, err := clapper.Parse(&foo)
```

the inner order of the tag does not matter at all. The given example gives command line options `-s` and `--some-value`, defaulting to value `42` if not given. `OptionalValue` is optional. 


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

## Trailing?

Clapper works different from clap and does not include `trailing` as a struct property. Trailing parameters are returned from the `Parse()` command. It is up to you to do whatever you like with them.
