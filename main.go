package main

import (
	"flag"
	"fmt"
	"os"

	"certview/cmd"
)

func main() {
	var (
		serverMode = flag.Bool("server", false, "Run in server mode")
		port       = flag.Int("port", 8080, "Server port (only in server mode)")
		help       = flag.Bool("help", false, "Show help")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "CertView - Certificate Analysis Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  CLI Mode:\n")
		fmt.Fprintf(os.Stderr, "    %s [options] <certificate-file|domain:port>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Server Mode:\n")
		fmt.Fprintf(os.Stderr, "    %s -server [-port=8080]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s cert.pem\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s google.com:443\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -server -port=8080\n", os.Args[0])
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *serverMode {
		cmd.RunServer(*port)
	} else {
		if flag.NArg() < 1 {
			fmt.Fprintf(os.Stderr, "Error: Missing certificate file or domain:port\n\n")
			flag.Usage()
			os.Exit(1)
		}
		input := flag.Arg(0)
		cmd.RunCLI(input)
	}
}