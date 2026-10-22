import math, sys

"""Draws board-sketch.svg from the route data below.
Edit a route's `seq` (one code per space) and re-run:  python3 tools/make_board.py
"""
OUT = sys.argv[1] if len(sys.argv) > 1 else "board-sketch.svg"

# Space codes: . plain, R +3 rep, R5 +5 rep, B boost, O obstacle, G gas, S speed trap, L:<name> landmark
ROUTES = {
    "A": dict(name="Wangan Expressway", seq=". R S . B . L:Daikoku_PA . S O . G . S R . L:Rainbow_Bridge . O S . B . S",
              pts=[(70,450),(110,290),(220,150),(400,105),(560,140),(620,280),(720,450)], side=-1, narrow=()),
    "B": dict(name="Ginza Avenue", seq=". B R . O . S R . L:Ginza B . G O . R . B S . L:Akihabara O R . B . R .",
              pts=[(70,450),(160,450),(225,345),(315,335),(380,450),(440,560),(530,565),(580,460),(640,450),(720,450)], side=-1, narrow=()),
    "C": dict(name="Back Alleys", seq=". R B R . O R L:Tsukiji B . G O B R . G R B . R O L:Ameyoko B G . R B O R . B R",
              pts=[(70,450),(95,600),(170,730),(290,810),(420,830),(530,780),(580,690),(620,600),(720,450)], side=1, narrow=tuple(range(5,13))+tuple(range(19,27))),
    "D": dict(name="Shuto C1 Loop", seq=". S . B S . O R S . G . S B . O S .",
              pts=[(720,450),(780,300),(920,215),(1100,215),(1250,290),(1350,450)], side=-1, narrow=()),
    "E": dict(name="Harajuku–Roppongi", seq=". R B . L:Harajuku R O . B S R . G B L:Roppongi O R . B R O .",
              pts=[(720,450),(770,610),(900,720),(1080,730),(1240,640),(1350,450)], side=1, narrow=()),
}
EXPECTED = {"A":24,"B":28,"C":32,"D":18,"E":22}

COL = {".":"#d6d6d6","R":"#f2c14e","R5":"#e0951f","B":"#43a047","O":"#f07f1f","G":"#2f7bd8","S":"#d63131","L":"#8e44ad"}
TXT = {".":"","R":"+3","R5":"+5","B":"B","O":"!","G":"G","S":"S","L":"★"}

def catmull(pts, n=400):
    P = [pts[0]] + pts + [pts[-1]]
    out = []
    for i in range(1, len(P)-2):
        p0,p1,p2,p3 = P[i-1],P[i],P[i+1],P[i+2]
        for k in range(n):
            t = k/n
            t2,t3 = t*t,t*t*t
            out.append(tuple(0.5*((2*p1[j]) + (-p0[j]+p2[j])*t + (2*p0[j]-5*p1[j]+4*p2[j]-p3[j])*t2 + (-p0[j]+3*p1[j]-3*p2[j]+p3[j])*t3) for j in (0,1)))
    out.append(pts[-1])
    return out

def resample(poly, count):
    # count spaces strictly between the two endpoints
    d = [0]
    for a,b in zip(poly, poly[1:]):
        d.append(d[-1] + math.dist(a,b))
    L = d[-1]
    res = []
    j = 0
    for i in range(1, count+1):
        M = 48
        target = M + (L-2*M)*(i-0.5)/count
        while d[j+1] < target: j += 1
        a,b = poly[j],poly[j+1]
        f = (target-d[j])/((d[j+1]-d[j]) or 1)
        p = (a[0]+(b[0]-a[0])*f, a[1]+(b[1]-a[1])*f)
        tx,ty = b[0]-a[0], b[1]-a[1]
        m = math.hypot(tx,ty) or 1
        res.append((p, (-ty/m, tx/m)))
    return res, (L-96)/count

R = 12
svg = []
allpos = {}
for key, r in ROUTES.items():
    seq = r["seq"].split()
    assert len(seq) == EXPECTED[key], (key, len(seq))
    poly = catmull(r["pts"])
    pos, spacing = resample(poly, len(seq))
    allpos[key] = [p for p,_ in pos]
    print(f"{key}: {len(seq)} spaces, spacing {spacing:.1f}", file=sys.stderr)
    d = "M" + " L".join(f"{x:.1f},{y:.1f}" for x,y in poly[::8]) + f" L{poly[-1][0]},{poly[-1][1]}"
    svg.append(f'<path d="{d}" fill="none" stroke="#3a3f4b" stroke-width="30" stroke-linecap="round" stroke-linejoin="round"/>')
    r["_pos"], r["_seq"] = pos, seq

