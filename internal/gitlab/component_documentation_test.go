package gitlab

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateDocumentationRendersInputTableCorrectly(t *testing.T) {
	t.Parallel()

	componentContent := `---
# Component Description
spec:
  inputs:
    string-with-default:
      type: "string"
      default: ""
    string-without-default:
      type: "string"
    array-with-default:
      type: "array"
      default: []
    array-without-default:
      type: "array"
    boolean-with-default:
      type: "boolean"
      default: false
    boolean-without-default:
      type: "boolean"
    number-with-default:
      type: "number"
      default: 0
    number-without-default:
      type: "number"
    string-with-options:
      type: "string"
      options:
        - "one"
        - "two"
    string-with-regex:
      type: "string"
      regex: "^test."
    input-with-description-only:
      description: "Input with description only"
    input-with-default-only:
      default: []
    input-without-anything: {}
...
`

	expectedMarkdownTable := "\n" +
		"| Name | Description | Type | Default | Options | Regex | Mandatory |\n" +
		"|------|-------------|------|---------|---------|-------|-----------|\n" +
		"| `array-with-default` |  | `array` | `[]` | `-` | `-` | No |\n" +
		"| `array-without-default` |  | `array` | `-` | `-` | `-` | Yes |\n" +
		"| `boolean-with-default` |  | `boolean` | `false` | `-` | `-` | No |\n" +
		"| `boolean-without-default` |  | `boolean` | `-` | `-` | `-` | Yes |\n" +
		"| `input-with-default-only` |  | `-` | `[]` | `-` | `-` | No |\n" +
		"| `input-with-description-only` | Input with description only | `-` | `-` | `-` | `-` | Yes |\n" +
		"| `input-without-anything` |  | `-` | `-` | `-` | `-` | Yes |\n" +
		"| `number-with-default` |  | `number` | `0` | `-` | `-` | No |\n" +
		"| `number-without-default` |  | `number` | `-` | `-` | `-` | Yes |\n" +
		"| `string-with-default` |  | `string` | `\"\"` | `-` | `-` | No |\n" +
		"| `string-with-options` |  | `string` | `-` | `[one two]` | `-` | Yes |\n" +
		"| `string-with-regex` |  | `string` | `-` | `-` | `^test.` | Yes |\n" +
		"| `string-without-default` |  | `string` | `-` | `-` | `-` | Yes |\n"

	filesystem := afero.NewMemMapFs()
	componentDir := "templates"
	componentPath := componentDir + "/first-component.yml"
	outputFilePath := "README.md"

	err := filesystem.Mkdir(componentDir, 0o644)
	require.NoError(t, err)
	err = afero.WriteFile(filesystem, componentPath, []byte(componentContent), 0o644)
	require.NoError(t, err)

	documentationGenerator := &RealDocumentationGenerator{}
	documentationGenerator.GenerateDocumentation(
		filesystem,
		"templates",
		"resources/default-template.md.gotmpl",
		"github.com/test",
		"1.0.0",
		"README.md",
		false,
	)

	outputExists, err := afero.Exists(filesystem, outputFilePath)
	require.NoError(t, err)
	assert.True(t, outputExists)

	outputContent, err := afero.ReadFile(filesystem, outputFilePath)
	require.NoError(t, err)
	assert.Contains(t, string(outputContent), expectedMarkdownTable)
}

