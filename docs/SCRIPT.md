# Speaker script: The Reality of Rootless Containers

ContainerDays 2026, Hamburg. 35 minute slot: 30 speaking, 5 Q&A.

How to use this document: it is written to be read once in rehearsal and then
never read again on stage. Each slide gets a short "SAY" block, in your own
words, not a script to recite; a "DO" block with the exact command to type,
copied verbatim from `scripts/demo.sh` and `cmd/nsdemo` in the companion repo;
and a "WATCH FOR" line where something can go visibly wrong. Timings are
cumulative from the start of the talk, not per slide, so you can glance at a
clock rather than a stopwatch.

All of this is also in the deck's own presenter notes, keyed to the same slide
numbers, so you should not need this document on stage at all. The operational
half, what to run and when, is in [`EVENT-DAY.md`](EVENT-DAY.md).

Slide numbers here are the deck's, S01 to S40.

Two hard checkpoints. If you are not starting Demo A by **11:00**, cut the
condo analogy on S07. If you are not out of the demo block by **22:00**,
skip straight to S34. Everything after S34 is what gets quoted
afterwards. Protect it.

Pane convention throughout: **LEFT is rootful** (root shell, red), **RIGHT is
rootless** (your normal user, green). This matches the deck's colour system
exactly, so you never have to explain which pane is which after S01.

---

## Before you walk on stage

On the laptop, from the companion repo:

```console
$ make vm                       # boot primary
$ make vm-v1                    # boot cgroupv1, demo 4's second half is live
$ make push                     # copy the repo into ~/crl in the VM
$ make push VARIANT=cgroupv1
```

`push` is not optional. QEMU shares no filesystem with the host, so nothing
else puts the code in the VM, and every command below runs from `~/crl`.

Then inside the VM, `./qemu/vm.sh ssh primary`:

```console
$ make build                    # produces ./nsdemo
$ ./scripts/demo.sh check       # pre-flight, the morning of the talk
$ ./scripts/stage.sh            # builds the two-pane tmux session
```

`stage.sh` runs `sudo -i` in the left pane, so the password is typed once,
before anyone is watching. Left is rootful, right is rootless, both already in
`~/crl`.

`check` reports your UID, the engine in use, `/etc/subuid` and `/etc/subgid`,
whether `newuidmap` is on PATH, whether `nsdemo` is built, the cgroup
controllers visible, and whether the demo image is already pulled. It also
asserts `kernel.apparmor_restrict_unprivileged_userns` is `0` directly, because
beat 3 of Demo A prints the same `operation not permitted` whether it is
working or being blocked. **Every line must be clean before you walk out.**

There are no recorded fallbacks. Demo 4's second half runs live on the
`cgroupv1` VM and demo 5 runs live on `primary`, so both boxes need to be up
before you start.

---

## Cold open (0:00, before S01, no slide showing)

Both panes already running, terminal full screen, no slides visible yet.

**DO**, right pane:

```console
$ ./scripts/demo.sh 1
```

Let it run. It starts a rootless container, prints `id` from inside it,
finds the PID on the host, prints the host's view of that same PID, then
prints the container's `uid_map`.

**SAY**, once the container's `id` and the host's `Uid:` line are both on
screen:

> Same process. The container says root. The host says something else. Both
> of them are telling the truth. The gap between those two answers is the
> entire subject of this talk.

Ten seconds. Do not explain further yet. Switch to S01.

**WATCH FOR**: if `podman` has not pre-pulled the image, this will hang on a
pull. You confirmed this at pre-flight; do not skip that step.

---

## S01 (1:30), layout: cover

**SAY**: your name, "Automattic", nothing else. Ten seconds, then advance.

## S02 (1:40), "Both answers are true"

**SAY**: let it sit for two seconds. Do not narrate it, the cold open already
made the point. This slide is the anchor you will gesture back to later.

## S03 (1:55), Three claims

**SAY** the three claims as written on the slide, then land the credibility
line that is in the presenter notes but not on the slide:

> I run rootless. I also think most of what people believe about it is
> wrong.

