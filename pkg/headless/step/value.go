package step

import (
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"strings"
)

func GetValueString(v *v1beta1.Value) string {
	if v == nil {
		return ""
	}

	if v.Str != nil {
		return *v.Str
	} else if len(v.StrList) > 0 {
		return strings.Join(v.StrList, ",")
	}

	return ""
}
