package tags_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-mods/tags"
	"reflect"
	"testing"
)

var simpleTagTest = []struct {
	tag   string
	key   string
	value string
}{
	{"", "", ""},
	{"json:\"id,omitempty\"", "json", "id,omitempty"},
	{"json:\",omitempty\"", "json", ",omitempty"},
	{"json:\"-\"", "json", "-"},
	{"gorm:\"embedded;embeddedPrefix:author_\"", "gorm", "embedded;embeddedPrefix:author_"},
	{"gorm:\",embedded;embeddedPrefix:author_\"", "gorm", ",embedded;embeddedPrefix:author_"},
	{"gorm:\"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;\"", "gorm", "constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"},
}

var simpleValueTest = []struct {
	tag     string
	name    string
	options []*tags.Option
}{
	{"test:\"key\"", "key", nil},
	{"test:\"key:value\"", "", []*tags.Option{{Key: "key", Value: "value"}}},
	{"test:\"key1:value1,key2:value2\"", "", []*tags.Option{{Key: "key1", Value: "value1"}, {Key: "key2", Value: "value2"}}},
	{"test:\"key1:value1;key2:value2\"", "", []*tags.Option{{Key: "key1", Value: "value1"}, {Key: "key2", Value: "value2"}}},
	{"json:\"id\"", "id", nil},
	{"json:\"id,omitempty\"", "id", []*tags.Option{{Key: "omitempty"}}},
	{"json:\",omitempty\"", "", []*tags.Option{{Key: "omitempty"}}},
	{"json:\"id,omitempty,default\"", "id", []*tags.Option{{Key: "omitempty"}, {Key: "default"}}},
	{"gorm:\"embedded;embeddedPrefix:author_\"", "embedded", []*tags.Option{{Key: "embeddedPrefix", Value: "author_"}}},
	{"gorm:\",embedded;embeddedPrefix:author_\"", "", []*tags.Option{{Key: "embedded"}, {Key: "embeddedPrefix", Value: "author_"}}},
	{"gorm:\"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;\"", "", []*tags.Option{{Key: "constraint", Value: "OnUpdate:CASCADE,OnDelete:SET NULL;"}}},
	{"gorm:\"name;type:varchar(255);not null\"", "name", []*tags.Option{{Key: "type", Value: "varchar(255)"}, {Key: "not null"}}},
	{"excel:\"array,split:;\"", "array", []*tags.Option{{Key: "split", Value: ";"}}},
	{"excel:\"array,split:,;\"", "array", []*tags.Option{{Key: "split", Value: ",;"}}},
	{"excel:\"array,split:,;|\"", "array", []*tags.Option{{Key: "split", Value: ",;|"}}},
}

var complexTagTest = struct {
	str  string
	tags []*tags.Tag
}{
	str: "json:\"id,omitempty,default\" gorm:\"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;\"",
	tags: []*tags.Tag{
		{
			Tag:   "json:\"id,omitempty,default\"",
			Key:   "json",
			Value: "id,omitempty,default",
			Name:  "id",
			Options: []*tags.Option{
				{Key: "omitempty", Value: nil},
				{Key: "default", Value: nil},
			},
		},
		{
			Tag:   "gorm:\"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;\"",
			Key:   "gorm",
			Value: "constraint:OnUpdate:CASCADE,OnDelete:SET NULL;",
			Options: []*tags.Option{
				{Key: "constraint", Value: "OnUpdate:CASCADE,OnDelete:SET NULL;"},
			},
		},
	},
}

func TestSimpleTag(t *testing.T) {
	for _, test := range simpleTagTest {
		tgs, err := tags.Parse(test.tag)

		if err == nil && tgs != nil {
			if tgs[0].Tag != test.tag {
				t.Errorf("Parse() got = %v, want %v", tgs[0].Tag, test.tag)
			}
			if tgs[0].Key != test.key {
				t.Errorf("Parse() got = %v, want %v", tgs[0].Key, test.key)
			}
			if tgs[0].Value != test.value {
				t.Errorf("Parse() got = %v, want %v", tgs[0].Value, test.value)
			}
		} else if err != nil {
			t.Errorf("Parse() error = %v", err)
		}
	}
}

