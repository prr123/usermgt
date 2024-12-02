package userLib

import (
	"fmt"
	"os"
	"bytes"
	"math/rand"
	"time"
	"github.com/goccy/go-yaml"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type provMap map[string]string

type usermgt struct {
	list map[string]provMap
	Dbg bool
	userFilnam string
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

func InitUserList(yamlFilnam string) (um *usermgt, err error) {

	var users usermgt
	list := make(map[string]provMap)

	ldat, err := os.ReadFile(yamlFilnam)
	if err != nil {return nil, fmt.Errorf("Read List: %v", err)}

	err = yaml.Unmarshal(ldat,&list)
	if err != nil {return nil, fmt.Errorf("UnMarshal: %v", err)}

	users.list = list
	users.userFilnam = yamlFilnam
	return &users, nil
}

func (ul *usermgt) SaveUserFile() (error) {

	filNam := []byte(ul.userFilnam)
	idx := bytes.Index(filNam, []byte(".yaml"))
	if idx < 1 {return fmt.Errorf("save: no yaml extension!")}
	oldFilnam := append(filNam[:idx], []byte("_old.yaml")...)
	err := os.Rename(ul.userFilnam, string(oldFilnam))
	if err != nil {return fmt.Errorf("rename: %v", err)}
	list:=ul.list
	out, err := yaml.Marshal(&list)
	if err != nil {return fmt.Errorf("Marshal: %v", err)}

	out = append([]byte("#user name list (revised)\n---\n"),out...)

	newFilnam := append(filNam[:idx], []byte("_new.yaml")...)
	err = os.WriteFile(string(newFilnam), out, 0666)
	if err != nil {return fmt.Errorf("write: %v", err)}

	return nil
}

func (ul *usermgt) ProcCmd(cmdStr, unam string) (error) {

//	list := ul.list
	switch cmdStr {
	// list users
	case "list":
		if unam == "*" || unam == "all" {
			users := ul.ListAllUsers()
//			if err != nil {return fmt.Errorf("list -- GetAllUsers: %v\n",err)}
			for i, unam := range users {
				fmt.Printf("--%d: %s\n", i, unam)
			}
			return nil
		}
		ok:= ul.ListUser(unam)
		if ok {
			fmt.Printf("dbg -- user: %s found!\n", unam)
		} else {
			fmt.Printf("dbg -- user: %s not found!\n", unam)
		}
		return nil

	case "get":
		fmt.Printf("dbg -- Cmd: get; User: %s\n", unam)
		token, err := ul.GetToken(unam)
		if err != nil {return fmt.Errorf("GetToken: %v", err)}
		fmt.Printf("dbg -- user: %s Token: %s\n", unam, token)
		return nil

	case "add":
		fmt.Printf("dbg -- Cmd: add; User: %s\n", unam)
		err := ul.AddUser(unam)
		if err != nil {return fmt.Errorf("AddUser: %v", err)}
		return nil

	case "rm":
		fmt.Printf("dbg -- Cmd: rm; User: %s\n", unam)
		err := ul.RmUser(unam)
		if err != nil {return fmt.Errorf("RmUser: %v", err)}
		return nil

	case "upd":
		fmt.Printf("dbg -- Cmd: upd; User: %s\n", unam)

		err := ul.UpdUser(unam)
		if err != nil {return fmt.Errorf("UpdUser: %v", err)}
		return nil

	default:
		return fmt.Errorf("unknown command: %s\n", cmdStr)

	}

	return nil
}

func (ul *usermgt) ListUser(unam string) (bool) {
	list := ul.list
	_, ok := list[unam]
	return ok
}

func (ul *usermgt) ListAllUsers() (users []string) {

	users = make([]string, len(ul.list))

	count := 0
	for unam := range ul.list {
    	users[count] = unam
		count++
	}

	return users
}

func (ul *usermgt) GetToken(unam string) (string, error){

	user, ok:= ul.list[unam]
	if !ok {return "", fmt.Errorf("GetToken -- %s not a user!", unam)}

	token, _ := user["token"]

	return token, nil
}

func (ul *usermgt) UpdUser(unam string) (error){

	list := ul.list
	_, ok:= list[unam]
	if !ok {return fmt.Errorf("UpdUser -- user %s does not exist!", unam)}

	valMap := make(map[string]string)

	valMap["token"] = genToken()
	list[unam] = valMap
//fmt.Printf("dbg -- %d\n",len(list))
	ul.list = list
	return nil
}


func (ul *usermgt) AddUser(unam string) (error){

	list := ul.list
	_, ok:= list[unam]
	if ok {return fmt.Errorf("AddUser -- user %s already exists!", unam)}

	valMap := make(map[string]string)

	valMap["token"] = genToken()
	list[unam] = valMap
//fmt.Printf("dbg -- %d\n",len(list))
	ul.list = list
	return nil
}

func (ul *usermgt) RmUser(unam string) (error){

	list := ul.list
	_, ok:= list[unam]
	if !ok {return fmt.Errorf("RmUser -- user %s does not exist!", unam)}

	delete(list, unam)
//fmt.Printf("dbg -- %d\n",len(list))
	ul.list = list
	return nil
}



func genToken() (string) {

	length:=6
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (ul *usermgt) PrintList() {

	list := ul.list
	fmt.Printf("*** api list: %d providers ***\n", len(list))
	for nam, val := range list {
		fmt.Printf("user: %s\n", nam)
		for key, kval := range val {
			fmt.Printf(" %s : %s\n", key, kval)
		}
	}
	fmt.Printf("*** end api list providers ***\n")
}
