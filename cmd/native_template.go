package cmd

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"gong/tui"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type projectConfig struct {
	Name      string
	DBDriver  string
	GoModPath string
}

var nativeTemplateCmd = &cobra.Command{
	Use:   "native-template <project_name>",
	Short: "Create a new project using the native Go template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := createNativeTemplate(args); err != nil {
			log.Fatalf("❌ Error: %v", err)
		}
	},
}

func init() {
	createCmd.AddCommand(nativeTemplateCmd)
}

func generateFolders(cfg projectConfig) error {
	dirs := []string{
		"cmd", "internal/config", "internal/database",
		"internal/models", "internal/repository", "internal/service",
	}
	for _, dir := range dirs {
		path := filepath.Join(cfg.Name, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}
	return nil
}

func generateFiles(cfg projectConfig) error {
	withDriver := cfg.DBDriver != "N/A"
	data := TemplateData{ProjectName: cfg.Name, DatabaseDriver: cfg.DBDriver}
	files := map[string]string{
		"cmd/main.go":               setMainTemplate(withDriver),
		"internal/config/config.go": setCfgTemplate(withDriver),
		"internal/database/db.go":   setDbTemplate(withDriver),
		".env":                      setEnvTemplate(withDriver),
	}
	for path, content := range files {
		fullPath := filepath.Join(cfg.Name, path)
		finalContent, err := renderTemplate(content, data)
		if err != nil {
			return fmt.Errorf("failed to render template %s: %w", path, err)
		}
		if err := os.WriteFile(fullPath, []byte(finalContent), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", path, err)
		}
	}
	return nil
}

func runGoCommands(cfg projectConfig) error {
	goModInitTask := func() error {
		cmd := exec.Command("go", "mod", "init", cfg.GoModPath)
		cmd.Dir = cfg.Name
		return cmd.Run()
	}

	if err := tui.RunSpinner("Initializing Go module...", goModInitTask); err != nil {
		return err
	}

	plugins := getPlugins(cfg.DBDriver)
	if len(plugins) > 0 {
		downloadTask := func() error {
			for _, pluginPath := range plugins {
				cmd := exec.Command("go", "get", pluginPath)
				cmd.Dir = cfg.Name
				if err := cmd.Run(); err != nil {
					return err
				}
			}
			return nil
		}
		if err := tui.RunSpinner(fmt.Sprintf("Downloading %s driver...", cfg.DBDriver), downloadTask); err != nil {
			return err
		}
	}
	return nil
}

func promptForDBDriver() (string, error) {
	var driver = dbOptions[0].Value
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a database driver to use:").
				Options(dbOptions...).
				Value(&driver),
		),
	)
	err := form.Run()
	if err != nil {
		return "", err
	}
	return driver, nil
}

func renderTemplate(content string, data TemplateData) (string, error) {
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

func createNativeTemplate(args []string) (err error) {
	projectName := args[0]

	defer func() {
		if err != nil {
			fmt.Printf("\n⚠️ An error occurred, rolling back changes for project '%s'...\n", projectName)
			_ = os.RemoveAll(projectName)
		}
	}()

	fmt.Printf("🚀 Creating new project: %s\n", projectName)

	selectedDb, err := promptForDBDriver()
	if err != nil {
		return err
	}

	cfg := projectConfig{
		Name:      projectName,
		DBDriver:  selectedDb,
		GoModPath: projectName,
	}

	if err := generateFolders(cfg); err != nil {
		return err
	}

	if err := generateFiles(cfg); err != nil {
		return err
	}

	if err := runGoCommands(cfg); err != nil {
		return err
	}

	fmt.Printf("\n✅ Project '%s' created successfully. Happy coding! 🎉\n", projectName)
	return nil
}
