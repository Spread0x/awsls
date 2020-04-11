// +build codegen

package aws

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/jckuester/awsls/gen/util"
)

// WriteClient writes Go code to initialize all AWS API Clients
// used by Terraform AWS resources.
func WriteClient(outputPath string, services []string) error {
	err := os.MkdirAll(outputPath, 0775)
	if err != nil {
		return fmt.Errorf("failed to create directory: %s", err)
	}

	err = util.WriteGoFile(
		filepath.Join(outputPath, "client.go"),
		util.CodeLayout,
		"",
		"aws",
		clientGoCode(services),
	)

	if err != nil {
		return fmt.Errorf("failed to write AWSClient file: %s", err)
	}

	return nil
}

func clientGoCode(services []string) string {
	var buf bytes.Buffer
	err := clientTmpl.Execute(&buf, services)
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(buf.String())
}

var clientTmpl = template.Must(template.New("client").Parse(`import (
"github.com/aws/aws-sdk-go-v2/aws/external"
{{range .}}"github.com/aws/aws-sdk-go-v2/service/{{.}}"
{{end}})
type Client struct {
{{range .}}{{.}}conn *{{.}}.Client
{{end}}}

func  NewClient() *Client {
cfg, err := external.LoadDefaultAWSConfig()
if err != nil {
	panic("failed to load config, " + err.Error())
}

return &Client{
{{range .}}{{.}}conn: {{.}}.New(cfg),
{{end}}}
}
`))
