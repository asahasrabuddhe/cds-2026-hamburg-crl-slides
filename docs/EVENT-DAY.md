# Event day

**The Reality of Rootless Containers**
ContainerDays 2026, Hamburg. 30 minutes plus 5 for questions.

This is the operational document. What you say is in [`SCRIPT.md`](SCRIPT.md)
and, as of this commit, also in the deck's own presenter notes, so you should
not need to open `SCRIPT.md` on stage at all. Design rationale is in
[`DESIGN.md`](DESIGN.md).

The companion repo is at
`~/Downloads/containerdays-2026/crl/cds-2026-hamburg-crl-code`, and this one is
its sibling. Every path below is relative to one of those two.

---

## Part 0: what is still outstanding

### Yours, before you travel

- **Decide about S07's photograph.** `public/images/` holds only a `.gitkeep`.
  The condo slide works as text and it is the first thing the script tells you
  to cut, so this is a decision rather than a blocker.

- **Confirm both repos are still private,** and remember that S40 links both of
  them. They are useless to the audience until you flip them, and that happens
  after the session, not before.

### What changed that affects what you type

- **`./rootless-demos.sh` is gone.** It was the pre-split harness and it never
  existed in this repo. Everything is `./scripts/demo.sh` now, and the deck no
  longer prints the old name at the audience on S26.
- **`make push` is a new step and it is not optional.** QEMU shares no
  filesystem with the host, cloud-init does not clone, and until now nothing at
  all put the code into the VM. `scripts/stage.sh` and every demo command
  expect the repository at `~/crl` in the guest.
- **Demo A works through the driver again.** `demo.sh a` looked for the binary
  at `./nsdemo/nsdemo`, one directory too deep, so it always printed "Not
  built". `demo.sh check` looked in the right place, so the pre-flight passed
  while the demo failed.
- **There are no recordings, and the script no longer pretends there are.**
  Demo 4's second half runs live on the `cgroupv1` VM and demo 5 runs live on
  `primary`. That is why two VMs come up at T-90 rather than one.
- **`npm run export:notes` exports notes.** It was passing `--with-toc`, which
  adds a PDF outline to the ordinary deck. Worth a run the night before now
  that all forty slides carry notes.
- **Slide numbers are the deck's, S01 to S40.** The script used to number to 38
  against a forty-slide deck and drifted by one from S28 onward. The notes and
  the script now agree.
- **Both repos are pushed and `vendor/crl-code` is current.** The pin matches the
  code repo's `HEAD`, so the deck ships against the code it describes. Nothing
  here is waiting on a push any more.
- **S23 was wrong about capabilities and is fixed.** It used to print
  `CapEff: 000001ffffffffff` for beat 1, which the program never produces. The
  kernel grants the full set at `clone`, then `execve` takes it back, because an
  unmapped process has euid 65534 and the exec counts as unprivileged. The slide
  now shows an empty `CapEff` against a full `CapBnd`, and the notes say the
  ceiling is unlimited and the holding is nothing.
- **S15 no longer claims `io` is delegated.** This image delegates `cpu memory
  pids`. Demo 4 uses `memory`, so it is unaffected, but do not claim `io` works
  rootless on this box.
- **S32's storage row is a size with no time.** The 12 seconds was measured
  against `golang:1.23` and `scripts/bench.sh` pulls `golang:1.27.0` now. The row
  is bandwidth-bound, so no single figure is defensible without naming the link.
  The 927 MB is a property of the image and stands.
- **The left pane runs Docker, not Podman.** `demo.sh` picks Docker whenever it
  is on `PATH` and the VM installs both. `demo.sh check` and `demo.sh 4` ask for
  Podman-only `--format` fields, so the left pane prints a Go template error
  before the useful output. It is cosmetic. `ENGINE=podman` suppresses it if the
  noise bothers you more than the inconsistency would.
