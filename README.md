# PotatoLang 

A strictly-typed, interpreted programming language built in Go. Potatolang features Go-style syntax, immutable-by-default variables, multiple return values, and strict runtime type enforcement. 

It is designed to be safe and predictable, gracefully catching type mismatches and assignment errors at runtime without crashing the underlying engine.

## Features

* **Immutable by Default:** Variables cannot be reassigned unless explicitly marked with the `mut` keyword.
* **Multiple Return Values:** Functions can return multiple values, which can be safely unpacked into multiple variables.
* **Strict Type Checking:** Variables lock in their type upon initialization. Reassigning a variable to a different type throws a graceful runtime error.
* **Safe Arity Validation:** Function calls are strictly checked against their defined parameters to prevent runtime panics.
* **First-Class Functions:** Functions can be passed around as arguments, assigned to variables, and reassigned (if marked with `mut`).

---

## Getting Started

### Prerequisites
* [Go](https://go.dev/doc/install) (Built and tested on `1.25.6` — 1.18+ should work)

### Installation
Clone the repository and build the engine:

```bash
git clone [https://github.com/PuggyPotato/PotatoLang.git](https://github.com/PuggyPotato/PotatoLang.git)
cd potatolang
go build -o potatolang
```

### Running the REPL
To start the interactive prompt, just run the executable:

```bash
./potatolang
```

---

## Language Tour

### Variables & Mutability
By default, all variables are immutable. You must use `let mut` if you want to reassign them later.

```go
let a = 5;
a = 10; // ERROR: 'a' is not mutable

let mut b = 5;
b = 10; // Success!

// Type boundaries are strictly enforced during reassignment
b = true; // ERROR: TypeError: cannot assign 'boolean' to variable of type 'number'
```

### Safe Variable Unpacking
The engine enforces strict matching for variables and return values. You cannot accidentally leak internal memory boxes or unpack the wrong number of variables.

```rust
let getTwo = func() -> number, number { 
    return 1, 2; 
};

let a = getTwo();       // ERROR: assignment mismatch: 1 variables but 2 value

let b, c, d = getTwo(); // ERROR: assignment mismatch: 3 variables but 2 value

let x, y = getTwo();    // x = 1, y = 2
```

### Function Reassignment
Functions are first-class citizens. If marked as mutable, they can be dynamically swapped at runtime, provided the reassigned value is also a function.

```rust
let mut attack = func() -> number { return 10; };

// Safely swap the function behavior
attack = func() -> number { return 50; };

```

---

## Error Handling
Potatolang prioritizes developer experience by returning clean, readable runtime errors instead of causing Go engine panics.

Examples of caught runtime errors:
* `wrong number of arguments: got 0, want 1.`
* `TypeError: cannot assign string to variable 'a'.`
* `function declaration cannot be unpacked.`
* `assignment mismatch: 2 variables but 1 value.`

---

## Acknowledgments
This project was heavily inspired by Thorsten Ball's fantastic book [*Writing An Interpreter in Go*](https://interpreterbook.com/) and its companion language, Monkey. 

PotatoLang builds upon that incredible foundation by introducing new, strict language mechanics like immutable-by-default variables, explicit `mut` tracking, multiple return values, multiple variable binding and strict runtime type enforcement.

## Contributing
Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/PuggyPotato/PotatoLang/issues).

## License
This project is [MIT](LICENSE) licensed.