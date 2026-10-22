# Too Fast — Simulation Report

20000 games per experiment · seed 1 · rules v0.3 (experiments 1–4)

## 1. Strategy tournament (4 players, random bots and seats)

A fair game gives every strategy about **25%**. Per-player averages per game.

| Strategy | Win % | Turn at Shibuya | Upgrades at Shibuya | Out of gas | Turns parked | Pulled over | Busted | Rep paid in fines |
|---|---|---|---|---|---|---|---|---|
| Showboat | 36.1% | 8.4 | 2.0 | 0.06 | 0.14 | 0.35 | 0.02 | 0.3 |
| Scanner Runner | 36.0% | 7.3 | 1.3 | 0.17 | 0.15 | 2.09 | 0.00 | 0.6 |
| Alley Builder | 30.8% | 8.8 | 2.3 | 0.09 | 0.11 | 0.56 | 0.01 | 0.4 |
| Balanced | 29.6% | 7.5 | 1.2 | 0.19 | 0.15 | 1.32 | 0.00 | 0.8 |
| Rusher | 27.0% | 7.5 | 1.0 | 0.10 | 0.55 | 3.27 | 0.00 | 1.0 |
| Random | 20.5% | 7.4 | 1.1 | 0.36 | 0.50 | 1.52 | 0.07 | 0.7 |
| Careless | 16.1% | 8.3 | 1.2 | 0.75 | 0.31 | 1.55 | 0.00 | 0.8 |
| Tourist | 4.3% | 5.9 | 0.0 | 0.44 | 0.07 | 0.87 | 0.00 | 0.5 |

### Upgrade timing

Average turn each upgrade is bought, and how often it gets bought at all (across all bots that want it).

| Upgrade | Avg turn bought | Bought in % of games |
|---|---|---|
| Nitrous | 7.9 | 20.9% |
| Police Scanner | 2.1 | 37.0% |
| Turbo | 3.5 | 79.1% |
| Racing Tires | 7.9 | 41.0% |

### Seat order

| Seat | Win % |
|---|---|
| Seat 1 | 26.0% |
| Seat 2 | 23.5% |
| Seat 3 | 22.8% |
| Seat 4 | 27.6% |

### Runaway leader check

| Arrival order at Shibuya | Win % |
|---|---|
| #1 into Shibuya | 57.2% |
| #2 into Shibuya | 28.2% |
| #3 into Shibuya | 12.4% |
| #4 into Shibuya | 4.9% |

| Upgrades owned at Shibuya | Share of players | Win % |
|---|---|---|
| 0 | 24.5% | 6.5% |
| 1 | 40.2% | 26.9% |
| 2 | 24.3% | 31.8% |
| 3 | 7.9% | 43.2% |
| 4 | 3.1% | 46.3% |

### Where players get stuck

| Per player, per game | Average |
|---|---|
| Turns parked (lost) | 0.25 |
| Out-of-gas events | 0.27 |
| Spaces pushed back when out of gas | 2.51 |
| Times pulled over | 1.44 |
| Times Busted | 0.01 |
| Oil slick spin-outs | 0.04 |
| Street duels | 1.08 |
| Turns started broke (0 Rep) | 12.6% |
| Players with a 3+ turn stall | 0.1% |
| Longest stall seen | 4 turns |
| Games that never finished | 0 |

A **stall** is a run of turns in a row where the car made no forward progress (parked, pushed back, out of gas).

| Why turns are lost | Per game (all players) |
|---|---|
| Busted | 0.06 |
| Couldn't pay fine | 0.73 |
| Tow Zone | 0.28 |
| Fixing flat | 0.01 |
| Fuel delivery | 0.00 |

| Longest stall in a game (turns) | Share of players |
|---|---|
| 0 | 37.4% |
| 1 | 60.6% |
| 2 | 2.0% |
| 3 | 0.1% |
| 4 | 0.0% |
| 5 | 0.0% |
| 6+ | 0.0% |

## 2. Route balance (identical bots, only the route differs)


### Bots buy Turbo → Tires → Scanner → Nitrous

Fair = about 25%.

| Route | Win % | Turn at Shibuya | Upgrades at Shibuya | Rep at Shibuya | Pulled over | Out of gas | Turns played |
|---|---|---|---|---|---|---|---|
| C→D | 30.4% | 8.5 | 2.0 | 2.7 | 0.89 | 0.06 | 11.8 |
| C→E | 25.9% | 8.5 | 2.0 | 2.7 | 0.26 | 0.06 | 11.8 |
| A→D | 25.8% | 7.2 | 0.8 | 2.0 | 2.96 | 0.11 | 11.7 |
| A→E | 23.5% | 7.3 | 0.8 | 2.1 | 2.09 | 0.11 | 11.8 |
| B→D | 22.8% | 7.4 | 1.2 | 2.6 | 2.00 | 0.18 | 11.8 |
| B→E | 21.5% | 7.4 | 1.1 | 2.6 | 1.24 | 0.19 | 11.8 |

