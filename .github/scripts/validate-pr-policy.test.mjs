import assert from "node:assert/strict";
import {readFileSync} from "node:fs";
import test from "node:test";

import {
  branchPattern,
  linkedIssueNumbers,
  runPolicy,
  validatePolicy,
} from "./validate-pr-policy.mjs";

const policyCases = JSON.parse(
  readFileSync(new URL("./fixtures/pr-policy-cases.json", import.meta.url), "utf8"),
);

test("validates pull-request policy fixtures", () => {
  for (const testCase of policyCases) {
    assert.deepEqual(validatePolicy(testCase.input), testCase.want, testCase.name);
  }
});

test("extracts unique GitHub closing references", () => {
  assert.deepEqual(linkedIssueNumbers("Closes #7; fixes #8; closes #7"), [7, 8]);
  assert.deepEqual(linkedIssueNumbers("References #7"), []);
  assert.match("ci/7-go-pr-checks", branchPattern);
  assert.doesNotMatch("ci/Go-PR-Checks", branchPattern);
});

test("accepts every permitted branch prefix and descriptor punctuation", () => {
  for (const branch of ["feat/dot.name", "fix/under_score", "chore/repeated--hyphen", "docs/trailing-", "style/name", "refactor/name", "perf/name", "test/name", "build/name", "ci/name", "revert/name"]) {
    assert.match(branch, branchPattern, branch);
  }
});

test("loads approval state through an injected GitHub client", async () => {
  let requestedPath;
  const errors = await runPolicy(
    {
      PR_BODY: "Closes #7",
      PR_HEAD_REF: "ci/7-go-pr-checks",
      PR_LABELS_JSON: '["type:chore"]',
      GITHUB_API_URL: "https://api.github.com",
      GITHUB_REPOSITORY: "kozz36/git-change-evidence",
      GITHUB_TOKEN: "test-token",
    },
    async (url) => {
      requestedPath = url.pathname;
      return {
        ok: true,
        json: async () => ({state: "open", labels: [{name: "status:approved"}]}),
      };
    },
  );

  assert.deepEqual(errors, []);
  assert.equal(requestedPath, "/repos/kozz36/git-change-evidence/issues/7");
});