This is what buys you the right to be critical for the next twenty minutes
without sounding like a hit piece.

## S04 (2:30), "UID 0 is a shorthand"

**SAY**: since Linux 2.2 the kernel has not really asked "is this UID 0". It
asks whether the process holds the capability, evaluated against the
namespace that owns the resource.

## S05 (2:55), The credential model (mermaid diagram)

**SAY**: walk the diagram top to bottom, in order: UID/GID is who you claim
to be, capability sets are what you may do, the user namespace owns the
other four namespace types, and the capability check runs against whoever
owns the resource, not against a global notion of root.

## S06 (3:35), "A user namespace is the only namespace that owns other namespaces"

**SAY**: that ownership is the entire mechanism. It is why an unprivileged
user can configure a network interface, mount a filesystem, and set a
hostname, all without ever holding root on the host.

## S07 (3:55), The condo committee

**CHECKPOINT**: if you are running behind, cut this slide entirely and go
straight to S08.

**SAY**: the analogy once, sixty seconds, do not linger. Elected to the
committee, real authority inside the compound, a person with a clipboard the
moment you step onto the main road. Let the table on the right do the
mapping work; do not read every row aloud, just the first and last.

## S08 (4:55), Five moving parts (progress marker, reused 5 times)

**SAY**: name all five once. This slide reappears with the highlight moved;
you do not need to re-explain it each time, just glance at it as you land on
each section.

## S09 (5:20), `/etc/subuid`

**SAY**: the administrator delegated 65,536 UIDs, starting at 100000, to me.
I may map them however I like inside a namespace I create.

**DO** nothing yet, this is a slide, not a terminal. The live version comes
in Demo A.

## S10 (5:50), The map has two lines

**SAY**, and pause after the first sentence:

> Container UID 0 is your own host UID. Not the start of the delegated
> range.

Then: a file written as "root" inside a rootless container lands on disk
owned by you. This is the fact almost everyone in the room has wrong,
including people who run rootless containers daily. Let it land before
moving on.

## S11 (6:35), Who writes that map? (mermaid diagram)

**SAY**: walk the diagram. Podman forks a child into an unmapped namespace,
then execs the setuid helper `newuidmap`, which checks `/etc/subuid` for the
caller and writes `/proc/<pid>/uid_map` on its behalf.

## S12 (7:05), "Two setuid-root binaries sit in your trust path" (amber)

**SAY**: `newuidmap` and `newgidmap`, from shadow-utils. Small, well
reviewed, but root, and they have had bugs, CVE-2018-7169 in `newgidmap`
before shadow 4.6. "Rootless" describes the runtime, not the installation.

This is amber slide one of three. Do not explain the colour, just let the
tone shift.

## S13 (7:35), Filesystem: three eras

**SAY**: most rootless folklore about slow builds dates from the first row,
`fuse-overlayfs`. Kernel 5.11 moved overlayfs itself inside a user
namespace. 5.12 added idmapped mounts.

## S14 (7:55), Idmapped mounts

**SAY**: the VFS shifts ownership on the fly. No recursive chown, no copy.
Then plant the seed deliberately:

> Remember this one. It is the exact kernel feature that unblocked stateful
> pods in Kubernetes, and we come back to it in the last five minutes.

## S15 (8:25), cgroups: three outcomes

**SAY**: read the three rows in order, cgroup v2 with delegation works,
cgroup v2 without delegation is silently unavailable, cgroup v1 is not
supported and ignored. Do not explain the mechanism yet, that is the next
slide.

## S16 (8:55), "not an error, it is a no-op with a warning" (amber)

**SAY**, slowly:

> A container you believed was capped is not capped.

Ten seconds of silence. This is amber slide two of three, and it is a
genuine production hazard, not a curiosity.

## S17 (9:20), veth vs tap

**SAY**: one end of a veth pair has to live in the host's network namespace.
An unprivileged user cannot put it there. So rootless containers run a
userspace TCP/IP stack instead, `pasta` or `slirp4netns`.

