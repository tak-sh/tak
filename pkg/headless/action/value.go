package action

import (
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
)

func GetValue(v *v1beta1.Value) any {
	if v == nil {
		return nil
	}

	if v.Str != nil {
		return *v.Str
	} else if len(v.StrList) > 0 {
		return v.StrList
	}

	return nil
}
