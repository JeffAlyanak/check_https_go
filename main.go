package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"git.computerassistance.co.uk/icinga-scripts/check_https_go/check"
)

func main() {
	// Handle cli arguments
	host := flag.String("h", "", "Fully-qualified domain name to check.")
	userAgent := flag.String("u", "check_https_go", "Custom user-agent string.")
	verbose := flag.Bool("v", false, "More verbose output includes details of any redirects.")
	redirects := flag.Int("r", 20, "Number of redirects to follow.")
	certwarn := flag.Int("w", 10, "Number of days for which the TLS certificate must be valid before a warning state is returned.")
	certcrit := flag.Int("c", 5, "Number of days for which the TLS certificate must be valid before a critical state is returned.")

	flag.Parse()

	if *host == "" {
		fmt.Println("Please provide a fully-qualified domain name.")
		os.Exit(3)
	}

	var h check.HTTPCheck
	h.URL = *host

	// Perform status code check, exit with additional info if error
	statuscodeResult := h.CheckStatus(*redirects, *userAgent)
	if statuscodeResult.Error != nil {
		fmt.Println(statuscodeResult.ReturnCode)
		if *verbose {
			fmt.Println("\nAdditional Info:")
			fmt.Println(statuscodeResult.VerboseValue)
		}
		os.Exit(3)
	}

	// Print basic info about status code check
	fmt.Println("HTTPS Check for https://" + h.URL)
	fmt.Println("Status Code: " + strconv.Itoa(statuscodeResult.Status) + " " + statuscodeResult.Value)

	// Check return code, exit with additional info if non-zero
	if statuscodeResult.ReturnCode != 0 {
		if *verbose {
			fmt.Println("\nAdditional Info:")
			fmt.Println(statuscodeResult.VerboseValue)
		}
		os.Exit(statuscodeResult.ReturnCode)
	}

	// Perform content check, exit with additional info if error
	contentResult := h.CheckContent()
	if contentResult.Error != nil {
		fmt.Println(statuscodeResult.Error)
		if *verbose {
			fmt.Println("\nAdditional Info:")
			fmt.Println(statuscodeResult.VerboseValue)
			fmt.Println(contentResult.VerboseValue)
		}
		os.Exit(3)
	}

	// Print basic info about content check
	fmt.Println("Content Check: " + contentResult.Value)

	// Check return code, exit with additional info if non-zero
	if contentResult.ReturnCode != 0 {
		if *verbose {
			fmt.Println("\nAdditional Info:")
			fmt.Println(statuscodeResult.VerboseValue)
			fmt.Println(contentResult.VerboseValue)
		}
		os.Exit(contentResult.ReturnCode)
	}

	// Perform content check, exit with additional info if error
	certResult := h.CheckCertificate(*certwarn, *certcrit)
	if certResult.Error != nil {
		fmt.Println(statuscodeResult.Error)
		if *verbose {
			fmt.Println("\nAdditional Info:")
			fmt.Println(statuscodeResult.VerboseValue)
			fmt.Println(contentResult.VerboseValue)
			fmt.Println(certResult.VerboseValue)
		}
		os.Exit(3)
	}

	// Print basic info about certificate check
	fmt.Println("Cert Check: " + certResult.Value)

	// Check return code, exit with additional info if non-zero
	if certResult.ReturnCode != 0 {
		if *verbose {
			fmt.Println("\nAdditional Info:")
			fmt.Println(statuscodeResult.VerboseValue)
			fmt.Println(contentResult.VerboseValue)
			fmt.Println(certResult.VerboseValue)
		}
		os.Exit(certResult.ReturnCode)
	}

	if *verbose {
		fmt.Println("\nAdditional Info:")
		fmt.Println(statuscodeResult.VerboseValue)
		fmt.Println(contentResult.VerboseValue)
		fmt.Println(certResult.VerboseValue)
	}
	os.Exit(0)
}
