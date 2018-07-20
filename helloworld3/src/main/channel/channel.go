package channel

const (
	maxWorldChannel  = 200
	maxNumPerChannel = 700
)

// World 世界频道
var World = newChannelMap()

// Team 组队频道
var Team = newChannelMap()

// Organize  社团频道
var Organize = newChannelMap()

// GetoutWorldChannel 离开世界频道(下线)
func GetoutWorldChannel(pid int64, worldID int32) {
	World.RLock()
	v, ok := World.Data[worldID]
	World.RUnlock()
	if ok {
		v.Del(pid)
	}
}

// PutIntoWorldChannel 加入世界频道(上线)
func PutIntoWorldChannel(pid int64, targetID int32) int32 {
	if targetID == 0 { //系统分配
		World.Lock()
		length := len(World.Data)
		var retval int32
		for i := 1; i <= length+1; i++ {
			v, ok := World.Data[int32(i)]
			if ok {
				if v.Length() < maxNumPerChannel {
					v.Add(pid)
					retval = int32(i)
					break
				}
			} else {
				//添加世界频道
				World.Add(int32(i), pid)
				retval = int32(i)
				break
			}
		}
		World.Unlock()
		return retval
	} else if targetID > 0 && targetID < maxWorldChannel { //玩家分配
		World.RLock()
		v, ok := World.Data[targetID]
		World.RUnlock()
		if ok == true {
			v.Add(pid)
		} else {
			World.Lock()
			v, ok = World.Data[targetID]
			if ok == true {
				v.Add(pid)
			} else {
				World.Add(targetID, pid)
			}
			World.Unlock()
		}
		return targetID
	}
	return 0
}

// PutIntoTeamChannel 加入队伍频道
func PutIntoTeamChannel(pid int64, teamID int32) int32 {
	Team.RLock()
	v, ok := Team.Data[teamID]
	Team.RUnlock()
	if ok == true {
		v.Add(pid)
	} else {
		Team.Lock()
		v, ok := Team.Data[teamID]
		if ok == true {
			v.Add(pid)
		} else {
			Team.Add(teamID, pid)
		}
		Team.Unlock()
	}
	return teamID
}

// GetoutTeamChannel 离开队伍频道(离开队伍)
func GetoutTeamChannel(pid int64, teamID int32) {
	Team.RLock()
	v, ok := Team.Data[teamID]
	Team.RUnlock()
	if ok {
		v.Del(pid)
		l := v.Length()
		if l == 0 {
			Team.Lock()
			v, ok = Team.Data[teamID]
			if ok && v.Length() == 0 {
				Team.Del(teamID)
			}
			Team.Unlock()
		}
	}
}

// PutIntoOrganizeChannel 加入社团频道(上线))
func PutIntoOrganizeChannel(pid int64, ograID int32) int32 {
	Organize.RLock()
	v, ok := Organize.Data[ograID]
	Organize.RUnlock()
	if ok == true {
		v.Add(pid)
	} else {
		Organize.Lock()
		v, ok := Organize.Data[ograID]
		if ok == true {
			v.Add(pid)
		} else {
			Organize.Add(ograID, pid)
		}
		Organize.Unlock()
	}
	return ograID
}

// GetoutOrganizeChannel 离开社团频道(下线)
func GetoutOrganizeChannel(pid int64, orgaID int32) {
	Organize.RLock()
	v, ok := Organize.Data[orgaID]
	Organize.RUnlock()
	if ok {
		v.Del(pid)
		l := v.Length()
		if l == 0 {
			Organize.Lock()
			v, ok = Organize.Data[orgaID]
			if ok && v.Length() == 0 {
				Organize.Del(orgaID)
			}
			Organize.Unlock()
		}
	}
}
