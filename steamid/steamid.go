package steamid

import (
	"fmt"
	"github.com/bang-go/util"
	"math"
	"regexp"
	"strconv"
)

type SteamID interface {
	IsValid() bool
	RenderSteamID64() int64
	RenderSteamID2(format bool) string
	RenderSteamID3() string
	String() string
	GetAccountID() int32
	GetInstance() Instance
	GetUniverse() Universe
	GetAccountType() AccountType
}

type steamIDEntity struct {
	universe    Universe
	accountType AccountType
	instance    Instance
	accountID   int32
}

var (
	regSteamId2 = regexp.MustCompile(`^STEAM_([0-5]):([0-1]):([0-9]+)$`)
	regSteamId3 = regexp.MustCompile(`^\[([a-zA-Z]):([0-5]):([0-9]+)(:[0-9]+)?\]$`)
)

// New 新建steamID 支持 64位id：76561199181487706， steam2：STEAM_0:0:610610989，steam3：[U:1:1221221978]，accountid:1221221978
func New(raw string) (sid SteamID, err error) {
	if regSteamId2.MatchString(raw) { //steam2
		match := regSteamId2.FindStringSubmatch(raw)
		return getSteamIdBySteam2(match)
	} else if regSteamId3.MatchString(raw) { //steam3
		match := regSteamId3.FindStringSubmatch(raw)
		return getSteamIdBySteam3(match)
	}
	var aid uint64
	aid, err = strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return
	}
	if aid < BaseSID {
		//accountID
		sid = getSteamIdByAccountID(aid)
		return
	}
	//64id
	sid = getSteamIdBy64Id(aid)
	return
}

func getSteamIdBySteam2(match []string) (sid *steamIDEntity, err error) {
	sid = &steamIDEntity{}
	iUniverse, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		return
	}
	authServer, err := strconv.ParseUint(match[2], 10, 64)
	if err != nil {
		return
	}
	accountID, err := strconv.ParseUint(match[3], 10, 64)
	if err != nil {
		return
	}
	// Games before orange box used to incorrectly display universe as 0
	uni := util.If(iUniverse > 0, Universe(iUniverse), UniversePublic)
	accountID = (accountID << 1) | authServer
	sid.universe = uni
	sid.accountType = AccountTypeIndividual
	sid.accountID = int32(accountID)
	sid.instance = InstanceDesktop
	return
}
func getSteamIdBySteam3(match []string) (sid *steamIDEntity, err error) {
	sid = &steamIDEntity{}
	iUniverse, err := strconv.ParseUint(match[2], 10, 64)
	if err != nil {
		return
	}
	accountID, err := strconv.ParseUint(match[3], 10, 64)
	if err != nil {
		return
	}
	sid.universe = Universe(iUniverse)
	sid.accountID = int32(accountID)

	typeChar := match[1]
	if len(match[4]) > 0 {
		var iInstance int
		iInstance, err = strconv.Atoi(string(match[4][1]))
		if err != nil {
			return
		}
		sid.instance = Instance(iInstance)
	} else if typeChar == "U" {
		sid.instance = InstanceDesktop
	}

	if typeChar == "C" {
		sid.instance = Instance(int(sid.instance) | int(InstanceFlagClan))
	} else if typeChar == "L" {
		sid.instance = Instance(int(sid.instance) | int(InstanceFlagLobby))
		sid.accountType = AccountTypeChat
	} else {
		sid.accountType = getAccountType(typeChar)
	}
	return
}
func getSteamIdBy64Id(rawId uint64) *steamIDEntity {
	// 76561199181487706->1221221978
	sid := &steamIDEntity{accountID: 0, instance: InstanceAll, accountType: AccountTypeInvalid, universe: UniverseInvalid}
	sid.accountID = int32((rawId & 0xFFFFFFFF) >> 0)
	sid.instance = Instance(rawId >> 32 & 0xFFFFF)
	sid.accountType = AccountType(rawId >> 52 & 0xF)
	sid.universe = Universe(rawId >> 56)
	return sid
}
func getSteamIdByAccountID(rawId uint64) *steamIDEntity {
	// 1221221978 -> 76561199181487706
	sid := &steamIDEntity{accountID: 0, instance: InstanceAll, accountType: AccountTypeInvalid, universe: UniverseInvalid}
	sid.universe = UniversePublic
	sid.accountType = AccountTypeIndividual
	sid.instance = InstanceDesktop
	sid.accountID = int32(rawId)
	return sid
}

