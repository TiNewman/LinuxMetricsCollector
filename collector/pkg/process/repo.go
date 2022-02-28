package process

type Repository interface {
	// select from database
	GetProcesses() []Process
	// create record
	RecordProcess(Process)
}
