package files_test

import (
	"testing"
	"text/template"

	. "github.com/FactomProject/enterprise-wallet/web/files"
)

func TestCustomParseGlob(t *testing.T) {
	temp := template.New("TestTemplate")
	temp = CustomParseGlob(temp, "templates/templateBottom.html")
	for _, temps := range temp.Templates() {
		if temps.Name() == "templateBottom" {
			return // We pass, as it was parsed
		}
	}

	t.Errorf("Template was not parsed")
}

func TestCustomParseFiles(t *testing.T) {
	var err error
	temp := template.New("TestTemplate")
	temp, err = CustomParseFile(temp, "templates/templateBottom.html")
	if err != nil {
		t.Fail()
	}
	for _, temps := range temp.Templates() {
		if temps.Name() == "templateBottom" {
			return // We pass, as it was parsed
		}
	}

	_, err = CustomParseFile(temp, "this-do-no-exist")
	if err == nil {
		t.Fail()
	}

	t.Errorf("Template was not parsed")
}
