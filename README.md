# Mini‑Go Projects

A collection of small projects and code samples written in Go, created for learning and practice.

---

## Project List
> (Replace with your actual project folders)

- `hello-world/` – the classic Go "Hello World".
- `calculator/` – a simple command-line calculator.
- `todo-api/` – a tiny RESTful API for managing a todo list.
- `web-scraper/` – a minimal web scraper for data extraction.

---

## Getting Started

### Requirements
- Go (version 1.18 or later)
- (Optional) Database like PostgreSQL if used in some projects

### Run
1. Clone the repository:
    ```bash
    git clone https://github.com/amingeek/mini-go.git
    ```

2. Navigate into any project and run it:
    ```bash
    cd mini-go/hello-world
    go run main.go
    ```

3. For projects with modules and dependencies:
    ```bash
    go mod tidy
    go run main.go
    ```

---

## Features

- Each project is **single-purpose and small**.
- Uses `go.mod` for dependency management.
- Covers core concepts such as:
  - HTTP routing & REST APIs
  - CLI input/output
  - Concurrency with goroutines and channels

---

## Suggested Project Structure
```
mini-go/
├── hello-world/
│   └── main.go
├── calculator/
│   └── calculator.go
├── todo-api/
│   ├── main.go
│   ├── handlers.go
│   └── go.mod
└── README.md   ← this file
```

---

## Contributing

If you have a fun or useful mini-project in Go that fits here, contributions are welcome!

1. Open an issue describing your idea.
2. Fork this repo and create a Pull Request.

---

## Quick Overview

| Feature | Description |
|---------|-------------|
| Goal | Quick learning of Go basics |
| Run | `go run main.go` (simple projects) or `go mod tidy && go run main.go` (with dependencies) |
| Contributions | PRs with useful mini-projects are welcome |
