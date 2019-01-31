// main.go

/*
This package implements reading the json from ssllabs-scan output
and generating a csv file.
*/
package main // import "github.com/keltia/tlscheck"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/keltia/cryptcheck"
	"github.com/keltia/observatory"
	"github.com/keltia/ssllabs"
)

var (
	// MyName is obvious
	MyName = filepath.Base(os.Args[0])

	contracts map[string]string
	tmpls     map[string]string

	logLevel = 0
)

const (
	// MyVersion uses semantic versioning.
	MyVersion = "0.63.0"
)

// checkOutput checks whether we want to specify an output file
func checkOutput(fOutput string) (fOutputFH *os.File) {
	var err error

	fOutputFH = os.Stdout

	// Open output file
	if fOutput != "" {
		verbose("Output file is %s\n", fOutput)

		if fOutput != "-" {
			fOutputFH, err = os.Create(fOutput)
			if err != nil {
				fatalf("Error creating %s\n", fOutput)
			}
		}
	}
	debug("output=%v\n", fOutputFH)
	return
}

// init is for pg connection and stuff
func init() {
	flag.Usage = Usage
	flag.Parse()
}

func readSiteList(fn string) ([]string, error) {
	buf, err := ioutil.ReadFile(fn)
	return strings.Split(string(buf), "\n"), err
}

func checkFlags(a []string) ([]string, error) {
	var (
		err      error
		siteList []string
	)

	// Basic argument check
	if fList != "" && len(a) != 0 {
		return nil, fmt.Errorf("You can't specify both -list and a single sitename!")
	}

	if a == nil || len(a) == 0 && fList == "" {
		return nil, fmt.Errorf("you must specify either -list or a sitename!")
	}

	if fList == "" {
		siteList = append(siteList, flag.Arg(0))
	} else {
		if siteList, err = readSiteList(fList); err != nil {
			return nil, err
		}
	}

	// Set logging level
	if fVerbose {
		logLevel = 1
	}

	if fDebug {
		fVerbose = true
		logLevel = 2
		debug("debug mode\n")
	}
	return siteList, nil
}

// main is the the starting point
func main() {
	// Announce ourselves
	fmt.Printf("%s version %s/j%d - Imirhil/%s SSLLabs/%s Mozilla/%s\n\n",
		filepath.Base(os.Args[0]), MyVersion, fJobs,
		cryptcheck.MyVersion, ssllabs.MyVersion, observatory.MyVersion)

	siteList, err := checkFlags(flag.Args())
	if err != nil {
		fatalf("Error: %v", err.Error())
	}

	allSites, err := processList(siteList)
	if err != nil {
		fatalf("Can't process %v: %s", siteList, err.Error())
	}

	err = loadResources(resourcesPath)
	if err != nil {
		fatalf("Can't load resources %s: %v", resourcesPath, err)
	}

	// Open output file
	fOutputFH := checkOutput(fOutput)

	// generate the final report & summary
	final, err := NewTLSReport(allSites)
	if err != nil {
		fatalf("error analyzing report: %v", err)
	}

	// Gather statistics for summaries
	cntrs := categoryCounts(allSites)
	https := httpCounts(final)

	verbose("SSLabs engine: %s\n", final.SSLLabs)

	switch fType {
	case "csv":
		err = WriteCSV(fOutputFH, final, cntrs, https)
		if err != nil {
			fatalf("WriteCSV failed: %v", err)
		}
	case "html":
		err = WriteHTML(fOutputFH, final, cntrs, https)
		if err != nil {
			fatalf("WriteHTML failed: %v", err)
		}
	default:
		// XXX Early debugging
		fmt.Printf("%#v\n", final)
		fmt.Printf("%s\n", displayCategories(cntrs))

	}
}
