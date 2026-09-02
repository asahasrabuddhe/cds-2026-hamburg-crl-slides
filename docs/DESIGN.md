# The Reality of Rootless Containers. Talk Design

**Event:** ContainerDays 2026, Hamburg (27 to 28 October)
**Speaker:** Ajitem Sahasrabuddhe
**Slot:** 35 minutes, 30 speaking + 5 Q&A.
**Audience:** platform engineers, DevOps, architects, people who already run containers in production.

---

## 1. The spine

Every talk needs one sentence the audience can repeat in the corridor afterwards.

> **Rootless does not remove root. It moves the boundary, shrinks the blast radius, and hands you a new set of problems, including a new attack surface of its own.**

Three claims follow from it, and the whole talk exists to prove them:

| Claim | Proven by |
|---|---|
| "Root" is not UID 0. Root is a set of capabilities, evaluated against a namespace. | Demo 2 and 3. CAP_SYS_ADMIN that can mount tmpfs but cannot touch a disk |
| Rootless shrinks the blast radius of an escape; it does not prevent the escape. | Demo 4 (`-v /etc`) and the CVE table |
| You pay for it, in features, in throughput, and in one new attack surface. | Demo 5, 6 and 7 |

Everything that does not serve one of those three claims gets cut. That rule alone will keep this talk at 30 minutes.

### Framing choice

The abstract promises "a clear, realistic understanding of what rootless improves, what it doesn't". So the tone is **the honest engineer's audit**, not the advocacy talk. Say early and plainly: *I run rootless. I also think most of what people believe about it is wrong.* That buys credibility for the criticism later and stops the talk reading as a hit piece.

### Timing gift

User namespaces went GA in Kubernetes v1.36, released 23 April 2026, on KEP-127, which has been open since 2016. By late October the room will have read the release notes and roughly nobody will have run it in anger. That is the perfect gap for this talk to fill, and it earns the last section a real reason to exist.

---

## 2. Structure

Following the Concept → Analogy → Architecture → Code → Tests → Real-world flow, mapped onto a 30-minute wall clock:

```
0:00  ┌─ COLD OPEN (live, no slides) ────────────────┐  2 min
      │  "Am I root? Yes. And no."                   │
2:00  ├─ CONCEPT: what root actually is ─────────────┤  3 min
      │  capabilities, not UID 0                     │
5:00  ├─ ANALOGY: the condo committee ───────────────┤  1 min
6:00  ├─ ARCHITECTURE: the five moving parts ────────┤  5 min
      │  IDs · filesystem · cgroups · network · exec │
11:00 ├─ DEMO A: build the mapping by hand ──────────┤  4 min
      │  unshare, uid_map, newuidmap, live          │
15:00 ├─ DEMOS 2–5: rootful vs rootless, side by side┤  7 min
      │  hard-stopped                                │
22:00 ├─ THE AUDIT: what improves, what doesn't ─────┤  2.5 min
24:30 ├─ THE TWIST: rootless as attack surface ──────┤  2.5 min
27:00 ├─ REAL WORLD: CI, K8s 1.36, the checklist ────┤  2.5 min
29:30 └─ TAKEAWAYS ──────────────────────────────────┘  0.5 min
30:00    Q&A
```

**Where the extra five minutes went.** The 30-minute build had no room to *show* the identity machinery, it only described it. Demo A (§7.0) fixes that, and it is the right place to spend the time: it is the one segment nobody else's rootless talk has, it turns the most abstract part of the talk into something the audience can type on the flight home, and it reuses your existing container-runtime material. Resist the temptation to spend the five minutes on more networking detail; that is depth the room will not retain.

**Checkpoints, the only two you need to watch.** If you are not starting Demo A by **11:00**, cut the condo analogy from the next run. If you are not out of the demo block by **22:00**, drop Demo 5 and go straight to the audit. Everything after 22:00 is the part people quote; protect it.

**Cut to 20 minutes** (for a meetup rerun): drop Demo A and Demo 4, compress the architecture to IDs + filesystem, keep the twist section intact.

---

## 3. The cold open

Do not open with your bio. Open with a terminal, already running, split into two panes.

**Left pane (rootless):**

```console
$ podman run --rm -it alpine sh
/ # id
uid=0(root) gid=0(root) groups=0(root),...
/ # whoami
root
```

**Right pane (host, same machine):**

```console
$ ps -eo user,pid,comm | grep -w sh
ajitem   48213 sh
$ cat /proc/48213/status | grep -E '^(Uid|Gid):'
Uid:    1000    1000    1000    1000
```

Then say, and let it sit for a beat:

> Same process. The container says root. The host says `ajitem`. **Both of them are telling the truth**, and the gap between those two answers is the entire subject of this talk.

Then one slide: title, name, handle. Ten seconds. Back to content.

