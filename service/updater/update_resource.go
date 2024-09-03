package updater

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	khttp "github.com/kkkkninezzz/maa-resource-updater/internal/http"
)

func downloadRemoteResource(maaResourceFileUrl string, outPath string) {
	client := khttp.Client()

	resp, err := client.Get(maaResourceFileUrl)
	if err != nil {
		log.Panicf("Error download remote resource: %v", err)
	}

	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error download remote resource, http status code: %v", resp.StatusCode)
	}

	// 创建本地文件
	out, err := os.Create(outPath)
	if err != nil {
		log.Panicf("Error create temp resource file: %v", err)
	}
	defer out.Close()

	// 复制响应体到文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Panicf("Error save temp resource file: %v", err)
	}

}

func unzipResource(zipFilePath, dest string) {

	// 打开 ZIP 文件
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		log.Panicf("Error open temp resource zip file: %v", err)
	}
	defer reader.Close()

	for _, f := range reader.File {
		// Construct the full path for the file
		path := filepath.Join(dest, f.Name)

		// Check for path traversal vulnerability
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			log.Panicf("Invalid file path: %s", path)
		}

		// Check if the file is a directory
		if f.FileInfo().IsDir() {
			// Create directories if they don't exist
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				log.Panicf("Error create directory: %s, %v", path, err)
			}
			continue
		}

		// Ensure the parent directory exists
		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			log.Panicf("Error create parent directory: %s, %v", path, err)
		}

		// Create the destination file
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Panicf("Error reate the destination file: %s, %v", path, err)
		}

		// Open the zipped file
		rc, err := f.Open()
		if err != nil {
			log.Panicf("Error Open the zipped file: %v", err)
		}

		// Copy the file content to the destination file
		_, err = io.Copy(outFile, rc)

		// Close the files
		outFile.Close()
		rc.Close()

		if err != nil {
			log.Panicf("Error copy file: %s, %v", path, err)
		}
	}

}

func generateTimestampedFilename(unzipDir string) string {
	// 获取当前时间
	now := time.Now()

	// 格式化时间为字符串，例如："20240903_224300"
	timestamp := now.Format("20060102_150405")

	// 创建带时间戳的文件名
	filename := fmt.Sprintf("MaaResourceTemp_%s.zip", timestamp)

	return filepath.Join(unzipDir, filename)
}

func UpdateResource(maaResourceFileUrl string, unzipDir string) {

	tempPath := generateTimestampedFilename(unzipDir)
	downloadRemoteResource(maaResourceFileUrl, tempPath)

	log.Println("Start unzip resource file")
	unzipResource(tempPath, unzipDir)

	log.Println("Start move resource file")
	moveResourceFiles(unzipDir)

	// 删除文件
	err := os.Remove(tempPath)
	if err != nil {
		log.Panicf("Error remove temp resouce file: %v", err)
	}
}

// 需要将解压的文件移动到指定目录
func moveResourceFiles(unzipDir string) {
	unzipedMaaPath := filepath.Join(unzipDir, "MaaResource-main")

	unzipedMaaResourcePath := filepath.Join(unzipedMaaPath, "resource")
	maaResoucePath := filepath.Join(unzipDir, "resource")
	err := moveDir(unzipedMaaResourcePath, maaResoucePath)

	if err != nil {
		log.Panicf("Error move temp resouce file: %v", err)
	}

	unzipedMaaCachePath := filepath.Join(unzipedMaaPath, "cache")
	maaCachePath := filepath.Join(unzipDir, "cache")
	err = moveDir(unzipedMaaCachePath, maaCachePath)

	if err != nil {
		log.Panicf("Error move temp cache file: %v", err)
	}

	err = os.RemoveAll(unzipedMaaPath)
	if err != nil {
		log.Panicf("Error remove MaaResource-main dir: %v", err)
	}
}

func moveDir(srcDir, destDir string) error {
	// Ensure the destination directory exists
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Walk through the source directory
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the source directory itself
		if path == srcDir {
			return nil
		}

		// Calculate the destination path for the current item
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, relPath)

		// If it's a directory, create the directory in the destination
		if info.IsDir() {
			err := os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			// If it's a file, move it to the destination
			err := os.Rename(path, destPath)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Optionally, you can delete the source directory after moving all files
	return os.RemoveAll(srcDir)
}
