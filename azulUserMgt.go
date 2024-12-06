// azulUserMgt

package main

import (
	"fmt"
	"log"
	"os"

	userLib "goDemo/usermgt/userLibDb"
	cliutil "github.com/prr123/utility/utilLib"
)


func main() {

    numarg := len(os.Args)
    flags:=[]string{"dbg", "cmd", "user", "val"}

    useStr := " /cmd={list,get,add,upd,rm} /user=username </val=value> [/dbg]"
    helpStr := "user management program"

    if numarg > len(flags) +1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(-1)
    }

    if numarg == 1 || (numarg > 1 && os.Args[1] == "help") {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage is: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }

    flagMap, err := cliutil.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    dbg:= false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}


    cmdStr := ""
    cmdval, ok := flagMap["cmd"]
    if !ok {
        log.Fatalf("error -- no command flag provided!\n")
    } else {
        if cmdval.(string) == "none" {log.Fatalf("error -- no command provided!\n")}
        cmdStr = cmdval.(string)
    }

    userStr := ""
   	uval, ok := flagMap["user"]
    if !ok {
        log.Fatalf("error -- no user flag provided!\n")
    } else {
        if uval.(string) == "none" {log.Fatalf("error -- no user name provided!\n")}
        userStr = uval.(string)
    }

    valStr := ""
    val, vok := flagMap["val"]
    if vok {
        if val.(string) == "none" {log.Fatalf("error -- no yaml file provided!\n")}
        valStr = val.(string)
    }

    dbFilnam := "/home/peter/dbData/users"

	if ok:= userLib.VerifyCmd(cmdStr); !ok {log.Fatalf("error -- invalid command: %s\n", cmdStr)}
//	if ok:= apiLib.VerifyApp(appStr); !ok {log.Fatalf("error -- invalid app: %s\n", appStr)}

    if dbg {
        fmt.Printf("cmd:  %s\n", cmdStr)
        fmt.Printf("user: %s\n", userStr)
        fmt.Printf("val:  %s\n", valStr)
		fmt.Printf("db:   %s\n", dbFilnam)
//			fmt.Printf("no valFilnam\n")
    }

	api, err := userLib.InitUserDbApi(dbFilnam)
	if err != nil {log.Fatalf("error -- InitApi: %v\n", err)}
	api.Dbg = dbg
	defer api.CloseDb()

	err = api.ProcCmd(cmdStr, userStr, valStr)
	if err != nil {log.Fatalf("error -- ProcCmd: %v\n", err)}

	fmt.Println("*** user mgt success ***")
}