*Why this works:* it is the shortest possible statement of the thesis, it is live, and it puts a concrete artefact, a PID, on screen before any theory arrives.

---

## 4. Concept: root is a claim, not an identity

The slide sequence here is short, punchy, and mostly diagram.

The point to land: since Linux 2.2 the kernel has not really asked "is this UID 0?". It asks "does this process hold the capability for this operation, **in the user namespace that owns the resource**?" UID 0 survives as a shorthand, a process with UID 0 in a namespace gets the full capability set within it.

The three-line version on a slide:

```
Rootful:   CAP_SYS_ADMIN in the initial user namespace   → the machine
Rootless:  CAP_SYS_ADMIN in a child user namespace       → your sandbox
The kernel:  same check, different owner
```

Then the credential model, as a diagram:

```
        process
          │
          ├── UID / GID              (who you claim to be)
          ├── capability sets        (what you may do)
          └── user namespace  ───────┐
                                     ↓
                          owns → mount ns, net ns, pid ns, ipc ns, uts ns
                                     │
                          capability check runs against THIS owner
```

Worth saying out loud, because it is the sentence most people have never heard: **a user namespace is the only namespace that owns other namespaces.** That ownership is the mechanism by which an unprivileged user gets to configure a network interface, mount a filesystem, and set a hostname without ever being root on the host.

---

## 5. Analogy: the condo management committee

One slide, sixty seconds, then move on.

You get elected to the management committee of your condo. Inside the compound you have real authority: you decide where the visitor parking goes, you can change the gate timings, you can repaint the lobby. Walk out of the gate and try to redirect traffic on the main road, and you are a person in a polo shirt with a clipboard.

The mapping, which should be on the slide as a table because the audience will photograph it:

| Condo | Container |
|---|---|
| The compound | The user namespace |
| Committee authority inside it | Capabilities in that namespace |
| The land title, held by the developer | The host's initial user namespace |
| Your flat number in the master land registry | Your host UID (1000), unchanged |
| Redecorating the lobby | `mount -t tmpfs`, `ip link add` |
| Redirecting traffic on the main road | `mount /dev/<disk>`, loading a kernel module |

The reason this beats the usual "sandbox" or "walled garden" image: it makes the *derived, granted, revocable* nature of the authority obvious. Nobody thinks the committee owns the land.

---

## 6. Architecture, the five moving parts

This is the technical core and the section that earns the "deep dive" label. Structure it as five numbered parts so the audience can track where they are. One slide per part, plus one code or config artefact per part.

### 6.1 Identity, subuid, subgid, and the setuid helpers

The whole system rests on two boring files:

```console
$ cat /etc/subuid
ajitem:100000:65536
$ cat /etc/subgid
ajitem:100000:65536
```

Read out loud: *the administrator has delegated 65,536 UIDs, starting at 100000, to me, and I may map them however I like inside a namespace I create.*

Inside the container:

```console
/ # cat /proc/self/uid_map
         0       1000          1
         1     100000      65536
```

Two lines, and the first one surprises almost everybody. **Container UID 0 is your own host UID**, not the start of the delegated range. Container UIDs 1 and upwards map into the `/etc/subuid` block. That is why a file written by "root" inside a rootless container lands on disk owned by you, and why a file written by UID 1000 inside lands on disk owned by 100999.

The trick to explain clearly, and the one thing to get right: an unprivileged process can create a user namespace, but it can only map **its own** UID into it, that first line, size 1. Mapping the range on the second line needs the setuid helpers `newuidmap` and `newgidmap` from shadow-utils, which read `/etc/subuid` and do the write on your behalf.

```
podman ──fork──▶ child in new userns (unmapped, uid 65534 "nobody")
   │
   └──exec──▶ newuidmap <pid> 0 1000 1  1 100000 65536
                  │
                  ├─ setuid root binary
                  ├─ checks /etc/subuid for the caller
                  └─ writes /proc/<pid>/uid_map
                            │
                            └──▶ child now reads as uid 0 inside, 1000 outside
```

Say the quiet part: **there are two setuid-root binaries in the trust path of your rootless stack.** They are small and well reviewed, but they are root, and they have had bugs (`newgidmap` before shadow 4.6, CVE-2018-7169, let a caller drop supplementary groups). "Rootless" describes the runtime, not the installation.

### 6.2 Filesystem, the chown tax and how idmapped mounts removed it

Three eras, one slide, because this is where most of the old folklore comes from:

| Era | Mechanism | Cost |
|---|---|---|
| Before kernel 5.11 | `fuse-overlayfs` in userspace | every metadata operation crosses into userspace; slow builds |
| Kernel 5.11+ | overlayfs mountable inside a user namespace | native speed for images |
| Kernel 5.12+ | **idmapped mounts** | a shared volume can be presented with a UID shift, no recursive `chown` |

