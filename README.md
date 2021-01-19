
## About
yamltostruct lets you define golang types in yaml format and unmarshals them into [AST declarations](https://golang.org/pkg/go/ast/#File).
***
### in:
```
name:
  first: string
  last: string
id: string
person:
  name: name
  age: int
  id: id
```
### out:
```
type id string
type name struct {
	first	string
	last	string
}
type person struct {
	age	int
	id	id
	name	name
}
```
<br/>

## Usage
```
package main

import (
        "fmt"
        "go/ast"
        "go/printer"
        "go/token"
        "os"

        "github.com/Java-Jonas/yamltostruct"
)

func logErrs(errs []error) {
        for _, err := range errs {
                fmt.Println(err.Error())
        }
}

func main() {
        yamlData := []byte(`
name:
  first: string
  last: string
id: string
person:
  name: name
  age: int
  id: id
`)

        decls, errs := yamltostruct.Unmarshal(yamlData)

        if len(errs) > 0 {
                logErrs(errs)
                return
        }

        // pretty print declaration information
        ast.Print(token.NewFileSet(), decls)
        fmt.Println("")
        // print generated code
        printer.Fprint(os.Stdout, token.NewFileSet(), decls)
}
```
<br/>


## Validation Error Messages
<br/> 

### structural:

| Error | Text | Meaning |
|---|---------|----------|
| ErrIllegalValue | value assigned to key "{KeyName}" in "{ParentObject}" is invalid | An invalid value was defined (nil, "", List, Object in Object). |
<br/> 

### syntactical:
| Error | Text | Meaning |
|---|---------|----------|
| ErrIllegalTypeName | illegal type name "{KeyName}" in "{ParentObject}" | A type was named without adhering to go's syntax limitations (e.g. "fo$o", "func", "<-+"). |
| ErrInvalidValueString | value "{ValueString}" assigned to "{KeyName}" in "{ParantObject}" is invalid | An invalid value was assigned to a key |
<br/> 

### logical:
| Error | Text | Meaning |
|---|---------|----------|
| ErrTypeNotFound | type with name "{TypeName}" in "{ParentObject}" was not found | A type was referenced as value but not defined anywhere in the YAML document. |
| ErrRecursiveTypeUsage | illegal recursive type detected for "{RecurringKeyNames}" | A recursive type was defined. |
| ErrInvalidMapKey | "{MapKey}" in "{ValueString}" is not a valid map key | An uncomparable type was chosen as map key. |
<br/> 


## Motivation
This was a project for me to get more comfortable with [TDD](https://en.wikipedia.org/wiki/Test-driven_development) and golang. I don't think the library itself is very useful and I am aware that there are ways I could have achieved the same functionality with a lot less effort. However this was more of a fun/educational project and it has fulfilled its purpose.
<br/>
<br/>

## Test of Time
In my humble opinion, maintenance is the most important feature in any software project. Revisiting old code and being content with the way it looks is pretty much a meme within the developer community. Every developer knows the struggles of working with hardly maintainable code, even if said developer is the author themselves. <br/>
Critiquing the maintainability of your own code is hard, unless enough time has passed and you don't completely remember the details of the logic. This is why I revisited this project to implement a new feature. <br/>
Here are a few things I noticed while revisiting this project:
- the logic was a lot more complex than I remembered it
- despite my urge to use TDD I had to get a rough mental image of the logic before I was able to write meaningful tests

### TODO:
- implement array support
