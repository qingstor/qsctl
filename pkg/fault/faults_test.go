package fault

import (
	"errors"
	"fmt"
	"testing"
)

func TestFault(t *testing.T) {
	v := errors.New("test")
	x := fmt.Errorf("\nxxxx:\n %w", v)
	t.Log(x)
}
