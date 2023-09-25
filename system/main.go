package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: generator <packagename> <filename>")
		return
	}
	packagename := os.Args[1]
	filename := os.Args[2]

	// Call the GenerateModelAndController function from the generate package
	if err := GenerateModelAndController(filename, packagename); err != nil {
		fmt.Println("Error generating files:", err)
	}
}

// GenerateModelAndController generates model and controller files.
func GenerateModelAndController(filename string, packagename string) error {
	modelDir := "models"
	controllerDir := "controllers"
	var modelPath string
	var controllerPath string
	if packagename == "both" || packagename == "models" {
		modelPath = filepath.Join(modelDir, filename+".go")
		if err := createFile(modelPath, "models"); err != nil {
			return err
		}
		fmt.Printf("Created model file: %s\n", modelPath)
	}
	if packagename == "both" || packagename == "controllers" {
		controllerPath = filepath.Join(controllerDir, filename+".go")
		if err := createFile(controllerPath, "controllers"); err != nil {
			return err
		}
		fmt.Printf("Created controller file: %s\n", controllerPath)
	}

	return nil
}

func createFile(filePath, typeName string) error {
	if _, err := ioutil.ReadFile(filePath); err == nil {
		return fmt.Errorf("file already exists: %s", filePath)
	}

	code := generateCode(typeName)
	return ioutil.WriteFile(filePath, []byte(code), 0644)
}

func generateCode(typeName string) string {
	if typeName == "models" {
		return fmt.Sprintf("package %s\n\ntype DefaultStructure struct {\n    //change struct name and  Define your fields here\n}\n\nfunc DefaultFunction(){\n // change function name and start writting code. \n //HAPPY CODING \n}\n", typeName)
	} else if typeName == "controllers" {
		return fmt.Sprintf("package %s\nimport \"net/http\"\n\nfunc DefaultFunction(w http.ResponseWriter, r *http.Request){\n // change function name and start writting code. \n //HAPPY CODING \n}\n", typeName)
	}
	return fmt.Sprintf("package %s\n\ntype %s struct {\n    // Define your fields here\n}\n", typeName, typeName)
}
