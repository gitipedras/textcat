package tc

import(
	"errors"
	"strings"
	"fmt"
)

func MakeError(msg ...any) error {
	if len(msg) == 0 {
		return nil
	}

	var b strings.Builder
	for i, m := range msg {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprint(m))
	}

	return errors.New(b.String())
}
