# Contributing to Diego

Thank you for contributing to Diego. Your time and help are appreciated. Here are some guidelines to help you get started.

To request support for a new service or file, submit a [feature request](https://github.com/ttybitnik/diego/issues/new?assignees=&labels=enhancement&projects=&template=feature_request.md&title=) with a small sample of the official exported file.

In case of unexpected behavior, please open a [bug report](https://github.com/ttybitnik/diego/issues/new?assignees=&labels=bug&projects=&template=bug_report.md&title=).

To contribute to `diego` development, refer to the sections below.

## Basic Steps

1. Fork the project.
1. Download your fork to your PC (`git clone https://github.com/your_username/diego && cd diego`)
1. Create your feature branch (`git checkout -b feat/new-feature`)
1. Make changes and test (`make run`)
1. Add them to staging (`git add .`)
1. Commit your changes (`git commit -m 'feat(import): add support for new-feature'`)
1. Push to the branch (`git push origin feat/new-feature`)
1. Create new pull request

## Testing and Running

Execute `make run` to test and check your changes. If you do not have `golangci-lint` installed on your system, comment out the `golangci-lint run` line in the `Makefile`.

## Diego Internals

Diego follows the **Port and Adapters Architecture (Hexagonal)**. Refer to commit [3eb8bf8](https://github.com/ttybitnik/diego/commit/3eb8bf8c4ff034c0383a258be3eda1b966aa1e86) for an overview of the files that need to be changed to implement support for a new service/file.

[Conventional Commits](https://www.conventionalcommits.org/) messages are welcome but not mandatory, since each pull request will be squashed during the merge process. They are used to automate [Semantic Versioning](https://semver.org/) for the releases.

### Makefile

- `make lint` :: lint the source code
- `make test` :: lint and test the source code
- `make build` :: lint, test, and build the binary
- `make run` :: lint, test, build, and run the binary
- `make deploy` :: lint, test, build, and deploy the application locally
- `make update` :: update module dependencies and call `make run`
- `make golden` :: generate/update golden files using current test results

### Generate Docs

To automate the process of updating the `docs/help` and `docs/man` files, set the `DIEGO_GENDOCS` environment variable to `1` before building and running the application.

The recommended approach is to temporarily set the variable by using the following command whenever necessary:
```bash
DIEGO_GENDOCS=1 make run
```
