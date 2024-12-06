package userLibDb

import (
	"fmt"
//	"os"
//	"bytes"
//	"math/rand"
//	"time"

	db "github.com/prr123/pogLib/pogLib"
//	"github.com/goccy/go-yaml"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//type provMap map[string]string

type userDbApi struct {
//	list map[string]provMap
	pogdb *db.PogDB
	Dbg bool
	dbFilnam string
}




func VerifyCmd(cmdStr string) (bool) {
    cmdList := []string{"list", "get", "add", "upd", "rm"}
    for i:=0; i< len(cmdList); i++ {
        if cmdStr == cmdList[i] {
            return true
        }
    }
    return false
}


func InitUserDbApi(dbdir string)(api *userDbApi, err error) {

    var apiObj userDbApi
    apiObj.pogdb, err = db.InitPogDb(dbdir)
    if err != nil {return nil, fmt.Errorf("InitApi-InitPodDb: %v\n", err)}
    return &apiObj, nil
}


func (api *userDbApi) ProcCmd(cmdStr, userNam, valStr string) (error) {

    switch cmdStr {

    // list users
    case "list":
        if userNam == "*" || userNam == "all" {
            UserList, err := api.ListAllUsers()
            if err != nil {return fmt.Errorf("ListAllUsers: %v\n",err)}
            for i, userNam := range UserList  {
                fmt.Printf("--%d: %s\n", i, userNam)
            }
            return nil
        }
        ok, err:= api.ListUser(userNam)
        if err != nil {return fmt.Errorf("ListUser: %v", err)}
        if ok {
            fmt.Printf("dbg -- api: %s found!\n", userNam)
        } else {
            fmt.Printf("dbg -- api: %s not found!\n", userNam)
        }
        return nil

    case "get":
        fmt.Printf("dbg -- Cmd: get; User: %s\n", userNam)
        token, err := api.GetUserInfo(userNam)
        if err != nil {return fmt.Errorf("GetToken: %v", err)}
        fmt.Printf("dbg -- User: %s Token: %s\n", userNam, token)
        return nil

    case "add":
        fmt.Printf("dbg -- Cmd: add; User: %s\n", userNam)
        err := api.AddUser(userNam, valStr)
        if err != nil {return fmt.Errorf("AddUser: %v", err)}
        return nil

    case "rm":
        fmt.Printf("dbg -- Cmd: rm; User: %s\n", userNam)
        err := api.RmUser(userNam)
        if err != nil {return fmt.Errorf("RmUser: %v", err)}
        return nil

    case "upd":
        fmt.Printf("dbg -- Cmd: upd; User: %s\n", userNam)
        err := api.UpdUser(userNam, valStr)
        if err != nil {return fmt.Errorf("UpUser: %v", err)}
        return nil

    default:
        return fmt.Errorf("unknown command: %s\n", cmdStr)

    }

    return nil
}


func (api *userDbApi) ListUser(userNam string) (bool, error) {
    db := api.pogdb
    ok, err := db.HasKey(userNam)
    if err != nil {return false, fmt.Errorf("DbHas: %v\n", err)}
    return ok, nil
}

func (api *userDbApi) ListAllUsers() (UserList []string, err error) {

    db := api.pogdb
    Usernum, err := db.DbCount()
    if err != nil {return UserList, fmt.Errorf("DbCount: %v\n", err)}

    UserList = make([]string, Usernum)

    count:=0
    for i:=0; i<Usernum; i++  {
        User,_, end, err := db.NextItem()
        if err != nil {return UserList, fmt.Errorf("NextItem: %v\n", err)}
        if end {break}
        UserList[i] = string(User)
        count++
    }
    return UserList[:count], nil
}


func (api *userDbApi) GetUserInfo(userNam string) (string, error){

    db := api.pogdb
    token, err := db.Read(userNam)
    if err != nil {return "", fmt.Errorf("GetToken: %v", err)}

    return string(token), nil
}

func (api *userDbApi) UpdUser(userNam, valStr string) (error){

    db := api.pogdb
    err := db.Upd(userNam, []byte(valStr))
    if err != nil {return fmt.Errorf("dbUpd: %v", err)}

    return nil
}


func (api *userDbApi) AddUser(userNam, token string) (error){

    db := api.pogdb
    err := db.Add(userNam, []byte(token))
    if err != nil {return fmt.Errorf("dbAdd: %v", err)}
    return nil
}

func (api *userDbApi) RmUser(userNam string) (error){

    db := api.pogdb
    err := db.Del(userNam)
    if err != nil {return fmt.Errorf("dbDel: %v", err)}
    return nil
}

func (api *userDbApi) DbClose() (error){
    db := api.pogdb
    err := db.Close()
    if err != nil {return fmt.Errorf("dbClose: %v", err)}
    return nil
}


func (api *userDbApi) PrintUserList() (error) {

    db := api.pogdb
    Usernum, err := db.DbCount()
    if err != nil {return fmt.Errorf("DbCount: %v\n", err)}

    count:=0
    for i:=0; i<Usernum; i++  {
        User, val, end, err := db.NextItem()
        if err != nil {return fmt.Errorf("NextItem: %v\n", err)}
        if end {break}
        fmt.Printf(" --%d: %s %s\n", i, User, string(val))
        count++
    }
    fmt.Printf("*** end user list ***\n")
	return nil
}
