package cmd

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

var nativeTemplateCmd = &cobra.Command{
	Use:   "native-template <project_name>",
	Short: "Create a fully template with go native",
	Long:  "Create a fully template with go native",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := createNativeTemplate(args); err != nil {
			log.Fatalf("❌ Error occured: %v", err)
		}
	},
}

func init() {
	createCmd.AddCommand(nativeTemplateCmd)
}

func generateFolders(projectName string, err error) (error, error, bool) {
	if err = os.Mkdir(projectName, 0755); err != nil {
		return nil, err, true
	}

	dirs := []string{
		"cmd",
		"internal/config",
		"internal/database",
		"internal/models",
		"internal/repository",
		"internal/service",
	}

	for _, dir := range dirs {
		path := filepath.Join(projectName, dir)
		if err = os.MkdirAll(path, 0755); err != nil {
			return nil, err, true
		}
	}
	return err, nil, false
}

func promptDriverDb() (string, error) {
	var driver = "N/A"

	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("Please select your database driver:").Options(dbOptions...).Value(&driver)))

	err := form.Run()

	if err != nil {
		return "", err
	}

	return driver, nil
}

func renderTemplate(content string, data TemplateData, err error) (string, error) {
	tmpl, err := template.New("file").Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func generateFiles(projectName string, selectedDb string, err error) (error, bool) {
	withoutDriver := selectedDb == "N/A"

	data := TemplateData{ProjectName: projectName, DatabaseDriver: selectedDb}

	files := map[string]string{
		"cmd/main.go":               setMainTemplate(withoutDriver),
		"internal/config/config.go": setCfgTemplate(withoutDriver),
		"internal/database/db.go":   setDbTemplate(withoutDriver),
		".env":                      setEnvTemplate(withoutDriver),
	}

	for path, content := range files {
		fullPath := filepath.Join(projectName, path)
		finalContent, err := renderTemplate(content, data, err)

		if err = os.WriteFile(fullPath, []byte(finalContent), 0644); err != nil {
			return err, false
		}
	}
	return nil, true
}

func createNativeTemplate(args []string) (err error) {
	projectName := args[0]

	defer func() {
		if err != nil {
			fmt.Printf("⚠️ Error!, rollback project '%s'...\n", projectName)
			_ = os.RemoveAll(projectName)
		}
	}()

	fmt.Printf("🚀 Creating project: %s\n", projectName)

	err, err2, done := generateFolders(projectName, err)
	if done {
		return err2
	}

	selectedDb, err := promptDriverDb()
	if err != nil {
		return err
	}

	err, _ = generateFiles(projectName, selectedDb, err)
	if err != nil {
		return err
	}

	gmCmd := exec.Command("go", "mod", "init", projectName)
	gmCmd.Dir = projectName

	err = gmCmd.Run()
	if err != nil {
		return err
	}

	plugins := getPlugins(selectedDb)

	for _, pluginPath := range plugins {
		gmCmd = exec.Command("go", "get", pluginPath)
		gmCmd.Dir = projectName

		err = gmCmd.Run()
		if err != nil {
			return err
		}
	}

	fmt.Printf("'%s' Created successfully, Happy coding 🎉 !", projectName)

	return nil
}
