#!/usr/bin/env python3
"""สร้าง golang-interview.html จาก PROBLEMS.md, SOLUTIONS.md, THEORY.md

ไม่ใช้ library ภายนอก: มี markdown converter ขนาดเล็ก (รองรับเฉพาะ syntax ที่ใช้ในไฟล์เหล่านี้)
และ syntax highlighter สำหรับ Go ในตัว

    python3 tools/build_html.py                  # → golang-interview.html
    python3 tools/build_html.py --fragment out.html  # ไม่มี <html>/<head>/<body> ครอบ
"""
import html
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# ---------------------------------------------------------------- Go highlighter

GO_KEYWORDS = set("""break case chan const continue default defer else fallthrough for func go goto
if import interface map package range return select struct switch type var""".split())
GO_BUILTINS = set("""append cap clear close complex copy delete imag len make max min new panic print
println real recover bool byte complex64 complex128 error float32 float64 int int8 int16 int32 int64
rune string uint uint8 uint16 uint32 uint64 uintptr any comparable true false nil iota""".split())

GO_TOKEN = re.compile(r"""
    (?P<comment>//[^\n]*|/\*.*?\*/)
  | (?P<string>"(?:\\.|[^"\\\n])*"|`[^`]*`|'(?:\\.|[^'\\\n])+')
  | (?P<number>\b(?:0[xX][0-9a-fA-F_]+|\d[\d_]*(?:\.\d+)?(?:[eE][+-]?\d+)?)\b)
  | (?P<ident>[A-Za-z_][A-Za-z0-9_]*)
""", re.S | re.X)


def highlight_go(code):
    out, pos = [], 0
    for m in GO_TOKEN.finditer(code):
        out.append(html.escape(code[pos:m.start()]))
        text = html.escape(m.group())
        kind = m.lastgroup
        if kind == "ident":
            if m.group() in GO_KEYWORDS:
                kind = "kw"
            elif m.group() in GO_BUILTINS:
                kind = "bi"
            elif code[m.end():m.end() + 1] == "(":
                kind = "fn"
            else:
                kind = None
        out.append(f'<span class="t-{kind}">{text}</span>' if kind else text)
        pos = m.end()
    out.append(html.escape(code[pos:]))
    return "".join(out)


def highlight_hash_comments(code):
    lines = []
    for line in code.split("\n"):
        m = re.search(r"(^|\s)(#.*)$", line)
        if m:
            i = m.start(2)
            lines.append(html.escape(line[:i]) + f'<span class="t-comment">{html.escape(line[i:])}</span>')
        else:
            lines.append(html.escape(line))
    return "\n".join(lines)


def code_block(code, lang):
    if lang == "go":
        body = highlight_go(code)
    elif lang in ("bash", "sh", "makefile"):
        body = highlight_hash_comments(code)
    else:
        body = html.escape(code)
    label = f'<span class="code-lang">{html.escape(lang)}</span>' if lang else ""
    return f'<div class="code">{label}<pre><code>{body}</code></pre></div>'


# ---------------------------------------------------------------- markdown

def slugify(text, used):
    base = re.sub(r"[^\w฀-๿]+", "-", text.lower()).strip("-") or "section"
    slug, n = base, 2
    while slug in used:
        slug, n = f"{base}-{n}", n + 1
    used.add(slug)
    return slug


def inline(text):
    # เก็บ code span ไว้ก่อน (แทนด้วย placeholder) เพื่อให้ **bold** ที่มี `code` ข้างในยังทำงาน
    spans = []

    def stash(m):
        raw = m.group()
        inner = raw[2:-2].strip() if raw.startswith("``") else raw[1:-1]
        spans.append(f"<code>{html.escape(inner)}</code>")
        return f"\x00{len(spans) - 1}\x00"

    s = re.sub(r"``.+?``|`[^`]+`", stash, text)
    s = html.escape(s, quote=False)
    s = re.sub(r"\*\*(.+?)\*\*", r"<strong>\1</strong>", s)

    def link(m):
        label, href = m.group(1), m.group(2)
        if href.startswith("http"):
            return f'<a href="{href}" target="_blank" rel="noopener">{label}</a>'
        if href.startswith("#"):
            return f'<a href="{href}">{label}</a>'
        return f'<span class="file-ref">{label}</span>'  # ไฟล์ในเครื่อง — แสดงเป็นชื่อไฟล์

    s = re.sub(r"\[([^\]]+)\]\(([^)]+)\)", link, s)
    return re.sub(r"\x00(\d+)\x00", lambda m: spans[int(m.group(1))], s)


LIST_RE = re.compile(r"^(\s*)([-*]|\d+\.)\s+(.*)$")


def render_list(lines):
    """lines: list ของบรรทัดใน list block (อาจซ้อนกันตาม indent)"""
    items = []  # (indent, ordered, text)
    for line in lines:
        m = LIST_RE.match(line)
        if m:
            items.append([len(m.group(1)), m.group(2)[-1] == ".", m.group(3)])
        elif items:
            items[-1][2] += " " + line.strip()

    def build(i, indent):
        ordered = items[i][1]
        tag = "ol" if ordered else "ul"
        out = [f"<{tag}>"]
        while i < len(items) and items[i][0] >= indent:
            if items[i][0] > indent:
                sub, i = build(i, items[i][0])
                out[-1] = out[-1][: -len("</li>")] + sub + "</li>"
                continue
            out.append(f"<li>{inline(items[i][2])}</li>")
            i += 1
        out.append(f"</{tag}>")
        return "".join(out), i

    html_out, _ = build(0, items[0][0])
    return html_out


