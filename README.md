# labdoc

<div align="center">
  <img src="./docs/img/icon.png" width="250" alt="labdoc icon">
    <p>
        Automatically generate documentation for GitLab CI/CD components and CI/CD pipelines.
    </p>
</div>

## What is `labdoc`?

`labdoc` currently focuses on generating documentation for [GitLab CI/CD components](https://docs.gitlab.com/ee/ci/components/).
For an example, check out the [`labdoc` component documentation](./templates/README.md), which is generated with `labdoc`.

In the future, the focus might shift to generating documentation for GitLab CI/CD Pipelines in general.

## Getting Started

### Install `labdoc`

#### Via Homebrew

```shell
brew install erNail/tap/labdoc
```

#### Via Binary

Check the [releases](https://github.com/erNail/labdoc/releases) for the available binaries.
Download the correct binary and add it to your `$PATH`.

#### Via Go

```shell
go install github.com/erNail/labdoc
```

#### Via Container

```shell
docker pull ernail/labdoc:<LATEST_GITHUB_RELEASE_VERSION>
```

#### From Source

Check out this repository and run the following:

```shell
go build
```

Add the resulting binary to your `$PATH`.

### Run `labdoc`

#### Prepare your GitLab CI/CD components

`labdoc` currently expects all your CI/CD components to be in a directory, with all files on the root level.
By default, `labdoc` will use the [`templates` directory](https://docs.gitlab.com/ee/ci/components/#directory-structure).

The documentation is generated from the `spec.inputs.*.description` keywords,
and from the comments above the `spec` and the job keywords. Below is a minimal example:

```yaml
---
# This comment will be used as description for the component
spec:
  inputs:
    my-input:
      description: "This is used as description for the input"
    my-other-input:
      description: >-
        This is a multiline input.
        Since this output is used in a table, the `>-` is used to remove any newline characters
...

---
# This comment will be used as description for the job
my-job:
  script: "echo Hello"
...
```

#### Generate Documentation

```shell
labdoc generate --repoUrl github.com/erNail/labdoc
```

This will generate a `README.md` in the `templates` directory.

The `--repoUrl` flag is required to generate the instructions on how to include your components.

#### Change the documentation output directory

```shell
labdoc generate --repoUrl github.com/erNail/labdoc --outputFile my/custom/path/README.md
```

#### Change the component directory

```shell
labdoc generate --repoUrl github.com/erNail/labdoc --componentDir my-components
```

#### Check if the documentation is up-to-date

```shell
labdoc generate --repoUrl github.com/erNail/labdoc --check
```

This command will not write the documentation to a file.
It will only check if there is already a documentation, and if the content would change.

If the content remains unchanged, the command will exit with code 0.
If there is no documentation, or the existing documentation would change, the command will exit with code 2.

#### Include the version in the usage instructions

```shell
labdoc generate --repoUrl github.com/erNail/labdoc --version 1.0.0
```

By default, `labdoc` will generate instructions on how to include your component in other CI/CD pipelines.
If no version is specified, it will use `latest` as the version to use for the include.

#### Custom Documentation Template

By default, `labdoc` will generate documentation based on the
[documentation template](./internal/gitlab/resources/default-template.md.gotmpl) located in this repository.

You can create your own template and use it to generate documentation.
Simply create a file that uses [Go Templating](https://pkg.go.dev/text/template) syntax and the [type `ComponentDocumentation`](./internal/gitlab/component_documentation.go),
then run the following:

```shell
labdoc generate --repoUrl github.com/erNail/labdoc --template templates/README.md.gotmpl
```

#### More Details

For more details about the `labdoc` command, run the following:

```shell
labdoc -h
```

### `pre-commit` Hook

You can run `labdoc` via [`pre-commit`](https://pre-commit.com/).
Add the following to your `.pre-commit-config.yml`:

```yaml
repos:
  - repo: "https://github.com/erNail/labdoc"
    rev: "<LATEST_GITHUB_RELEASE_VERSION>"
    hooks:
      - id: "labdoc-generate"
        args:
          - "--repoUrl=gitlab.com/erNail/labdoc"
```

## Limitations

- `labdoc` currently only supports reading CI/CD component files from a directory with all files on the root level.
  These file names will be used as the component names.
- `labdoc` currently expects all components to define `spec:inputs` and at least one job.
  Not defining one or the other can lead to unwanted behavior.
- As a result of this, `labdoc` is currently not able to handle components that only include other components
  or GitLab CI/CD files.

## Planned Features

Please check the open [GitHub Issues](https://github.com/erNail/homebrew-tap/issues)
to get an overview of the planned features.

## Development

### Dependencies

To use all of the functionality listed below,
you need to install all dependencies listed in the [dev container setup script](.devcontainer/postCreateCommand.sh).
If you are using this repositories dev container, you already have all dependencies installed.

### Testing

```shell
task test
```

### Linting

```shell
task lint
```

### Running

```shell
task run -- --help
```

```shell
task run-generate
```

### Building

```shell
task build
```

### Building Container Images

```shell
task build-image
```

### Test GitHub Actions

```shell
task test-github-actions
```

### Test Release

```shell
task release-test
```
