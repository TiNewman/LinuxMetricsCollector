package process

type Repository interface {
	// retrieve from database
	GetProcesses() []Process
	// create record
	PutNewProcess(Process) (int64, error)
}
