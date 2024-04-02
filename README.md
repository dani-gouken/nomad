#Nomad

## A simple interpreter build in Go

Nomad is a toy programming language i built to learn go and programming language theory (lexical analysis, parsing, ...)

It compiles the program to bytecode which is later interpreted by a stack-based virtual machine

Examples are in the `examples` folder

### example

```
[string] parts :: [string]{"daniel", "nghokeng", "stéphane"};
string fullname :: " ";

for int i :: 0; i <= 2; i++ {
    fullname  = fullname + parts[i];
    if i < 2 {
        fullname  = " " + fullname
    }
}

print fullname
```

## Test it

`go run main.go somefile.nd`

## Todo
- [x] Variables
- [x] Math
- [x] Literal types (inte)
- [x] Literal types
- [x] Control flow
- [x] Array
- [ ] Object

#
