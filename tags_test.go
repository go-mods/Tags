package tags_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-mods/tags"
)

// checkOptions fails the test when got differs from want
func checkOptions(t *testing.T, got, want []*tags.Option) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("Options = %s, want %s", formatOptions(got), formatOptions(want))
	}
	for i := range got {
		if got[i].Key != want[i].Key || got[i].Value != want[i].Value {
			t.Errorf("Options = %s, want %s", formatOptions(got), formatOptions(want))
			return
		}
	}
}

func formatOptions(options []*tags.Option) string {
	parts := make([]string, 0, len(options))
	for _, o := range options {
		parts = append(parts, fmt.Sprintf("{%q: %#v}", o.Key, o.Value))
	}
	return "[" + strings.Join(parts, " ") + "]"
}

func TestParse_keyAndValue(t *testing.T) {
	tests := []struct {
		name  string
		tag   string
		key   string
		value string
	}{
		{name: "name and option", tag: `json:"id,omitempty"`, key: "json", value: "id,omitempty"},
		{name: "option only", tag: `json:",omitempty"`, key: "json", value: ",omitempty"},
		{name: "dash", tag: `json:"-"`, key: "json", value: "-"},
		{name: "semicolon", tag: `gorm:"embedded;embeddedPrefix:author_"`, key: "gorm", value: "embedded;embeddedPrefix:author_"},
		{name: "constraint", tag: `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`, key: "gorm", value: "constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"},
		{name: "escaped quotes", tag: `excel:"q,default:\"x\""`, key: "excel", value: `q,default:"x"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tgs, err := tags.Parse(tt.tag)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if len(tgs) != 1 {
				t.Fatalf("Parse() returned %d tags, want 1", len(tgs))
			}
			if tgs[0].Tag != tt.tag || tgs[0].Key != tt.key || tgs[0].Value != tt.value {
				t.Errorf("Parse() = {Tag: %q, Key: %q, Value: %q}, want {%q, %q, %q}",
					tgs[0].Tag, tgs[0].Key, tgs[0].Value, tt.tag, tt.key, tt.value)
			}
		})
	}
}

func TestParse_empty(t *testing.T) {
	for _, s := range []string{"", "  "} {
		tgs, err := tags.Parse(s)
		if tgs != nil || err != nil {
			t.Errorf("Parse(%q) = %v, %v, want nil, nil", s, tgs, err)
		}
	}
}

func TestParse_nameAndOptions(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		want    string
		options []*tags.Option
	}{
		{name: "name only", tag: `test:"key"`, want: "key"},
		{name: "key value", tag: `test:"key:value"`, options: []*tags.Option{{Key: "key", Value: "value"}}},
		{name: "key values with comma", tag: `test:"key1:value1,key2:value2"`, options: []*tags.Option{{Key: "key1", Value: "value1"}, {Key: "key2", Value: "value2"}}},
		{name: "key values with semicolon", tag: `test:"key1:value1;key2:value2"`, options: []*tags.Option{{Key: "key1", Value: "value1"}, {Key: "key2", Value: "value2"}}},
		{name: "json name", tag: `json:"id"`, want: "id"},
		{name: "json omitempty", tag: `json:"id,omitempty"`, want: "id", options: []*tags.Option{{Key: "omitempty"}}},
		{name: "json without name", tag: `json:",omitempty"`, options: []*tags.Option{{Key: "omitempty"}}},
		{name: "json two options", tag: `json:"id,omitempty,default"`, want: "id", options: []*tags.Option{{Key: "omitempty"}, {Key: "default"}}},
		{name: "json dash", tag: `json:"-"`, want: "-"},
		{name: "gorm semicolon", tag: `gorm:"embedded;embeddedPrefix:author_"`, want: "embedded", options: []*tags.Option{{Key: "embeddedPrefix", Value: "author_"}}},
		{name: "gorm without name", tag: `gorm:",embedded;embeddedPrefix:author_"`, options: []*tags.Option{{Key: "embedded"}, {Key: "embeddedPrefix", Value: "author_"}}},
		{name: "gorm constraint", tag: `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`, options: []*tags.Option{{Key: "constraint", Value: "OnUpdate:CASCADE,OnDelete:SET NULL;"}}},
		{name: "gorm option with space", tag: `gorm:"name;type:varchar(255);not null"`, want: "name", options: []*tags.Option{{Key: "type", Value: "varchar(255)"}, {Key: "not null"}}},
		{name: "gorm enum is not split", tag: `gorm:"type:enum('a', 'b')"`, options: []*tags.Option{{Key: "type", Value: "enum('a', 'b')"}}},
		{name: "split semicolon", tag: `excel:"array,split:;"`, want: "array", options: []*tags.Option{{Key: "split", Value: ";"}}},
		{name: "split comma semicolon", tag: `excel:"array,split:,;"`, want: "array", options: []*tags.Option{{Key: "split", Value: ",;"}}},
		{name: "split three separators", tag: `excel:"array,split:,;|"`, want: "array", options: []*tags.Option{{Key: "split", Value: ",;|"}}},
		{name: "name with space", tag: `excel:"First Name"`, want: "First Name"},
		{name: "non-ASCII name", tag: `excel:"Año"`, want: "Año"},
		{name: "non-ASCII name and value", tag: `excel:"Größe,default:keine Größe"`, want: "Größe", options: []*tags.Option{{Key: "default", Value: "keine Größe"}}},
		{name: "name with space and option", tag: `excel:"Postal code,required"`, want: "Postal code", options: []*tags.Option{{Key: "required"}}},
		{name: "non-ASCII name and split", tag: `excel:"Números;split:|"`, want: "Números", options: []*tags.Option{{Key: "split", Value: "|"}}},
		{name: "column with space", tag: `excel:"column:Last Name"`, options: []*tags.Option{{Key: "column", Value: "Last Name"}}},
		{name: "non-ASCII option", tag: `excel:"x,a,état"`, want: "x", options: []*tags.Option{{Key: "a"}, {Key: "état"}}},
		{name: "non-ASCII option key", tag: `excel:"x,clé:val"`, want: "x", options: []*tags.Option{{Key: "clé", Value: "val"}}},
		{name: "time format keeps the next option", tag: `excel:"Time,format:15:04,default:now"`, want: "Time", options: []*tags.Option{{Key: "format", Value: "15:04"}, {Key: "default", Value: "now"}}},
		{name: "time layout with letters", tag: `excel:"Time,format:hh:mm,default:now"`, want: "Time", options: []*tags.Option{{Key: "format", Value: "hh:mm"}, {Key: "default", Value: "now"}}},
		{name: "escaped quotes", tag: `excel:"q,default:\"x\""`, want: "q", options: []*tags.Option{{Key: "default", Value: `"x"`}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tgs, err := tags.Parse(tt.tag)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if len(tgs) != 1 {
				t.Fatalf("Parse() returned %d tags, want 1", len(tgs))
			}
			if tgs[0].Name != tt.want {
				t.Errorf("Name = %q, want %q", tgs[0].Name, tt.want)
			}
			checkOptions(t, tgs[0].Options, tt.options)
		})
	}
}

func TestParse_severalTags(t *testing.T) {
	tgs, err := tags.Parse(`json:"id,omitempty,default" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	want := []*tags.Tag{
		{
			Tag:     `json:"id,omitempty,default"`,
			Key:     "json",
			Value:   "id,omitempty,default",
			Name:    "id",
			Options: []*tags.Option{{Key: "omitempty"}, {Key: "default"}},
		},
		{
			Tag:     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`,
			Key:     "gorm",
			Value:   "constraint:OnUpdate:CASCADE,OnDelete:SET NULL;",
			Options: []*tags.Option{{Key: "constraint", Value: "OnUpdate:CASCADE,OnDelete:SET NULL;"}},
		},
	}
	if len(tgs) != len(want) {
		t.Fatalf("Parse() returned %d tags, want %d", len(tgs), len(want))
	}
	for i, tag := range tgs {
		t.Run(want[i].Key, func(t *testing.T) {
			if tag.Tag != want[i].Tag || tag.Key != want[i].Key || tag.Value != want[i].Value || tag.Name != want[i].Name {
				t.Errorf("Parse() = {%q, %q, %q, %q}, want {%q, %q, %q, %q}",
					tag.Tag, tag.Key, tag.Value, tag.Name, want[i].Tag, want[i].Key, want[i].Value, want[i].Name)
			}
			checkOptions(t, tag.Options, want[i].Options)
		})
	}
}

func TestLookup(t *testing.T) {
	type Employee struct {
		ID    int    `json:"id" xml:"id"`
		Name  string `json:"name,string" xml:"name"`
		Quote string `excel:"q,default:\"x\""`
	}
	typ := reflect.TypeOf(Employee{})

	tests := []struct {
		name    string
		field   int
		key     string
		want    string
		options []*tags.Option
	}{
		{name: "name only", field: 0, key: "json", want: "id"},
		{name: "name and option", field: 1, key: "json", want: "name", options: []*tags.Option{{Key: "string"}}},
		{name: "second key", field: 1, key: "xml", want: "name"},
		{name: "escaped quotes", field: 2, key: "excel", want: "q", options: []*tags.Option{{Key: "default", Value: `"x"`}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag := tags.Lookup(typ.Field(tt.field), tt.key)
			if tag == nil {
				t.Fatal("Lookup() = nil")
			}
			if tag.Key != tt.key || tag.Name != tt.want {
				t.Errorf("Lookup() = {Key: %q, Name: %q}, want {%q, %q}", tag.Key, tag.Name, tt.key, tt.want)
			}
			checkOptions(t, tag.Options, tt.options)
		})
	}

	t.Run("missing key", func(t *testing.T) {
		if tag := tags.Lookup(typ.Field(0), "excel"); tag != nil {
			t.Errorf("Lookup() = %v, want nil", tag)
		}
	})
}

func TestGetOption(t *testing.T) {
	tgs, _ := tags.Parse(`excel:"Name,required,default:anonymous"`)
	tag := tgs[0]
	if !tag.HasOption("required") || tag.GetOption("required").Value != nil {
		t.Errorf("option required: HasOption = %v, GetOption = %v", tag.HasOption("required"), tag.GetOption("required"))
	}
	if o := tag.GetOption("default"); o == nil || o.Value != "anonymous" {
		t.Errorf("GetOption(default) = %v, want the value anonymous", o)
	}
	if tag.HasOption("split") || tag.GetOption("split") != nil {
		t.Error("HasOption(split) = true, want false")
	}
}

func ExampleParse() {
	type Employee struct {
		ID   int    `json:"id" excel:"Id,required"`
		Name string `json:"name,omitempty" excel:"First Name,default:anonymous"`
	}

	field, _ := reflect.TypeOf(Employee{}).FieldByName("Name")
	tgs, _ := tags.Parse(string(field.Tag))
	for _, tag := range tgs {
		fmt.Printf("%s: name %q", tag.Key, tag.Name)
		for _, o := range tag.Options {
			fmt.Printf(", option %s=%v", o.Key, o.Value)
		}
		fmt.Println()
	}
	// Output:
	// json: name "name", option omitempty=<nil>
	// excel: name "First Name", option default=anonymous
}

func ExampleLookup() {
	type Employee struct {
		Name string `json:"name" excel:"column:First Name,default:anonymous"`
	}

	field, _ := reflect.TypeOf(Employee{}).FieldByName("Name")
	tag := tags.Lookup(field, "excel")
	fmt.Println(tag.GetOption("column").Value, "/", tag.GetOption("default").Value)
	fmt.Println(tags.Lookup(field, "xml") == nil)
	// Output:
	// First Name / anonymous
	// true
}

func FuzzParse(f *testing.F) {
	for _, seed := range []string{
		`json:"id,omitempty"`,
		`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`,
		`excel:"array,split:,;|" json:"-"`,
		`excel:"q,default:\"x\""`,
		`excel:"x,a,état"`,
		"0 \u000b:\"\"", // vertical tab before the colon: not a key
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		tgs, err := tags.Parse(s)
		if err != nil {
			t.Fatalf("Parse(%q) error = %v", s, err)
		}
		for _, tag := range tgs {
			if tag.Key == "" {
				t.Errorf("Parse(%q) returned a tag without key: %+v", s, tag)
			}
			for _, o := range tag.Options {
				if o == nil {
					t.Errorf("Parse(%q) returned a nil option", s)
				}
			}
		}
	})
}
