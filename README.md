# the-code-is-self-documented (tcisd)

[![Go Reference](https://pkg.go.dev/badge/github.com/idelchi/tcisd.svg)](https://pkg.go.dev/github.com/idelchi/tcisd)
[![Go Report Card](https://goreportcard.com/badge/github.com/idelchi/tcisd)](https://goreportcard.com/report/github.com/idelchi/tcisd)
[![Build Status](https://github.com/idelchi/tcisd/actions/workflows/github-actions.yml/badge.svg)](https://github.com/idelchi/tcisd/actions/workflows/github-actions.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`tcisd` is a tool that liberates your code from the burden of documentation. Because obviously, code speaks for itself!

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Command Line Flags](#command-line-flags)
- [Default Exclusion Patterns](#default-exclusion-patterns)
- [Disclaimer](#disclaimer)

## Overview

Let's be honest: Comments are just lies waiting to happen. Code that requires explanation through comments is simply code that should be rewritten. The real 10x developers amongst us know that variable names like `x`, `temp`, and `data` are self-explanatory, and anyone who can't understand your undocumented 500-line functions simply isn't trying hard enough.

`tcisd` solves this problem by purging your codebase of these wasteful explanations, preserving the elegant obscurity and job security that true programming artisans strive for.

It supports:

- Detecting and removing comments from Go, Python, and Dockerfiles
- Recursive searching through your project
- Parallel processing for maximum efficiency

## Installation

### From source

```sh
go install github.com/idelchi/tcisd@latest
```

### From installation script

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/tcisd/refs/heads/dev/install.sh | sh -s -- -d ~/.local/bin
```

## Usage

```sh
tcisd [flags] command [flags] [path ...]
```

### Commands

- `lint`: For developers who want to identify the bad practice of over-commenting their code
- `format`: For those ready to streamline their code by removing unnecessary explanations

### Global Flags

| Flag            | Description                     |
| --------------- | ------------------------------- |
| `-s, --show`    | Show the configuration and exit |
| `-h, --help`    | Help for tcisd                  |
| `-v, --version` | Version for tcisd               |

### Command Flags

| Flag             | Description                                    | Default             |
| ---------------- | ---------------------------------------------- | ------------------- |
| `-p, --pattern`  | File pattern to match (doublestar format)      | Based on file types |
| `-t, --type`     | File types to process (go, python, dockerfile) | All supported types |
| `-e, --exclude`  | Patterns to exclude                            | -                   |
| `-a, --hidden`   | Include hidden files and directories           | `false`             |
| `-j, --parallel` | Number of parallel workers to use              | Number of CPUs      |
| `-h, --help`     | Help for the command                           | -                   |

## Default Behaviors

- If no file types are specified, all supported types (go, python, dockerfile) are used
- If no patterns are specified, patterns are automatically generated based on the selected file types:
  - Go: `**/*.go`
  - Python: `**/*.py`
  - Dockerfile: `**/Dockerfile` and `**/Dockerfile.*`
- At least one path must be provided for both commands

## Examples

```sh
# Identify all the comments lurking in your Go code
tcisd lint --type="go" .

# Purge Python comments from a specific directory
tcisd format --type="python" src/

# Process only Dockerfiles in a project
tcisd format --type="dockerfile" .

# Process specific file types with custom patterns
tcisd lint --type="go" --pattern="**/cmd/*.go" --pattern="**/internal/*.go" .
```

## Default Exclusion Patterns

To avoid disasters, tcisd automatically excludes several patterns:

- `**/*.exe`
- `**/.git/**`
- `**/node_modules/**`
- `**/vendor/**`
- `**/.task/**`
- `**/.cache/**`
- Hidden files and directories (unless explicitly included with `-a`)
- Binary files
- The executable itself

## Custom Comment Removers

Advanced users looking to spread the joy of uncommented code to other languages can easily extend `tcisd` by implementing the `Remover` interface and registering it for their language:

```go
type Remover interface {
    // Process removes comments from the given lines of code.
    // It returns the processed lines and a list of issues found.
    Process(lines []string) ([]string, []string)
}

// Register a custom remover
remover.Register("ruby", &RubyRemover{})
```

## Philosophy

Remember the sacred principles:

1. If you need comments to explain your code, your code is wrong
2. Documentation is for people who don't understand their own code
3. True intellect is measured by how obscure your variable names are
4. If you can't understand the code without comments, you don't deserve to understand it with them
5. Future you and your colleagues will thank you for making them really think about the code

## Disclaimer

> **Warning**
> Use of this tool implies that you're the type of developer who values mystery and job security over maintainability. Users are advised against using this tool on production code bases where other humans might need to understand what's happening. Side effects may include increased job security, confused colleagues, and the silent tears of future maintainers.
