package cmd

import (
	"github.com/erNail/labdoc/internal/gitlab"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command for the CLI application.
// This command serves as the entry point and parent for all other commands.
//
// Returns:
//   - *cobra.Command: A pointer to the newly created cobra.Command.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "labdoc",
		Short: "Generate Markdown documentation from GitLab CI/CD Components",
		Long:  "A CLI tool for generating Markdown documentation from GitLab CI/CD Components",
	}

	filesystem := afero.NewOsFs()
	documentationGenerator := &gitlab.RealDocumentationGenerator{}
	rootCmd.AddCommand(NewGenerateCmd(filesystem, documentationGenerator))

	return rootCmd
}

func Execute() {
	cmd := NewRootCmd()

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
