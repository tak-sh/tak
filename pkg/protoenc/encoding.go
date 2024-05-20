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

	ProtoUnmarshal = &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
)

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
		j, err := yaml.YAMLToJSON(b)
		if err != nil {
			return err
		}

		return ProtoUnmarshal.Unmarshal(j, msg)
	}
	return fmt.Errorf("%s is not a supported file extension (.yaml, .yml, .json)", name)
}
