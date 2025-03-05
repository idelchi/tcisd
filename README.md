# the-code-is-self-documented (tcisd)

[![Go Reference](https://pkg.go.dev/badge/github.com/idelchi/tcisd.svg)](https://pkg.go.dev/github.com/idelchi/tcisd)
[![Go Report Card](https://goreportcard.com/badge/github.com/idelchi/tcisd)](https://goreportcard.com/report/github.com/idelchi/tcisd)
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

`tcisd` solves this problem by ruthlessly purging your codebase of these wasteful explanations, preserving the elegant obscurity and job security that true programming artisans strive for.

It supports:

- Detecting and removing comments from Go, Python, and Bash files
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

- `lint`: For cowards who just want to know where the comments are without removing them
- `format`: For true believers who want to immediately purify their code by removing all comments

### Global Flags

| Flag            | Description                     |
| --------------- | ------------------------------- |
| `-s, --show`    | Show the configuration and exit |
| `-h, --help`    | Help for tcisd                  |
| `-v, --version` | Version for tcisd               |

### Command Flags

| Flag            | Description                                        | Default   |
| --------------- | -------------------------------------------------- | --------- |
| `-p, --pattern` | File pattern to match (doublestar format)          | `**/*.go` |
| `-t, --type`    | File types to process (go, bash, python)           | `go`      |
| `-e, --exclude` | Patterns to exclude                                | -         |
| `-a, --hidden`  | Include hidden files and directories               | `false`   |
| `-d, --dry-run` | Show what would be changed without modifying files | `false`   |
| `-h, --help`    | Help for the command                               | -         |

## Examples

```sh
# Identify all the comments lurking in your Go code
tcisd lint --pattern="**/*.go" .

# Purge Python comments from a specific directory
tcisd format --type="python" --pattern="**/*.py" src/

# Check for Bash comments while excluding test scripts
tcisd lint --type="bash" --pattern="**/*.sh" --exclude="**/test/*.sh" .

# Perform a dry run to see what comments would be removed
tcisd format --dry-run --pattern="**/*.go" .

# Process all supported file types in a project
tcisd format --type="go" --type="python" --type="bash" .
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
