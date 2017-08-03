package acsengine

//go:generate go-bindata -nometadata -pkg $GOPACKAGE -prefix ../../parts/ -o templates.go ../../parts/ ../../parts/kubernetes/agentpool
//go:generate gofmt -s -l -w templates.go
// fileloader use go-bindata (https://github.com/jteeuwen/go-bindata)
// go-bindata is the way we handle embedded files, like binary, template, etc.
