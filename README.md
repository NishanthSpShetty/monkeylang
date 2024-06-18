## monkeylang


Interpreter written in Go.

## Build

```
make build
```

## Usage

To run you can use the binary built in previous step
```
monkey
```

OR

```
make run
```

This will start monkey repl

## language

* C-like syntax
* Dynamic typing
* Variable bindings
* Integers and booleans
* Arithmetic expressions
* Arrays and maps
* Built-in functions
* First-class and higher-order functions
* Closures - TODO


Sample snippets
```
let age = 1;
let name = "Monkey";
let result = 10 * (20 / 2);
let add = fn(a, b) { return a + b; };
add(1, 2);

let fibonacci = fn(x) {
    if (x == 0) {
        0
    } else {
        if (x == 1) {
            1
        } else {
            fibonacci(x - 1) + fibonacci(x - 2);
        }
    }
};

fn(x,Y){ x+y ;} (3,4) 
let firstname = "Nishanth"
let lastname= "Shetty"
let fullname = firstname + " " + lastname

let ar = [1,2, "hello", 2/2, fn(x) { x* 2; }(3)]
ar[2]
len(ar)
let a = {1:"1", 2:"2", 3:"3"};
a[1]
a[4]
let map = {1:"one", 2:"two", 3:"three"}
map[1]
```
_all the above snippets are valid monkey lang, try executing them in a repl_


## Developement and Testing

### Test
```
make test   
```

## Reference

This implementation is based on the Thorsten Ball Book. [Writing Interpreter In Go](https://edu.anarcho-copy.org/Programming%20Languages/Go/writing%20an%20INTERPRETER%20in%20go.pdf)
