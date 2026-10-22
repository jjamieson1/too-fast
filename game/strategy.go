package game

import "math"

type NitrousMode int

const (
	NitrousNever NitrousMode = iota
	NitrousAlways
	NitrousSmart // only when no speed trap is in range, or with a scanner, or to finish
)

// Strategy is a bot's decision profile. Every bot plays cards sensibly;
// the fields below are the choices that differ between play styles.
type Strategy struct {
	Name        string
	Leg1        string // "A", "B", "C", or "" for random
	Leg2        string // "D", "E", "" = D with a Police Scanner else E, "?" = random
	BuyOrder    []Upgrade
	Reserve     int // Rep kept back for fines after buying
	FuelStop    int // stop at a passing gas station when fuel <= this; -1 = plan ahead
	NitrousMode NitrousMode
	UseBridge   bool
	FixFlatDist int // park to fix a flat if the next gas station is farther than this; 0 = never
}

func (s *Strategy) chooseRoute(g *Game, p *Player) string {
	if p.Leg == 1 {
		if s.Leg1 == "" {
			return g.Board.Leg1[g.Rng.IntN(len(g.Board.Leg1))]
		}
		return s.Leg1
	}
	switch s.Leg2 {
	case "":
		if p.Up[Scanner] {
			return "D"
		}
		return "E"
	case "?":
		return g.Board.Leg2[g.Rng.IntN(len(g.Board.Leg2))]
	}
	return s.Leg2
}

// buyUpgrades buys in priority order, saving up for the next item rather
// than skipping to a cheaper one. With a discount (Daikoku PA) it buys at most one.
func (s *Strategy) buyUpgrades(g *Game, p *Player, discount int) {
	for _, u := range s.BuyOrder {
		if p.Up[u] {
			continue
		}
		price := max(1, g.Cfg.Price[u]-discount)
		if p.Rep-price < s.Reserve {
			return
		}
		p.Rep -= price
		p.Up[u] = true
		if u == Nitrous {
			p.NitrousFull = true
		}
		p.Stats.UpgradeTurn[u] = p.Stats.Turns
		g.logf("%s buys %s for %d", p.Name, u, price)
		if discount > 0 {
			return
		}
	}
}

func (s *Strategy) startOfTurn(g *Game, p *Player) {
	if p.has(FakePlates) && p.Heat >= 3 {
		g.play(p, FakePlates)
		p.Heat = max(0, p.Heat-3)
	}
	if p.has(OilSlick) {
		if q := g.leader(p); q != nil && q.Route != "" && g.Remaining(q) < g.Remaining(p) {
			if l, ok := g.oilSpot(q); ok {
				g.play(p, OilSlick)
				g.Oil[l] = true
				g.logf("%s drops an oil slick in front of %s", p.Name, q.Name)
			}
		}
	}
}

// oilSpot picks the empty Plain space where target q is most likely to land.
func (g *Game) oilSpot(q *Player) (Loc, bool) {
	r := g.Board.Routes[q.Route]
	want := int(math.Round(g.expectedMove(q)))
	best, bestD := Loc{}, 99
	for pos := q.Pos + 1; pos <= min(r.Len(), q.Pos+9); pos++ {
		l := Loc{q.Leg, q.Route, pos}
		if r.At(pos).Kind != Plain || g.Oil[l] || g.someoneAt(l) {
			continue
		}
		if d := abs(pos - q.Pos - want); d < bestD {
			best, bestD = l, d
		}
	}
	return best, bestD < 99
}

func (g *Game) someoneAt(l Loc) bool {
	for _, q := range g.Players {
		if !q.Finished && g.loc(q) == l {
			return true
		}
	}
	return false
}

func (g *Game) expectedMove(p *Player) float64 {
	m := 3.5
	if p.Up[Turbo] {
		m++
	}
	if p.Up[Tires] {
		m++
	}
	if p.Flat {
		m -= 2
	}
	return math.Max(1, m)
}

func (s *Strategy) fixFlatNow(g *Game, p *Player) bool {
	return s.FixFlatDist > 0 && g.distToFuel(p, p.Pos) > s.FixFlatDist
}

