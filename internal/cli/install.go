package cli

import (
	"fmt"
	"path/filepath"

	"github.com/answerdev/answer/configs"
	"github.com/answerdev/answer/i18n"
	"github.com/answerdev/answer/pkg/dir"
	"github.com/answerdev/answer/pkg/writer"
)

const (
	DefaultConfigFileName = "config.yaml"
	DefaultCacheFileName  = "cache.db"
)

var (
	ConfigFilePath = "/conf/"
	UploadFilePath = "/upfiles/"
	I18nPath       = "/i18n/"
	CachePath      = "/cache/"
)

// InstallAllInitialEnvironment install all initial environment
func InstallAllInitialEnvironment(dataDirPath string) {
	ConfigFilePath = filepath.Join(dataDirPath, ConfigFilePath)
	UploadFilePath = filepath.Join(dataDirPath, UploadFilePath)
	I18nPath = filepath.Join(dataDirPath, I18nPath)
	CachePath = filepath.Join(dataDirPath, CachePath)

	installUploadDir()
	installI18nBundle()
	fmt.Println("install all initial environment done")
}

func InstallConfigFile(configFilePath string) error {
	if len(configFilePath) == 0 {
		configFilePath = filepath.Join(ConfigFilePath, DefaultConfigFileName)
	}
	fmt.Println("[config-file] try to create at ", configFilePath)

	// if config file already exists do nothing.
	if CheckConfigFile(configFilePath) {
		fmt.Printf("[config-file] %s already exists\n", configFilePath)
		return nil
	}

	if err := dir.CreateDirIfNotExist(ConfigFilePath); err != nil {
		fmt.Printf("[config-file] create directory fail %s\n", err.Error())
		return fmt.Errorf("create directory fail %s", err.Error())
	}
	fmt.Printf("[config-file] create directory success, config file is %s\n", configFilePath)

	if err := writer.WriteFile(configFilePath, string(configs.Config)); err != nil {
		fmt.Printf("[config-file] install fail %s\n", err.Error())
		return fmt.Errorf("write file failed %s", err)
	}
	fmt.Printf("[config-file] install success\n")
	return nil
}

func installUploadDir() {
	fmt.Println("[upload-dir] try to install...")
	if err := dir.CreateDirIfNotExist(UploadFilePath); err != nil {
		fmt.Printf("[upload-dir] install fail %s\n", err.Error())
	} else {
		fmt.Printf("[upload-dir] install success, upload directory is %s\n", UploadFilePath)
	}
}

func installI18nBundle() {
	fmt.Println("[i18n] try to install i18n bundle...")
	if err := dir.CreateDirIfNotExist(I18nPath); err != nil {
		fmt.Println(err.Error())
		return
	}

	i18nList, err := i18n.I18n.ReadDir(".")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("[i18n] find i18n bundle %d\n", len(i18nList))
	for _, item := range i18nList {
		path := filepath.Join(I18nPath, item.Name())
		content, err := i18n.I18n.ReadFile(item.Name())
		if err != nil {
			continue
		}
		fmt.Printf("[i18n] install %s bundle...\n", item.Name())
		err = writer.WriteFile(path, string(content))
		if err != nil {
			fmt.Printf("[i18n] install %s bundle fail: %s\n", item.Name(), err.Error())
		} else {
			fmt.Printf("[i18n] install %s bundle success\n", item.Name())
		}
	}
}
