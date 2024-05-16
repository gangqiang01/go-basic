package v1

// running: 0, error: 1, success: 2
const (
	//log status
	StatusRunning int32 = 0
	StatusError   int32 = 1
	StatusSuccess int32 = 2

	//run status
	StatusEnable  string = "enable"
	StatusDisable string = "disable"
)
