# go-VersionBump

![VersionBump Gopher](assets//versionbump_gopher-250.png)

**Latest Version:** v0.3.0 ([Download Binary](https://github.com/ptgoetz/go-versionbump/releases/tag/v0.3.0))

VersionBump is a powerful command-line tool designed to streamline the process of version management in your projects. 
By automating version bumping, VersionBump ensures that your project’s version numbers are always up-to-date across all 
relevant files, reducing the risk of human error and saving you valuable time.

## Key Features

- **Automated Version Bumping**: Automatically updates version numbers in specified files, ensuring consistency and 
  accuracy.
- **Git Integration**: Seamlessly integrates with Git to commit and tag changes, making version control effortless.
- **GPG Integration**: Supports GPG signing of git commits and tags for enhanced security and authenticity.
- **Interactive Mode**: Prompts for user confirmation before making changes, with options to disable prompts for a fully
  automated experience.
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

#### Existing Projects with Similar Functionality

The following Python projects drove and inspired the development of VersionBump. 
- **[bumpversion](https://github.com/peritus/bumpversion)**: No longer maintained. Requires Python.
- **[bump2version](https://github.com/c4urself/bump2version)**: No longer maintained. Requires Python.
- **[bump-my-version](https://github.com/callowayproject/bump-my-version)**: This is the closest to VersionBump in terms
  of intended functionality. Requires Python.

**Why Not Use Python Tools?**

The problem at hand is essentially a text-based search and replace operation, with some extra external tool calls for git 
integration. Dealing with Python dependencies, virtual environments, and compatibility issues is overkill for this
problem, especially when a tool's dependencies require switching between virtual environments so as not to conflict with
your prject's dependencies.

### Do No Harm
By default VersionBump will do its best to not not make any changes to your project unless you approve them. You will
be prompted to confirm the changes before they are made. By default VersionBump will run in "interactive" and will 
prompt you to approve all changes and extensively log what actions it's performing.  To make VersionBump truly silent 
and prompt-less, you have use the `--no-prompt` and `--silent` flags.

If anything goes wrong, VersionBump will not make any changes to your project and will exit with a non-zero error code.

## Installation

### With Go
When installed with `go install`, it provides a `versionbump` binary that can be run from the command line.

```shell
go install github.com/ptgoetz/go-versionbump/cmd/versionbump
```

### Without Go
If you don't have Go installed and just want the binary executable, you can download a prebuilt binaries from  
[here](https://github.com/ptgoetz/go-versionbump/releases/tag/v0.3.0).

VersionBump binary distribution archives include the `README.md` and `versionbump[.exe]` files:

```console
$ unzip -l versionbump-v0.3.0-darwin-arm64.zip 
Archive:  versionbump-v0.3.0-darwin-arm64.zip
  Length      Date    Time    Name
---------  ---------- -----   ----
  6082130  09-13-2024 21:23   versionbump
    16573  09-13-2024 21:22   README.md
---------                     -------
  6098703                     2 files

$ unzip -l versionbump-v0.3.0-windows-arm64.zip 
Archive:  versionbump-v0.3.0-windows-arm64.zip
  Length      Date    Time    Name
---------  ---------- -----   ----
  6112256  09-13-2024 21:23   versionbump.exe
    16573  09-13-2024 21:22   README.md
---------                     -------
  6128829                     2 files
 
$ tar -ztvf versionbump-v0.3.0-linux-arm64.tgz
-rwxr-xr-x  0 tgoetz staff 6103765 Sep 13 21:23 versionbump
-rw-r--r--  0 tgoetz staff   16573 Sep 13 21:23 README.md


```

When installing from a binary archive, you should place the `versionbump[.exe]` binary file in a directory in your 
system path.

## Usage
Run VersionBump without any arguments to see the available flags and commands:

```console
$ versionbump
                                              
VersionBump is a command-line tool designed to automate version string management in projects.

Usage:
  versionbump [flags]
  versionbump [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Show the effective configuration of the project.
  help        Help about any command
  major       Bump the major version number (e.g. 1.2.3 -> 2.0.0).
  minor       Bump the minor version number (e.g. 1.2.3 -> 1.3.0).
  patch       Bump the patch version number (e.g. 1.2.3 -> 1.2.4).
  reset       Reset the project version to the specified value.
  show        Show potential versioning paths for the project version or a specific version.

Flags:
  -h, --help      help for versionbump
  -V, --version   Show the version of Config and exit.

Use "versionbump [command] --help" for more information about a command.
```

The commands `major`, `minor` `patch` and `reset` support the following flags:'
- `-c`, `-config`: Path to the configuration file (default: `./versionbump.yaml`).
- `-no-prompt`: Do not prompt the user for confirmation before making changes.
- `-no-git`: Do not commit or tag the changes in a Git repository.
- `-no-color`: Disable colorized output.
- `-q`, `-quiet`: Disable verbose logging.

The commands `config` and `show` support the following flags:
- `-c`, `-config`: Path to the configuration file (default: `./versionbump.yaml`).
- `-no-color`: Disable colorized output.

## Configuration
The configuration file (**Default:** `versionbump.yaml`) defines the version bump settings:

```yaml
version: "0.0.0"         # (REQUIRED) The current version of the project.

# Git settings are optional. All default to `false`.
git-commit: false        # Whether to create a git commit for the version bump.
git-tag: false           # Whether to create a git tag for the version bump.
git-sign: false          # Whether to sign the git commit/tag.

files:                   # The files to update with the new version.
   - path: "version.go"   # The path to the file to update.
     replace: 
      - "v{version}" # The search string to replace in the file.

   - path: "README.md"
     replace: 
      - "Latest Version: {version}"
```

- `version`: REQUIRED The current version of the project This must be a [Semantic Versioning](https://semver.org/) 
             `major.minor.patch` string.
- `git-commit`: (Optional) Whether to `git commit` the changes.
- `git-tag`: (Optional) Whether to tag the commit (implies `git-commit`).
- `git-sign`: (Optional) Whether to sign the commit/tag with GPG.
- `files`: (Required) A list of files to update with the new version number.
   - `path`: The path to the file. **Note**: Relative file paths are relative to the config file parent directory. 
             Absolute paths are used as-is.
   - `replace`: A list of strings to replace with the new version number. Use `{version}` as a placeholder.

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

## Examples

### Configuration File
```yaml
version: "0.1.9"          # The current version of the project.
git-commit: true          # Whether to create a git commit for the version bump.
git-tag: true             # Whether to create a git tag for the version bump.
git-sign: true            # Whether to sign the git commit/tag.

files:                    # The files to update with the new version (i.e. "Tracked files").
  - path: "main.go"       # The path to the file to update.
    replace: 
      - "v{version}" # The search string to replace in the file.

  - path: "README.md"
    replace: 
      - "**Current Version:** v{version}"
```

### Default (Verbose) Output with Prompts
In the following scenario, the project is not a git repository but git features are enabled, so VersionBump will 
offer to initialize a git repository in the project directory. VersionBump will add tracked files to the git repository
and perform an initial commit before continuing.

```console
$ versionbump patch
VersionBump v0.3.0
Configuration file: versionbump.yaml
Project root directory: /Users/tgoetz/Projects/ptgoetz/test-project
Checking git configuration...
Git version: 2.39.3 (Apple Git-146)
The project directory is not a git repository.
Do you want to initialize a git repository in the project directory? [y/N]: y
Initialized Git repository.
Adding tracked files...
Tracked Files:
  - main.go
  - README.md
  - versionbump.yaml
Performing initial commit.
Current branch: main
Checking for existing tag...
GPG signing of git commits is enabled. Checking configuration...
Git commits will be signed with GPG key: ACEFE18DD2322E1E84587A148DE03962E80B8FFD
Tracked Files:
  - main.go
  - README.md
  - versionbump.yaml
Bumping version part: patch
Will bump version 0.1.10 --> 0.1.11
main.go
     Find: "v0.1.10"
  Replace: "v0.1.11"
    Found 1 replacement(s)
README.md
     Find: "**Current Version:** v0.1.10"
  Replace: "**Current Version:** v0.1.11"
    Found 1 replacement(s)
versionbump.yaml
     Find: "version: "0.1.10""
  Replace: "version: "0.1.11""
    Found 1 replacement(s)
Proceed with the changes? [y/N]: y
Updated file: main.go
Updated file: README.md
Updated file: versionbump.yaml
Commit Message: Bump version 0.1.10 --> 0.1.11
Tag Message: Release version 0.1.11
Tag Name: v0.1.11
Do you want to commit the changes to the git repository? [y/N]: y
Committing changes...
Committed changes with message: Bump version 0.1.10 --> 0.1.11
Tagging changes...
Tag 'v0.1.11' created with message: Release version 0.1.11

```
### Suppressing Prompts and Verbose Output
```console
$ versionbump -no-prompt -quiet patch
# No output

$ echo $?
0 # Success

$ git log --show-signature --name-status HEAD^..HEAD # Show last commit
commit e695bb7aaa8d4f7b6c821eb13d15fe4c658a929f (HEAD -> main, tag: v0.1.12)
gpg: Signature made Fri Sep 13 19:08:11 2024 EDT
gpg:                using RSA key ACEFE18DD2322E1E84587A148DE03962E80B8FFD
gpg: Good signature from "P. Taylor Goetz <ptgoetz@apache.org>" [ultimate]
gpg:                 aka "P. Taylor Goetz <ptgoetz@gmail.com>" [ultimate]
Author: P. Taylor Goetz <ptgoetz@gmail.com>
Date:   Fri Sep 13 19:08:11 2024 -0400

    Bump version 0.1.11 --> 0.1.12

M       README.md
M       main.go
M       versionbump.yaml

```

### Show Command
Without parameters, the `show` command will display the potential versioning paths for the project version:
```console
$ versionbump show
Potential versioning paths for project version: 0.1.7
0.1.7 ── bump ─┬─ major ─ 1.0.0
               ├─ minor ─ 0.2.0
               ╰─ patch ─ 0.1.8
```

You can also specify any version identifier to see the potential versioning paths:
```console
versionbump show 1.2.3 
Potential versioning paths for version: 1.2.3
1.2.3 ── bump ─┬─ major ─ 2.0.0
               ├─ minor ─ 1.3.0
               ╰─ patch ─ 1.2.4
```

### Config Command

The `config` command will display the effective configuration of the project. This will show default values for any
configuration settings that are not explicitly set in the configuration file.

```console
$ versionbump config
Config file: versionbump.yaml
Project root: /Users/tgoetz/Projects/ptgoetz/test-project
Effective Configuration YAML:
version: 0.1.7
git-commit: true
git-commit-template: Bump version {old} --> {new}
git-sign: true
git-tag: true
git-tag-template: v{new}
git-tag-message-template: Release version {new}
files:
    - path: main.go
      replace: v{version}
    - path: README.md
      replace: '**Current Version:** v{version}'
    - path: versionbump.yaml
      replace: 'version: "{version}"'
```

## Failure Modes and Errors
VersionBump does its best to prevent leaving your project in an inconsistent state. Before making any changes, it will
perform a series of "pre-flight" checks to ensure that the version bump can be completed successfully. If any errors are
detected, VersionBump will exit with a non-zero error code and will not make any changes to your project.

If VersionBump is run in `--no-prompt` mode, it will exit with an error if any of the pre-flight checks fail. If it is
run in interactive mode (default), it will prompt the user to confirm whether to proceed with the version bump.

If git integration is enabled in the VersionBump configuration, VersionBump will also exit with an error if it detects 
that any git operations (e.g., committing or tagging) will fail (e.g. the project directory is not a git repository).
When running in interactive mode, VersionBump will prompt the user to correct git issues it can fix (e.g. initializing
a git repository).


### Standard Pre-Flight Checks

- **Configuration File**: VersionBump will check that the configuration file exists and is read/write. If the file is
  missing or cannot be read or written, VersionBump will exit with an error.
- **Version Number**: VersionBump will check that the version number in the configuration file is a valid semantic
  version number. If the version number is invalid, VersionBump will exit with an error. Note that VersionBump will 
  normalize the version strings to a semantic version number before proceeding. For example the string `"1.2.003"` will
  be normalized to `1.2.3`.
- **Tracked Files**: VersionBump will check that all tracked files in the configuration file exist and are read/write. 
  If any files are missing or cannot be read, VersionBump will exit with an error.
- **At Least One Replacement**: VersionBump will check that at least one replacement will be made in each tracked file
  replacement. 
  If no replacements would be made, VersionBump will exit with an error.

### Git Pre-Flight Checks

- **Git Installed**: VersionBump will check that the `git` command is available in the system path. If the `git` command
  is not available, VersionBump will exit with an error.
- **Git Repository**: If git integration is enabled, VersionBump will check that the project directory is a git
  repository. If the project directory is not a git repository, VersionBump will exit with an error.

  In interactive mode, VersionBump will prompt the user to initialize a git repository in the project directory. It will
  also add all tracked files to the git repository and commit them with the message "Initial commit".
- **Git Clean**: VersionBump will check that the git repository is clean (i.e., no uncommitted changes). If the git
  repository is not clean, VersionBump will exit with an error.
- **Git Tagging**: If git tagging is enabled, VersionBump will check that the tag name does not already exist in the git
  repository. If the tag name already exists, VersionBump will exit with an error.

### GPG Pre-Flight Checks

If signing of git commits and tags is enabled, either in the VersionBump or git configuration, VersionBump will perform 
the following additional checks:

- **GPG Signing Key**: VersionBump will check that a GPG signing key is available in the git configuration 
  (`git config --get user.signingkey`). If no GPG signing key is available, VersionBump will exit with an error.
- **Sign/Don't Sign Conflict**: If signing is disabled in the VersionBump configuration, but enabled in the git 
  configuration, VersionBump will log a warning message and continue. VersionBump will not override the git 
  configuration for signing.

## Contributing
If you want to hack and/or contribute to VersionBump, look at the [DEVELOPER.md](DEVELOPER.md) file for more 
information.
