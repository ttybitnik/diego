# Contributing to Diego

Thank you for contributing to Diego. Your time and help are appreciated. Here are some guidelines to help you get started.

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

Diego follows the **Port and Adapters Architecture (Hexagonal)**. Refer to commit [0e48f4d](#placeholder) for an overview of the files that need to be changed to implement support for a new service/file.

[Conventional Commits](https://www.conventionalcommits.org/) messages are welcome but not mandatory, since each pull request will be squashed during the merge process. They are used to automate [Semantic Versioning](https://semver.org/) for the releases.
