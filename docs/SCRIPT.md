# Speaker script: The Reality of Rootless Containers

ContainerDays 2026, Hamburg. 35 minute slot: 30 speaking, 5 Q&A.

How to use this document: it is written to be read once in rehearsal and then
never read again on stage. Each slide gets a short "SAY" block, in your own
words, not a script to recite; a "DO" block with the exact command to type,
copied verbatim from `code/rootless-demos.sh` and `code/nsdemo`; and a "WATCH
FOR" line where something can go visibly wrong. Timings are cumulative from
the start of the talk, not per slide, so you can glance at a clock rather
than a stopwatch.

Two hard checkpoints. If you are not starting Demo A by **11:00**, cut the
condo analogy on slide 7. If you are not out of the demo block by **22:00**,
skip straight to slide 33. Everything after slide 33 is what gets quoted
afterwards. Protect it.

Pane convention throughout: **LEFT is rootful** (root shell, red), **RIGHT is
rootless** (your normal user, green). This matches the deck's colour system
exactly, so you never have to explain which pane is which after slide 1.

---

## Before you walk on stage

```console
# In the companion repo, both panes open in tmux, split vertically.
LEFT:   sudo -i
LEFT:   cd ~/demos
RIGHT:  cd ~/demos

# Pre-flight, run in the RIGHT pane only, the morning of the talk.
$ ./rootless-demos.sh check
```

`check` reports your UID, the engine in use, `/etc/subuid` and `/etc/subgid`,
whether `newuidmap` is on PATH, whether `nsdemo` is built, the cgroup
controllers visible, and whether the demo image is already pulled. **Every
line must be clean before you walk out.** If `nsdemo` is missing:

```console
$ (cd nsdemo && go vet ./... && go test ./... && go build -o nsdemo .)
```

Confirm the recordings exist as your fallback for Demo 5 and, if you do not
have a second cgroup v1 machine, for Demo 4:

```console
$ ls recordings/*.cast
```

---

## Cold open (0:00, before slide 1, no slide showing)

Both panes already running, terminal full screen, no slides visible yet.

**DO**, right pane:

```console
$ ./rootless-demos.sh 1
```

Let it run. It starts a rootless container, prints `id` from inside it,
finds the PID on the host, prints the host's view of that same PID, then
prints the container's `uid_map`.

**SAY**, once the container's `id` and the host's `Uid:` line are both on
screen:

> Same process. The container says root. The host says something else. Both
> of them are telling the truth. The gap between those two answers is the
> entire subject of this talk.

Ten seconds. Do not explain further yet. Switch to slide 1.

**WATCH FOR**: if `podman` has not pre-pulled the image, this will hang on a
pull. You confirmed this at pre-flight; do not skip that step.

---

## Slide 1 (0:02), layout: cover

**SAY**: your name, "Automattic", nothing else. Ten seconds, then advance.

## Slide 2 (0:02), "Both answers are true"

**SAY**: let it sit for two seconds. Do not narrate it, the cold open already
made the point. This slide is the anchor you will gesture back to later.

## Slide 3 (0:03), Three claims

**SAY** the three claims as written on the slide, then land the credibility
line that is in the presenter notes but not on the slide:

> I run rootless. I also think most of what people believe about it is
> wrong.

This is what buys you the right to be critical for the next twenty minutes
without sounding like a hit piece.

## Slide 4 (0:05), "UID 0 is a shorthand"

**SAY**: since Linux 2.2 the kernel has not really asked "is this UID 0". It
asks whether the process holds the capability, evaluated against the
namespace that owns the resource.

## Slide 5 (0:06), The credential model (mermaid diagram)

**SAY**: walk the diagram top to bottom, in order: UID/GID is who you claim
to be, capability sets are what you may do, the user namespace owns the
other four namespace types, and the capability check runs against whoever
owns the resource, not against a global notion of root.

## Slide 6 (0:07), "A user namespace is the only namespace that owns other namespaces"

**SAY**: that ownership is the entire mechanism. It is why an unprivileged
user can configure a network interface, mount a filesystem, and set a
hostname, all without ever holding root on the host.

## Slide 7 (0:08), The condo committee

