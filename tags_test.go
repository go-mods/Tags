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
	{"gorm:\"embedded;embeddedPrefix:author_\"", "gorm", "embedded;embeddedPrefix:author_"},
	{"gorm:\"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;\"", "gorm", "constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"},
}

var simpleValueTest = []struct {
	tag     string
	name    string
	options []*tags.Option
}{
	{"json:\"id\"", "id", nil},
	{"json:\"id,omitempty\"", "id", []*tags.Option{{Key: "omitempty"}}},
	{"json:\",omitempty\"", "", []*tags.Option{{Key: "omitempty"}}},
	{"json:\"id,omitempty,default\"", "id", []*tags.Option{{Key: "omitempty"}, {Key: "default"}}},
	{"gorm:\"embedded;embeddedPrefix:author_\"", "embedded", []*tags.Option{{Key: "embeddedPrefix", Value: "author_"}}},
	{"gorm:\"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;\"", "", []*tags.Option{{Key: "constraint", Value: "OnUpdate:CASCADE"}, {Key: "OnDelete", Value: "SET NULL;"}}},
	{"test:\"key:value\"", "", []*tags.Option{{Key: "key", Value: "value"}}},
	{"test:\"key1:value1,key2:value2\"", "", []*tags.Option{{Key: "key1", Value: "value1"}, {Key: "key2", Value: "value2"}}},
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
				{Key: "omitempty", Value: ""},
				{Key: "default", Value: ""},
			},
		},
		{
			Tag:   "gorm:\"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;\"",
			Key:   "gorm",
			Value: "constraint:OnUpdate:CASCADE,OnDelete:SET NULL;",
			Name:  "constraint",
			Options: []*tags.Option{
				{Key: "OnUpdate", Value: "CASCADE"},
				{Key: "OnDelete", Value: "SET NULL;"},
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
			if tgs[0].Tag != test.tag {
				t.Errorf("Parse() got = %v, want %v", tgs[0].Tag, test.tag)
			}
			if tgs[0].Name != test.name {
				t.Errorf("Parse() got = %v, want %v", tgs[0].Key, test.name)
			}
			if test.options != nil {
				if len(test.options) >= 1 && tgs[0].Options[0].Key != test.options[0].Key {
					t.Errorf("Parse() got = %v, want %v", tgs[0].Options[0].Key, test.options[0].Key)
				}
				if len(test.options) >= 1 && tgs[0].Options[0].Value != test.options[0].Value {
					t.Errorf("Parse() got = %v, want %v", tgs[0].Options[0].Value, test.options[0].Value)
				}
				if len(test.options) >= 2 && tgs[0].Options[1].Key != test.options[1].Key {
					t.Errorf("Parse() got = %v, want %v", tgs[0].Options[1], test.options[1])
				}
				if len(test.options) >= 2 && tgs[0].Options[1].Value != test.options[1].Value {
					t.Errorf("Parse() got = %v, want %v", tgs[0].Options[1], test.options[1])
				}
			}

		} else if err != nil {
			t.Errorf("Parse() error = %v", err)
		}
	}
}

func TestComplexTag(t *testing.T) {
	_, err := tags.Parse(complexTagTest.str)

	if err != nil {
		t.Errorf("Parse() error = %v", err)
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
			panic(err)
		}

		// iterate over all tags
		fmt.Println(fmt.Sprintf("// Tags for field: %s", field.Name))
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
