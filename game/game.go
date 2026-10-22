package game

import (
	"fmt"
	"math/rand/v2"
)

type ParkCause int

const (
	ParkBusted ParkCause = iota
	ParkFine
	ParkTow
	ParkFlatFix
	ParkDelivery
	NumParkCauses
)

func (c ParkCause) String() string {
	return [...]string{"Busted", "Couldn't pay fine", "Tow Zone", "Fixing flat", "Fuel delivery"}[c]
}

// Stats collects everything the simulator reports on, per player per game.
type Stats struct {
	Turns, ParkedTurns  int
	ParkCauses          [NumParkCauses]int
	OutOfGas            int
	OutOfGasSpaces      int // spaces pushed back when out of gas
	FuelStops           int // voluntary stops at a passing gas station
	PulledOver, Fines   int
	TrapPassCatches     int
	Busted, Bribes      int
	DuelsWon, DuelsLost int
	Obstacles           [numObstacles]int
	Boosts              [numBoosts]int
	UpgradeTurn         [numUpgrades]int // turn number bought, 0 = never
	RepEarned           int
	ShibuyaTurn         int
	ShibuyaRank         int // 1 = first to arrive
	ShibuyaRep          int
	ShibuyaUpgrades     int
	Leg1, Leg2          string
	MaxStall, stall     int // turns in a row without forward progress
	BrokeTurns          int // turns started with 0 Rep
	SpacesBack          int
	NitrousFired        int
	OilSpinouts         int
}

type Player struct {
	ID       int
	Name     string
	Strat    *Strategy
	Leg      int    // 1 or 2
	Route    string // "" = on the hub (Start in leg 1, Shibuya in leg 2)
	Pos      int    // 1-based space on Route
	Finished bool

	Fuel, Heat, Rep int
	Up              [numUpgrades]bool
	NitrousFull     bool
	Hand            []Boost
	Parked          bool
	Flat            bool
	NextMod         int
	Jammed          bool
	tipActive       bool

	Stats Stats
}

func (p *Player) Upgrades() int {
	n := 0
	for _, u := range p.Up {
		if u {
			n++
		}
	}
	return n
}

func (p *Player) has(b Boost) bool {
	for _, c := range p.Hand {
		if c == b {
			return true
		}
	}
	return false
}

type Loc struct {
	Leg   int
	Route string
	Pos   int
}

type Game struct {
	Cfg       Config
	Board     *Board
	Rng       *rand.Rand
	Players   []*Player
	Boosts    *Deck[Boost]
	Obstacles *Deck[Obstacle]
	Oil       map[Loc]bool
	Shibuya   int // cars that have reached Shibuya
	Round     int
	Winner    int // player ID, -1 while racing (or if the game timed out)
	Log       func(string)
}

func New(cfg Config, strats []*Strategy, seed uint64) *Game {
	rng := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
	g := &Game{
		Cfg: cfg, Board: NewBoardWith(cfg.RouteOverride, cfg.NarrowOverride), Rng: rng, Oil: map[Loc]bool{}, Winner: -1,
		Boosts:    NewDeck(cfg.BoostCounts, rng, allBoosts()),
		Obstacles: NewDeck(cfg.ObstacleCounts, rng, allObstacles()),
	}
	for i, s := range strats {
		g.Players = append(g.Players, &Player{
			ID: i, Name: fmt.Sprintf("P%d %s", i+1, s.Name), Strat: s,
			Leg: 1, Fuel: cfg.FuelMax, Rep: cfg.StartRep + seatRep(cfg, i),
		})
	}
	return g
}

func seatRep(cfg Config, seat int) int {
	if cfg.SeatRep != nil {
		if seat < len(cfg.SeatRep) {
			return cfg.SeatRep[seat]
		}
		return 0
	}
	return seat * cfg.SeatRepBonus
}

func (g *Game) logf(format string, a ...any) {
	if g.Log != nil {
		g.Log(fmt.Sprintf(format, a...))
	}
}

func (g *Game) d6() int { return g.Rng.IntN(6) + 1 }

