package utils

import (
	"io/fs"
	"os"
	"server/global"

	"gopkg.in/yaml.v3"
)

const configFile = "config.yaml"

// LoadYAML 从文件中读取 YAML 数据并返回字节数组
func LoadYAML() ([]byte, error) {
	// ioutil.ReadFile(configFile)
	return os.ReadFile(configFile) // 我自己再用viper配置一遍到时候
}

// SaveYAML 将全局配置对象保存为 YAML 格式到文件,并不需要再进行一个热部署，将yaml文件配置再重新读取到config里面
// 然后我感觉热部署其实是适用于在后端直接配置文件进行更改，但是前端的配置文件，我感觉还是用热部署来更新前端的配置文件好像没什么大用处
// 既然前后端是分离开的那么就用吧
func SaveYAML() error {
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, byteData, fs.ModePerm)
}
