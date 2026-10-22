# Too Fast — Rulebook (Draft v0.3)

*2–4 players · ages 10+ · 30–45 minutes*

The midnight race across Tokyo starts at the Odaiba waterfront and ends at Tokyo Tower. The first car to the tower wins. But drive too fast and the police will notice. Forget to fill up and you won't get there at all.

---

## 1. Components

- 1 game board
- 4 car tokens
- 4 **player dashboards** (see §2)
- 25 Boost cards (green backs)
- 25 Obstacle cards (orange backs)
- 16 Upgrade cards (4 each of Nitrous, Police Scanner, Turbo, Racing Tires)
- Reputation chips ("1" and "5")
- 6 oil slick tokens
- 1 six-sided die

## 2. Your dashboard

Each player holds a large dashboard card. It's your car's control panel:

- **Fuel gauge** (a dial marked **5 → 0**): how many turns of driving you have left.
- **Heat dial** (a tachometer marked **0 → 5**): how much police attention you've drawn. **5 is the red line.**
- **Four upgrade slots:** Nitrous, Police Scanner, Turbo and Racing Tires. Slide each Upgrade card into its slot when you buy it.

A sketch is in [dashboard-sketch.svg](dashboard-sketch.svg). For a prototype, cut the dials out of card stock and pin them on with paper fasteners (brads), so they spin.

## 3. Setup

1. Put the board in the middle of the table.
2. Shuffle the Boost deck and the Obstacle deck separately, and put them face down next to the board.
3. Sort the Upgrade cards into four face-up piles. This is the **shop**.
4. Each player takes a car and a dashboard. Turn your **fuel to 5** and your **Heat to 0**.
5. Put all cars on the **Start** space at Odaiba.
6. The player who most recently rode in a car goes first. Play continues clockwise.
7. Each player takes **4 Reputation** (Rep). Going later is a disadvantage, so the **3rd player takes 1 extra Rep** (5 total) and the **4th player takes 2 extra** (6 total).

## 4. How to win

The **first car to reach or pass Tokyo Tower wins immediately**. You don't need an exact roll.

## 5. Your turn

Each turn has four steps, in this order.

### Step 1: Before you roll
- If your car is **parked** (lying on its side, see §10), stand it back up. Your turn ends. *Parked cars don't burn fuel.*
- You may **buy upgrades** from the shop (see §7).
- You may play Boost cards marked **Before roll**.

### Step 2: Burn fuel, roll and move
**First, burn fuel:** turn your fuel gauge down by 1.
- If your gauge is **already at 0**, you're **Out of Gas** (see §8). Do that instead of rolling. Your turn ends.

**Then you may fire your Nitrous**, if you have it and it's full.

**Then roll** the die and add your modifiers:
- **+1** for Turbo
- **+1** for Racing Tires
- **+1** for Slipstream (see §11)
- **+3** for Nitrous, if you fired it
- **−2** for a Flat Tire, plus any other card effects

The total is your **movement**. Movement is never less than 1.

**Natural rolls** (the number on the die, before modifiers):
- **Rolled a 6?** Gain **+1 Heat**. You're flooring it.
- **Rolled a 1?** Gain **+1 Rep**. You're keeping it cool.

**Move** your car forward that many spaces.
- **Forks:** at the start, and again at Shibuya Crossing, choose which route to take.
- **Speed traps:** see §9.
- **Shibuya Crossing:** you **must stop** there, even with movement left over (see §12).
- **Narrow streets:** see §11.
- **Gas stations:** if you pass one, you may stop there and end your move early. This is often a smart move when your fuel is low.

### Step 3: Resolve your space
Do what the space you stopped on says (see §6).

### Step 4: Clash
If you ended on the same space as another car, you have a **Street Duel** (see §11).

Then it's the next player's turn.

## 6. Spaces

| Space | What happens when you land on it |
|---|---|
| **Plain** | Nothing. |
| **Rep (+3)** | Gain 3 Rep. |
| **Boost** | Draw a Boost card and keep it in your hand. |
| **Obstacle** | Draw an Obstacle card and resolve it immediately. |
| **Gas Station** | Pit stop: **fill your fuel to 5**, lose 2 Heat, refill your Nitrous, and fix a Flat Tire. |
| **Speed Trap** | You're **pulled over** (see §9). |
| **Landmark** | Do what the landmark says (see [BOARD.md](BOARD.md)). |
| **Oil Slick** token | Spin out: move back 3 spaces and discard the token. Racing Tires are immune. |

**Moved by a card or landmark?**
- If you're moved **forward**, resolve the space you land on.
- If you're moved **backward**, don't resolve it (unless the rule says to).
- Speed traps only trigger during your normal roll movement.

**Moving backward:** you can never move back past Shibuya Crossing or the Start. Stop there. If you're pushed back onto one of them, you choose your route again the next time you leave.

## 7. Upgrades

At the start of your turn, you may buy upgrades by paying Rep to the bank. Slide each upgrade into its slot on your dashboard. You can only own **one of each**.

| Upgrade | Cost | Effect |
|---|---|---|
| **Nitrous** | 3 | After you burn fuel and before you roll, you may fire it: **+3** to this roll and **+1 Heat**. Then turn the card sideways in its slot (empty). It refills when you stop at a gas station, or you can pay 2 Rep at the start of your turn to refill it. |
| **Police Scanner** | 3 | Speed traps only catch you if you **land** on them. Passing through is always safe. |
| **Turbo** | 5 | **+1** to every roll. |
| **Racing Tires** | 6 | **+1** to every roll. Immune to oil slicks and Rain. |

If the shop runs out of an upgrade, nobody else can buy it.

## 8. Fuel