### Bots buy Scanner first → Turbo → Tires → Nitrous

Fair = about 25%.

| Route | Win % | Turn at Shibuya | Upgrades at Shibuya | Rep at Shibuya | Pulled over | Out of gas | Turns played |
|---|---|---|---|---|---|---|---|
| A→D | 34.0% | 7.3 | 1.3 | 2.0 | 2.09 | 0.16 | 11.9 |
| C→D | 27.4% | 8.7 | 2.3 | 2.7 | 0.53 | 0.10 | 12.0 |
| B→D | 27.0% | 7.3 | 1.5 | 2.4 | 1.23 | 0.27 | 12.0 |
| A→E | 24.1% | 7.2 | 1.3 | 1.9 | 1.38 | 0.17 | 12.0 |
| B→E | 19.2% | 7.2 | 1.5 | 2.4 | 0.76 | 0.28 | 12.0 |
| C→E | 18.3% | 8.6 | 2.3 | 2.7 | 0.15 | 0.10 | 12.0 |

## 3. Game length

| Players | Rounds (10th pct) | Median | 90th pct | Longest | Median play time |
|---|---|---|---|---|---|
| 2 | 12 | 14 | 16 | 27 | ~14 min |
| 3 | 11 | 13 | 15 | 23 | ~20 min |
| 4 | 11 | 12 | 14 | 21 | ~24 min |

Play time assumes about 30 seconds per turn. Change it with `-sec-per-turn`.

## 4. Card impact

Win rate of players who drew each card at least once (4 players, so the baseline is 25%).
Back Alley drivers draw more cards, so treat small differences with care.

| Obstacle | Drawn per game | Win % after drawing |
|---|---|---|
| Tow Zone | 0.31 | 13.6% |
| Flat Tire | 0.46 | 17.6% |
| Wrong Turn | 0.31 | 17.8% |
| Fuel Leak | 0.47 | 18.1% |
| Construction | 0.30 | 19.6% |
| Traffic Jam | 0.46 | 19.7% |
| Engine Trouble | 0.31 | 20.6% |
| Rival Crew | 0.31 | 22.2% |
| Rain | 0.31 | 23.8% |
| Police Checkpoint | 0.45 | 24.6% |
| Caught on Camera | 0.15 | 27.1% |

| Boost | Drawn per game | Win % after drawing |
|---|---|---|
| Shortcut | 0.50 | 32.6% |
| Empty Street | 0.50 | 30.2% |
| Sponsor Deal | 0.51 | 29.7% |
| Slipstream | 0.50 | 24.9% |
| Green Wave | 0.51 | 23.2% |
| Jerry Can | 0.50 | 23.1% |
| Mechanic on Call | 0.51 | 22.1% |
| Bribe a Cop | 0.76 | 22.1% |
| Scanner Tip | 0.51 | 22.0% |
| Drift King | 0.25 | 21.5% |
| Oil Slick | 0.76 | 21.3% |
| Fake Plates | 0.51 | 20.4% |

## 5. Rule variants

Same 4-player tournament, one rule changed at a time.

| Rules | Median rounds | Out of gas / player | Turns parked / player | Best bot | Worst bot | Careless win % | Tourist win % |
|---|---|---|---|---|---|---|---|
| Current rules (v0.3) | 12 | 0.27 | 0.25 | Showboat 35.9% | Tourist 3.9% | 16.3% | 3.9% |
| Fuel delivery variant | 12 | 0.26 | 0.38 | Showboat 35.7% | Tourist 4.3% | 17.0% | 4.3% |
| Thirsty nitrous variant | 12 | 0.28 | 0.24 | Showboat 35.9% | Tourist 4.2% | 16.6% | 4.2% |
| Fuel tank 4 | 13 | 0.60 | 0.24 | Showboat 35.8% | Tourist 4.1% | 21.1% | 4.1% |
| Fuel tank 6 | 12 | 0.15 | 0.25 | Showboat 37.5% | Tourist 3.5% | 19.2% | 3.5% |
| No refuel at Shibuya | 14 | 0.91 | 0.25 | Showboat 44.2% | Tourist 3.7% | 14.0% | 3.7% |
| Speed traps catch at 7+ (not 6+) | 12 | 0.27 | 0.25 | Showboat 35.9% | Tourist 3.9% | 16.3% | 3.9% |
| Start with 6 Rep | 12 | 0.24 | 0.23 | Showboat 41.5% | Tourist 2.2% | 14.1% | 2.2% |