func TestGenerateDocumentationRendersMultipleJobsAndComponentsCorrectly(t *testing.T) {
	t.Parallel()

	firstComponentContent := `---
# First Component
spec:
  inputs:
    stage:
...
---
# First Component first job
first-component-first-job: {}

# First Component second job
first-component-second-job: {}
`

	secondComponentContent := `---
# Second Component
spec:
  inputs:
    stage:
...
---
# Second Component first job
second-component-first-job: {}

# Second Component second job
second-component-second-job: {}
...
`

	expectedMarkdown := `# Components Documentation

## Components

### first-component

First Component
` +
		"\n#### Usage of component `first-component`" +
		"\n" +
		"\nYou can add this component to an existing `.gitlab-ci.yml` file by using the `include:` keyword.\n" +
		"\n" +
		"```yaml\n" +
		"include:\n" +
		"  - component: \"github.com/test/first-component@1.0.0\"\n" +
		"    inputs: {}\n" +
		"```\n" +
		`
You can configure the component with the inputs documented below.
` +
		"\n#### Inputs of component `first-component`\n" +
		`
| Name | Description | Type | Default | Options | Regex | Mandatory |
|------|-------------|------|---------|---------|-------|-----------|
` +
		"| `stage` |  | `-` | `-` | `-` | `-` | Yes |\n" +
		"\n" +
		"#### Jobs of component `first-component`\n" +
		`
The component will add the following jobs to your CI/CD Pipeline.
` +
		"\n##### `first-component-first-job`\n" +
		`
First Component first job
` +
		"\n##### `first-component-second-job`\n" +
		`
First Component second job

### second-component

Second Component
` +
		"\n#### Usage of component `second-component`" +
		"\n" +
		"\nYou can add this component to an existing `.gitlab-ci.yml` file by using the `include:` keyword.\n" +
		"\n" +
		"```yaml\n" +
		"include:\n" +
		"  - component: \"github.com/test/second-component@1.0.0\"\n" +
		"    inputs: {}\n" +
		"```\n" +
		`
You can configure the component with the inputs documented below.
` +
		"\n#### Inputs of component `second-component`\n" +
		`
| Name | Description | Type | Default | Options | Regex | Mandatory |
|------|-------------|------|---------|---------|-------|-----------|
` +
		"| `stage` |  | `-` | `-` | `-` | `-` | Yes |\n" +
		"\n" +
		"#### Jobs of component `second-component`\n" +
		`
The component will add the following jobs to your CI/CD Pipeline.
` +
		"\n##### `second-component-first-job`\n" +
		`
Second Component first job
` +
		"\n##### `second-component-second-job`\n" +
		`
Second Component second job
`

	filesystem := afero.NewMemMapFs()
	componentDir := "templates"
	firstComponentPath := componentDir + "/first-component.yml"
	secondComponentPath := componentDir + "/second-component.yml"
	outputFilePath := "README.md"

	err := filesystem.Mkdir(componentDir, 0o644)
	require.NoError(t, err)
	err = afero.WriteFile(filesystem, firstComponentPath, []byte(firstComponentContent), 0o644)
	require.NoError(t, err)
	err = afero.WriteFile(filesystem, secondComponentPath, []byte(secondComponentContent), 0o644)
	require.NoError(t, err)

	documentationGenerator := &RealDocumentationGenerator{}
	documentationGenerator.GenerateDocumentation(
		filesystem,
		"templates",
		"resources/default-template.md.gotmpl",
		"github.com/test",
		"1.0.0",
		"README.md",
		false,
	)

	outputExists, err := afero.Exists(filesystem, outputFilePath)
	require.NoError(t, err)
	assert.True(t, outputExists)

	outputContent, err := afero.ReadFile(filesystem, outputFilePath)
	require.NoError(t, err)
	assert.Equal(t, expectedMarkdown, string(outputContent))
}

func TestGenerateDocumentationUsesCustomTemplate(t *testing.T) {
	t.Parallel()

	componentContent := `---
# This is a custom template
spec:
  inputs:
    stage:
...
`

	customTemplateContent := `
{{- range $component := .Components }}
Description: {{ $component.Description }}
{{- end }}
`

	filesystem := afero.NewMemMapFs()
	componentDir := "templates"
	componentPath := componentDir + "/first-component.yml"
	outputFilePath := "README.md"
	customTemplateFilePath := "my-template.yml"

	err := filesystem.Mkdir(componentDir, 0o644)
	require.NoError(t, err)
	err = afero.WriteFile(filesystem, componentPath, []byte(componentContent), 0o644)
	require.NoError(t, err)
	err = afero.WriteFile(filesystem, customTemplateFilePath, []byte(customTemplateContent), 0o644)
	require.NoError(t, err)

	documentationGenerator := &RealDocumentationGenerator{}
	documentationGenerator.GenerateDocumentation(
		filesystem,
		"templates",
		customTemplateFilePath,
		"github.com/test",
		"1.0.0",
		"README.md",
		false)

	outputExists, err := afero.Exists(filesystem, outputFilePath)
	require.NoError(t, err)
	assert.True(t, outputExists)

	outputContent, err := afero.ReadFile(filesystem, outputFilePath)
	require.NoError(t, err)
	assert.Contains(t, string(outputContent), "Description: This is a custom template")
}

func TestBuildComponentDocumentationFromComponentsBuildsSortedComponentDocumentation(t *testing.T) {
	t.Parallel()

	components := []Component{
		{
			Name:        "ComponentB",
			Description: "Second component",
			Inputs: []Input{
				{Name: "InputB", Description: "Second input"},
				{Name: "InputA", Description: "First input"},
			},
			Jobs: []Job{
				{Name: "JobB", Comment: "Second job"},
				{Name: "JobA", Comment: "First job"},
			},
		},
		{
			Name:        "ComponentA",
			Description: "First component",
			Inputs: []Input{
				{Name: "InputB", Description: "Second input"},
				{Name: "InputA", Description: "First input"},
			},
			Jobs: []Job{
				{Name: "JobB", Comment: "Second job"},
				{Name: "JobA", Comment: "First job"},
			},
		},
	}

	repoURL := "https://example.com/repo"
	version := "1.0.0"

	// Expected sorted components, inputs, and jobs
	expectedComponents := []Component{
		{
			Name:        "ComponentA",
			Description: "First component",
			Inputs: []Input{
				{Name: "InputA", Description: "First input"},
				{Name: "InputB", Description: "Second input"},
			},
			Jobs: []Job{
				{Name: "JobA", Comment: "First job"},
				{Name: "JobB", Comment: "Second job"},
			},
		},
		{
			Name:        "ComponentB",
			Description: "Second component",
			Inputs: []Input{
				{Name: "InputA", Description: "First input"},
				{Name: "InputB", Description: "Second input"},
			},
			Jobs: []Job{
				{Name: "JobA", Comment: "First job"},
				{Name: "JobB", Comment: "Second job"},
			},
		},
	}

	expected := ComponentsDocumentation{
		Components: expectedComponents,
		RepoURL:    repoURL,
		Version:    version,
	}

	result := buildComponentDocumentationFromComponents(components, repoURL, version)

	assert.Equal(t, expected, result)
}

