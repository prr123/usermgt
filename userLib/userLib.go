package userLib

import (
	"fmt"
	"os"
	"bytes"
	"github.com/goccy/go-yaml"
)

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

func ProcCmd(cmdStr string) (error) {

	switch cmdStr {
	case "list":


	case "add":


	case "rm":


	case "upd":

	default:
		return fmt.Errorf("unknown command: %s\n", cmdStr)

	}

	return nil
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

func (ul *usermgt) GetAllUsers() (users []string) {

	users = make([]string, len(ul.list))

	count := 0
	for unam := range ul.list {
    	users[count] = unam
		count++
	}

	return users
}
func (ul *usermgt) GetUserToken(unam string) (string, bool){

	user, ok:= ul.list[unam]
	if !ok {return "", false}

	token, _ := user["token"]

	return token, true
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
