# SecureJS

SecureJS is a powerful tool designed to collect all related links from a target website, perform requests on these links (primarily JavaScript files), and scan for sensitive information such as tokens, keys, passwords, AKSKs, and more.

## Table of Contents

- [SecureJS](#securejs)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Usage](#usage)
    - [Example](#example)
  - [Configuration](#configuration)
  - [Project Structure](#project-structure)

## Features

- **Comprehensive Crawling**: Simulates browser visits to collect all links and JavaScript files from the target.
- **Secondary Requests**: Performs additional requests on collected resources for deeper analysis.
- **Customizable Matching Rules**: Supports custom rules defined in `config.yaml` to identify sensitive information.
- **Flexible Output Formats**: Outputs results in CSV, JSON, or plain text formats.
- **Easy Configuration**: Simplifies setup and customization through a configuration file.

## Usage

SecureJS can be executed via the command line with various options to customize its behavior.

### Example

```bash
./SecureJS -u https://example.com -o results.csv
```

```bash
./SecureJS -l targets.txt -o results.csv -t 30
```
## Configuration

SecureJS uses a `config/config.yaml` file to define custom matching rules and other project-level configurations.

## Project Structure

```
SecureJS/
├── cmd/
│   └── root.go             # Entry point for command-line arguments handling (-u, -l, -t, etc.)
│
├── internal/
│   ├── crawler/
│   │   ├── crawler.go      # Crawler logic, simulates browser access, collects all links and JS files
│   │   └── linkfind.go     # Extracts all links and JS from the response body of the target page
│   │
│   ├── parser/
│   │   └── parser.go       # Performs secondary requests on all collected links and JS files
│   │
│   ├── matcher/
│   │   └── matcher.go      # Reads and parses custom rules from config.yaml and matches against response bodies
│   │
│   └── output/
│       └── output.go       # Outputs results to files in CSV, JSON, or text formats
│
├── config/
│   ├── config.go           # Handles loading and parsing of the configuration file (config.yaml)
│   └── config.yaml         # Custom rules and other project-level configurations
│
├── go.mod                  # Go Modules management file
├── go.sum                  # Go Modules checksum file
└── main.go                 # Main program entry point, initializes and starts the application
```