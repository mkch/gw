package main

import (
	"flag"
	"fmt"
)

type addIconFlagSet struct {
	flags       *flag.FlagSet
	File        string
	GroupID     int
	FirstIconID int
}

func (f *addIconFlagSet) Init() {
	f.flags = flag.NewFlagSet("icon", flag.ExitOnError)
	f.flags.StringVar(&f.File, "res", "icon.ico", "The icon file(*.ico)")
	f.flags.IntVar(&f.GroupID, "group-id", 1, "The id of icon group(RT_GROUP_ICON) resource")
	f.flags.IntVar(&f.FirstIconID, "icon-id", 1, "The first id of icon(RT_ICON) resource")
}

func (f *addIconFlagSet) Parse(arguments []string) {
	err := f.flags.Parse(arguments)
	if err != nil {
		panic(err)
	}
}

func (f *addIconFlagSet) Usage() {
	fmt.Fprintln(f.flags.Output(), "usage: addres icon [flags] path_of_exe_file")
	f.flags.PrintDefaults()
}

func (f *addIconFlagSet) Args() []string {
	return f.flags.Args()
}

type addManifestFlagSet struct {
	flags *flag.FlagSet
	File  string
}

func (f *addManifestFlagSet) Init() {
	f.flags = flag.NewFlagSet("manifest", flag.ExitOnError)
	f.flags.StringVar(&f.File, "res", "manifest.xml", "The icon file(*.ico)")
}

func (f *addManifestFlagSet) Parse(arguments []string) {
	err := f.flags.Parse(arguments)
	if err != nil {
		panic(err)
	}
}

func (f *addManifestFlagSet) Usage() {
	fmt.Fprintln(f.flags.Output(), "usage: addres manifest [flags] path_of_exe_file")
	f.flags.PrintDefaults()
}

func (f *addManifestFlagSet) Args() []string {
	return f.flags.Args()
}