**CHECKPOINT**: if you are running behind, cut this slide entirely and go
straight to slide 8.

**SAY**: the analogy once, sixty seconds, do not linger. Elected to the
committee, real authority inside the compound, a person with a clipboard the
moment you step onto the main road. Let the table on the right do the
mapping work; do not read every row aloud, just the first and last.

## Slide 8 (0:09), Five moving parts (progress marker, reused 5 times)

**SAY**: name all five once. This slide reappears with the highlight moved;
you do not need to re-explain it each time, just glance at it as you land on
each section.

## Slide 9 (0:10), `/etc/subuid`

**SAY**: the administrator delegated 65,536 UIDs, starting at 100000, to me.
I may map them however I like inside a namespace I create.

**DO** nothing yet, this is a slide, not a terminal. The live version comes
in Demo A.

## Slide 10 (0:11), The map has two lines

**SAY**, and pause after the first sentence:

> Container UID 0 is your own host UID. Not the start of the delegated
> range.

Then: a file written as "root" inside a rootless container lands on disk
owned by you. This is the fact almost everyone in the room has wrong,
including people who run rootless containers daily. Let it land before
moving on.

## Slide 11 (0:12), Who writes that map? (mermaid diagram)

**SAY**: walk the diagram. Podman forks a child into an unmapped namespace,
then execs the setuid helper `newuidmap`, which checks `/etc/subuid` for the
caller and writes `/proc/<pid>/uid_map` on its behalf.

## Slide 12 (0:13), "Two setuid-root binaries sit in your trust path" (amber)

**SAY**: `newuidmap` and `newgidmap`, from shadow-utils. Small, well
reviewed, but root, and they have had bugs, CVE-2018-7169 in `newgidmap`
before shadow 4.6. "Rootless" describes the runtime, not the installation.

This is amber slide one of three. Do not explain the colour, just let the
tone shift.

## Slide 13 (0:14), Filesystem: three eras

**SAY**: most rootless folklore about slow builds dates from the first row,
`fuse-overlayfs`. Kernel 5.11 moved overlayfs itself inside a user
namespace. 5.12 added idmapped mounts.

## Slide 14 (0:15), Idmapped mounts

**SAY**: the VFS shifts ownership on the fly. No recursive chown, no copy.
Then plant the seed deliberately:

> Remember this one. It is the exact kernel feature that unblocked stateful
> pods in Kubernetes, and we come back to it in the last five minutes.

## Slide 15 (0:16), cgroups: three outcomes

**SAY**: read the three rows in order, cgroup v2 with delegation works,
cgroup v2 without delegation is silently unavailable, cgroup v1 is not
supported and ignored. Do not explain the mechanism yet, that is the next
slide.

## Slide 15 (0:17), "not an error, it is a no-op with a warning" (amber)

**SAY**, slowly:

> A container you believed was capped is not capped.

Ten seconds of silence. This is amber slide two of three, and it is a
genuine production hazard, not a curiosity.

## Slide 17 (0:18), veth vs tap

**SAY**: one end of a veth pair has to live in the host's network namespace.
An unprivileged user cannot put it there. So rootless containers run a
userspace TCP/IP stack instead, `pasta` or `slirp4netns`.

## Slide 18 (0:19), What userspace networking costs

**SAY**: run down the five bullets at a clip, this is a list slide, do not
dwell. The one worth a half-beat extra: the source IP the container sees is
often the gateway, not the real client, which silently breaks IP allow-lists
and access logs.

## Slide 19 (0:20), pasta vs slirp4netns

**SAY**: `pasta` copies the host's actual addresses and routes rather than
inventing a fake subnet, and it is meaningfully faster. If you benchmarked
rootless networking before 2024, benchmark it again.

## Slide 20 (0:20), "Membership of the docker group has always been root"

**SAY**: fifteen seconds. No privileged daemon in Podman rootless. Docker
rootless still has a daemon, but it is yours, socket under
`$XDG_RUNTIME_DIR`. This deletes an entire category of operational mistake.

**CHECKPOINT**: you should be at roughly 11:00 now. If you are past it, skip
the condo analogy note is moot, you are already at Demo A. Move straight in.

---

## Slide 21 (11:00), DEMO A title card

