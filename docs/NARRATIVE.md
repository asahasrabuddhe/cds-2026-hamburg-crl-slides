# The story

**The Reality of Rootless Containers**, ContainerDays 2026, Hamburg.

What you say is in [`SCRIPT.md`](SCRIPT.md) and in the deck's presenter notes.
What you do on the day is in [`EVENT-DAY.md`](EVENT-DAY.md). This document is
neither. It is the shape underneath both, so that when you lose your place you
can find it again by remembering what the talk is about rather than which slide
comes next.

---

## The story in one sentence

A word you have trusted your whole career turns out to be a shorthand, and once
you see what it is shorthand for, you can finally price what rootless costs you.

The word is root.

## What the audience should leave holding

They walked in believing rootless is a security button you switch on. They walk
out with this instead:

> Rootless did not stop the mistake. It made the mistake survivable.

Everything else in the talk exists to earn that sentence and to say honestly
what it does not cover.

---

## The three movements

The forty slides are three movements, not five sections. If you can name which
movement you are in, you can improvise your way back to the deck.

### One: the contradiction (cold open to S03, about two minutes)

Two witnesses give opposite testimony about the same process and both pass the
polygraph. The container swears it is root. The host swears it is `ajitem`.
Neither is lying.

You open with this before a single slide is up, because a contradiction is more
interesting than an agenda. The room does not yet know what a user namespace is
and does not need to. They only need to want the answer.

Then you stake your claim on S03, and it is the sentence that buys you
permission to be critical for the next twenty minutes:

> I run rootless. I also think most of what people believe about it is wrong.

**If you are lost here, say:** same process, two answers, both true.

### Two: the explanation (S04 to S25, about thirteen minutes)

The contradiction dissolves the moment you learn that the kernel stopped asking
"is this UID 0" in 1999. It asks whether the process holds a capability, and it
asks that against the namespace that owns the thing you are touching.

This movement has two halves and they do the same work twice, at different
depths.

**You explain it (S04 to S20).** The credential model, then the one structural
fact everything rests on, which is that a user namespace is the only namespace
that owns other namespaces. Then the condo committee, sixty seconds, because
nobody thinks a management committee owns the land. Then the five moving parts,
walked through in order.

**You prove it (S21 to S25).** Then you stop asserting and build the thing by
hand. This is the turn from lecture to demonstration, and it is why the talk is
not a blog post. Four beats, and the fourth is the peak of the whole talk:

> Same PID. No setuid call. The kernel changed its mind about who it is.

That line is the answer to the cold open, arriving thirteen minutes later, and
the room should feel the circuit close.

**If you are lost here, say:** root is a claim, and the kernel evaluates it
against a namespace.

### Three: the bill (S26 to S39, about twelve minutes)

Now that they believe you about the mechanism, you charge for it.

**The contrast (S26 to S32).** One script, two panes, and the only difference is
where you typed it. This is the most persuasive staging device in the talk
because it removes every "trust me, the other side does this." The peak is demo
3, where the left pane writes an unauthenticated root account into the host's
password file and the right pane cannot.

**The audit (S33 to S36).** What rootless genuinely fixes, what it does not
touch, and then the turn that the whole talk has been walking toward. It does
not merely fail to fix some things. It adds something.

**The choice (S37 to S39).** Kubernetes finally caught up, so this is not
exotic any more. Then a verdict table, then five things to remember, and you
stop talking.

**If you are lost here, say:** smaller blast radius, real bill, know which one
you are buying.

---

## The three motifs

Three things recur. They are what turn a list of subsystems into a story, and
they are the parts an audience remembers a week later.

### The two answers

Set up in the cold open, and every time it returns the room gets a small jolt of
recognition.

| Where | How it returns |
|---|---|
| Cold open | Container says root, host says `ajitem` |
| S10 | Container UID 0 is your own host UID, not the start of the range |
| S24 | `0 1000 1`, three fields, the entire security model |
| S28 | `WROTE` on the left, `DENIED` on the right |
| S29 | The mount succeeded in both panes |

Do not announce these as callbacks. Just let the shape repeat.

### The three uncomfortable truths

The amber and red slides, and the only ones that use those colours. By the third
one the room reads the colour before you speak, which is exactly what you want.

