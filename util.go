package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"
)

// Go doesn't really have a defacto package for handling environment variables
// the way it does for flags, so here are some small helpers to bridge the gap
// to a degree.

func envifyFlag(name string) string {
	return strings.Replace(strings.ToUpper(name), "-", "_", -1)
}

func eStringVar(p *string, name, desc string) {
	flag.StringVar(p, name, *p, desc)

	name = envifyFlag(name)
	if v := os.Getenv(name); v != "" {
		*p = v
	}
}

func eDurationVar(p *time.Duration, name, desc string) {
	flag.DurationVar(p, name, *p, desc)

	name = envifyFlag(name)
	if v := os.Getenv(name); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalln("env:", name+":", err)
		}
		*p = d
	}
}
