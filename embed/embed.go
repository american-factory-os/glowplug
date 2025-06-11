// / Package embed provides access to embedded files.
package embed

import "embed"

var (
	//go:embed assets/index.html
	embeddedFiles embed.FS
	indexHTML     []byte
)

func init() {
	bytes, err := embeddedFiles.ReadFile("assets/index.html") // Ensure the embedded file is initialized
	if err != nil {
		panic("failed to read embedded file: " + err.Error())
	}
	indexHTML = bytes
}

func GetIndexHTML() []byte {
	return indexHTML
}
