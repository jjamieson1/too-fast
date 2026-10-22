// Command sim plays thousands of Too Fast games with bot strategies and
// prints a balance report (markdown).
//
//	go run ./cmd/sim                    # full report, 20,000 games per experiment
//	go run ./cmd/sim -games 2000        # quicker
//	go run ./cmd/sim -exp routes        # one experiment
//	go run ./cmd/sim -trace -seed 7     # play-by-play of a single game
package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/jjamieson1/too-fast/game"
)

var (
	games   = flag.Int("games", 20000, "games per experiment")
	seed    = flag.Uint64("seed", 1, "random seed")
	exp     = flag.String("exp", "all", "experiment: all, strategies, routes, length, cards, variants, fixes")
	trace   = flag.Bool("trace", false, "print a play-by-play of one game and exit")
	rules   = flag.String("rules", "v0.3", "rule set for experiments 1–4: v0.3 (current) or v0.2")
	secTurn = flag.Float64("sec-per-turn", 30, "assumed real seconds per player turn, for time estimates")
)

// ── bots ────────────────────────────────────────────────────────────

type bot struct {
	name string
	make func(r *rand.Rand) *game.Strategy
}

func fixed(s game.Strategy) func(*rand.Rand) *game.Strategy {
	return func(*rand.Rand) *game.Strategy { c := s; return &c }
}

var (
	N, S, T, R = game.Nitrous, game.Scanner, game.Turbo, game.Tires
)

func roster() []bot {
	return []bot{
		{"Rusher", fixed(game.Strategy{Name: "Rusher", Leg1: "A", Leg2: "D", BuyOrder: []game.Upgrade{T, R, N},
			FuelStop: -1, NitrousMode: game.NitrousAlways, UseBridge: true})},
		{"Scanner Runner", fixed(game.Strategy{Name: "Scanner Runner", Leg1: "A", Leg2: "D", BuyOrder: []game.Upgrade{S, T, R, N},
			Reserve: 1, FuelStop: -1, NitrousMode: game.NitrousSmart, UseBridge: true, FixFlatDist: 10})},
		{"Alley Builder", fixed(game.Strategy{Name: "Alley Builder", Leg1: "C", Leg2: "", BuyOrder: []game.Upgrade{S, T, R, N},
			Reserve: 1, FuelStop: -1, NitrousMode: game.NitrousSmart, FixFlatDist: 10})},
		{"Balanced", fixed(game.Strategy{Name: "Balanced", Leg1: "B", Leg2: "E", BuyOrder: []game.Upgrade{T, R, N},
			Reserve: 1, FuelStop: -1, NitrousMode: game.NitrousSmart, FixFlatDist: 10})},
		{"Showboat", fixed(game.Strategy{Name: "Showboat", Leg1: "C", Leg2: "E", BuyOrder: []game.Upgrade{T, R, N, S},
			Reserve: 1, FuelStop: -1, NitrousMode: game.NitrousSmart, FixFlatDist: 10})},
		{"Tourist", fixed(game.Strategy{Name: "Tourist", Leg1: "B", Leg2: "E", FuelStop: -1})},
		{"Careless", fixed(game.Strategy{Name: "Careless", Leg1: "B", Leg2: "E", BuyOrder: []game.Upgrade{T, R},
			FuelStop: 0, NitrousMode: game.NitrousSmart})},
		{"Random", func(r *rand.Rand) *game.Strategy {
			order := []game.Upgrade{N, S, T, R}
			r.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })
			return &game.Strategy{Name: "Random", Leg2: "?", BuyOrder: order[:r.IntN(5)],
				FuelStop: r.IntN(3), NitrousMode: game.NitrousAlways, UseBridge: r.IntN(2) == 0}
		}},
	}
}

// standard is the bot used when we only want to vary the route.
func standard(name, leg1, leg2 string, order []game.Upgrade) *game.Strategy {
	return &game.Strategy{Name: name, Leg1: leg1, Leg2: leg2, BuyOrder: order, Reserve: 1,
		FuelStop: -1, NitrousMode: game.NitrousSmart, UseBridge: true, FixFlatDist: 10}
}

