package game

import "math/rand/v2"

type Boost int

const (
	EmptyStreet Boost = iota
	BribeACop
	OilSlick
	GreenWave
	DriftKing
	ScannerTip
	SlipstreamCard
	Shortcut
	FakePlates
	SponsorDeal
	MechanicOnCall
	JerryCan
	numBoosts
)

var boostNames = [...]string{"Empty Street", "Bribe a Cop", "Oil Slick", "Green Wave", "Drift King", "Scanner Tip",
	"Slipstream", "Shortcut", "Fake Plates", "Sponsor Deal", "Mechanic on Call", "Jerry Can"}

func (b Boost) String() string { return boostNames[b] }

type Obstacle int

const (
	FlatTire Obstacle = iota
	FuelLeak
	PoliceCheckpoint
	TrafficJam
	Construction
	WrongTurn
	EngineTrouble
	Rain
	RivalCrew
	TowZone
	CaughtOnCamera
	numObstacles
)

var obstacleNames = [...]string{"Flat Tire", "Fuel Leak", "Police Checkpoint", "Traffic Jam", "Construction",
	"Wrong Turn", "Engine Trouble", "Rain", "Rival Crew", "Tow Zone", "Caught on Camera"}

func (o Obstacle) String() string { return obstacleNames[o] }

// Card counts from CARDS.md.
var DefaultBoostCounts = map[Boost]int{
	EmptyStreet: 2, BribeACop: 3, OilSlick: 3, GreenWave: 2, DriftKing: 1, ScannerTip: 2,
	SlipstreamCard: 2, Shortcut: 2, FakePlates: 2, SponsorDeal: 2, MechanicOnCall: 2, JerryCan: 2,
}

var DefaultObstacleCounts = map[Obstacle]int{
	FlatTire: 3, FuelLeak: 3, PoliceCheckpoint: 3, TrafficJam: 3, Construction: 2, WrongTurn: 2,
	EngineTrouble: 2, Rain: 2, RivalCrew: 2, TowZone: 2, CaughtOnCamera: 1,
}

// Deck is a draw pile plus discard pile; the discards are reshuffled when it runs out.
type Deck[T comparable] struct {
	draw, discard []T
	rng           *rand.Rand
}

func NewDeck[T comparable](counts map[T]int, rng *rand.Rand, order []T) *Deck[T] {
	d := &Deck[T]{rng: rng}
	for _, c := range order { // fixed order so a seed always gives the same game
		for i := 0; i < counts[c]; i++ {
			d.draw = append(d.draw, c)
		}
	}
	d.shuffle()
	return d
}

func (d *Deck[T]) shuffle() {
	d.rng.Shuffle(len(d.draw), func(i, j int) { d.draw[i], d.draw[j] = d.draw[j], d.draw[i] })
}

// Draw returns false only if every card is in a player's hand.
func (d *Deck[T]) Draw() (T, bool) {
	if len(d.draw) == 0 {
		d.draw, d.discard = d.discard, nil
		d.shuffle()
	}
	var zero T
	if len(d.draw) == 0 {
		return zero, false
	}
	c := d.draw[len(d.draw)-1]
	d.draw = d.draw[:len(d.draw)-1]
	return c, true
}

func (d *Deck[T]) Discard(c T) { d.discard = append(d.discard, c) }

func allBoosts() []Boost {
	var o []Boost
	for i := Boost(0); i < numBoosts; i++ {
		o = append(o, i)
	}
	return o
}

func allObstacles() []Obstacle {
	var o []Obstacle
	for i := Obstacle(0); i < numObstacles; i++ {
		o = append(o, i)
	}
	return o
}
