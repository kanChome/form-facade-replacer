// pipeline.go: ファイル/ディレクトリ走査と1ファイル処理のパイプライン。
package ffr

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type ReplacementConfig struct {
	TargetPath     string
	IsFile         bool
	ProcessedFiles []string
	FileCount      int
}

// processBladeFiles はディレクトリ（または単一ファイル）を走査して置換処理を行う。
func processBladeFiles(config *ReplacementConfig) error {
	if config.IsFile {
		return processSingleFile(config, config.TargetPath)
	}

	return filepath.WalkDir(config.TargetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".blade.php") {
			return processSingleFile(config, path)
		}
		return nil
	})
}

// processSingleFile は1ファイルの置換と進捗集計を行う。
func processSingleFile(config *ReplacementConfig, filePath string) error {
	hasFormFacade, err := containsFormFacade(filePath)
	if err != nil {
		return fmt.Errorf("ファイル %s のチェックに失敗しました: %v", filePath, err)
	}

	if hasFormFacade {
		fmt.Printf("処理中: %s\n", filePath)
		err := replaceFormPatterns(filePath)
		if err != nil {
			return fmt.Errorf("ファイル %s の処理に失敗しました: %v", filePath, err)
		}
		config.ProcessedFiles = append(config.ProcessedFiles, filePath)
		config.FileCount++
		fmt.Printf(" - 処理完了: %s\n", filePath)
	}
	return nil
}

// containsFormFacade はファイル内に "Form::" が存在するかを高速に判定する。
func containsFormFacade(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Form::") {
			return true, nil
		}
	}
	return false, scanner.Err()
}