// ── running games ───────────────────────────────────────────────────

// tournament plays n games of `players` seats, each filled with distinct
// random bots from the roster in random seat order.
func tournament(cfg game.Config, bots []bot, players, n int, s uint64) []*game.Game {
	r := rand.New(rand.NewPCG(s, 77))
	out := make([]*game.Game, 0, n)
	for i := 0; i < n; i++ {
		idx := r.Perm(len(bots))[:players]
		var strats []*game.Strategy
		for _, j := range idx {
			strats = append(strats, bots[j].make(r))
		}
		g := game.New(cfg, strats, s*1_000_003+uint64(i))
		g.Run()
		out = append(out, g)
	}
	return out
}

// ── tallies ─────────────────────────────────────────────────────────

type tally struct {
	apps, wins                                  int
	shibTurn, shibUp, shibRep, firstShib        float64
	oog, oogSpaces, parked, pulled, busted      float64
	fines, stall3, broke, turns, duels, spinout float64
	upTurn                                      [4]float64
	upCount                                     [4]float64
}

func (t *tally) add(g *game.Game, p *game.Player) {
	t.apps++
	if g.Winner == p.ID {
		t.wins++
	}
	st := p.Stats
	if st.ShibuyaTurn > 0 {
		t.shibTurn += float64(st.ShibuyaTurn)
		t.shibUp += float64(st.ShibuyaUpgrades)
		t.shibRep += float64(st.ShibuyaRep)
		if st.ShibuyaRank == 1 {
			t.firstShib++
		}
	}
	t.oog += float64(st.OutOfGas)
	t.oogSpaces += float64(st.OutOfGasSpaces)
	t.parked += float64(st.ParkedTurns)
	t.pulled += float64(st.PulledOver)
	t.busted += float64(st.Busted)
	t.fines += float64(st.Fines)
	t.turns += float64(st.Turns)
	t.broke += float64(st.BrokeTurns)
	t.duels += float64(st.DuelsWon + st.DuelsLost)
	t.spinout += float64(st.OilSpinouts)
	if st.MaxStall >= 3 {
		t.stall3++
	}
	for u, turn := range st.UpgradeTurn {
		if turn > 0 {
			t.upTurn[u] += float64(turn)
			t.upCount[u]++
		}
	}
}

func (t *tally) per(x float64) float64 { return x / float64(max(1, t.apps)) }
func (t *tally) winPct() float64       { return 100 * float64(t.wins) / float64(max(1, t.apps)) }

func byKey(gs []*game.Game, key func(*game.Player) string) (map[string]*tally, []string) {
	m := map[string]*tally{}
	for _, g := range gs {
		for _, p := range g.Players {
			k := key(p)
			if m[k] == nil {
				m[k] = &tally{}
			}
			m[k].add(g, p)
		}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return m[keys[i]].winPct() > m[keys[j]].winPct() })
	return m, keys
}

// ── output helpers ──────────────────────────────────────────────────

func h(level int, s string) { fmt.Printf("\n%s %s\n\n", strings.Repeat("#", level), s) }

func table(head []string, rows [][]string) {
	fmt.Println("| " + strings.Join(head, " | ") + " |")
	sep := make([]string, len(head))
	for i := range sep {
		sep[i] = "---"
	}
	fmt.Println("|" + strings.Join(sep, "|") + "|")
	for _, r := range rows {
		fmt.Println("| " + strings.Join(r, " | ") + " |")
	}
}

func f1(x float64) string  { return fmt.Sprintf("%.1f", x) }
func f2(x float64) string  { return fmt.Sprintf("%.2f", x) }
func pct(x float64) string { return fmt.Sprintf("%.1f%%", x) }

func pctile(xs []int, p float64) int {
	s := slices.Clone(xs)
	slices.Sort(s)
	return s[int(p*float64(len(s)-1))]
}

func rounds(gs []*game.Game) []int {
	var r []int
	for _, g := range gs {
		r = append(r, g.Round)
	}
	return r
}

func timeouts(gs []*game.Game) int {
	n := 0
	for _, g := range gs {
		if g.Winner < 0 {
			n++
		}
	}
	return n
}