func TestReadTemplateFileReadsEmbeddedFile(t *testing.T) {
	t.Parallel()

	filesystem := afero.NewMemMapFs()

	fileContent := readTemplateFile("resources/default-template.md.gotmpl", filesystem)

	assert.Contains(t, fileContent, "You can add this component to an existing `.gitlab-ci.yml` file")
}

func TestReadTemplateFileReadsCustomFile(t *testing.T) {
	t.Parallel()

	customTemplateContent := "# Component Documentation"
	filesystem := afero.NewMemMapFs()
	err := afero.WriteFile(filesystem, "template.md", []byte(customTemplateContent), 0o644)
	require.NoError(t, err)

	fileContent := readTemplateFile("template.md", filesystem)

	assert.Contains(t, fileContent, customTemplateContent)
}

func TestWriteDocumentation(t *testing.T) {
	t.Parallel()

	filesystem := afero.NewMemMapFs()
	outputFilePath := "README.md"
	documentationContent := "# Sample Documentation"

	writeDocumentation(filesystem, outputFilePath, documentationContent)

	outputExists, err := afero.Exists(filesystem, outputFilePath)
	require.NoError(t, err)
	assert.True(t, outputExists)

	outputContent, err := afero.ReadFile(filesystem, outputFilePath)
	require.NoError(t, err)
	assert.Equal(t, documentationContent, string(outputContent))
}

func TestCompareExistingDocumentationShowsNoDifference(t *testing.T) {
	t.Parallel()

	filesystem := afero.NewMemMapFs()
	outputFilePath := "README.md"
	documentationContent := "# Sample Documentation"

	err := afero.WriteFile(filesystem, outputFilePath, []byte(documentationContent), 0o644)
	require.NoError(t, err)

	err = compareExistingDocumentation(filesystem, outputFilePath, documentationContent)
	require.NoError(t, err)
}

func TestCompareExistingDocumentationShowsDifference(t *testing.T) {
	t.Parallel()

	filesystem := afero.NewMemMapFs()
	outputFilePath := "README.md"
	oldDocumentationContent := "# Old Documentation"
	newDocumentationContent := "# New Documentation"

	err := afero.WriteFile(filesystem, outputFilePath, []byte(oldDocumentationContent), 0o644)
	require.NoError(t, err)

	err = compareExistingDocumentation(filesystem, outputFilePath, newDocumentationContent)
	require.Error(t, err)
	assert.Equal(t, "documentation is not up-to-date. changes have been detected", err.Error())
}

func TestRenderDocumentationContentWithValidTemplate(t *testing.T) {
	t.Parallel()

	filesystem := afero.NewMemMapFs()
	templateFilePath := "template.md"
	templateContent := `# Components Documentation
{{ range .Components }}
## {{ .Name }}
Description: {{ .Description }}
{{ end }}
`
	err := afero.WriteFile(filesystem, templateFilePath, []byte(templateContent), 0o644)
	require.NoError(t, err)

	components := []Component{
		{Name: "Component1", Description: "Description1"},
		{Name: "Component2", Description: "Description2"},
	}
	componentsDocumentation := ComponentsDocumentation{
		Components: components,
		RepoURL:    "http://repo.url",
		Version:    "v1.0",
	}

	expectedContent := `# Components Documentation

## Component1
Description: Description1

## Component2
Description: Description2

`
	actualContent := renderDocumentationContent(componentsDocumentation, templateFilePath, filesystem)
	assert.Equal(t, expectedContent, actualContent)
}

func TestSortJobsCorrectlySortsJobs(t *testing.T) {
	t.Parallel()

	jobA := Job{
		Name: "A",
	}
	jobB := Job{
		Name: "B",
	}
	jobs := []Job{
		jobB,
		jobA,
	}

	expectedJobs := []Job{
		jobA,
		jobB,
	}

	sortJobs(jobs)
	assert.Equal(t, expectedJobs, jobs)
}

func TestSortSpecInputsCorrectlySortsInputs(t *testing.T) {
	t.Parallel()

	inputA := Input{
		Name: "A",
	}
	inputB := Input{
		Name: "B",
	}
	inputs := []Input{
		inputB,
		inputA,
	}

	expectedInputs := []Input{
		inputA,
		inputB,
	}

	sortInputs(inputs)
	assert.Equal(t, expectedInputs, inputs)
}

func TestSortComponentsCorrectlySortsComponents(t *testing.T) {
	t.Parallel()

	components := []Component{
		{Name: "beta"},
		{Name: "alpha"},
		{Name: "gamma"},
	}

	expectedComponents := []Component{
		{Name: "alpha"},
		{Name: "beta"},
		{Name: "gamma"},
	}

	actualComponents := sortComponents(components)

	assert.Equal(t, expectedComponents, actualComponents)
}
