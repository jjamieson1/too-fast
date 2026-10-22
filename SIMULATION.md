# Too Fast — Simulation Findings

A Go program plays thousands of games of Too Fast with bot players, so we can spot balance problems before real playtesting.

- **Rules tested:** [RULES.md](RULES.md) v0.2, then v0.3 (v0.2 plus the fixes below), and the board in [BOARD.md](BOARD.md)
- **Full reports:** [sim-results/report-v0.3.md](sim-results/report-v0.3.md) (current rules, all experiments, including the fix tests) and [sim-results/report-v0.2.md](sim-results/report-v0.2.md) (the old rules, for comparison)
- **Volume:** 20,000 games per experiment

## TL;DR

These findings are for the **v0.2** rules. The recommended fixes are now in **v0.3** (see [Adopted in v0.3](#adopted-in-v03-proposal-f) below).

1. **The Back Alleys are too strong, and the Wangan Expressway is too weak.** With identical bots, Alley players win about 42% of the time, while Expressway → C1 players win about 10% (fair is 25%).
2. **Going first is a big advantage.** Seat 1 wins 30%, seat 4 wins 20%.
3. **The first car into Shibuya wins 64% of games.** Comebacks are rare.
4. **Nobody gets truly stuck.** No game ever stalled, and only 0.1% of players lost 3+ turns in a row. Most lost turns come from being **broke at a speed trap**.
5. **Fuel works well.** Players who plan their stops rarely run out (0.06–0.27 times per game). Careless players run out about 0.75 times per game, so fuel rewards planning without wrecking anyone's game.
6. **Heat barely matters yet.** Players get Busted only about once every 50 games, because Heat resets whenever you're pulled over and drops at every gas station.
7. **Games are a bit short:** about 13 rounds, or roughly 25 minutes with 4 players at 30 seconds a turn. The target is 30–45 minutes, though real players with cards and chatter will probably take longer per turn.

**The fix ("Proposal F", now in v0.3):** six small rule changes. Together they bring every route and play style to between 19% and 36% wins, remove the seat advantage, and halve the turns lost to unpaid fines.

---

## How to run it

```bash
go run ./cmd/sim
```

That runs the full report with 20,000 games per experiment and takes about 10 seconds. Other useful commands:

```bash
go run ./cmd/sim -games 2000 -exp routes
```

```bash
go run ./cmd/sim -rules v0.2
```

```bash
go run ./cmd/sim -trace -seed 7
```

```bash
go test ./...
```

- `-exp` picks one experiment: `strategies`, `routes`, `length`, `cards`, `variants` or `fixes`.
- `-rules v0.2` runs experiments 1–4 with the old rules, for comparison. The default is v0.3.
- `-trace` prints a turn-by-turn play log of a single game. It's useful for seeing exactly where a player got stuck.
- `go test ./...` checks the engine rules, and checks that the board data matches [tools/make_board.py](tools/make_board.py).

**Code layout:**
- [game/](game/): the rules engine. Every tunable number lives in [game/config.go](game/config.go), so you can test a change without touching the engine.
- [game/strategy.go](game/strategy.go): how bots make decisions.
- [cmd/sim/main.go](cmd/sim/main.go): the bot lineup and the experiments.

### The bots

All bots play their cards sensibly. For example, they bribe cops when a fine is 2 or more, drop oil slicks in front of the leader, and cancel bad obstacles with *Mechanic on Call*. They differ in the choices that define a play style:

| Bot | Route | Buys (in order) | Notes |
|---|---|---|---|
| **Rusher** | A → D | Turbo, Tires, Nitrous | Always fires Nitrous |
| **Scanner Runner** | A → D | Scanner, Turbo, Tires, Nitrous | Expressway with a Police Scanner |
| **Alley Builder** | C → D or E | Scanner, Turbo, Tires, Nitrous | Takes C1 if it owns a Scanner |
| **Balanced** | B → E | Turbo, Tires, Nitrous | |
| **Showboat** | C → E | Turbo, Tires, Nitrous, Scanner | |
| **Tourist** | B → E | Nothing | Checks whether upgrades matter |
| **Careless** | B → E | Turbo, Tires | Only stops for gas when empty |
| **Random** | Random | Random | |

---

## Findings in detail (v0.2 rules)

### 1. Route balance: the alleys win too much

With every player using the same bot, so that only the route differs:

| Route | Win % | Upgrades by Shibuya | Times pulled over |
|---|---|---|---|
| C → E (Alleys → Harajuku) | **42%** | 2.1 | 0.3 |
| C → D (Alleys → C1) | **40%** | 2.1 | 1.0 |
| B → E | 25% | 0.9 | 1.6 |
| A → E | 17% | 0.5 | 3.0 |
| B → D | 16% | 0.9 | 2.5 |
| A → D (Expressway → C1) | **9%** | 0.5 | 3.9 |

**Why it happens:**
- **The alleys cost only about 1.5 turns more** to reach Shibuya (8.5 turns vs 7.0).
- **That extra time buys roughly 2 upgrades** by halfway, worth +1 to +2 movement on every remaining turn.
- **The Expressway's 5 speed traps** (10 counting C1) stop a fast car over and over. Each catch ends your move early, and the fines keep you too poor to buy upgrades.

Buying a Police Scanner first makes the Expressway competitive (A → D rises to 30%). That's a nice combo, but the expressway only works if you know to do it.

### 2. Seat order: seat 1 wins 30%, seat 4 wins 20%

This is the usual first-player advantage in a race game, but 10 points is large.

### 3. The halfway leader usually wins

| Arrival at Shibuya | Win % |
|---|---|
| 1st | 64% |
| 2nd | 24% |
| 3rd | 11% |
| 4th | 5% |

The +3 Rep bonus for arriving first makes this worse. Leg 2 is short (18–22 spaces), so there isn't much time to catch up.

### 4. Upgrades matter a lot, maybe slightly too much

| Upgrades at Shibuya | Win % |
|---|---|
| 0 | 6% |
| 1 | 26% |
| 2 | 38% |
| 3 | 51% |
| 4 | 59% |

The Tourist bot, which never buys anything, wins only 6% of games. It's good that upgrades matter, but this is also why the Rep-rich alleys come out on top.

### 5. Where players get stuck

**The good news:** there are no deadlocks. All games finished, and the longest stall seen was 4 turns. Fewer than 1 in 1,000 players lost 3 or more turns in a row.

**The feel-bad spot is being broke.** Players start 15% of their turns with 0 Rep. At 0 Rep, even a 1-Rep fine parks you, which makes "couldn't pay a fine" the number one cause of lost turns: about 1.4 per game, compared with about 0.3 per game from *Tow Zone*. A broke player loses a turn even at Heat 0, which feels unfair for someone driving carefully.

### 6. Fuel

| | Out of gas per player per game |
|---|---|
| Bots that plan their gas stops | 0.06–0.27 |
| Careless bot | 0.73 |
| Everyone (current rules) | 0.31 |

Players who run out lose about 9 spaces when it happens. This is about right: it punishes careless driving without being brutal. Ginza Avenue's 16-space last stretch didn't cause problems.

| Fuel variant | Out of gas per player |
|---|---|
| Tank of 4 | 0.73 (too punishing) |
| Tank of 6 | 0.18 |
| No refuel at Shibuya | 0.88 (keep the refuel) |

### 7. Heat

Players get **Busted** about once every 50 games, and the average fine is only about 1 Rep. Heat rarely climbs past 1–2, because:
- being pulled over resets it to 0,
- gas stations take away 2, and
- the only steady source is rolling a natural 6.

Two ideas the simulator tested:
- **Gaining +1 Heat for passing a trap safely** made the Expressway even worse.
- **Natural 5s and 6s adding Heat** barely changed anything.

Heat needs a design rethink rather than a number tweak. One idea: **don't reset Heat when pulled over, and raise the fines** (for example, a fine of 2× your Heat).

### 8. Cards

| Card | Win % after drawing (baseline 25%) | Note |
|---|---|---|
| **Shortcut** | 36% | Strongest card. It can jump 10+ spaces to the next gas station. |
| Empty Street | 31% | Strong |
| Sponsor Deal | 31% | Strong. Rep is powerful. |
| **Tow Zone** | 17% | Harshest card. It's the only one that skips a turn outright. |
| Caught on Camera | 30% | Actually *helps* you, since +2 Rep is worth more than +2 Heat. Another sign Heat is too weak. |
| Oil Slick | 22% | Rarely hits: about 0.16 spin-outs per game vs 0.75 slicks drawn |

### 9. Game length

| Players | Median rounds | Time at 30 s/turn |
|---|---|---|
| 2 | 14 | ~14 min |
| 3 | 13 | ~20 min |
| 4 | 13 | ~26 min |

If real playtests also run short, lengthen each route by 3–4 spaces.

---

## Adopted in v0.3: Proposal F

The simulator tested each fix on its own, then in combinations (experiment 6 in the report). The best set was:

| # | Change | Fixes |
|---|---|---|
| 1 | **Back Alley +5 spaces become +3** | Alleys too rich |
| 2 | **Speed traps catch cars passing at 7+ movement** (was 6+) | Expressway too punishing |
| 3 | **At Heat 0 you get a warning, not a fine** | Broke, careful players losing turns |
| 4 | **Seats 3 and 4 start with +1 and +2 Rep** | First-player advantage |
| 5 | **No bonus for the first car into Shibuya** | Halfway leader winning too often |
| 6 | **Shortcut moves 6 spaces at most** | Strongest card |

| | v0.2 | v0.3 (Proposal F) |
|---|---|---|
| Route spread (worst → best win %) | 9% → 42% | 22% → 31% |
| Play-style spread (5 real styles) | 8% → 45% | 27% → 36% |
| Seat 1 vs seat 4 | 30% vs 20% | 26% vs 28% |
| 1st into Shibuya wins | 64% | 57% |
| Turns lost to unpaid fines | 0.34 per player | 0.18 per player |
| Median length | 13 rounds | 12 rounds |

**Things that didn't help:**
- **Catch-up Rep at Shibuya for late arrivals:** it cut the leader's win rate to 49%, but it rewarded the slow alleys and broke the route balance.
- **Oil slicks that trigger when driven over:** this hurt the Expressway more.
- **Pulled over lowering Heat by 2 instead of resetting it:** this did nothing on its own, because Heat is rarely above 2.
- **Fuel delivery and thirsty nitrous variants:** both had almost no effect. Bots rarely chose them.

## Rule questions found while coding

Writing the engine exposed spots where the v0.2 rules didn't say what happens. **v0.3 now answers questions 1–5 in RULES.md.** Question 6 is still open:

1. **Nitrous vs. fuel order:** you fire Nitrous in step 1, but fuel is checked in step 2. Can you waste Nitrous by firing it and then running out of gas? *v0.3: Nitrous now fires after you burn fuel.*
2. **Busted mid-move:** if a natural 6 busts you, do you still move this turn? *Sim: yes, you finish this move and sit out next turn.*
3. **Ginza and Sponsor Deal:** if you draw a Sponsor Deal at Ginza, does it count as your "keep 1"? *v0.3: yes.*
4. **Pushed back to Start or Shibuya:** can you pick a different route next turn? *Sim: yes.*
5. **Traffic Jam and fuel:** does a jammed turn still burn fuel? *Sim: yes.*
6. **Penalties stacking:** do Engine Trouble or Rain stack with a Traffic Jam? *Sim: Traffic Jam ignores them, and they're used up.*

## Limitations

- **Bots aren't people.** Real players will find combos and make mistakes the bots don't. Use these numbers to decide **what to watch in playtests**, not as final answers.
- **The simulator can't measure fun.** A rule can be balanced and still feel bad (or the reverse).
- **Play time assumes 30 seconds per turn.** Time a real game and pass `-sec-per-turn` to match.
