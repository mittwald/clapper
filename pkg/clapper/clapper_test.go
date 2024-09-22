package clapper

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneralParse(t *testing.T) {
	type Foo struct {
		FooBar     bool    `clapper:"short=x,long,help='Yes maybe no',default=false"`
		SomeString string  `clapper:"long,help='A string'"`
		Defaulted  int     `clapper:"short=d,long=defaulted,default=10"`
		Float      float64 `clapper:"short=f,default=3.14"`
		Pointer    *int    `clapper:"short=p,long=pointer"`
		Nothing    bool
	}

	var foo Foo
	trailing, err := Parse(&foo, "-x", "--some-string", "foo")
	require.NoError(t, err)

	assert.Empty(t, trailing)
	assert.Equal(t, true, foo.FooBar)
	assert.Equal(t, "foo", foo.SomeString)
	assert.Equal(t, false, foo.Nothing)
	assert.Equal(t, 10, foo.Defaulted)
	assert.Equal(t, 3.14, foo.Float)
	assert.Nil(t, foo.Pointer)

	help, err := HelpDefault(&foo)
	require.NoError(t, err)

	fmt.Println(help)
}

func TestShortLongMix_LongAlwaysWins(t *testing.T) {
	type Foo struct {
		ShortAndLong int `clapper:"short,long=longparam"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "--longparam", "1", "-S", "2")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, 1, foo.ShortAndLong)

	trailing, err = Parse(&foo, "-S", "2", "--longparam", "1")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, 1, foo.ShortAndLong)
}

func TestShortOnly(t *testing.T) {
	type Foo struct {
		ShortAndLong int `clapper:"short,long"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "-S", "1")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, 1, foo.ShortAndLong)
}

func TestShortOnlyWithOverride(t *testing.T) {
	type Foo struct {
		ShortAndLong int `clapper:"short=s,long"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "-s", "1")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, 1, foo.ShortAndLong)
}

func TestLongOnly(t *testing.T) {
	type Foo struct {
		ShortAndLong int `clapper:"short,long"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "--short-and-long", "1")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, 1, foo.ShortAndLong)
}

func TestLongOnlyWithOverride(t *testing.T) {
	type Foo struct {
		ShortAndLong int `clapper:"short,long=foo"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "--foo", "1")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, 1, foo.ShortAndLong)
}

func TestNoShort(t *testing.T) {
	type Foo struct {
		ShortAndLong int `clapper:"long"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "--short-and-long", "1")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, 1, foo.ShortAndLong)

	_, err = Parse(&foo, "-S", "1")
	assert.ErrorIs(t, err, ErrRequiredValueNotGiven)
}

func TestNoLong(t *testing.T) {
	type Foo struct {
		ShortAndLong int `clapper:"short"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "-S", "1")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, 1, foo.ShortAndLong)

	_, err = Parse(&foo, "--short-and-long", "1")
	assert.ErrorIs(t, err, ErrRequiredValueNotGiven)
}

func TestMissingOverrideIsIgnored(t *testing.T) {
	type Foo struct {
		ShortAndLong int `clapper:"short="`
	}

	var foo Foo
	trailing, err := Parse(&foo, "-S", "1")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, 1, foo.ShortAndLong)
}

func TestLongOverrideOfShortFails(t *testing.T) {
	type Foo struct {
		ShortAndLong int `clapper:"short=nope"`
	}

	var foo Foo
	_, err := Parse(&foo, "-n", "1")
	require.ErrorIs(t, err, NewParseError(
		ErrShortOverrideCanOnlyBeOneLetter, 0, "ShortAndLong", "short=nope",
	))
}

func TestShortOverrideOfLongFails(t *testing.T) {
	type Foo struct {
		ShortAndLong int `clapper:"long=s"`
	}

	var foo Foo
	_, err := Parse(&foo, "-s", "1")
	assert.ErrorIs(t, err, NewParseError(
		ErrLongMustBeMoreThanOne, 0, "ShortAndLong", "long=s",
	))
}

func TestMultiShorts(t *testing.T) {
	type Foo struct {
		A bool `clapper:"short=a"`
		B bool `clapper:"short=b"`
		C bool `clapper:"short=c"`
		D bool `clapper:"short=d,default=false"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "-abc")
	require.NoError(t, err)
	assert.Empty(t, trailing)

	assert.True(t, foo.A)
	assert.True(t, foo.B)
	assert.True(t, foo.C)

	assert.False(t, foo.D)
}

