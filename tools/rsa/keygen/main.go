package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mf-financial/cayenne-wallet-util/storage"

	"github.com/mf-financial/cayenne-wallet-util/rsa"
)

const (
	privateFileName = "private.pem"
	publicFileName  = "public.pem"
)

var (
	parentFolderPath string
	storageType      string
	bucketName       string
)

// RSAの鍵生成 のためのコマンド。
// 複数回実行すると鍵が上書きされるので要注意
// Usage of ./keygen:
//   -b string
//         bucket name
//   -o string
//         parent folder to output (default ".")
//   -s string
//         storage type (default "local")

func main() {
	parseFlag()

	privateKey, err := rsa.GeneratePrivateKey()
	if err != nil {
		fmt.Printf("GeneratePrivateKey Error: %v", err)
		return
	}
	publicKey := rsa.GeneratePublickey(privateKey)
	privateBytes := rsa.EncodePrivateKeyToPem(privateKey)
	publicBytes := rsa.EncodePublicKeyToPem(publicKey)

	uploader, err := getUploader()
	if err != nil {
		fmt.Printf("getUploader Error: %v", err)
		return
	}
	// public keyとprivate keyがあるので親階層のフォルダまで
	privataKeyPath, publicKeyPath, err := getUploadPath()
	if err != nil {
		fmt.Printf("getUploadPath Error: %v", err)
		return
	}

	if err := uploader.Upload(privataKeyPath, privateBytes); err != nil {
		fmt.Printf("Upload PrivateKey Error: %v", err)
		return
	}
	if err := uploader.Upload(publicKeyPath, publicBytes); err != nil {
		fmt.Printf("Upload PrivateKey Error: %v", err)
		return
	}
}

func parseFlag() {
	flag.StringVar(&parentFolderPath, "o", ".", "parent folder to output")
	flag.StringVar(&storageType, "s", "local", "storage type")
	flag.StringVar(&bucketName, "b", "", "bucket name")
	flag.Parse()
}

func getUploader() (storage.Uploader, error) {
	switch strings.ToLower(storageType) {
	case "gcs":
		return storage.NewGCSUploader(context.Background(), bucketName)
	case "local":
		return storage.NewLocalUploader(), nil
	default:
		return nil, errors.New("unrecognized storage type")
	}
}

func getUploadPath() (privatePath string, publicPath string, err error) {
	if parentFolderPath == "" {
		return privateFileName, publicFileName, nil
	}
	switch strings.ToLower(storageType) {
	case "gcs":
		return fmt.Sprintf("%s/%s", parentFolderPath, privateFileName), fmt.Sprintf("%s/%s", parentFolderPath, publicFileName), nil
	case "local":
		privateFullPath, err := filepath.Abs(fmt.Sprintf("%s/%s", parentFolderPath, privateFileName))
		if err != nil {
			return "", "", err
		}
		publicFullPath, err := filepath.Abs(fmt.Sprintf("%s/%s", parentFolderPath, publicFileName))
		if err != nil {
			return "", "", err
		}
		return privateFullPath, publicFullPath, nil

	default:
		return "", "", nil
	}
}