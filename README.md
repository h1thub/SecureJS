# SecureJS

SecureJS is a powerful tool designed to collect all related links from a target website, perform requests on these links (primarily JavaScript files), and scan for sensitive information such as tokens, keys, passwords, AKSKs, and more.

## Table of Contents

- [SecureJS](#securejs)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Steps](#steps)
  - [Usage](#usage)
    - [Example](#example)
  - [Configuration](#configuration)
    - [Sample `config.yaml`](#sample-configyaml)
    - [Loading Configuration](#loading-configuration)
  - [Project Structure](#project-structure)

## Features

- **Comprehensive Crawling**: Simulates browser visits to collect all links and JavaScript files from the target.
- **Secondary Requests**: Performs additional requests on collected resources for deeper analysis.
- **Customizable Matching Rules**: Supports custom rules defined in `config.yaml` to identify sensitive information.
- **Flexible Output Formats**: Outputs results in CSV, JSON, or plain text formats.
- **Easy Configuration**: Simplifies setup and customization through a configuration file.

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) 1.16 or later

### Steps

1. **Clone the Repository**

   ```bash
   git clone
   cd SecureJS
   ```

2. **Build the Application**

   ```bash
   go build
   ```

3. **Verify Installation**

   ```bash
   ./SecureJS -h
   ```

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

SecureJS uses a `config.yaml` file to define custom matching rules and other project-level configurations.

### Sample `config.yaml`

```yaml
rules:
  - name: Sensitive Field
    f_regex: (?i)\[?["']?[0-9A-Za-z_-]{0,15}(?:key|secret|token|config|auth|access|admin|ticket)[0-9A-Za-z_-]{0,15}["']?\]?\s*(?:=|:|\)\.val\()\s*\[?\{?["']([^"']{8,256})["']?(?::|,)?

  - name: Password Field
    f_regex: ((|\\)(|'|")(|[\w]{1,10})([p](ass|wd|asswd|assword))(|[\w]{1,10})(|\\)(|'|")(:|=|\)\.val\()(|)(|\\)('|")([^'"]+?)(|\\)('|")(|,|\)))

  - name: JSON Web Token
    f_regex: (eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9._-]{10,}|eyJ[A-Za-z0-9_\/+-]{10,}\.[A-Za-z0-9._\/+-]{10,})

  - name: Cloud Key
    f_regex: (?i)(?:AWSAccessKeyId=[A-Z0-9]{16,32}|access[-_]?key[-_]?(?:id|secret)|LTAI[a-z0-9]{12,20})
```

### Loading Configuration

The configuration is automatically loaded from the `config/config.yaml` file. Ensure that your custom rules are correctly defined to match the sensitive information you aim to identify.

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