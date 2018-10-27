package boot

import (
	"sinago/rsync"
	"strings"
	"time"
)

var (
	maxInit int
	cmsPath []string
	frontPath string
)

var (
	rsyncBin string
	baseDir string
	pathList []string
	moduleName string
	params []string
)

var (
	dustbinSecond time.Duration
)

func InitRsync() (*rsync.Rsync) {

	var r = new(rsync.Rsync)

	maxInit = Viper.GetInt("Rsync.maxInit")
	cmsPath = Viper.GetStringSlice("HealthKeeper.cmsPath")
	frontPath = Viper.GetString("HealthKeeper.IPSPath")

	rsyncBin = Viper.GetString("Rsync.rsyncBin")
	baseDir = Viper.GetString("Rsync.baseDir")
	pathList = Viper.GetStringSlice("Rsync.pathList")
	moduleName = Viper.GetString("Rsync.moduleName")
	params = Viper.GetStringSlice("Rsync.rsyncParams")

	dustbinSecond = Viper.GetDuration("Rsync.dustbinSecond")

	rsyncParams := Viper.GetStringSlice("Rsync.rsyncParams")
	r.SetWorkers(Viper.GetInt("Rsync.maxProcesses"), rsyncParams)

	r.Prepare(moduleName,baseDir,rsyncBin,rsyncParams...)

	r.OnFailed(func(task *rsync.Task) {
		// 过滤相关WARNING
		switch {
		case strings.Contains(task.Err, "Connection reset by peer"):
			r.AddRecycle(task)
		case strings.Contains(task.Err, "update discarded (will try again)"):
			r.AddRecycle(task)
		case strings.Contains(task.Err, "Connection refused"):
			// 不进行add 等待recovery即可
		case strings.Contains(task.Err, "io timeout after 3 seconds"):
			r.AddRecycle(task)
		default :
			MailPool.Add(time.Now().Format("[2006-01-02 15:04:05]")+
				"Rsync: " + "fileSyncGo_Rsync" +
				task.GetTarget() + " : " + task.Err+
				" : "+task.GetFileName())
			r.AddRecycle(task)
		}
		Logger.Warn(
			time.Now().Format("[2006-01-02 15:04:05]")+
				"Rsync: " + "fileSyncGo_Rsync" +
				task.GetTarget() + " : " + task.Err+
				" : "+task.GetFileName())
	})

	r.OnRecycle(func() {
		for {
			time.Sleep(dustbinSecond * time.Second)
			if !IsInit {
				continue
			}
			for _, v := range r.RecycleTasks {
				readyList := Hk.List(Viper.GetInt("HealthKeeper.status.ready"))
				if _,ok := readyList[v.GetTarget()]; ok {
					r.Collect(v)
					Logger.Notice("Rsync: " + "recycleBin recover" + v.GetTarget() + " : " + v.GetFileName())
				}
			}
			r.DiscardRecycle()
		}
	})
	return r
}
