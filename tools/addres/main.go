package main

import (
	"flag"
	"fmt"
	"os"
	"unsafe"

	"github.com/mkch/gw/util/icon"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"golang.org/x/sys/windows"
)

var addIconFlags addIconFlagSet
var addManifestFlags addManifestFlagSet

func main() {
	addIconFlags.Init()
	addManifestFlags.Init()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `addres is a tool for adding resource to exe file.
Usage:

    addres <command> [flags] exe_file
	
The commands are:

    icon	add an icon resource
    manifest	add a manifest resource

 Use "addres help <command>" for more information about a command.
 `)
	}
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		printError("not enough arguments.\n")
		flag.Usage()
		os.Exit(1)
	}
	cmd := args[0]
	switch cmd {
	case "help":
		if len(args) > 2 {
			printError("too many arguments.\n")
			flag.Usage()
			os.Exit(1)
		}
		if len(args) < 2 {
			flag.Usage()
			os.Exit(1)
		}
		commandUsage(args[1])
	case "icon":
		addIcon(args[1:])
	case "manifest":
		addManifest(args[1:])
	default:
		printError("unknown command %v.\n", cmd)
		flag.Usage()
		os.Exit(1)
	}
}

func commandUsage(cmd string) {
	switch cmd {
	case "icon":
		addIconFlags.Usage()
	case "manifest":
		addManifestFlags.Usage()
	}
}

func addManifest(arguments []string) {
	addManifestFlags.Parse(arguments)
	if addManifestFlags.File == "" {
		printError("no manifest file.\n")
		os.Exit(1)
	}
	args := addManifestFlags.Args()
	if len(args) < 1 {
		printError("not enough arguments.\n")
		addManifestFlags.Usage()
		os.Exit(1)
	}

	var err error
	var manifestData []byte
	if manifestData, err = os.ReadFile(addManifestFlags.File); err != nil {
		printError("%v\n", err)
		os.Exit(2)
	}

	exeFile := args[0]
	var exeFileBuf []win32.WCHAR
	win32util.CString(exeFile, &exeFileBuf)
	var update win32.HUPDATE
	if update, err = win32.BeginUpdateResourceW(&exeFileBuf[0], false); err != nil {
		printError("failed to update resource: %v\n", err)
		os.Exit(3)
	}
	defer func() {
		if err = win32.EndUpdateResourceW(update, false); err != nil {
			printError("failed to update resource: %v\n", err)
			os.Exit(3)
		}
	}()

	var module windows.Handle
	if module, err = windows.LoadLibraryEx(exeFile, 0, windows.LOAD_LIBRARY_AS_DATAFILE); err != nil {
		printError("failed to load %v: %v\n", exeFile, err)
		os.Exit(3)
	}
	//defer windows.FreeLibrary(module)

	// RT_MANIFEST must have id 1.
	const MANIFEST_RES_ID = 1

	if h, _ := win32.FindResourceW(win32.HMODULE(module), MANIFEST_RES_ID, win32.RT_MANIFEST); h != 0 {
		printError("manifest resource %v already exists.\n", MANIFEST_RES_ID)
		os.Exit(2)
	}

	if err = win32.UpdateResourceW(update,
		win32.RT_MANIFEST, MANIFEST_RES_ID, 0,
		unsafe.Pointer(&manifestData[0]), win32.DWORD(len(manifestData))); err != nil {
		printError("failed to update resource: %v\n", err)
		os.Exit(3)
	}
}