The pre-5.12 problem, stated concretely: your image layers on disk are owned by 100000+, not by root. A 3 GB image extracted for the first time used to mean touching every inode. Idmapped mounts moved the translation into the VFS, the same directory appears as owned by `0` in the container and by `100000` on the host, with no on-disk change.

This is also the exact kernel feature that let user-namespaced pods use volumes (KEP-127), so it sets up the Kubernetes section later. Plant the seed here, harvest it at 26:00.

### 6.3 cgroups, where rootless quietly stops working

The honest slide of the section, and the one that gets a laugh of recognition.

```
cgroup v2 + systemd delegation      →  limits work
cgroup v2, no delegation            →  limits silently unavailable
cgroup v1                           →  "Resource limits are not supported
                                        and ignored on cgroups V1 rootless"
```

The mechanism: systemd starts a `user@<uid>.service` slice and delegates a subtree to it. Only the controllers listed in `Delegate=` are yours.

```console
$ systemctl show user@$(id -u).service -p Delegate -p DelegateControllers
$ cat /sys/fs/cgroup/user.slice/user-$(id -u).slice/cgroup.controllers
cpu io memory pids
$ podman info --format '{{.Host.CgroupsVersion}} {{.Host.CgroupControllers}}'
```

The line to say: **on a cgroup v1 host, `--memory=512m` is not an error. It is a no-op with a warning.** A container you believed was capped is not capped. That is a genuine production hazard and it is worth ten seconds of silence.

### 6.4 Network, nobody gave you a tap device

An unprivileged user cannot create a veth pair that reaches the host bridge, because one end of it lives in the host's network namespace. So rootless containers use a userspace TCP/IP stack.

```
  ROOTFUL                              ROOTLESS
  container netns                      container netns
      │ veth                               │ tap
      ▼                                    ▼
  host bridge (docker0/cni0)          pasta / slirp4netns
      │                                    │  userspace TCP/IP
      ▼                                    ▼
  host netns → NIC                    host socket → NIC
                                      (owned by your UID)
```

The consequences, which is what the audience actually needs:

- Packets are copied through a userspace process. Throughput and latency both suffer, measure it, don't guess (see §8).
- Ports below 1024 are refused unless the admin lowers `net.ipv4.ip_unprivileged_port_start`.
- `ping` works only if your GID falls inside `net.ipv4.ping_group_range`.
- The source IP seen by the container is often the gateway, not the real client, depending on the port-forwarding driver. This breaks IP allow-lists and access logs.
- Most CNI plugins assume host privileges and simply do not apply.