// ── experiments ─────────────────────────────────────────────────────

func expStrategies(cfg game.Config) {
	h(2, "1. Strategy tournament (4 players, random bots and seats)")
	gs := tournament(cfg, roster(), 4, *games, *seed)
	m, keys := byKey(gs, func(p *game.Player) string { return p.Strat.Name })
	fmt.Println("A fair game gives every strategy about **25%**. Per-player averages per game.")
	fmt.Println()
	var rows [][]string
	for _, k := range keys {
		t := m[k]
		rows = append(rows, []string{k, pct(t.winPct()), f1(t.per(t.shibTurn)), f1(t.per(t.shibUp)),
			f2(t.per(t.oog)), f2(t.per(t.parked)), f2(t.per(t.pulled)), f2(t.per(t.busted)), f1(t.per(t.fines))})
	}
	table([]string{"Strategy", "Win %", "Turn at Shibuya", "Upgrades at Shibuya", "Out of gas", "Turns parked", "Pulled over", "Busted", "Rep paid in fines"}, rows)

	h(3, "Upgrade timing")
	fmt.Println("Average turn each upgrade is bought, and how often it gets bought at all (across all bots that want it).")
	fmt.Println()
	all := &tally{}
	for _, g := range gs {
		for _, p := range g.Players {
			if p.Strat.Name != "Tourist" {
				all.add(g, p)
			}
		}
	}
	rows = nil
	for u, name := range []string{"Nitrous", "Police Scanner", "Turbo", "Racing Tires"} {
		rows = append(rows, []string{name, f1(all.upTurn[u] / max(1, all.upCount[u])), pct(100 * all.upCount[u] / float64(all.apps))})
	}
	table([]string{"Upgrade", "Avg turn bought", "Bought in % of games"}, rows)

	h(3, "Seat order")
	seat, skeys := byKey(gs, func(p *game.Player) string { return fmt.Sprintf("Seat %d", p.ID+1) })
	sort.Strings(skeys)
	rows = nil
	for _, k := range skeys {
		rows = append(rows, []string{k, pct(seat[k].winPct())})
	}
	table([]string{"Seat", "Win %"}, rows)

	h(3, "Runaway leader check")
	upW, upN := map[int]int{}, map[int]int{}
	rankW, rankN := map[int]int{}, map[int]int{}
	for _, g := range gs {
		for _, p := range g.Players {
			u, rk := p.Stats.ShibuyaUpgrades, p.Stats.ShibuyaRank
			upN[u]++
			if rk > 0 {
				rankN[rk]++
			}
			if g.Winner == p.ID {
				upW[u]++
				if rk > 0 {
					rankW[rk]++
				}
			}
		}
	}
	rows = nil
	for rk := 1; rk <= 4; rk++ {
		rows = append(rows, []string{fmt.Sprintf("#%d into Shibuya", rk), pct(100 * float64(rankW[rk]) / float64(max(1, rankN[rk])))})
	}
	table([]string{"Arrival order at Shibuya", "Win %"}, rows)
	fmt.Println()
	rows = nil
	for u := 0; u <= 4; u++ {
		if upN[u] > 0 {
			rows = append(rows, []string{fmt.Sprint(u), pct(100 * float64(upN[u]) / float64(4**games)), pct(100 * float64(upW[u]) / float64(upN[u]))})
		}
	}
	table([]string{"Upgrades owned at Shibuya", "Share of players", "Win %"}, rows)

	h(3, "Where players get stuck")
	var stuckRows [][]string
	tot := &tally{}
	var causes [game.NumParkCauses]int
	stallHist := map[int]int{}
	worst := 0
	for _, g := range gs {
		for _, p := range g.Players {
			tot.add(g, p)
			for c, n := range p.Stats.ParkCauses {
				causes[c] += n
			}
			stallHist[min(p.Stats.MaxStall, 6)]++
			worst = max(worst, p.Stats.MaxStall)
		}
	}
	pg := float64(tot.apps)
	stuckRows = append(stuckRows,
		[]string{"Turns parked (lost)", f2(tot.parked / pg)},
		[]string{"Out-of-gas events", f2(tot.oog / pg)},
		[]string{"Spaces pushed back when out of gas", f2(tot.oogSpaces / pg)},
		[]string{"Times pulled over", f2(tot.pulled / pg)},
		[]string{"Times Busted", f2(tot.busted / pg)},
		[]string{"Oil slick spin-outs", f2(tot.spinout / pg)},
		[]string{"Street duels", f2(tot.duels / pg)},
		[]string{"Turns started broke (0 Rep)", pct(100 * tot.broke / tot.turns)},
		[]string{"Players with a 3+ turn stall", pct(100 * tot.stall3 / pg)},
		[]string{"Longest stall seen", fmt.Sprintf("%d turns", worst)},
		[]string{"Games that never finished", fmt.Sprint(timeouts(gs))},
	)
	table([]string{"Per player, per game", "Average"}, stuckRows)
	fmt.Println()
	fmt.Println("A **stall** is a run of turns in a row where the car made no forward progress (parked, pushed back, out of gas).")
	fmt.Println()
	rows = nil
	for c := game.ParkCause(0); c < game.NumParkCauses; c++ {
		rows = append(rows, []string{c.String(), f2(float64(causes[c]) / float64(len(gs)))})
	}
	table([]string{"Why turns are lost", "Per game (all players)"}, rows)
	fmt.Println()
	rows = nil
	for s := 0; s <= 6; s++ {
		label := fmt.Sprint(s)
		if s == 6 {
			label = "6+"
		}
		rows = append(rows, []string{label, pct(100 * float64(stallHist[s]) / pg)})
	}
	table([]string{"Longest stall in a game (turns)", "Share of players"}, rows)
}

