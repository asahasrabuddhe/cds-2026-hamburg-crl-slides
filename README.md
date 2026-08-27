# The Reality of Rootless Containers

Slides for **The Reality of Rootless Containers**, ContainerDays 2026.
Thirty minutes, plus five for questions.

Built with [Slidev](https://sli.dev) and the private
`slidev-theme-architectural-console` theme. The talk's code lives in a sibling
repository,
[`cds-2026-hamburg-crl-code`](https://github.com/asahasrabuddhe/cds-2026-hamburg-crl-code),
checked out here as a git submodule at `vendor/crl-code`.

## Setup

```bash
git clone --recurse-submodules git@github.com:asahasrabuddhe/cds-2026-hamburg-crl-slides.git
cd cds-2026-hamburg-crl-slides
nvm use                       # Node 24
cp .env.example .env          # fills in from 1Password, see below
op run --env-file=.env -- npm install
```

The theme is a private npm package. `NPM_TOKEN_THEME` in `.env` is an `op://`
reference resolved by `op run`, never a literal token on disk. `.npmrc` pins
the registry to `registry.npmjs.org` explicitly, because this machine's global
npm config routes elsewhere.

## Commands

```bash
npm run dev            # slidev --open
npm run build          # static build
npm run export         # PDF, every click state
npm run export:notes   # PDF with presenter notes
npm run snippets       # pull Go snippets out of vendor/crl-code
```

`npm run export` passes `--with-clicks`, so each click state becomes its own
page. Twelve slides use click reveals, so the exported PDF is longer than the
slide count, and the organisers get every state rather than only the final
one.

`@slidev/cli` is pinned to exactly `52.14.1` rather than a caret range. A deck
is a fixed artefact delivered on a fixed date: a minor version arriving the
week of the talk can change layout metrics or export behaviour, and there is
no upside to finding that out on the morning. Bump it deliberately, after an
export, or not at all.

## The crossing file

Slide 22 shows the `SysProcAttr` literal that creates the user namespace. That
code is not written here. It lives in the companion repository and is pulled
in by anchor comment, so the deck cannot drift from the program that runs on
stage.

```bash
npm run snippets       # reads vendor/crl-code, writes snippets/*.go
```

`cmd/snippets` scans `slides.md` for `<<< @/snippets/<name>.go` imports, finds
the matching `// snippet:start <name>` region in the submodule, strips the
common indentation, and writes the file. It fails loudly if the deck
references an anchor that no longer exists, because a silently missing snippet
is worse than a broken build.

Before the submodule exists, point it at a sibling checkout:

```bash
CRL_CODE_DIR=../cds-2026-hamburg-crl-code npm run snippets
```

There is no sync script beyond this, and there should not be. The submodule
pin is the sync: it records exactly which commit of the code the deck was
built against.

## Presenter notes

Notes are HTML comments at the end of each slide in `slides.md`. Press `p` in
the dev server for presenter mode, which puts them on your laptop screen while
the projector keeps the slide. `npm run export:notes` bakes them into a PDF
for rehearsal.

## Stage layout, and the cold open

The deck is not the whole screen. It occupies the top two thirds; a tmux
session fills the bottom third and stays there for the entire talk. Never
switch windows.

That matters most at the very start, because **the talk opens in the terminal,
not on a slide.** The cold open runs for about ninety seconds with no slide up
at all: two panes, `id` inside a rootless container saying `uid=0(root)` on one
side and `/proc/<pid>/status` on the host saying `Uid: 1000` on the other.
Slide 1 comes up afterwards, and slide 2, "Both answers are true", is the
thesis line landing on the two answers the audience has just watched.

Slide 2 is load-bearing on the cold open. Skip the cold open and it reads as a
non-sequitur.

Panes are set up per the design: left labelled `ROOTFUL (uid 0)` in red, right
`ROOTLESS (uid 1000)` in green, font at 20pt minimum, prompt stripped to `#`
and `$`. `scripts/stage.sh` in the companion repository builds that session.
Demo A is the exception: a single full-screen pane, not split.

## On-stage timing checkpoints

Two, and they are the whole timing discipline.

| Time | Checkpoint | If you are late |
|---|---|---|
| 11:00 | Demo A starts | Cut the condo analogy, slide 7 |
| 22:00 | Out of the demo block | Drop Demo 5, go straight to the audit at slide 33 |

Everything after 22:00 is the part people quote. Protect it.

## Slides that are deliberately incomplete

| Slide | What it needs | Note |
|---|---|---|
| 7 | A photograph for the condo analogy | The layout has **zero** vertical headroom as it stands, so the photo cannot simply be added. Something has to give, most likely shortening the left column's three paragraphs to two |
| 36 | Refresh against the measured AppArmor behaviour | Ubuntu's restriction turned out to be an allowlist rather than a switch, and Podman still works under it. See `docs/environment.md` in the code repository |

Slide 32's benchmark numbers are filled in, measured on the demo VM. Read the
caveat on the slide: it is a VM, so the ratios are the finding, not the
absolute figures.

## Rendering

Checked at 16:9, 1280x720, across all 40 slides and every click state. No
slide overflows. Three fixes were needed to get there and they live in
`style.css`: the theme centres table headers but left-aligns table cells, its
table cell padding is 1px which runs adjacent columns together, and slide 11's
mermaid diagram needed its scale dropped from 0.8 to 0.58.

## Licence

CC BY-SA 4.0. See [`LICENSE`](LICENSE).