- **S12.** Two setuid-root binaries sit in your trust path. Rootless describes
  the runtime, not the installation.
- **S16.** `--memory=512m` is not an error. It is a no-op with a warning, and a
  container you believed was capped is not capped.
- **S35.** Rootless adds an attack surface.

They escalate. The first is a caveat, the second is a production hazard, the
third undermines the premise. Spacing them across the talk is what keeps it from
reading as either a sales pitch or a hit piece.

### The seed and the payback

You plant one thing early and collect it late, and it is the only long-range
setup in the deck.

On S14 you say idmapped mounts are the kernel feature that unblocked stateful
pods, and that you will come back to it. On S37 you do. That gap is seventeen
minutes, and paying it off tells the room the talk was built rather than
assembled.

---

## Where each slide sits in the story

| Slides | Beat | What it is doing |
|---|---|---|
| Cold open | The contradiction | Two true answers, no slide up |
| S01 to S03 | The claim | Your name, then the three claims and the right to be critical |
| S04 to S06 | The reveal | Root was never UID 0, and a userns owns other namespaces |
| S07 | The handhold | The condo committee, sixty seconds, first thing to cut |
| S08 | The map | Five moving parts, reused four more times |
| S09 to S12 | Identity | Two text files, a two-line map, and the first uncomfortable truth |
| S13 to S14 | Filesystem | Three eras, and the seed you collect on S37 |
| S15 to S16 | cgroups | Three outcomes, and the second uncomfortable truth |
| S17 to S19 | Network | No tap device, and what userspace costs |
| S20 | Execution | The `docker` group, and the strongest plain argument for rootless |
| S21 to S25 | The proof by hand | Lecture becomes demonstration, peaking on beat 4 |
| S26 to S29 | The proof by contrast | Same command, two panes, peaking on demo 3 |
| S30 to S32 | The measured cost | Limits that are not, networking that costs, your own numbers |
| S33 to S34 | The honest audit | What it fixes, what it does not touch |
| S35 to S36 | The turn | It adds a surface, and hardening guides restrict the feature it needs |
| S37 | The payback | Kubernetes, and the S14 seed collected |
| S38 to S39 | The choice | A verdict table, five takeaways, then stop |

---

## Why the demos are where they are

Each demo answers a question the slides just raised, which is why none of them
feel like an interlude.

- **Demo 1**, in the cold open, poses the contradiction rather than resolving
  it. It is the only demo that asks instead of answers.
- **Demo A** resolves it, and it has to come after the mechanism slides or the
  room has no vocabulary for what they are watching.
- **Demo 2** proves S27's claim that a capability is genuine while the set of
  objects it applies to is not.
- **Demo 3** is the money demo, and it is the one moment where the audience sees
  a real consequence rather than a printed value.
- **Demo 4** turns a slide's assertion about silent failure into something they
  watch fail silently.
- **Demo 5** exists mostly so the numbers on S32 are yours rather than quoted.

## The two sentences that have to land

Everything else can be paraphrased. These two are worth saying exactly:

> Same PID. No setuid call. The kernel changed its mind about who it is.

> Rootless did not stop the mistake. It made the mistake survivable.

The first closes the explanation. The second closes the talk. Pause after both.

## The closing move

You end on the five takeaways rather than a thank-you slide, and you leave that
slide up through questions so the room is looking at your argument rather than
at your handles while they think of one.

The last line, verbatim, and then stop talking:

> Not a security button. A smaller, sharper blast radius, bought with real
> trade-offs. That is a good deal, as long as you know you are making it.

That sentence is why the talk is not a hit piece. You spent thirty minutes
listing what rootless costs, and you land on recommending it anyway, with your
eyes open. An audience will forgive you almost anything if you are fair at the
end.

---

## If it all goes wrong

The demos can die and the talk still works, because the argument does not live
in the terminal. It lives in three sentences:

1. Root is a capability set evaluated against a namespace, not UID 0.
2. Rootless shrinks the blast radius. It does not prevent the mistake.
3. It needs a kernel feature that is itself an attack surface.

If every VM dies, say those three, walk the static slides that carry them, and
take questions early. Nobody in the room will know what they missed.
