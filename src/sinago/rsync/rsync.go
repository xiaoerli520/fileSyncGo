package rsync

import (
	"strings"
	"fmt"
	"bytes"
	"sync"
	"sinago/cmd"
	"os/exec"
)

type FailFunc func (task *Task)
type FinishFunc func (task *Task)
type RecycleFunc func()
type Tasks []*Task

type Rsync struct {

	params       []string
	moduleName   string
	baseDir      string
	bin          string

	RecycleTasks Tasks

	failFunc     FailFunc
	finishFunc   FinishFunc
	recycler     RecycleFunc

	Workers      Worker

	Lock         sync.RWMutex


}

func (rsync *Rsync) Prepare(
	moduleName string,
	baseDir string, bin string, params ...string) {

	rsync.baseDir = baseDir
	rsync.params = params
	rsync.moduleName = moduleName
	rsync.bin = bin
}

func (rsync *Rsync) Monitor(task *Task) {
	fmt.Println(rsync.params, task.fileName, task.targetIP+"::"+rsync.moduleName)
}

func (rsync *Rsync) Sync(task *Task) (bool, error) {
	var buffer bytes.Buffer
	buffer.WriteString(task.targetIP)
	buffer.WriteString("::")
	buffer.WriteString(rsync.moduleName)
	rsyncParam := buildParams(rsync.params, task.fileName, buffer.String())
	out, err := cmd.ExecShellAt(rsync.baseDir, rsync.bin, rsyncParam)
	task.Result = out
	task.Err = err
	if len(err) > 0 {
		if rsync.failFunc != nil {
			go func(task *Task) {
				rsync.failFunc(task)
			}(task)
		}
	}
	return true, nil
}

func FullSync(bin string, baseDir string, srcs []string, dstIP string, module string, params []string) (bool, string) {
	var stdout string
	var stderr string
	var target bytes.Buffer
	var sources []string
	target.WriteString(dstIP)
	target.WriteString("::")
	target.WriteString(module)
	// 合并 src  减少 rsync 数量
	fullParam := []string{"-auz", "--omit-dir-times", "--ignore-errors","--exclude", "*.git", "--exclude", "*.svn", "--exclude", "*.tmp", "--exclude", ".~tmp~/"}
	for _, value := range srcs {
		sources = append(sources, baseDir+value)
	}
	fullParam = append(fullParam, sources...)
	fullParam = append(fullParam, target.String())
	stdout, stderr = cmd.ExecShell(bin, fullParam...)
	fmt.Println(bin, fullParam)
	if len(stderr) > 0 {
		return false, stderr
	}
	return true, stdout
}

func IsFullSyncing (target string) (isSync bool) {
	findCommand := "ps aux | grep 'rsync -auz' | grep '"+target+"'"
	out,_ := exec.Command("bash","-c", findCommand).Output()
	outString := string(out)
	for _, v := range strings.Split(outString, "\n") {
		if strings.Contains(v, "grep") {
			continue
		}
		if len(v) < 5 {
			continue
		}
		if strings.Contains(v, "rsync -auz") {
			isSync = true
			break
		}
	}
	return isSync
}

// 判断当前是否正在同步，防止端口占用导致Dail失败
func IsRsyncing (target string) (isSync bool) {
	findCommand := "ps aux | grep 'rsync -avuzR' | grep '"+target+"'"
	out,_ := exec.Command("bash","-c", findCommand).Output()
	outString := string(out)
	for _, v := range strings.Split(outString, "\n") {
		if strings.Contains(v, "grep") {
			continue
		}
		if len(v) < 5 {
			continue
		}
		if strings.Contains(v, target) {
			isSync = true
			break
		}
	}
	return
}

func ShowDiff(bin string, baseDir string, srcs []string, dstIP string, module string) (bool, string) {
	var stdout string
	var stderr string
	for k := range srcs {
		srcs[k] = strings.TrimLeft(srcs[k], "/")
	}
	diffParam := []string{"-auR", "--out-format", "%n", "--omit-dir-times", "--exclude", "*.git", "--exclude", "*.svn", "--exclude", "*.tmp", "--exclude", ".~tmp~/", "--dry-run"}
	diffParam = append(diffParam, srcs...)
	diffParam = append(diffParam, dstIP+"::"+module)

	stdout, stderr = cmd.ExecShellAt(baseDir ,bin, diffParam)
	if len(stderr) > 0 {
		return false, stderr
	}
	return true,stdout
}

func DirectSync(bin string, baseDir string, params []string,fileName string, dstIP string, module string) (bool, string) {
	rsyncParam := params
	rsyncParam = append(rsyncParam,fileName, dstIP+"::"+module)
	result, stderr := cmd.ExecShellAt(baseDir, bin, rsyncParam)
	if len(stderr) > 0 {
		return false, stderr
	}
	return true, result
}

func buildParams(params []string, src string, dst string) []string {
	s := make([]string, 0)
	s = append(s, params...)
	s = append(s, src)
	s = append(s, dst)
	return s
}

func (r *Rsync) AddRecycle(task *Task) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	r.RecycleTasks = append(r.RecycleTasks, task)
}

func (r *Rsync) DiscardRecycle() {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	r.RecycleTasks = []*Task{}
}

func (r *Rsync) OnFailed(onFailed FailFunc) {
	r.failFunc = onFailed
}

func (r *Rsync) OnFinish(onFinish FinishFunc) {
	r.finishFunc = onFinish
}

// 回收机制
func (r *Rsync) OnRecycle(onRecycle RecycleFunc) {
	r.recycler = onRecycle
}



