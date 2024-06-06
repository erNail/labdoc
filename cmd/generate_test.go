package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockDocumentationGenerator struct {
	mock.Mock
}

func (m *MockDocumentationGenerator) GenerateDocumentation(filesystem afero.Fs,
	componentDirectory string,
	templateFilePath string,
	repoURL string,
	componentVersion string,
	outputFilePath string,
	checkOnly bool,
) {
	m.Called(filesystem, componentDirectory, templateFilePath, repoURL, componentVersion, outputFilePath, checkOnly)
}

func TestGenerateCmdThrowsErrorIfRepoUrlIsNotSet(t *testing.T) {
	t.Parallel()

	cmd := NewGenerateCmd(afero.NewMemMapFs(), new(MockDocumentationGenerator))
	cmd.SetArgs([]string{"--version='1.0.0'"})

	err := cmd.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), `required flag(s) "repoUrl" not set`)
}

func TestGenerateCmdIsSuccessfulWithRequiredParametersSet(t *testing.T) {
	t.Parallel()

	filesystem := afero.NewMemMapFs()
	mockDocumentationGenerator := new(MockDocumentationGenerator)
	mockDocumentationGenerator.On(
		"GenerateDocumentation",
		filesystem,
		"templates",
		"resources/default-template.md.gotmpl",
		"github.com/test",
		"latest",
		"templates/README.md",
		false,
	).Return()

	cmd := NewGenerateCmd(filesystem, mockDocumentationGenerator)
	cmd.SetArgs([]string{"--repoUrl=github.com/test"})

	err := cmd.Execute()

	require.NoError(t, err)
	mockDocumentationGenerator.AssertExpectations(t)
}
