// Package game is a rules engine for the Too Fast board game. It follows
// RULES.md (v0.2) and BOARD.md so bots can play thousands of games.
package game

import (
	"fmt"
	"strings"
)

type SpaceKind int

const (
	Plain SpaceKind = iota
	Rep3
	Rep5
	BoostSpace
	ObstacleSpace
	Gas
	Trap
	Landmark
)

func (k SpaceKind) String() string {
	return [...]string{"Plain", "+3 Rep", "+5 Rep", "Boost", "Obstacle", "Gas", "Speed trap", "Landmark"}[k]
}

type Space struct {
	Kind   SpaceKind
	Name   string // landmark name
	Narrow bool
}

type Route struct {
	ID     string
	Name   string
	Leg    int
	Spaces []Space // index 0 is space 1 on the board
}

func (r *Route) Len() int { return len(r.Spaces) }

// At returns the space at a 1-based board position.
func (r *Route) At(pos int) Space { return r.Spaces[pos-1] }

// RouteSeqs mirrors the `seq` strings in tools/make_board.py.
// board_test.go checks that the two stay in sync.
var RouteSeqs = []struct {
	ID, Name string
	Leg      int
	Seq      string
	Narrow   [][2]int
}{
	{"A", "Wangan Expressway", 1, ". R S . B . L:Daikoku_PA . S O . G . S R . L:Rainbow_Bridge . O S . B . S", nil},
	{"B", "Ginza Avenue", 1, ". B R . O . S R . L:Ginza B . G O . R . B S . L:Akihabara O R . B . R .", nil},
	{"C", "Back Alleys", 1, ". R B R . O R L:Tsukiji B . G O B R . G R B . R O L:Ameyoko B G . R B O R . B R", [][2]int{{5, 12}, {19, 26}}},
	{"D", "Shuto C1 Loop", 2, ". S . B S . O R S . G . S B . O S .", nil},
	{"E", "Harajuku–Roppongi", 2, ". R B . L:Harajuku R O . B S R . G B L:Roppongi O R . B R O .", nil},
}

type Board struct {
	Routes map[string]*Route
	Leg1   []string
	Leg2   []string
}

func ParseSeq(seq string) ([]Space, error) {
	var out []Space
	for _, tok := range strings.Fields(seq) {
		var s Space
		switch {
		case tok == ".":
			s.Kind = Plain
		case tok == "R":
			s.Kind = Rep3
		case tok == "R5":
			s.Kind = Rep5
		case tok == "B":
			s.Kind = BoostSpace
		case tok == "O":
			s.Kind = ObstacleSpace
		case tok == "G":
			s.Kind = Gas
		case tok == "S":
			s.Kind = Trap
		case strings.HasPrefix(tok, "L:"):
			s.Kind = Landmark
			s.Name = strings.ReplaceAll(tok[2:], "_", " ")
		default:
			return nil, fmt.Errorf("unknown space code %q", tok)
		}
		out = append(out, s)
	}
	return out, nil
}

func NewBoard() *Board { return NewBoardWith(nil, nil) }

// NewBoardWith builds the board, replacing any routes named in override.
func NewBoardWith(override map[string]string, narrow map[string][][2]int) *Board {
	b := &Board{Routes: map[string]*Route{}}
	for _, rs := range RouteSeqs {
		seq, nar := rs.Seq, rs.Narrow
		if s, ok := override[rs.ID]; ok {
			seq = s
		}
		if n, ok := narrow[rs.ID]; ok {
			nar = n
		}
		sp, err := ParseSeq(seq)
		if err != nil {
			panic(err)
		}
		for _, n := range nar {
			for i := n[0]; i <= n[1]; i++ {
				sp[i-1].Narrow = true
			}
		}
		b.Routes[rs.ID] = &Route{ID: rs.ID, Name: rs.Name, Leg: rs.Leg, Spaces: sp}
		if rs.Leg == 1 {
			b.Leg1 = append(b.Leg1, rs.ID)
		} else {
			b.Leg2 = append(b.Leg2, rs.ID)
		}
	}
	return b
}

// AvgLeg2 is the mean leg-2 length, used to compare progress across legs.
func (b *Board) AvgLeg2() float64 {
	t := 0
	for _, id := range b.Leg2 {
		t += b.Routes[id].Len()
	}
	return float64(t) / float64(len(b.Leg2))
}
