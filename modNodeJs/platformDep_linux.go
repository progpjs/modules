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
	stat.Dev = osInfo.Dev
	stat.Rdev = osInfo.Rdev
	stat.Ino = osInfo.Ino
	stat.Mode = uint16(osInfo.Mode)
	stat.Nlink = osInfo.Nlink
	stat.Blksize = uint64(osInfo.Blksize)
	stat.Blocks = uint64(osInfo.Blocks)

	var v int64

	v, _ = osInfo.Atim.Unix()
	stat.ATimeMs = uint64(v)
	stat.Atime = time.Unix(osInfo.Atim.Sec, osInfo.Atim.Nsec).UTC().Format(time.RFC3339Nano)

	v, _ = osInfo.Mtim.Unix()
	stat.MTimeMs = uint64(v)
	stat.Mtime = time.Unix(osInfo.Mtim.Sec, osInfo.Mtim.Nsec).UTC().Format(time.RFC3339Nano)

	v, _ = osInfo.Ctim.Unix()
	stat.CTimeMs = uint64(v)
	stat.Ctime = time.Unix(osInfo.Ctim.Sec, osInfo.Ctim.Nsec).UTC().Format(time.RFC3339Nano)

	v, _ = osInfo.Ctim.Unix()
	stat.BirthtimeMs = uint64(v)
	stat.Birthtime = time.Unix(osInfo.Ctim.Sec, osInfo.Ctim.Nsec).UTC().Format(time.RFC3339Nano)
}
