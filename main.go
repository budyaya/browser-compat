package main

import (
	"browser-compat/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gorm.io/datatypes"
)

var (
	Verbose bool
	Dir     string
	rootCmd = &cobra.Command{
		Use:          "parse",
		Short:        "Start parse browser compat",
		Example:      "browser-compat parse -d ../browser-compat-data",
		SilenceUsage: false,
		Args:         cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			RunCmd()
		},
	}
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&Dir, "dir", "d", "../browser-compat-data", "The browser-compat-data project path")
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func RunCmd() {
	log.Println("dir", Dir)
	Dir, _ = filepath.Abs(Dir)
	// 遍历当前目录下的所有文件和子目录
	filepath.Walk(Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Panicln(err)
			return err
		}
		path, _ = filepath.Abs(path)
		subpath := strings.TrimPrefix(path, Dir)
		//log.Println("subpath", subpath)
		if info.IsDir() {
			// 忽略隐藏目录和文件
			if info.Name()[0] == '.' {
				return filepath.SkipDir
			}
			if subpath == "browsers" {
				return filepath.SkipDir
			}
			if subpath == "docs" {
				return filepath.SkipDir
			}
			if subpath == "lint" {
				return filepath.SkipDir
			}
			if subpath == "release_notes" {
				return filepath.SkipDir
			}
			if subpath == "schemas" {
				return filepath.SkipDir
			}
			if subpath == "scripts" {
				return filepath.SkipDir
			}
			if subpath == "svg" {
				return filepath.SkipDir
			}
			if subpath == "types" {
				return filepath.SkipDir
			}
			if subpath == "utils" {
				return filepath.SkipDir
			}
		}
		if info.Name()[0] == '.' {
			return nil
		}
		// 检查文件是否为JSON文件
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			//log.Println("Processing file:", path)
			// 读取JSON文件内容
			jsonstr, err := os.ReadFile(path)
			if err != nil {
				log.Panicln(err)
				return err
			}
			// 解析JSON数据
			parseJson(path, []string{}, jsonstr)
		}
		return nil
	})
}

// 递归处理 json
// 遍历 json 数据
// 解析 json 数据
func parseJson(path string, keys []string, jsonstr []byte) error {
	api := strings.Join(keys, ".")
	//log.Println("api:", api)
	// 解析JSON数据
	var result map[string]datatypes.JSON
	err := json.Unmarshal(jsonstr, &result)
	if err != nil {
		log.Println(string(jsonstr), err)
		return err
	}
	datastr, ok := result["__compat"]
	if ok && len(keys) > 0 {
		var data map[string]datatypes.JSON
		err := json.Unmarshal(datastr, &data)
		if err != nil {
			log.Println(string(datastr), err)
			return err
		}
		compat := models.BrowserCompatData{
			Type:    keys[0],
			MdnUrl:  string(data["mdn_url"]),
			SpecUrl: string(data["spec_url"]),
			Api:     api,
		}
		for dk, dv := range data {
			switch dk {
			case "tags":
				json.Unmarshal(data["tags"], &compat.Tags)
			case "status":
				json.Unmarshal(data["status"], &compat.Status)
			case "support":
				var supports map[string]datatypes.JSON
				err = json.Unmarshal(dv, &supports)
				if err != nil {
					log.Println(string(dv), err)
					return err
				}
				for sk, sv := range supports {

					var version map[string]interface{}
					err = json.Unmarshal(sv, &version)
					if err != nil {
						var versions []map[string]interface{}
						err = json.Unmarshal(sv, &versions)
						if err != nil {
							compat.Browser = sk
							compat.BrowserVersion = string(sv)
							err = models.Create(compat)
							if err != nil {
								log.Println(api, err)
								return err
							}
						}
						for _, version = range versions {
							compat.Browser = sk
							compat.BrowserVersion = fmt.Sprintf("%v", version["version_added"])
							err = models.Create(compat)
							if err != nil {
								log.Println(api, err)
								return err
							}
						}
					} else {
						compat.Browser = sk
						compat.BrowserVersion = fmt.Sprintf("%v", version["version_added"])
						err = models.Create(compat)
						if err != nil {
							log.Println(api, err)
							return err
						}
					}
				}
			}
		}

		return nil
	}
	for k, v := range result {
		keys = append(keys, k)
		err := parseJson(path, keys, v)
		if err != nil {
			log.Println(path, err)
			return err
		}
	}
	return nil
}