**SAY**: "Podman did this in one command. What did it actually do?"

Switch to a full screen terminal, single pane, right side only. This is the
one segment in the talk with no split screen.

## Slide 22 (11:15), A Go program cannot unshare itself

**SAY**: `CLONE_NEWUSER` is only legal for a single-threaded process. The Go
runtime is multi-threaded before `main` starts. A Go program can never
unshare itself into a user namespace, it has to fork and exec, which is
exactly what `runc` does. Point at the three highlighted fields in the code
block, `ContainerID: 0, HostID: os.Getuid(), Size: 1`, and say: that struct
literal *is* the `uid_map` line from slide 10.

## Beat 1 (11:45)

**DO**:

```console
$ ./nsdemo/nsdemo 1
```

Expected output: `uid=65534`, empty `uid_map`, a full `CapEff` bitmask.

**SAY** once it prints:

> Not root. Not you. Nobody. And look at CapEff, full capability set, held
> by nobody. Creating a user namespace grants the full set inside it,
> mapping or no mapping.

Advance to slide 23.

## Slide 23 (12:15), "Full capability set, held by nobody"

Statement slide, let it sit for two seconds, then straight back to the
terminal for beat 2.

## Beat 2 (12:30)

**DO**:

```console
$ ./nsdemo/nsdemo 2
```

Expected output: `uid=0`, `uid_map` reads `0 <your uid> 1`, then three probe
lines, `DENIED` on `/etc/shadow`, `ALLOWED` on a tmpfs mount, `ALLOWED` on a
write into `$HOME`.

**SAY**:

> Root on a host file, denied. Root in my own home directory, allowed.
> Remember this, it is the audit section of the whole talk, arriving
> eighteen minutes early. We come back to it at slide 28.

## Slide 24 (13:15), `0 1000 1`

**SAY**: three fields, that is the entire security model made visible.
Two seconds, advance.

## Beat 3 (13:30)

**CHECKPOINT**: if you are running long, skip this beat, advance straight to
slide 25 and say the one line from the slide notes instead of running it.

**DO**:

```console
$ ./nsdemo/nsdemo 3
```

Expected output: a request for a 65536-wide range, then `cmd.Start()` fails
with `operation not permitted`.

**SAY**: the failure comes from `Start()`, not from inside the child,
because the write to `uid_map` happens during process creation. Widening the
map beyond your own ID needs `CAP_SETUID` in the parent namespace, which you
do not have. This is why `/etc/subuid` and the helpers from slide 9 exist.

## Beat 4 (14:15), the best beat

**DO**:

```console
$ ./nsdemo/nsdemo 4
```

Expected output, in order: the child reports as unmapped and blocked, two
lines showing `newuidmap` and `newgidmap` being invoked with the exact
`ContainerID HostID Size` triples, then the same PID reporting again, now as
`uid=0` with a two-line map.

**SAY**, this is the closing line for the whole segment, say it exactly:

> Same PID. No setuid call. Nothing about this process changed. The kernel
> changed its mind about who it is.

## Slide 25 (15:00), "Same PID. No setuid call." (statement, echoes what you just said)

Let it sit. Do not re-explain, you already said it. This is the visual
confirmation. Advance immediately to slide 26.

**CHECKPOINT**: you should be at 15:00. Demo A budget was four minutes from
11:00; if beat 3 was cut you are ahead, good, bank the time for Demo 3 later.

---

## Slide 26 (15:00), DEMO section card: "One script. Two panes."

**SAY**: point at the two lines on the slide.

> Left pane, sudo, this script, argument three. Right pane, same file, same
> argument, no sudo. The only difference is which pane I typed it in.

Say this once. Do not repeat it for every subsequent demo, the audience has
it now.

Switch to the split screen, both panes visible.

## Demo 1 recap (15:20)

You already ran `./rootless-demos.sh 1` in the cold open. Do not rerun it. If
you want a physical beat here, just gesture at the earlier output still
scrolled in the terminal history, or say the one line:

> You saw this already, at the very start. Container 0 to host 1000 on the
> right, container 0 to host 0 on the left. No translation, no boundary.

## Demo 2 (15:30)

