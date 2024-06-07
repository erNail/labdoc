# Components Documentation

## Components

### labdoc-generate

A GitLab CI/CD Component for generating Markdown documentation from GitLab CI/CD Components
using `labdoc`.

#### Usage

You can add this component to an existing `.gitlab-ci.yml` file by using the `include:` keyword.

```yaml
include:
  - component: "gitlab.com/erNail/labdoc/labdoc-generate@latest"
    inputs: {}
```

You can configure the component with the inputs documented below.

#### Inputs

| Name | Description | Type | Default |
|------|-------------|------|---------|
| `additional-labdoc-parameters` | Additional parameters to add to the `labdoc generate` command. If you want this job to only check if your existing documentation is up-to-date, use the `--check` flag. | `string` | `` |
| `image` | The image to use for running `labdoc`. | `string` | `ernail/labdoc:1.0.0` |
| `labdoc-generate-job-extends` | The jobs that the job that generates the documentation should inherit from | `array` | `[]` |
| `labdoc-generate-job-name` | The name of the job that generates the documentation. | `string` | `labdoc-generate-job` |
| `output-file-path` | The path and name of the rendered file to be created. | `string` | `templates/README.md` |
| `repo-url` | The repository URL from which to include the GitLab CI/CD Component. Will be used in the documentation. | `string` | `<no value>` |
| `stage` | The stage of the jobs for generating the documentation. | `string` | `docs` |

#### Jobs

The component will add the following jobs to your CI/CD Pipeline

##### `$[[ inputs.labdoc-generate-job-name ]]`

Generates Markdown documentation from GitLab CI/CD Components.
The generated documentation will be uploaded as an artifact at `$[[ inputs.output-file-path ]]`.
