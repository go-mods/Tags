package tags

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type parser struct {
	Tags string
}

func Lookup(field reflect.StructField, key string) *Tag {
	if value, ok := field.Tag.Lookup(key); ok {
		tags, err := Parse(fmt.Sprintf("%s:\"%s\"", key, value))
		if tags == nil || err != nil {
			return nil
		}
		return tags[0]
	}
	return nil
}

func Parse(tags string) ([]*Tag, error) {
	tags = strings.TrimSpace(tags)
	if len(tags) == 0 {
		return nil, nil
	}

	p := parser{Tags: tags}
	return p.parse()
}

func (p *parser) parse() ([]*Tag, error) {
	// Tags list
	var tags []*Tag

	// Split tags onto multiples tag
	// ie: json:"id,omitempty" gorm:"embedded;embeddedPrefix:author_"
	// to json:"id,omitempty" and gorm:"embedded;embeddedPrefix:author_"
	regex := *regexp.MustCompile(`(.*?):"(.*?)"`)
	matches := regex.FindAllStringSubmatch(p.Tags, -1)

	// Loop throw all tags
	for i := range matches {
		tag, err := p.parseTag(matches[i][0], matches[i][1], matches[i][2])
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (p *parser) parseTag(tag string, key string, value string) (*Tag, error) {
	// cleanup
	tag = strings.TrimSpace(tag)
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)

	// Create the Tag to return
	t := &Tag{
		Tag:   tag,
		Key:   key,
		Value: value,
	}

	// Get Name and []Options from Value
	n, opts, err := p.parseValue(value)
	if err != nil {
		return nil, err
	}
	t.Name = n
	t.Options = opts

	return t, nil
}

func (p *parser) parseValue(value string) (name string, options []*Option, err error) {
	// cleanup
	value = strings.TrimSpace(value)

	// value only have a name
	// ie: id
	regex := *regexp.MustCompile(`^(\w*)$`)
	matches := regex.FindAllStringSubmatch(value, -1)
	if len(matches) == 1 {
		name = matches[0][1]
		return
	}

	// Value has options delimited with a comma
	// ie: id,omitempty
	// ie: id,omitempty,default
	// ie: embedded;embeddedPrefix:author_
	regex = *regexp.MustCompile(`^(\w*)[,;](.*)$`)
	matches = regex.FindAllStringSubmatch(value, -1)
	if len(matches) == 1 {
		name = matches[0][1]
		options, err = p.parseOptions(matches[0][2])
		return
	}

	// If no value is found then this is an option
	options, err = p.parseOptions(value)

	if len(value) > 0 && len(name) == 0 && len(options) == 0 {
		name = value
	}

	return
}

func (p *parser) parseOptions(option string) (options []*Option, err error) {
	// cleanup
	option = strings.TrimSpace(option)

	// Split options onto multiples options
	// ie: embedded;embeddedPrefix:author_
	// key: embedded, value: nil
	// key: embeddedPrefix, value: author_
	// ie: array,split:;
	// Name: array
	// key: split, value: ;
	regexSplitter := regexp.MustCompile(`;(\b)`)
	for _, ops := range regexSplitter.Split(option, -1) {

		// ops can be empty
		if len(ops) == 0 {
			continue
		}

		// Special case for options like `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		// which handle multiple values for the same key
		// ie: OnUpdate:CASCADE and OnDelete:SET NULL; from gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"
		regex := *regexp.MustCompile(`^(.*):((\w*):(.*)[,](\w*):(.*))$`)
		matches := regex.FindAllStringSubmatch(ops, -1)
		if len(matches) == 1 {
			o := Option{
				Key:   matches[0][1],
				Value: matches[0][2],
			}
			options = append(options, &o)
			continue
		}

		// Back to normal case
		regexSplitter := regexp.MustCompile(`,(\b)`)
		for _, op := range regexSplitter.Split(ops, -1) {
			// option with keys and values
			// ie: embeddedPrefix:author_ from gorm:"embedded;embeddedPrefix:author_"
			// ie: OnUpdate:CASCADE and OnDelete:SET NULL; from gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"
			regex = *regexp.MustCompile(`^(\w*):(.*)$`)
			matches = regex.FindAllStringSubmatch(op, -1)
			if len(matches) == 1 {
				o := Option{
					Key:   matches[0][1],
					Value: matches[0][2],
				}
				options = append(options, &o)
				continue
			}

			// option with only one key (no value)
			// ie: omitempty from json:"id,omitempty"
			regex = *regexp.MustCompile(`^(.*)$`)
			matches = regex.FindAllStringSubmatch(op, -1)
			if len(matches) == 1 {
				o := Option{
					Key: matches[0][1],
				}
				options = append(options, &o)
				continue
			}
		}
	}

	return
}
