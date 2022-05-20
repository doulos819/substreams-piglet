package manifest

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/jhump/protoreflect/desc/protoparse"
	statikfs "github.com/rakyll/statik/fs"
	pbsubstreams "github.com/streamingfast/substreams/pb/sf/substreams/v1"
	_ "github.com/streamingfast/substreams/pb/statik"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func loadProtobufs(pkg *pbsubstreams.Package, manif *Manifest) error {
	// System protos
	systemFiles, err := readSystemProtobufs()
	if err != nil {
		return err
	}
	for _, file := range systemFiles.File {
		pkg.ProtoFiles = append(pkg.ProtoFiles, file)
	}

	var importPaths []string
	for _, imp := range manif.Protobuf.ImportPaths {
		importPaths = append(importPaths, filepath.Join(manif.Workdir, imp))
	}
	// User-specified protos
	parser := protoparse.Parser{
		ImportPaths:           importPaths,
		IncludeSourceCodeInfo: true,
	}

	customFiles, err := parser.ParseFiles(manif.Protobuf.Files...)
	if err != nil {
		return fmt.Errorf("error parsing proto files %q (import paths: %q): %w", manif.Protobuf.Files, importPaths, err)
	}
	for _, fd := range customFiles {
		pkg.ProtoFiles = append(pkg.ProtoFiles, fd.AsFileDescriptorProto())
	}

	return nil
}

func readSystemProtobufs() (*descriptorpb.FileDescriptorSet, error) {
	sfs, err := statikfs.New()
	if err != nil {
		return nil, err
	}
	staticFDS, err := sfs.Open("/system.pb") // see generation in pb/generate.sh
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(staticFDS)
	if err != nil {
		return nil, err
	}
	fds := &descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(b, fds)
	if err != nil {
		return nil, err
	}
	return fds, nil
}