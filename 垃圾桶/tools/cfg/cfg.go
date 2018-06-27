package cfg

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	_map  map[string]string
	_lock sync.RWMutex
)

func init() {
	Reload("E:/GitHub/GameServers/NanoCopy")
}

// Get .
func Get() map[string]string {
	_lock.RLock()
	defer _lock.RUnlock()
	return _map
}

// Reload .
func Reload(basePath string) {
	var baseConfigPath = basePath + "/data/server/config_base.ini"     //os.Getenv("GOGAMESERVER_PATH") + "/data/server/config_base.ini"
	var serverConfigPath = basePath + "/data/server/config_server.ini" //os.Getenv("GOGAMESERVER_PATH") + "/data/server/config_server.ini"

	_lock.Lock()
	_map = make(map[string]string)
	_loadConfig(baseConfigPath)
	_loadConfig(serverConfigPath)
	_lock.Unlock()
}

func _loadConfig(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(path, err)
		return
	}

	re := regexp.MustCompile(`[\t ]*([0-9A-Za-z_]+)[\t ]*=[\t ]*([^\t\n\f\r# ]+)[\t #]*`)

	// using scanner to read config file
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// expression match
		slice := re.FindStringSubmatch(line)

		if slice != nil {
			_map[slice[1]] = slice[2]
		}
	}

	return
}

// GetValue .
func GetValue(key string) string {
	config := Get()
	return config[key]
}

// GetUint16 .
func GetUint16(key string) uint16 {
	result, _ := strconv.Atoi(GetValue(key))
	return uint16(result)
}

// GetInt .
func GetInt(key string) int {
	result, _ := strconv.Atoi(GetValue(key))
	return result
}
