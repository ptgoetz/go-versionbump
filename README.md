# go-versionbump

**Latest Version:** v0.1.0

VersionBump is a powerful command-line tool designed to streamline the process of version management in your projects. 
By automating version bumping, VersionBump ensures that your project’s version numbers are always up-to-date across all 
relevant files, reducing the risk of human error and saving you valuable time.

## Key Features

- **Automated Version Bumping**: Automatically updates version numbers in specified files, ensuring consistency and accuracy.
- **Git Integration**: Seamlessly integrates with Git to commit and tag changes, making version control effortless.
- **Dry Run Mode**: Preview changes without making any modifications, giving you full control over the process.
- **User Confirmation**: Prompts for user confirmation before making changes, with options to disable prompts for a fully automated experience.
- **Verbose Logging**: Detailed logging for debugging and verification, with options to enable or disable as needed.
- **Customizable Configuration**: Flexible configuration options to tailor VersionBump to your specific needs.


## Rationale
Any project that relies on version strings embedded in code and/or configuration files can get unwieldy pretty quickly
if you have to manually update those version strings. VersionBump is designed to automate this process so you can focus
on writing code instead of updating version numbers.

### Why Not [Insert Tool Here]?
There are many tools available that can automate version bumping, but VersionBump is designed to be simple, flexible,
and unobtrusive. It is a single binary with no external dependencies, and it is easy to configure and use. It is also
designed to be as safe as possible, with built-in safeguards to prevent accidental changes to your project.

With VersionBump you'll never have to switch between virtual environments, install dependencies, or worry about
compatibility issues. It is a simple, lightweight tool that gets the job done without any fuss, and will work with any
project that uses version strings in code or configuration files.

### Do No Harm
By default VersionBump will do its best to not not make any changes to your project unless you approve them. You will
be prompted to confirm the changes before they are made. You can also run VersionBump in `--dry-run` mode to see what
changes would be made without actually making them. To make VersionBump truly silent and prompt-less, you have use the
`--no-prompt` and `--silent` flags.

If anything goes wrong, VersionBump will not make any changes to your project and will exit with a non-zero error code.

## Installation
When installed with `go install`, it provides a `versionbump` binary that can be run from the command line. In the 
future we will also provide pre-built binaries for common platforms and CPU architectures.

```shell
go install github.com/ptgoetz/go-versionbump/cmd/versionbump
```

## Usage
Run VersionBump with the desired options:

```sh
# Bump the version
./versionbump [--config path/to/versionbump.yaml][--dry-run] [--no-prompt] [--quiet] bump-part

# Reset the version in all tracked files
./versionbump [--config path/to/versionbump.yaml] --reset version
```
- `bump-part`: The part of the version number to bump (`major`, `minor`, or `patch`).
- `--V`: Print the VersionBump version and exit.
- `--config`: Path to the configuration file (default: `./versionbump.yaml`).
- `--dry-run`: Perform a dry run without making any changes.
- `--no-prompt`: Do not prompt the user for confirmation before making changes.
- `--no-git`: Do not commit or tag the changes in a Git repository.
- `--quiet`: Disable verbose logging.
- `--reset [version]`: Reset the version number to the specified value.

## Configuration
The configuration file (**Default:** `versionbump.yaml`) defines the version bump settings:

```yaml
version: "0.0.0"         # (REQUIRED) The current version of the project.

# Git settings are optional. All default too `false`.
git-commit: false        # Whether to create a git commit for the version bump.
git-tag: false           # Whether to create a git tag for the version bump.

files:                   # The files to update with the new version.
   - path: "version.go"   # The path to the file to update.
     replace: "v{version}" # The search string to replace in the file.

   - path: "README.md"
     replace: "Latest Version: {version}"
```

- `version`: REQUIRED The current version of the project This must be a [Semantic Versioning](https://semver.org/) 
             `major.minor.patch` string.
- `git-commit`: (Optional) Whether to `git commit` the changes.
- `git-tag`: (Optional) Whether to tag the commit (implies `git-commit`).
- `files`: (Required) A list of files to update with the new version number.
   - `path`: The path to the file. **Note**: Relative file paths are relative to the config file parent directory. 
             Absolute paths are used as-is.
   - `replace`: The string to replace with the new version number. Use `{version}` as a placeholder.

**Important Note:**
The specified or default configuration file is implicitly included as a file that will undergo version replacement. It
serves as the source of truth for the version number. VersionBump will always include it as a file to update with the
new version number.

### Git Message Templates
VersionBump will use the following templates for the commit and tag messages. You can customize these templates in the
YAML configuration file.

| Template Type      | Default Value                  | Config YAML Key            |
|--------------------|--------------------------------|----------------------------|
| Commit Message     | `Bump version {old} --> {new}` | `git-commit-template`      |
| Tag Name           | `v{new}`                       | `git-tag-template`         |
| Tag Message        | `Release version {new}`        | `git-tag-message-template` |

The following placeholders can be used in the templates:
- `{old}`: The old semantic version number.
- `{new}`: The new semantic version number.

## Sample Ouput

### Configuration File
```yaml
version: "0.1.7"          # The current version of the project.
git-commit: true          # Whether to create a git commit for the version bump.
git-tag: true             # Whether to create a git tag for the version bump.

files:                    # The files to update with the new version.
  - path: "main.go"       # The path to the file to update.
    replace: "v{version}" # The search string to replace in the file.

  - path: "README.md"
    replace: "**Current Version:** v{version}"
```

### Default (Verbose) Output with Prompts
```text
$ versionbump patch
VersionBump v0.0.0
Config path: versionbump.yaml
Git version: 2.39.3 (Apple Git-146)
Project root: /Users/tgoetz/Projects/ptgoetz/test-project
Current Version: 0.1.6
Tracked Files:
  - main.go
  - README.md
  - versionbump.yaml
Bumping version part: patch
Will bump version 0.1.6 --> 0.1.7
main.go
     Find: "v0.1.6"
  Replace: "v0.1.7"
    Found 1 replacement(s)
README.md
     Find: "**Current Version:** v0.1.6"
  Replace: "**Current Version:** v0.1.7"
    Found 1 replacement(s)
versionbump.yaml
     Find: "version: "0.1.6""
  Replace: "version: "0.1.7""
    Found 1 replacement(s)
The following files will be updated:
  - main.go
  - README.md
  - versionbump.yaml
Do you want to proceed with the changes? [y/N]: y
Updated file: main.go
Updated file: README.md
Updated file: versionbump.yaml
Do you want to commit the changes to the git repository? [y/N]: y
Committing changes...
Commit message: 'Bump version 0.1.6 --> 0.1.7'
Tagging changes...
Tag 'v0.1.7' created with message: 'Release  0.1.7'
```
### Suppressing Prompts and Verbose Output
```shell
$ versionbump --no-prompt --quiet patch
# No output

$ echo $?
0 # Success

$ git log --name-status HEAD^..HEAD # Show last commit
commit c78b938150ab0e1c70178bf9e6c5a82f6c762830 (HEAD -> main, tag: v0.1.8)
Author: P. Taylor Goetz <ptgoetz@gmail.com>
Date:   Sat Aug 24 14:15:24 2024 -0400
    Bump version 0.1.7 --> 0.1.8

M       README.md
M       main.go
M       versionbump.yaml
```

## Development
If you want to hack and/or contribute to VersionBump, look at the [DEVELOPER.md](DEVELOPER.md) file for more 
information.