func expRoutes(cfg game.Config) {
	h(2, "2. Route balance (identical bots, only the route differs)")
	combos := [][2]string{{"A", "D"}, {"A", "E"}, {"B", "D"}, {"B", "E"}, {"C", "D"}, {"C", "E"}}
	for _, v := range []struct {
		title string
		order []game.Upgrade
	}{
		{"Bots buy Turbo → Tires → Scanner → Nitrous", []game.Upgrade{T, R, S, N}},
		{"Bots buy Scanner first → Turbo → Tires → Nitrous", []game.Upgrade{S, T, R, N}},
	} {
		h(3, v.title)
		var bots []bot
		for _, c := range combos {
			name := c[0] + "→" + c[1]
			leg1, leg2 := c[0], c[1]
			order := v.order
			bots = append(bots, bot{name, func(*rand.Rand) *game.Strategy { return standard(name, leg1, leg2, order) }})
		}
		gs := tournament(cfg, bots, 4, *games, *seed+11)
		m, keys := byKey(gs, func(p *game.Player) string { return p.Strat.Name })
		fmt.Println("Fair = about 25%.")
		fmt.Println()
		var rows [][]string
		for _, k := range keys {
			t := m[k]
			rows = append(rows, []string{k, pct(t.winPct()), f1(t.per(t.shibTurn)), f1(t.per(t.shibUp)), f1(t.per(t.shibRep)),
				f2(t.per(t.pulled)), f2(t.per(t.oog)), f1(t.per(t.turns))})
		}
		table([]string{"Route", "Win %", "Turn at Shibuya", "Upgrades at Shibuya", "Rep at Shibuya", "Pulled over", "Out of gas", "Turns played"}, rows)
	}
}

func expLength(cfg game.Config) {
	h(2, "3. Game length")
	var rows [][]string
	for _, n := range []int{2, 3, 4} {
		gs := tournament(cfg, roster(), n, *games, *seed+uint64(n))
		r := rounds(gs)
		med := pctile(r, 0.5)
		rows = append(rows, []string{fmt.Sprint(n), fmt.Sprint(pctile(r, 0.1)), fmt.Sprint(med), fmt.Sprint(pctile(r, 0.9)), fmt.Sprint(pctile(r, 1)),
			fmt.Sprintf("~%.0f min", float64(med*n)**secTurn/60)})
	}
	table([]string{"Players", "Rounds (10th pct)", "Median", "90th pct", "Longest", "Median play time"}, rows)
	fmt.Printf("\nPlay time assumes about %.0f seconds per turn. Change it with `-sec-per-turn`.\n", *secTurn)
}