// Run plays until someone wins or MaxRounds is hit.
func (g *Game) Run() {
	for g.Winner < 0 && g.Round < g.Cfg.MaxRounds {
		g.Round++
		g.logf("── Round %d ──", g.Round)
		for _, p := range g.Players {
			g.turn(p)
			if g.Winner >= 0 {
				return
			}
		}
	}
}

func (g *Game) turn(p *Player) {
	s := p.Strat
	p.Stats.Turns++
	if p.Rep == 0 {
		p.Stats.BrokeTurns++
	}
	before := g.Remaining(p)
	defer func() { g.trackStall(p, before); p.tipActive = false }()

	if p.Parked {
		p.Parked = false
		p.Stats.ParkedTurns++
		g.logf("%s is parked and sits out", p.Name)
		return
	}
	slip := g.slipstreaming(p)

	// Step 1: before you roll.
	s.startOfTurn(g, p)
	if p.Flat && s.fixFlatNow(g, p) {
		p.Flat = false
		g.Obstacles.Discard(FlatTire)
		p.Stats.ParkedTurns++
		p.Stats.ParkCauses[ParkFlatFix]++
		g.logf("%s spends the turn fixing a flat", p.Name)
		return
	}
	s.buyUpgrades(g, p, 0)
	if p.Up[Nitrous] && !p.NitrousFull && p.Rep >= g.Cfg.NitrousRefill+s.Reserve && s.NitrousMode != NitrousNever {
		p.Rep -= g.Cfg.NitrousRefill
		p.NitrousFull = true
	}
	bonus := 0
	if slip {
		bonus++
	}
	if p.has(ScannerTip) && !p.Up[Scanner] && g.trapAhead(p, 12) {
		g.play(p, ScannerTip)
		p.tipActive = true
	}
	if p.has(GreenWave) {
		g.play(p, GreenWave)
		bonus += 2
	}
	drift := false
	if p.has(DriftKing) && p.Heat < g.Cfg.HeatMax-2 {
		g.play(p, DriftKing)
		drift = true
	}
	if p.has(SlipstreamCard) {
		if q, d := g.nearestAhead(p); q != nil && d > 1 && d <= 6 {
			g.play(p, SlipstreamCard)
			g.logf("%s slipstreams up behind %s", p.Name, q.Name)
			g.moveForward(p, d-1, false)
		}
	}

	// Step 2: burn fuel, roll and move.
	if p.Fuel == 0 && p.has(JerryCan) {
		g.play(p, JerryCan)
		p.Fuel = min(g.Cfg.FuelMax, p.Fuel+2)
	}
	if p.Fuel == 0 {
		g.outOfGas(p)
		return
	}
	p.Fuel--
	fired := false
	if p.Up[Nitrous] && p.NitrousFull && s.fireNitrous(g, p) {
		fired = true
		p.NitrousFull = false
		p.Stats.NitrousFired++
		bonus += g.Cfg.NitrousBonus
		if g.Cfg.ThirstyNitrous {
			p.Fuel = max(0, p.Fuel-1)
		}
	}

	var move int
	if p.Jammed {
		p.Jammed = false
		move = 2
		g.logf("%s is stuck in traffic: moves 2", p.Name)
	} else {
		r := g.d6()
		if drift {
			r = max(r, g.d6())
		}
		for _, u := range []Upgrade{Turbo, Tires} {
			if p.Up[u] {
				bonus++
			}
		}
		if p.Flat {
			bonus -= 2
		}
		move = max(1, r+bonus+p.NextMod)
		g.logf("%s rolls %d → moves %d (fuel %d, heat %d, rep %d)", p.Name, r, move, p.Fuel, p.Heat, p.Rep)
		if fired {
			g.addHeat(p, 1)
		}
		if drift {
			g.addHeat(p, 1)
		}
		if r >= g.Cfg.HeatRollMin {
			g.addHeat(p, 1)
		}
		if r == 1 {
			g.gainRep(p, 1)
		}
	}
	p.NextMod = 0
	g.moveForward(p, move, true)
	if p.Finished || p.Route == "" {
		return
	}

	// After-move cards move you on before you resolve anything.
	if s.useAfterMove(g, p) && (p.Finished || p.Route == "") {
		return
	}

	// Step 3 and 4.
	g.resolve(p)
	g.duel(p)
}

