# fileSyncGo

Golang编写的文件实时同步工具。使用Rsync配合inotifywait进行实时同步，以及多目标机自动增活减活，自动全量同步补齐。

# 特点

- dev、prod，多环境配置

- HealthKeeper：控制多同步目标增活减活，自动对恢复过来的机器全量同步，补齐差漏文件。

- Semaphore：信号量控制并发，合理资源占用。

- MacOSX(FsWatch)、Linux(inotifywait)双平台支持。目前不支持Windows。

- SmartPool：参考Redis的RDB持久化机制，在x秒内发生x次错误，可上报错误（Email、Sms）。

- dustbin：自动收集同步失败的记录，并且按配置时间间隔再次同步。

- 多目标状态自定义hook。

- 借助golang的并发模型，资源占用更低，速度更快，并发同步量更大。

- 完善的Logger机制

- 基于YML的配置文件，使用Viper进行配置读取、管理

# 实践

新浪-模板发布系统，实测修改1个文件，1秒可同步50余台目标机，并且灵活增减活，免去手工运维。

# 配置

配置文件位于./config/config.YML 

```
AppName: fileSyncGo
LogPath: /data1/ms/log/fileSyncGo/
envPath: /data1/ms/fileSyncGo/.env
# 邮件策略，在x秒发送x个错误触发报警hook。例如这个就是90s内发生30个错误则触发hook。
# mailStrategies: { 90: 30, 120: 20, 150: 10, 300: 5 }
# 可以通过email、SMS等多种形式，通知使用方错误信息
# NoticeEmail: ["guoqingzhe@hotmail.com"]

Rsync:
  rsyncParams: ["-avuzR", "-slt","--force", "--from0" ,  "--omit-dir-times","--timeout", "3", "--ignore-errors","--exclude", "*.git", "--exclude", "*.svn", "--exclude", "*.tmp", "--exclude", ".~tmp~/"]
  # 想要监控的目录列表
  pathList: [/template, /template.test, /resource]
  moduleName: comos
  baseDir: /data1/ms
  rsyncBin: /usr/bin/rsync
  port: 873
  timeout: 10
  # 最大协程数量
  maxProcesses: 160
  # 在最初全量同步时的rsync进程最大数量
  maxInit: 12
  # 在机器从dead->ready，全量同步时最大的rsync进程最大数量
  maxRecovery: 6
  # dustbin自动再次同步失败记录的时间间隔
  dustbinSecond: 5

# HealthKeeper的状态配置项
HealthKeeper:
	# 状态以及其代码
  status:
    ready: 0
    dead: 1
    recovery: 2
    unhealth: 3
  # prod模式下，定时动态获取targetsIP的脚本（目前使用PHP脚本）
  IPSPath: /data1/ms/fileSyncGo/loadIPs.php
  # dev模式下的测试IP
  devIPS: [7.7.7.7]
  # 各种状态下的Check可达性间隔
  sleepTimes:
    total: 4
    ready: 5
    dead: 3
    unhealth: 3
    recovery: 3
    unhealthTryTime: 5
    unhealthFatalTime: 3

```

# 相关说明

## target状态

- ready：目标机状态正常且文件都是最新的，可以实时同步

- dead：目标机确认不可达，不会实时同步

- unhealth：目标机状态不能确认（对目标机的几次Ping有失败），待进一步确认，不会实时同步

- recovery：unhealth的目标机后被确认是正常的，需要进行一步全量同步，补齐差的文件，然后转入ready状态进行实时同步。该状态不会实时同步。

## target状态维护

四个状态会分别启动四个协程，不断循环去维护当前状态的targets。

### readyChecker

对ready的目标机进行固定间隔的Ping，如果有问题，则转入Unhealthy进行进一步确定。

### deadChecker

对dead的目标机进行固定间隔的Ping，如果Ping通了目标机，则转入Unhealthy进行进一步确定。

### unhealthChecker

对Unhealth的目标机进行固定间隔的多次Ping，并且根据Ping的成功与总次数的比率，确定机器复活进入recovery状态还是dead。

### recoveryChecker

recovery中的目标机会自动进行全量同步，补全差漏的文件，然后转入ready进行实时同步，如果全量同步失败，则转入Unhealthy进行进一步确定。

# 现状

目前项目结构还尚且比较混乱，没有运用一些优秀的结构，且运用方式比较固定。





