# Too Fast — Simulation Report

20000 games per experiment · seed 1 · rules v0.2 (experiments 1–4)

## 1. Strategy tournament (4 players, random bots and seats)

A fair game gives every strategy about **25%**. Per-player averages per game.

| Strategy | Win % | Turn at Shibuya | Upgrades at Shibuya | Out of gas | Turns parked | Pulled over | Busted | Rep paid in fines |
|---|---|---|---|---|---|---|---|---|
| Showboat | 45.1% | 8.5 | 2.1 | 0.07 | 0.17 | 0.38 | 0.06 | 0.4 |
| Alley Builder | 41.9% | 8.9 | 2.4 | 0.09 | 0.14 | 0.56 | 0.02 | 0.6 |
| Scanner Runner | 37.5% | 7.4 | 1.1 | 0.19 | 0.35 | 2.09 | 0.00 | 1.9 |
| Balanced | 27.2% | 7.5 | 0.9 | 0.27 | 0.18 | 1.68 | 0.00 | 1.7 |
| Random | 19.2% | 7.3 | 1.0 | 0.39 | 0.68 | 1.88 | 0.08 | 1.5 |
| Careless | 15.6% | 7.7 | 1.0 | 0.73 | 0.45 | 1.82 | 0.00 | 1.5 |
| Rusher | 7.8% | 7.3 | 0.6 | 0.20 | 1.09 | 3.95 | 0.00 | 2.8 |
| Tourist | 6.0% | 5.6 | 0.0 | 0.52 | 0.08 | 1.32 | 0.00 | 1.4 |

### Upgrade timing

Average turn each upgrade is bought, and how often it gets bought at all (across all bots that want it).

| Upgrade | Avg turn bought | Bought in % of games |
|---|---|---|
| Nitrous | 7.9 | 22.3% |
| Police Scanner | 2.4 | 38.5% |
| Turbo | 4.6 | 75.2% |
| Racing Tires | 8.1 | 37.9% |

### Seat order

| Seat | Win % |
|---|---|
| Seat 1 | 30.3% |
| Seat 2 | 26.6% |
| Seat 3 | 22.8% |
| Seat 4 | 20.3% |

### Runaway leader check

| Arrival order at Shibuya | Win % |
|---|---|
| #1 into Shibuya | 64.2% |
| #2 into Shibuya | 24.0% |
| #3 into Shibuya | 11.3% |
| #4 into Shibuya | 4.8% |

| Upgrades owned at Shibuya | Share of players | Win % |
|---|---|---|
| 0 | 32.4% | 6.1% |
| 1 | 37.8% | 26.0% |
| 2 | 18.0% | 38.3% |
| 3 | 7.4% | 51.0% |
| 4 | 4.3% | 58.5% |

### Where players get stuck

| Per player, per game | Average |
|---|---|
| Turns parked (lost) | 0.39 |
| Out-of-gas events | 0.31 |
| Spaces pushed back when out of gas | 2.84 |
| Times pulled over | 1.71 |
| Times Busted | 0.02 |
| Oil slick spin-outs | 0.04 |
| Street duels | 1.07 |
| Turns started broke (0 Rep) | 14.6% |
| Players with a 3+ turn stall | 0.1% |
| Longest stall seen | 4 turns |
| Games that never finished | 0 |

A **stall** is a run of turns in a row where the car made no forward progress (parked, pushed back, out of gas).

| Why turns are lost | Per game (all players) |
|---|---|
| Busted | 0.08 |
| Couldn't pay fine | 1.37 |
| Tow Zone | 0.28 |
| Fixing flat | 0.01 |
| Fuel delivery | 0.00 |

| Longest stall in a game (turns) | Share of players |
|---|---|
| 0 | 30.6% |
| 1 | 66.5% |
| 2 | 2.8% |
| 3 | 0.1% |
| 4 | 0.0% |
| 5 | 0.0% |
| 6+ | 0.0% |
## 2. Route balance (identical bots, only the route differs)