def render_table(lines):
    rows = [[c.strip() for c in l.strip().strip("|").split("|")] for l in lines]
    head, body = rows[0], rows[2:]
    out = ['<div class="table-wrap"><table><thead><tr>']
    out += [f"<th>{inline(c)}</th>" for c in head]
    out.append("</tr></thead><tbody>")
    for r in body:
        out.append("<tr>" + "".join(f"<td>{inline(c)}</td>" for c in r) + "</tr>")
    out.append("</tbody></table></div>")
    return "".join(out)


def include_code(path):
    full = os.path.join(ROOT, path)
    with open(full, encoding="utf-8") as f:
        src = f.read().rstrip("\n")
    return (f'<details class="full-code"><summary>ดูโค้ดเต็ม · <code>{html.escape(path)}</code></summary>'
            f"{code_block(src, 'go')}</details>")


def markdown(md, prefix):
    lines = md.split("\n")
    out, toc, used = [], [], set()
    i = 0
    while i < len(lines):
        line = lines[i]
        s = line.strip()

        if s.startswith("```"):
            lang = s[3:].strip()
            j = i + 1
            while j < len(lines) and not lines[j].strip().startswith("```"):
                j += 1
            out.append(code_block("\n".join(lines[i + 1:j]), lang))
            i = j + 1
            continue

        m = re.match(r"^<!-- code: (.+?) -->$", s)
        if m:
            out.append(include_code(m.group(1)))
            i += 1
            continue

        if s.startswith("<"):  # raw HTML (details/summary)
            out.append(s)
            i += 1
            continue

        m = re.match(r"^(#{1,4})\s+(.*)$", line)
        if m:
            level, text = len(m.group(1)), m.group(2)
            slug = prefix + "-" + slugify(text, used)
            if level == 1:
                out.append(f'<h1 id="{slug}">{inline(text)}</h1>')
            else:
                out.append(f'<h{level} id="{slug}">{inline(text)}</h{level}>')
                if level <= 3:
                    toc.append((level, slug, re.sub(r"`", "", text)))
            i += 1
            continue

        if s == "---":
            out.append("<hr>")
            i += 1
            continue

        if s.startswith("|"):
            j = i
            while j < len(lines) and lines[j].strip().startswith("|"):
                j += 1
            out.append(render_table(lines[i:j]))
            i = j
            continue

        if s.startswith(">"):
            j = i
            buf = []
            while j < len(lines) and lines[j].strip().startswith(">"):
                buf.append(lines[j].strip()[1:].strip())
                j += 1
            out.append(f"<blockquote><p>{'<br>'.join(inline(b) for b in buf)}</p></blockquote>")
            i = j
            continue

        if LIST_RE.match(line):
            j = i
            while j < len(lines) and lines[j].strip() and (LIST_RE.match(lines[j]) or lines[j].startswith(" ")):
                j += 1
            out.append(render_list(lines[i:j]))
            i = j
            continue

        if not s:
            i += 1
            continue

        j = i
        buf = []
        while j < len(lines):
            t = lines[j].strip()
            if not t or t.startswith(("```", "#", "|", ">", "<", "---")) or LIST_RE.match(lines[j]):
                break
            buf.append(t)
            j += 1
        i = j
        if buf[0].startswith("**โค้ดเต็ม:**"):
            continue  # ใน HTML มีปุ่ม "ดูโค้ดเต็ม" อยู่แล้ว
        out.append(f"<p>{inline(' '.join(buf))}</p>")

    return "\n".join(out), toc


# ---------------------------------------------------------------- page

def read(name):
    with open(os.path.join(ROOT, name), encoding="utf-8") as f:
        return f.read()


def toc_html(toc, skip_first_level=None):
    items = []
    for level, slug, text in toc:
        if level == 2:
            items.append(f'<li><a href="#{slug}">{html.escape(text)}</a></li>')
    return "<ol class=\"toc-list\">" + "".join(items) + "</ol>"


def build(fragment):
    tabs = [
        ("problems", "โจทย์", "PROBLEMS.md"),
        ("solutions", "เฉลย", "SOLUTIONS.md"),
        ("theory", "ทฤษฎี 50 ข้อ", "THEORY.md"),
    ]
    panels, navs, buttons = [], [], []
    for key, label, fname in tabs:
        body, toc = markdown(read(fname), key)
        panels.append(f'<article class="panel prose" id="panel-{key}" data-panel="{key}">{body}</article>')
        navs.append(f'<nav class="toc" data-toc="{key}" aria-label="สารบัญ {label}">'
                    f'<p class="toc-title">สารบัญ</p>{toc_html(toc)}</nav>')
        buttons.append(f'<button type="button" role="tab" class="tab" id="tab-{key}" data-tab="{key}" '
                       f'aria-controls="panel-{key}">{label}</button>')

    with open(os.path.join(ROOT, "tools", "template.html"), encoding="utf-8") as f:
        tpl = f.read()
    page = (tpl.replace("{{TABS}}", "\n".join(buttons))
               .replace("{{TOCS}}", "\n".join(navs))
               .replace("{{PANELS}}", "\n".join(panels)))
    if not fragment:
        page = ('<!doctype html>\n<html lang="th">\n<head>\n<meta charset="utf-8">\n'
                '<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">\n'
                + page.replace("<!--BODY-->", "</head>\n<body>", 1) + "\n</body>\n</html>\n")
    else:
        page = page.replace("<!--BODY-->", "", 1)
    return page


if __name__ == "__main__":
    fragment = "--fragment" in sys.argv
    out = sys.argv[sys.argv.index("--fragment") + 1] if fragment else os.path.join(ROOT, "golang-interview.html")
    with open(out, "w", encoding="utf-8") as f:
        f.write(build(fragment))
    print("wrote", out)
