# Embedded CPython 3.13 feasibility — bounded probe (#147)

**Assessment: UNDETERMINED; do not adopt.** The first pinned 3.13 build attempt stopped before compilation because the upstream adapter calls a CPython 3.14-only build helper. This is an integration blocker, **not** evidence that CPython 3.13 cannot be embedded. No 3.13 WASM, Go binding, static binary, AST output, or resource measurement was produced. The existing Go census and any CNSIC pilot are unchanged.

## Frozen candidate and boundary

| Component | Observed pin | Role |
| --- | --- | --- |
| `goccy/python-wasm` | `f6b10c6adc09e4334f748147a54be46bbddbea4a` | Build recipe and C++/Go bridge; its original submodule/config pin is CPython 3.14.6. |
| CPython | tag `v3.13.13`, peeled commit `01104ce1beb3135c2e0c01ec835b994c1f55a1c0` | Target grammar/runtime, detached checkout in external scratch. |
| `ghcr.io/goccy/wasmify:v0.6.17` | image digest `sha256:fe54e565f081924301a45df6fb1f32a607acd5755d058b9fbe1dee35a0931a0d` | Toolchain, pulled before the offline build attempt. |
| Host | Linux/amd64; Docker 29.8.1; image Go 1.25.14, Python 3.12.3 | **Not** proof of Go 1.25.10 compatibility; no Go executable was built. |

The recipe is a **modified external scratch checkout**, not a dependency of GCE. `wasmify.json` alone was changed: `upstream.commit` and library target `3.14` → `3.13`; `HostSockets`/`HostSubprocess` → `false`; stale generated phase markers removed. Its SHA-256 after the edit was `b4d928a37d791ffafb2fb171acdb35f99c5e30140731822160860abd600a4304`. No CPython source patch was applied, and no file from the scratch checkout was imported into this repository. The source clone itself pins the upstream commit even though the recipe's `.gitmodules` points to the existing submodule path.

### Reproduce the first stop

In fresh scratch **outside GCE**, fetch the exact build recipe commit and CPython tag, verify the peeled commit, then work from `python-wasm/`. Do not run `git submodule update`, which would restore upstream 3.14:

```sh
scratch=$(mktemp -d)
git init -q "$scratch/python-wasm"
git -C "$scratch/python-wasm" remote add origin https://github.com/goccy/python-wasm.git
git -C "$scratch/python-wasm" fetch --depth=1 --filter=blob:none origin f6b10c6adc09e4334f748147a54be46bbddbea4a
git -C "$scratch/python-wasm" checkout --detach FETCH_HEAD
git clone --filter=blob:none --depth=1 --branch v3.13.13 \
  https://github.com/python/cpython.git "$scratch/python-wasm/cpython"
test "$(git -C "$scratch/python-wasm/cpython" rev-parse HEAD)" = 01104ce1beb3135c2e0c01ec835b994c1f55a1c0
cd "$scratch/python-wasm"
```

Apply exactly these changes to `wasmify.json` (the last line removes stale 3.14 phase markers):

```sh
python3 - <<'PY'
import json
from pathlib import Path
p = Path('wasmify.json')
c = json.loads(p.read_text())
c['upstream']['commit'] = '01104ce1beb3135c2e0c01ec835b994c1f55a1c0'
c['targets'][0]['name'] = 'libpython3.13'
c['targets'][0]['build_target'] = 'libpython3.13.a'
c['targets'][0]['description'] = 'Static CPython 3.13 interpreter library (wasm32-wasi)'
c['user_selection']['target_name'] = 'libpython3.13'
c['bridge']['HostSockets'] = False
c['bridge']['HostSubprocess'] = False
c.pop('phases', None)
p.write_text(json.dumps(c, indent=2) + '\n')
PY
sha256sum wasmify.json  # expect b4d928a37d791ffafb2fb171acdb35f99c5e30140731822160860abd600a4304
```

The local command used `--network none` after pulling the pinned image, a four-CPU affinity and an 8 GiB container cap to avoid competing with host workloads:

```sh
docker run --rm --network none --cpuset-cpus=0-3 --memory=8g \
  -v "$PWD":/work -w /work -e WASMIFY_NON_INTERACTIVE=1 \
  ghcr.io/goccy/wasmify@sha256:fe54e565f081924301a45df6fb1f32a607acd5755d058b9fbe1dee35a0931a0d \
  bash -c 'make tools && WASMIFY_NON_INTERACTIVE=1 WASMIFY_NO_EMSCRIPTEN_DEFINE=1 WASMIFY_NO_POSIX_COMPAT=1 bash scripts/wasi-configure.sh'
```

