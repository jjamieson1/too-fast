package game

import (
	"os"
	"regexp"
	"testing"
)

func testGame(t *testing.T, n int) *Game {
	t.Helper()
	var s []*Strategy
	for i := 0; i < n; i++ {
		s = append(s, &Strategy{Name: "bot", Leg1: "A", Leg2: "D", FuelStop: 0})
	}
	return New(DefaultConfig(), s, 1)
}

func place(p *Player, leg int, route string, pos int) { p.Leg, p.Route, p.Pos = leg, route, pos }

func TestBoardMatchesSketchScript(t *testing.T) {
	src, err := os.ReadFile("../tools/make_board.py")
	if err != nil {
		t.Fatal(err)
	}
	re := regexp.MustCompile(`"([A-E])": dict\(name="[^"]+", seq="([^"]+)"`)
	found := map[string]string{}
	for _, m := range re.FindAllStringSubmatch(string(src), -1) {
		found[m[1]] = m[2]
	}
	for _, rs := range RouteSeqs {
		if found[rs.ID] != rs.Seq {
			t.Errorf("route %s differs from tools/make_board.py:\n go: %s\n py: %s", rs.ID, rs.Seq, found[rs.ID])
		}
	}
	want := map[string]int{"A": 24, "B": 28, "C": 32, "D": 18, "E": 22}
	for id, n := range want {
		if got := NewBoard().Routes[id].Len(); got != n {
			t.Errorf("route %s has %d spaces, want %d", id, got, n)
		}
	}
}

func TestDeckSizes(t *testing.T) {
	sum := func(m map[Boost]int) (n int) {
		for _, v := range m {
			n += v
		}
		return
	}
	if n := sum(DefaultBoostCounts); n != 25 {
		t.Errorf("boost deck has %d cards, want 25", n)
	}
	n := 0
	for _, v := range DefaultObstacleCounts {
		n += v
	}
	if n != 25 {
		t.Errorf("obstacle deck has %d cards, want 25", n)
	}
}

func TestSpeedTrapCatchesFastCarsPassingThrough(t *testing.T) {
	g := testGame(t, 1)
	p := g.Players[0]
	place(p, 1, "A", 1) // trap at A3
	g.moveForward(p, 7, true)
	if p.Pos != 3 {
		t.Fatalf("moving 7 past a trap should stop on it (A3), got A%d", p.Pos)
	}
	place(p, 1, "A", 1)
	g.moveForward(p, 6, true)
	if p.Pos != 7 {
		t.Fatalf("moving 6 should pass the trap in v0.3, got A%d", p.Pos)
	}
	place(p, 1, "A", 1)
	p.Up[Scanner] = true
	g.moveForward(p, 7, true)
	if p.Pos != 8 {
		t.Fatalf("a Police Scanner should pass safely, got A%d", p.Pos)
	}
}

func TestPulledOverFineAndParking(t *testing.T) {
	g := testGame(t, 1)
	p := g.Players[0]
	p.Heat, p.Rep = 3, 10
	g.pullOver(p)
	if p.Rep != 7 || p.Heat != 0 || p.Parked {
		t.Fatalf("fine should equal Heat: rep=%d heat=%d parked=%v", p.Rep, p.Heat, p.Parked)
	}
	p.Heat, p.Rep = 4, 2
	g.pullOver(p)
	if p.Rep != 0 || !p.Parked {
		t.Fatalf("can't pay in full → pay all and park: rep=%d parked=%v", p.Rep, p.Parked)
	}
	p.Parked, p.Heat, p.Rep = false, 0, 0
	g.pullOver(p)
	if p.Parked || p.Rep != 0 {
		t.Fatal("at Heat 0 a broke driver should just get a warning")
	}
}

func TestSeatRep(t *testing.T) {
	g := testGame(t, 4)
	for i, want := range []int{4, 4, 5, 6} {
		if got := g.Players[i].Rep; got != want {
			t.Errorf("seat %d starts with %d Rep, want %d", i+1, got, want)
		}
	}
}

func TestV02ConfigRestoresOldAlleys(t *testing.T) {
	c := V02Config()
	g := New(c, []*Strategy{{Name: "x"}}, 1)
	if g.Board.Routes["C"].At(7).Kind != Rep5 {
		t.Fatal("v0.2 alleys should have +5 Rep at C7")
	}
	if NewBoard().Routes["C"].At(7).Kind != Rep3 {
		t.Fatal("v0.3 alleys should have +3 Rep at C7")
	}
}

