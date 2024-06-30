package cron

import "sync"

type manager struct {
	tasks []cronTask
}

type cronTask interface {
	// TH
	Run(Done func())
}

func NewCronManager() *manager {
	return &manager{}
}

func (m *manager) EquipmentTask(task cronTask) {
	m.tasks = append(m.tasks, task)
}

func (m *manager) Run() {
	wg := sync.WaitGroup{}
	wg.Add(len(m.tasks))
	for _, task := range m.tasks {
		go task.Run(func() {
			wg.Done()
		})
	}
	wg.Wait()
}
