// cli.go

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/keltia/cryptcheck"
	"github.com/keltia/observatory"
	"github.com/keltia/ssllabs"
)

var (
	fDebug bool

	fJobs int

	fList    string
	fType    string
	fOutput  string
	fSummary string

	fVerbose bool
)

const (
	cliUsage = `%s version %s - Imirhil/%s SSLLabs/%s Mozilla/%s

Usage: %s [-hvDIM] [-j n] [-t text|csv|html] [-s file] [-o file] [-file fn] [site]

`
)

// Usage string override.
var Usage = func() {
	fmt.Fprintf(os.Stderr, cliUsage, MyName,
		MyVersion, cryptcheck.MyVersion, ssllabs.MyVersion, observatory.MyVersion,
		MyName, MyName)
	flag.PrintDefaults()
}

func init() {
	// Main switches
	flag.StringVar(&fList, "list", "", "Specify the list of sites")

	flag.StringVar(&fType, "t", "csv", "Type of report")

	flag.IntVar(&fJobs, "j", runtime.NumCPU(), "# of parallel jobs")
	flag.StringVar(&fOutput, "o", "-", "Save into file (default stdout)")
	flag.StringVar(&fSummary, "s", "summaries", "Save summary there")

	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
}
