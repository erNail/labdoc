package yamlutils

import (
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// ReadYamlFilesFromDirectory reads all YAML files from the specified directory
// and returns a map where the keys are file paths and the values are the file contents.
//
// Parameters:
//   - filesystem: An interface for interacting with the file system.
//   - directory: The directory from which to read YAML files.
//
// Returns:
//   - map[string][]byte: A map of file paths to their respective contents.
func ReadYamlFilesFromDirectory(filesystem afero.Fs, directory string) map[string][]byte {
	filePathToContentMap := make(map[string][]byte)

	yamlFilePaths, err := afero.Glob(filesystem, filepath.Join(directory, "*.yml"))
	if err != nil {
		log.Fatal(err)
	}

	for _, yamlFilePath := range yamlFilePaths {
		yamlFileContent, err := afero.ReadFile(filesystem, yamlFilePath)
		if err != nil {
			log.Fatal(err)
		}

		filePathToContentMap[yamlFilePath] = yamlFileContent
	}

	return filePathToContentMap
}

// RemoveYamlDocumentSeparators removes YAML document separators ("---" and "...")
// from the given YAML content.
//
// Parameters:
//   - yamlContent: The content of the YAML file.
//
// Returns:
//   - []byte: The cleaned YAML content without document separators.
func RemoveYamlDocumentSeparators(yamlContent []byte) []byte {
	re := regexp.MustCompile(`(?m)^---$|^...$`)
	cleanedYamlContent := re.ReplaceAllString(string(yamlContent), "")

	return []byte(cleanedYamlContent)
}

// FormatCommentAsPlainText removes `#` characters from the comment and handles
// multiline whitespaces, returning a cleaned comment string.
//
// Parameters:
//   - comment: The comment string to format.
//
// Returns:
//   - string: The cleaned comment string.
func FormatCommentAsPlainText(comment string) string {
	lines := strings.Split(comment, "\n")

	for i, line := range lines {
		lines[i] = strings.TrimSpace(strings.TrimPrefix(line, "#"))
	}

	cleanedComment := strings.Join(lines, "\n")

	return cleanedComment
}