// moveForward moves p up to n spaces. isRoll marks normal roll movement,
// the only kind that triggers speed traps and optional gas-station stops.
func (g *Game) moveForward(p *Player, n int, isRoll bool) {
	if p.Finished {
		return
	}
	if p.Route == "" {
		p.Route = p.Strat.chooseRoute(g, p)
		p.Pos = 0
		if p.Leg == 1 {
			p.Stats.Leg1 = p.Route
		} else {
			p.Stats.Leg2 = p.Route
		}
	}
	r := g.Board.Routes[p.Route]
	for i := 1; i <= n; i++ {
		if p.Pos == r.Len() {
			if p.Leg == 1 {
				g.arriveShibuya(p)
			} else {
				g.finish(p)
			}
			return
		}
		p.Pos++
		if i == n {
			break
		}
		sp := r.At(p.Pos)
		if isRoll && sp.Kind == Trap {
			if n >= g.Cfg.TrapSpeed && !p.Up[Scanner] && !p.tipActive {
				p.Stats.TrapPassCatches++
				g.logf("%s blows past a speed trap at %s%d and gets caught", p.Name, p.Route, p.Pos)
				break
			}
			if g.Cfg.TrapPassHeat > 0 {
				g.addHeat(p, g.Cfg.TrapPassHeat)
			}
		}
		if isRoll && g.Cfg.OilOnPass && g.Oil[g.loc(p)] {
			break // spin out here; resolve() handles the oil
		}
		if isRoll && sp.Kind == Gas && p.Strat.stopForGas(g, p) {
			p.Stats.FuelStops++
			g.logf("%s pulls into the gas station at %s%d early", p.Name, p.Route, p.Pos)
			break
		}
	}
	for p.Pos > 0 && r.At(p.Pos).Narrow && g.occupiedByOther(p) {
		p.Pos--
	}
}

func (g *Game) moveBack(p *Player, n int) {
	if p.Route == "" || p.Finished {
		return
	}
	lost := min(n, p.Pos)
	p.Pos -= n
	p.Stats.SpacesBack += lost
	if p.Pos <= 0 {
		p.Route, p.Pos = "", 0
	}
}

func (g *Game) arriveShibuya(p *Player) {
	p.Leg, p.Route, p.Pos = 2, "", 0
	if g.Cfg.ShibuyaRefuel {
		p.Fuel = g.Cfg.FuelMax
	}
	g.Shibuya++
	if g.Shibuya == 1 {
		g.gainRep(p, g.Cfg.ShibuyaBonus)
	} else {
		g.gainRep(p, g.Cfg.ShibuyaCatchUp*(g.Shibuya-1))
	}
	p.Stats.ShibuyaTurn = p.Stats.Turns
	p.Stats.ShibuyaRank = g.Shibuya
	p.Stats.ShibuyaRep = p.Rep
	p.Stats.ShibuyaUpgrades = p.Upgrades()
	g.logf("%s reaches Shibuya Crossing (#%d) on turn %d", p.Name, g.Shibuya, p.Stats.Turns)
}

func (g *Game) finish(p *Player) {
	p.Finished = true
	g.Winner = p.ID
	g.logf("🏁 %s reaches Tokyo Tower and wins on round %d", p.Name, g.Round)
}

func (g *Game) resolve(p *Player) {
	if p.Finished || p.Route == "" {
		return
	}
	l := g.loc(p)
	if g.Oil[l] {
		delete(g.Oil, l)
		if !p.Up[Tires] {
			p.Stats.OilSpinouts++
			g.logf("%s spins out on an oil slick", p.Name)
			g.moveBack(p, 3)
			return
		}
	}
	sp := g.Board.Routes[p.Route].At(p.Pos)
	switch sp.Kind {
	case Rep3:
		g.gainRep(p, g.Cfg.RepSpace)
	case Rep5:
		g.gainRep(p, g.Cfg.Rep5Space)
	case BoostSpace:
		g.drawBoost(p)
	case ObstacleSpace:
		g.drawObstacle(p)
	case Gas:
		g.pitStop(p)
	case Trap:
		if !p.tipActive {
			g.pullOver(p)
		}
	case Landmark:
		g.landmark(p, sp.Name)
	}
}