- You start with a **full tank (5)**.
- Each turn you drive, you **burn 1 fuel** before you roll. This includes Traffic Jam turns, when you move without rolling.
- **Fill up to 5** by stopping at a **Gas Station** or at **Shibuya Crossing**.
- Some cards add or remove fuel. Fuel never goes above 5 or below 0.

### Out of Gas
If your gauge is at **0** when you need to burn fuel, you're out of gas:

1. Push your car **back to the nearest Gas Station behind you** on your route. **Start** and **Shibuya Crossing** count as gas stations for this.
2. Fill your fuel to 5 and resolve the pit stop.
3. That's your whole turn. You don't roll.

> **Tip:** Watch your gauge! At 1 or 2 fuel, stop at the next gas station you pass, even if it costs you a few spaces.

## 9. Heat and the police

Your **Heat dial** goes from 0 to 5 and tracks how much attention you've drawn.

**Gaining Heat:** rolling a natural 6, firing Nitrous, and some cards and landmarks.

**Losing Heat:** gas stations (−2), getting pulled over (resets to 0), and some cards and landmarks.

### Speed traps
You're **caught** by a speed trap if either:
- you **land** on it, or
- you **pass through** it with **movement of 7 or more** that turn.

If you're caught while passing through, stop on the trap. Your movement ends there.

### Pulled over
- **At Heat 0, you get a warning.** No fine, nothing happens. Careful drivers are left alone.
- **At Heat 1 or more:**
  1. Pay Rep equal to your Heat.
  2. Reset your Heat to **0**.
  3. If you can't pay the full amount, pay everything you have and **park** your car (you lose your next turn).

### Busted
If your Heat ever hits the **red line (5)**, you're **Busted** right away:
- **Park** your car (you lose your next turn). If you're in the middle of a move, finish it first.
- Set your Heat to **2**.

*Bribe a Cop* (a Boost card) can cancel being pulled over or Busted.

## 10. Parking (losing a turn)

When something tells you to **lose your next turn**, lay your car on its side. At the start of your next turn, stand it back up. That's your whole turn. You don't burn fuel.

Things that park you:
- getting Busted
- not being able to pay a fine
- the *Tow Zone* card
- choosing to fix a Flat Tire on the spot

## 11. Interaction

### Slipstream
If you **start your turn** on the space directly behind another car on the same route, you get **+1** to your roll.

### Street Duel
If you end your move on a space with another car (except Start and Shibuya Crossing):
- Both players roll the die. Upgrades don't count. Reroll ties.
- The loser gives the winner **2 Rep** (or everything they have, if less than 2).
- If there are several cars on the space, the moving player picks one to duel.

### Narrow streets
Some Back Alley spaces are **narrow** (shown with dashed outlines). Only one car fits on a narrow space.
- You **can** drive through a narrow space that has a car on it.
- If your move would **end** on an occupied narrow space, stop on the space before it instead.
- There are no duels on narrow spaces.
- Cars moved **backward** (by Out of Gas, a card or an oil slick) ignore this rule.

### Oil slicks
Some Boost cards let you place an oil slick token on any empty Plain space. The next car to **land** on it spins out (see §6).

## 12. Shibuya Crossing

- **Every car must stop here.** Leftover movement is lost.
- **Fill your fuel to 5.** Every car leaves Shibuya with a full tank.
- No duels here. It's crowded, and everyone is waiting for the light.
- On your next turn, choose your leg 2 route: **Shuto C1 Loop** or **Harajuku–Roppongi**.

## 13. Cards

### Boost cards
- When you draw one, keep it in your hand. You can hold up to **3** (discard down to 3 if you go over).
- Each card says when you can play it:
  - **Before roll**
  - **After move**
  - **Instant**: any time, even on another player's turn
- Discard cards after you play them.

### Obstacle cards
Resolve them **immediately** when you draw them, then discard them. Cards with a lasting effect (like *Flat Tire*) stay face up next to your dashboard until the effect ends.

### Empty decks
If a deck runs out, shuffle its discard pile to make a new deck.

## 14. Optional variants

- **Street Cred mode:** the first car to the tower gets +10 Rep, and then the game ends. The most Rep wins. This rewards the Back Alleys more.
- **Short race:** skip leg 1. Start everyone at Shibuya Crossing with 8 Rep. Plays in about 15 minutes.
- **Hot pursuit:** start everyone at 1 Heat.
- **Fuel delivery:** when you're Out of Gas, you may instead pay 3 Rep to have fuel delivered. Fill to 5 where you are, and park (lose your next turn).
- **Thirsty nitrous:** firing Nitrous also burns 1 extra fuel.

## 15. Changes in v0.3

These come from 20,000-game balance simulations (see [SIMULATION.md](SIMULATION.md)).

| Change | Why |
|---|---|
| Back Alley +5 spaces are now +3 | The alleys were too rich: alley drivers won about 42% of games. |
| Speed traps catch cars passing at **7+** movement (was 6+) | The Wangan Expressway was too punishing: Expressway → C1 drivers won about 9% of games. |
| **Heat 0 = warning**, no fine | Careful but broke drivers were losing turns to 1-Rep fines. |
| 3rd and 4th players start with +1 / +2 Rep | Seat 1 was winning 30% of games and seat 4 only 20%. |
| No bonus for the first car into Shibuya | The halfway leader was winning 64% of games. |
| *Shortcut* moves 6 spaces at most | It was the strongest card. |
| Nitrous fires **after** you burn fuel | So you can't waste Nitrous on a turn you run out of gas. |
| Clarified: Busted mid-move, Traffic Jam fuel, re-choosing routes | Gaps found while writing the simulator. |
