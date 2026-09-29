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

The bounded local attempt exited **0** at 17:34:18 UTC after starting 17:28:20 UTC. Full scratch log SHA-256: `ab6b6ddcc1c9ba1f00765e9dc416b4106016f8db7054563eae8dc4124c2f4a32` (539,806 bytes). The emitted record reports **309 WASM build steps** (305 compiles, 4 archives), an optimized `python.wasm` of **5,370,565 bytes** (SHA-256 `9782ceed534c10c929f2ef99c5c19456fc2fbc4584fd317b198a55f0153a67ed`), generated bridge `build/wasm2go/internal/python.go` (43,000 bytes, SHA-256 `81d2b45b409ad3e40a5468a41e1b257167e4d38e9e608c9d317da5bff4f38f77`), and a generated `github.com/goccy/pythonwasm2go` bundle (~227.7 MB unpacked, `go 1.25.0`). Compiler emitted macro-redefinition warnings and wasm2go emitted 35 SSA fixpoint-cap messages; no `: error:` lines, but warnings are not independently certified harmless. Generation alone did not prove an executable; P2 assembled and ran the matching bridge, bundle and stdlib below. No artifact provenance attestation was issued for this local run.

## P2: isolated Go 1.25.10 runtime smoke

Pinned [`goccy/go-python` v0.4.0](https://github.com/goccy/go-python/tree/v0.4.0) at `61b2ecf5dd19125d662abda3909d3d2bb520ebe8` supplied the host API **only**. Its checked-in 3.14 bridge/stdlib were replaced **in external scratch** with this build's generated bridge (hash above) and deterministic 3.13 stdlib zip SHA-256 `b35db8ab3b63c6a3bb3b27b4bf98bde467b2362745fd13b1e55901068846f623` (546 files). The generated bundle is a local Go module, not a GCE dependency. To reproduce after P1, set `$GCE_ROOT` as above and run:

```sh
git clone --depth=1 --branch v0.4.0 https://github.com/goccy/go-python.git "$scratch/go-python"
test "$(git -C "$scratch/go-python" rev-parse HEAD)" = 61b2ecf5dd19125d662abda3909d3d2bb520ebe8
python3 "$scratch/python-wasm/scripts/make-stdlib-zip.py" "$scratch/python-wasm/cpython/Lib" "$scratch/python-wasm/python_stdlib.zip"
cp "$scratch/python-wasm/build/wasm2go/internal/python.go" "$scratch/go-python/internal/python.go"
cp "$scratch/python-wasm/python_stdlib.zip" "$scratch/go-python/fs/stdlib.zip"
cp -R "$GCE_ROOT/internal/embeddedcpythonprobe/smoke" "$scratch/smoke"
docker run --rm --network none --cpuset-cpus=0-3 --memory=8g -v "$scratch":/work -w /work/smoke \
  -e GOPROXY=off -e GOSUMDB=off -e CGO_ENABLED=0 \
  golang@sha256:154bd7001b6eb339e88c964442c0ad6ed5e53f09844cc818a41ce4ecb3ce3b43 \
  bash -c 'go mod edit -replace=github.com/goccy/go-python=/work/go-python -replace=github.com/goccy/pythonwasm2go=/work/python-wasm/build/wasm2go/internal/internal/wasm2go && go mod tidy && go test -count=1 ./... && go build -p 1 -o embedded313 .'
docker run --rm --network none --read-only --cap-drop=ALL --memory=256m --cpuset-cpus=0-3 \
  -v "$scratch/smoke/embedded313":/probe:ro --entrypoint /probe \
  busybox@sha256:bdf57e528e45e4433820e045b29b4597825a1c9e38353532d90a01445013f82e
```

The committed nested smoke module has **no local replacement**: tests fail closed on upstream's 3.14 runtime if accidentally run without the generated 3.13 artifacts. Root `go test ./...` excludes it; the Docker command explicitly runs all four focused tests. TDD: before implementation, tests failed to compile (`undefined: runtimeVersion`, `undefined: inspect`; red log SHA-256 `1b5ac8c97a74c2b1ac8cdce3e89fd33bc67698cfa557e6b804156e396f377c1d`). After implementation all four pass (~0.36 s). The pinned exact-shape predicate has **1 site** in `smoke/main.go`; broadening it to accept an `Attribute` receiver yielded **2** candidates rather than **1**, killing `TestExactCallShapeNeverNestedAttributeOrIndexed` (mutation log SHA-256 `4ec6499dc095d1798a96bc88a427cb57803d9b0deee2c6d434abda2845fe2bfa`). Restored source passes. `ast.parse` accepts top-level `return 42`; compile rejects it as `SyntaxError`. A `raise` in scanned source remains unexecuted. 3.13 type alias and generic stub fixtures compile.

The `CGO_ENABLED=0` Go **1.25.10** binary is statically linked, **38,237,673 bytes** (SHA-256 `ce670bf826f9bd49793b88d380ba259488ec845134183fc233facd7a9a93f90b`). It ran in read-only, networkless BusyBox with no Python executable. A **single** run reported CPython `3.13.13` (dirty source-build suffix), `direct_calls=1`, startup **40.65 ms**, post-GC `HeapAlloc` **38,689,368 bytes** and boot `TotalAlloc` delta **52,787,456 bytes**. These are sample figures, not throughput or worst-case bounds.

| Contract axis | Evidence status |
| --- | --- |
| Go 1.25.10, no cgo/installed Python, correct embedded version | **PASS for this x86-64 smoke only.** No production dependency or cross-platform proof. |
| Parse + compile without executing scanned source; exact `Name.Attribute` | **PASS for four fixtures only.** No claim of exhaustive 3.13 grammar or all invalid forms. |
| Physical byte spans, multi-file atomicity, inventory/results | **UNKNOWN**; the probe is not a census implementation. |
| Startup and heap | **MEASURED once**; concurrency, adversarial size/depth and hard cancellation bounds **UNKNOWN**. |
| Capability isolation and supply-chain provenance | **UNKNOWN**; sockets/subprocess disabled in the recipe and Config{} denies host access, but final imports and isolation have not been audited. Local generated assets lack release attestations; CPython-derived bytes require PSF notices if distributed. |

**Independent audit:** a separate read-only Pi session matched frozen HEAD `c45ee4d`, patch/log/WASM/bridge/stdlib/binary and red/mutation hashes against scratch, and inspected the smoke code. Its readback could not prove exit status or elapsed time from a log alone; those are observations of the bounded parent command. It did not rerun the build or certify the remaining axes.

**Decision: candidate generation and smoke PASS, parser suitability UNKNOWN.** No dependency choice, census implementation, CNSIC pilot or delivery authority follows. The unproven axes require separate approval and independent tests before adoption.
