package onboarding

import (
	"io/fs"
	"os"
	"path/filepath"
)

var lockfiles = []struct {
	file string
	cmd  []string
}{
	{"package-lock.json", []string{"npm", "ci"}},
	{"pnpm-lock.yaml", []string{"pnpm", "install", "--frozen-lockfile"}},
	{"yarn.lock", []string{"yarn", "install", "--frozen-lockfile"}},
	{"bun.lock", []string{"bun", "install", "--frozen-lockfile"}},
	{"bun.lockb", []string{"bun", "install", "--frozen-lockfile"}},
}

// NPMTask detects Node.js projects with lockfiles.
type NPMTask struct{}

func (NPMTask) Name() string { return "npm" }

func (NPMTask) Detect(projectDir string) []Workload {
	var workloads []Workload
	for _, dir := range findPackageJSONDirs(projectDir) {
		for _, lf := range lockfiles {
			lfPath := filepath.Join(dir, lf.file)
			if _, err := os.Stat(lfPath); err == nil {
				rel, _ := filepath.Rel(projectDir, dir)
				if rel == "." {
					rel = filepath.Base(projectDir)
				}
				workloads = append(workloads, Workload{
					Label:      rel,
					Command:    lf.cmd,
					Dir:        dir,
					HashInputs: []string{lfPath},
					Phase:      PostContainer,
				})
				break
			}
		}
	}
	return workloads
}

// findPackageJSONDirs walks projectDir and returns directories containing
// package.json files (skipping node_modules and other non-relevant dirs).
func findPackageJSONDirs(projectDir string) []string {
	skip := map[string]bool{
		".git": true, ".venv": true, "__pycache__": true,
		"vendor": true, "target": true, "dist": true,
	}
	var results []string
	filepath.WalkDir(projectDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		name := d.Name()
		if name == "node_modules" {
			return filepath.SkipDir
		}
		if path != projectDir && skip[name] {
			return filepath.SkipDir
		}
		if _, err := os.Stat(filepath.Join(path, "package.json")); err == nil {
			results = append(results, path)
		}
		return nil
	})
	return results
}
