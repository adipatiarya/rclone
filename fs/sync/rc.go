package sync

import (
	"github.com/ncw/rclone/fs/rc"
)

func init() {
	for _, name := range []string{"sync", "copy", "move"} {
		name := name
		moveHelp := ""
		if name == "move" {
			moveHelp = "- deleteEmptySrcDirs - delete empty src directories if set\n"
		}
		rc.Add(rc.Call{
			Path:         "sync/" + name,
			AuthRequired: true,
			Fn: func(in rc.Params) (rc.Params, error) {
				return rcSyncCopyMove(in, name)
			},
			Title: name + " a directory from source remote to destination remote",
			Help: `This takes the following parameters

- srcFs - a remote name string eg "drive:src" for the source
- dstFs - a remote name string eg "drive:dst" for the destination
` + moveHelp + `
This returns
- jobid - ID of async job to query with job/status

See the [` + name + ` command](/commands/rclone_` + name + `/) command for more information on the above.`,
		})
	}
}

// Sync/Copy/Move a file
func rcSyncCopyMove(in rc.Params, name string) (out rc.Params, err error) {
	srcFs, err := rc.GetFsNamed(in, "srcFs")
	if err != nil {
		return nil, err
	}
	dstFs, err := rc.GetFsNamed(in, "dstFs")
	if err != nil {
		return nil, err
	}
	switch name {
	case "sync":
		return nil, Sync(dstFs, srcFs)
	case "copy":
		return nil, CopyDir(dstFs, srcFs)
	case "move":
		deleteEmptySrcDirs, err := in.GetBool("deleteEmptySrcDirs")
		if rc.NotErrParamNotFound(err) {
			return nil, err
		}
		return nil, MoveDir(dstFs, srcFs, deleteEmptySrcDirs)
	}
	panic("unknown rcSyncCopyMove type")
}