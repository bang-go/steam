package steamid

type AccountType int

// Steam account types.
const (
	AccountTypeInvalid AccountType = iota
	AccountTypeIndividual
	AccountTypeMultiSeat
	AccountTypeGameServer
	AccountTypeAnonGameServer
	AccountTypePending
	AccountTypeContentServer
	AccountTypeClan
	AccountTypeChat
	AccountTypeP2PSuperSeeder
	AccountTypeAnonUser
)

func getAccountType(t string) AccountType {
	if t == "I" {
		return AccountTypeInvalid
	} else if t == "U" {
		return AccountTypeIndividual
	} else if t == "M" {
		return AccountTypeMultiSeat
	} else if t == "G" {
		return AccountTypeGameServer
	} else if t == "A" {
		return AccountTypeAnonGameServer
	} else if t == "P" {
		return AccountTypePending
	} else if t == "C" {
		return AccountTypeContentServer
	} else if t == "g" {
		return AccountTypeClan
	} else if t == "T" {
		return AccountTypeChat
	} else if t == "a" {
		return AccountTypeP2PSuperSeeder
	} else {
		return AccountTypeInvalid
	}
}

func getAccountString(t AccountType) string {
	if t == AccountTypeInvalid {
		return "I"
	} else if t == AccountTypeIndividual {
		return "U"
	} else if t == AccountTypeMultiSeat {
		return "M"
	} else if t == AccountTypeGameServer {
		return "G"
	} else if t == AccountTypeAnonGameServer {
		return "A"
	} else if t == AccountTypePending {
		return "P"
	} else if t == AccountTypeContentServer {
		return "C"
	} else if t == AccountTypeClan {
		return "g"
	} else if t == AccountTypeChat {
		return "T"
	} else if t == AccountTypeP2PSuperSeeder {
		return "a"
	} else {
		return "I"
	}
}

func (t AccountType) ToString() string {
	return getAccountString(t)
}
