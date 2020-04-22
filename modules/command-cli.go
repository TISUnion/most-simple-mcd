package modules

import (
	"bytes"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/utils"
	"path/filepath"
	"strings"
	"text/template"
)

// 插件模板参数
type pluginTmlp struct {
	Dirname, ENName, ZHName, Description, Command, HelpDescription string
	IsGlobal                                                       bool
}

// 创建一个插件文件
func CreatePluginTmplFile(dirname, enName, zhName, description, command, helpDescription string, isGlobal bool) error {
	tmplParam := pluginTmlp{
		Dirname:         dirname,
		ENName:          enName,
		ZHName:          zhName,
		Description:     description,
		Command:         command,
		HelpDescription: helpDescription,
		IsGlobal:        isGlobal,
	}
	tmplFuncMap := template.FuncMap{
		// 注册函数title, strings.Title会将单词首字母大写
		"filename2packagename": filename2packagename,
	}
	tmpl, err := template.New("pluginTmpl").Funcs(tmplFuncMap).Parse(constant.PLUGIN_TMPL)
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, tmplParam); err != nil {
		return err
	}
	tmplBytes := buf.Bytes()
	pluginPath := filepath.Join(GetConfVal(constant.WORKSPACE), constant.PLUGIN_DIR, dirname, dirname+".go")
	fileObj, err := utils.CreateFile(pluginPath)
	if err != nil {
		return err
	}
	defer fileObj.Close()
	_, err = fileObj.Write(tmplBytes)
	if err != nil {
		return err
	}
	return nil
}

// 文件名转包名
func filename2packagename(filename string) string {
	return strings.ReplaceAll(filename, "-", "_")
}
