package gitlab

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateDocumentationUsesEmbeddedTemplate(t *testing.T) {
	t.Parallel()

	firstComponentContent := `---
# This is a description of the first component
spec:
  inputs:
    stage:
      description: "The stage of the jobs"
      type: "string"
      default: "test"
...
---
# This is the first job of the first component
first-component-first-job: {}

# This is the second job of the first component
first-component-second-job: {}
`

	secondComponentContent := `---
# This is a description of the second component
spec:
  inputs:
    stage:
      description: "The stage of the jobs"
      type: "string"
      default: "test"
...
---
# This is the first job of the second component
second-component-first-job: {}

# This is the second job of the second component
second-component-second-job: {}
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
	assert.Contains(t, string(outputContent), "### first-component")
	assert.Contains(t, string(outputContent), "### second-component")
	assert.Contains(t, string(outputContent), "##### `second-component-second-job`")
}

func TestGenerateDocumentationUsesCustomTemplate(t *testing.T) {
	t.Parallel()

	firstComponentContent := `---
# This is a description of the first component
spec:
  inputs:
    stage:
      description: "The stage of the jobs"
      type: "string"
      default: "test"
...
---
# This is the first job of the first component
first-component-first-job: {}

# This is the second job of the first component
first-component-second-job: {}
`

	secondComponentContent := `---
# This is a description of the second component
spec:
  inputs:
    stage:
      description: "The stage of the jobs"
      type: "string"
      default: "test"
...
---
# This is the first job of the second component
second-component-first-job: {}

# This is the second job of the second component
second-component-second-job: {}
`

	customTemplateContent := `
	# Components Documentation

	{{ range $component := .Components }}
	This is component {{ $component.Name }}
	{{ end }}
	`

	filesystem := afero.NewMemMapFs()
	componentDir := "templates"
	firstComponentPath := componentDir + "/first-component.yml"
	secondComponentPath := componentDir + "/second-component.yml"
	outputFilePath := "README.md"
	customTemplateFilePath := "my-template.yml"

	err := filesystem.Mkdir(componentDir, 0o644)
	require.NoError(t, err)
	err = afero.WriteFile(filesystem, firstComponentPath, []byte(firstComponentContent), 0o644)
	require.NoError(t, err)
	err = afero.WriteFile(filesystem, secondComponentPath, []byte(secondComponentContent), 0o644)
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
	assert.Contains(t, string(outputContent), "This is component first-component")
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
