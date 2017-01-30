package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
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
			URL string `xml:",href"`
			Rel string `xml:",rel"`
		} `xml:"link"`
		Updated   string `xml:"updated"`
		Published string `xml:"published"`
	}
)

const (
	atomTimeFormat = "2006-01-02T15:04:25.000000-07:00"
	atomHeader     = `<?xml version='1.0' encoding='UTF-8'?>
<feed xmlns="http://www.w3.org/2005/Atom" xml:lang="en">
	<id>http://npmjs.org/</id>
	<title>Node Package Feed</title>
	<subtitle>Provides a feed of updates to your npm dependencies.</subtitle>
	<link href="/feed.rss" rel="self"/>
	<generator url="https://github.com/Benzinga/npm-feed">npm-feed</generator>
	<updated>%s</updated>
`
	atomFooter = `</feed>
`
)

func atom(rels []Release) []byte {
	// Prepare buffer.
	b := bytes.Buffer{}
	b.Grow(65535)

	// Get updated time.
	updated := time.Now().Format(atomTimeFormat)

	// Write header.
	b.WriteString(fmt.Sprintf(atomHeader, updated))

	// Format entries.
	entries := make([]AtomEntry, len(rels))
	for i, rel := range rels {
		guid := sha1.Sum([]byte(fmt.Sprintf("%s#%s#%d", rel.Name, rel.Version, rel.Date.Unix())))
		entries[i].ID = hex.EncodeToString(guid[:])
		entries[i].Title = rel.Name + " " + rel.Version
		entries[i].Link.URL = getNodePackageURL(rel.Name)
		entries[i].Link.Rel = "alternate"
		entries[i].Updated = rel.Date.Format(rssTimeFormat)
		entries[i].Published = rel.Date.Format(rssTimeFormat)
	}

	// Marshal items.
	enc := xml.NewEncoder(&b)
	enc.Indent("\t", "\t")
	err := enc.Encode(entries)
	if err != nil {
		panic(err)
	}

	// Write footer.
	b.WriteString(rssFooter)

	return b.Bytes()
}
