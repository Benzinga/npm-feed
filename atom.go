package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"
)

type (
	// AtomEntry represents an <entry> tag in Atom.
	AtomEntry struct {
		XMLName struct{} `xml:"entry"`
		ID      string   `xml:"id"`
		Title   string   `xml:"title"`
		Link    struct {
			URL string `xml:"href,attr"`
			Rel string `xml:"rel,attr"`
		} `xml:"link"`
		Updated   string `xml:"updated"`
		Published string `xml:"published"`
	}
)

const (
	atomTimeFormat = time.RFC3339
	atomHeader     = `<?xml version='1.0' encoding='UTF-8'?>
<feed xmlns="http://www.w3.org/2005/Atom" xml:lang="en">
	<id>http://npmjs.org/</id>
	<title>Node Package Feed</title>
	<subtitle>Provides a feed of updates to your npm dependencies.</subtitle>
	<link href="%s" rel="self"/>
	<generator uri="https://github.com/Benzinga/npm-feed">npm-feed</generator>
	<updated>%s</updated>
`
	atomFooter = `</feed>
`
)

func atom(rels []Release, uri string) []byte {
	// Prepare buffer.
	b := bytes.Buffer{}
	b.Grow(65535)

	// Get updated time.
	updated := time.Now().Format(atomTimeFormat)

	// Write header.
	b.WriteString(fmt.Sprintf(atomHeader, uri, updated))

	// Format entries.
	entries := make([]AtomEntry, len(rels))
	for i, rel := range rels {
		entries[i].ID = getNodePackageURL(rel.Name)
		entries[i].Title = rel.Name + " " + rel.Version
		entries[i].Link.URL = getNodePackageURL(rel.Name)
		entries[i].Link.Rel = "alternate"
		entries[i].Updated = rel.Date.Format(atomTimeFormat)
		entries[i].Published = rel.Date.Format(atomTimeFormat)
	}

	// Marshal items.
	enc := xml.NewEncoder(&b)
	enc.Indent("\t", "\t")
	err := enc.Encode(entries)
	if err != nil {
		panic(err)
	}

	// Write footer.
	b.WriteString(atomFooter)

	return b.Bytes()
}
