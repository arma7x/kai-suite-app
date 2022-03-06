package global

import (
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
)

var (
	ROOT_PATH string
	WINDOW fyne.Window
	DB *buntdb.DB
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	ROOT_PATH = filepath.Dir(ex)
	var errdB error
	DB, errdB = buntdb.Open(ResolvePath("database.db"))
	if errdB != nil {
		log.Fatal(errdB)
	}
	DB.CreateIndex("key", "people:*", buntdb.IndexString)
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
