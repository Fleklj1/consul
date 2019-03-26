// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// This program generates internet protocol constants by reading IANA
// protocol registries.
//
// Usage:
//	go run gentest.go > iana_test.go
package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"go/format"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var registries = []struct {
	url   string
	parse func(io.Writer, io.Reader) error
}{
	{
		"http://www.iana.org/assignments/dscp-registry/dscp-registry.xml",
		parseDSCPRegistry,
	},
	{
		"http://www.iana.org/assignments/ipv4-tos-byte/ipv4-tos-byte.xml",
		parseTOSTCByte,
	},
}

func main() {
	var bb bytes.Buffer
	fmt.Fprintf(&bb, "// go run gentest.go\n")
	fmt.Fprintf(&bb, "// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT\n\n")
	fmt.Fprintf(&bb, "package ipv6_test\n\n")
	for _, r := range registries {
		resp, err := http.Get(r.url)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "got HTTP status code %v for %v\n", resp.StatusCode, r.url)
			os.Exit(1)
		}
		if err := r.parse(&bb, resp.Body); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprintf(&bb, "\n")
	}
	b, err := format.Source(bb.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Stdout.Write(b)
}

func parseDSCPRegistry(w io.Writer, r io.Reader) error {
	dec := xml.NewDecoder(r)
	var dr dscpRegistry
	if err := dec.Decode(&dr); err != nil {
		return err
	}
	drs := dr.escape()
	fmt.Fprintf(w, "// %s, Updated: %s\n", dr.Title, dr.Updated)
	fmt.Fprintf(w, "const (\n")
	for _, dr := range drs {
		fmt.Fprintf(w, "DiffServ%s = %#x", dr.Name, dr.Value)
		fmt.Fprintf(w, "// %s\n", dr.OrigName)
	}
	fmt.Fprintf(w, ")\n")
	return nil
}

type dscpRegistry struct {
	XMLName     xml.Name `xml:"registry"`
	Title       string   `xml:"title"`
	Updated     string   `xml:"updated"`
	Note        string   `xml:"note"`
	RegTitle    string   `xml:"registry>title"`
	PoolRecords []struct {
		Name  string `xml:"name"`
		Space string `xml:"space"`
	} `xml:"registry>record"`
	Records []struct {
		Name  string `xml:"name"`
		Space string `xml:"space"`
	} `xml:"registry>registry>record"`
}

type canonDSCPRecord struct {
	OrigName string
	Name     string
	Value    int
}

func (drr *dscpRegistry) escape() []canonDSCPRecord {
	drs := make([]canonDSCPRecord, len(drr.Records))
	sr := strings.NewReplacer(
		"+", "",
		"-", "",
		"/", "",
		".", "",
		" ", "",
	)
	for i, dr := range drr.Records {
		s := strings.TrimSpace(dr.Name)
		drs[i].OrigName = s
		drs[i].Name = sr.Replace(s)
		n, err := strconv.ParseUint(dr.Space, 2, 8)
		if err != nil {
			continue
		}
		drs[i].Value = int(n) << 2
	}
	return drs
}

func parseTOSTCByte(w io.Writer, r io.Reader) error {
	dec := xml.NewDecoder(r)
	var ttb tosTCByte
	if err := dec.Decode(&ttb); err != nil {
		return err
	}
	trs := ttb.escape()
	fmt.Fprintf(w, "// %s, Updated: %s\n", ttb.Title, ttb.Updated)
	fmt.Fprintf(w, "const (\n")
	for _, tr := range trs {
		fmt.Fprintf(w, "%s = %#x", tr.Keyword, tr.Value)
		fmt.Fprintf(w, "// %s\n", tr.OrigKeyword)
	}
	fmt.Fprintf(w, ")\n")
	return nil
}

type tosTCByte struct {
	XMLName  xml.Name `xml:"registry"`
	Title    string   `xml:"title"`
	Updated  string   `xml:"updated"`
	Note     string   `xml:"note"`
	RegTitle string   `xml:"registry>title"`
	Records  []struct {
		Binary  string `xml:"binary"`
		Keyword string `xml:"keyword"`
	} `xml:"registry>record"`
}

type canonTOSTCByteRecord struct {
	OrigKeyword string
	Keyword     string
	Value       int
}

func (ttb *tosTCByte) escape() []canonTOSTCByteRecord {
	trs := make([]canonTOSTCByteRecord, len(ttb.Records))
	sr := strings.NewReplacer(
		"Capable", "",
		"(", "",
		")", "",
		"+", "",
		"-", "",
		"/", "",
		".", "",
		" ", "",
	)
	for i, tr := range ttb.Records {
		s := strings.TrimSpace(tr.Keyword)
		trs[i].OrigKeyword = s
		ss := strings.Split(s, " ")
		if len(ss) > 1 {
			trs[i].Keyword = strings.Join(ss[1:], " ")
		} else {
			trs[i].Keyword = ss[0]
		}
		trs[i].Keyword = sr.Replace(trs[i].Keyword)
		n, err := strconv.ParseUint(tr.Binary, 2, 8)
		if err != nil {
			continue
		}
		trs[i].Value = int(n)
	}
	return trs
}
