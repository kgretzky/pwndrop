package daemon

import "strings"

// Daemon interface has a standard set of methods/commands
type Daemon interface {

	// Install the service into the system
	Install(exec_path string, args ...string) (string, error)

	// Remove the service and all corresponding files from the system
	Remove() (string, error)

	// Start the service
	Start() (string, error)

	// Stop the service
	Stop() (string, error)

	// Status - check the service status
	Status() (string, error)

	// Run - run executable service
	Run(e Executable) (string, error)
}

// Executable interface defines controlling methods of executable service
type Executable interface {
	// Start - non-blocking start service
	Start()
	// Stop - non-blocking stop service
	Stop()
	// Run - blocking run service
	Run()
}

// New - Create a new daemon
//
// name: name of the service
//
// description: any explanation, what is the service, its purpose
func New(name, description string, dependencies ...string) (Daemon, error) {
	return newDaemon(strings.Join(strings.Fields(name), "_"), description, dependencies)
}
