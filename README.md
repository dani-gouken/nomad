#Nomad

## A simple interpreter build in Go

Nomad is an interpreted programming language

It compiles the program to bytecode which is later interpreted by a stack-based virtual machine

Examples are in the `examples` folder

### example

```
[string] parts :: [string]{"daniel", "nghokeng", "stéphane"}
string fullname :: ""
int size :: len parts

for int i :: 0; i < size; i++ {
    fullname :: fullname + parts[i];
    bool last :: i = size - 1
    if !last {
        fullname :: fullname + " "
    }
}
print fullname
```
### structs

```
type Header  :: {
    string name :: ""
    string value :: ""
}

type HttpResponse :: {
    int status :: 0
    string body :: ""
    [Header] headers :: [Header]{
        new Header { name :: "Content-Type", value :: "application/json" }
        new Header { name :: "Content-Lenght", value :: "0" }
    }
}

type HttpStatus :: {
    int OK :: 200 
    int ValidationError :: 422 
}

HttpStatus status :: new HttpStatus{}

HttpResponse res :: new HttpResponse{
    status :: status.OK
    body :: "{\"msg\": \"Hello world\"}"
}

print res.status
print res.body
```

## Test it

`go run main.go somefile.nd`

## Todo
- [x] Variables
- [x] Math
- [x] Literal types
- [x] Control flow
- [x] Array
- [ ] Object
- [ ] Advanced types

#
