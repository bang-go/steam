package steamid

type Universe int

// Steam universes. Each universe is a self-contained Steam instance.
const (
	UniverseInvalid  Universe = iota //steamID is Invalid
	UniversePublic                   //steamID is Public
	UniverseBeta                     // steamID is Beta
	UniverseInternal                 //steamID is Internal
	UniverseDev                      // steamID is Dev
)
