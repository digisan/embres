package embres

import "testing"

func TestFileBytes(t *testing.T) {
	SetResAliasAsKey("Abc", "./go.sum")
	PrintFileBytes("embres", "MapRes", "./test", false, "./go.sum")
	// fPln(string(MapRes["Abc"]))
}

func TestCreateResDirBytesFile(t *testing.T) {
	GenerateDirBytes("embres", "MapResDir", "./", "./testDir", true, "Git")
	// fPln(string(MapResDir["auto0002"]))
}
