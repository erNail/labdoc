package yamlutils

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadYamlFilesReadsAllYamlFilesInTheCurrentWorkingDirectory(t *testing.T) {
	t.Parallel()

	filesystem := afero.NewMemMapFs()
	firstFileName := "first.yml"
	secondFileName := "second.yml"
	firstFileContent := "test: first"
	secondFileContent := "test: second"

	err := afero.WriteFile(filesystem, firstFileName, []byte(firstFileContent), 0o644)
	require.NoError(t, err)
	err = afero.WriteFile(filesystem, secondFileName, []byte(secondFileContent), 0o644)
	require.NoError(t, err)

	expectedFileContentMap := map[string][]byte{
		"first.yml":  []byte(firstFileContent),
		"second.yml": []byte(secondFileContent),
	}

	actualFileContentMap := ReadYamlFilesFromDirectory(filesystem, "")
	assert.Equal(t, expectedFileContentMap, actualFileContentMap)
}

func TestReadYamlFilesReadsAllYamlFilesInASubdirectory(t *testing.T) {
	t.Parallel()

	filesystem := afero.NewMemMapFs()
	dirName := "templates/"
	firstFileName := "first.yml"
	secondFileName := "second.yml"
	firstFileContent := "test: first"
	secondFileContent := "test: second"

	err := filesystem.Mkdir(dirName, 0o755)
	require.NoError(t, err)

	err = afero.WriteFile(filesystem, dirName+firstFileName, []byte(firstFileContent), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(filesystem, dirName+secondFileName, []byte(secondFileContent), 0o644)
	require.NoError(t, err)

	expectedFileContentMap := map[string][]byte{
		"templates/first.yml":  []byte(firstFileContent),
		"templates/second.yml": []byte(secondFileContent),
	}

	actualFileContentMap := ReadYamlFilesFromDirectory(filesystem, "templates")
	assert.Equal(t, expectedFileContentMap, actualFileContentMap)
}

func TestRemoveYamlDocumentSeparatorsRemovesAllSeparators(t *testing.T) {
	t.Parallel()

	yamlFileContent := `---
This
---
should
...
be
---
without
...
---
separators`

	expectedYamlContent := "\nThis\n\nshould\n\nbe\n\nwithout\n\n\nseparators"
	cleanedYamlContent := RemoveYamlDocumentSeparators([]byte(yamlFileContent))
	assert.YAMLEq(t, expectedYamlContent, string(cleanedYamlContent))
}

func TestRemoveYamlDocumentSeparatorsKeepsSeparatorsIfWithinString(t *testing.T) {
	t.Parallel()

	yamlFileContent := `---
This will stay: --- ...
...`

	expectedYamlContent := "\nThis will stay: --- ...\n"
	cleanedYamlContent := RemoveYamlDocumentSeparators([]byte(yamlFileContent))
	assert.YAMLEq(t, expectedYamlContent, string(cleanedYamlContent))
}

func TestFormatCommentAsPlainTextFormatsSingleLineComment(t *testing.T) {
	t.Parallel()

	commentWithSpace := "# Comment with space"
	commentWithoutSpace := "#Comment without space"
	cleanedCommentWithSpace := FormatCommentAsPlainText(commentWithSpace)
	cleanedCommentWithoutSpace := FormatCommentAsPlainText(commentWithoutSpace)

	assert.Equal(t, "Comment with space", cleanedCommentWithSpace)
	assert.Equal(t, "Comment without space", cleanedCommentWithoutSpace)
}

func TestFormatCommentAsPlainTextFormatsMultiLineComment(t *testing.T) {
	t.Parallel()

	multilineComment := `# Comment with space
#Which is continued on the next line
#
# And also handles empty lines`

	expectedFormattedComment := "Comment with space\nWhich is continued on the next line\n\nAnd also handles empty lines"

	actualFormattedComment := FormatCommentAsPlainText(multilineComment)

	assert.Equal(t, expectedFormattedComment, actualFormattedComment)
}

func TestFormatCommentAsPlainTextOnlyRemovesHashSymbolInTheBeginning(t *testing.T) {
	t.Parallel()

	multilineComment := `# This is a comment with a # in it
# This is a comment with an [anchor link](#anchor) in it`

	expectedFormattedComment := "This is a comment with a # in it\nThis is a comment with an [anchor link](#anchor) in it"

	actualFormattedComment := FormatCommentAsPlainText(multilineComment)

	assert.Equal(t, expectedFormattedComment, actualFormattedComment)
}

func TestReadYamlFilesFromDirectoryWithNoYamlFilesReturnsNoYamlFiles(t *testing.T) {
	t.Parallel()

	filesystem := afero.NewMemMapFs()
	err := filesystem.Mkdir("empty_dir", 0o755)
	require.NoError(t, err)

	actualFileContentMap := ReadYamlFilesFromDirectory(filesystem, "empty_dir")
	assert.Empty(t, actualFileContentMap)
}
