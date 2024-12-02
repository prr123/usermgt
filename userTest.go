//

package main

import (
	"fmt"
	"log"
	"os"

	ul "goDemo/usermgt/userLib"
	cliutil "github.com/prr123/utility/utilLib"
)


func main() {

    numarg := len(os.Args)
    flags:=[]string{"dbg", "cmd", "user", "val"}

    useStr := " /cmd={list,get,add,upd,rm} /user=username </val={yamlfile}> [/dbg]"
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

    valStr := "users"
    val, vok := flagMap["val"]
    if vok {
        if val.(string) == "none" {log.Fatalf("error -- no yaml file provided!\n")}
        valStr = val.(string)
    }

    valFilnam := "userData/" + valStr + ".yaml"

	if ok:= ul.VerifyCmd(cmdStr); !ok {log.Fatalf("error -- invalid command: %s\n", cmdStr)}

    if dbg {
        fmt.Printf("cmd:  %s\n", cmdStr)
        fmt.Printf("user: %s\n", userStr)
		fmt.Printf("val:  %s\n", valFilnam)
//			fmt.Printf("no valFilnam\n")
    }

	um, err := ul.InitUserList(valFilnam)
	if err != nil {log.Fatalf("error -- InitUserList: %v\n", err)}

	um.Dbg = dbg
	if dbg {um.PrintList()}

	fmt.Println("*** api success ***")
}

