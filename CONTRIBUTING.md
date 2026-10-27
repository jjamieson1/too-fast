# Contributing to Too Fast

Thanks for playing, testing or changing *Too Fast*! 🏁

The game is free to use under [CC BY-NC-SA 4.0](LICENSE). The license requires that any changed version you share stays under the same license, with credit to this project. **On top of that, we ask that you send your changes back here**, so every player gets the best version of the game. The license doesn't require it, but it's how the game gets better.

## Ways to help

### 🎲 Playtest and tell us how it went
This is the most useful thing you can do. Print the board and cards (see [BOARD.md](BOARD.md) and [CARDS.md](CARDS.md)), play a race, and [open an issue](https://github.com/jjamieson1/too-fast/issues) with:

- How many players, and how long the game took
- Who won, and which routes and upgrades they used
- Anything confusing, unfair or un-fun
- Any moment that made the table cheer or groan

The open questions we most want answered are listed in [IDEA.md](IDEA.md#open-questions-for-playtesting).

### 💡 Suggest a change
Got an idea for a new card, landmark or rule? [Open an issue](https://github.com/jjamieson1/too-fast/issues) describing it and why it would make the game better.

### 🔧 Send a change
1. Fork the repo and create a branch.
2. Make your change. Keep the docs in sync:
   - **Rules:** [RULES.md](RULES.md)
   - **Cards:** [CARDS.md](CARDS.md), and the deck counts in `game/cards.go`
   - **Board:** the `seq` lines in [tools/make_board.py](tools/make_board.py) **and** `game/board.go`, then regenerate the sketch with `python3 tools/make_board.py`
3. If your change affects balance, run the simulator and include the before and after numbers in your pull request:

   ```bash
   go test ./...
   ```

   ```bash
   go run ./cmd/sim
   ```

   See [SIMULATION.md](SIMULATION.md) for how to read the report.
4. Open a pull request explaining what changed and why.

## Licensing of contributions

By sending a contribution (a pull request, card idea, artwork or rule text), you agree that it's licensed under the project's [CC BY-NC-SA 4.0](LICENSE) license. You also agree that the project owner may include it in any future official edition of *Too Fast*, including one released under different terms.

## Be kind

Be respectful in issues and pull requests. Critique the game, not the people. Everyone here is trying to make a better race.
