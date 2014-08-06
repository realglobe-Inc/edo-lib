package util

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

// 標準の FlagSet だと ContinueOnError にしても、
// エラーのところで読み込みを終えてしまうので、
// 最後まで読み込めるようにする。
type FlagSet struct {
	*flag.FlagSet
	name string
	flag.ErrorHandling
	output io.Writer
}

func NewFlagSet(name string, errorHandling flag.ErrorHandling) *FlagSet {
	return &FlagSet{
		flag.NewFlagSet(name, errorHandling),
		name,
		errorHandling,
		nil,
	}
}

func (f *FlagSet) SetOutput(output io.Writer) {
	f.output = output
	f.FlagSet.SetOutput(output)
}

func (f *FlagSet) CompleteParse(args []string) {
	usage := f.FlagSet.Usage
	f.FlagSet.Usage = func() {}
	f.FlagSet.Init(f.name, flag.ContinueOnError)

	errTokens := []string{}
	for curArgs := args; len(curArgs) > 0; curArgs = f.FlagSet.Args() {
		buff := &bytes.Buffer{}
		f.FlagSet.SetOutput(buff)
		f.FlagSet.Parse(curArgs)
		if buff.Len() == 0 {
			// 正常に最後まで読み取った。
			break
		}
		errTokens = append(errTokens, curArgs[len(curArgs)-len(f.FlagSet.Args())-1])
	}

	f.FlagSet.Usage = usage
	f.FlagSet.Init(f.name, f.ErrorHandling)
	f.FlagSet.SetOutput(f.output)

	if len(errTokens) > 0 {
		fmt.Fprintln(os.Stderr, "flag provided but not defined:")
		for _, errToken := range errTokens {
			fmt.Fprintln(os.Stderr, "  "+errToken)
		}
		f.FlagSet.Usage()
	}
}
