// Package files - ファイル操作に関するパッケージ
package files

import "os"

// Exists - ファイルの存在チェックをする
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