### Bots buy Turbo → Tires → Scanner → Nitrous

Fair = about 25%.

| Route | Win % | Turn at Shibuya | Upgrades at Shibuya | Rep at Shibuya | Pulled over | Out of gas | Turns played |
|---|---|---|---|---|---|---|---|
| C→E | 42.5% | 8.5 | 2.1 | 3.3 | 0.28 | 0.07 | 12.2 |
| C→D | 39.8% | 8.6 | 2.1 | 3.4 | 1.01 | 0.07 | 12.1 |
| B→E | 24.9% | 7.4 | 0.9 | 3.0 | 1.57 | 0.26 | 12.2 |
| A→E | 16.9% | 7.1 | 0.5 | 1.9 | 3.00 | 0.24 | 12.2 |
| B→D | 16.3% | 7.3 | 0.9 | 3.0 | 2.45 | 0.27 | 12.2 |
| A→D | 9.4% | 7.1 | 0.5 | 1.9 | 3.85 | 0.24 | 12.2 |

### Bots buy Scanner first → Turbo → Tires → Nitrous

Fair = about 25%.

| Route | Win % | Turn at Shibuya | Upgrades at Shibuya | Rep at Shibuya | Pulled over | Out of gas | Turns played |
|---|---|---|---|---|---|---|---|
| C→D | 30.7% | 8.5 | 2.4 | 3.0 | 0.50 | 0.09 | 11.8 |
| A→D | 29.8% | 7.3 | 1.1 | 2.7 | 2.01 | 0.18 | 11.9 |
| A→E | 24.4% | 7.2 | 1.1 | 2.6 | 1.40 | 0.18 | 11.9 |
| B→D | 23.3% | 6.8 | 1.3 | 2.8 | 1.20 | 0.32 | 11.9 |
| C→E | 23.3% | 8.4 | 2.3 | 3.0 | 0.15 | 0.10 | 11.9 |
| B→E | 18.3% | 6.9 | 1.3 | 2.7 | 0.76 | 0.31 | 11.9 |
## 3. Game length

| Players | Rounds (10th pct) | Median | 90th pct | Longest | Median play time |
|---|---|---|---|---|---|
| 2 | 11 | 14 | 17 | 33 | ~14 min |
| 3 | 11 | 13 | 15 | 27 | ~20 min |
| 4 | 11 | 13 | 15 | 22 | ~26 min |

Play time assumes about 30 seconds per turn. Change it with `-sec-per-turn`.
## 4. Card impact

Win rate of players who drew each card at least once (4 players, so the baseline is 25%).
Back Alley drivers draw more cards, so treat small differences with care.

| Obstacle | Drawn per game | Win % after drawing |
|---|---|---|
| Tow Zone | 0.31 | 17.3% |
| Wrong Turn | 0.31 | 19.4% |
| Fuel Leak | 0.46 | 20.1% |
| Flat Tire | 0.47 | 21.3% |
| Traffic Jam | 0.47 | 21.8% |
| Construction | 0.30 | 22.5% |
| Rival Crew | 0.31 | 24.0% |
| Engine Trouble | 0.31 | 24.5% |
| Rain | 0.31 | 24.9% |
| Police Checkpoint | 0.45 | 27.0% |
| Caught on Camera | 0.15 | 30.1% |

| Boost | Drawn per game | Win % after drawing |
|---|---|---|
| Shortcut | 0.51 | 35.5% |
| Empty Street | 0.50 | 31.1% |
| Sponsor Deal | 0.51 | 30.9% |
| Slipstream | 0.50 | 25.1% |
| Green Wave | 0.50 | 24.4% |
| Jerry Can | 0.50 | 24.2% |
| Mechanic on Call | 0.51 | 24.0% |
| Bribe a Cop | 0.76 | 23.5% |
| Scanner Tip | 0.51 | 23.3% |
| Drift King | 0.25 | 23.2% |
| Oil Slick | 0.75 | 22.4% |
| Fake Plates | 0.51 | 21.9% |
