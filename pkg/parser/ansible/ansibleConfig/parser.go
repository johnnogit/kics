package ansibleConfig

import (
	"strconv"
	"strings"

	"github.com/Checkmarx/kics/pkg/model"
	"github.com/bigkevmcd/go-configparser"
)

// Parser defines a parser type
type Parser struct {
}

func (p *Parser) Resolve(fileContent []byte, filename string) ([]byte, error) {
	return fileContent, nil
}

// Parse parses .cfg/.conf file and returns it as a Document
func (p *Parser) Parse(filePath string, fileContent []byte) ([]model.Document, []int, error) {
	model.NewIgnore.Reset()

	reader := strings.NewReader(string(fileContent))
	configparser.Delimiters("=")
	config, err := configparser.ParseReader(reader)
	if err != nil {
		return nil, nil, err
	}

	doc := refactorConfig(config)

	return []model.Document{*doc}, []int{}, nil
}

// refactorConfig removes all extra information
func refactorConfig(config *configparser.ConfigParser) (doc *model.Document) {
	doc = emptyDocument()
	for _, section := range config.Sections() {
		dict, err := config.Items(section)
		if err != nil {
			continue
		}
		dictRefact := make(map[string]interface{})
		for key, value := range dict {
			if intValue, err := strconv.Atoi(value); err == nil {
				dictRefact[key] = intValue
				continue
			}
			dictRefact[key] = value
		}
		(*doc)[section] = dictRefact
	}
	return doc
}

// SupportedExtensions returns extensions supported by this parser, which are only ini extension
func (p *Parser) SupportedExtensions() []string {
	return []string{".cfg", ".conf"}
}

// SupportedTypes returns types supported by this parser, which is ansible
func (p *Parser) SupportedTypes() map[string]bool {
	return map[string]bool{
		"ansible": true,
	}
}

// GetKind returns CFG constant kind
func (p *Parser) GetKind() model.FileKind {
	return model.KindCFG
}

// GetCommentToken return the comment token of CFG/CONF - #
func (p *Parser) GetCommentToken() string {
	return "#"
}

// GetResolvedFiles returns resolved files
func (p *Parser) GetResolvedFiles() map[string]model.ResolvedFile {
	return make(map[string]model.ResolvedFile)
}

// StringifyContent converts original content into string formatted version
func (p *Parser) StringifyContent(content []byte) (string, error) {
	return string(content), nil
}

func emptyDocument() *model.Document {
	return &model.Document{}
}
