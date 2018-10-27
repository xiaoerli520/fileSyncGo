package healthKeeper

import (
	"fmt"
	"sync"
)

type Checker func(this *HealthKeeper)

type Symptom struct {

	Code       int

	Phenomenon string

}

type StatusEvent struct {

	Member string

	Symptom *Symptom

}

type StatusFunc struct {

	Func func(hk *HealthKeeper, se *StatusEvent)

	Symptom *Symptom

}


type HealthKeeper struct {

	IKeeper

	Healthy map[string]*Symptom

	Lock sync.RWMutex

	checkers []Checker

	onStartFunc func(args ...interface{})

	OnStatusFuncs []StatusFunc
}


func (k *HealthKeeper) SetChecker (newChecker Checker) {
	k.Lock.Lock()
	defer k.Lock.Unlock()
	k.checkers = append(k.checkers, newChecker)

}

func (k *HealthKeeper) OnStart (onStart func(args ...interface{})) {
	k.onStartFunc = onStart
}

func (k *HealthKeeper) Start () {

	k.Healthy = make(map[string]*Symptom)

	if nil != k.onStartFunc {
		k.onStartFunc()
	}
	for _, checker := range k.checkers {
		go func(checker Checker) {
			checker(k)
		}(checker)
	}
}

func (k *HealthKeeper) List (Code int) map[string]*Symptom {
	k.Lock.RLock()
	defer k.Lock.RUnlock()
	var returnMap = make(map[string]*Symptom)
	for k, v := range k.Healthy {
		if v.Code == Code {
			returnMap[k] = v
		}
	}
	return returnMap
}

func (k *HealthKeeper) OnStatus (onStatus func(hk *HealthKeeper, se *StatusEvent), code int) {
	k.OnStatusFuncs = append(k.OnStatusFuncs, StatusFunc{onStatus, &Symptom{code, ""}})
}

func (k *HealthKeeper) SetMember (member string, symptom *Symptom) {
	k.Lock.Lock()
	defer k.Lock.Unlock()
	if nil != k.OnStatusFuncs {
		for _, OnStatus := range k.OnStatusFuncs {
			if OnStatus.Symptom.Code == symptom.Code {
				OnStatus.Func(k, &StatusEvent{Member:member, Symptom:symptom})
			}
		}
	}
	k.Healthy[member] = symptom
}

func (k *HealthKeeper) DelMember(key string) {
	k.Lock.Lock()
	defer k.Lock.Unlock()
	delete(k.Healthy, key)
}

func (k *HealthKeeper) Monitor () {
	for key, v := range k.Healthy {
		fmt.Println(key, " STATUS IS ",v.Code)
	}
}