func (g *Game) pitStop(p *Player) {
	p.Fuel = g.Cfg.FuelMax
	p.Heat = max(0, p.Heat-g.Cfg.GasHeatDrop)
	if p.Up[Nitrous] {
		p.NitrousFull = true
	}
	if p.Flat {
		p.Flat = false
		g.Obstacles.Discard(FlatTire)
	}
}

func (g *Game) outOfGas(p *Player) {
	p.Stats.OutOfGas++
	if g.Cfg.FuelDelivery && p.Rep >= g.Cfg.FuelDeliveryCost && p.Strat.wantsDelivery(g, p) {
		p.Rep -= g.Cfg.FuelDeliveryCost
		p.Fuel = g.Cfg.FuelMax
		g.park(p, ParkDelivery)
		g.logf("%s is out of gas and pays for a delivery", p.Name)
		return
	}
	lost := 0
	if p.Route != "" {
		r := g.Board.Routes[p.Route]
		pos := p.Pos
		for pos >= 1 && r.At(pos).Kind != Gas {
			pos--
		}
		lost = p.Pos - pos
		if pos == 0 {
			p.Route = ""
		}
		p.Pos = pos
	}
	p.Stats.OutOfGasSpaces += lost
	p.Stats.SpacesBack += lost
	g.logf("%s runs out of gas and pushes back %d spaces", p.Name, lost)
	g.pitStop(p)
}

func (g *Game) pullOver(p *Player) {
	p.Stats.PulledOver++
	fine := max(g.Cfg.MinFine, p.Heat)
	if fine == 0 {
		g.logf("%s gets off with a warning", p.Name)
		return
	}
	if p.has(BribeACop) && (fine >= 2 || p.Rep < fine) {
		g.play(p, BribeACop)
		p.Stats.Bribes++
		p.Heat = 1
		g.logf("%s bribes their way out of a ticket", p.Name)
		return
	}
	pay := min(fine, p.Rep)
	p.Rep -= pay
	p.Stats.Fines += pay
	g.logf("%s is pulled over: fine %d, paid %d", p.Name, fine, pay)
	if pay < fine {
		g.park(p, ParkFine)
	}
	if g.Cfg.PulledHeatDrop > 0 {
		p.Heat = max(0, p.Heat-g.Cfg.PulledHeatDrop)
	} else {
		p.Heat = 0
	}
}

func (g *Game) addHeat(p *Player, n int) {
	p.Heat += n
	if p.Heat < g.Cfg.HeatMax {
		return
	}
	if p.has(BribeACop) {
		g.play(p, BribeACop)
		p.Stats.Bribes++
		p.Heat = 1
		return
	}
	p.Stats.Busted++
	p.Heat = g.Cfg.BustedHeat
	g.park(p, ParkBusted)
	g.logf("%s hits the red line: BUSTED", p.Name)
}

func (g *Game) park(p *Player, c ParkCause) {
	if !p.Parked {
		p.Parked = true
		p.Stats.ParkCauses[c]++
	}
}

func (g *Game) gainRep(p *Player, n int) {
	p.Rep += n
	p.Stats.RepEarned += n
}

func (g *Game) play(p *Player, b Boost) {
	for i, c := range p.Hand {
		if c == b {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			g.Boosts.Discard(b)
			return
		}
	}
	panic("playing a card not in hand: " + b.String())
}

func (g *Game) drawBoost(p *Player) {
	c, ok := g.Boosts.Draw()
	if !ok {
		return
	}
	p.Stats.Boosts[c]++
	g.logf("%s draws Boost: %s", p.Name, c)
	if c == SponsorDeal { // no reason to hold it
		g.Boosts.Discard(c)
		g.gainRep(p, 4)
		return
	}
	p.Hand = append(p.Hand, c)
	for len(p.Hand) > g.Cfg.HandLimit {
		g.play(p, p.Strat.worstCard(p.Hand))
	}
}