**DO**, either pane, the script prints both automatically because it is run
once per pane in sequence, so run it in LEFT first, then RIGHT, back to
back:

```console
LEFT:   ./rootless-demos.sh 2
RIGHT:  ./rootless-demos.sh 2
```

Expected: `tmpfs mount: OK` on both. Block-device mount and `mknod` both
succeed on the left and print `DENIED` on the right.

**SAY** after both have run:

> tmpfs is allowed everywhere, it carries a flag that says a user namespace
> can mount it. The real disk and the device node are not, even under
> `--privileged` on the right. Privileged in rootless mode grants everything
> your namespace holds, which still is not device access.

## Slide 27 (16:30), "tmpfs yes. ext4 no."

Statement slide confirming what you just showed. Two seconds, advance.

## Demo 3 (16:45), the money demo

**SAY** before running it:

> Same mistake on both sides, a volume mount most of you have written by
> accident at some point.

**DO**:

```console
LEFT:   ./rootless-demos.sh 3
RIGHT:  ./rootless-demos.sh 3
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

## Slide 28 (18:00), "The mount succeeded in both panes"

Statement slide. This is the sentence people will quote. Pause two full
seconds before advancing. Do not add anything, the slide plus your delivery
a moment ago already did the work.

## Slide 29 (18:15), DEMO 4 section card

**SAY**: "ask for a limit, then ask the container what it got."

## Demo 4 (18:30)

**DO**:

```console
LEFT:   ./rootless-demos.sh 4
RIGHT:  ./rootless-demos.sh 4
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

**IF NO SECOND MACHINE**: skip the live run and play the recording instead:

```console
$ asciinema play recordings/04-cgroups-v1.cast
```

State plainly that this is a recording from a cgroup v1 box, do not present
it as live.

## Slide 30 (19:30), DEMO 5 section card

**SAY**: "userspace networking, and its bill." Immediately follow with:

> This one is pre-recorded. I am not benchmarking network throughput live on
> conference wifi, and neither should you.

## Demo 5 (19:45)

**DO**:

```console
$ asciinema play recordings/05-network.cast
```

While it plays, narrate over it rather than reading the terminal verbatim:
packets copied through a userspace process cost throughput and latency,
ports below 1024 need a sysctl that applies to the whole host, not just your
container.

**Optional live beat**, only if time allows, run this for real in the RIGHT
pane, it is quick and safe:

```console
RIGHT:  ./rootless-demos.sh 5
```

Expected: publishing port 80 fails on the right with a rootlessport error,
succeeds on the left.

## Slide 31 (21:00), Measured on my hardware

**SAY**: point at the table, state your own numbers from memory rather than
reading each cell. If the table is not filled in yet, say so honestly:

> I have not published exact figures for this deck yet, the shape of the
> result is what matters, pasta beats slirp4netns, both lose to a native
> veth pair on throughput.

**CHECKPOINT**: you should be at 22:00 now. If you are past it, skip
straight to slide 33 and drop the "what rootless genuinely fixes" table to a
verbal one-liner: "there is a table in the deck, the short version is CI
runners and developer laptops are where this pays off cleanly."

---

## Slide 32 (22:00), What rootless genuinely fixes

**SAY**: this is a table slide, written to be photographed. Do not read
every row. Say the categories once, `docker` group risk, known escape CVEs,
careless volume mounts, multi-tenant CI, then move on. Give the room three
seconds of silence with the slide up before advancing, people are taking
photos.

## Slide 33 (22:45), What it does not touch

**SAY**: same treatment, do not read every row. The one to say out loud
rather than just point at:

> Your own data. SSH keys, cloud credentials, source code, your kubeconfig.
> None of that is protected by any of this.

## Slide 35 (23:30), "Rootless adds an attack surface" (amber, third and last)

**SAY**: pause before this slide even appears if you can control the
advance timing. This is the twist section and the sentence that gets
quoted. Say it plainly, no hedging:

> Rootless adds an attack surface.

Let it sit for two full seconds. This is amber slide three of three; by now
the room reads the colour before you speak.

## Slide 35 (24:00), The sysctls

**DO**, if you want to show these live rather than as a slide, both are
read-only and safe to run on the actual demo box:

