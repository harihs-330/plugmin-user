## Installation

### Ubuntu

To install `golangci-lint` using `snap` on the latest version of Ubuntu, follow these steps:

1. **Update your package list**:
    ```bash
    sudo apt update
    ```

2. **Install `golangci-lint` using `snap`**:
    ```bash
    sudo snap install golangci-lint --classic
    ```

3. **Verify the installation**:
    ```bash
    golangci-lint --version
    ```

This will install `golangci-lint` via `snap`, ensuring you have the latest version and that it's properly integrated into your system.


## How to Run ?

To run `golangci-lint` execute:
```bash
golangci-lint --version
```
    

It's an equivalent of executing:

```bash
golangci-lint run ./...
```

You can choose which directories or files to analyze:

```bash
golangci-lint run dir1 dir2/...
golangci-lint run file1.go
```

Pass `-E/--enable` to enable linter and `-D/--disable` to disable:

```bash
golangci-lint run --disable-all -E errcheck
```

## Configuration

GolangCI-Lint looks for config files in the following paths from the current working directory:

- .golangci.yml
- .golangci.yaml
- .golangci.toml
- .golangci.json

**note**: `GolangCI-Lint` also searches for config files in all directories from the directory of the first analyzed path up to the root. If no configuration file has been found, `GolangCI-Lint` will try to find one in your home directory. To see which config file is being used and where it was sourced from `run golangci-lint with -v` option.

The configuration file can be validated with the JSON Schema: https://golangci-lint.run/jsonschema/golangci.jsonschema.json



```yml
# Options for analysis running.
run:
  # See the dedicated "run" documentation section.
  option: value
# output configuration options
output:
  # See the dedicated "output" documentation section.
  option: value
# All available settings of specific linters.
linters-settings:
  # See the dedicated "linters-settings" documentation section.
  option: value
linters:
  # See the dedicated "linters" documentation section.
  option: value
issues:
  # See the dedicated "issues" documentation section.
  option: value
severity:
  # See the dedicated "severity" documentation section.
  option: value
```


If you don't need to check the linting for a specific file, add the following at the top of the file:

```go
// nolint
```

To know more about available options visit https://golangci-lint.run/usage/configuration/