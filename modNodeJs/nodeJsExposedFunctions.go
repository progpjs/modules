/*
 * (C) Copyright 2024 Johan Michel PIQUET, France (https://johanpiquet.fr/).
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package modNodeJs

import (
	"encoding/json"
	"errors"
	"github.com/progpjs/libProgpScripts"
	"github.com/progpjs/progpAPI"
	"io"
	"os"
	"runtime"
	"syscall"
	"time"
)

func registerExportedFunctions() {
	rg := libProgpScripts.GetFunctionRegistry()
	myMod := rg.UseGoNamespace("github.com/progpjs/modules/modNodeJs")

	//region node:process

	modProcess := myMod.UseCustomGroup("nodejsModProcess")
	modProcess.AddFunction("kill", "JsProcessKill", JsProcessKill)
	modProcess.AddFunction("cwd", "JsProcessCwd", JsProcessCwd)
	modProcess.AddFunction("env", "JsProcessEnv", JsProcessEnv)
	modProcess.AddFunction("arch", "JsProcessArch", JsProcessArch)
	modProcess.AddFunction("platform", "JsProcessPlatform", JsProcessPlatform)
	modProcess.AddFunction("argv", "JsProcessArgV", JsProcessArgV)
	modProcess.AddFunction("exit", "JsProcessExit", JsProcessExit)
	modProcess.AddFunction("pid", "JsProcessPID", JsProcessPID)
	modProcess.AddFunction("ppid", "JsProcessParentPID", JsProcessParentPID)
	modProcess.AddFunction("chdir", "JsProcessChDir", JsProcessChDir)
	modProcess.AddFunction("getuid", "JsProcessGetUid", JsProcessGetUid)
	modProcess.AddAsyncFunction("nextTick", "JsProcessNextTickAsync", JsProcessNextTickAsync)

	//endregion

	//region node:os

	modOS := myMod.UseCustomGroup("nodejsModOS")
	modOS.AddFunction("homeDir", "JsOsHomeDir", JsOsHomeDir)
	modOS.AddFunction("hostName", "JsOsHostName", JsOsHostName)
	modOS.AddFunction("tempDir", "JsOsTempDir", JsOsTempDir)

	//endregion

	//region node:fs

	modFS := myMod.UseCustomGroup("nodejsModFS")

	//region Sync

	modFS.AddFunction("existsSync", "JsFsExistsSync", JsFsExistsSync)
	modFS.AddFunction("statSync", "JsFsStatSync", JsFsStatSync)
	modFS.AddFunction("accessSync", "JsFsAccessSync", JsFsAccessSync)
	modFS.AddFunction("chmodSync", "JsChmodSync", JsChmodSync)
	modFS.AddFunction("chownSync", "JsChownSync", JsChownSync)
	modFS.AddFunction("truncateSync", "JsTruncateSync", JsTruncateSync)
	modFS.AddFunction("readFileUtf8Sync", "JsReadFileUtf8Sync", JsReadFileUtf8Sync)
	modFS.AddFunction("readFileBytesSync", "JsReadFileBytesSync", JsReadFileBytesSync)
	modFS.AddFunction("copyFileSync", "JsCopyFileSync", JsCopyFileSync)
	modFS.AddFunction("linkSync", "JsLinkSync", JsLinkSync)
	modFS.AddFunction("symlinkSync", "JsSymLinkSync", JsSymLinkSync)
	modFS.AddFunction("unlinkSync", "JsUnlink", JsUnlink)

	//endregion

	//endregion
}

//region node:process	(nodejsModProcess)

func JsProcessCwd() string {
	cwd, _ := os.Getwd()
	return cwd
}

func JsProcessEnv() progpAPI.StringBuffer {
	res := os.Environ()
	b, _ := json.Marshal(res)
	return b
}

func JsProcessArch() string {
	// Apple MAC: arm64
	return runtime.GOARCH
}

func JsProcessPlatform() string {
	return runtime.GOOS
}

func JsProcessArgV() []string {
	return os.Args
}

func JsProcessExit(code int) {
	os.Exit(code)
}

func JsProcessPID() int {
	return os.Getpid()
}

func JsProcessParentPID() int {
	return os.Getppid()
}

func JsProcessChDir(dir string) error {
	return os.Chdir(dir)
}

func JsProcessGetUid() int {
	return os.Getuid()
}

func JsProcessNextTickAsync(fct progpAPI.ScriptFunction) {
	progpAPI.SafeGoRoutine(func() {
		fct.CallWithUndefined()
	})
}

func JsProcessKill(pid int, signal int) error {
	err := syscall.Kill(pid, syscall.Signal(signal))

	if signal == 0 {
		// Don't throw error fi signal is 0
		// which allows testing if process exists.
		// It's a node.js special case.
		//
		return nil
	}

	return err
}

//endregion

//region node:os (nodejsModOS)

func JsOsHomeDir() (string, error) {
	dirname, err := os.UserHomeDir()
	return dirname, err
}

func JsOsHostName() (string, error) {
	name, err := os.Hostname()
	return name, err
}

func JsOsTempDir() string {
	return os.TempDir()
}

//endregion

//region node:fs (nodejsModFS)

//region Structs & Enums

type FsFileState struct {
	Mode    uint16 `json:"mode"`
	Size    int64  `json:"size"`
	Gid     uint32 `json:"gid"`
	Uid     uint32 `json:"uid"`
	Dev     int32  `json:"dev"`
	Rdev    int32  `json:"rdev"`
	Ino     uint64 `json:"ino"`
	Nlink   uint16 `json:"nlink"`
	Blksize int32  `json:"blksize"`
	Blocks  int64  `json:"blocks"`

	ATimeMs int64  `json:"atimeMs"`
	Atime   string `json:"atime"`

	MTimeMs int64  `json:"mtimeMs"`
	Mtime   string `json:"mtime"`

	CTimeMs int64  `json:"ctimeMs"`
	Ctime   string `json:"ctime"`

	BirthtimeMs int64  `json:"birthtimeMs"`
	Birthtime   string `json:"birthtime"`
}

type FsConst uint32

const (
	FSCONST__UV_FS_SYMLINK_DIR               FsConst = 1
	FSCONST__UV_FS_SYMLINK_JUNCTION                  = 2
	FSCONST__O_RDONLY                                = 0
	FSCONST__O_WRONLY                                = 1
	FSCONST__O_RDWR                                  = 2
	FSCONST__UV_DIRENT_UNKNOWN                       = 0
	FSCONST__UV_DIRENT_FILE                          = 1
	FSCONST__UV_DIRENT_DIR                           = 2
	FSCONST__UV_DIRENT_LINK                          = 3
	FSCONST__UV_DIRENT_FIFO                          = 4
	FSCONST__UV_DIRENT_SOCKET                        = 5
	FSCONST__UV_DIRENT_CHAR                          = 6
	FSCONST__UV_DIRENT_BLOCK                         = 7
	FSCONST__EXTENSIONLESS_FORMAT_JAVASCRIPT         = 0
	FSCONST__EXTENSIONLESS_FORMAT_WASM               = 1
	FSCONST__S_IFMT                                  = 61440
	FSCONST__S_IFREG                                 = 32768
	FSCONST__S_IFDIR                                 = 16384
	FSCONST__S_IFCHR                                 = 8192
	FSCONST__S_IFBLK                                 = 24576
	FSCONST__S_IFIFO                                 = 4096
	FSCONST__S_IFLNK                                 = 40960
	FSCONST__S_IFSOCK                                = 49152
	FSCONST__O_CREAT                                 = 512
	FSCONST__O_EXCL                                  = 2048
	FSCONST__UV_FS_O_FILEMAP                         = 0
	FSCONST__O_NOCTTY                                = 131072
	FSCONST__O_TRUNC                                 = 1024
	FSCONST__O_APPEND                                = 8
	FSCONST__O_DIRECTORY                             = 1048576
	FSCONST__O_NOFOLLOW                              = 256
	FSCONST__O_SYNC                                  = 128
	FSCONST__O_DSYNC                                 = 4194304
	FSCONST__O_SYMLINK                               = 2097152
	FSCONST__O_NONBLOCK                              = 4
	FSCONST__S_IRWXU                                 = 448
	FSCONST__S_IRUSR                                 = 256
	FSCONST__S_IWUSR                                 = 128
	FSCONST__S_IXUSR                                 = 64
	FSCONST__S_IRWXG                                 = 56
	FSCONST__S_IRGRP                                 = 32
	FSCONST__S_IWGRP                                 = 16
	FSCONST__S_IXGRP                                 = 8
	FSCONST__S_IRWXO                                 = 7
	FSCONST__S_IROTH                                 = 4
	FSCONST__S_IWOTH                                 = 2
	FSCONST__S_IXOTH                                 = 1
	FSCONST__F_OK                                    = 0
	FSCONST__R_OK                                    = 4
	FSCONST__W_OK                                    = 2
	FSCONST__X_OK                                    = 1
	FSCONST__UV_FS_COPYFILE_EXCL                     = 1
	FSCONST__COPYFILE_EXCL                           = 1
	FSCONST__UV_FS_COPYFILE_FICLONE                  = 2
	FSCONST__COPYFILE_FICLONE                        = 2
	FSCONST__UV_FS_COPYFILE_FICLONE_FORCE            = 4
	FSCONST__COPYFILE_FICLONE_FORCE                  = 4
)

//endregion

//region Sync

func JsFsExistsSync(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func JsFsStatSync(path string, throwErrorIfMissing bool) (*FsFileState, error) {
	info, err := os.Stat(path)

	if err != nil {
		if throwErrorIfMissing {
			return nil, err
		}

		return nil, nil
	}

	osInfo, isUnixFS := info.Sys().(*syscall.Stat_t)

	stat := &FsFileState{}

	stat.Size = info.Size()

	if isUnixFS {
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

	return stat, nil
}

func JsFsAccessSync(path string, mode int) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// https://phoenixnap.com/kb/what-is-umask
	// User Group Other

	perm := info.Mode().Perm()

	// 0444  =>  4: user can read, 4: group can read, 4: others can read
	// 0444  =>  can read
	// 0555	 =>  can read & execute
	// 0666	 =>  can read & write
	// 0777	 =>  can read & write & execute
	// 0222	 =>  can write
	// 0333	 =>  can write & execute
	// 0111  =>  can execute

	canRead := (perm&0444 == 0444) || (perm&0555 == 0555) || (perm&0666 == 0666) || (perm&0777 == 0777)
	canWrite := (perm&0222 == 0222) || (perm&0666 == 0666) || (perm&0777 == 0777) || (perm&0333 == 0333)
	canExecute := (perm&0111 == 0111) || (perm&0555 == 0555) || (perm&0777 == 0777) || (perm&0333 == 0333)

	if mode == FSCONST__F_OK {
		// F_OK allows testing if file exists.
		return nil
	} else {
		if mode&FSCONST__R_OK == FSCONST__R_OK {
			if !canRead {
				return errors.New("can't read")
			}
		}

		if mode&FSCONST__W_OK == FSCONST__W_OK {
			if !canWrite {
				return errors.New("can't write")
			}
		}

		if mode&FSCONST__X_OK == FSCONST__X_OK {
			if !canExecute {
				return errors.New("can't execute")
			}
		}
	}

	return nil
}

func JsChmodSync(path string, mode uint32) error {
	return os.Chmod(path, os.FileMode(mode))
}

func JsChownSync(path string, uid int, gid int) error {
	return os.Chown(path, uid, gid)
}

func JsTruncateSync(path string, length int64) error {
	fd, err := os.OpenFile(path, os.O_WRONLY, 0222)
	if err != nil {
		return err
	}

	defer fd.Close()

	err = fd.Truncate(length)
	if err != nil {
		return err
	}

	// fd.Seek(0,0)

	err = fd.Sync()
	return err

}

func JsReadFileUtf8Sync(path string) (progpAPI.StringBuffer, error) {
	bytes, err := os.ReadFile(path)
	return bytes, err
}

func JsReadFileBytesSync(path string) ([]byte, error) {
	bytes, err := os.ReadFile(path)
	return bytes, err
}

func JsCopyFileSync(sourcePath, destPath string) error {
	sourceFileStat, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return errors.New("can't copy file")
	}

	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func JsLinkSync(existingPath, newPath string) error {
	return os.Link(existingPath, newPath)
}

func JsSymLinkSync(existingPath, newPath string) error {
	return os.Symlink(existingPath, newPath)
}

func JsUnlink(filePath string) error {
	return os.Remove(filePath)
}

//endregion

//endregion
