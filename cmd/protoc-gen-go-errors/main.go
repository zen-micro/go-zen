package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var versionFlag = flag.Bool("version", false, "print the version and exit")

var version string

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Printf("Version: %s\n", version)
		return
	}
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate || !needGenerate(f.Proto.SourceCodeInfo.Location) {
				continue
			}

			generateFile(gen, f)
		}
		return nil
	})
}
