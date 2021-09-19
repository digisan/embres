package embres

import (
	"unicode"
)

// PrintLineKeyMap :
func PrintLineKeyMap(pkg, outmap, savepath, file string) string {
	failP1OnErrWhen(
		!rPkgNameRule.MatchString(pkg),
		"%v", fEf("[%s] error: Not Valid Package Name", pkg),
	)
	failP1OnErrWhen(
		!rNameRule.MatchString(outmap) || unicode.IsLower(rune(outmap[0])) || outmap[0] == '_',
		"%v", fEf("[%s] error: Not Valid Exportable Variable Name", outmap),
	)

	content, err := fileLineScan(file, func(line string) (bool, string) {
		if sTrim(line, " \t\r\n") == "" {
			return false, ""
		}
		return true, fSf(`	"%s": struct{}{},`, line)
	}, "")
	failOnErr("%v", err)

	content = fSf("package %s\n\n// %s : auto-generated\nvar %s = map[string]struct{}{\n%s\n}", pkg, outmap, outmap, content)

	if savepath != "" {
		savepath = sTrimSuffix(savepath, ".go") + ".go"
		mustWriteFile(savepath, []byte(content))
	}

	return content
}
