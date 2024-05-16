package grp

import (
	"github.com/panjf2000/ants/v2"
	"k8s.io/klog/v2"
)

type GoRoutinePool struct {
	pool *ants.Pool
}

func NewGoRoutinePool(maxGoRoutine int) *GoRoutinePool {
	p, err := ants.NewPool(maxGoRoutine)
	if err != nil {
		klog.Errorf("New pool with err %v", err)
		return nil
	}

	pool := &GoRoutinePool{
		pool: p,
	}

	return pool
}

func (gpr *GoRoutinePool) Run(task func()) error {
	return gpr.pool.Submit(task)
}

func (gpr *GoRoutinePool) Close() {
	gpr.pool.Release()
}
