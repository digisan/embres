package embres

import "testing"

func TestPrintLineKeyMap(t *testing.T) {
	PrintLineKeyMap("embres", "MyMap", "./test1", "./go.mod")
}
