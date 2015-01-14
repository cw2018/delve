package proctl

/*
#include <stdlib.h>
#include <libproc.h>

char *
findExecutable(int pid) {
	char *pathbuf = (char *)malloc(sizeof(char)*PROC_PIDPATHINFO_MAXSIZE);
	proc_pidpath(pid, pathbuf, sizeof(pathbuf));
	return pathbuf;
}

static const unsigned char info_plist[]
__attribute__ ((section ("__TEXT,__info_plist"),used)) =
  "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
  "<!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\""
  " \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n"
  "<plist version=\"1.0\">\n"
  "<dict>\n"
  "  <key>CFBundleIdentifier</key>\n"
  "  <string>org.dlv</string>\n"
  "  <key>CFBundleName</key>\n"
  "  <string>delve</string>\n"
  "  <key>CFBundleVersion</key>\n"
  "  <string>1.0</string>\n"
  "  <key>SecTaskAccess</key>\n"
  "  <array>\n"
  "    <string>allowed</string>\n"
  "    <string>debug</string>\n"
  "  </array>\n"
  "</dict>\n"
  "</plist>\n";
*/
import "C"
import (
	"debug/elf"
	"sync"
	"syscall"
)

func (dbp *DebuggedProcess) addThread(tid int) (*ThreadContext, error) {
	dbp.Threads[tid] = &ThreadContext{
		Id:      tid,
		Process: dbp,
	}
	return dbp.Threads[tid], nil
}

// Finds the executable and then uses it
// to parse the following information:
// * Dwarf .debug_frame section
// * Dwarf .debug_line section
// * Go symbol table.
func (dbp *DebuggedProcess) LoadInformation() error {
	var (
		wg  sync.WaitGroup
		exe *elf.File
		err error
	)

	exe, err = dbp.findExecutable()
	if err != nil {
		return err
	}

	wg.Add(2)
	go dbp.parseDebugFrame(exe, &wg)
	go dbp.obtainGoSymbols(exe, &wg)
	wg.Wait()

	return nil
}

func stopped(pid int) bool {
	return &ps, nil
}

func (dbp *DebuggedProcess) findExecutable() (string, error) {
	pathptr, err := C.findExecutable(dbp.Pid)
	path := string(pathptr)
	C.free(path)
	if err != syscall.Errno(0) {
		return path, err
	}
	return path, nil
}
