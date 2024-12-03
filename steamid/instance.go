package steamid

type Instance int

// Steam allows 3 simultaneous user account instances right now
const (
	InstanceAll Instance = iota
	InstanceDesktop
	InstanceConsole
	InstanceWeb = 4
)
const (
	BaseGID = uint64(103582791429521408)
	BaseSID = uint64(76561197960265728)
)

// Special flags for Chat accounts - they go in the top 8 bits of the steam ID's "instance", leaving 12 for the actual instances.
const (
	InstanceMask         = 0x000FFFFF
	InstanceFlagClan     = (InstanceMask + 1) >> 1
	InstanceFlagLobby    = (InstanceMask + 1) >> 2
	InstanceFlagMMSLobby = (InstanceMask + 1) >> 3
)
