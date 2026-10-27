"""Draws the README images: promo-banner.svg and cards-sketch.svg.
Re-run after editing:  python3 tools/make_promo.py
"""
import random

FONT = 'font-family="Helvetica, Arial, sans-serif"'


def esc(s):
    return s.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;")


def wrap(text, width):
    lines, cur = [], ""
    for w in text.split():
        if len(cur) + len(w) + 1 > width and cur:
            lines.append(cur)
            cur = w
        else:
            cur = f"{cur} {w}".strip()
    if cur:
        lines.append(cur)
    return lines


# ── banner ──────────────────────────────────────────────────────────
def banner():
    W, H = 1200, 400
    rnd = random.Random(7)
    o = [f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {W} {H}" width="{W}" height="{H}" {FONT}>',
         '<defs>'
         '<linearGradient id="sky" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#0b0820"/><stop offset="1" stop-color="#2a0f3d"/></linearGradient>'
         '<linearGradient id="road" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#1a1a24"/><stop offset="1" stop-color="#0d0d12"/></linearGradient>'
         '<linearGradient id="fade" x1="0" y1="0" x2="1" y2="0"><stop offset="0" stop-color="#0b0820" stop-opacity="0.95"/><stop offset="0.55" stop-color="#0b0820" stop-opacity="0.85"/><stop offset="0.75" stop-color="#0b0820" stop-opacity="0"/></linearGradient>'
         '<filter id="glow" x="-20%" y="-20%" width="140%" height="140%"><feGaussianBlur stdDeviation="4" result="b"/>'
         '<feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge></filter>'
         '</defs>',
         f'<rect width="{W}" height="{H}" fill="url(#sky)"/>']
    # skyline
    x = 0
    while x < W:
        w = rnd.randint(30, 70)
        h = rnd.randint(60, 190)
        o.append(f'<rect x="{x}" y="{290 - h}" width="{w}" height="{h}" fill="#150f2a"/>')
        for wy in range(290 - h + 10, 280, 16):
            for wx in range(x + 6, x + w - 6, 12):
                if rnd.random() < 0.35:
                    c = rnd.choice(["#ffd36b", "#7fe3ff", "#ff6fae"])
                    o.append(f'<rect x="{wx}" y="{wy}" width="5" height="7" fill="{c}" opacity="0.75"/>')
        x += w + rnd.randint(2, 8)
    # Tokyo Tower
    tx = 980
    o.append(f'<g stroke="#ff5a36" stroke-width="4" fill="none" filter="url(#glow)">'
             f'<path d="M{tx-60},290 L{tx-8},60 L{tx},20 L{tx+8},60 L{tx+60},290"/>'
             f'<path d="M{tx-40},200 L{tx+40},200 M{tx-26},140 L{tx+26},140 M{tx-16},95 L{tx+16},95"/>'
             f'<path d="M{tx-50},245 L{tx+50},245"/></g>')
    # dark fade behind the title so the text stays readable
    o.append(f'<rect width="{W}" height="290" fill="url(#fade)"/>')
    # road + speed lines
    o.append(f'<rect y="290" width="{W}" height="{H-290}" fill="url(#road)"/>')
    o.append(f'<line x1="0" y1="292" x2="{W}" y2="292" stroke="#ff2fa0" stroke-width="3" filter="url(#glow)"/>')
    for i in range(14):
        y = rnd.randint(305, 390)
        x1 = rnd.randint(-100, W)
        o.append(f'<line x1="{x1}" y1="{y}" x2="{x1 + rnd.randint(80, 260)}" y2="{y}" stroke="#7fe3ff" stroke-width="{rnd.choice([1, 2, 2, 3])}" opacity="{rnd.uniform(0.25, 0.7):.2f}"/>')
    for x1 in range(40, W, 140):
        o.append(f'<rect x="{x1}" y="343" width="70" height="5" fill="#f2c14e" opacity="0.8"/>')
    # title
    o.append('<text x="70" y="175" font-size="118" font-weight="900" font-style="italic" fill="#ffffff" filter="url(#glow)" letter-spacing="4">TOO FAST</text>')
    o.append('<text x="76" y="225" font-size="28" font-weight="700" fill="#7fe3ff" letter-spacing="2">GO FAST. GET FAMOUS. DON\'T GET CAUGHT.</text>')
    o.append('<text x="78" y="262" font-size="18" fill="#d6c8ff">A street-racing board game through Tokyo · 2–4 players · Ages 10+</text>')
    o.append('</svg>')
    open("promo-banner.svg", "w").write("".join(o))