func expCards(cfg game.Config) {
	h(2, "4. Card impact")
	fmt.Println("Win rate of players who drew each card at least once (4 players, so the baseline is 25%).")
	fmt.Println("Back Alley drivers draw more cards, so treat small differences with care.")
	fmt.Println()
	gs := tournament(cfg, roster(), 4, *games, *seed+5)
	type wn struct{ w, n, drawn int }
	ob := map[game.Obstacle]*wn{}
	bo := map[game.Boost]*wn{}
	for _, g := range gs {
		for _, p := range g.Players {
			for c, k := range p.Stats.Obstacles {
				o := game.Obstacle(c)
				if ob[o] == nil {
					ob[o] = &wn{}
				}
				ob[o].drawn += k
				if k > 0 {
					ob[o].n++
					if g.Winner == p.ID {
						ob[o].w++
					}
				}
			}
			for c, k := range p.Stats.Boosts {
				b := game.Boost(c)
				if bo[b] == nil {
					bo[b] = &wn{}
				}
				bo[b].drawn += k
				if k > 0 {
					bo[b].n++
					if g.Winner == p.ID {
						bo[b].w++
					}
				}
			}
		}
	}
	var rows [][]string
	for o, v := range ob {
		rows = append(rows, []string{o.String(), f2(float64(v.drawn) / float64(len(gs))), pct(100 * float64(v.w) / float64(max(1, v.n)))})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i][2] < rows[j][2] })
	table([]string{"Obstacle", "Drawn per game", "Win % after drawing"}, rows)
	fmt.Println()
	rows = nil
	for b, v := range bo {
		rows = append(rows, []string{b.String(), f2(float64(v.drawn) / float64(len(gs))), pct(100 * float64(v.w) / float64(max(1, v.n)))})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i][2] > rows[j][2] })
	table([]string{"Boost", "Drawn per game", "Win % after drawing"}, rows)
}

func expVariants() {
	h(2, "5. Rule variants")
	fmt.Println("Same 4-player tournament, one rule changed at a time.")
	fmt.Println()
	type variant struct {
		name string
		mod  func(*game.Config)
	}
	vs := []variant{
		{"Current rules (v0.3)", func(*game.Config) {}},
		{"Fuel delivery variant", func(c *game.Config) { c.FuelDelivery = true }},
		{"Thirsty nitrous variant", func(c *game.Config) { c.ThirstyNitrous = true }},
		{"Fuel tank 4", func(c *game.Config) { c.FuelMax = 4 }},
		{"Fuel tank 6", func(c *game.Config) { c.FuelMax = 6 }},
		{"No refuel at Shibuya", func(c *game.Config) { c.ShibuyaRefuel = false }},
		{"Speed traps catch at 7+ (not 6+)", func(c *game.Config) { c.TrapSpeed = 7 }},
		{"Start with 6 Rep", func(c *game.Config) { c.StartRep = 6 }},
	}
	var rows [][]string
	for _, v := range vs {
		cfg := game.DefaultConfig()
		v.mod(&cfg)
		gs := tournament(cfg, roster(), 4, *games, *seed+99)
		m, _ := byKey(gs, func(p *game.Player) string { return p.Strat.Name })
		all := &tally{}
		for _, t := range m {
			all.apps += t.apps
			all.oog += t.oog
			all.parked += t.parked
		}
		best, worst := "", ""
		for k, t := range m {
			if best == "" || t.winPct() > m[best].winPct() {
				best = k
			}
			if worst == "" || t.winPct() < m[worst].winPct() {
				worst = k
			}
		}
		rows = append(rows, []string{v.name, fmt.Sprint(pctile(rounds(gs), 0.5)), f2(all.per(all.oog)), f2(all.per(all.parked)),
			fmt.Sprintf("%s %s", best, pct(m[best].winPct())), fmt.Sprintf("%s %s", worst, pct(m[worst].winPct())),
			pct(m["Careless"].winPct()), pct(m["Tourist"].winPct())})
	}
	table([]string{"Rules", "Median rounds", "Out of gas / player", "Turns parked / player", "Best bot", "Worst bot", "Careless win %", "Tourist win %"}, rows)
}

