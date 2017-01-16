package files_test

import (
	"fmt"
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
		fmt.Println(temps.Name())
		if temps.Name() == "templateBottom" {
			return // We pass, as it was parsed
		}
	}

	t.Errorf("Template was not parsed")
}
