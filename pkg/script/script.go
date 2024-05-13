package script

import "github.com/tak-sh/tak/generated/go/script/v1beta1"

func Compile(s *v1beta1.Script) (*CompiledScript, error) {
	return nil, nil
}

func Validate(s *v1beta1.Script) error {
	return nil
}

type CompiledScript struct {
	script *v1beta1.Script
}