func main() {
	flag.Parse()
	cfg := game.DefaultConfig()
	switch *rules {
	case "v0.3", "proposal-f":
	case "v0.2":
		cfg = game.V02Config()
	default:
		fmt.Fprintln(os.Stderr, "unknown -rules", *rules)
		os.Exit(2)
	}
	if *trace {
		r := rand.New(rand.NewPCG(*seed, 77))
		bots := roster()
		var strats []*game.Strategy
		for _, j := range r.Perm(len(bots))[:4] {
			strats = append(strats, bots[j].make(r))
		}
		g := game.New(cfg, strats, *seed)
		g.Log = func(s string) { fmt.Println(s) }
		g.Run()
		return
	}
	fmt.Println("# Too Fast — Simulation Report")
	fmt.Printf("\n%d games per experiment · seed %d · rules %s (experiments 1–4)\n", *games, *seed, *rules)
	run := func(name string, f func()) {
		if *exp == "all" || *exp == name {
			f()
		}
	}
	run("strategies", func() { expStrategies(cfg) })
	run("routes", func() { expRoutes(cfg) })
	run("length", func() { expLength(cfg) })
	run("cards", func() { expCards(cfg) })
	run("variants", expVariants)
	run("fixes", expFixes)
	if !slices.Contains([]string{"all", "strategies", "routes", "length", "cards", "variants", "fixes"}, *exp) {
		fmt.Fprintln(os.Stderr, "unknown -exp", *exp)
		os.Exit(2)
	}
}

// ── fix candidates ──────────────────────────────────────────────────

type fix struct {
	name string
	mod  func(*game.Config)
}

func alleyRep3(c *game.Config) {
	if seq, ok := c.RouteOverride["C"]; ok {
		c.RouteOverride["C"] = strings.ReplaceAll(seq, "R5", "R")
	}
}

var fixes = []fix{
	{"v0.2 (before any fixes)", func(*game.Config) {}},
	{"Alley +5 spaces become +3", alleyRep3},
	{"Speed traps catch at 7+", func(c *game.Config) { c.TrapSpeed = 7 }},
	{"Heat 0 = warning, no fine", func(c *game.Config) { c.MinFine = 0 }},
	{"Later seats start +1 Rep each", func(c *game.Config) { c.SeatRepBonus = 1 }},
	{"Shortcut moves 6 max", func(c *game.Config) { c.ShortcutMax = 6 }},
	{"Passing a trap safely: +1 Heat", func(c *game.Config) { c.TrapPassHeat = 1 }},
	{"Pulled over: Heat −2 (not reset)", func(c *game.Config) { c.PulledHeatDrop = 2 }},
	{"Oil slicks trigger when driven over", func(c *game.Config) { c.OilOnPass = true }},
	{"Seats 3 and 4 start +1 Rep", func(c *game.Config) { c.SeatRep = []int{0, 0, 1, 1} }},
	{"No bonus for 1st into Shibuya", func(c *game.Config) { c.ShibuyaBonus = 0 }},
	{"Shibuya catch-up: +2 Rep per place behind 1st", func(c *game.Config) { c.ShibuyaBonus = 0; c.ShibuyaCatchUp = 2 }},
	{"Natural 5 or 6 adds Heat", func(c *game.Config) { c.HeatRollMin = 5 }},
	{"**Proposal A**: alley +3, traps 7+, warning, seat Rep, Shortcut 6", proposalA},
	{"**Proposal B**: A + passing a trap safely gives +1 Heat", func(c *game.Config) { proposalA(c); c.TrapPassHeat = 1 }},
	{"**Proposal C**: A, but seats 3–4 +1 Rep and Shibuya catch-up", proposalC},
	{"**Proposal E**: A, but seat Rep +0/+0/+1/+2", proposalE},
	{"**Proposal F**: E + no bonus for 1st into Shibuya", func(c *game.Config) { proposalE(c); c.ShibuyaBonus = 0 }},
}

func proposalE(c *game.Config) {
	proposalA(c)
	c.SeatRepBonus = 0
	c.SeatRep = []int{0, 0, 1, 2}
}

