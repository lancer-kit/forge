package scaffolder

//go:generate go-bindata -nomemcopy -pkg defaultp -prefix templates/ -ignore schema.yml -o project/provider/defaultp/bindata.go templates/...