```console
$ sysctl kernel.unprivileged_userns_clone 2>/dev/null || echo "not on this distro"
$ sysctl kernel.apparmor_restrict_unprivileged_userns 2>/dev/null || echo "not on this distro"
```

**SAY**: a long line of local privilege-escalation bugs, filesystem mount
parsers, `nf_tables`, overlayfs copy-up, were reachable by unprivileged
users only because they could enter a user namespace first. Then the line
from the slide, verbatim:

> The feature that makes rootless containers possible is the same feature
> hardening guides tell you to disable.

No exploit code, no live demonstration of any CVE. Show the sysctls, name
the shape of the risk, stop there.

## Slide 36 (25:00), Kubernetes finally caught up

**SAY**: GA in v1.36, April 2026, after ten years of KEP-127. Point back at
slide 14: idmapped mounts are the exact kernel feature that made stateful
pods practical under this. State the caveats from the slide in one breath,
modern kernel, CRI support, UID range per pod caps pods-per-node, and the
container still runs as UID 0 inside, `runAsNonRoot` still matters.

## Slide 37 (26:00), So: should you?

**SAY**: do not read the whole table. Land on two rows only: CI runners are
the strongest yes, host devices and low ports are the clearest no. The rest
is there for people to photograph and read later.

## Slide 38 (27:00), Five things to take away

**SAY**: read all five, they are short, this is the slide you leave up
through Q&A. On the fifth point, slow down:

> Rootless needs unprivileged user namespaces, which is itself an attack
> surface. Know your threat model, then choose.

Then the closing line, verbatim, and stop talking immediately after:

> Not a security button. A smaller, sharper blast radius, bought with real
> trade-offs. That is a good deal, as long as you know you are making it.

**Do not advance to the thank-you slide yet.** Leave slide 38 up. Take
questions against it.

---

## Q&A (30:00 to 35:00)

Jump back to slide 38 if the deck has moved past it (press `o` for the
overview grid if you need to navigate quickly, then click back in).

Five prepared answers, in the order you are most likely to be asked:

1. **"Is rootless slower?"** Networking, yes, and storage on old kernels.
   CPU-bound workloads are unaffected. Cite your own numbers from slide 31.

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

If someone asks for the repos: they are not yet linked from a slide on
purpose, so state the URLs verbally, `github.com/asahasrabuddhe/cds-2026-hamburg-crl-slides` and
`github.com/asahasrabuddhe/cds-2026-hamburg-crl-code`, and mention that the safety note about Demo 3 is
in the companion README before anyone runs it unsupervised.

---

## Full command reference, in run order

Every command this script asks you to type, in one place, for a final dry
run the day before.

```console
# pre-flight, right pane, morning of
./rootless-demos.sh check

# cold open, right pane only, before slide 1
./rootless-demos.sh 1

# Demo A, right pane only, full screen, after slide 21
./nsdemo/nsdemo 1
./nsdemo/nsdemo 2
./nsdemo/nsdemo 3        # cut first under time pressure
./nsdemo/nsdemo 4

# Demos 1 to 5, split screen, left then right for each
./rootless-demos.sh 2    # left, then right
./rootless-demos.sh 3    # left, then right, then verify cleanup:
grep backdoor /etc/passwd || echo "clean"
./rootless-demos.sh 4    # left, then right; or on cgroup v1 box
asciinema play recordings/04-cgroups-v1.cast   # if no v1 box
asciinema play recordings/05-network.cast      # always, for demo 5
./rootless-demos.sh 5    # optional live extra, right pane only

# read-only, optional, during the sysctl slide
sysctl kernel.unprivileged_userns_clone
sysctl kernel.apparmor_restrict_unprivileged_userns
```

## Things to fill in before this script is final

- Slide 31's benchmark numbers. Run the benchmarks in the design document,
  section 12, and put the real figures on the slide and in this script's
  "SAY" line.
- The two repo URLs in the Q&A section.
- Confirm whether you have a second cgroup v1 machine for Demo 4. This
  script currently branches on that; if the answer is settled, delete the
  branch you will not use.
- Time this script against a real clock at least twice before the talk. The
  checkpoints at 11:00 and 22:00 are estimates from the design document, not
  guarantees, adjust them after your first full timed run.
