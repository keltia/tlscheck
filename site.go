package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/keltia/cryptcheck"
	"github.com/keltia/observatory"
	"github.com/keltia/ssllabs"
)

const (
	DefaultKeySize = 2048
	DefaultAlg     = "RSA"
	DefaultSig     = "SHA256withRSA"
)

var (
	fnImirhil func(site ssllabs.Host) string
	fnMozilla func(site ssllabs.Host) string

	DefaultIssuer = regexp.MustCompile(`(?i:GlobalSign)`)
)

func fixTimestamp(ts int64) (int64, int64) {
	return ts / 1000, ts % 1000
}

func checkSweet32(det ssllabs.EndpointDetails) (yes bool) {
	if len(det.Suites) != 0 {
		ciphers := det.Suites[0].List
		for _, cipher := range ciphers {
			if strings.Contains(cipher.Name, "DES") {
				return true
			}
		}
	}
	return false
}

func getGrade(site ssllabs.Host, fn func(site ssllabs.Host) string) string {
	return fn(site)
}

func hasExpired(t int64) bool {
	return time.Now().After(time.Unix(fixTimestamp(t)))
}

func findServerType(site ssllabs.Host) int {
	// Should be obvious, 2nd field is only present if no valid cert is found
	if len(site.Certs) == 0 || len(site.CertHostnames) != 0 {
		return TypeHTTP
	}

	// Check the Mozilla report
	if yes, _ := moz.IsHTTPSonly(site.Host); yes {
		return TypeHTTPSok
	}
	return TypeHTTPSnok
}

func initAPIs() {
	irml = cryptcheck.NewClient(cryptcheck.Config{
		Log:     logLevel,
		Refresh: true,
		Timeout: 30,
	})

	fnImirhil = func(site ssllabs.Host) string {
		debug("  imirhil\n")
		score, err := irml.GetScore(site.Host)
		if err != nil {
			verbose("cryptcheck error: %s (%s)\n", site.Host, err.Error())
		}
		return score
	}

	moz, _ = observatory.NewClient(observatory.Config{
		Log:     logLevel,
		Timeout: 30,
	})

	fnMozilla = func(site ssllabs.Host) string {
		debug("  observatory\n")
		score, err := moz.GetGrade(site.Host)
		if err != nil {
			verbose("Mozilla error: %s (%s)\n", site.Host, err.Error())
		}
		return score
	}
}

func NewTLSSite(site ssllabs.Host) TLSSite {
	var current TLSSite

	initAPIs()
	if site.Endpoints == nil {
		verbose("Site %s has no endpoint\n", site.Host)
		current = TLSSite{
			Name:     site.Host,
			Contract: contracts[site.Host],
		}
	} else {
		endp := site.Endpoints[0]
		det := endp.Details

		fmt.Printf("Host: %s\n", site.Host)

		protos := []string{}
		for _, p := range det.Protocols {
			protos = append(protos, fmt.Sprintf("%sv%s", p.Name, p.Version))
		}

		// FIll in all details
		current = TLSSite{
			Name:       site.Host,
			Contract:   contracts[site.Host],
			Grade:      endp.Grade,
			CryptCheck: getGrade(site, fnImirhil),
			Mozilla:    getGrade(site, fnMozilla),
			Protocols:  strings.Join(protos, ","),
			PFS:        det.ForwardSecrecy >= 2,
			OCSP:       det.OcspStapling,
			HSTS:       checkHSTS(det),
			Sweet32:    checkSweet32(det),
			Type:       findServerType(site),
		}

		// Handle case where we have a DNS entry but no connection
		if len(site.Certs) != 0 {
			cert := site.Certs[0]
			current.DefKey = cert.KeySize == DefaultKeySize && cert.KeyAlg == DefaultAlg

			current.DefCA = checkIssuer(cert, DefaultIssuer)
			current.DefSig = cert.SigAlg == DefaultSig
			current.IsExpired = hasExpired(cert.NotAfter)
		}

		if len(det.CertChains) != 0 {
			current.PathIssues = det.CertChains[0].Issues != 0
		}
	}
	return current
}

func checkIssuer(cert ssllabs.Cert, ours *regexp.Regexp) bool {
	return ours.MatchString(cert.IssuerSubject)
}

func checkHSTS(det ssllabs.EndpointDetails) int64 {
	if det.HstsPolicy.Status == "present" {
		return det.HstsPolicy.MaxAge
	}
	return -1
}