## S18 (9:55), What userspace networking costs

**SAY**: run down the five bullets at a clip, this is a list slide, do not
dwell. The one worth a half-beat extra: the source IP the container sees is
often the gateway, not the real client, which silently breaks IP allow-lists
and access logs.

## S19 (10:20), pasta vs slirp4netns

**SAY**: `pasta` copies the host's actual addresses and routes rather than
inventing a fake subnet, and it is meaningfully faster. If you benchmarked
rootless networking before 2024, benchmark it again.

## S20 (10:40), "Membership of the docker group has always been root"

**SAY**: fifteen seconds. No privileged daemon in Podman rootless. Docker
rootless still has a daemon, but it is yours, socket under
`$XDG_RUNTIME_DIR`. This deletes an entire category of operational mistake.

**CHECKPOINT**: you should be at roughly 11:00 now. If you are past it, skip
the condo analogy note is moot, you are already at Demo A. Move straight in.

---

## S21 (11:00), DEMO A title card

**SAY**: "Podman did this in one command. What did it actually do?"

Switch to a full screen terminal, single pane, right side only. This is the
one segment in the talk with no split screen.

## S22 (11:15), A Go program cannot unshare itself

**SAY**: `CLONE_NEWUSER` is only legal for a single-threaded process. The Go
runtime is multi-threaded before `main` starts. A Go program can never
unshare itself into a user namespace, it has to fork and exec, which is
exactly what `runc` does. Point at the three fields in the code block,
`ContainerID: 0, HostID: uid, Size: 1`, and say: that struct literal *is*
the `uid_map` line from S10.

## Beat 1 (11:45)

**DO**:

```console
$ ./nsdemo 1
```

Expected output: `uid=65534`, empty `uid_map`, `CapPrm` and `CapEff` all
zeroes, `CapBnd` `000001ffffffffff`.

**SAY** once it prints:

> Not root. Not you. Nobody. And look at the two capability sets. The
> bounding set is full and the effective set is empty. The ceiling is
> unlimited, the holding is nothing.

Advance to S23.

## S23 (12:15), "Full capability set, held by nobody"

Statement slide, let it sit for two seconds, then straight back to the
terminal for beat 2.

## Beat 2 (12:30)

**DO**:

```console
$ ./nsdemo 2
```

Expected output: `uid=0`, `uid_map` reads `0 <your uid> 1`, then three probe
lines, `DENIED` on `/etc/shadow`, `ALLOWED` on a tmpfs mount, `ALLOWED` on a
write into `$HOME`.

**SAY**:

> Root on a host file, denied. Root in my own home directory, allowed.
> Remember this, it is the audit section of the whole talk, arriving
> eighteen minutes early. We come back to it at S28.

## S24 (13:15), `0 1000 1`

**SAY**: three fields, that is the entire security model made visible.
Two seconds, advance.

## Beat 3 (13:30)

**CHECKPOINT**: if you are running long, skip this beat, advance straight to
S25 and say the one line from the slide notes instead of running it.

**DO**:

```console
$ ./nsdemo 3
```

Expected output: a request for a 65536-wide range, then `cmd.Start()` fails
with `operation not permitted`.

**SAY**: the failure comes from `Start()`, not from inside the child,
because the write to `uid_map` happens during process creation. Widening the
map beyond your own ID needs `CAP_SETUID` in the parent namespace, which you
do not have. This is why `/etc/subuid` and the helpers from S09 exist.

## Beat 4 (14:15), the best beat

**DO**:

```console
$ ./nsdemo 4
```

Expected output, in order: the child reports as unmapped and blocked, two
lines showing `newuidmap` and `newgidmap` being invoked with the exact
`ContainerID HostID Size` triples, then the same PID reporting again, now as
`uid=0` with a two-line map.

**SAY**, this is the closing line for the whole segment, say it exactly:

> Same PID. No setuid call. Nothing about this process changed. The kernel
> changed its mind about who it is.

## S25 (15:00), "Same PID. No setuid call." (statement, echoes what you just said)

