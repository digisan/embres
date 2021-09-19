package embres

import (
	"fmt"
	"strings"

	"github.com/digisan/gotk/io"
	"github.com/digisan/logkit"
)

var (
	fSf             = fmt.Sprintf
	fEf             = fmt.Errorf
	sReplace        = strings.Replace
	sReplaceAll     = strings.ReplaceAll
	sTitle          = strings.Title
	sTrim           = strings.Trim
	sTrimLeft       = strings.TrimLeft
	sTrimSuffix     = strings.TrimSuffix
	sContains       = strings.Contains
	failOnErr       = logkit.FailOnErr
	failP1OnErr     = logkit.FailP1OnErr
	failP1OnErrWhen = logkit.FailP1OnErrWhen
	mustWriteFile   = io.MustWriteFile
	fileLineScan    = io.FileLineScan
	rmTailFromLast  = logkit.RmTailFromLast

	replAllOnAny = func(s string, olds []string, new string) string {
		for _, old := range olds {
			s = sReplaceAll(s, old, new)
		}
		return s
	}
)