func (g *Game) drawObstacle(p *Player) {
	c, ok := g.Obstacles.Draw()
	if !ok {
		return
	}
	p.Stats.Obstacles[c]++
	g.logf("%s draws Obstacle: %s", p.Name, c)
	if p.has(MechanicOnCall) && p.Strat.cancelObstacle(g, p, c) {
		g.play(p, MechanicOnCall)
		g.Obstacles.Discard(c)
		return
	}
	if c != FlatTire || p.Flat {
		g.Obstacles.Discard(c)
	}
	switch c {
	case FlatTire:
		p.Flat = true
	case FuelLeak:
		p.Fuel = max(0, p.Fuel-2)
	case PoliceCheckpoint:
		if p.Heat >= 1 {
			g.pullOver(p)
		} else {
			g.gainRep(p, 1)
		}
	case TrafficJam:
		p.Jammed = true
	case Construction:
		g.moveBack(p, 2)
	case WrongTurn:
		g.moveBack(p, 3)
	case EngineTrouble:
		p.NextMod -= 2
	case Rain:
		if !p.Up[Tires] {
			p.NextMod--
		}
	case RivalCrew:
		if q := g.closestBehind(p); q != nil {
			amt := min(2, p.Rep)
			p.Rep -= amt
			q.Rep += amt
		}
	case TowZone:
		g.park(p, ParkTow)
	case CaughtOnCamera:
		g.gainRep(p, 2)
		g.addHeat(p, 2)
	}
}

func (g *Game) landmark(p *Player, name string) {
	switch name {
	case "Daikoku PA":
		p.Strat.buyUpgrades(g, p, g.Cfg.DaikokuDiscount)
	case "Rainbow Bridge":
		if p.Strat.UseBridge && p.Rep >= g.Cfg.BridgeCost {
			p.Rep -= g.Cfg.BridgeCost
			g.logf("%s pays the Rainbow Bridge toll", p.Name)
			g.moveForward(p, g.Cfg.BridgeJump, false)
			g.resolve(p)
		}
	case "Ginza":
		if p.Rep >= 1 && len(p.Hand) < g.Cfg.HandLimit {
			p.Rep--
			n, sponsors := len(p.Hand), p.Stats.Boosts[SponsorDeal]
			g.drawBoost(p)
			g.drawBoost(p)
			switch {
			case p.Stats.Boosts[SponsorDeal] > sponsors: // a Sponsor Deal counts as the card you keep
				for len(p.Hand) > n {
					g.play(p, p.Hand[len(p.Hand)-1])
				}
			case len(p.Hand) == n+2: // keep the better of the two new cards
				g.play(p, p.Strat.worstCard(p.Hand[n:]))
			}
		}
	case "Akihabara":
		if q := g.leader(p); q != nil && q.Route != "" {
			g.logf("%s hacks the lights on %s", p.Name, q.Name)
			g.moveBack(q, 2)
		}
	case "Tsukiji":
		g.gainRep(p, 3)
		p.NextMod--
	case "Ameyoko":
		p.Heat = 0
	case "Harajuku":
		g.gainRep(p, p.Upgrades())
	case "Roppongi":
		g.gainRep(p, g.carsBehind(p))
		g.addHeat(p, 1)
	}
}

func (g *Game) duel(p *Player) {
	if p.Route == "" || p.Finished || g.Board.Routes[p.Route].At(p.Pos).Narrow {
		return
	}
	var opp *Player
	for _, q := range g.Players {
		if q != p && !q.Finished && g.loc(q) == g.loc(p) && (opp == nil || q.Rep > opp.Rep) {
			opp = q
		}
	}
	if opp == nil {
		return
	}
	a, b := g.d6(), g.d6()
	for a == b {
		a, b = g.d6(), g.d6()
	}
	win, lose := p, opp
	if b > a {
		win, lose = opp, p
	}
	amt := min(g.Cfg.DuelStake, lose.Rep)
	lose.Rep -= amt
	win.Rep += amt
	win.Stats.DuelsWon++
	lose.Stats.DuelsLost++
	g.logf("Street duel! %s beats %s", win.Name, lose.Name)
}

