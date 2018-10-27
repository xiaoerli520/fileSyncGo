package rsync

type Task struct {

	fileName string

	targetIP string

	Result string

	Err string
}

func (r *Rsync) NewTask(fileName string, targetIP string) (*Task) {
	var task = new(Task)
	task.fileName = fileName
	task.targetIP = targetIP
	return task
}

func (t *Task) GetFileName() string {
	return t.fileName
}

func (t *Task) GetTarget() string {
	return t.targetIP
}


