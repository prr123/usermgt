package apiLib

import (
	"fmt"
	"os"
	"github.com/goccy/go-yaml"
)

type provMap map[string]string

type usermgt struct {
	list map[string]provMap
	dbg bool
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


func (ul *usermgt) InitUserList(yamlFil string) (err error) {

	ul.list = make(map[string]provMap)

	ldat, err := os.ReadFile(yamlFil)
	if err != nil {return fmt.Errorf("Read List: %v", err)}

	err = yaml.Unmarshal(ldat,ul.list)
	if err != nil {return fmt.Errorf("UnMarshal: %v", err)}

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
