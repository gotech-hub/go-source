package codegen

type Option struct {
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
}