## 6. Fix candidates (each applied to v0.2; Proposal F became v0.3)

Each row changes one rule (or a set of them). Lower spreads and gaps mean a fairer game.
**Route spread** uses identical bots that only differ in route. **Play-style spread** is the full tournament, comparing the five real play styles (Rusher, Scanner Runner, Alley Builder, Balanced, Showboat).

| Rules | Route spread | Play-style spread | Seat 1 − Seat 4 | 1st at Shibuya wins | Turns lost to unpaid fines / player | Busted / player | Median rounds |
|---|---|---|---|---|---|---|---|
| v0.2 (before any fixes) | 32 pts (A→D 10% … C→E 42%) | 38 pts (Rusher … Showboat) | 10.4 pts | 64.2% | 0.34 | 0.02 | 13 |
| Alley +5 spaces become +3 | 28 pts (A→D 11% … C→E 39%) | 33 pts (Rusher … Showboat) | 11.0 pts | 64.4% | 0.35 | 0.01 | 13 |
| Speed traps catch at 7+ | 13 pts (B→D 20% … C→D 33%) | 16 pts (Rusher … Showboat) | 12.0 pts | 63.0% | 0.26 | 0.02 | 12 |
| Heat 0 = warning, no fine | 27 pts (A→D 12% … C→D 39%) | 32 pts (Rusher … Showboat) | 11.0 pts | 63.4% | 0.16 | 0.02 | 13 |
| Later seats start +1 Rep each | 37 pts (A→D 7% … C→E 44%) | 46 pts (Rusher … Showboat) | -5.5 pts | 66.5% | 0.40 | 0.03 | 12 |
| Shortcut moves 6 max | 31 pts (A→D 10% … C→E 41%) | 37 pts (Rusher … Showboat) | 10.6 pts | 64.4% | 0.34 | 0.02 | 13 |
| Passing a trap safely: +1 Heat | 37 pts (A→D 8% … C→E 45%) | 44 pts (Rusher … Showboat) | 10.2 pts | 62.2% | 0.45 | 0.17 | 13 |
| Pulled over: Heat −2 (not reset) | 32 pts (A→D 10% … C→E 42%) | 38 pts (Rusher … Showboat) | 10.4 pts | 64.1% | 0.34 | 0.02 | 13 |
| Oil slicks trigger when driven over | 36 pts (A→D 8% … C→E 44%) | 41 pts (Rusher … Showboat) | 9.6 pts | 62.2% | 0.35 | 0.02 | 13 |
| Seats 3 and 4 start +1 Rep | 33 pts (A→D 9% … C→E 42%) | 40 pts (Rusher … Showboat) | 6.6 pts | 64.6% | 0.42 | 0.02 | 13 |
| No bonus for 1st into Shibuya | 37 pts (A→D 9% … C→E 46%) | 39 pts (Rusher … Showboat) | 8.8 pts | 56.3% | 0.38 | 0.01 | 13 |
| Shibuya catch-up: +2 Rep per place behind 1st | 41 pts (A→D 7% … C→D 48%) | 46 pts (Rusher … Alley Builder) | 7.6 pts | 48.9% | 0.34 | 0.02 | 13 |
| Natural 5 or 6 adds Heat | 32 pts (A→D 10% … C→E 42%) | 38 pts (Rusher … Showboat) | 10.4 pts | 64.1% | 0.40 | 0.05 | 13 |
| **Proposal A**: alley +3, traps 7+, warning, seat Rep, Shortcut 6 | 8 pts (B→E 23% … C→D 30%) | 17 pts (Rusher … Scanner Runner) | -3.2 pts | 64.3% | 0.18 | 0.02 | 12 |
| **Proposal B**: A + passing a trap safely gives +1 Heat | 14 pts (A→D 17% … C→E 31%) | 25 pts (Rusher … Showboat) | -3.4 pts | 61.9% | 0.31 | 0.20 | 12 |
| **Proposal C**: A, but seats 3–4 +1 Rep and Shibuya catch-up | 19 pts (A→E 18% … C→D 38%) | 16 pts (Rusher … Showboat) | 6.7 pts | 49.8% | 0.19 | 0.02 | 12 |
| **Proposal E**: A, but seat Rep +0/+0/+1/+2 | 7 pts (B→E 22% … C→D 29%) | 13 pts (Rusher … Scanner Runner) | 0.1 pts | 64.1% | 0.17 | 0.02 | 12 |
| **Proposal F**: E + no bonus for 1st into Shibuya | 9 pts (B→E 22% … C→D 31%) | 8 pts (Rusher … Showboat) | -0.7 pts | 58.2% | 0.18 | 0.01 | 12 |
