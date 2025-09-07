package ffr

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Run はCLIエントリポイント。戻り値はプロセス終了コード。
func Run(args []string) int {
	config := &ReplacementConfig{
		TargetPath:     "",
		IsFile:         false,
		ProcessedFiles: make([]string, 0),
		FileCount:      0,
	}

	if len(args) < 2 {
		fmt.Println("エラー: ファイルまたはディレクトリを指定してください。")
		printUsage()
		return 1
	}

	arg := args[1]
	if arg == "--help" || arg == "-h" {
		printUsage()
		return 0
	}

	if arg == "--version" || arg == "-v" {
		printVersion()
		return 0
	}

	config.TargetPath = arg

	info, err := os.Stat(config.TargetPath)
	if err != nil {
		log.Printf("エラー: '%s' が存在しません。", config.TargetPath)
		return 1
	}

	config.IsFile = !info.IsDir()

	if config.IsFile {
		if !strings.HasSuffix(config.TargetPath, ".blade.php") {
			log.Printf("エラー: '%s' は.blade.phpファイルではありません。", config.TargetPath)
			return 1
		}
		fmt.Printf("Form Facade置換を開始します (ファイル): %s\n", config.TargetPath)
	} else {
		fmt.Printf("Form Facade置換を開始します (ディレクトリ): %s\n", config.TargetPath)
	}

	err = processBladeFiles(config)
	if err != nil {
		log.Printf("ファイル処理中にエラーが発生しました: %v", err)
		return 1
	}

	printSummary(config)
	return 0
}
