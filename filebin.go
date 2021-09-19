package embres

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var (
	mPathAlias   = make(map[string][]string)
	mAliasPath   = make(map[string]string)
	rNameRule    = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	rPkgNameRule = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
)

// SetResAliasAsKey :
func SetResAliasAsKey(alias, file string) {
	_, err := os.Stat(file)
	failP1OnErr("%v", err)
	fpAbs, err := filepath.Abs(file)
	failP1OnErr("%v", err)

	_, exist := mAliasPath[alias]
	failP1OnErrWhen(exist, "%v", fEf("alias [%s] is already used by [%s]", alias, mAliasPath[alias]))
	mAliasPath[alias], mPathAlias[fpAbs] = fpAbs, append(mPathAlias[fpAbs], alias)
}

// PrintFileBytes :
func PrintFileBytes(pkg, outmap, savepath string, keepext bool, files ...string) string {
	failP1OnErrWhen(
		!rPkgNameRule.MatchString(pkg),
		"%v", fEf("[%s] error: Not Valid Package Name", pkg),
	)
	failP1OnErrWhen(
		!rNameRule.MatchString(outmap) || unicode.IsLower(rune(outmap[0])) || outmap[0] == '_',
		"%v", fEf("[%s] error: Not Valid Exportable Variable Name", outmap),
	)

	sb := strings.Builder{}
	for _, file := range files {
		fpAbs, err := filepath.Abs(file)
		failP1OnErr("%v", err)

		if _, exist := mPathAlias[fpAbs]; !exist {
			bytesName := rmTailFromLast(file, ".")
			bytesName = replAllOnAny(bytesName, []string{".", "-"}, "")
			bytesName = replAllOnAny(bytesName, []string{"/"}, "_")
			suffix := ""
			if keepext {
				suffix = "_" + sTrimLeft(filepath.Ext(file), ".")
			}
			bytesName = sTrimLeft(bytesName, " \t_") + suffix
			mPathAlias[fpAbs] = append([]string{}, sTitle(bytesName))
		}

		bytes, err := os.ReadFile(file)
		failP1OnErr("%v", err)

		for _, alias := range mPathAlias[fpAbs] {
			sb.WriteString(fSf("\t\"%s\": {\n\t\t", alias))
			for i, v := range bytes {
				if i > 0 {
					if i%16 == 0 {
						sb.WriteString(",\n\t\t")
					} else {
						sb.WriteString(", ")
					}
				}
				sb.WriteString(fSf("0x%02x", v))
			}
			sb.WriteString(",\n\t},\n")
		}
	}

	content := fSf("package %s\n\n// %s : auto-generated\nvar %s = map[string][]byte{\n%s}\n", pkg, outmap, outmap, sb.String())

	// deal with `"_":` & `"":`
	I := 0
	for _, old := range []string{`"_":`, `"":`} {
		for sContains(content, old) {
			content = sReplace(content, old, fSf("\"%04d\":", I), 1)
			I++
		}
	}

	if savepath != "" {
		savepath = sTrimSuffix(savepath, ".go") + ".go"
		mustWriteFile(savepath, []byte(content))
	}
	return content
}

// GenerateDirBytes :
func GenerateDirBytes(pkg, outmap, dir, savepath string, keepext bool, trimkey ...string) string {
	fDir, err := os.Open(dir)
	failP1OnErr("%v", err)
	dInfo, err := fDir.Stat()
	failP1OnErr("%v", err)
	failP1OnErrWhen(!dInfo.IsDir(), "%v", fEf("input dir is invalid"))

	resGrp := []string{}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.HasSuffix(info.Name(), ".go") {
			resGrp = append(resGrp, path)
		}
		return nil
	})
	failP1OnErr("%v", err)
	resfile := PrintFileBytes(pkg, outmap, "", keepext, resGrp...)

	// trim keys
	for _, seg := range trimkey {
		resfile = replAllOnAny(resfile, []string{`"` + seg + `_`, `_` + seg + `"`}, `"`)
		resfile = sReplaceAll(resfile, `_`+seg+`_`, `_`)
	}

	if savepath != "" {
		savepath = sTrimSuffix(savepath, ".go") + ".go"
		mustWriteFile(savepath, []byte(resfile))
	}

	return resfile
}
