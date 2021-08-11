// Copyright (c) 2015-2021 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/minio/cli"
)

var (
	// CA root certificates, a nil value means system certs pool will be used
	globalRootCAs *x509.CertPool
)

//ClientURL used for endpoint validation/processing
type ClientURL struct {
	Type            ClientURLType
	Scheme          string
	Host            string
	Path            string
	SchemeSeparator string
	Separator       rune
}

// ClientURLType - enum of different url types
type ClientURLType int

// enum types
const (
	objectStorage = iota // MinIO and S3 compatible cloud storage
	fileSystem           // POSIX compatible file systems
)

func validateEndpoint(ctx *cli.Context, endpoint string) (bool, string, http.RoundTripper) {
	// Creates a parsed URL.
	targetURL := newClientURL(endpoint)
	// By default enable HTTPs.
	useTLS := true
	if targetURL.Scheme == "http" {
		useTLS = false
	}
	var transport http.RoundTripper

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 15 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost:   256,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		// Set this value so that the underlying transport round-tripper
		// doesn't try to auto decode the body of objects with
		// content-encoding set to `gzip`.
		//
		// Refer:
		//    https://golang.org/src/net/http/transport.go?h=roundTrip#L1843
		DisableCompression: true,
	}

	if useTLS {
		// Keep TLS config.
		tlsConfig := &tls.Config{
			RootCAs: globalRootCAs,
			// Can't use SSLv3 because of POODLE and BEAST
			// Can't use TLSv1.0 because of POODLE and BEAST using CBC cipher
			// Can't use TLSv1.1 because of RC4 cipher usage
			MinVersion: tls.VersionTLS12,
		}
		if ctx.Bool("insecure") || ctx.GlobalBool("insecure") {

			tlsConfig.InsecureSkipVerify = true
		}
		transport = tr

	}

	var url string
	if targetURL.Type == 1 {
		url = targetURL.Path
		useTLS = false
	} else {
		url = targetURL.Host
	}
	return useTLS, url, transport

}

//Parses the given URL
func newClientURL(urlStr string) *ClientURL {
	scheme, rest := getScheme(urlStr)
	if strings.HasPrefix(rest, "//") {
		// if rest has '//' prefix, skip them
		var authority string
		authority, rest = splitSpecial(rest[2:], "/", false)
		if rest == "" {
			rest = "/"
		}
		host := getHost(authority)
		if host != "" && (scheme == "http" || scheme == "https") {
			return &ClientURL{
				Scheme:          scheme,
				Type:            objectStorage,
				Host:            host,
				Path:            rest,
				SchemeSeparator: "://",
				Separator:       '/',
			}
		}
	}
	return &ClientURL{
		Type:      fileSystem,
		Path:      rest,
		Separator: filepath.Separator,
	}
}

func getScheme(rawurl string) (scheme, path string) {
	urlSplits := strings.Split(rawurl, "://")
	if len(urlSplits) == 2 {
		scheme, uri := urlSplits[0], "//"+urlSplits[1]
		// ignore numbers in scheme
		validScheme := regexp.MustCompile("^[a-zA-Z]+$")
		if uri != "" {
			if validScheme.MatchString(scheme) {
				return scheme, uri
			}
		}
	}
	return "", rawurl
}

func splitSpecial(s string, delimiter string, cutdelimiter bool) (string, string) {
	i := strings.Index(s, delimiter)
	if i < 0 {
		// if delimiter not found return as is.
		return s, ""
	}
	// if delimiter should be removed, remove it.
	if cutdelimiter {
		return s[0:i], s[i+len(delimiter):]
	}
	// return split strings with delimiter
	return s[0:i], s[i:]
}

func getHost(authority string) (host string) {
	i := strings.LastIndex(authority, "@")
	if i >= 0 {
		// TODO support, username@password style userinfo, useful for ftp support.
		return
	}
	return authority
}