- **The `hardened` VM has no `make`.** That variant omits the `make` and `git`
  packages, and `vm.sh push` copies git-tracked files only, so `nsdemo` does not
  cross either. If you intend to prove the S36 allowlist point live, build it in
  the guest with `go build -o nsdemo ./cmd/nsdemo`.

---

## Part 1: the night before

In the hotel, on mains power.

1. **Run the benchmarks, if you are going to re-run them at all.**
   `./scripts/bench.sh` pulls `networkstatic/iperf3` and `golang:1.27.0`, and
   those are the only two images cloud-init does not pre-pull. It is the one
   thing here that needs working internet inside the VM, so it does not happen
   at the venue. The iperf3 and cold-start rows on S32 are good. The storage
   row is not, see Part 0.

2. **Full dry run, both repos, start to finish.** Not a spot check. Every
   command in the script's own reference, in the order it lists them.

   This one is not optional any more. Demo 2, demo 4 and demo 5 all changed
   behaviour, so what you rehearsed before is not what the box does now. Demo 1
   and demo 3 in particular have not been run since the change.

3. **`./scripts/demo.sh all` twice, on `primary`.** The second run proves the
   demos are idempotent, which is the property you are relying on at 18:00 when
   demo 3 has half-died. Confirm `grep backdoor /etc/passwd` comes back empty
   after both.

4. **Export the notes and read them.** `npm run export:notes`. This is the
   first talk where the deck carries the whole script, so read it once as a
   document rather than clicking through it.

5. **Time yourself, twice.** The pacing card in Part 6 is derived from the two
   hard checkpoints, not measured, so every mark on it is a guess. Two timed
   run-throughs turn it into a real card. Move the numbers in the table and in
   the deck's notes together.

6. **Charge everything.** Laptop, clicker, phone.

7. **Every notification off.** Slack, Mail, Calendar, iMessage. Focus on,
   screen saver off, auto-lock off, auto-update off.

---

## Part 2: the morning, ninety minutes out

Ninety minutes is deliberate. The VM disks are warm, so a boot is a boot rather
than a re-provision, but cloud-init still has to settle and you want the
pre-flight clean with time to spare.

### T-90: bring both VMs up and push the code

```bash
cd ~/Downloads/containerdays-2026/crl/cds-2026-hamburg-crl-code
make vm                      # primary
make vm-v1                   # cgroupv1, demo 4's second half
make push
make push VARIANT=cgroupv1
```

`vm.sh` waits for cloud-init and prints `ready` per box. Two VMs, not one:
demo 4's second half is live on `cgroupv1` and there is no recording to fall
back on.

Bring up `hardened` as well only if you intend to prove the S36 allowlist point
live. It is a good five seconds if someone pushes back, and dead weight if
nobody does. Note that this one has never been booted, so it is a full
provision rather than a warm start, and it needs `go build -o nsdemo
./cmd/nsdemo` rather than `make build`. Do it the night before if you want it
at all.

This is also where the stale `qemu/ssh_config` from the repo move heals itself:
`push` and `up` both rewrite it from the current directory. Do not use the raw
`ssh -F qemu/ssh_config` form before either has run.

### T-75: build and pre-flight, inside the VM

```bash
./qemu/vm.sh ssh primary
make build
make check
./scripts/demo.sh check
```

**Every line of `check` must be clean before you walk out.** It asserts
`kernel.apparmor_restrict_unprivileged_userns` is `0` directly rather than
inferring it, because beat 3 of Demo A prints `operation not permitted` whether
it is working correctly or being blocked, so beat 3 is the one thing in the
talk you cannot use as a signal.

`check` also sweeps a leftover `backdoor:` line out of `/etc/passwd`. If it
reports one, something interrupted a previous demo 3 and you want to know that
now rather than at 18:00.

### T-70: build the stage

```bash
./scripts/stage.sh
tmux attach -t talk
```