# overlap check between routes
mind = 1e9
keys = list(allpos)
for i,a in enumerate(keys):
    for b in keys[i+1:]:
        for p in allpos[a]:
            for q in allpos[b]:
                mind = min(mind, math.dist(p,q))
print(f"min cross-route distance {mind:.1f}", file=sys.stderr)

for key, r in ROUTES.items():
    for idx, ((p, nrm), code) in enumerate(zip(r["_pos"], r["_seq"]), start=1):
        x,y = p
        label = None
        if code.startswith("L:"):
            label = code[2:].replace("_"," "); code = "L"
        dash = ' stroke-dasharray="3,2"' if idx in r["narrow"] else ""
        sw = 2.2 if idx in r["narrow"] else 1.5
        svg.append(f'<circle cx="{x:.1f}" cy="{y:.1f}" r="{R}" fill="{COL[code]}" stroke="#111"{dash} stroke-width="{sw}"/>')
        t = TXT[code]
        if t:
            fc = "#111" if code in ("R","R5",".") else "#fff"
            svg.append(f'<text x="{x:.1f}" y="{y+4:.1f}" font-size="11" font-weight="700" text-anchor="middle" fill="{fc}">{t}</text>')
        s = r["side"]
        # number outside, small
        nx, ny = nrm
        if (ny > 0) != (s > 0): nx, ny = -nx, -ny
        svg.append(f'<text x="{x+nx*21:.1f}" y="{y+ny*21+3:.1f}" font-size="8" text-anchor="middle" fill="#9aa0aa">{idx}</text>')
        if label:
            svg.append(f'<text x="{x+nx*42:.1f}" y="{y+ny*42+4:.1f}" font-size="13" font-weight="700" text-anchor="middle" fill="#c58cf0">{label}</text>')

# route name labels at midpoint
for key, r in ROUTES.items():
    x,y = {"A":(400,60),"B":(330,610),"C":(330,885),"D":(1035,165),"E":(1035,800)}[key]
    svg.append(f'<text x="{x}" y="{y}" font-size="15" font-style="italic" text-anchor="middle" fill="#e8e8e8">{key} · {r["name"]} ({len(r["_seq"])})</text>')

def node(x,y,title,sub,col):
    return (f'<circle cx="{x}" cy="{y}" r="34" fill="{col}" stroke="#fff" stroke-width="3"/>'
            f'<text x="{x}" y="{y+1}" font-size="12" font-weight="800" text-anchor="middle" fill="#fff">{title}</text>'
            f'<text x="{x}" y="{y+15}" font-size="9" text-anchor="middle" fill="#f0f0f0">{sub}</text>')

svg.append(node(70,450,"START","Odaiba","#222"))
svg.append(node(720,450,"SHIBUYA","Must stop · Refuel","#c2185b"))
svg.append(node(1350,450,"FINISH","Tokyo Tower","#e53935"))

legend = [(".","Plain"),("R","+3 Rep"),("B","Boost card"),("O","Obstacle card"),("G","Gas station"),("S","Speed trap"),("L","Landmark")]
lx, ly = 1010, 40
svg.append(f'<rect x="{lx-15}" y="{ly-20}" width="380" height="120" rx="10" fill="#1d2129" stroke="#444"/>')
for i,(c,name) in enumerate(legend):
    x = lx + (i%2)*190; y = ly + (i//2)*24
    svg.append(f'<circle cx="{x}" cy="{y}" r="9" fill="{COL[c]}" stroke="#111"/>')
    svg.append(f'<text x="{x+16}" y="{y+4}" font-size="12" fill="#e8e8e8">{name}</text>')
svg.append(f'<circle cx="{lx}" cy="{ly+96}" r="9" fill="#d6d6d6" stroke="#111" stroke-dasharray="3,2" stroke-width="2"/>')
svg.append(f'<text x="{lx+16}" y="{ly+100}" font-size="12" fill="#e8e8e8">Narrow street (one car per space)</text>')

W,H = 1420, 900
head = (f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {W} {H}" width="{W}" height="{H}" font-family="Helvetica, Arial, sans-serif">'
        f'<rect width="{W}" height="{H}" fill="#14171d"/>'
        f'<text x="40" y="55" font-size="40" font-weight="900" fill="#fff" font-style="italic">TOO FAST</text>'
        f'<text x="42" y="80" font-size="14" fill="#9aa0aa">Board sketch v0.3 · Tokyo (stylized, not to scale)</text>')
open(OUT,"w").write(head + "".join(svg) + "</svg>")

