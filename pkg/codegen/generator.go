package codegen

import (
	"errors"
	"strings"
)

var (
	ErrInvalidLength       = errors.New("invalid length, it should be greater than 6")
	ErrInvalidCount        = errors.New("invalid count, it should be greater than 0")
	ErrInvalidCharset      = errors.New("invalid charset, it should be alphanumeric")
	ErrInvalidSetCharset   = errors.New("invalid charset, it should not be empty")
	ErrInvalidPrefix       = errors.New("invalid prefix, it should be alphanumeric")
	ErrInvalidSetPrefix    = errors.New("invalid prefix, it should not be empty")
	ErrInvalidSuffix       = errors.New("invalid suffix, it should be alphanumeric")
	ErrorInvalidSetSuffix  = errors.New("invalid suffix, it should not be empty")
	ErrInvalidPattern      = errors.New("invalid pattern, it cannot be longer than 255")
	ErrorInvalidSetPattern = errors.New("invalid pattern, it should not be empty")

	ErrNotFeasible = errors.New("not feasible to generate requested number of codes")
)

const (
	// Charset types
	CharsetNumbers      = "0123456789"
	CharsetAlphabetic   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetAlphanumeric = CharsetNumbers + CharsetAlphabetic

	// Default minimums
	minCount  uint64 = 1
	minLength uint8  = 6

	patternChar = "#"
)

type CodeGen struct {
	// Length of the code to generate
	Length uint8

	// Count of the codes
	// How many codes to generate
	Count uint64

	// Charset to use
	// `CharsetNumbers`, `CharsetAlphabetic`, and `CharsetAlphanumeric`
	// are already defined, and you can use them.
	Charset string

	// Prefix of the code
	Prefix string

	// Suffix of the code
	Suffix string

	// Pattern of the code
	// # is the placeholder for the charset
	Pattern string

	// Error
	err error
}

func NewWithOptions(opt Option) (*CodeGen, error) {
	// Check length
	if opt.Length < minLength {
		return nil, ErrInvalidLength
	}

	// Check count
	if opt.Count == 0 {
		return nil, ErrInvalidCount
	}

	// Check charset
	if len(opt.Charset) > 0 && !isAlphanumeric(opt.Charset) {
		return nil, ErrInvalidCharset
	}

	if len(opt.Charset) == 0 {
		opt.Charset = CharsetAlphanumeric
	}

	// Check prefix
	if len(opt.Prefix) > 0 && !isAlphanumeric(opt.Prefix) {
		return nil, ErrInvalidPrefix
	}

	// Check suffix
	if len(opt.Suffix) > 0 && !isAlphanumeric(opt.Suffix) {
		return nil, ErrInvalidSuffix
	}

	// Check pattern
	if len(opt.Pattern) > 0 {
		numPatternChar := numberOfChar(opt.Pattern, patternChar)
		if numPatternChar > 255 {
			return nil, ErrInvalidPattern
		}

		if opt.Length == 0 || opt.Length != uint8(numPatternChar) {
			opt.Length = uint8(numPatternChar)
		}
	} else {
		opt.Pattern = repeatStr(opt.Length, patternChar)
	}

	return &CodeGen{
		Length:  opt.Length,
		Count:   opt.Count,
		Charset: opt.Charset,
		Prefix:  opt.Prefix,
		Suffix:  opt.Suffix,
		Pattern: opt.Pattern,
	}, nil
}

func New() *CodeGen {
	return &CodeGen{
		Length:  minLength,
		Count:   minCount,
		Charset: CharsetAlphanumeric,
		Pattern: repeatStr(minLength, patternChar),
	}
}

func (c *CodeGen) Generate() ([]string, error) {
	if c.err != nil {
		return nil, c.err
	}

	if !isFeasible(c.Charset, c.Pattern, patternChar, c.Count) {
		return nil, ErrNotFeasible
	}

	codes := make([]string, c.Count)

	var i uint64
	for i = 0; i < c.Count; i++ {
		codes[i] = c.one()
	}

	return codes, nil
}

// one generates one code
func (c *CodeGen) one() string {
	pts := strings.Split(c.Pattern, "")
	for i, v := range pts {
		if v == patternChar {
			pts[i] = randomChar([]byte(c.Charset))
		}
	}

	return c.Prefix + strings.Join(pts, "") + c.Suffix
}

// SetLength sets the length of the code, overrides the pattern
func (c *CodeGen) SetLength(length uint8) *CodeGen {
	if length < minLength {
		c.err = ErrInvalidLength
		return c
	}

	c.Length = length
	c.Pattern = repeatStr(c.Length, patternChar)
	return c
}

// SetCount sets the count of the codes to generate
func (c *CodeGen) SetCount(count uint64) *CodeGen {
	if count == 0 {
		c.err = ErrInvalidCount
		return c
	}

	c.Count = count
	return c
}

// SetCharset sets the charset of the code
func (c *CodeGen) SetCharset(charset string) *CodeGen {
	if len(charset) == 0 {
		c.err = ErrInvalidSetCharset
		return c
	}

	if !isAlphanumeric(charset) {
		c.err = ErrInvalidCharset
		return c
	}

	c.Charset = charset
	return c
}

// SetPrefix sets the prefix of the code
func (c *CodeGen) SetPrefix(prefix string) *CodeGen {
	if len(prefix) == 0 {
		c.err = ErrInvalidSetPrefix
		return c
	}

	if !isAlphanumeric(prefix) {
		c.err = ErrInvalidPrefix
		return c
	}

	c.Prefix = prefix
	return c
}

// SetSuffix sets the suffix of the code
func (c *CodeGen) SetSuffix(suffix string) *CodeGen {
	if len(suffix) == 0 {
		c.err = ErrorInvalidSetSuffix
		return c
	}

	if !isAlphanumeric(suffix) {
		c.err = ErrInvalidSuffix
		return c
	}

	c.Suffix = suffix
	return c
}

// SetPattern sets the pattern of the code, overrides the length
func (c *CodeGen) SetPattern(pattern string) *CodeGen {
	if len(pattern) == 0 {
		c.err = ErrorInvalidSetPattern
		return c
	}

	numPatternChar := numberOfChar(pattern, patternChar)
	if numPatternChar > 255 {
		c.err = ErrInvalidPattern
		return c
	}

	if c.Length == 0 || c.Length != uint8(numPatternChar) {
		c.Length = uint8(numPatternChar)
	}

	c.Pattern = pattern
	return c
}
