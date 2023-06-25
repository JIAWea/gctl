package parser

import (
	"bytes"
	"path"
	"text/template"

	"github.com/ml444/gctl/util"
	log "github.com/ml444/glog"
)

// func ParseTemplateToFile(pd *ProtoData, basePath, tempDir, tempName string, funcMap template.FuncMap) error {
//	fPath := filepath.Join(basePath, pd.Options["go_package"], strings.TrimSuffix(tempName, viper.GetString("template_format_suffix")))
//	return GenerateTemplate(fPath, tempDir, tempName, pd, funcMap)
// }

func GenerateTemplate(fPath string, tempFile, tempName string, data interface{}, funcMap template.FuncMap) error {
	var err error
	f, err := util.OpenFile(fPath)
	if err != nil {
		return err
	}
	temp := template.New(path.Base(tempName))
	if funcMap != nil {
		temp.Funcs(funcMap)
	}
	// temp, err = temp.ParseFS(tfs, tempFile)
	temp, err = temp.ParseFiles(tempFile)
	if err != nil {
		log.Error(err)
		return err
	}
	err = temp.Execute(f, data)
	if err != nil {
		log.Printf("Can't generate file %s,Error :%v\n", fPath, err)
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

var svcMethodTemp = `
{{$pn := .PackageName}}

{{ range $i, $svc := .ServiceList }}
{{$sn := SnakeToCamel $svc.ServiceName}}
{{$svcName := Concat $sn "Service" }}

{{ range $j, $v := $svc.RpcList }}
func (s {{ $svcName }}) {{$v.Name}}(ctx context.Context, req *{{$pn}}.{{$v.RequestType}}) (*{{$pn}}.{{$v.ResponseType}}, error) {
	var rsp {{$pn}}.{{$v.ResponseType}}
	return &rsp, nil
}
{{ end }}
{{ end }}
`

func GenerateServiceMethodContent(pd *ParseData, funcMap template.FuncMap) ([]byte, error) {

	temp := template.New("svcMethodTemp")
	if funcMap != nil {
		temp.Funcs(funcMap)
	}
	temp, err := temp.Parse(svcMethodTemp)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var buffer bytes.Buffer
	err = temp.Execute(&buffer, pd)
	if err != nil {
		log.Printf("Can't generate file content,Error :%v\n", err)
		return nil, err
	}
	return buffer.Bytes(), nil
}
