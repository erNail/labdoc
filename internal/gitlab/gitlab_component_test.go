package gitlab

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestUnmarshalYAMLUnmarshalsToGitLabCiConfig(t *testing.T) {
	t.Parallel()

	yamlFileContent := `---
spec:
  inputs:
    string-input:
      description: "String description"
      type: "string"
      default: "test"
# Job comment
job: {}
`

	var gitlabCiConfig CiConfig

	err := yaml.Unmarshal([]byte(yamlFileContent), &gitlabCiConfig)
	require.NoError(t, err)
	assert.Len(t, gitlabCiConfig.Jobs, 1)
	assert.Equal(t, "job", gitlabCiConfig.Jobs[0].Name)
	assert.Len(t, gitlabCiConfig.Spec.Inputs, 1)
	assert.Equal(t, "string-input", gitlabCiConfig.Spec.Inputs[0].Name)
}

func TestParseYamlWithoutSeparatorsCreatesGitLabCiConfig(t *testing.T) {
	t.Parallel()

	yamlFileContent := `---
# Spec comment
spec:
  inputs:
    string-input:
      description: "String description"
      type: "string"
      default: "test"
    array-input:
      type: "array"
      default: []
    bool-input:
      description: "Boolean description"
      type: "boolean"
      default: false
    number-input:
      description: "Number description"
      default: 1
    options-input:
      options:
        - 1
    regex-input:
      regex: 'regex'
...
---
# First job comment
first-job: {}
# Second job comment
second-job: {}
`

	expectedGitLabCiConfig := CiConfig{
		Spec: Spec{
			Inputs: []Input{
				{
					Name:        "string-input",
					Description: "String description",
					Type:        "string",
					Default:     "test",
				},
				{
					Name:    "array-input",
					Type:    "array",
					Default: []interface{}{},
				},
				{
					Name:        "bool-input",
					Description: "Boolean description",
					Type:        "boolean",
					Default:     false,
				},
				{
					Name:        "number-input",
					Description: "Number description",
					Default:     1,
				},
				{
					Name:    "options-input",
					Options: []interface{}{1},
				},
				{
					Name:  "regex-input",
					Regex: "regex",
				},
			},
			Comment: "Spec comment",
		},
		Jobs: []Job{
			{
				Name:    "first-job",
				Comment: "First job comment",
			},
			{
				Name:    "second-job",
				Comment: "Second job comment",
			},
		},
	}

	actualGitLabCiConfig := parseYamlFileWithoutSeparatorsToGitLabCiConfig([]byte(yamlFileContent))
	assert.ElementsMatch(t, expectedGitLabCiConfig.Jobs, actualGitLabCiConfig.Jobs)
	assert.ElementsMatch(t, expectedGitLabCiConfig.Spec.Inputs, actualGitLabCiConfig.Spec.Inputs)
}

func TestNewComponentFromGitLabCiConfigCreatesValidGitLabCiConfig(t *testing.T) {
	t.Parallel()

	jobs := []Job{
		{
			Name:    "first-job",
			Comment: "First job comment",
		},
		{
			Name:    "second-job",
			Comment: "Second job comment",
		},
	}

	inputs := []Input{
		{
			Name:        "string-input",
			Description: "String description",
			Type:        "string",
			Default:     "test",
		},
		{
			Name:    "array-input",
			Type:    "array",
			Default: []interface{}{},
		},
		{
			Name:        "bool-input",
			Description: "Boolean description",
			Type:        "boolean",
			Default:     false,
		},
		{
			Name:        "number-input",
			Description: "Number description",
			Default:     1,
		},
		{
			Name:    "options-input",
			Options: []interface{}{1},
		},
		{
			Name:  "regex-input",
			Regex: "regex",
		},
	}

	gitlabCiConfig := CiConfig{
		Spec: Spec{
			Inputs:  inputs,
			Comment: "Spec comment",
		},
		Jobs: jobs,
	}

	expectedComponent := Component{
		Jobs:        jobs,
		Inputs:      inputs,
		Description: "Spec comment",
		Name:        "component-name",
	}
	actualComponent := newComponentFromGitLabCiConfig(gitlabCiConfig, "component-name")
	assert.Equal(t, expectedComponent, actualComponent)
}

func TestIsJobMappingNodeReturnsFalseOnNonJobTopLevelKeyword(t *testing.T) {
	t.Parallel()

	assert.False(t, isJobMappingNode("image", yaml.Node{}))
	assert.False(t, isJobMappingNode("workflow", yaml.Node{}))
	assert.False(t, isJobMappingNode("spec", yaml.Node{}))
}

func TestIsJobMappingNodeReturnsFalseOnNonMappingNode(t *testing.T) {
	t.Parallel()

	valueNode := yaml.Node{Kind: yaml.ScalarNode}
	assert.False(t, isJobMappingNode("job", valueNode))
}

func TestGenerateComponentNameFromFilePathResultsInCorrectName(t *testing.T) {
	t.Parallel()

	expectedName := "file"
	actualName := generateComponentNameFromFilePath("this/is/my/file.yml")
	assert.Equal(t, expectedName, actualName)

	actualName = generateComponentNameFromFilePath("file.yml")
	assert.Equal(t, expectedName, actualName)
}
