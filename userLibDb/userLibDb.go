package userLibDb

import (
	"fmt"
	"os"
	"strings"
	db "github.com/prr123/pogLib/pogLib"
	"github.com/goccy/go-json"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//type provMap map[string]string

type UserData struct {
	Short string `json:"short,omitempty"`
	Role string `json:"role,omitempty"`
	Token string `json:"token,omitempty"`
}

type userDbApi struct {
//	list map[string]provMap
	pogdb *db.PogDB
	Dbg bool
	dbFilnam string
}


func EncodeToJson(dat UserData)(jsonData []byte, err error) {

	jsonData, err =json.Marshal(dat)
	return jsonData, err
}

func DecodefromJson(jsonData []byte)(dat UserData, err error) {
	err = json.Unmarshal(jsonData, &dat)
	return dat, err
}



func VerifyCmd(cmdStr string) (bool) {
    cmdList := []string{"list", "get", "getUserInfo", "add", "addUser", "upd", "updUser", "rm", "clear","printDbraw" }
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

	dbg := api.Dbg
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
        if dbg {fmt.Printf("dbg -- Cmd: get; User: %s\n", userNam)}
        valStr, err := api.GetUserStr(userNam)
        if err != nil {return fmt.Errorf("GetUserStr: %v", err)}
        if dbg {fmt.Printf("dbg -- User: %s val: %s\n", userNam, valStr)}
        return nil

    case "getUserInfo":
        if dbg {fmt.Printf("dbg -- Cmd: getUser; User: %s\n", userNam)}
        userInfo, err := api.GetUserInfo(userNam)
        if err != nil {return fmt.Errorf("GetUserInfo: %v", err)}
        if dbg {fmt.Printf("dbg -- User: %s Info: %v\n", userNam,userInfo )}
        return nil

    case "add":
        if dbg {fmt.Printf("dbg -- Cmd: add; User: %s\n", userNam)}
        err := api.AddUserStr(userNam, valStr)
        if err != nil {return fmt.Errorf("AddUserStr: %v", err)}
        return nil

    case "addUser":
        if dbg {fmt.Printf("dbg -- Cmd: addUser; User: %s, val: %s\n", userNam, valStr)}
		userInfo, err := api.Decode(valStr)
        if err != nil {return fmt.Errorf("AddUser Decode: %v", err)}
        err = api.AddUserInfo(userNam, *userInfo)
        if err != nil {return fmt.Errorf("AddUser: %v", err)}
        return nil

    case "rm":
        if dbg {fmt.Printf("dbg -- Cmd: rm; User: %s\n", userNam)}
        err := api.RmUser(userNam)
        if err != nil {return fmt.Errorf("RmUser: %v", err)}
        return nil

    case "upd":
        if dbg {fmt.Printf("dbg -- Cmd: upd; User: %s\n", userNam)}
        err := api.UpdUserStr(userNam, valStr)
        if err != nil {return fmt.Errorf("UpdUserStr: %v", err)}
        return nil

    case "updUser":
        if dbg {fmt.Printf("dbg -- Cmd: updUser; User: %s\n", userNam)}
		userInfo, err := api.Decode(valStr)
		if err != nil {return fmt.Errorf("UpdUser-Decode: %v", err)}
        err = api.UpdUserInfo(userNam, userInfo)
        if err != nil {return fmt.Errorf("UpdUserStr: %v", err)}
        return nil

    case "clear":
        if dbg {fmt.Printf("dbg -- Cmd: clear\n")}
        err := api.ClearDb()
        if err != nil {return fmt.Errorf("ClearDb: %v", err)}

	case "printDbraw":
        if dbg {fmt.Printf("dbg -- Cmd: printDbraw\n")}
        err := api.PrintDb(-1, true)
        if err != nil {return fmt.Errorf("PrintDbraw: %v", err)}

	case "printDb":
        if dbg {fmt.Printf("dbg -- Cmd: printDb\n")}
        err := api.PrintDb(-1, false)
        if err != nil {return fmt.Errorf("PrintDb: %v", err)}

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


func (api *userDbApi) GetUserStr(userNam string) (string, error){

    db := api.pogdb
    token, err := db.Read(userNam)
    if err != nil {return "", fmt.Errorf("GetToken: %v", err)}
    return string(token), nil
}

func (api *userDbApi) GetUserData(userNam string) ([]byte, error){
    db := api.pogdb
    token, err := db.Read(userNam)
    if err != nil {return nil, fmt.Errorf("GetToken: %v", err)}
    return token, nil
}

func (api *userDbApi) GetUserInfo(userNam string) (uInfo *UserData,err error){

	var userInfo UserData
    db := api.pogdb
    token, err := db.Read(userNam)
    if err != nil {return nil, fmt.Errorf("GetToken: %v", err)}
	err = json.Unmarshal(token, &userInfo)
	if err != nil {return nil, fmt.Errorf("Unmarshal: %v", err)}

    return &userInfo, nil
}

func (api *userDbApi) UpdUserStr(userNam, valStr string) (error){

    db := api.pogdb
    err := db.Upd(userNam, []byte(valStr))
    if err != nil {return fmt.Errorf("dbUpd: %v", err)}
    return nil
}

func (api *userDbApi) UpdUser(userNam string, valdata []byte) (error){

    db := api.pogdb
    err := db.Upd(userNam, valdata)
    if err != nil {return fmt.Errorf("dbUpd: %v", err)}
    return nil
}

func (api *userDbApi) UpdUserInfo(userNam string, userInfo *UserData) (error){

    db := api.pogdb
	valdata, err := json.Marshal(*userInfo)
	if err != nil {return fmt.Errorf("Marshal: %v", err)}
    err = db.Upd(userNam, valdata)
    if err != nil {return fmt.Errorf("dbUpd: %v", err)}
    return nil
}

func (api *userDbApi) AddUserStr(userNam, token string) (error){

    db := api.pogdb
    err := db.Add(userNam, []byte(token))
    if err != nil {return fmt.Errorf("dbAdd: %v", err)}
    return nil
}

func (api *userDbApi) AddUser(userNam string, token []byte) (error){

    db := api.pogdb
    err := db.Add(userNam, token)
    if err != nil {return fmt.Errorf("dbAdd: %v", err)}
    return nil
}

func (api *userDbApi) AddUserInfo(userNam string, UserInfo UserData) (error){

//fmt.Printf("dbg -- userInfo: %v, short: %s\n", UserInfo, UserInfo.Short)
    db := api.pogdb
	token, err := json.Marshal(UserInfo)
	if err != nil {return fmt.Errorf("Marshal: %v", err)}
//fmt.Printf("dbg -- token: %s\n", token)
    err = db.Add(userNam, token)
    if err != nil {return fmt.Errorf("dbAdd: %v", err)}
    return nil
}

func (api *userDbApi) RmUser(userNam string) (error){

    db := api.pogdb
    err := db.Del(userNam)
    if err != nil {return fmt.Errorf("dbDel: %v", err)}
    return nil
}

func (api *userDbApi) CloseDb() (error){
    db := api.pogdb
    err := db.Close()
    if err != nil {return fmt.Errorf("dbClose: %v", err)}
    return nil
}

//        err := api.ClearDb()
func (api *userDbApi) ClearDb() (error){

	err := os.RemoveAll(api.dbFilnam)
	return err
}

func (api *userDbApi) Decode(inp string) (userInfo *UserData, err error){

	res := strings.Split(inp, ",")
	if len(res) != 3 {return nil, fmt.Errorf("insufficient strings!")}
	userInfo = &UserData{res[0], res[1], res[2]}
	return userInfo, nil
}

func (api *userDbApi) PrintDb(limit int, raw bool) (error) {

    db := api.pogdb
    Usernum, err := db.DbCount()
    if err != nil {return fmt.Errorf("DbCount: %v\n", err)}

	if limit > -1 && limit < Usernum {Usernum = limit}

    count:=0
    for i:=0; i<Usernum; i++  {
        User, val, end, err := db.NextItem()
        if err != nil {return fmt.Errorf("NextItem: %v\n", err)}
        if end {break}
        if raw {
			fmt.Printf(" --%d: %s %s\n", i, User, string(val))
		} else {
			userInfo:= &UserData{}
			err =json.Unmarshal(val, &userInfo)
			if err != nil {return fmt.Errorf("--%d %s PrintDb Unmarshal: %v", i, User, err)}
			fmt.Printf(" --%d: %s\n", i, User)
			fmt.Printf("         short: %s\n", userInfo.Short)
			fmt.Printf("         role:  %s\n", userInfo.Role)
			fmt.Printf("         token: %s\n", userInfo.Token)
		}
        count++
    }
    fmt.Printf("*** end user list ***\n")
	return nil
}

func (api *userDbApi) PrintUserInfo(userInfo *UserData) (error) {

	fmt.Println("***** UserInfo *****")
	fmt.Printf("  short: %s\n", userInfo.Short)
	fmt.Printf("  role:  %s\n", userInfo.Role)
	fmt.Printf("  token: %s\n", userInfo.Token)
	fmt.Println("*** End UserInfo ***")
	return nil
}
