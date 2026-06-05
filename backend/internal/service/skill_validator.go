package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	MaxUnzipSize int64 = 10 * 1024 * 1024 // 10MB
	MaxFileCount int   = 100
	SkillMDFile        = "SKILL.md"
	// Deprecated: 仅用于迁移过渡，迁移完成后移除
	HubBasePath = "../data/skills/hub"
)

type ValidationResult struct {
	Valid       bool     `json:"valid"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	FileCount   int      `json:"file_count,omitempty"`
	TotalSize   int64    `json:"total_size,omitempty"`
	TmpDir      string   `json:"tmp_dir,omitempty"`
	Errors      []string `json:"errors,omitempty"`
}

type frontmatterData struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// ValidateZip 解压到临时目录并校验
func ValidateZip(zipData []byte) (*ValidationResult, string, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return &ValidationResult{Valid: false, Errors: []string{"invalid zip file"}}, "", err
	}

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "skill-upload-*")
	if err != nil {
		return nil, "", fmt.Errorf("create temp dir: %w", err)
	}

	var (
		totalSize  int64
		fileCount  int
		errors     []string
		hasSkillMD bool
		fmName     string
		fmDesc     string
	)

	for _, f := range reader.File {
		// 路径安全检查
		if strings.Contains(f.Name, "..") {
			errors = append(errors, "path traversal detected")
			continue
		}
		if filepath.IsAbs(f.Name) {
			errors = append(errors, "absolute path not allowed")
			continue
		}

		// 符号链接检查
		if f.Mode()&os.ModeSymlink != 0 {
			errors = append(errors, "symbolic links not allowed")
			continue
		}

		destPath := filepath.Join(tmpDir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(destPath, 0755)
			continue
		}

		// 确保父目录存在
		os.MkdirAll(filepath.Dir(destPath), 0755)

		rc, err := f.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("cannot open %s: %v", f.Name, err))
			continue
		}

		var buf bytes.Buffer
		n, err := io.Copy(&buf, rc)
		rc.Close()
		if err != nil {
			errors = append(errors, fmt.Sprintf("cannot read %s: %v", f.Name, err))
			continue
		}

		totalSize += n
		fileCount++

		// Zip bomb 检查
		if totalSize > MaxUnzipSize {
			errors = append(errors, fmt.Sprintf("zip bomb: total size exceeds %dMB", MaxUnzipSize/1024/1024))
			os.RemoveAll(tmpDir)
			return &ValidationResult{Valid: false, Errors: errors}, tmpDir, nil
		}
		if fileCount > MaxFileCount {
			errors = append(errors, fmt.Sprintf("too many files: exceeds %d", MaxFileCount))
			os.RemoveAll(tmpDir)
			return &ValidationResult{Valid: false, Errors: errors}, tmpDir, nil
		}

		// 写入文件
		if err := os.WriteFile(destPath, buf.Bytes(), 0644); err != nil {
			errors = append(errors, fmt.Sprintf("cannot write %s: %v", f.Name, err))
			continue
		}

		// 检查 SKILL.md（根目录或一级子目录，如 skill-name/SKILL.md）
		if filepath.Base(f.Name) == SkillMDFile && !strings.Contains(filepath.Dir(f.Name), "/") {
			hasSkillMD = true
			name, desc, parseErr := parseFrontmatter(buf.Bytes())
			if parseErr != nil {
				errors = append(errors, parseErr.Error())
			} else {
				fmName = name
				fmDesc = desc
			}
		}
	}

	if !hasSkillMD {
		errors = append(errors, "missing SKILL.md")
	}

	if len(errors) > 0 {
		return &ValidationResult{Valid: false, Errors: errors}, tmpDir, nil
	}

	return &ValidationResult{
		Valid:       true,
		Name:        fmName,
		Description: fmDesc,
		FileCount:   fileCount,
		TotalSize:   totalSize,
		TmpDir:      tmpDir,
	}, tmpDir, nil
}

// parseFrontmatter 解析 SKILL.md 的 YAML frontmatter
func parseFrontmatter(data []byte) (name, description string, err error) {
	content := string(data)
	if !strings.HasPrefix(content, "---") {
		return "", "", fmt.Errorf("missing frontmatter")
	}

	end := strings.Index(content[3:], "---")
	if end == -1 {
		return "", "", fmt.Errorf("missing frontmatter")
	}

	fmContent := strings.TrimSpace(content[3 : end+3])
	var fm frontmatterData
	if err := yaml.Unmarshal([]byte(fmContent), &fm); err != nil {
		return "", "", fmt.Errorf("invalid frontmatter: %v", err)
	}

	if fm.Name == "" {
		return "", "", fmt.Errorf("missing name field")
	}

	return fm.Name, fm.Description, nil
}

// PackValidatedSkillDir repacks the validated skill directory into a zip blob.
func PackValidatedSkillDir(name string, tmpDir string) ([]byte, error) {
	// 确定源目录（zip 可能有或没有外层目录）
	srcDir := filepath.Join(tmpDir, name)
	if info, err := os.Stat(srcDir); err != nil || !info.IsDir() {
		srcDir = tmpDir
	}

	// 将已校验的文件重新打包为 zip
	zipData, err := zipDir(srcDir)
	if err != nil {
		return nil, fmt.Errorf("pack skill files: %w", err)
	}

	return zipData, nil
}

// zipDir 将目录打包为 zip 字节流
func zipDir(src string) ([]byte, error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error at %s: %w", path, err)
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("relative path error: %w", err)
		}
		rel = filepath.ToSlash(rel) // 统一正斜杠，跨平台一致
		f, err := w.Create(rel)
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		_, err = io.Copy(f, in)
		return err
	})
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
