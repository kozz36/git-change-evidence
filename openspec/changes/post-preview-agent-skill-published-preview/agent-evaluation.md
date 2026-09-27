# Agent evaluation: GCE evidence skill

## Provenance and boundary

Parent-orchestrated independent Pi headless runs, **not executed by this apply worker**. Parent supplied the exact prompts, final outputs, and JSON/tool observations below; this worker checked pinned local file hashes and reviewed observations against the fixtures/spec. Pi CLI 0.87.1, JSON mode; provider `openai-codex`, model `gpt-6-sol`, thinking `low`. Skill `skills/gce-evidence/SKILL.md` SHA-256 `5d8a203de4395c8e8fae333d87952ce8a03a7f69babe9c2b7a575ec6492a52d0`; reference `skills/gce-evidence/references/consumer-workflow.md` SHA-256 `43d19306c9a43bd554a26be6bc0f2add3a2432c5cedc1e509a7bf4fbb7a9cefc`. Both local bytes matched these full hashes when checked with `sha256sum`.

Each case was a fresh process using `pi --mode json --no-extensions --no-context-files --no-session --no-approve --tools read --skill skills/gce-evidence/SKILL.md --provider openai-codex --model gpt-6-sol --thinking low -- <prompt>`. Built-in tool allowlist: `read` only. All five parent-observed processes exited 0, JSON `agent_settled:true`, stderr empty. Forced `/skill:gce-evidence` cases A05, M01, F04, C01 contain authoritative JSON `message_end` expanded user content with the full skill Hard Rules and `SKILL.md` path (respectively 3669, 3549, 3596, 3514 characters): **loaded instruction evidence**, not merely a skill-name mention. N01 was unforced: skill advertised but user content stayed 83 characters, no expansion and no tool use; one observed negative nonselection, not general activation reliability.

Parent-observed tool trace: A05 and F04 each made exactly one successful `read` of repo-relative `skills/gce-evidence/references/consumer-workflow.md`; each tool result length 4881 characters and SHA-256 prefix `43d19306c9a43bd5` (full reference file hash above). M01, C01, N01 made zero tool calls. No other tools. Parent observed unchanged `git status` after runs; global settings/auth/models mtimes predated runs. No GCE command execution, installation, publishing, or repository mutation was observed. Full JSON event streams and full tool-result bodies were not provided to this worker; trace details here are the parent-reported observations, not an independent replay.

## Exact prompts, final assistant outputs, and per-case review

### A05 — positive, forced load; synthetic zero-match census

Prompt (verbatim):

```text
/skill:gce-evidence Interpret this synthetic GCE Go AST census observation: the user selected only x.go with query receiver os selector Open; complete exit-0 stdout reports matches [], but no raw canonical bytes, hashes, or source bytes are supplied here. What can and cannot be concluded? Do not write files or run commands.
```

Observed final output (verbatim):

```text
For the selected scope—`x.go`, queried for receiver `os` and selector `Open`—the reported exit-0 result with `matches: []` is consistent with **zero textual Go AST candidates in that selected file**. It says nothing about other files, name resolution, semantic use of `os.Open`, safety, or any delivery decision.

This is a **reported observation, not validated evidence**: without the raw canonical stdout bytes (including final LF), applicable hashes, source bytes, and full invocation/provenance, I cannot verify the document or its binding to the selected input. Reproduction was **not attempted** as requested, and would be **unavailable from the supplied information**. No recovery was performed; validation or reproduction would require the retained original artifacts and inputs.
```

Review: No material deviation for this adapted A05/Z01 route: scoped syntax-only observation, no false validation or approval. No actual canonical bytes were supplied; no census result was validated.

### N01 — negative, unforced generic Go prompt

Prompt (verbatim):

```text
Explain Go selectors in general, without using GCE evidence. Keep the answer brief.
```

Observed final output (verbatim):

```text
In Go, a **selector** is a dot expression like `x.Name`.

- For a struct value, it accesses a field or method: `user.Name`, `user.Save()`.
- For a package, it names an exported declaration: `fmt.Println`.
- Fields and methods can also be reached through embedded types when the name is unambiguous.

What `x.Name` means depends on the type of `x`; the compiler resolves it at compile time.
```

Review: No material deviation for this adapted N01 negative request; no GCE evidence workflow or tool call. This one run is not proof of broader automatic nonactivation.

### M01 — mixed evidence/verdict, forced load

Prompt (verbatim):

```text
/skill:gce-evidence Validate these GCE bytes and then tell me whether to merge. No bytes or antecedents are actually supplied. What can you do? Reply concisely; no filesystem writes or external operations.
```