# ── sample cards ────────────────────────────────────────────────────
CARDS = [
    ("BOOST", "#2e9d4f", "Bribe a Cop", "Instant", "Cancel being Pulled Over or Busted. Set your Heat to 1.", "“I didn't see anything.”"),
    ("BOOST", "#2e9d4f", "Oil Slick", "Instant", "Place an oil slick token on any empty Plain space.", "Oops. Was that me?"),
    ("BOOST", "#2e9d4f", "Empty Street", "After move", "Move forward 3 more spaces, then resolve that space.", "3 a.m. Not a soul in sight."),
    ("OBSTACLE", "#e0701a", "Flat Tire", "Resolve now", "−2 to your rolls until you stop at a Gas Station. Or park for a turn to fix it now.", "Pssssssss…"),
    ("OBSTACLE", "#e0701a", "Police Checkpoint", "Resolve now", "Heat 1+: you're Pulled Over. Heat 0: gain 1 Rep (clean record).", "License and registration, please."),
    ("UPGRADE", "#8e44ad", "Nitrous", "Cost: 3 Rep", "After burning fuel, before rolling: +3 to this roll, +1 Heat. Refills at gas stations.", "Push the red button."),
]


def card(x, y, kind, col, name, timing, effect, flavor):
    w, h = 220, 320
    o = [f'<g transform="translate({x},{y})">',
         f'<rect width="{w}" height="{h}" rx="16" fill="#1b1f27" stroke="{col}" stroke-width="4"/>',
         f'<rect x="10" y="10" width="{w-20}" height="44" rx="8" fill="{col}"/>',
         f'<text x="{w/2}" y="39" font-size="16" font-weight="900" text-anchor="middle" fill="#fff" letter-spacing="3">{kind}</text>',
         f'<text x="{w/2}" y="92" font-size="21" font-weight="800" text-anchor="middle" fill="#fff">{esc(name)}</text>',
         f'<text x="{w/2}" y="116" font-size="12" text-anchor="middle" fill="{col}" font-weight="700">{esc(timing.upper())}</text>',
         f'<line x1="24" y1="132" x2="{w-24}" y2="132" stroke="#3a3f4b"/>']
    for i, line in enumerate(wrap(effect, 24)):
        o.append(f'<text x="{w/2}" y="{162 + i*21}" font-size="15" text-anchor="middle" fill="#e8e8e8">{esc(line)}</text>')
    for i, line in enumerate(wrap(flavor, 28)):
        o.append(f'<text x="{w/2}" y="{272 + i*18}" font-size="13" font-style="italic" text-anchor="middle" fill="#9aa0aa">{esc(line)}</text>')
    o.append('</g>')
    return "".join(o)


def cards():
    W, H = 3 * 220 + 4 * 30, 2 * 320 + 3 * 30
    o = [f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {W} {H}" width="{W}" height="{H}" {FONT}>',
         f'<rect width="{W}" height="{H}" fill="#14171d"/>']
    for i, c in enumerate(CARDS):
        o.append(card(30 + (i % 3) * 250, 30 + (i // 3) * 350, *c))
    o.append('</svg>')
    open("cards-sketch.svg", "w").write("".join(o))


banner()
cards()
