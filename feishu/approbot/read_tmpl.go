package approbot

import (
	"embed"
)

//go:embed tmpl/*.json
var staticFiles embed.FS

func readAllTemplateJsonFiles(folder string) (map[string]string, error) {
	// 读取目录下的所有 JSON 文件
	files, err := staticFiles.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	var tmplJsonMap = make(map[string]string)
	for _, file := range files {
		key := file.Name()
		data, err := staticFiles.ReadFile(folder + "/" + file.Name())
		if err != nil {
			continue
		}
		tmplJsonMap[key] = string(data)
	}
	return tmplJsonMap, nil
}