`pasta` (Podman's default since 5.0) is a meaningful improvement on `slirp4netns`, it copies the host's addresses and routes rather than inventing a fake subnet, and it is faster. Docker's rootless mode ships `slirp4netns` by default with `pasta` available in recent versions.

### 6.5 Execution, what is left of the runtime

Short slide, mostly for completeness: `crun`/`runc` still create the same namespaces, apply the same seccomp profile, and still cannot ask for anything the parent user namespace does not hold. `conmon` supervises. No daemon runs as root, with Docker rootless there is still a daemon, just an unprivileged one owned by you, with its socket under `$XDG_RUNTIME_DIR` instead of `/var/run/docker.sock`.

Worth one line, because it is one of the strongest real arguments for rootless: **membership of the `docker` group has always been equivalent to root on that host.** Rootless deletes that entire category of mistake.

---

## 7. The demos

Eleven minutes total: one four-minute build-it-by-hand segment, then seven minutes of side-by-side comparison. Every demo answers a question posed on the previous slide, and the comparison demos run the **same command twice**, once rootful, once rootless, because the contrast is the content.

### The staging device: one script, two panes

`scripts/demo.sh` detects which side it is running on and runs identical commands either way. Left pane is a root shell, right pane is yours:

```console
# LEFT  (red)     sudo -i ; cd ~/crl ; ./scripts/demo.sh 3
$ RIGHT (green)             cd ~/crl ; ./scripts/demo.sh 3
```

Say it out loud the first time and then never again: **the same file, the same argument. The only difference is which pane I typed it in.** That removes every "trust me, the other side does this" from the talk, you are not describing the rootful behaviour, you are running it. It also halves what you have to remember on stage, because there is one command per demo rather than two.

Where the two sides need different commentary, and they do, four times, the script prints it for you. Do not read those lines aloud; they are cue cards, not a script.

### How much of the talk is actually containers?

Worth being deliberate about, since the title says containers and Demo A does not use one. The split across the 11 demo minutes:

| Segment | Containers on screen? |
|---|---|
| Cold open (2 min, before the demo block) | Yes, rootless container, both views of one PID |
| Demo A (4 min) | No, bare namespaces and Go |
| Demos 1–5 (7 min) | Yes, every one is rootful vs rootless, side by side |

So containers are on screen for roughly nine of the thirteen minutes you spend in a terminal, and Demo A is the four minutes that explain *why* the other nine behave as they do. That ordering is deliberate: run Demo A after the architecture section and the comparison demos stop needing explanation. If you ever have to cut, cut Demo A, never the comparisons. The comparisons are the argument; Demo A is the proof.

Layout: tmux, two panes, left labelled `ROOTFUL (uid 0)` in red, right labelled `ROOTLESS (uid 1000)` in green. Font at 20pt minimum. Prompt stripped down to `#` and `$`.

### Demo A: Build the mapping by hand, in Go (4 min), **the segment that makes this talk yours**

Question on the slide: *Podman did this in one command. What did it actually do?*

Run `nsdemo`, a ~250-line Go program (in `nsdemo/`). Single full-screen terminal. Four beats, one subcommand each.

**The Go detail to open with**, because it earns the language choice in ten seconds and half the room will not know it:

> `CLONE_NEWUSER` is only legal for a single-threaded process. The Go runtime is multi-threaded before `main` starts. **A Go program can never unshare itself into a user namespace**, it has to fork and exec. Which is exactly what `runc` does, and exactly why every beat here re-execs `/proc/self/exe`.

Put that on the slide next to the four lines that do the work:

```go
cmd := exec.Command("/proc/self/exe", "child")
cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS,
    UidMappings: []syscall.SysProcIDMap{
        {ContainerID: 0, HostID: os.Getuid(), Size: 1},
    },
}
```

**Beat 1: a namespace with no map at all.**

```console
$ ./nsdemo 1
  ── inside the new user namespace ──
  uid=65534  euid=65534  gid=65534
  uid_map: (empty, no mapping exists)
  CapPrm:  0000000000000000
  CapEff:  0000000000000000
  CapBnd:  000001ffffffffff
```

Two things to say. *This is what a container looks like before anyone tells the kernel who you are: not root, not you, nobody.* And then the nastier one, look at the two sets together. **Full capability set, held by nobody.** The bounding set is unlimited and the effective set is empty. See the amendment below for why the effective set arrives empty rather than full.

**Beat 2: map yourself, and become root.**

```console
$ ./nsdemo 2
  parent: uid 1000 → container uid 0, size 1
  ── inside the new user namespace ──
  uid=0  euid=0  gid=0
  uid_map: 0 1000 1

  ── what is this 'root' worth? ──
  DENIED   append to /etc/shadow    permission denied
  ALLOWED  mount a tmpfs
  ALLOWED  write into $HOME
```

Three lines of output, three lessons, and they arrive together: root is denied on a host file, allowed to mount, and allowed everywhere your own UID already reached. The probe output is the whole audit section of the talk arriving eighteen minutes early, say so, and tell them you will come back to it.

**Beat 3: why one ID is not enough.**

```console
$ ./nsdemo 3
  parent: asking for 0 → 100000, size 65536
  cmd.Start(): fork/exec /proc/self/exe: operation not permitted
```

The failure comes from `cmd.Start()`, not from the child, because writing `uid_map` happens during process creation. Widening the map beyond your own ID needs `CAP_SETUID` in the parent namespace, which you do not have. This is the payoff for `/etc/subuid` in §6.1.

**Beat 4: what Podman really does.** The best beat, and the reason the Go version beats the shell version outright.

```console
$ ./nsdemo 4
  parent: child is pid 48213, currently unmapped
  ── inside the new user namespace ──
  uid=65534  uid_map: (empty, no mapping exists)

  child: blocked, the map does not exist yet
  parent: newuidmap 48213 0 1000 1 1 100000 65536
  parent: newgidmap 48213 0 1000 1 1 100000 65536

  ── same process, one moment later, after newuidmap ran ──
  uid=0  euid=0  gid=0
  uid_map: 0 1000 1  |  1 100000 65536
```

The dance: start the child unmapped, hold it at a pipe barrier, call the setuid helpers, release it. That is precisely what Podman, Docker rootless and runc do, and you cannot show it with `unshare` at all, `podman unshare` hides it behind a black box.

Then the line that closes the segment, and the reason to build the beat this way:

> Same PID. No `setuid` call. Nothing about this process changed. **The kernel changed its mind about who it is.** That is all rootless is: a text file, a setuid helper, and a kernel that has always cared about namespaces more than it cared about UID 0.

*Timing discipline:* four commands, roughly 3:30 when it runs clean. If dry runs drift past four minutes, cut Beat 3 to a slide. Beat 4's error path implies it.

#### Why Go rather than `unshare`

You could do beats 1–3 with `unshare --user` and `unshare --user --map-root-user`. Both versions work; they are not equivalent.

| | Shell (`unshare`) | Go (`nsdemo`) |
|---|---|---|
| Setup cost | none, on any Linux box | one binary to build and copy |
| Beat 4 (paused child + newuidmap) | not possible, `podman unshare` hides it | shown exactly as runc does it |
| What the audience reads | output only | output, and four lines of syscall arguments |
| Repertoire fit | generic | connects to your container-runtime talk; the room has a GoLang track |
| Risk on stage | typos in a live shell | none, one binary, one argument |
| Failure mode | retype and hope | swap to the recording |

The deciding argument is Beat 4. The mechanism that makes rootless containers work, an unprivileged parent, a privileged helper, and a child held at a barrier, has no shell demonstration that is not a black box. `SysProcIDMap` is also self-documenting in a way `unshare` flags never are: `{ContainerID: 0, HostID: 1000, Size: 1}` *is* the `uid_map` line, in a struct.

Keep `unshare --user; id` in your back pocket as a one-line answer during Q&A.

### Demo 1. Who am I? (60 s)

The cold open showed the rootless half. Run it properly here, both panes, because it is the cleanest single frame in the talk:

```console
$ ./scripts/demo.sh 1
```

Inside the container, both panes print `uid=0(root)`. On the host, `/proc/<pid>/status` prints `Uid: 0` on the left and `Uid: 1000` on the right. Same image, same command, same claim, one of them is backed by the kernel and one of them is a translation.

Then the `uid_map` line, which only the right pane has at all. Rootful containers have no `uid_map` worth reading, because there is nothing to translate.

### Demo 2. Capabilities are real, but local (90 s)

Question on the slide: *If I have CAP_SYS_ADMIN, can I mount things?*

```console
$ ./scripts/demo.sh 2
```

Three commands, run identically in both panes: mount a tmpfs, mount a real disk, create a device node. Rootful does all three. Rootless mounts the tmpfs happily and fails the other two with `EPERM`. The explanation is one sentence: **tmpfs is flagged `FS_USERNS_MOUNT` in the kernel; ext4 is not.** The capability is genuine, the set of objects it applies to is not.

Follow-up jab, five seconds, no extra slide, the script runs it in both panes:

```console
# LEFT  (root):     mknod /dev/evil b 8 0   →  mknod: OK
$ RIGHT (rootless): mknod /dev/evil b 8 0   →  mknod: DENIED
```

`--privileged` in rootless mode does not mean what people think it means. On the left it means what it says. On the right it grants everything your user namespace has, which is still not device access.

### Demo 3: The blast radius (2 min), **the money demo**

Question on the slide: *Same mistake, two outcomes.*

Run the classic misconfiguration in both panes at once:

```console
$ ./scripts/demo.sh 3
```

First it reads:

```
LEFT  (root):      root:$6$....
RIGHT (rootless):  cat: can't open '/host/shadow': Permission denied
```

Then, and this is the half that makes it stick, it writes:

```
-v /etc:/host  →  echo "backdoor::0:0::/root:/bin/sh" >> /host/passwd
LEFT:   WROTE
RIGHT:  DENIED
```

Let the left pane sit for a second. That host now has an unauthenticated root account, from one flag. (The script deletes the line immediately afterwards, but say out loud that on a real host, nothing would have noticed.)

Rootful: host compromised. Rootless: `Permission denied`, because the container's UID 0 is host UID 1000, you, and `/etc/passwd` is owned by host UID 0.

The line to say: **the mount succeeded in both panes. Rootless did not stop the mistake, it made the mistake survivable.** That distinction is the talk in one sentence, so pause on it.

Then flip it, because honesty is the whole framing. The script's third step reads a fake credentials file out of a directory you own, and it succeeds in *both* panes. Your own data is inside the blast radius either way, and for a developer laptop or a CI runner your own data *is* the crown jewels.

### Demo 4. Limits that are not limits (90 s)

Question on the slide: *What did I lose?*

```console
$ ./scripts/demo.sh 4
```

The script asks for `--memory=64m` and then asks the container itself what limit it got, by reading `/sys/fs/cgroup/memory.max` from inside. `67108864` means the limit applied. Anything else means it did not, and on a cgroup v1 rootless host you get a warning rather than an error, which is the dangerous part.

Show it on a v2 host with delegation first (both panes agree), then on a v1 host (the panes diverge). If a second VM is too much on the day, this is the demo to pre-record.

If a second VM is too much work on the day, this is the demo to pre-record.

### Demo 5. The cost of userspace networking (90 s, pre-recorded)

`iperf3` through `pasta`, through `slirp4netns`, and rootful with a veth pair. **Measure these on your own hardware and put your own numbers on the slide**, quoting someone else's benchmark in a talk about honesty would be a bad look. Table skeleton in §8.

Add the second, smaller surprise:

```console
$ podman run --rm -p 80:80 nginx
Error: rootlessport cannot expose privileged port 80
```

and the fix, with its caveat: `sysctl net.ipv4.ip_unprivileged_port_start=80` lowers it for *every* process on the box, not just yours.

### Demo hygiene

- Everything runs from a local VM. No conference wi-fi in the path. No image pulls on stage, pre-pull and show `podman images` at the start.
- Every demo has a `reset` function that runs before it, so any demo can be skipped or repeated out of order.
- `asciinema` recording of every demo, embedded in the slides, one keystroke away. If a demo fails twice, switch to the recording and keep talking. Never debug on stage.
- Practise the whole demo block with the laptop on battery, on a beamer, at 1080p.

---

## 8. The audit, three slides that people will photograph

### 8.1 What rootless genuinely fixes

| Threat | Why rootless helps |
|---|---|
| `docker` group = root | No privileged daemon, no root-owned socket |
| Escape to host root (e.g. CVE-2019-5736, runc `/proc/self/exe` overwrite) | The host runtime binary is not writable by your UID |
| Leaked file descriptor escapes (e.g. CVE-2024-21626) | Reachable files are only those your UID can already reach |
| Careless `-v /:/host` | Reads and writes are bounded by your UID's own permissions |
| Multi-tenant CI runners on shared hosts | Each job runs under a distinct UID range |

### 8.2 What it does not touch

| Threat | Still fully exposed |
|---|---|
| Kernel LPE | One shared kernel. A kernel bug is a kernel bug. |
| Your own data | SSH keys, cloud credentials, source code, the `~/.kube/config` |
| Egress abuse | Crypto-mining, data exfiltration, joining a botnet, no root needed |
| Supply-chain compromise | A malicious image runs as you, and that is enough |
| Denial of service | Especially where cgroup limits are silently ignored (Demo 4) |
| Lateral movement to other containers of the same user | Same UID range, same everything |

### 8.3 The twist, rootless *adds* an attack surface

This is the section that will be quoted, so give it room and state it carefully.

To run rootless containers, the host must permit unprivileged user namespace creation. That single permission hands every local user `CAP_SYS_ADMIN` over a namespace, and thereby reachability into kernel subsystems that were previously only reachable by root. A long line of local privilege-escalation bugs (in the filesystem-mount parser, in nf_tables, in overlayfs copy-up) were exploitable by unprivileged users **only because** they could enter a user namespace first.

Which is why the distributions started pushing back:

```console
# Debian family, the blunt instrument
$ sysctl kernel.unprivileged_userns_clone

# Ubuntu 23.10+, the surgical one
$ sysctl kernel.apparmor_restrict_unprivileged_userns
```

The uncomfortable conclusion, said plainly: **the feature that makes rootless containers possible is the same feature hardening guides tell you to disable.** Rootless containers do not eliminate risk; they trade a large, well-understood risk for a smaller, sharper, newer one. Whether that trade is good depends on your threat model, and if you cannot state your threat model, you cannot answer the question.

No exploit code, no weaponised demo. Show the sysctls, show the capability surface, and let the audience draw the line.

---

## 9. Real world, where this lands in 2026

### Kubernetes finally caught up

User namespaces reached GA in **Kubernetes v1.36 (23 April 2026)**, on KEP-127 (open since 2016), alpha in 1.25, beta in 1.30, on by default from 1.33, GA in 1.36. Persistent volumes under a user namespace only became practical once idmapped mounts landed in kernel 5.12, which is the callback to §6.2.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: honest-pod
spec:
  hostUsers: false          # ← one field, open since 2016
  containers:
    - name: app
      image: ghcr.io/example/app:1.0
```

The caveats worth naming, because this is a talk about reality:

- The node needs a modern kernel and a CRI runtime that supports the CRI extension.
- Volume types and some device configurations remain restricted.
- Each pod consumes a UID range on the node, which caps pods per node.
- The container still runs as UID 0 *inside*. Everything you already knew about `runAsNonRoot` still applies.

### The decision table, the last content slide

| Context | Verdict |
|---|---|
| CI/CD runners doing untrusted builds | **Strong yes.** Unprivileged image builds are the best argument rootless has. |
| Developer laptops | **Yes.** Loses almost nothing, removes the root daemon. |
| Shared multi-tenant hosts / HPC | **Yes**, with per-user subuid ranges. |
| Kubernetes worker nodes | **Yes, with `hostUsers: false`**, as of 1.36, no longer exotic. |
| Network-heavy workloads (proxies, ingress, high PPS) | **Measure first.** Userspace networking may cost you more than the security is worth. |
| Anything needing host devices, GPUs, low ports, exotic CNI | **Probably not.** Fight the platform, lose the platform. |
| A host where you must disable unprivileged user namespaces | **No, and that is a coherent position.** |

---

## 10. Takeaways slide

Five lines. No animation. Leave it up through Q&A.

1. Root is a capability set evaluated against a user namespace, not UID 0.
2. Rootless shrinks the blast radius. It does not prevent the mistake.
3. Your own data is inside the blast radius. For a laptop or a CI runner, that is most of what matters.
4. What you lose is concrete: cgroup limits on old hosts, network throughput, low ports, some devices.
5. Rootless requires unprivileged user namespaces, which is itself an attack surface. Know your threat model, then choose.

Closing line, then stop talking:

> Rootless is not a security button. It is a smaller, sharper blast radius, bought with real trade-offs. That is a good deal, as long as you know you are making it.

---

## 11. Slide plan

38 slides for 30 minutes. One idea per slide. No slide with more than 20 words except the four tables.

| # | Slide | Type | Notes |
|---|---|---|---|
|, | *(cold open, no slide)* | live | tmux, two panes |
| 1 | Title + name + handle | title | 10 seconds |
| 2 | "Both answers are true" | statement | the thesis |
| 3 | The three claims | list | the roadmap |
| 4 | UID 0 is a shorthand | statement | |
| 5 | Capability check diagram | diagram | the credential model |
| 6 | User namespaces own other namespaces | statement | the key insight |
| 7 | The condo committee | analogy | photo + table |
| 8 | Five moving parts | agenda | numbered, reused as a progress marker |
| 9 | `/etc/subuid` | terminal | |
| 10 | container 0 → *you*, container 1+ → 100000+ | diagram | the two-line map; the live version comes in Demo A |
| 11 | newuidmap flow | diagram | the setuid path |
| 12 | "Two setuid-root binaries in your trust path" | statement | first uncomfortable truth |
| 13 | Filesystem: three eras | table | fuse-overlayfs → overlayfs → idmapped |
| 14 | Idmapped mounts, visually | diagram | plant the K8s seed |
| 15 | cgroups: three outcomes | table | |
| 16 | "A no-op with a warning" | statement | second uncomfortable truth |
| 17 | Network: veth vs tap | diagram | side by side |
| 18 | What userspace networking costs | list | five bullets |
| 19 | pasta vs slirp4netns | table | |
| 20 | Runtime + no root daemon | diagram | the `docker` group point |
| 21 | DEMO A title card, "What did Podman actually do?" | demo | full screen, single pane |
| 22 | Why a Go program can't unshare itself | code | 4 lines of `SysProcAttr` + the threading reason |
| 23 | "Full capabilities, held by nobody" | statement | after beat 1 |
| 24 | Three fields are the whole model: `0 1000 1` | statement | after beat 2 |
| 25 | "The kernel changed its mind about who it is" | statement | closes beat 4 |
| 26 | DEMO 2 title card | demo | capabilities |
| 27 | Why tmpfs and not ext4 | statement | `FS_USERNS_MOUNT` |
| 28 | DEMO 3 title card | demo | blast radius |
| 29 | "It made the mistake survivable" | statement | the money line |
| 30 | DEMO 4 title card | demo | cgroups |
| 31 | DEMO 5 title card | demo | networking, pre-recorded |
| 32 | Your measured numbers | table | fill in from §12 |
| 33 | What rootless fixes | table | photographed |
| 34 | What it does not touch | table | photographed |
| 35 | The twist: userns as attack surface | statement | |
| 36 | The two sysctls | terminal | |
| 37 | K8s 1.36 GA + `hostUsers: false` | code | |
| 38 | The decision table + five takeaways | table | stays up for Q&A, merge these two so the last slide carries both |

**Visual direction:** dark background, one accent colour, real terminal screenshots rather than styled code blocks, and no stock photography of shipping containers. Red for the rootful pane, green for the rootless pane, used consistently from slide 1 so the audience reads the colour before the text.

---

## 12. Numbers to measure before the talk

Fill this in from your own lab. Do not ship the talk with anyone else's benchmark.

| Measurement | Rootful | Rootless (pasta) | Rootless (slirp4netns) |
|---|---|---|---|
| `iperf3` TCP throughput, container → host | | | |
| `iperf3` UDP, packets/s | | | |
| Round-trip latency, p99 | | | |
| Container start-to-`echo`, cold | | | |
| `podman build` of a 5-layer Go image, cold cache | | | |
| First-pull extract time, 1 GB image (native overlayfs vs fuse-overlayfs) | | | |

---

## 13. Q&A, prepare these five

1. *"Is rootless slower?"*. Networking and, on old kernels, storage. CPU-bound workloads are unaffected. Give your numbers.
2. *"Can I run Kubernetes rootless?"*. Two different questions. Rootless kubelet (usernetes, k3s rootless) is niche. Pods with `hostUsers: false` on a normal cluster is GA and boring, which is the good kind of answer.
3. *"Does rootless replace gVisor / Kata / Firecracker?"*. No. Rootless narrows what an escape yields; those change what an escape has to break through. They compose.
4. *"We disabled unprivileged user namespaces for hardening. Now what?"*. A coherent policy. You have chosen rootful with tight seccomp/AppArmor and no `docker` group. Say so out loud rather than pretending both are possible.
5. *"Is `--privileged` safe in rootless?"*. Safer, not safe. It grants everything your namespace holds. It still cannot reach host devices, and it still gives full access to everything your UID owns.

---

## 14. Preparation plan

| When | Task |
|---|---|
| T-8 weeks | Build the demo VM(s); pin kernel, Podman, Docker versions; record `podman info` for the appendix |
| T-6 weeks | Run the benchmarks in §12; the numbers shape slides 18, 19, 27 |
| T-5 weeks | First full draft of slides; cut to 38 |
| T-4 weeks | Script and record all asciinema fallbacks; rehearse Demo A until it runs in 3:30 from muscle memory |
| T-3 weeks | Dry run to a colleague who runs Docker rootful in production, they will find the weak claims |
| T-2 weeks | Timed run against the 11:00 and 22:00 checkpoints. Target 28:30 of content, not 30, every talk grows on stage |
| T-1 week | Freeze the VM image, snapshot it, copy it to a second laptop |
| T-1 day | Full run on the actual laptop, on battery, at 1080p, with the display mirrored |

---

## 15. Follow-on ideas

- **Blog series, three parts:** "Root is not UID 0", "The chown tax and how idmapped mounts killed it", "Rootless as an attack surface". The last one is the piece with the most reach.
- **Companion talk:** *Build a rootless container in 100 lines of Go*, `clone(CLONE_NEWUSER)`, write `uid_map`, `pivot_root`, exec. Reuses your existing container-runtime material and would fit KCD or GopherCon.
- **Workshop version, 90 minutes:** the same content with attendees on their own VMs, ending with each participant breaking out of a rootful container and failing to break out of a rootless one.
- **Talk title for later:** *"The Sysctl That Hardening Guides and Container Runtimes Disagree About"*.

---

## Amendments made while building the repositories

The deck is authoritative wherever this document and `slides.md` disagree.
These are the differences found and how they were settled.

**The three amber slides are 12, 15 and 35**, not 12, 16 and 34. This
document and `SPEAKER-SCRIPT.md` were written against a 38 slide deck; the
deck is now 40 slides. Slide 16, the `--memory=512m` no-op, stays red rather
than amber, which is what the deck already did.

**Beat 1 cannot print a full `CapEff`.** Section 7 shows
`CapEff: 000001ffffffffff` for a namespace with no mapping. The kernel does
grant the full set at `clone` time, and that was verified with a C program
that clones and reads `/proc/self/status` without execing. But `execve`
recalculates capabilities, and an unmapped process has euid 65534 rather than
0, so the exec is treated as unprivileged and the effective set arrives empty.
The bounding set survives. Since Go must fork and exec, which is this talk's
own argument, beat 1 now prints `CapPrm`, `CapEff` and `CapBnd` together, and
the gap between the empty effective set and the full bounding set is the
point. Section 7 above and S23 in the deck have both been corrected to match,
so the old `CapEff: 000001ffffffffff` now survives only where it is quoted as
the thing that was wrong.

**Beat 4 no longer probes.** Its child execs while still unmapped, so it holds
no capabilities, and a later `uid_map` write does not restore them. Probing
there reported a weaker root than beat 2 for a reason that belongs in beat 1.
Podman and runc map first and exec second, which is the ordering this beat
cannot show.

**The AppArmor restriction is an allowlist, not a switch.** Section 8.3 frames
`kernel.apparmor_restrict_unprivileged_userns` as the feature hardening guides
disable. Measured on stock Ubuntu 24.04 with the sysctl at its shipped value
of 1, Podman still works, because `/etc/apparmor.d/podman` carries a `userns,`
rule, as do 90 other profiles on the image. What fails is `unshare` and the Go
program in the companion repository, neither of which has a profile. Ubuntu
did not break rootless containers, it permitted the runtimes it knows about
and refused everything else. The trade is real but it has moved: the allowlist
is the boundary now.

**No demo is pre-recorded.** Section 7 offers pre-recording for demo 5 and for
the cgroup v1 half of demo 4. Both now run live, demo 4's second half on a
second QEMU VM booted with `systemd.unified_cgroup_hierarchy=0`. There are no
asciinema recordings in either repository.

**Section 12's table is filled in.** See `docs/benchmarks.md` in the companion
repository. Read the ratios, not the absolute figures: they were measured in a
VM and there is no wire.

**Smaller corrections.** Section 1 cites "Demo 4 (`-v /etc`)" where that is
Demo 3, and refers to demos 6 and 7, which do not exist. Demo 5 points at
section 8 for the benchmark skeleton, which is in section 12.