func TestBustedAtRedLine(t *testing.T) {
	g := testGame(t, 1)
	p := g.Players[0]
	p.Heat = 4
	g.addHeat(p, 1)
	if !p.Parked || p.Heat != 2 || p.Stats.Busted != 1 {
		t.Fatalf("heat 5 should bust: parked=%v heat=%d", p.Parked, p.Heat)
	}
}

func TestShibuyaMustStopRefuelAndBonus(t *testing.T) {
	g := testGame(t, 2)
	a, b := g.Players[0], g.Players[1]
	place(a, 1, "A", 22)
	a.Fuel, a.Rep = 1, 0
	g.Cfg.ShibuyaBonus = 3 // v0.2 bonus, to check only the first car gets it
	g.moveForward(a, 6, false)
	if a.Leg != 2 || a.Route != "" || a.Fuel != 5 || a.Rep != 3 {
		t.Fatalf("first at Shibuya: leg=%d route=%q fuel=%d rep=%d", a.Leg, a.Route, a.Fuel, a.Rep)
	}
	place(b, 1, "A", 24)
	b.Rep = 0
	g.moveForward(b, 3, false)
	if b.Rep != 0 {
		t.Fatalf("second car should not get the bonus, rep=%d", b.Rep)
	}
}

func TestOutOfGasPushesBackToStation(t *testing.T) {
	g := testGame(t, 1)
	p := g.Players[0]
	place(p, 1, "A", 20) // gas at A12
	p.Fuel, p.Heat = 0, 3
	g.outOfGas(p)
	if p.Pos != 12 || p.Fuel != 5 || p.Heat != 1 || p.Stats.OutOfGasSpaces != 8 {
		t.Fatalf("pos=%d fuel=%d heat=%d lost=%d", p.Pos, p.Fuel, p.Heat, p.Stats.OutOfGasSpaces)
	}
	place(p, 1, "A", 5) // no station behind → back to Start
	p.Fuel = 0
	g.outOfGas(p)
	if p.Route != "" || p.Fuel != 5 {
		t.Fatalf("should be back at Start, route=%q fuel=%d", p.Route, p.Fuel)
	}
}

func TestNarrowStreetBlocks(t *testing.T) {
	g := testGame(t, 2)
	a, b := g.Players[0], g.Players[1]
	place(a, 1, "C", 8) // narrow
	place(b, 1, "C", 5)
	g.moveForward(b, 3, false)
	if b.Pos != 7 {
		t.Fatalf("should stop behind the car on a narrow space, got C%d", b.Pos)
	}
}

func TestFinishWins(t *testing.T) {
	g := testGame(t, 1)
	p := g.Players[0]
	place(p, 2, "D", 16)
	g.moveForward(p, 5, false)
	if !p.Finished || g.Winner != 0 {
		t.Fatal("passing the end of leg 2 should win")
	}
}

func TestSameSeedSameGame(t *testing.T) {
	run := func() int {
		s := []*Strategy{
			{Name: "a", Leg1: "A", Leg2: "D", BuyOrder: []Upgrade{Turbo}, FuelStop: -1, NitrousMode: NitrousSmart},
			{Name: "b", Leg1: "C", Leg2: "E", BuyOrder: []Upgrade{Scanner}, FuelStop: -1},
		}
		g := New(DefaultConfig(), s, 42)
		g.Run()
		return g.Round*10 + g.Winner
	}
	if run() != run() {
		t.Fatal("same seed gave different games")
	}
}

func TestManyGamesFinish(t *testing.T) {
	for seed := uint64(0); seed < 500; seed++ {
		s := []*Strategy{
			{Name: "a", Leg1: "A", Leg2: "D", BuyOrder: []Upgrade{Turbo, Tires}, FuelStop: 0, NitrousMode: NitrousAlways},
			{Name: "b", Leg1: "B", Leg2: "E", FuelStop: -1},
			{Name: "c", Leg1: "C", Leg2: "", BuyOrder: []Upgrade{Scanner, Nitrous}, FuelStop: -1, NitrousMode: NitrousSmart},
			{Name: "d", Leg2: "?", BuyOrder: []Upgrade{Nitrous}, FuelStop: 1, NitrousMode: NitrousAlways},
		}
		g := New(DefaultConfig(), s, seed)
		g.Run()
		if g.Winner < 0 {
			t.Fatalf("seed %d never finished", seed)
		}
		for _, p := range g.Players {
			if p.Fuel < 0 || p.Fuel > 5 || p.Heat < 0 || p.Heat > 5 || p.Rep < 0 || len(p.Hand) > 3 {
				t.Fatalf("seed %d bad state %+v", seed, p)
			}
		}
	}
}
