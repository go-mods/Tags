package tags

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	// tagRegex matches a key:"value" pair of a struct tag. As for reflect.StructTag, the key holds
	// no control character, space, : or " and the value is a Go quoted string,
	// so it may hold escaped quotes: excel:"q,default:\"x\""
	tagRegex = regexp.MustCompile(`([^\x00-\x20\x7f:"]+):"((?:[^"\\]|\\.)*)"`)

	// nameOnlyRegex matches a value made of a name only: any character but the separators , ; and :
	nameOnlyRegex = regexp.MustCompile(`^([^,;:]*)$`)

	// nameOptionsRegex matches a name followed by options: id,omitempty or embedded;embeddedPrefix:author_
	nameOptionsRegex = regexp.MustCompile(`^([^,;:]*)[,;](.*)$`)

	// keyValueRegex matches an option with a key and a value: embeddedPrefix:author_
	keyValueRegex = regexp.MustCompile(`^([\p{L}\p{N}_]*):(.*)$`)
)

// constraintKey is the key of the options which hold several key:value pairs,
// like gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"
const constraintKey = "constraint"

type parser struct {
	Tags string
}

// Lookup returns the tag of the field for the given key, or nil when the field has no tag for this key.
//
//	field, _ := reflect.TypeOf(User{}).FieldByName("Name")
//	tag := tags.Lookup(field, "json") // `json:"name,omitempty"` gives Name "name" and the option omitempty
func Lookup(field reflect.StructField, key string) *Tag {
	value, ok := field.Tag.Lookup(key)
	if !ok {
		return nil
	}
	p := parser{}
	return p.parseTag(key+":"+strconv.Quote(value), key, value)
}

// Parse parses a struct tag string, holding one or several key:"value" pairs, and returns one Tag per pair.
//
// Each value is split into a Name and Options, see Tag. It returns nil when the string is empty.
// The error is always nil: it is kept for compatibility.
//
//	tgs, _ := tags.Parse(`json:"id,omitempty" gorm:"embedded;embeddedPrefix:author_"`)
//	// tgs[0]: Key "json", Name "id", Options [omitempty]
//	// tgs[1]: Key "gorm", Name "embedded", Options [embeddedPrefix:author_]
func Parse(tags string) ([]*Tag, error) {
	tags = strings.TrimSpace(tags)
	if len(tags) == 0 {
		return nil, nil
	}

	p := parser{Tags: tags}
	return p.parse(), nil
}

func (p *parser) parse() []*Tag {
	var tags []*Tag

	// Split the string into tags:
	// json:"id,omitempty" gorm:"embedded;embeddedPrefix:author_"
	// gives json:"id,omitempty" and gorm:"embedded;embeddedPrefix:author_"
	for _, match := range tagRegex.FindAllStringSubmatch(p.Tags, -1) {
		value, err := strconv.Unquote(`"` + match[2] + `"`)
		if err != nil {
			// Not a valid Go string (invalid escape): the value is kept as written
			value = match[2]
		}
		tags = append(tags, p.parseTag(match[0], match[1], value))
	}

	return tags
}

func (p *parser) parseTag(tag string, key string, value string) *Tag {
	value = strings.TrimSpace(value)
	name, options := p.parseValue(value)
	return &Tag{
		Tag:     strings.TrimSpace(tag),
		Key:     key,
		Value:   value,
		Name:    name,
		Options: options,
	}
}

func (p *parser) parseValue(value string) (name string, options []*Option) {
	// The value is a name only: id, First Name
	if matches := nameOnlyRegex.FindStringSubmatch(value); matches != nil {
		return strings.TrimSpace(matches[1]), nil
	}

	// The value is a name followed by options:
	// id,omitempty
	// id,omitempty,default
	// embedded;embeddedPrefix:author_
	// Name,default:anonymous
	if matches := nameOptionsRegex.FindStringSubmatch(value); matches != nil {
		return strings.TrimSpace(matches[1]), p.parseOptions(matches[2])
	}

	// Otherwise, the value holds options only: column:Last Name
	return "", p.parseOptions(value)
}

func (p *parser) parseOptions(option string) (options []*Option) {
	option = strings.TrimSpace(option)

	// Options are separated by ; then by ,
	// embedded;embeddedPrefix:author_ gives embedded and embeddedPrefix:author_
	// array,split:; gives array and split:; (see splitOptions)
	for _, ops := range splitOptions(option, ';') {
		if len(ops) == 0 {
			continue
		}

		// Options which hold several key:value pairs are not split:
		// constraint:OnUpdate:CASCADE,OnDelete:SET NULL; gives the key constraint
		// and the value OnUpdate:CASCADE,OnDelete:SET NULL;
		if key, value, ok := strings.Cut(ops, ":"); ok && key == constraintKey {
			options = append(options, &Option{Key: key, Value: value})
			continue
		}

		for _, op := range splitOptions(ops, ',') {
			// Option with a key and a value: embeddedPrefix:author_
			if matches := keyValueRegex.FindStringSubmatch(op); matches != nil {
				options = append(options, &Option{Key: matches[1], Value: matches[2]})
				continue
			}
			// Option with a key only: omitempty
			options = append(options, &Option{Key: op})
		}
	}

	return options
}

// splitOptions splits s around each separator followed by a letter, a digit or an underscore.
// Other separators belong to a value: split:,;| keeps ",;|", type:enum('a', 'b') is not split,
// and a trailing separator is kept.
func splitOptions(s string, sep byte) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] != sep {
			continue
		}
		if next, _ := utf8.DecodeRuneInString(s[i+1:]); isWordRune(next) {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	return append(parts, s[start:])
}

// isWordRune reports whether r is a letter, a digit or an underscore, in any script
func isWordRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
