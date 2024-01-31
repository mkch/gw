package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

var t = template.Must(template.New("methods").Funcs(
	template.FuncMap{
		"loop": func(from, to int) (result []int) {
			for i := from; i <= to; i++ {
				result = append(result, i)
			}
			return
		},
		"add": func(n, m int) int {
			return n + m
		},
	}).Parse(`// Code generated by .scripts/genmethods.go DO NOT EDIT.

package mscom

import (
	"unsafe"

	"golang.org/x/sys/windows"
)
{{range $i := loop 0 . }}
func newMethod{{$i}}() (h method) {
	h.nArg = {{$i}}
	h.ptr = windows.NewCallback(func(this unsafe.Pointer{{range $j := loop 1 $i -}}, arg{{$j}} uintptr{{end}}) uintptr {
		f := mtdMap.Methods(this).Method(h.ptr).(func({{if $i}}uintptr{{end}} {{- range $j := loop 2 $i}}, uintptr{{end}}) uintptr)
		return f({{if $i}}arg1{{end}} {{- range $j := loop 2 $i}}, arg{{$j}}{{end}})
	})
	return
}
{{- end}}

var methodFactory = map[int]func() method{
{{- range $i := loop 0 . }}
	{{$i}}:	newMethod{{$i}},
{{- end}}
}
`))

func main() {
	var count uint
	flag.UintVar(&count, "count", 10, "count of methods")
	var output string
	flag.StringVar(&output, "o", "methods.go", "The output file name")
	flag.Parse()
	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		errorExit(err.Error() + "\n")
	}
	defer f.Close()
	if err := t.Execute(f, int(count)); err != nil {
		panic(err)
	}

	if message, err := exec.Command("go", "fmt", output).CombinedOutput(); err != nil {
		errorExit(fmt.Sprintf(`Error occurred when executing "go fmt"`+"\n%v", string(message)))
	}
}

func errorExit(message string) {
	os.Stderr.WriteString(message)
	os.Exit(1)
}