Let it sit. Do not re-explain, you already said it. This is the visual
confirmation. Advance immediately to S26.

**CHECKPOINT**: you should be at 15:00. Demo A budget was four minutes from
11:00; if beat 3 was cut you are ahead, good, bank the time for Demo 3 later.

---

## S26 (15:00), DEMO section card: "One script. Two panes."

**SAY**: point at the two lines on the slide.

> Left pane, sudo, this script, argument three. Right pane, same file, same
> argument, no sudo. The only difference is which pane I typed it in.

Say this once. Do not repeat it for every subsequent demo, the audience has
it now.

Switch to the split screen, both panes visible.

## Demo 1 recap (15:20)

You already ran `./scripts/demo.sh 1` in the cold open. Do not rerun it. If
you want a physical beat here, just gesture at the earlier output still
scrolled in the terminal history, or say the one line:

> You saw this already, at the very start. Container 0 to host 1000 on the
> right, container 0 to host 0 on the left. No translation, no boundary.

## Demo 2 (15:30)

**DO**, either pane, the script prints both automatically because it is run
once per pane in sequence, so run it in LEFT first, then RIGHT, back to
back:

```console
LEFT:   ./scripts/demo.sh 2
RIGHT:  ./scripts/demo.sh 2
```

Expected: `tmpfs mount: OK` on both. Block-device mount and `mknod` both
succeed on the left and print `DENIED` on the right.

**SAY** after both have run:

> tmpfs is allowed everywhere, it carries a flag that says a user namespace
> can mount it. The real disk and the device node are not, even under
> `--privileged` on the right. Privileged in rootless mode grants everything
> your namespace holds, which still is not device access.

## S27 (16:30), "tmpfs yes. ext4 no."

Statement slide confirming what you just showed. Two seconds, advance.

## S28 (16:45), demo 3, the money demo

**SAY** before running it:

> Same mistake on both sides, a volume mount most of you have written by
> accident at some point.

**DO**:

```console
LEFT:   ./scripts/demo.sh 3
RIGHT:  ./scripts/demo.sh 3
```

Expected, read step: LEFT prints the first line of `/etc/shadow`. RIGHT
prints `DENIED`.

Expected, write step: LEFT prints `WROTE`, and the script cleans the
backdoor line from `/etc/passwd` immediately afterwards. RIGHT prints
`DENIED`.

**SAY**, the moment LEFT prints `WROTE`, let it sit for one full second
before speaking:

> That host now has an unauthenticated root account. One flag did it.

Then, gesturing at both outputs together:

> The mount succeeded in both panes. Rootless did not stop the mistake. It
> made the mistake survivable.

**DO**, the third step runs automatically as part of the same invocation,
both panes will print the fake AWS credentials file:

Expected: both LEFT and RIGHT print `not-a-real-key`.

**SAY**:

> And here is the honest half. Your own data is in the blast radius either
> way. For a laptop or a CI runner, this is most of what an attacker wanted
> in the first place.

**WATCH FOR**: confirm on LEFT, after the demo, that `backdoor` is actually
gone from `/etc/passwd`. The script's trap should have handled it, but this
is the one command in the whole talk with a real consequence if it fails.

```console
LEFT:   grep backdoor /etc/passwd || echo "clean"
```

## S29 (18:00), "The mount succeeded in both panes"

Statement slide. This is the sentence people will quote. Pause two full
seconds before advancing. Do not add anything, the slide plus your delivery
a moment ago already did the work.

## S30 (18:15), DEMO 4 section card

**SAY**: "ask for a limit, then ask the container what it got."

## Demo 4 (18:30)

**DO**:

```console
LEFT:   ./scripts/demo.sh 4
RIGHT:  ./scripts/demo.sh 4
```

Expected on a cgroup v2 host with delegation: both panes read
`67108864` from inside the container, the limit applied on both sides.

Expected if you have the second, cgroup v1 machine available: repeat on that
box, RIGHT prints a warning and no working limit, LEFT is unaffected because
root does not go through the same delegation path.

