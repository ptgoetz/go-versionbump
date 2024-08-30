# Development

## Project Structure

- `cmd/`: Application entry points (main packages)
- `pkg/`: Public libraries (importable by other projects)
- `internal/`: Internal application code (not exposed publicly)
- `config/`: Configuration files or modules
- `test/integration/`: Integration tests
- `test/unit/`: Unit tests

## Getting Started

1. **Initialize the Project:**
    - Run `make init-project` to install all necessary Go tools and tidy up the dependencies.

2. **Install Dependencies:**
    - Run `make deps` to install and update the project dependencies.

3. **Format Code:**
    - Use `make format` to automatically format all Go source files according to Go standards.

4. **Lint Code:**
    - Use `make lint` to check the code for any style issues or potential bugs. This requires `golangci-lint` to be installed, which is handled by `make init-project`.

5. **Run the Application:**
    - Use `make run` to build and run the application.

6. **Run Tests:**
    - Use `make test` to run the unit tests.
    - Use `make test-integration` to run the integration tests (requires integration tests to be set up).
    - Use `make test-all` to run all tests (unit and integration).

7. **Build the Application:**
    - Use `make build` to compile the application into a binary for your current OS and architecture.

8. **Create a Binary Distribution:**
    - Use `make dist` to create a `.tgz` and `.zip` archive of the application binary for your current OS and architecture. These archives will be placed in the `dist/` directory.

9. **Create Distributions for Multiple Platforms:**
    - Use `make dist-all` to create `.tgz` and `.zip` archives for the most common OS and architecture combinations (Linux, macOS, and Windows on amd64 and arm64 architectures). These archives will be placed in the `dist/` directory.

10. **Clean Up:**
    - Use `make clean` to remove previous build artifacts and test cache.

11. **Tidy Dependencies:**
    - Use `make tidy` to clean up the `go.mod` and `go.sum` files by removing unused dependencies.

## Eating Our Own Dog Food
The VersionBump project uses itself to manage its version strings. The configuration file `versionbump.yaml` contains 
the current version number and the files that need to be updated with the new version. The `version.go` file contains 
the version number as a constant, and the `README.md` file contains the version number in the "Latest Version" section.