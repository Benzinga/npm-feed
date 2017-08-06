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
	// RSSItem represents an <item> tag in RSS.
	RSSItem struct {
		XMLName struct{} `xml:"item"`
		Title   string   `xml:"title"`
		PubDate string   `xml:"pubDate"`
		GUID    struct {
			Value       string `xml:",chardata"`
			IsPermaLink bool   `xml:"isPermaLink,attr"`
		} `xml:"guid"`
		Link string `xml:"link"`
	}
)

const (
	rssTimeFormat = "Mon, 02 Jan 2006 15:04:05 -0700"
	rssHeader     = `<?xml version='1.0' encoding='UTF-8'?>
<rss xmlns:atom="http://www.w3.org/2005/Atom" xmlns:content="http://purl.org/rss/1.0/modules/content/" version="2.0">
	<channel>
		<title>Node Package Feed</title>
		<description>Provides a feed of updates to your npm dependencies.</description>
		<link>%s</link>
		<atom:link href="%s" rel="self"/>
		<docs>http://www.rssboard.org/rss-specification</docs>
		<generator>npm-feed</generator>
		<language>en</language>
		<lastBuildDate>%s</lastBuildDate>
	`
	rssFooter = `	</channel>
</rss>
`
)

func rss(rels []Release, uri string) []byte {
	// Prepare buffer.
	b := bytes.Buffer{}
	b.Grow(65535)

	// Get updated time.
	updated := time.Now().Format(rssTimeFormat)

	// Write header.
	b.WriteString(fmt.Sprintf(rssHeader, uri, uri, updated))

	// Format items.
	items := make([]RSSItem, len(rels))
	for i, rel := range rels {
		guid := sha1.Sum([]byte(fmt.Sprintf("%s#%s#%d", rel.Name, rel.Version, rel.Date.Unix())))
		items[i].Title = rel.Name + " " + rel.Version
		items[i].Link = getNodePackageURL(rel.Name)
		items[i].GUID.Value = hex.EncodeToString(guid[:])
		items[i].GUID.IsPermaLink = false
		items[i].PubDate = rel.Date.Format(rssTimeFormat)
	}

	// Marshal items.
	enc := xml.NewEncoder(&b)
	enc.Indent("\t\t", "\t")
	err := enc.Encode(items)
	if err != nil {
		panic(err)
	}

	// Write footer.
	b.WriteString(rssFooter)

	return b.Bytes()
}
