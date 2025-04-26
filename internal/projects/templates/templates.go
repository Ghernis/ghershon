package templates

import (
	"embed"
)

//go:embed *.tmpl
var FS embed.FS

//func ParseTemplate(name string) ([]byte, error){
//	temp,err:= FS.ReadFile(name)
//	if err != nil{
//		fmt.Errorf("Loading template: %w",err)
//	}
//	return temp,nil
//}
