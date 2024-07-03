package protoenc

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io/fs"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

var (
	ProtoMarshal = &protojson.MarshalOptions{
		Multiline: true,
	}

	ProtoMarshalPretty = &protojson.MarshalOptions{
		Indent: "  ",
	}

	ProtoUnmarshal = &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
)

func MarshalYAML(msg proto.Message) ([]byte, error) {
	b, err := ProtoMarshal.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return yaml.JSONToYAML(b)
}

func UnmarshalFile(msg proto.Message, name string, f fs.FS) error {
	ext := filepath.Ext(name)
	b, err := fs.ReadFile(f, name)
	if err != nil {
		return err
	}

	switch ext {
	case ".json":
		return ProtoUnmarshal.Unmarshal(b, msg)
	case ".yaml", ".yml":
		return UnmarshalYAML(msg, b)
	}
	return fmt.Errorf("%s is not a supported file extension (.yaml, .yml, .json)", name)
}

func UnmarshalYAML(msg proto.Message, b []byte) error {
	j, err := yaml.YAMLToJSON(b)
	if err != nil {
		return err
	}

	return ProtoUnmarshal.Unmarshal(j, msg)
}