Left pane is rootful and red, right is rootless and green, both in `~/crl`.
Build it at T-70 rather than on stage so that a tmux that comes up wrong is a
problem you have seventy minutes to fix. Cloud-init grants the demo user
`NOPASSWD:ALL`, so the `sudo -i` in the left pane never prompts.

Set the terminal font to 20pt or larger before attaching, and read it from the
back row of the room later.

### T-65: the cgroupv1 window

In a second terminal on the laptop:

```bash
cd ~/Downloads/containerdays-2026/crl/cds-2026-hamburg-crl-code
./qemu/vm.sh ssh cgroupv1
cd ~/crl && make build && ./scripts/demo.sh check
```

Leave it logged in and sitting at a prompt. At 18:30 you switch to this window
rather than dialling out in front of the room. You only need the rootless side
here, since the point of the second half is that rootless gets a warning.

### T-60: the deck

```bash
cd ../cds-2026-hamburg-crl-slides
nvm use
npm run dev
```

`op run` is only needed for `npm install`. The theme is already in
`node_modules`, so the dev server needs no token.

Presenter view on the laptop screen at `http://localhost:3030/presenter`, slide
view on the projector. Note there are two launch configs with two different
ports, `:3030` in this repo and `:3055` in the `crl/` parent, so if a stray dev
server is already running you may be looking at the wrong one.

Click through all forty slides once. You are checking four things: the two
mermaid diagrams render, the Go snippet on S22 loads from `snippets/`, the
fonts are the theme's rather than the system's, and every slide shows notes in
the presenter pane.

### T-45: the room

- Test the projector at the resolution you will actually present at, not the
  one it defaults to.
- Deck in the top two thirds, terminal strip in the bottom third. Set it now
  and never switch windows again, except to the cgroupv1 window at 18:30.
- Read the terminal from the back row. If you cannot, raise the font until you
  can.
- Test the clicker from where you will actually stand.

### T-30: rehearse the two that must not fail

Demo A and demo 3. Not the whole talk, just those two, end to end, out loud.

Demo A because it is four minutes of live terminal with no fallback, and demo 3
because it is the only command in the talk with a real consequence if the
cleanup does not fire.

```bash
RIGHT:  ./scripts/demo.sh a
LEFT:   ./scripts/demo.sh 3
RIGHT:  ./scripts/demo.sh 3
LEFT:   grep backdoor /etc/passwd || echo "clean"
```

### T-15: reset to a clean slate

```bash
RIGHT:  ./scripts/demo.sh check
```

That clears containers, the scratch directory and any leftover passwd line.
Leave both VMs up. Leave the dev server running. Everything else goes.

### T-5: last look

- Both panes clear, prompts visible, both in `~/crl`.
- Deck on S01, presenter view showing the S01 notes.
- cgroupv1 window logged in and idle.
- Water within reach.
- Phone silent and out of sight.

---

## Part 3: the command reference

Same order as the script. Everything below runs inside the VM, from `~/crl`.

### Cold open, before S01, no slide up

```bash
RIGHT:  ./scripts/demo.sh 1
```

### Demo A, after S21, full screen, single pane

```bash
RIGHT:  ./nsdemo 1
RIGHT:  ./nsdemo 2
RIGHT:  ./nsdemo 3        # cut this first under time pressure
RIGHT:  ./nsdemo 4
```

### Demo 2, after S26, split screen

```bash
LEFT:   ./scripts/demo.sh 2
RIGHT:  ./scripts/demo.sh 2
```

### Demo 3, on S28

```bash
LEFT:   ./scripts/demo.sh 3
RIGHT:  ./scripts/demo.sh 3
LEFT:   grep backdoor /etc/passwd || echo "clean"
```

### Demo 4, on S30

```bash
LEFT:   ./scripts/demo.sh 4
RIGHT:  ./scripts/demo.sh 4
```

then, in the cgroupv1 window:

```bash
./scripts/demo.sh 4
```

### Demo 5, on S31

```bash
RIGHT:  ./scripts/demo.sh 5
```

