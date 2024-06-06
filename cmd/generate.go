package cmd

import (
	"github.com/erNail/labdoc/internal/gitlab"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// NewGenerateCmd creates a new command for generating documentation
// for GitLab CI/CD components. It processes YAML files from a specified
// directory using a provided or default template.
//
// Parameters:
//   - filesystem: An interface for interacting with the file system.
//   - documentationGenerator: An interface for generating documentation.
//
// Returns:
//   - *cobra.Command: A pointer to the newly created cobra.Command.
func NewGenerateCmd(filesystem afero.Fs, documentationGenerator gitlab.DocumentationGenerator) *cobra.Command {
	var (
		repoURL          string
		componentVersion string
		componentDir     string
		templateFilePath string
		outputFilePath   string
		checkOnly        bool
	)

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate documentation for GitLab CI/CD components",
		Long:  `Generate documentation for GitLab CI/CD components from a directory of CI/CD components`,
		Run: func(_ *cobra.Command, _ []string) {
			documentationGenerator.GenerateDocumentation(
				filesystem,
				componentDir,
				templateFilePath,
				repoURL,
				componentVersion,
				outputFilePath,
				checkOnly,
			)
		},
	}

	repoURLFlag := "repoUrl"

	generateCmd.Flags().StringVarP(
		&repoURL, repoURLFlag, "r", "",
		"The repository URL from which to include the GitLab CI/CD Component. Will be used in the documentation (required)",
	)
	generateCmd.Flags().StringVarP(
		&componentVersion, "version", "v", "latest",
		"The current version or ref of the GitLab CI/CD Component. Will be used in the documentation",
	)
	generateCmd.Flags().StringVarP(
		&componentDir, "componentDir", "d", "templates",
		"The directory containing the GitLab CI/CD components",
	)
	generateCmd.Flags().StringVarP(
		&templateFilePath, "template", "t", "resources/default-template.md.gotmpl",
		"The template file from which the documentation is generated",
	)
	generateCmd.Flags().StringVarP(
		&outputFilePath, "outputFile", "o", "templates/README.md",
		"The path and name of the rendered file to be created",
	)
	generateCmd.Flags().BoolVarP(
		&checkOnly, "check", "c", false,
		"If set, will check if the documentation is up-to-date. If not, the application will exit with exit code 2",
	)

	err := generateCmd.MarkFlagRequired(repoURLFlag)
	if err != nil {
		log.Fatal(err)
	}

	return generateCmd
}
