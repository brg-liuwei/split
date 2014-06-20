package syncio

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
)

type SyncFile struct {
	abspath string
	fp      *os.File
	ref     int32
	mutex   sync.Mutex
}

type syncMap struct {
	m     map[string]*SyncFile
	mutex sync.RWMutex
}

var smap syncMap = syncMap{m: make(map[string]*SyncFile)}
var syncMutex sync.Mutex

func (smap *syncMap) add(abspath string, sfile *SyncFile) {
	smap.mutex.Lock()
	defer smap.mutex.Unlock()
	smap.m[abspath] = sfile
}

func (smap *syncMap) del(abspath string) {
	smap.mutex.Lock()
	defer smap.mutex.Unlock()
	delete(smap.m, abspath)
}

func (smap *syncMap) get(abspath string) *SyncFile {
	smap.mutex.RLock()
	defer smap.mutex.RUnlock()
	sfile, ok := smap.m[abspath]
	if ok {
		return sfile
	}
	return nil
}

func NewSyncFile(path string) *SyncFile {
	syncMutex.Lock()
	defer syncMutex.Unlock()
	abspath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("NewSyncFile get abspath err: ", err)
		return nil
	}
	sfile := smap.get(abspath)
	if sfile != nil {
		atomic.AddInt32(&sfile.ref, 1)
		return sfile
	}
	f, e := os.OpenFile(abspath, os.O_CREATE|os.O_RDWR, 0666)
	if e != nil {
		fmt.Println("NewSyncFile open file err: ", e, ", file: ", abspath)
		return nil
	}
	sfile = new(SyncFile)
	sfile.abspath = abspath
	sfile.fp = f
	sfile.ref = 0
	smap.add(abspath, sfile)
	return sfile
}

func (s *SyncFile) Read(buf []byte, off int64) (n int, err error) {
	return syscall.Pread(int(s.fp.Fd()), buf, off)
}

func (s *SyncFile) Write(buf []byte, off int64) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return syscall.Pwrite(int(s.fp.Fd()), buf, off)
}

func (s *SyncFile) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	smap.del(s.abspath)
	if atomic.AddInt32(&s.ref, -1) == 0 {
		return s.fp.Close()
	}
	return nil
}