// ── position helpers ────────────────────────────────────────────────

func (g *Game) loc(p *Player) Loc { return Loc{p.Leg, p.Route, p.Pos} }

// Remaining estimates spaces left to the finish, so cars on different
// routes can be ranked.
func (g *Game) Remaining(p *Player) float64 {
	if p.Finished {
		return 0
	}
	leg2 := g.Board.AvgLeg2() + 1
	if p.Leg == 2 {
		if p.Route == "" {
			return leg2
		}
		return float64(g.Board.Routes[p.Route].Len() - p.Pos + 1)
	}
	if p.Route == "" {
		return 28 + 1 + leg2
	}
	return float64(g.Board.Routes[p.Route].Len()-p.Pos+1) + leg2
}

func (g *Game) trackStall(p *Player, before float64) {
	if !p.Finished && g.Remaining(p) >= before-0.01 {
		p.Stats.stall++
		p.Stats.MaxStall = max(p.Stats.MaxStall, p.Stats.stall)
	} else {
		p.Stats.stall = 0
	}
}

func (g *Game) others(p *Player) []*Player {
	var o []*Player
	for _, q := range g.Players {
		if q != p && !q.Finished {
			o = append(o, q)
		}
	}
	return o
}

func (g *Game) leader(p *Player) *Player {
	var best *Player
	for _, q := range g.others(p) {
		if best == nil || g.Remaining(q) < g.Remaining(best) {
			best = q
		}
	}
	return best
}

func (g *Game) closestBehind(p *Player) *Player {
	var best *Player
	me := g.Remaining(p)
	for _, q := range g.others(p) {
		if r := g.Remaining(q); r > me && (best == nil || r < g.Remaining(best)) {
			best = q
		}
	}
	return best
}

func (g *Game) carsBehind(p *Player) int {
	n := 0
	for _, q := range g.others(p) {
		if g.Remaining(q) > g.Remaining(p) {
			n++
		}
	}
	return n
}

func (g *Game) occupiedByOther(p *Player) bool {
	for _, q := range g.others(p) {
		if g.loc(q) == g.loc(p) {
			return true
		}
	}
	return false
}

func (g *Game) slipstreaming(p *Player) bool {
	if p.Route == "" {
		return false
	}
	for _, q := range g.others(p) {
		if q.Leg == p.Leg && q.Route == p.Route && q.Pos == p.Pos+1 {
			return true
		}
	}
	return false
}

// nearestAhead finds the closest car ahead on p's route and its distance.
func (g *Game) nearestAhead(p *Player) (*Player, int) {
	if p.Route == "" {
		return nil, 0
	}
	var best *Player
	for _, q := range g.others(p) {
		if q.Leg == p.Leg && q.Route == p.Route && q.Pos > p.Pos && (best == nil || q.Pos < best.Pos) {
			best = q
		}
	}
	if best == nil {
		return nil, 0
	}
	return best, best.Pos - p.Pos
}

// trapAhead reports a speed trap within n spaces ahead (or on the next
// route if p is waiting on a hub).
func (g *Game) trapAhead(p *Player, n int) bool {
	if p.Route == "" {
		return true
	}
	r := g.Board.Routes[p.Route]
	for i := p.Pos + 1; i <= min(r.Len(), p.Pos+n); i++ {
		if r.At(i).Kind == Trap {
			return true
		}
	}
	return false
}

// distToFuel is how far p is from the next place to refuel ahead
// (a gas station, Shibuya, or the finish).
func (g *Game) distToFuel(p *Player, from int) int {
	if p.Route == "" {
		return 8
	}
	r := g.Board.Routes[p.Route]
	for i := from + 1; i <= r.Len(); i++ {
		if r.At(i).Kind == Gas {
			return i - from
		}
	}
	return r.Len() + 1 - from
}
