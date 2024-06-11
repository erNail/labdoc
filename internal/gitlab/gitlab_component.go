package gitlab

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/erNail/labdoc/internal/yamlutils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// CiConfig represents the GitLab CI configuration.
type CiConfig struct {
	Spec Spec  `yaml:"spec"`
	Jobs []Job `yaml:"-"`
}

// Spec defines the "spec" keyword of the GitLab CI configuration.
type Spec struct {
	Inputs  []Input `yaml:"inputs"`
	Comment string  `yaml:"-"`
}

// Input represents an input parameter for the GitLab CI spec.
type Input struct {
	Name        string        `yaml:"-"`
	Description string        `yaml:"description,omitempty"`
	Type        string        `yaml:"type,omitempty"`
	Default     interface{}   `yaml:"default,omitempty"`
	Options     []interface{} `yaml:"options,omitempty"`
	Regex       string        `yaml:"regex,omitempty"`
}

// Job represents a job in the GitLab CI configuration.
type Job struct {
	Name    string
	Comment string
}

// Component represents a GitLab CI component.
type Component struct {
	Jobs        []Job
	Description string
	Name        string
	Inputs      []Input
}

// UnmarshalYAML is called when using yaml.Unmarshal on a GitLabCiConfig type.
//
// Parameters:
//   - node: The YAML node to unmarshal.
//
// Returns:
//   - error: An error if unmarshalling fails.
func (gitlabCiConfig *CiConfig) UnmarshalYAML(node *yaml.Node) error {
	gitlabCiConfig.Jobs = []Job{}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		key := keyNode.Value

		if key == "spec" {
			gitlabCiConfig.Spec.Inputs = parseSpecInputs(*valueNode)
			gitlabCiConfig.Spec.Comment = yamlutils.FormatCommentAsPlainText(keyNode.HeadComment)
		} else if isJobMappingNode(key, *valueNode) {
			job := Job{
				Name:    key,
				Comment: yamlutils.FormatCommentAsPlainText(keyNode.HeadComment),
			}
			gitlabCiConfig.Jobs = append(gitlabCiConfig.Jobs, job)
		}
	}

	return nil
}

// parseSpecInputs parses the inputs of a spec node.
//
// Parameters:
//   - specNode: The YAML node containing the spec inputs.
//
// Returns:
//   - []Input: A slice of Input structs.
func parseSpecInputs(specNode yaml.Node) []Input {
	inputs := []Input{}

	var temp struct {
		Inputs map[string]yaml.Node `yaml:"inputs"`
	}

	if err := specNode.Decode(&temp); err != nil {
		log.Fatal(err)
	}

	for key, inputNode := range temp.Inputs {
		var input Input
		if err := inputNode.Decode(&input); err != nil {
			log.Fatal(err)
		}

		input.Name = key
		inputs = append(inputs, input)
	}

	return inputs
}

// isJobMappingNode checks if a given YAML node represents a job mapping node.
//
// Parameters:
//   - key: The key of the YAML node.
//   - valueNode: The value node of the YAML node.
//
// Returns:
//   - bool: True if the node is a job mapping node, false otherwise.
func isJobMappingNode(key string, valueNode yaml.Node) bool {
	nonJobTopLevelKeywords := []string{
		"stages",
		"default",
		"include",
		"workflow",
		"spec",
		"image",
		"services",
		"cache",
		"before_script",
		"after_script",
	}

	if valueNode.Kind != yaml.MappingNode {
		return false
	}

	if slices.Contains(nonJobTopLevelKeywords, key) {
		return false
	}

	return true
}

// parseYamlFileWithoutSeparatorsToGitLabCiConfig parses a YAML file content into a CiConfig,
// removing YAML document separators.
//
// Parameters:
//   - yamlContent: The content of the YAML file.
//
// Returns:
//   - CiConfig: The parsed CiConfig struct.
func parseYamlFileWithoutSeparatorsToGitLabCiConfig(yamlContent []byte) CiConfig {
	yamlContentWithoutSeparators := yamlutils.RemoveYamlDocumentSeparators(yamlContent)

	var gitlabCiConfig CiConfig

	err := yaml.Unmarshal(yamlContentWithoutSeparators, &gitlabCiConfig)
	if err != nil {
		log.Fatalf("failed to parse YAML: %v", err)
	}

	return gitlabCiConfig
}

// newComponentFromGitLabCiConfig creates a new Component from a CiConfig and a component name.
//
// Parameters:
//   - gitlabCiConfig: The CiConfig struct.
//   - componentName: The name of the component.
//
// Returns:
//   - Component: The constructed Component struct.
func newComponentFromGitLabCiConfig(gitlabCiConfig CiConfig, componentName string) Component {
	component := Component{
		Jobs:        gitlabCiConfig.Jobs,
		Inputs:      gitlabCiConfig.Spec.Inputs,
		Description: gitlabCiConfig.Spec.Comment,
		Name:        componentName,
	}

	return component
}

// generateComponentNameFromFilePath generates a component name from a file path.
//
// Parameters:
//   - filePath: The file path.
//
// Returns:
//   - string: The generated component name.
func generateComponentNameFromFilePath(filePath string) string {
	filenameWithoutPath := filepath.Base(filePath)
	filename := strings.TrimSuffix(filenameWithoutPath, filepath.Ext(filenameWithoutPath))

	return filename
}