func addIcon(arguments []string) {
	addIconFlags.Parse(arguments)
	if addIconFlags.File == "" {
		printError("no icon file.\n")
		addIconFlags.Usage()
		os.Exit(1)
	}
	args := addIconFlags.Args()
	if len(args) < 1 {
		printError("not enough arguments.\n")
		addIconFlags.Usage()
		os.Exit(1)
	}

	var err error
	var iconReader *os.File
	if iconReader, err = os.Open(addIconFlags.File); err != nil {
		printError("%v\n", err)
		os.Exit(2)
	}
	defer iconReader.Close()
	var iconData *icon.Icon
	if iconData, err = icon.Read(iconReader); err != nil {
		printError("invalid icon file: %v.\n", err)
		os.Exit(2)
	}

	exeFile := args[0]
	var exeFileBuf []win32.WCHAR
	win32util.CString(exeFile, &exeFileBuf)
	var update win32.HUPDATE
	if update, err = win32.BeginUpdateResourceW(&exeFileBuf[0], false); err != nil {
		printError("failed to update resource: %v\n", err)
		os.Exit(3)
	}
	defer func() {
		if err = win32.EndUpdateResourceW(update, false); err != nil {
			printError("failed to update resource: %v\n", err)
			os.Exit(3)
		}
	}()

	var module windows.Handle
	if module, err = windows.LoadLibraryEx(exeFile, 0, windows.LOAD_LIBRARY_AS_DATAFILE); err != nil {
		printError("failed to load %v: %v\n", exeFile, err)
		os.Exit(3)
	}
	defer windows.FreeLibrary(module)

	var iconIDs []uint16
	// Add icons.
	for i := range iconData.Images {
		id := uintptr(addIconFlags.FirstIconID + i)
		if h, _ := win32.FindResourceW(win32.HMODULE(module), id, win32.RT_ICON); h != 0 {
			printError("icon resource %v already exists.\n", id)
			os.Exit(2)
		}
		if err = win32.UpdateResourceW(update, win32.RT_ICON, id, 0,
			unsafe.Pointer(&iconData.Images[i].Data[0]), win32.DWORD(len(iconData.Images[i].Data))); err != nil {
			printError("failed to update resource: %v\n", err)
			os.Exit(3)
		}
		iconIDs = append(iconIDs, uint16(id))
	}
	// Add icon group.
	if h, _ := win32.FindResourceW(win32.HMODULE(module), uintptr(addIconFlags.GroupID), win32.RT_GROUP_ICON); h != 0 {
		printError("group icon resource %v already exists.\n", addIconFlags.GroupID)
		os.Exit(2)
	}
	groupIconData := make([]byte, unsafe.Sizeof(icon.IconDirHeader{})+uintptr(len(iconData.Images))*unsafe.Sizeof(icon.GrpIconDirEntry{}))
	hdr := (*icon.IconDirHeader)(unsafe.Pointer(&groupIconData[0]))
	*hdr.Type() = iconData.Type
	*hdr.Count() = uint16(len(iconData.Images))
	for i := range iconData.Images {
		entry := (*icon.GrpIconDirEntry)(unsafe.Pointer(uintptr(unsafe.Pointer(hdr)) + unsafe.Sizeof(icon.IconDirHeader{}) + uintptr(i)*unsafe.Sizeof(icon.GrpIconDirEntry{})))
		*entry.Width() = *iconData.Images[i].Entry.Width()
		*entry.Height() = *iconData.Images[i].Entry.Height()
		*entry.ColorCount() = *iconData.Images[i].Entry.ColorCount()
		*entry.Planes() = *iconData.Images[i].Entry.Planes()
		*entry.BitCount() = *iconData.Images[i].Entry.BitCount()
		*entry.BytesInRes() = *iconData.Images[i].Entry.BytesInRes()
		*entry.ID() = iconIDs[i]
	}
	if err = win32.UpdateResourceW(update, win32.RT_GROUP_ICON, uintptr(addIconFlags.GroupID), 0, unsafe.Pointer(&groupIconData[0]), win32.DWORD(len(groupIconData))); err != nil {
		printError("failed to update resource: %v\n", err)
		os.Exit(3)
	}
}

func printError(format string, a ...any) {
	fmt.Fprintf(flag.CommandLine.Output(), format, a...)
}