**SAY**:

> 67108864 means the limit landed. Anything else means it did not, and on a
> cgroup v1 rootless host you get a warning, not an error. A container you
> believed was capped is not capped.

The `cgroupv1` VM is already up and already logged in from T-90, so switch to
that window rather than dialling out in front of the room.

## S31 (19:30), DEMO 5 section card

**SAY**: "userspace networking, and its bill." Immediately follow with:

> The throughput numbers were measured beforehand. I am not benchmarking
> network throughput live on conference wifi, and neither should you.

## Demo 5 (19:45)

**DO**, live, quick and safe:

```console
RIGHT:  ./scripts/demo.sh 5
```

Expected: publishing port 80 fails on the right with a `rootlessport` error and
succeeds on the left. The script also prints which of `pasta` and
`slirp4netns` is in use.

Narrate rather than reading the terminal: packets copied through a userspace
process cost throughput and latency, and ports below 1024 need a sysctl that
applies to the whole host rather than to your container. The throughput figures
themselves are on the next slide, from `scripts/bench.sh`.

## S32 (21:00), Measured on my hardware

**SAY**: point at the table and state the numbers from memory rather than
reading each cell. `pasta` beats `slirp4netns`, and both lose to a native veth
pair on throughput.

Say the caveat out loud: this was measured inside a VM, so the ratios are the
result and the absolute figures are not. Method and full output are in
`docs/benchmarks.md` in the companion repo.

**CHECKPOINT**: you should be at 22:00 now. If you are past it, skip
straight to S34 and drop the "what rootless genuinely fixes" table to a
verbal one-liner: "there is a table in the deck, the short version is CI
runners and developer laptops are where this pays off cleanly."

---

## S33 (22:00), What rootless genuinely fixes

**SAY**: this is a table slide, written to be photographed. Do not read
every row. Say the categories once, `docker` group risk, known escape CVEs,
careless volume mounts, multi-tenant CI, then move on. Give the room three
seconds of silence with the slide up before advancing, people are taking
photos.

## S34 (22:45), What it does not touch

**SAY**: same treatment, do not read every row. The one to say out loud
rather than just point at:

> Your own data. SSH keys, cloud credentials, source code, your kubeconfig.
> None of that is protected by any of this.

## S35 (23:30), "Rootless adds an attack surface" (amber, third and last)

**SAY**: pause before this slide even appears if you can control the
advance timing. This is the twist section and the sentence that gets
quoted. Say it plainly, no hedging:

> Rootless adds an attack surface.

Let it sit for two full seconds. This is amber slide three of three; by now
the room reads the colour before you speak.

## S36 (24:00), The sysctls

**DO**, if you want to show these live rather than as a slide, both are
read-only and safe to run on the actual demo box:

```console
$ sysctl kernel.unprivileged_userns_clone 2>/dev/null || echo "not on this distro"
$ sysctl kernel.apparmor_restrict_unprivileged_userns 2>/dev/null || echo "not on this distro"
```

**SAY**: a long line of local privilege-escalation bugs, filesystem mount
parsers, `nf_tables`, overlayfs copy-up, were reachable by unprivileged
users only because they could enter a user namespace first.

Then the correction, which is yours and measured, and which calls back to
Demo A: Ubuntu's sysctl is an allowlist rather than a switch. With it at its
shipped value of 1, Podman still runs rootless containers, because
`/etc/apparmor.d/podman` grants `userns`, as do ninety other profiles on the
image. What gets refused is `unshare` by hand and the program from Demo A,
because nobody wrote them a profile. They watched that refusal happen
already without knowing what it was.

Then the line from the slide, verbatim:

> The feature that makes rootless containers possible is the one hardening
> guides restrict, and what you are trusting instead is a list of 91
> binaries.

If anyone pushes back, it is five seconds on the hardened VM: the sysctl
reads 1, `podman run` works, `./nsdemo 2` is refused.

No exploit code, no live demonstration of any CVE. Show the sysctls, name
the shape of the risk, stop there.

