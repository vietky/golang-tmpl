package templatewriter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/vietky/golang-tmpl/libraries/logger"
)

var log = logger.GetLogger("templatewriter")

type TemplateWriter struct {
	SourceFolder string
	templates    *template.Template
}

func NewTemplateWriter(source string) (*TemplateWriter, error) {
	tw := &TemplateWriter{
		SourceFolder: source,
	}
	err := tw.init()
	if err != nil {
		return nil, err
	}
	return tw, err
}

func (tw *TemplateWriter) init() error {
	fileList, err := tw.listAllFiles(tw.SourceFolder)
	if err != nil {
		return err
	}
	t := template.New("templatewriter")
	tpl, err := t.ParseFiles(fileList...)
	if err != nil {
		return err
	}
	tw.templates = tpl
	return nil
}

func (*TemplateWriter) listAllFiles(searchDir string) ([]string, error) {
	fileList := make([]string, 0)
	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return err
		}
		fileList = append(fileList, path)
		return err
	})

	if e != nil {
		return nil, e
	}

	return fileList, nil
}

func (tw *TemplateWriter) Write(targetFolder string, data map[string]string) (err error) {
	if _, err = os.Stat(targetFolder); os.IsNotExist(err) {
		err = os.Mkdir(targetFolder, 0777)
		if err != nil {
			return err
		}
	}
	targetFolder = strings.Trim(targetFolder, "/")

	for _, tpl := range tw.templates.Templates() {

		// arr := strings.Split(tpl, "\/")
		// templateName = arr[len(arr)-1]
		templateName := tpl.Name()
		generatedFilePath := fmt.Sprintf("%s/%s", targetFolder, templateName)
		// fmt.Println("templateName", templateName, generatedFilePath)
		w, err := os.Create(generatedFilePath)
		if err != nil {
			return err
		}
		err = tw.templates.ExecuteTemplate(w, templateName, data)
		if err != nil {
			return err
		}
	}
	return err
}
