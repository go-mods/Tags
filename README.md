# Tags

A simple tags parser for golang's struct

## Install

```shell
go get github.com/go-mods/tags
```

## Example 
```go
type Employee struct {
    Id   int    `json:"id" xml:"id"`
    Name string `json:"name,string" xml:"name"`
    Age  int    `json:"age,omitempty,int" xml:"age"`
}

// Loop throw all fields
for i := 0; i < reflect.TypeOf(Employee{}).NumField(); i++ {
    // get the field
    field := reflect.TypeOf(Employee{}).Field(i)
    // get the tag field
    tag := field.Tag

    // parse it
    tags, err := tags.Parse(string(tag))
    if err != nil {
        panic(err)
    }

    // iterate over all tags
    for _, t := range tags {
        fmt.Printf("tag: %+v\n", t)
    }
}
```
