# Too Fast — Design Doc

> A street-racing board game through the streets of Tokyo. Go fast, get famous, don't get caught.

Your original notes are saved in [notes/IDEA-original.md](notes/IDEA-original.md).

| | |
|---|---|
| **Players** | 2–4 |
| **Age** | 10+ |
| **Play time** | 30–45 minutes (target) |
| **Win condition** | First car to reach or pass the finish line at Tokyo Tower wins |

Related docs:
- [RULES.md](RULES.md): the rulebook draft
- [CARDS.md](CARDS.md): every Boost, Obstacle and Upgrade card
- [BOARD.md](BOARD.md): the board layout, space by space, plus [board-sketch.svg](board-sketch.svg)

---

## The core idea: going too fast draws heat

Every player tracks **Heat** (0–5), which measures how much police attention they've drawn.

- You gain Heat by rolling a natural 6, using nitrous, and from some cards and spaces.
- When you're **pulled over** at a speed trap, you pay Reputation equal to your Heat.
- If your Heat hits 5, you're **Busted**: you lose your next turn.

This gives every turn a push-your-luck choice. It also slows down the leader, who is usually the one driving hardest.

## The core loop

1. **Before you roll:** buy upgrades with Reputation and play Boost cards.
2. **Roll and move:** d6 plus your upgrades. Watch for speed traps.
3. **Land:** resolve the space (gain Reputation, draw a card, visit a landmark…).
4. **Clash:** if you end on another car's space, you have a Street Duel.

## The map: pick your route

The race runs in two legs, with Shibuya Crossing in the middle.

**Leg 1: Odaiba (start) → Shibuya Crossing**

| Route | Spaces | Character |
|---|---|---|
| **Wangan Expressway** | 24 | Fastest, but full of speed traps. Has Daikoku PA (cheap upgrades) and Rainbow Bridge (a paid shortcut). |
| **Ginza Avenue** | 28 | Balanced. Good card draws, a few traps. |
| **Back Alleys** | 32 | Slowest, but lots of Reputation. Narrow streets let you block other cars. No police. |

Every car **must stop** at Shibuya Crossing. The first car to arrive gets a bonus.

**Leg 2: Shibuya Crossing → Tokyo Tower (finish)**

| Route | Spaces | Character |
|---|---|---|
| **Shuto C1 Loop** | 18 | Short, but swarming with speed traps. |
| **Harajuku–Roppongi** | 22 | Longer, with landmarks that reward showing off. |

**The intended tension:** take the alleys in leg 1 to build Reputation and buy upgrades, then blast down the expressway in leg 2. Or go fast early and hope your lead holds.

## The economy: Reputation

- **Start with:** 4 Reputation (Rep).
- **Earn it from:** Rep spaces, landmarks, Sponsor cards, winning duels, and rolling a natural 1 ("keeping it cool").
- **Spend it on:** upgrades, fines, and some landmark actions.

## Upgrades (max one of each)

| Upgrade | Cost | Effect |
|---|---|---|
| **Nitrous** | 3 | One-shot burst: +3 to a roll, +1 Heat. Refills at gas stations. |
| **Police Scanner** | 3 | Speed traps catch you only if you *land* on them. |
| **Turbo** | 5 | +1 to every roll. |
| **Racing Tires** | 6 | +1 to every roll; immune to oil slicks and rain. |

Permanent bonuses cap at +2. That keeps an early leader from running away with the game.

## Catch-up and interaction

- **Slipstream:** start your turn 1 space behind another car and you get +1 to your roll.
- **Street Duel:** end your move on another car's space and you both roll. The winner takes 2 Rep from the loser.
- **Narrow streets:** in the Back Alleys, only one car fits per space, so you can block other cars.
- **Oil slicks:** Boost cards let you drop oil slicks for rivals to spin out on.
- **Rival Crew** (Obstacle card): you give Rep to the player behind you.

## Changes from the original notes

| Original | Now | Why |
|---|---|---|
| Flat tire appeared twice (miss a turn / re-roll with a minus) | One **Flat Tire** card: −2 to rolls until you stop at a gas station | Removes the duplicate, and gives gas stations a job |
| Construction = miss a turn | Construction = detour back 2 spaces | Missing turns isn't fun. Now only *Tow Zone* and getting *Busted* skip a turn outright. A few other cases do it only if you can't pay or choose to. |
| Ran out of gas, return to the nearest station | Kept as the **Out of Gas** card | No fuel gauge to track. Lighter bookkeeping. |
| Turbo / tires / nitrous each +1 | Each upgrade does something different | Real choices instead of three copies of +1 |
| Oil slick: back 4 spaces | Back 3 spaces; racing tires are immune | Slightly softer, and gives tires a purpose |
| Speed trap: pay Rep | Pay Rep equal to your Heat | Ties into the Heat system |
| "Traffic break Anna…" (unfinished) | Became **Green Wave** (+2 to your roll) | My best guess at the idea. Change it if you meant something else! |

## Components (prototype list)

- Game board (see [BOARD.md](BOARD.md))
- 4 car tokens (Hot Wheels–scale cars work well)
- 4 player dashboards (Heat track 0–5, 4 upgrade slots)
- 25 Boost cards, 25 Obstacle cards, 16 Upgrade cards (4 of each)
- Reputation chips: about 40 "1" chips and 10 "5" chips
- 4 Heat markers (small cubes)
- 6 oil slick tokens
- 1 six-sided die (or one per player)

## Open questions for playtesting

- [ ] Does the game land in 30–45 minutes? If it runs long, shorten routes by 2–4 spaces.
- [ ] Is the Wangan Expressway too strong or too weak? Adjust the speed trap count.
- [ ] Is Rep income high enough that most players buy 1–2 upgrades by Shibuya?
- [ ] Do players reach Heat 5 too often (or never)?
- [ ] Are any cards "feel-bad" cards that people groan at?
- [ ] Variant idea: **Street Cred mode**, where the first to finish gets +10 Rep and the most Rep wins.
