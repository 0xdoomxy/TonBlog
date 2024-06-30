package cron

type manager struct {
	tasks []cronTask
}

type cronTask interface {
	// TH
	Run()
}

func NewCronManager() *manager {
	return &manager{}
}

func (m *manager) EquipmentTask(task ...cronTask) {
	m.tasks = append(m.tasks, task...)
}

func (m *manager) Run() {
	for _, task := range m.tasks {
		go task.Run()
	}
}
