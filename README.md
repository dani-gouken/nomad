#Nomad

## A simple interpreter build in Go

Nomad is an interpreted programming language

It compiles the program to bytecode which is later interpreted by a stack-based virtual machine

Examples are in the `examples` folder

### example
#### Fib sequence
```
auto fib :: func(int n) int {
    if n = 0 {
        return 0
    } elif n = 1 {
        return 1
    } else {
        return fib(n-2) + fib(n-1)
    }
}
print fib(16) // 967
```
### 2nd degree equation solver

```
type Equation :: {
    auto a :: 0.0
    float b :: 0.0
    float c :: 0.0
}

auto solve_2nd :: func(Equation eq) [float] {
    auto delta ::  pow(eq.b, 2) - (4.0 * eq.a * eq.c)
    if delta < 0.0 {
        return [float]{};
    }
    if delta > 0.0 {
        return [float]{
            (-eq.b - sqrt(delta))/(2.0 * eq.a), 
            (-eq.b + sqrt(delta))/(2.0 * eq.a)
        }
    }
    return [float]{ (-eq.b)/(2.0 * eq.a)}
} 

auto eq :: new Equation{
    a :: 4.0, b :: 0.0, c :: -16.0,
}
print solve_2nd(eq)
```

## Test it

`go run main.go examples/fib.nd`

## Todo
- [x] Variables
- [x] Math
- [x] Literal types
- [x] Control flow
- [x] Array
- [x] Object
- [x] Advanced types
- [x] type checking
- [x] function
- [ ] interface with go
#
