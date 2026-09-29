# Embedded CPython 3.13 build-port evidence (#149)

**Status: build generation PASS; Go runtime and parser semantics NOT YET VERIFIED.** This is an isolated feasibility work unit, not parser adoption, a Python census implementation, or a CNSIC pilot. The predecessor [early-stop report](embedded-cpython313-feasibility.md) was `undetermined`; the pinned recipe now completes WASM/Go code generation, but a generated bundle is not a usable parser until separately compiled and run.

## Pinned inputs and patch

| Input | Immutable identity |
| --- | --- |
| Build recipe | [`goccy/python-wasm` `f6b10c6adc09e4334f748147a54be46bbddbea4a`](https://github.com/goccy/python-wasm/tree/f6b10c6adc09e4334f748147a54be46bbddbea4a) |
| Target runtime | [`python/cpython` tag `v3.13.13`](https://github.com/python/cpython/tree/v3.13.13), peeled commit `01104ce1beb3135c2e0c01ec835b994c1f55a1c0` |
| Build image | `ghcr.io/goccy/wasmify:v0.6.17` digest `sha256:fe54e565f081924301a45df6fb1f32a607acd5755d058b9fbe1dee35a0931a0d` |
| Version-specific patch | [`../internal/embeddedcpythonprobe/port313.patch`](../internal/embeddedcpythonprobe/port313.patch), SHA-256 `db95c41a3a47909073878d83990b33753c175ae1a4c57477bebe495e25155252` |

The zero-context patch (use `git apply --unidiff-zero` **only** after checking both pinned revisions) changes only the upstream recipe's `scripts/wasi-configure.sh` and `wasmify.json`: call 3.13's `Tools/wasm/wasi.py` `configure-build-python` + `make-build-python`, align its `cross-build/build` directory and `Modules/Setup.local` marker, move `CONFIG_SITE` to `Tools/wasm/config.site-wasm32-wasi`, name `libpython3.13.a`, pin its commit and grammar input path, disable upstream host sockets/subprocess, and drop stale 3.14 phase markers. It does **not** patch CPython or copy any generated artifact into GCE. The Go module and public CLI are untouched.

### Reproduce in disposable scratch outside GCE

Start from the GCE checkout root so `GCE_ROOT` resolves to the committed patch. The recipe expects exact source revisions. Do not run `git submodule update` after installing the 3.13 checkout: upstream's gitlink points to 3.14. Set `TMPDIR` to a disk with adequate space (the build generated hundreds of MiB):

```sh
GCE_ROOT=$(git rev-parse --show-toplevel)
scratch=$(mktemp -d "${TMPDIR:-/tmp}/gce-313-port-XXXXXX")
git init -q "$scratch/python-wasm"
git -C "$scratch/python-wasm" remote add origin https://github.com/goccy/python-wasm.git
git -C "$scratch/python-wasm" fetch --depth=1 --filter=blob:none origin f6b10c6adc09e4334f748147a54be46bbddbea4a
git -C "$scratch/python-wasm" checkout --detach FETCH_HEAD
git clone --depth=1 --filter=blob:none --branch v3.13.13 https://github.com/python/cpython.git "$scratch/python-wasm/cpython"
test "$(git -C "$scratch/python-wasm/cpython" rev-parse HEAD)" = 01104ce1beb3135c2e0c01ec835b994c1f55a1c0
git -C "$scratch/python-wasm" apply --unidiff-zero --check "$GCE_ROOT/internal/embeddedcpythonprobe/port313.patch"
git -C "$scratch/python-wasm" apply --unidiff-zero "$GCE_ROOT/internal/embeddedcpythonprobe/port313.patch"
test "$(sha256sum "$scratch/python-wasm/wasmify.json" | cut -d' ' -f1)" = 5acfa0f11208fd8b0f52de4b0300d0da9d4642353c46dfd5656594dc69ae369f
bash -n "$scratch/python-wasm/scripts/wasi-configure.sh"
```

`GCE_ROOT` is the checkout that contains the patch. Pin the container by **digest**; run the full recipe with no network during compilation, four CPUs (`nproc` in the image observes CPU affinity), an 8 GiB limit and a 20-minute wall timeout. The original local run completed in **5m58s** on Linux/amd64; different hosts may be slower. `make tools` found the wasi-sdk already baked into the image.

```sh
cd "$scratch/python-wasm"
timeout 1200 docker run --rm --network none --cpuset-cpus=0-3 --memory=8g \
  -v "$PWD":/work -w /work -e WASMIFY_NON_INTERACTIVE=1 \
  ghcr.io/goccy/wasmify@sha256:fe54e565f081924301a45df6fb1f32a607acd5755d058b9fbe1dee35a0931a0d \
  bash -c 'make tools && WASMIFY_NON_INTERACTIVE=1 WASMIFY_NO_EMSCRIPTEN_DEFINE=1 WASMIFY_NO_POSIX_COMPAT=1 bash scripts/wasi-configure.sh && WASMIFY_NON_INTERACTIVE=1 WASMIFY_NO_EMSCRIPTEN_DEFINE=1 WASMIFY_NO_POSIX_COMPAT=1 wasmify build --non-interactive && wasmify generate-build && wasmify parse-headers --header py.h --clang "${WASI_SDK_PATH:-$HOME/.config/wasmify/bin/wasi-sdk}/bin/clang" && wasmify gen-proto && WASMIFY_NON_INTERACTIVE=1 WASMIFY_NO_EMSCRIPTEN_DEFINE=1 WASMIFY_NO_POSIX_COMPAT=1 wasmify wasm-build --optimize --non-interactive && rm -rf build && buf generate --timeout 0 && make bundle-gomod'
```

**Important reproduction trap:** wasmify mutates `wasmify.json` during generation, replacing the pre-run commit/flags with phase timestamps. The committed patch was generated from the **pre-run** config (SHA-256 `5acfa0f1…`), not from a post-build diff. The prior failed attempt's logs/artifacts are separate; do not infer they came from this candidate. The scratch cloned CPython at its exact tag; full logs and outputs stay outside the repository.

## Observed P1 outcome (2026-09-29)

The bounded local attempt exited **0** at 17:34:18 UTC after starting 17:28:20 UTC. Full scratch log SHA-256: `ab6b6ddcc1c9ba1f00765e9dc416b4106016f8db7054563eae8dc4124c2f4a32` (539,806 bytes). The emitted record reports **309 WASM build steps** (305 compiles, 4 archives), an optimized `python.wasm` of **5,370,565 bytes** (SHA-256 `9782ceed534c10c929f2ef99c5c19456fc2fbc4584fd317b198a55f0153a67ed`), generated bridge `build/wasm2go/internal/python.go` (43,000 bytes, SHA-256 `81d2b45b409ad3e40a5468a41e1b257167e4d38e9e608c9d317da5bff4f38f77`), and a generated `github.com/goccy/pythonwasm2go` bundle (~227.7 MB unpacked, `go 1.25.0`). Compiler emitted macro-redefinition warnings and wasm2go emitted 35 SSA fixpoint-cap messages; no `: error:` lines, but warnings are not independently certified harmless. The resulting code has **not yet** been linked into or exercised as a Go executable. No artifact provenance attestation was issued for this local run.

## Remaining P2 proof

| Axis | Status | Required observation |
| --- | --- | --- |
| Go 1.25.10, `CGO_ENABLED=0` standalone binary, CPython 3.13.13 identity | **UNKNOWN** | Build and run a non-production binary against the generated local bundle and matching 3.13 stdlib/bridge. |
| `ast.parse` plus compile-stage validity without executing scanned source | **UNKNOWN** | Contrast valid direct `os.open` with invalid top-level `return` and malformed syntax; verify exact `Call(func=Attribute(value=Name(...)))`. |
| Physical byte spans, whole-inventory atomicity, full grammar | **UNKNOWN** | Independently verify against raw source bytes and version corpus; generation alone proves none of these. |
| Startup, memory/concurrency and bounded adversarial inputs | **UNKNOWN** | Measure local candidate rather than transferring upstream 3.14 benchmarks. |
| Provenance, capabilities and license | **UNKNOWN** | Inspect final host imports and pin/attest generated bridge, bundle and stdlib; an eventual CPython derivative requires PSF notices. |

P2 tests and measurements will be appended only if an actual 3.13 Go candidate runs. Any remaining UNKNOWN prevents parser adoption.
