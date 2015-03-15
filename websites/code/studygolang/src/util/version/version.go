// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

/*
sets version information for the binary where it is imported.
The version can be retrieved either from the -version command line argument,
or from the /version/ http endpoint.

To include in a project simply import the package and call version.Init().

The version and compile date is stored in version and date variables and
are supposed to be set during compile time. Typically this is done by the
install(bash/bat). Or date is binary modify time.

To set these manually use -ldflags together with -X, like in this example:

	go install -ldflags "-X util/version.Version v1.0"

*/

package version

import (
	"flag"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
)

var showVersion = flag.Bool("version", false, "Print version of this binary")

var (
	Version string
	Date    string
)

func printVersion(w io.Writer, version string, date string) {
	fmt.Fprintf(w, "Version: %s\n", version)
	fmt.Fprintf(w, "Binary: %s\n", os.Args[0])
	fmt.Fprintf(w, "Compile date: %s\n", date)
}

// initializes the version flag and /version/ http endpoint.
// Note that this method will call flag.Parse if the flags are not already parsed.
func init() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if Date == "" {
		fileInfo, err := os.Stat(os.Args[0])
		if err != nil {
			panic(err)
		}

		Date = fileInfo.ModTime().String()
	}

	if showVersion != nil && *showVersion {
		printVersion(os.Stdout, Version, Date)
		os.Exit(0)
	}

	http.Handle("/version", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		printVersion(w, html.EscapeString(Version), html.EscapeString(Date))
	}))
}
