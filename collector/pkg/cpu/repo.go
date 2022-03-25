package cpu

type Repository interface {
	// retrieve from database
	GetCPUs() []CPU
	// create record
}