func TestDefaultDoesNotOverride(t *testing.T) {
	type Foo struct {
		A bool `clapper:"short=a,default=false"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "-a")
	require.NoError(t, err)

	assert.Empty(t, trailing)
	assert.True(t, foo.A)
}

func TestIntSlice(t *testing.T) {
	type Foo struct {
		Slice []int `clapper:"short=a,long"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "-a", "1", "-a", "2")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, []int{1, 2}, foo.Slice)
}

func TestIntSliceSeq(t *testing.T) {
	type Foo struct {
		Slice []int `clapper:"short=a,long"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "-a", "1", "2", "3", "-b")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, []int{1, 2, 3}, foo.Slice)
}

func TestEqualsSplitsArgutments(t *testing.T) {
	type Foo struct {
		A string `clapper:"short=a"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "-a=hello")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, "hello", foo.A)
}

func TestStringSlice(t *testing.T) {
	type Foo struct {
		Slice []string `clapper:"short=a,long"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "-a", "hello", "-a", "world")
	require.NoError(t, err)
	assert.Empty(t, trailing)
	assert.Equal(t, []string{"hello", "world"}, foo.Slice)
}

func TestTrailingValuesAfterSingle(t *testing.T) {
	type Foo struct {
		SomeThing string `clapper:"long"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "--some-thing", "hello", "foo", "bar")
	require.NoError(t, err)
	require.NotEmpty(t, trailing)

	assert.Equal(t, []string{"foo", "bar"}, trailing)
	assert.Equal(t, "hello", foo.SomeThing)
}

func TestTrailingValuesAfterSlice_WillNotBeTakenAsTrailing(t *testing.T) {
	type Foo struct {
		SomeThing []string `clapper:"long"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "--some-thing", "hello", "foo", "bar")
	require.NoError(t, err)
	require.Empty(t, trailing)

	assert.Equal(t, []string{"hello", "foo", "bar"}, foo.SomeThing)
}

func TestPointerSet(t *testing.T) {
	type Foo struct {
		A *int `clapper:"short=a"`
	}

	var foo Foo
	_, err := Parse(&foo, "-a", "123")
	require.NoError(t, err)
	require.NotNil(t, foo.A)

	assert.Equal(t, 123, *foo.A)
}

func TestPointerMissing(t *testing.T) {
	type Foo struct {
		A *int `clapper:"short=a"`
	}

	var foo Foo
	_, err := Parse(&foo, "--nope")
	require.NoError(t, err)

	assert.Nil(t, foo.A)
}

func TestPointerMissingWithPresetDoesNotTouch(t *testing.T) {
	type Foo struct {
		A *int `clapper:"short=a"`
	}

	bar := 123
	foo := Foo{
		A: &bar,
	}
	_, err := Parse(&foo, "--nope")
	require.NoError(t, err)
	require.NotNil(t, foo.A)

	assert.Equal(t, 123, *foo.A)
}

func TestPointerMissingWithDefault(t *testing.T) {
	type Foo struct {
		A *int `clapper:"short=a,default=234"`
	}

	var foo Foo
	_, err := Parse(&foo, "--nope")
	require.NoError(t, err)
	require.NotNil(t, foo.A)

	assert.Equal(t, 234, *foo.A)
}

func TestNoArgsTakesOsArgs(t *testing.T) {
	type Foo struct {
		A int `clapper:"short=a,default=123"`
	}

	var foo Foo
	os.Args = []string{"program", "-a", "456"}
	_, err := Parse(&foo)
	require.NoError(t, err)

	assert.Equal(t, 456, foo.A)
}

func TestQuotedValues(t *testing.T) {
	type Foo struct {
		SomeThing string `clapper:"long"`
	}

	var foo Foo
	os.Args = []string{"program", "--some-thing", "hello world this is a test"}
	trailing, err := Parse(&foo)
	require.NoError(t, err)

	assert.Empty(t, trailing)
	assert.Equal(t, "hello world this is a test", foo.SomeThing)
}

func TestNonSetBool(t *testing.T) {
	type Foo struct {
		Flag bool `clapper:"short"`
	}

	var foo Foo
	trailing, err := Parse(&foo, "-nope")
	require.NoError(t, err)

	assert.Empty(t, trailing)
	assert.Equal(t, false, foo.Flag)
}
