package modNodeJs

import (
	"os"
	"syscall"
	"time"
)

func fillFsStat(fi os.FileInfo, stat *FsFileState) {
	// Darwin and Linux version aren't the same for this fields.
	osInfo := fi.Sys().(*syscall.Stat_t)

	stat.Gid = osInfo.Gid
	stat.Uid = osInfo.Uid
	stat.Dev = osInfo.Dev
	stat.Rdev = osInfo.Rdev
	stat.Ino = osInfo.Ino
	stat.Mode = osInfo.Mode
	stat.Nlink = osInfo.Nlink
	stat.Blksize = osInfo.Blksize
	stat.Blocks = osInfo.Blocks

	stat.ATimeMs, _ = osInfo.Atimespec.Unix()
	stat.Atime = time.Unix(osInfo.Atimespec.Sec, osInfo.Atimespec.Nsec).UTC().Format(time.RFC3339Nano)

	stat.MTimeMs, _ = osInfo.Mtimespec.Unix()
	stat.Mtime = time.Unix(osInfo.Mtimespec.Sec, osInfo.Mtimespec.Nsec).UTC().Format(time.RFC3339Nano)

	stat.CTimeMs, _ = osInfo.Ctimespec.Unix()
	stat.Ctime = time.Unix(osInfo.Ctimespec.Sec, osInfo.Ctimespec.Nsec).UTC().Format(time.RFC3339Nano)

	stat.BirthtimeMs, _ = osInfo.Birthtimespec.Unix()
	stat.Birthtime = time.Unix(osInfo.Birthtimespec.Sec, osInfo.Birthtimespec.Nsec).UTC().Format(time.RFC3339Nano)
}