## S37 (25:00), Kubernetes finally caught up

**SAY**: GA in v1.36, April 2026, on KEP-127, open since 2016. Point back at
S14: idmapped mounts are the exact kernel feature that let user-namespaced pods
use volumes under this. State the caveats from the slide in one breath,
modern kernel, CRI support, UID range per pod caps pods-per-node, and the
container still runs as UID 0 inside, `runAsNonRoot` still matters.

## S38 (26:00), So: should you?

**SAY**: do not read the whole table. Land on two rows only: CI runners are
the strongest yes, host devices and low ports are the clearest no. The rest
is there for people to photograph and read later.

## S39 (27:00), Five things to take away

**SAY**: read all five, they are short, this is the slide you leave up
through Q&A. On the fifth point, slow down:

> Rootless needs unprivileged user namespaces, which is itself an attack
> surface. Know your threat model, then choose.

Then the closing line, verbatim, and stop talking immediately after:

> Not a security button. A smaller, sharper blast radius, bought with real
> trade-offs. That is a good deal, as long as you know you are making it.

**Do not advance to the thank-you slide yet.** Leave S39 up. Take
questions against it.

---

## Q&A (30:00 to 35:00)

Jump back to S39 if the deck has moved past it (press `o` for the
overview grid if you need to navigate quickly, then click back in).

Five prepared answers, in the order you are most likely to be asked:

1. **"Is rootless slower?"** Networking, yes, and storage on old kernels.
   CPU-bound workloads are unaffected. Cite your own numbers from S32.

2. **"Can I run Kubernetes rootless?"** Two different questions. A fully
   rootless kubelet, usernetes, k3s rootless, is niche. Pods with
   `hostUsers: false` on an ordinary cluster is GA and boring, which is the
   good kind of answer.

3. **"Does rootless replace gVisor, Kata, Firecracker?"** No, they compose.
   Rootless narrows what an escape yields. Those change what an escape has
   to break through in the first place.

4. **"We disabled unprivileged user namespaces for hardening, now what?"**
   That is a coherent position, not a mistake. You have chosen rootful with
   tight seccomp and AppArmor and no `docker` group membership. Say so out
   loud rather than pretending both are simultaneously possible.

5. **"Is `--privileged` safe in rootless mode?"** Safer, not safe. It grants
   everything your namespace holds. Demo 2 already showed it cannot reach
   host devices, but it gives full access to everything your own UID owns.

Both repos are linked on S40, so point at the slide rather than reading URLs
out. Mention that the safety note about demo 3 is in the companion README
before anyone runs it unsupervised.

---

## Full command reference, in run order

Every command this script asks you to type, in one place, for a final dry run
the day before. Everything from `make build` down runs inside the VM, from
`~/crl`.

```console
# on the laptop, T-90
make vm
make vm-v1
make push
make push VARIANT=cgroupv1

# inside the VM, once
make build
make check

# pre-flight, right pane, morning of
./scripts/demo.sh check
./scripts/stage.sh

# cold open, right pane only, before S01
./scripts/demo.sh 1

# Demo A, right pane only, full screen, after S21
./nsdemo 1
./nsdemo 2
./nsdemo 3        # cut first under time pressure
./nsdemo 4

# demos 2 to 5, split screen, left then right for each
./scripts/demo.sh 2    # left, then right
./scripts/demo.sh 3    # left, then right, then verify cleanup:
grep backdoor /etc/passwd || echo "clean"
./scripts/demo.sh 4    # left, then right on primary, then right on cgroupv1
./scripts/demo.sh 5    # right pane only

# read-only, optional, during the sysctl slide
sysctl kernel.unprivileged_userns_clone
sysctl kernel.apparmor_restrict_unprivileged_userns
```

## Still open

- Time this script against a real clock at least twice. The per-slide marks
  are derived from the two checkpoints at 11:00 and 22:00, not measured, so
  they will move after the first full timed run. The deck's presenter notes
  carry the same marks, so change both together.
- S07 wants a photograph and `public/images/` is empty. The slide works
  without one.
