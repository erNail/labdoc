package gitlab

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"text/template"

	"github.com/erNail/labdoc/internal/yamlutils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// ComponentsDocumentation represents the data needed to document GitLab CI/CD components.
type ComponentsDocumentation struct {
	RepoURL    string
	Version    string
	Components []Component
}

// DocumentationGenerator defines the interface for generating documentation.
type DocumentationGenerator interface {
	GenerateDocumentation(
		filesystem afero.Fs,
		componentDirectory string,
		templateFilePath string,
		repoURL string,
		componentVersion string,
		outputFilePath string,
		checkOnly bool)
}

// RealDocumentationGenerator implements the DocumentationGenerator interface.
type RealDocumentationGenerator struct{}

// GenerateDocumentation generates documentation for GitLab CI/CD components.
// It processes the components in the specified directory using a given template.
//
// Parameters:
//   - filesystem: An interface for interacting with the file system.
//   - componentDirectory: The directory containing the component YAML files.
//   - templateFilePath: The path to the template file used for generating documentation.
//   - repoURL: The URL of the repository containing the components.
//   - componentVersion: The version or ref of the components to document.
//   - outputFilePath: The path where the generated documentation will be saved.
//   - checkOnly: If true, checks if the documentation is up-to-date without writing the file.
func (r *RealDocumentationGenerator) GenerateDocumentation(
	filesystem afero.Fs,
	componentDirectory string,
	templateFilePath string,
	repoURL string,
	componentVersion string,
	outputFilePath string,
	checkOnly bool,
) {
	log.Info("Generating documentation...")

	filePathContentMap := yamlutils.ReadYamlFilesFromDirectory(filesystem, componentDirectory)
	if len(filePathContentMap) == 0 {
		log.WithField("componentsDir", componentDirectory).Fatal("No files found in directory")
	}

	components := []Component{}

	for filePath, componentFileContent := range filePathContentMap {
		gitlabCiConfig := parseYamlFileWithoutSeparatorsToGitLabCiConfig(componentFileContent)
		component := newComponentFromGitLabCiConfig(gitlabCiConfig, generateComponentNameFromFilePath(filePath))
		components = append(components, component)
	}

	log.WithField("componentCount", len(components)).Info("Found components")
	componentsDocumentation := buildComponentDocumentationFromComponents(components, repoURL, componentVersion)
	documentationContent := renderDocumentationContent(componentsDocumentation, templateFilePath, filesystem)

	if checkOnly {
		err := compareExistingDocumentation(filesystem, outputFilePath, documentationContent)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		writeDocumentation(filesystem, outputFilePath, documentationContent)
	}
}

// buildComponentDocumentationFromComponents creates a ComponentsDocumentation
// struct from the given components, repository URL, and version.
//
// Parameters:
//   - components: A slice of Component structs to document.
//   - repoURL: The URL of the repository containing the components.
//   - version: The version or ref of the components to document.
//
// Returns:
//   - ComponentsDocumentation: The constructed ComponentsDocumentation struct.
func buildComponentDocumentationFromComponents(
	components []Component,
	repoURL string,
	version string,
) ComponentsDocumentation {
	for _, component := range components {
		component.Inputs = sortInputs(component.Inputs)
		component.Jobs = sortJobs(component.Jobs)
	}

	componentDocumentation := ComponentsDocumentation{
		Components: components,
		RepoURL:    repoURL,
		Version:    version,
	}

	return componentDocumentation
}

// renderDocumentationContent renders the documentation content using the
// specified template and component documentation data.
//
// Parameters:
//   - componentsDocumentation: The data for the components to document.
//   - templateFilePath: The path to the template file used for generating documentation.
//   - filesystem: An interface for interacting with the file system.
//
// Returns:
//   - string: The rendered documentation content.
func renderDocumentationContent(
	componentsDocumentation ComponentsDocumentation,
	templateFilePath string,
	filesystem afero.Fs,
) string {
	templateFileContent := readTemplateFile(templateFilePath, filesystem)

	tmpl, err := template.New(templateFilePath).Parse(templateFileContent)
	if err != nil {
		log.Fatal(err)
	}

	buffer := new(bytes.Buffer)

	err = tmpl.Execute(buffer, componentsDocumentation)
	if err != nil {
		log.Fatal(err)
	}

	return buffer.String()
}

//go:embed resources/default-template.md.gotmpl
var embedFs embed.FS

// readTemplateFile reads the content of the template file from the given
// file system. If the templateFilePath is the default template, it reads
// from the embedded file system.
//
// Parameters:
//   - templateFilePath: The path to the template file used for generating documentation.
//   - filesystem: An interface for interacting with the file system.
//
// Returns:
//   - string: The content of the template file.
func readTemplateFile(templateFilePath string, filesystem afero.Fs) string {
	if templateFilePath == "resources/default-template.md.gotmpl" {
		log.Info("Using default template")

		templateFileContent, err := fs.ReadFile(embedFs, templateFilePath)
		if err != nil {
			log.Fatal(err)
		}

		return string(templateFileContent)
	}

	log.WithField("filePath", templateFilePath).Info("Using custom template")

	templateFileContent, err := afero.ReadFile(filesystem, templateFilePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(templateFileContent)
}

// writeDocumentation writes the generated documentation content to the specified
// output file path.
//
// Parameters:
//   - filesystem: An interface for interacting with the file system.
//   - outputFilePath: The path where the generated documentation will be saved.
//   - documentationContent: The content of the generated documentation.
func writeDocumentation(filesystem afero.Fs, outputFilePath string, documentationContent string) {
	err := afero.WriteFile(filesystem, outputFilePath, []byte(documentationContent), 0o644)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Generated documentation!")
}

// compareExistingDocumentation compares the existing documentation content with the new content.
// If they differ, it returns an error indicating that the documentation is not up-to-date.
//
// Parameters:
//   - filesystem: An interface for interacting with the file system.
//   - outputFilePath: The path to the existing documentation file.
//   - newDocumentationContent: The new documentation content to compare.
//
// Returns:
//   - error: An error if the documentation is not up-to-date.
func compareExistingDocumentation(filesystem afero.Fs, outputFilePath string, newDocumentationContent string) error {
	log.Info("Running in check mode. No file will be written.")

	oldDocumentationContent, err := afero.ReadFile(filesystem, outputFilePath)
	if err != nil {
		return fmt.Errorf("documentation does not exist: %w", err)
	}

	if string(oldDocumentationContent) != newDocumentationContent {
		err := errors.New("documentation is not up-to-date. changes have been detected")

		return err
	}

	log.Info("Your documentation is up-to-date!")

	return nil
}

// sortJobs sorts a slice of Job structs by their name.
//
// Parameters:
//   - jobs: A slice of Job structs.
//
// Returns:
//   - []Job: The sorted slice of Job structs.
func sortJobs(jobs []Job) []Job {
	slices.SortFunc(jobs, func(a Job, b Job) int {
		if a.Name < b.Name {
			return -1
		}

		if a.Name > b.Name {
			return 1
		}

		return 0
	})

	return jobs
}

// sortInputs sorts a slice of Input structs by their name.
//
// Parameters:
//   - inputs: A slice of Input structs.
//
// Returns:
//   - []Input: The sorted slice of Input structs.
func sortInputs(inputs []Input) []Input {
	slices.SortFunc(inputs, func(a Input, b Input) int {
		if a.Name < b.Name {
			return -1
		}

		if a.Name > b.Name {
			return 1
		}

		return 0
	})

	return inputs
}