func proposalC(c *game.Config) {
	proposalA(c)
	c.SeatRepBonus = 0
	c.SeatRep = []int{0, 0, 1, 1}
	c.ShibuyaBonus = 0
	c.ShibuyaCatchUp = 2
}

// sensible bots: the real play styles, leaving out the deliberately weak ones.
var sensible = map[string]bool{"Rusher": true, "Scanner Runner": true, "Alley Builder": true, "Balanced": true, "Showboat": true}

func proposalA(c *game.Config) {
	alleyRep3(c)
	c.TrapSpeed = 7
	c.MinFine = 0
	c.SeatRepBonus = 1
	c.ShortcutMax = 6
}

func routeBots() []bot {
	var bots []bot
	for _, c := range [][2]string{{"A", "D"}, {"A", "E"}, {"B", "D"}, {"B", "E"}, {"C", "D"}, {"C", "E"}} {
		name, leg1, leg2 := c[0]+"→"+c[1], c[0], c[1]
		bots = append(bots, bot{name, func(*rand.Rand) *game.Strategy { return standard(name, leg1, leg2, []game.Upgrade{T, R, S, N}) }})
	}
	return bots
}

func spread(m map[string]*tally) (lo, hi string) {
	for k, t := range m {
		if lo == "" || t.winPct() < m[lo].winPct() {
			lo = k
		}
		if hi == "" || t.winPct() > m[hi].winPct() {
			hi = k
		}
	}
	return
}

func expFixes() {
	h(2, "6. Fix candidates (each applied to v0.2; Proposal F became v0.3)")
	fmt.Println("Each row changes one rule (or a set of them). Lower spreads and gaps mean a fairer game.")
	fmt.Println("**Route spread** uses identical bots that only differ in route. **Play-style spread** is the full tournament, comparing the five real play styles (Rusher, Scanner Runner, Alley Builder, Balanced, Showboat).")
	fmt.Println()
	var rows [][]string
	for _, f := range fixes {
		cfg := game.V02Config()
		f.mod(&cfg)
		rg := tournament(cfg, routeBots(), 4, *games, *seed+21)
		rm, _ := byKey(rg, func(p *game.Player) string { return p.Strat.Name })
		rlo, rhi := spread(rm)

		sg := tournament(cfg, roster(), 4, *games, *seed+31)
		sm, _ := byKey(sg, func(p *game.Player) string { return p.Strat.Name })
		sens := map[string]*tally{}
		for k, t := range sm {
			if sensible[k] {
				sens[k] = t
			}
		}
		slo, shi := spread(sens)
		seat, _ := byKey(sg, func(p *game.Player) string { return fmt.Sprint(p.ID) })
		firstW, firstN, fineParks, busted, players := 0, 0, 0, 0, 0
		for _, g := range sg {
			for _, p := range g.Players {
				players++
				busted += p.Stats.Busted
				fineParks += p.Stats.ParkCauses[game.ParkFine]
				if p.Stats.ShibuyaRank == 1 {
					firstN++
					if g.Winner == p.ID {
						firstW++
					}
				}
			}
		}
		rows = append(rows, []string{f.name,
			fmt.Sprintf("%.0f pts (%s %.0f%% … %s %.0f%%)", rm[rhi].winPct()-rm[rlo].winPct(), rlo, rm[rlo].winPct(), rhi, rm[rhi].winPct()),
			fmt.Sprintf("%.0f pts (%s … %s)", sm[shi].winPct()-sm[slo].winPct(), slo, shi),
			fmt.Sprintf("%.1f pts", seat["0"].winPct()-seat["3"].winPct()),
			pct(100 * float64(firstW) / float64(max(1, firstN))),
			f2(float64(fineParks) / float64(players)),
			f2(float64(busted) / float64(players)),
			fmt.Sprint(pctile(rounds(sg), 0.5)),
		})
	}
	table([]string{"Rules", "Route spread", "Play-style spread", "Seat 1 − Seat 4", "1st at Shibuya wins", "Turns lost to unpaid fines / player", "Busted / player", "Median rounds"}, rows)
}
