# Tags

A simple tags parser for golang's struct

## Install

```shell
go get github.com/go-mods/tags
```

## Example 
```go
package main

import (
    "fmt"
    "github.com/go-mods/tags"
    "reflect"
)

func main() {
    type Employee struct {
        Id   int    `json:"id" xml:"id" excel:"id"`
        Name string `json:"name,string" xml:"name" excel:"name"`
        Age int `json:"age,omitempty" xml:"age" excel:"column:age"`
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

// OUTPUT

// Tags for field: Id
{"Tag":"json:\"id\"","Key":"json","Value":"id","Name":"id","Options":null}
{"Tag":"xml:\"id\"","Key":"xml","Value":"id","Name":"id","Options":null}
{"Tag":"excel:\"id\"","Key":"excel","Value":"id","Name":"id","Options":null}

// Tags for field: Name
{"Tag":"json:\"name,string\"","Key":"json","Value":"name,string","Name":"name","Options":[{"Key":"string","Value":""}]}
{"Tag":"xml:\"name\"","Key":"xml","Value":"name","Name":"name","Options":null}
{"Tag":"excel:\"name\"","Key":"excel","Value":"name","Name":"name","Options":null}

// Tags for field: Age
{"Tag":"json:\"age,omitempty\"","Key":"json","Value":"age,omitempty","Name":"age","Options":[{"Key":"omitempty","Value":""}]}
{"Tag":"xml:\"age\"","Key":"xml","Value":"age","Name":"age","Options":null}
{"Tag":"excel:\"column:age\"","Key":"excel","Value":"column:age","Name":"","Options":[{"Key":"column","Value":"age"}]}

--- PASS: TestExample (0.00s)
PASS


Process finished with the exit code 0

--- PASS: TestExample (0.00s)
PASS


Process finished with the exit code 0

```