The original attempt requested the **full upstream pipeline**, but `&&` stopped at this configure step (exit **2**, 2026-09-29 14:49 UTC). The shorter command above reproduces the same first stop. The observed output was:

```text
wasmify ensure-tools ./cpython --output-dir .
[ensure-tools] wasi-sdk: already installed
[ensure-tools] All tools ready.
== build triple: x86_64-pc-linux-gnu
== wasi sdk:     /root/.config/wasmify/bin/wasi-sdk
== building host build-python
python3: can't open file '/work/cpython/Tools/wasm/wasi': [Errno 2] No such file or directory
```

The full local log stayed **outside the repository** (`build-313-attempt1.log`, SHA-256 `15f0938037718fe63649e9aaddfb9c6578b9805ed848e8413876932f002a752b`). An independent tree inspection confirms CPython 3.13 instead has `Tools/wasm/wasi.py`, `Tools/wasm/wasm_build.py` and `Tools/wasm/config.site-wasm32-wasi`; it has no `Tools/wasm/wasi` path or `Tools/wasm/wasi/config.site-wasm32-wasi`. The 3.13 `wasi.py` **does** default to `wasm32-wasip1`, matching the borrowed script's host triple; the helper path and config-site layout are the proven mismatches. This is **version-specific build orchestration**, not a failed Python syntax test. An adapter port would need a separately budgeted compatibility pass across configure, paths, patches, bridge/API, and generated stdlib before an end-to-end build can be claimed.

## Contract matrix

| Axis from `docs/python-census-design.md` | Result | Evidence / missing proof |
| --- | --- | --- |
| Go 1.25.10, no installed Python/cgo at **runtime** | **UNKNOWN** | Upstream modules advertise Go 1.25, but this 3.13 candidate generated no binary. Build-time Python and C/C++ toolchains are distinct from runtime dependencies. |
| Exact `.py`/`.pyi` Python 3.13 grammar, `Call(func=Attribute(value=Name(...)))` | **UNKNOWN** | No 3.13 embedding built. `ast.parse` alone would not enforce compile-stage validity: CPython 3.13 accepts an AST for top-level `return 42`, then compilation rejects it. A future extractor must compile without executing scanned source. |
| Physical UTF-8 source-byte call/function spans, Unicode/CRLF | **UNKNOWN** | CPython AST documents UTF-8 byte columns and exclusive ends; converting them to absolute bytes and verifying original slices is not tested here. |
| Invalid/unsupported syntax, atomic multi-file result | **UNKNOWN** | No parser bridge, no `.pyi` or malformed fixture executed, and no partial-output path implemented. |
| Limits, adversarial size/depth, cancellation, heap/stack and concurrency | **UNKNOWN** | No local measurements. Upstream `go-python` v0.4.0 reports roughly 52 MiB per live instance and 10.5 ms startup on a different machine/runtime (**3.14.6**, Apple M5); these are not GCE benchmarks or budgets. |
| Isolation, capability reduction, distribution and licenses | **UNKNOWN** | Socket/subprocess bridge flags were disabled in the scratch recipe, but no built artifact verifies the final imports or filesystem behavior. The recipe/image is large (downloaded image ~2.9 GB locally); no GCE binary size measured. A shipped CPython derivative would require the PSF license/attribution and verified generated-artifact provenance. |

**Conclusion:** `undetermined`, not `infeasible` or `feasible`. The first integration mismatch is reproducible. A follow-up would be a **newly scoped build-port work unit**, not a silent third-party dependency adoption or an enlargement of #147. The owner must decide whether its cost is justified before any parser adoption or Python census implementation.

References: [CPython 3.13 AST parsing versus compilation](https://docs.python.org/3.13/library/ast.html#ast.parse), [`goccy/python-wasm` build inputs](https://github.com/goccy/python-wasm/tree/f6b10c6adc09e4334f748147a54be46bbddbea4a), [CPython `v3.13.13` WASM build tools](https://github.com/python/cpython/tree/v3.13.13/Tools/wasm), [`goccy/go-python` v0.4.0 resource/licensing notes](https://github.com/goccy/go-python/tree/v0.4.0).
