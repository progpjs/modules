package modNodeJs

import (
	"os"
	"syscall"
	"time"
)

func fillFsStat(fi os.FileInfo, stat *FsFileState) {
	// Darwin and Linux version aren't the same for this fields.
	osInfo := fi.Sys().(*syscall.Stat_t)

	stat.Gid = uint64(osInfo.Gid)
	stat.Uid = uint64(osInfo.Uid)
	stat.Dev = uint64(osInfo.Dev)
	stat.Rdev = uint64(osInfo.Rdev)
	stat.Ino = osInfo.Ino
	stat.Mode = osInfo.Mode
	stat.Nlink = uint64(osInfo.Nlink)
	stat.Blksize = uint64(osInfo.Blksize)
	stat.Blocks = uint64(osInfo.Blocks)

	var v int64

	v, _ = osInfo.Atimespec.Unix()
	stat.ATimeMs = uint64(v)
	stat.Atime = time.Unix(osInfo.Atimespec.Sec, osInfo.Atimespec.Nsec).UTC().Format(time.RFC3339Nano)

	v, _ = osInfo.Mtimespec.Unix()
	stat.MTimeMs = uint64(v)
	stat.Mtime = time.Unix(osInfo.Mtimespec.Sec, osInfo.Mtimespec.Nsec).UTC().Format(time.RFC3339Nano)

	v, _ = osInfo.Ctimespec.Unix()
	stat.CTimeMs = uint64(v)
	stat.Ctime = time.Unix(osInfo.Ctimespec.Sec, osInfo.Ctimespec.Nsec).UTC().Format(time.RFC3339Nano)

	v, _ = osInfo.Birthtimespec.Unix()
	stat.BirthtimeMs = uint64(v)
	stat.Birthtime = time.Unix(osInfo.Birthtimespec.Sec, osInfo.Birthtimespec.Nsec).UTC().Format(time.RFC3339Nano)
}