### The sysctls, on S36, optional and read-only

```bash
sysctl kernel.unprivileged_userns_clone
sysctl kernel.apparmor_restrict_unprivileged_userns
```

---

## Part 4: when something dies

- **A demo hangs.** Count ten seconds, then narrate the expected output off the
  slide and move on. There are no recordings, so there is nothing to cut to.
  Do not debug on stage. S27, S29, S30 and S31 are static slides and carry the
  point on their own.
- **A namespace demo fails.** Check the sysctl before anything else. It must
  read `0` on `primary` and `cgroupv1`.

  ```bash
  sysctl kernel.apparmor_restrict_unprivileged_userns
  ```

  Provisioning writes `/etc/sysctl.d/99-crl-demo.conf` to set it. If it reads
  `1`, either provisioning did not finish or something restored the default,
  and nothing rootless will work until it is `0`.
- **Beat 3 fails and you are not sure whether it is meant to.** It is meant to.
  Beat 3 prints `operation not permitted` on a healthy box, which is the whole
  lesson. It prints the same thing when the sysctl is blocking you. Beats 1, 2
  and 4 fail loudly, so judge from those.
- **Demo 3 leaves the backdoor line behind.** The script arms a trap, but if an
  interrupt beat it:

  ```bash
  grep '^backdoor:' /etc/passwd && sudo sed -i '/^backdoor:/d' /etc/passwd
  ```

  Do this before you do anything else. It is a real unauthenticated root
  account on a real machine, throwaway VM or not.
- **The left pane behaves unexpectedly.** `demo.sh` picks Docker over Podman on
  the rootful side whenever `docker` is on `PATH`, and the VM installs both.
  The two engines have separate image stores, so a container that exists in one
  is invisible to the other. `ENGINE=podman ./scripts/demo.sh 2` forces it.
- **The tmux session is gone.** `./scripts/stage.sh` rebuilds it. That costs
  you the `sudo` password in front of the room, which is annoying and not fatal.
- **A whole VM dies.** `./qemu/vm.sh down primary && ./qemu/vm.sh up primary`,
  then `make push` and `make build` again. That is minutes, not seconds, so it
  is a between-talks recovery rather than an on-stage one. On stage, narrate off
  the slides.
- **You finish early.** Take questions. Do not pad.

The script's ranked cuts, in the order to spend them: S07's condo analogy
first, then beat 3 of Demo A, then S33 reduced to a one-line summary.

---

## Part 5: after the session

1. Flip both repos to public. S40 points at them and they are useless until you
   do.
2. Post the deck link.
3. Tear down: `make vm-down` brings all three VMs down.
4. Confirm no `backdoor:` line survived anywhere.

---

## Part 6: pacing card

Cumulative, from zero. If a number on the left is ahead of the clock, talk
faster from there rather than trimming. The two hard checkpoints are 11:00 and
22:00.

| Time | Beat | Slide |
|---|---|---|
| 0:00 | Cold open, no slide up | |
| 1:30 | Title | S01 |
| 1:55 | Three claims | S03 |
| 5:20 | Identity, the two files | S09 |
| 7:35 | Filesystem | S13 |
| 8:25 | cgroups | S15 |
| 9:20 | Network | S17 |
| **11:00** | **Demo A, hard checkpoint** | S21 |
| 15:00 | Demos 1 to 5 | S26 |
| 16:45 | Demo 3, the money demo | S28 |
| 18:15 | Demo 4 | S30 |
| 19:30 | Demo 5 | S31 |
| 21:00 | The numbers | S32 |
| **22:00** | **Out of the demos, hard checkpoint** | S33 |
| 23:30 | The twist | S35 |
| 25:00 | Kubernetes | S37 |
| 27:00 | Five takeaways | S39 |
| 30:00 | Questions, against S39 | |

These marks are derived from the two checkpoints, not measured. Time yourself
twice and move them, in this table and in the deck's notes together.
