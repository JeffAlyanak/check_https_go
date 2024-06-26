package check

import (
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jeffalyanak/check_https_go/tlsmap"
)

// HTTPCheck value
type HTTPCheck struct {
	URL string
}

// PerfData holds the Icinga/Nagios format Performance Data
type PerfData struct {
	timer  time.Time
	String string
}

// StartTimer adds a new key/value Performance Data metric
func (p *PerfData) StartTimer(t time.Time) {
	p.timer = t
}

// Add adds a new key/value Performance Data metric
func (p *PerfData) Add(key string, value string, uom string) {
	if p.String == "" {
		p.String += "|"
	}
	p.String += key + "=" + value + uom + " "
}

// Get returns all the key/value Performance Data metrics
// Stops the timer and adds that to the key/value pairs.
func (p *PerfData) Get() string {
	duration := time.Since(p.timer)
	durationstr := strconv.Itoa(int(duration.Milliseconds()))

	if p.String == "" {
		p.String += "|"
	}
	return p.String + "checks_took" + "=" + durationstr + "ms"
}

// parseStatusCodes takes a comma-seperated string of HTTP status codes
// parses and returns a slice of int as well as any errors.
func parseStatusCodes(userStatusCodes string) ([]int, error) {
	if userStatusCodes == "" {
		return nil, nil
	}

	statusCodes := strings.Split(userStatusCodes, ",")
	result := make([]int, len(statusCodes))
	for i, code := range statusCodes {
		status, err := strconv.Atoi(strings.TrimSpace(code))
		if err != nil {
			return nil, err
		}
		result[i] = status
	}
	return result, nil
}

// CheckStatus function runs a check of the HTTP status code and returns the result.
func (h *HTTPCheck) CheckStatus(redirects int, userAgent string, timeoutduration int, userStatusCodes string) Result {
	var resp *http.Response
	var r Result
	var url string = "https://" + h.URL
	var statusCodes = []int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226}

	// Create request for domain with a User-Agent header.
	tr := &http.Transport{}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeoutduration) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse

		},
	}

	// Follow redirects up to redirect limit
	for i := 0; i < redirects; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			r.Error = err
			return r
		}

		req.Header.Set("User-Agent", userAgent)
		resp, err = client.Do(req)
		if err != nil {
			r.Error = err
			return r
		}

		if resp.StatusCode == 301 || resp.StatusCode == 302 || resp.StatusCode == 307 {
			l := resp.Header.Get("Location")
			r.VerboseValue += url + " redirected (" + resp.Status + ") to " + l + "\n"

			// If the location is relative to the domain
			if !strings.HasPrefix(l, "http") {
				// If the new location is relative to the old location, simply add it
				// to the end of the previous location.
				if !strings.HasPrefix(l, "/") {
					l = url + "/" + l
					// Otherwise—if the new location is relative to the webroot—extract the
					// root from the old location and then add the new location to it.
					// Error/break if unable to parse the URL from the old location.
				} else {
					re, _ := regexp.Compile("https?://[0-9a-z-.]+")
					found := re.FindAllString(url, -1)

					if len(found) > 0 {
						l = found[0] + l
					} else {
						r.Error = errors.New("check http status: could not parse a valid URL")
						break
					}
				}
			}
			url = l
		} else {
			break
		}
	}

	r.Status = resp.StatusCode
	r.Value = http.StatusText(resp.StatusCode)

	if userStatusCodes != "" {
		parsedStatusCodes, err := parseStatusCodes(userStatusCodes)
		if err != nil {
			r.Error = err
			return r
		}
		statusCodes = parsedStatusCodes
	}

	statusCodeGood := false
	for _, code := range statusCodes {
		if code == r.Status {
			statusCodeGood = true
			break
		}
	}

	// Return code is 0 (OK) for any 2xx status code.
	if statusCodeGood {
		r.ReturnCode = 0

	} else {
		r.ReturnCode = 2
	}

	return r
}

// CheckContent function runs a check of returned body content and returns the result.
func (h *HTTPCheck) CheckContent(checkString string) Result {
	var r Result

	resp, err := http.Get("https://" + h.URL)
	if err != nil {
		r.Error = err
		return r
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.Error = err
		return r
	}

	if string(body) != "" {
		if strings.ContainsAny(string(body), checkString) {
			r.ReturnCode = 0
			r.Value = "Expected content returned: " + checkString
		} else {
			r.ReturnCode = 3
			r.Value = "Unknown content returned"
		}
	} else {
		r.ReturnCode = 2
		r.Value = "No content returned"
	}

	// Verbose output includes the first five lines of the body
	lines := strings.Split(string(body), "\n")

	r.VerboseValue = "Returned " + strconv.Itoa(len(lines)) + " lines of content.\n"

	return r
}

// CheckCertificate function runs a check of TLS certificate and returns the result.
func (h *HTTPCheck) CheckCertificate(warn int, crit int) Result {
	var r Result
	var c *x509.Certificate

	resp, err := http.Get("https://" + h.URL)
	if err != nil {
		r.Error = err
		return r
	}

	certs := resp.TLS.PeerCertificates
	if len(certs) > 0 {
		c = certs[0]
	} else {
		r.Error = errors.New("TLS error: no certificates returned")
		return r
	}

	warndate := time.Now().AddDate(0, 0, warn)
	critdate := time.Now().AddDate(0, 0, crit)

	if critdate.After(c.NotAfter) {
		r.ReturnCode = 2
		r.Value = "Cert critical"
	} else if warndate.After(c.NotAfter) {
		r.ReturnCode = 1
		r.Value = "Cert warning"
	} else {
		r.ReturnCode = 0
		r.Value = "Cert okay"
	}
	r.Value = r.Value + ", valid until " + c.NotAfter.Format("January 02, 2006 15:04")

	// Verbose info on TLS version and cipher suite
	r.VerboseValue += "TLS Version used:  " + tlsmap.TLSVersion(resp.TLS.Version) + "\n"
	r.VerboseValue += "Cipher suite used: " + tlsmap.CipherSuite(resp.TLS.CipherSuite) + "\n"

	return r
}