// RenderSteamID64 获取steamID64 76561199181487706
func (s *steamIDEntity) RenderSteamID64() int64 {
	var data = util.VarAddr(int64(0))
	set(data, 56, 0xFF, int64(s.universe))       //0xFF=255
	set(data, 52, 0xF, int64(s.accountType))     //0xF=15
	set(data, 32, 0xFFFFF, int64(s.instance))    //0xFFFFF=1048575
	set(data, 0, 0xFFFFFFFF, int64(s.accountID)) //0xFFFFFFFF=4294967295
	return *data
}

// RenderSteamID2 获取steamID2  STEAM_0:0:610610989
func (s *steamIDEntity) RenderSteamID2(format bool) string {
	if s.accountType != AccountTypeIndividual {
		return ""
	}
	uni := s.universe
	if !format && uni == 1 {
		uni = 0
	}
	return fmt.Sprintf("STEAM_%d:%d:%d", uni, s.accountID&1, int64(math.Floor(float64(s.accountID)/2)))
}

// RenderSteamID3 获取steamID3 [U:1:1221221978]
func (s *steamIDEntity) RenderSteamID3() string {
	char := s.accountType.ToString()
	if s.instance&InstanceFlagClan > 0 {
		char = "c"
	} else if s.instance&InstanceFlagLobby > 0 {
		char = "L"
	}
	doInstance := s.accountType == AccountTypeAnonGameServer ||
		s.accountType == AccountTypeMultiSeat ||
		(s.accountType == AccountTypeIndividual && s.instance != InstanceDesktop)
	if !doInstance {
		return fmt.Sprintf("[%s:%d:%d]", char, s.universe, s.accountID)
	} else {
		return fmt.Sprintf("[%s:%d:%d:%d]", char, s.universe, s.accountID, s.instance)
	}
}

// IsValid Check whether this steamIDEntity is valid (according to Steam's rules)
func (s *steamIDEntity) IsValid() bool {
	if s.accountType <= AccountTypeInvalid || s.accountType > AccountTypeAnonUser {
		return false
	}

	if s.universe <= UniverseInvalid || s.universe > UniverseDev {
		return false
	}

	if s.accountType == AccountTypeIndividual && (s.accountID == 0 || s.instance > InstanceWeb) {
		return false
	}

	if s.accountType == AccountTypeClan && (s.accountID == 0 || s.instance != InstanceAll) {
		return false
	}

	if s.accountType == AccountTypeGameServer && s.accountID == 0 {
		return false
	}

	return true
}
func (s *steamIDEntity) String() string {
	return fmt.Sprintf("AccountID: %d, AccountType: %d, Universe: %d, Instance: %d", s.accountID, s.accountType, s.universe, s.instance)
}

func (s *steamIDEntity) GetAccountID() int32 {
	return s.accountID
}
func (s *steamIDEntity) GetInstance() Instance {
	return s.instance
}
func (s *steamIDEntity) GetUniverse() Universe {
	return s.universe
}

func (s *steamIDEntity) GetAccountType() AccountType {
	return s.accountType
}

func set(data *int64, BitOffset int64, ValueMask int64, Value int64) {
	// 左移 ValueMask
	shiftedValueMask := ValueMask << uint(BitOffset)
	// 计算 maskComplement
	maskComplement := ^shiftedValueMask
	// 计算 shiftedValue
	shiftedValue := (Value & ValueMask) << uint(BitOffset)
	// 更新 Data
	*data = (*data & maskComplement) | shiftedValue
}
