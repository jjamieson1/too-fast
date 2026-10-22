package game

type Upgrade int

const (
	Nitrous Upgrade = iota
	Scanner
	Turbo
	Tires
	numUpgrades
)

func (u Upgrade) String() string {
	return [...]string{"Nitrous", "Police Scanner", "Turbo", "Racing Tires"}[u]
}

// Config holds every tunable number, so balance changes can be tested
// without touching the engine. DefaultConfig matches RULES.md v0.3;
// V02Config is the earlier rule set, kept for comparison.
type Config struct {
	StartRep        int
	FuelMax         int
	HeatMax         int // reaching this = Busted
	BustedHeat      int // Heat after being Busted
	TrapSpeed       int // passing a trap with movement >= this gets you caught
	GasHeatDrop     int
	NitrousBonus    int
	NitrousRefill   int
	Price           [numUpgrades]int
	DaikokuDiscount int
	BridgeCost      int
	BridgeJump      int
	ShibuyaRefuel   bool
	ShibuyaBonus    int
	DuelStake       int
	HandLimit       int
	RepSpace        int
	Rep5Space       int

	// Optional variants from RULES.md §14.
	FuelDelivery     bool // pay FuelDeliveryCost, refuel in place, lose a turn
	FuelDeliveryCost int
	ThirstyNitrous   bool // firing Nitrous burns 1 extra fuel

	// Knobs added while testing balance fixes.
	MinFine        int               // smallest fine at a speed trap (0 = Heat-0 drivers get a warning)
	SeatRepBonus   int               // extra starting Rep per seat after the first
	SeatRep        []int             // explicit extra starting Rep by seat (overrides SeatRepBonus)
	ShibuyaCatchUp int               // Rep for each place behind the first car into Shibuya
	HeatRollMin    int               // natural rolls >= this add +1 Heat (v0.2: 6)
	ShortcutMax    int               // cap on how far Shortcut moves you; 0 = no cap
	PulledHeatDrop int               // Heat lost when pulled over; 0 = reset to 0
	TrapPassHeat   int               // Heat gained for driving past a trap without being caught
	OilOnPass      bool              // oil slicks also trigger when you drive over them
	RouteOverride  map[string]string // replace a route's space list (same codes as make_board.py)
	NarrowOverride map[string][][2]int

	BoostCounts    map[Boost]int
	ObstacleCounts map[Obstacle]int

	MaxRounds int // safety stop; a game that hits this counts as stuck
}

func DefaultConfig() Config {
	return Config{
		StartRep:         4,
		FuelMax:          5,
		HeatMax:          5,
		BustedHeat:       2,
		TrapSpeed:        7,
		GasHeatDrop:      2,
		NitrousBonus:     3,
		NitrousRefill:    2,
		Price:            [numUpgrades]int{Nitrous: 3, Scanner: 3, Turbo: 5, Tires: 6},
		DaikokuDiscount:  2,
		BridgeCost:       2,
		BridgeJump:       4,
		ShibuyaRefuel:    true,
		ShibuyaBonus:     0,
		DuelStake:        2,
		HandLimit:        3,
		RepSpace:         3,
		Rep5Space:        5,
		FuelDeliveryCost: 3,
		MinFine:          0,
		SeatRep:          []int{0, 0, 1, 2},
		ShortcutMax:      6,
		HeatRollMin:      6,
		BoostCounts:      copyMap(DefaultBoostCounts),
		ObstacleCounts:   copyMap(DefaultObstacleCounts),
		MaxRounds:        150,
	}
}

func copyMap[K comparable](m map[K]int) map[K]int {
	o := make(map[K]int, len(m))
	for k, v := range m {
		o[k] = v
	}
	return o
}

// V02AlleySeq is the Back Alleys route before v0.3 (three +5 Rep spaces).
const V02AlleySeq = ". R B R . O R5 L:Tsukiji B . G O B R . G R5 B . R O L:Ameyoko B G . R5 B O R . B R"

// V02Config returns the rules as they were in RULES.md v0.2.
func V02Config() Config {
	c := DefaultConfig()
	c.TrapSpeed = 6
	c.ShibuyaBonus = 3
	c.MinFine = 1
	c.SeatRep = nil
	c.ShortcutMax = 0
	c.RouteOverride = map[string]string{"C": V02AlleySeq}
	return c
}
