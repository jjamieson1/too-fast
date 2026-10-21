"""Draws dashboard-sketch.svg, the player dashboard card.
Re-run after editing:  python3 tools/make_dashboard.py
"""
import math, sys

OUT = sys.argv[1] if len(sys.argv) > 1 else "dashboard-sketch.svg"
W, H = 920, 560
START, SWEEP = 225, 270  # dial runs from bottom-left, clockwise, to bottom-right


def pt(cx, cy, r, deg):
    a = math.radians(deg)
    return cx + r * math.cos(a), cy - r * math.sin(a)


def dial(cx, cy, title, values, colors, needle_at, caption):
    out = [f'<circle cx="{cx}" cy="{cy}" r="150" fill="#0b0d11" stroke="#5a6070" stroke-width="6"/>']
    n = len(values)
    for i, (v, c) in enumerate(zip(values, colors)):
        a0 = START - SWEEP * i / n
        a1 = START - SWEEP * (i + 1) / n
        x0, y0 = pt(cx, cy, 128, a0 - 2)
        x1, y1 = pt(cx, cy, 128, a1 + 2)
        out.append(f'<path d="M{x0:.1f},{y0:.1f} A128,128 0 0,1 {x1:.1f},{y1:.1f}" fill="none" stroke="{c}" stroke-width="18"/>')
        tx, ty = pt(cx, cy, 92, (a0 + a1) / 2)
        out.append(f'<text x="{tx:.1f}" y="{ty + 9:.1f}" font-size="26" font-weight="800" text-anchor="middle" fill="#f2f2f2">{v}</text>')
    i = values.index(needle_at)
    na = START - SWEEP * (i + 0.5) / n
    nx, ny = pt(cx, cy, 66, na)
    out.append(f'<line x1="{cx}" y1="{cy}" x2="{nx:.1f}" y2="{ny:.1f}" stroke="#ff5a36" stroke-width="6" stroke-linecap="round"/>')
    out.append(f'<circle cx="{cx}" cy="{cy}" r="14" fill="#ccc" stroke="#333" stroke-width="3"/>')
    out.append(f'<text x="{cx}" y="{cy + 58}" font-size="22" font-weight="900" text-anchor="middle" fill="#fff" letter-spacing="2">{title}</text>')
    out.append(f'<text x="{cx}" y="{cy + 82}" font-size="12" text-anchor="middle" fill="#aab">{caption}</text>')
    return out


svg = [f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {W} {H}" width="{W}" height="{H}" font-family="Helvetica, Arial, sans-serif">',
       f'<rect width="{W}" height="{H}" fill="#14171d"/>',
       f'<rect x="10" y="10" width="{W - 20}" height="{H - 20}" rx="36" fill="#22262f" stroke="#444a57" stroke-width="3"/>',
       '<text x="460" y="52" font-size="26" font-weight="900" font-style="italic" text-anchor="middle" fill="#fff">TOO FAST</text>',
       '<text x="460" y="72" font-size="12" text-anchor="middle" fill="#9aa0aa">PLAYER DASHBOARD · sketch v0.1</text>']

# Fuel: E (0) on the left, F (5) on the right, like a real gauge. Needle starts at 5.
svg += dial(210, 220, "FUEL", ["0", "1", "2", "3", "4", "5"],
            ["#d63131", "#f07f1f", "#f2c14e", "#9ccc65", "#43a047", "#43a047"], "5",
            "−1 each turn")
# Heat: tachometer with a red line at 5. Needle starts at 0.
svg += dial(710, 220, "HEAT", ["0", "1", "2", "3", "4", "5"],
            ["#5a6070", "#5a6070", "#8a8f9c", "#f2c14e", "#f07f1f", "#d63131"], "0",
            "red line = BUSTED")

# Centre reminder panel
svg.append('<rect x="378" y="110" width="164" height="210" rx="12" fill="#0b0d11" stroke="#444a57"/>')
for i, line in enumerate(["EACH TURN", "1 · Fuel −1", "   (at 0: back to", "   last gas station)", "2 · Roll + upgrades", "3 · Natural 6:", "   Heat +1", "4 · Natural 1:", "   Rep +1"]):
    bold = 'font-weight="800" fill="#fff"' if i == 0 else 'fill="#cfd3dc"'
    x = 408 if line.startswith(" ") else 392
    svg.append(f'<text x="{x}" y="{136 + i * 21}" font-size="13" {bold}>{line.strip()}</text>')

# Upgrade slots
slots = [("NITROUS", "3 Rep", "+3 once · +1 Heat"), ("POLICE SCANNER", "3 Rep", "Traps only if you land"),
         ("TURBO", "5 Rep", "+1 every roll"), ("RACING TIRES", "6 Rep", "+1 · no slicks / rain")]
for i, (name, cost, eff) in enumerate(slots):
    x = 40 + i * 215
    svg.append(f'<rect x="{x}" y="380" width="195" height="150" rx="10" fill="none" stroke="#7a8090" stroke-width="2.5" stroke-dasharray="8,6"/>')
    svg.append(f'<text x="{x + 97}" y="430" font-size="16" font-weight="800" text-anchor="middle" fill="#e8e8e8">{name}</text>')
    svg.append(f'<text x="{x + 97}" y="455" font-size="13" text-anchor="middle" fill="#f2c14e">{cost}</text>')
    svg.append(f'<text x="{x + 97}" y="478" font-size="12" text-anchor="middle" fill="#aab">{eff}</text>')
    svg.append(f'<text x="{x + 97}" y="515" font-size="10" text-anchor="middle" fill="#6b7180">upgrade card slot</text>')

svg.append('</svg>')
open(OUT, "w").write("".join(svg))