func (s *Strategy) fireNitrous(g *Game, p *Player) bool {
	switch s.NitrousMode {
	case NitrousAlways:
		return true
	case NitrousSmart:
		if p.Leg == 2 && p.Route != "" && g.Remaining(p) <= g.expectedMove(p)+float64(g.Cfg.NitrousBonus)+1 {
			return true // go for the win
		}
		return p.Up[Scanner] || p.tipActive || !g.trapAhead(p, 12)
	}
	return false
}

func (s *Strategy) stopForGas(g *Game, p *Player) bool {
	if s.FuelStop >= 0 {
		return p.Fuel <= s.FuelStop
	}
	// Plan ahead: can the fuel left reach the next refuel point beyond this one?
	need := float64(g.distToFuel(p, p.Pos))
	reach := float64(p.Fuel) * g.expectedMove(p) * 0.8
	return reach < need
}

func (s *Strategy) wantsDelivery(g *Game, p *Player) bool {
	if p.Route == "" {
		return false
	}
	r := g.Board.Routes[p.Route]
	back := 0
	for pos := p.Pos; pos >= 1 && r.At(pos).Kind != Gas; pos-- {
		back++
	}
	return back >= 5
}

// spaceValue is a rough "how good is it to end here" score.
func (g *Game) spaceValue(p *Player, pos int) float64 {
	r := g.Board.Routes[p.Route]
	if pos > r.Len() {
		if p.Leg == 2 {
			return 100
		}
		return 2
	}
	if g.Oil[Loc{p.Leg, p.Route, pos}] && !p.Up[Tires] {
		return -3
	}
	switch r.At(pos).Kind {
	case Trap:
		if p.tipActive {
			return 0
		}
		return -1 - float64(max(1, p.Heat))
	case ObstacleSpace:
		return -1.5
	case Rep3, Rep5:
		return 2
	case BoostSpace:
		return 1
	case Gas:
		if p.Fuel <= 2 {
			return 2
		}
		return 0.5
	case Landmark:
		return 0.5
	}
	return 0
}

// useAfterMove plays Empty Street or Shortcut when it gets you off a bad
// space or wins the race. Returns true if a card moved the car.
func (s *Strategy) useAfterMove(g *Game, p *Player) bool {
	here := g.spaceValue(p, p.Pos)
	r := g.Board.Routes[p.Route]
	if p.has(EmptyStreet) {
		dest := g.spaceValue(p, p.Pos+3)
		if dest >= 100 || (here < 0 && dest > here) || len(p.Hand) >= g.Cfg.HandLimit {
			g.play(p, EmptyStreet)
			g.moveForward(p, 3, false)
			return true
		}
	}
	if p.has(Shortcut) {
		d := g.distToFuel(p, p.Pos)
		if p.Pos+d > r.Len() { // no station ahead: card says move 3
			d = 3
		}
		if g.Cfg.ShortcutMax > 0 {
			d = min(d, g.Cfg.ShortcutMax)
		}
		if d >= 4 || (here < 0 && d >= 2) || p.Fuel <= 1 {
			g.play(p, Shortcut)
			g.moveForward(p, d, false)
			return true
		}
	}
	return false
}

func (s *Strategy) cancelObstacle(g *Game, p *Player, c Obstacle) bool {
	switch c {
	case TowZone, FlatTire, TrafficJam, WrongTurn, Construction:
		return true
	case PoliceCheckpoint:
		return p.Heat >= 2
	case FuelLeak:
		return p.Fuel <= 2
	}
	return false
}

var cardKeep = map[Boost]int{
	BribeACop: 10, MechanicOnCall: 9, EmptyStreet: 8, JerryCan: 7, Shortcut: 6,
	ScannerTip: 5, GreenWave: 5, FakePlates: 4, DriftKing: 4, SlipstreamCard: 3, OilSlick: 3, SponsorDeal: 11,
}

func (s *Strategy) worstCard(hand []Boost) Boost {
	w := hand[0]
	for _, c := range hand[1:] {
		if cardKeep[c] < cardKeep[w] {
			w = c
		}
	}
	return w
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
