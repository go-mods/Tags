# Tags

[![Go Reference](https://pkg.go.dev/badge/github.com/go-mods/tags.svg)](https://pkg.go.dev/github.com/go-mods/tags)
[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-mods/tags)](https://goreportcard.com/report/github.com/go-mods/tags)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A parser for Go struct tags. Each `key:"value"` pair is split into a **name** and **options**, so that a package
reading its own tags (`excel:"First Name,default:anonymous"`) gets `First Name` and the option `default=anonymous`
without writing its own parser. No dependency outside the standard library.

## Installation

Prerequisite: Go 1.20 or later.

```shell
go get github.com/go-mods/tags
```

## Usage

```go
package main

import (
    "fmt"
    "reflect"

    "github.com/go-mods/tags"
)

type Employee struct {
    ID   int    `json:"id" excel:"Id,required"`
    Name string `json:"name,omitempty" excel:"First Name,default:anonymous"`
}

func main() {
    // Parse every tag of a field
    field, _ := reflect.TypeOf(Employee{}).FieldByName("Name")
    tgs, _ := tags.Parse(string(field.Tag))
    for _, tag := range tgs {
        fmt.Printf("%s: name %q", tag.Key, tag.Name)
        for _, o := range tag.Options {
            fmt.Printf(", option %s=%v", o.Key, o.Value)
        }
        fmt.Println()
    }
    // json: name "name", option omitempty=<nil>
    // excel: name "First Name", option default=anonymous

    // Read the tag of one key: nil when the field has no tag for this key
    tag := tags.Lookup(field, "excel")
    fmt.Println(tag.HasOption("required"), tag.GetOption("default").Value)
    // false anonymous
}
```

| Function or method | Returns |
|---|---|
| `tags.Parse(s string) ([]*tags.Tag, error)` | One `Tag` per `key:"value"` pair, nil for an empty string. The error is always nil. |
| `tags.Lookup(field reflect.StructField, key string) *tags.Tag` | The tag of the field for this key, or nil |
| `(*Tag).HasOption(key string) bool` | Whether the tag has the option |
| `(*Tag).GetOption(key string) *tags.Option` | The first option with this key, or nil. `Option.Value` is nil for an option without value. |

## Splitting rules

| Value | Name | Options |
|---|---|---|
| `id` | `id` | — |
| `id,omitempty` | `id` | `omitempty` |
| `embedded;embeddedPrefix:author_` | `embedded` | `embeddedPrefix=author_` |
| `First Name,default:anonymous` | `First Name` | `default=anonymous` |
| `Año;split:\|` | `Año` | `split=\|` |
| `,omitempty` | — | `omitempty` |
| `column:Last Name` | — | `column=Last Name` |
| `array,split:,;\|` | `array` | `split=,;\|` |
| `type:enum('a', 'b')` | — | `type=enum('a', 'b')` |
| `Time,format:15:04,default:now` | `Time` | `format=15:04`, `default=now` |
| `constraint:OnUpdate:CASCADE,OnDelete:SET NULL;` | — | `constraint=OnUpdate:CASCADE,OnDelete:SET NULL;` |

- The **name** is the part before the first `,` or `;`. It may hold spaces and non-ASCII letters, but no `,` `;` or `:`. A value holding a `:` before any `,` or `;` has no name, only options.
- **Options** are separated by `,` or `;`. An option is a key (`omitempty`) or a `key:value` pair; the value is the string after the first `:`.
- A `,` or `;` is a separator only when a letter, a digit or `_` follows it. Otherwise it belongs to the value: `split:,;|` keeps `,;|` and `enum('a', 'b')` is not split.
- The `constraint` option keeps all its `key:value` pairs as a single value.
- As in `reflect.StructTag`, a value may hold escaped quotes: `` `excel:"q,default:\"x\""` `` gives `default="x"`.

## Development

Prerequisite: [Task](https://taskfile.dev/installation/).

```shell
task go-tools:install   # golangci-lint and editorconfig-checker, at their pinned versions
task go:default         # clean, format, vet, tidy, lint and test
task go:test
task go:lint
```

## License

MIT — see [LICENSE](LICENSE).
