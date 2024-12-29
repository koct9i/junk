package main

import (
	"debug/buildinfo"
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	pkgpath "path"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

// isVersionElement reports whether s is a well-formed path version element:
// v2, v3, v10, etc, but not v0, v05, v1.
func isVersionElement(s string) bool {
	if len(s) < 2 || s[0] != 'v' || s[1] == '0' || s[1] == '1' && len(s) == 2 {
		return false
	}
	for i := 1; i < len(s); i++ {
		if s[i] < '0' || '9' < s[i] {
			return false
		}
	}
	return true
}

type Tool struct {
	Name    string
	Package string
	Version string
}

func listTools(source string) ([]Tool, error) {
	var tools []Tool

	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedModule,
		BuildFlags: []string{"-tags", "tools"},
	}

	pkgs, err := packages.Load(cfg, source)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		log.Println("Package", pkg.Name)
		for _, imp := range pkg.Imports {
			// log.Println("Import", imp.PkgPath)
			if imp.PkgPath == "golang.org/x/tools/go/packages" {
				continue
			}
			_, name := pkgpath.Split(imp.PkgPath)
			if isVersionElement(name) {
				_, name = pkgpath.Split(pkgpath.Dir(imp.PkgPath))
			}
			if imp.Module != nil {
				log.Println("Import", name, imp.Module.Version)
				tools = append(tools, Tool{
					Name:    name,
					Package: imp.PkgPath,
					Version: imp.Module.Version,
				})
			}
		}
	}
	return tools, nil
}

func getToolVersion(filename string) (string, error) {
	info, err := buildinfo.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return info.Main.Version, nil
}

func installTool(pkg, version, binDir string) error {
	cmd := exec.Command("go", "install", pkg+"@"+version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("Run", cmd)
	// log.Println("Env", cmd.Env)
	return cmd.Run()
}

func main() {
	binDir, err := filepath.Abs("bin")
	os.Unsetenv("GOROOT")
	os.Setenv("GOBIN", binDir)
	if err != nil {
		log.Fatal("Cannot get absolute path", err)
	}
	tools, err := listTools("./...")
	if err != nil {
		log.Fatal("Cannot list tools", err)
	}
	for _, tool := range tools {
		version, err := getToolVersion(filepath.Join(binDir, tool.Name))
		if err == nil && version == tool.Version {
			log.Println("Uptodate", tool.Name, tool.Version)
			continue
		}
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			log.Fatal("Cannot get tool version", err)
		}
		log.Println("Install", tool.Name, tool.Version)
		err = installTool(tool.Package, tool.Version, binDir)
		if err != nil {
			log.Fatalln("Cannot install tool", err)
		}
	}
}