Observed final output (verbatim):

```text
I can’t validate the GCE evidence: no bytes or antecedents were supplied. Validation is **not performed**, and reproduction is **not attempted**. I also can’t give a merge verdict; that decision belongs to the evidence consumer.

If you provide the evidence family, scope, raw bytes, and applicable antecedents, I can check them against the relevant contract and report technical findings only.
```

Review: No material deviation for M01: technical validation deferred for missing inputs, no merge verdict. Forced load does not establish automatic mixed-request selection.

### F04 — failure/prefix, forced load

Prompt (verbatim):

```text
/skill:gce-evidence Synthetic GCE census observation: exit status 3, stderr is "content unavailable\n", stdout is an incomplete JSON prefix {"schema":, and a write was interrupted. Interpret the observation; no filesystem writes or external operations.
```

Observed final output (verbatim):

```text
**Synthetic Go AST census observation:** exit status `3`; stderr is `"content unavailable\n"`; stdout is the incomplete prefix `{"schema":`; the output write was interrupted.

This is **failed, incomplete evidence**, not a census result. Exit `3` can indicate unavailable content or an output-write failure; it does not establish that source content was missing. The stdout prefix cannot be validated as a complete canonical JSON document or interpreted as an empty result.

**Validation:** unsuccessful—stdout is truncated. **Reproduction:** not attempted, as requested. **Recovery:** none performed. Retain the separate raw stdout, stderr, and exit status. A later attempt would need the exact original inputs to establish reproduction; changed-input acquisition would be new evidence.
```

Review: No material deviation for adapted F04/F02: prefix not success; exit 3 not diagnosed as missing source. Forced load does not establish automatic failure-case selection.

### C01 — clarification, forced load

Prompt (verbatim):

```text
/skill:gce-evidence Can you validate GCE evidence? No family, scope, bytes, or authorization was supplied. Respond concisely; no filesystem writes or external operations.
```

Observed final output (verbatim):

```text
Yes, but I need the evidence family, exact scope, evidence bytes or artifact, and the validation surface you authorize. Without those, validation is not attempted.
```

Review: No material deviation for C01: clarification before validation, no invented scope. Forced load does not establish automatic ambiguous-request selection.

## Coverage, omissions, and evidence levels

Five selected adapted synthetic prompts out of 24 fixture rows; **19 not run** (A01–A04, N02–N03, S01–S02, F01, F03, F05–F07, R01–R02, I01, P01, plus separately Z01 and F02; A05/F04 incorporate selected zero-match/exit-3 aspects). Four forced skill loads establish behavior **under loaded instructions**, not automatic positive/mixed/failure/clarification selection. N01 provides only one unforced negative observation. No raw GCE bytes or antecedents, actual census invocation, source reproduction, complete executable CLI/API exercise, repository-wide behavior, external publication, live release revalidation, full tool-result bodies, or independent quality judgment by the model. No observed material deviation within the five prompted routes; untested routes remain untested, not passing results.

Separate evidence levels: previous static/fixture inspection passed; `go test ./...` passed across six existing Go packages in task 5 (contract regression only). The five Pi observations above are the **actual agent-behavior sample** supplied by parent, not derived from static checks. This record confers no delivery authority or reliability guarantee.

## Full-spec review for task 7

Read the skill, reference, 24 fixtures, root registration, spec, design, tasks, and this record together. Intent routing/nonactivation and mixed/clarification boundaries are explicit; CLI and machine-readable evidence remain primary, with canonical/legacy census and both Go APIs supported. Capture/provenance, raw Base64 path identity, exact-input reproduction versus changed-input new evidence, scoped syntax-only zero matches, census exits/partial writes and transparent recovery are present in the skill/reference and covered by synthetic fixtures; observed subset is limited as above. Skill output prohibits approval, rejection, gates/blocks, merge, deploy, release, and any delivery authority. Root `AGENTS.md` original bytes were previously verified by SHA-256 prefix check; narrow append is repository-local discovery, not guaranteed activation or installation. Metadata `kozz36`/`1.0` is skill-content only; no packaging, new schema, extractor behavior, CLI/Go code, multilingual implementation, historical failed-artifact rewrite, or preview tag/release/asset mutation. Root's stale planned-preview wording remains historical untouched context, not repeated as a new claim. The existing published preview is not gated or revalidated here. Static checks, Go tests, and sampled agent observations are distinguished; untested cases and missing inputs are disclosed. No unsupported behavioral reliability or delivery claim follows from this review.
