package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"text/template"
	"time"
	"unicode"
)

const service = `Mono<{{.Method | Capitalize}}Response> {{.Method}}({{.Method | Capitalize}}Request {{.Method}}Request);`

const serviceImpl = `
/**
 * {{.Method | Blank}}
 *
 * @param {{.Method}}Request
 * @return Mono<{{.Method | Capitalize}}Response>
 * @DateTime : {{ now.Format "2006/01/02 3:04 PM" }}
 * @Summary : 
 */
@Override
public Mono<{{.Method | Capitalize}}Response> {{.Method}}({{.Method | Capitalize}}Request {{.Method}}Request) {
	final String path = "/{{.Method}}";
	return webClient.callGet(npcGatewayUrl, vpcPath + path, {{.Method}}Request, {{.Method | Capitalize}}Response.class);
}`

const controller = `
/**
 * {{.Method | Blank}}
 *
 * @param {{.Method}}Request
 * @return Mono<{{.Method | Capitalize}}Response>
 * @DateTime : {{ now.Format "2006/01/02 3:04 PM" }}
 * @Summary : 
 */
@ApiOperation(value = "")
	@ApiImplicitParams({
			@ApiImplicitParam(name = "{{.Method}}Request", value = "", required = true, dataTypeClass = {{.Method | Capitalize}}Request.class, paramType = "body")
	})
@PostMapping("/{{.Method}}")
public Mono<{{.Method | Capitalize}}Response> {{.Method}}(@RequestBody {{.Method | Capitalize}}Request {{.Method}}Request) {
	return service.{{.Method}}({{.Method}}Request);
}`

const test = `
POST http://localhost:8080/api/v1/networking/vpc/{{.Method}}
Content-Type: application/json

{}
###`

func ToCapitalize(s string) string {
	for i, v := range s {
		return string(unicode.ToUpper(v)) + s[i+1:]
	}
	return s
}

func Blank(s string) string {
	for i := 0; i < len(s); i++ {
		if 'A' <= s[i] && s[i] <= 'Z' {
			s = s[:i] + " " + strings.ToLower(string(s[i])) + s[i+1:]
			i++
		}
	}
	return s
}

type Method struct {
	Method string
}

type MethodBox struct {
	Methods []Method
}

func (box *MethodBox) AddMethod(method Method) {
	box.Methods = append(box.Methods, method)
}

func main() {
	funcMap := template.FuncMap{
		"Capitalize": ToCapitalize,
		"now":        time.Now,
		"Blank":      Blank,
	}

	var inputMethod string
	var whichComp string

	flag.StringVar(&inputMethod, "methods", "", "methods to create struct")
	flag.StringVar(&whichComp, "component", "", "templates to create")
	flag.Parse()

	Methods := MethodBox{}

	for _, m := range strings.Split(inputMethod, " ") {
		method := Method{
			Method: m,
		}
		Methods.AddMethod(method)
	}

	switch whichComp {
	case "serviceImpl":
		whichComp = serviceImpl
	case "controller":
		whichComp = controller
	case "service":
		whichComp = service
	case "test":
		whichComp = test
	default:
	}

	t := template.Must(template.New("").Funcs(funcMap).Parse(whichComp))

	for _, method := range Methods.Methods {
		buffer := bytes.Buffer{}
		if err := t.Execute(&buffer, method); err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", buffer.String())
	}
}
