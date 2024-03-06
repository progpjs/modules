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

	stat.ATimeMs, _ = osInfo.Atim.Unix()
	stat.Atime = time.Unix(osInfo.Atim.Sec, osInfo.Atim.Nsec).UTC().Format(time.RFC3339Nano)

	stat.MTimeMs, _ = osInfo.Mtim.Unix()
	stat.Mtime = time.Unix(osInfo.Mtim.Sec, osInfo.Mtim.Nsec).UTC().Format(time.RFC3339Nano)

	stat.CTimeMs, _ = osInfo.Ctim.Unix()
	stat.Ctime = time.Unix(osInfo.Ctim.Sec, osInfo.Ctim.Nsec).UTC().Format(time.RFC3339Nano)

	// Birth time doesn't exist so take create time.
	stat.BirthtimeMs, _ = osInfo.Ctim.Unix()
	stat.Birthtime = time.Unix(osInfo.Ctim.Sec, osInfo.Ctim.Nsec).UTC().Format(time.RFC3339Nano)
}
