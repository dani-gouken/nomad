#Nomad

## A simple interpreter build in Go

Nomad is a toy programming language i built to learn go and programming language theory (lexical analysis, parsing, ...)

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
