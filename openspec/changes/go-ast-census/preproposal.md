{
  "schema": "gentle-ai.sdd-preproposal/v1",
  "revision": 3,
  "exploration_reference": "openspec/changes/go-ast-census/explore.md",
  "exploration_read_status": "Reference retained from directly read preproposal input; exploration path was not a carried locator and was not accessed in this bounded continuation.",
  "retained_intent": "Official Go documentation only: verify Go 1.25 AST selector-call normalization including IndexExpr ambiguity, physical byte positions ignoring //line, full-file parse behavior and resource-limit caveats. No open-web. Product decisions confirmed; planning only.",
  "research_request": {
    "classes": ["documentation"],
    "questions": [
      "What AST nodes represent calls, selectors, parentheses, and ambiguous index/type-instantiation forms in Go 1.25?",
      "How are exact half-open byte spans and physical 1-based byte line/columns obtained without //line remapping, including UTF-8 and CRLF?",
      "What full-file parse error and parser object-resolution behavior supports a syntax-only all-or-nothing census?",
      "What can input/count limits guarantee and what parser-memory isolation cannot they guarantee?"
    ],
    "status": "done"
  },
  "admission": "done",
  "class_admission": {"documentation": "available"},
  "observed_grants": {"documentation": ["fetch_content"]},
  "admission_basis": "Actual runtime-injected documentation availability and callable fetch_content reflect host validation of active non-SDK tool provenance against the exact selected extension path. Official Go go1.25.10 source collection actually executed. No open-web selected or used.",
  "historical_denial": {"revision": 2, "outcome": "blocked", "context": "Earlier executor incorrectly required an additional registry verification API and persisted denial before clarification. Fresh human authorization allowed one same-scope continuation. Both actual revision2 JSON artifacts were read and host identity confirmation observed before evidence collection. Historical failed parent query remains recorded in research revision3; it is not evidence of current tool denial."},
  "evidence_references": [
    {"research_revision": 3, "claim": "Q1-A", "source_ids": ["S1", "S2"]},
    {"research_revision": 3, "claim": "Q2-A", "source_ids": ["S1", "S3", "S5"]},
    {"research_revision": 3, "claim": "Q3-A", "source_ids": ["S2", "S4"]},
    {"research_revision": 3, "claim": "Q4-A", "source_ids": ["S2", "S4"]}
  ],
  "research_reference": {"path": "/data/Projects/git-change-evidence-worktrees/go-ast-census/openspec/changes/go-ast-census/research.md", "revision": 3, "outcome": "done"},
  "product_decisions": {
    "status": "confirmed",
    "authority": "parent/user only",
    "scope": "User-approved full contract and acceptance in exploration. Future alignment of AGENTS.md and openspec/config.yaml layout.public_api ONLY confirmed; do not edit those now or alter testing rubric. No implementation/delivery authorization.",
    "preserved_constraints": [
      "Planning for the approved public census subpackage only; exact textual selector-call candidates, no import/type/alias/shadowing resolution or semantic-completeness claim.",
      "Accept existing InventoryDocumentV1, full supplied raw path/content records, exact versioned query and explicit positive limits. Validate inventory antecedent from owned public accessor state and canonical-byte digest without adding PolicyDocumentV1 or duplicating/round-tripping the inventory decoder.",
      "Validate both directions of the exact path set, uniqueness, content SHA-256 and byte lengths before any parsing. Zero/invalid inventory fails; a valid empty inventory succeeds empty.",
      "Unwrap only repeated ParenExpr, IndexExpr and IndexListExpr callee wrappers. IndexExpr ambiguity intentionally admits indexed-value calls as candidates. Require an identifier receiver and exact receiver/selector spelling.",
      "Use complete files and AST traversal only; errors or limits discard all results. Retain exact half-open call and normalized-selector byte spans and physical 1-based byte coordinates ignoring line directives, including UTF-8 and CRLF.",
      "Retain inventory digest, exact versioned query, extractor version, raw path, file digest and digest of the exact full-call source fragment. Sort by raw path bytes then call-start position and deterministic final position tie-breaker. Preserve defensive input/output ownership.",
      "Exact default bounds: 256 files, 1 MiB per file, 8 MiB total input and 10000 matches. No parser-memory isolation claim.",
      "Preserve the full exploration acceptance matrix, including decoys, repeated and legacy sites, wrapper ambiguity, Unicode/CRLF/multiline/line-directive positions, binding failures, empty-versus-zero inventory, deterministic ranges and external ownership tests. No implementation or tests are executed here.",
      "Future production ceiling is 400 novel lines; the session review budget is separately 400 changed lines. Test estimates create ask-on-risk delivery sizing concerns; no chain strategy or size exception is inferred.",
      "No store, DSL, CLI, serialization, adapter, mutation, runner, receipt, Python pilot, U10/U11, dependency resolution, lifecycle marker or delivery work is authorized by this research."
    ]
  },
  "proposal_ready": true,
  "readiness_scope": "All selected research questions have source-backed answers and product decisions remain parent-confirmed. This readiness is conditional on successful host-validated newer readback of both written artifacts. It permits parent gatekeeper consideration of planning only, not native proposal admission, implementation or delivery authority.",
  "risks": [
    "Official source evidence was not runtime-tested; no quantitative parser-memory bound was established.",
    "Retrieval times are a host-observed UTC interval rather than exact per-call timestamps.",
    "The separately referenced exploration was outside the two-locator read scope; no new independent verification of that file is claimed.",
    "Future review-budget risk remains governed by ask-on-risk."
  ],
  "skill_resolution": "none"
}
