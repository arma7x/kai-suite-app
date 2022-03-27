package global

import (
	"crypto/rand"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"net"
	"strings"
	"errors"
	"strconv"
	"fyne.io/fyne/v2"
	"github.com/tidwall/buntdb"
	"kai-suite/types"
)

var (
	ROOT_PATH string
	APP fyne.App
	WINDOW fyne.Window
	VISIBILITY = true
	THEME = 0
	CONTACTS_DB *buntdb.DB
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	ROOT_PATH = filepath.Dir(ex)
	CONTACTS_DB, err = buntdb.Open(ResolvePath("db/contacts.db"))
	if err != nil {
		log.Warn(err)
	}
	CONTACTS_DB.CreateIndex("people_local", "local:people:*", buntdb.IndexString)
	CONTACTS_DB.CreateIndex("metadata_local", "metadata:local:people:*", buntdb.IndexString)
}

func ResolvePath(dirs... string) string {
	return filepath.FromSlash(fmt.Sprintf("%s%s%s", ROOT_PATH, "/", strings.Join(dirs, "/")))
}

func CheckIPAddress(ip, port string) (string, error) {
	ipAddr := strings.Join([]string{ip, port}, ":")
	if net.ParseIP(ip) == nil {
		return ipAddr, errors.New(strings.Join([]string{"Error:", ip, "is invalid IP address"}, " "))
	}
	p, err := strconv.Atoi(port);
	if err != nil {
		return ipAddr, errors.New(strings.Join([]string{"Error:", port, "is invalid port number"}, " "))
	}
	if (p <= 1024) {
		return ipAddr, errors.New(strings.Join([]string{"Error:", "Port", port, "must greater than", "1024"}, " "))
	}
	return ipAddr, nil 
}

func InitDatabaseIndex(accounts map[string]*types.UserInfoAndToken) {
	for key, _ := range accounts {
		index := strings.Join([]string{key, "people", "*"}, ":")
		indexName := strings.Join([]string{"people", key}, "_")
		CONTACTS_DB.CreateIndex(indexName, index, buntdb.IndexString)
		metadataIndex := strings.Join([]string{"metadata", key, "people", "*"}, ":")
		metadataIndexName := strings.Join([]string{"metadata", key}, "_")
		CONTACTS_DB.CreateIndex(metadataIndexName, metadataIndex, buntdb.IndexString)
	}
}

func RandomID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X%X%X%X%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
