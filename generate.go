// +build ignore

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

func get4(b []byte) int32 {
	if len(b) < 4 {
		return 0
	}
	return int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
}

const (
	INT_LEN = 4
)

const dataGoTempl = `package phoneregloc

var (
    total_len   = int32({{.totalLen}})
    firstoffset = int32({{.firstOffset}})
    content     = []byte{
	{{- range $slice := .content }}
        {{$slice}}
    {{- end -}}
}
)
`

func main() {
	output, err := os.OpenFile("./phone.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.FileMode(0755))
	if err != nil {
		panic(err)
	}
	defer output.Close()
	content, err := ioutil.ReadFile("./phone.dat")
	if err != nil {
		panic(err)
	}
	contents := make([]string, 0, len(content)/8+1)
	for i := 0; i < len(content); i += 8 {
		end := i + 8
		if end >= len(content) {
			end = len(content) - 1
		}
		s := ""
		for j, b := range content[i:end] {
			if j > 0 {
				s += " "
			}
			s += fmt.Sprintf("0x%02x,", b)
		}
		contents = append(contents, s)
	}
	totalLen := int32(len(content))
	firstoffset := get4(content[INT_LEN : INT_LEN*2])
	templ := template.New("t")
	templ.Parse(dataGoTempl)
	err = templ.Execute(output, map[string]interface{}{
		"totalLen":    totalLen,
		"firstOffset": firstoffset,
		"content":     contents,
	})
	if err != nil {
		panic(err)
	}
}
