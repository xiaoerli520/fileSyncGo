package boot

import (
	"sinago/healthKeeper"
	"fmt"
	"strings"
	"filesync/notification"
	"sinago/cmd"
	"time"
	"sinago/utils"
	"sinago/rsync"
)

// if is first full synced
var IsInit bool

// front IPS
var targetIPS []string

// cms IPS
var Cms    []string

// Global HealhtKeeper (single)
var Hk     *healthKeeper.HealthKeeper

func InitKeeper() *healthKeeper.HealthKeeper {

	// Semaphore controls Concurrency Num of goroutine
	var sem = utils.NewSemaphore(maxInit)

	// Init Status Code
	var (
		hk = new(healthKeeper.HealthKeeper)
		ready = Viper.GetInt("HealthKeeper.status.ready")
		dead = Viper.GetInt("HealthKeeper.status.dead")
		recovery = Viper.GetInt("HealthKeeper.status.recovery")
		unhealth = Viper.GetInt("HealthKeeper.status.unhealth")
	)

	// Init Symptom
	var (
		statusReady = &healthKeeper.Symptom{ready, "节点正常"}
		statusDead = &healthKeeper.Symptom{dead, "节点不可达"}
		statusUnhealthy = &healthKeeper.Symptom{unhealth, "节点不健康"}
		statusRecovery = &healthKeeper.Symptom{recovery, "节点恢复中"}
	)

	// Init The dormancy time of different states
	var (
		totalSleep = Viper.GetDuration("HealthKeeper.sleepTimes.total")*time.Second
		readySleep = Viper.GetDuration("HealthKeeper.sleepTimes.ready")*time.Second
		deadSleep = Viper.GetDuration("HealthKeeper.sleepTimes.dead")*time.Second
		unhealthSleep = Viper.GetDuration("HealthKeeper.sleepTimes.unhealth")*time.Second
		recoverySleep = Viper.GetDuration("HealthKeeper.sleepTimes.recovery")*time.Second
		unhealthTryTime = Viper.GetInt("HealthKeeper.sleepTimes.unhealthTryTime")
		unhealthFatalTime = Viper.GetInt("HealthKeeper.sleepTimes.unhealthFatalTime")
		rsyncTimeout = Viper.GetDuration("Rsync.timeout")*time.Second
	)

	// Set OnStart Func
	hk.OnStart(func(args ...interface{}) {
		fmt.Println("Current Environment is ",Env)
		if strings.Contains(Env, "dev") {
			devIPS := Viper.GetStringSlice("HealthKeeper.devIPS")
			fmt.Println(devIPS)
			Logger.Info("[Dev] Inited IPS : " + strings.Join(devIPS, ","))
			for _, v := range devIPS {
				hk.SetMember(v, &healthKeeper.Symptom{ready, ""})
			}
		} else {
			result, stderr := cmd.ExecShell("php", frontPath)
			if len(stderr) > 0 {
				err, res := notification.Email(EmailObject, "fileSyncGo_IP列表获取失败", EmailList, stderr)
				Logger.Err("fileSyncGo_IP列表获取失败" + stderr)
				if !err {
					Logger.Err("邮件发送失败" + res)
				}
				panic("fileSyncGo_IP列表获取失败" + stderr)
				return
			}
			Logger.Info("Inited targetIPS IPS : " + result)
			targetIPS = strings.Split(result, ",")
			Logger.Info("Inited All IPS : " + strings.Join(targetIPS, ","))
			for _, v := range targetIPS {
				Logger.Info("Logged IP : " + v)
				hk.SetMember(v, &healthKeeper.Symptom{recovery, ""})
			}
		}
	})

	hk.SetChecker(func(hk *healthKeeper.HealthKeeper) {
		for {
			time.Sleep(readySleep)
			aliveList := hk.List(ready)
			for k := range aliveList {
				time.Sleep(3 * time.Second)
				if rsync.IsRsyncing(k) {
					continue
				}
				if !utils.IsReachable("tcp", k, 873, rsyncTimeout) {
					Logger.Warn("IP : " + k + " Being Unhealthy")
					hk.SetMember(k, statusUnhealthy)
				}
			}
		}
	})

	hk.SetChecker(func(hk *healthKeeper.HealthKeeper) {
		for {
			time.Sleep(deadSleep)
			deadList := hk.List(dead)
			for k := range deadList {
				if utils.IsReachable("tcp", k, 873, rsyncTimeout) {
					Logger.Warn("IP : " + k + " Being Unhealthy")
					hk.SetMember(k, statusUnhealthy)
				}
			}
		}
	})

	hk.SetChecker(func(hk *healthKeeper.HealthKeeper) {
		for {
			time.Sleep(unhealthSleep)
			unhealthList := hk.List(unhealth)

			for k := range unhealthList {
				var failTime = 0
				for i := 1; i <= unhealthTryTime; i++ {
					time.Sleep(unhealthSleep)
					if !utils.IsReachable("tcp", k, 873, rsyncTimeout) {
						failTime++
					}
				}
				if failTime > unhealthFatalTime {
					Logger.Err("IP : " + k + " Dead")
					Logger.Warn("IP : " + k + " Dead")
					hk.SetMember(k, statusDead)
				} else {
					Logger.Warn("IP : " + k + " Recovery")
					Logger.Notice("IP : " + k + " Recovery")
					hk.SetMember(k, statusRecovery)
				}
			}
		}
	})

	hk.SetChecker(func(hk *healthKeeper.HealthKeeper) {
		for {
			time.Sleep(totalSleep)
			if strings.Contains(Env, "prod") {
				IPS, _ := cmd.ExecShell("php", frontPath)
				IPSArr := strings.Split(IPS, ",")
				IPSArr = utils.RemoveDupAndEmpty(IPSArr)
				// 减机器
				for k := range hk.Healthy {
					posFront := utils.ArrayFind(k, IPSArr)
					posCms := utils.ArrayFind(k, cmsPath)
					if posFront == -1 && posCms == -1 { // 都 不存在 就delete
						Logger.Notice("IP : " + k + " has been Deleted")
						hk.DelMember(k)
					}
				}
				for _, v := range IPSArr {
					_, ok := hk.Healthy[v]
					if !ok {
						// 原来不存在 需要先 RECOVERY
						Logger.Notice("IP : " + v + " Added and need Recovery")
						hk.SetMember(v, statusRecovery)
					}
				}
			}
		}
	})

	hk.SetChecker(func(this *healthKeeper.HealthKeeper) {
		for {
			time.Sleep(recoverySleep)
			recoveryList := hk.List(recovery)
			sem.Discard()
			if len(recoveryList) <= 0 {
				continue
			}
			for k := range recoveryList {
				if rsync.IsFullSyncing(k) { // 如果当前有rsync进程 则取消再次同步的操作
					continue
				}
				sem.P()
				go func(k string) {
					result, stderr := rsync.FullSync(rsyncBin, baseDir, pathList, k, moduleName, params)
					if result {
						Logger.Notice("IP : " + k + " Recovery Success")
						fmt.Println(k, "全量同步成功")
						hk.SetMember(k, statusReady)
					} else {
						Logger.Err("IP : " + k + " Recovery Fail")
						fmt.Println(k, "全量同步失败")
						if IsInit {
							hk.SetMember(k, statusUnhealthy)
						}
						MailPool.Add(k+" : "+stderr)
					}
					defer sem.V()
				}(k)
			}
			sem.Wait()
			if !IsInit {
				IsInit = true
			}
		}
	})


	hk.OnStatus(func(hk *healthKeeper.HealthKeeper, se *healthKeeper.StatusEvent) {
		err, res := notification.Email(EmailObject, "fileSyncGo_IP_DEAD", EmailList, se.Member+se.Symptom.Phenomenon+" : 节点不可用")
		Logger.Err("IP : " + se.Member + " 节点不可用")
		if !err {
			Logger.Err("邮件发送失败" + res)
		}
	}, dead)

	hk.OnStatus(func(hk *healthKeeper.HealthKeeper, se *healthKeeper.StatusEvent) {
		Logger.Warn("IP : " + se.Member + " 节点不健康")
		MailPool.Add(se.Member+se.Symptom.Phenomenon+" : 节点不健康")
	}, unhealth)

	hk.OnStatus(func(hk *healthKeeper.HealthKeeper, se *healthKeeper.StatusEvent) {
		if IsInit {
			Logger.Warn("IP : " + se.Member + " 节点恢复中")
			MailPool.Add(se.Member+se.Symptom.Phenomenon+" : 节点恢复中")
		}
	}, recovery)

	Hk = hk

	return hk
}


