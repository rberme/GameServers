package global

import "strings"

var (
	// ServerName 服务器名
	ServerName string
	// ServerID 服务器ID
	ServerID uint32 // = 0
)

// GetTrueServerName ...
func GetTrueServerName() string {
	return strings.Split(ServerName, "[")[0]
}

// LocalServer ...
func LocalServer() bool {
	return GetTrueServerName() == "LocalServer"
}

// IsWorldServer ...
func IsWorldServer() bool {
	return GetTrueServerName() == "WorldServer"
}

// IsGameServer ...
func IsGameServer() bool {
	return GetTrueServerName() == "GameServer"
}

// IsLoginServer ...
func IsLoginServer() bool {

	return GetTrueServerName() == "LoginServer"

}