func TestSimpleValue(t *testing.T) {
	for _, test := range simpleValueTest {
		tgs, err := tags.Parse(test.tag)

		if err == nil && tgs != nil {
			if len(tgs) != 1 {
				t.Errorf("Parse() got = %v, want %v", len(tgs), 1)
			}
			if tgs[0].Name != test.name {
				t.Errorf("Parse() got = %v, want %v", tgs[0].Name, test.name)
			}
			if len(tgs[0].Options) != len(test.options) {
				t.Errorf("Parse() got = %v, want %v", len(tgs[0].Options), len(test.options))
			}
			for i, opt := range tgs[0].Options {
				if opt.Key != test.options[i].Key {
					t.Errorf("Parse() got = %v, want %v", opt.Key, test.options[i].Key)
				}
				if opt.Value != test.options[i].Value {
					t.Errorf("Parse() got = %v, want %v", opt.Value, test.options[i].Value)
				}
			}
		} else if err != nil {
			t.Errorf("Parse() error = %v", err)
		}
	}
}

func TestComplexTag(t *testing.T) {
	tgs, err := tags.Parse(complexTagTest.str)

	if err == nil && tgs != nil {
		if len(tgs) != len(complexTagTest.tags) {
			t.Errorf("Parse() got = %v, want %v", len(tgs), len(complexTagTest.tags))
		}
		for i, tag := range tgs {
			if tag.Tag != complexTagTest.tags[i].Tag {
				t.Errorf("Parse() got = %v, want %v", tag.Tag, complexTagTest.tags[i].Tag)
			}
			if tag.Key != complexTagTest.tags[i].Key {
				t.Errorf("Parse() got = %v, want %v", tag.Key, complexTagTest.tags[i].Key)
			}
			if tag.Value != complexTagTest.tags[i].Value {
				t.Errorf("Parse() got = %v, want %v", tag.Value, complexTagTest.tags[i].Value)
			}
			if tag.Name != complexTagTest.tags[i].Name {
				t.Errorf("Parse() got = %v, want %v", tag.Name, complexTagTest.tags[i].Name)
			}
			if len(tag.Options) != len(complexTagTest.tags[i].Options) {
				t.Errorf("Parse() got = %v, want %v", len(tag.Options), len(complexTagTest.tags[i].Options))
			}
			for j, opt := range tag.Options {
				if opt.Key != complexTagTest.tags[i].Options[j].Key {
					t.Errorf("Parse() got = %v, want %v", opt.Key, complexTagTest.tags[i].Options[j].Key)
				}
				if opt.Value != complexTagTest.tags[i].Options[j].Value {
					t.Errorf("Parse() got = %v, want %v", opt.Value, complexTagTest.tags[i].Options[j].Value)
				}
			}
		}
	}
}

func TestExample(t *testing.T) {

	type Employee struct {
		Id   int    `json:"id" xml:"id" excel:"id"`
		Name string `json:"name,string" xml:"name" excel:"name"`
		Age  int    `json:"age,omitempty" xml:"age" excel:"column:age"`
	}

	// Loop throw all fields
	for i := 0; i < reflect.TypeOf(Employee{}).NumField(); i++ {
		// get the field
		field := reflect.TypeOf(Employee{}).Field(i)
		// get the tag field
		tag := field.Tag

		// parse it
		tgs, err := tags.Parse(string(tag))
		if err != nil {
			t.Errorf("Parse() error = %v", err)
		}

		// iterate over all tags
		fmt.Printf("// Tags for field: %s", field.Name)
		for _, t := range tgs {
			out, _ := json.Marshal(t)
			fmt.Println(string(out))
		}
		fmt.Println("")
	}
}

func TestLookup(t *testing.T) {

	type Employee struct {
		Id   int    `json:"id" xml:"id"`
		Name string `json:"name,string" xml:"name"`
		Age  int    `json:"age,omitempty" xml:"age"`
	}

	if tags.Lookup(reflect.TypeOf(Employee{}).Field(0), "json") == nil {
		t.Error("Lookup() got = nil")
	}

	if tags.Lookup(reflect.TypeOf(Employee{}).Field(1), "json") == nil {
		t.Error("Lookup() got = nil")
	}

	if tags.Lookup(reflect.TypeOf(Employee{}).Field(1), "xml") == nil {
		t.Error("Lookup() got = nil")
	}
}
