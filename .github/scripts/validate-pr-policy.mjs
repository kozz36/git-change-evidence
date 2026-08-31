import {pathToFileURL} from "node:url";

export const branchPattern = /^(feat|fix|chore|docs|style|refactor|perf|test|build|ci|revert)\/[a-z0-9._-]+$/;

const closingIssuePattern = /\b(?:close[sd]?|fix(?:e[sd])?|resolve[sd]?)\s+#(\d+)\b/gi;
const typeLabelPattern = /^type:(?:bug|feature|docs|refactor|chore|breaking-change)$/;

export function linkedIssueNumbers(body) {
  const numbers = new Set();
  for (const match of String(body).matchAll(closingIssuePattern)) {
    const number = Number(match[1]);
    if (Number.isSafeInteger(number) && number > 0) {
      numbers.add(number);
    }
  }
  return [...numbers];
}

export function validatePolicy({body, labels, branch, issue}) {
  const errors = [];
  if (linkedIssueNumbers(body).length !== 1) {
    errors.push("issue_linkage");
  }

  const typeLabels = labels.filter((label) => label.startsWith("type:"));
  if (typeLabels.length !== 1) {
    errors.push("type_label_count");
  } else if (!typeLabelPattern.test(typeLabels[0])) {
    errors.push("type_label_invalid");
  }

  if (!branchPattern.test(branch)) {
    errors.push("branch_name");
  }
  if (issue?.isPullRequest) {
    errors.push("issue_kind");
  } else if (issue?.state !== "open" || !issue.labels.includes("status:approved")) {
    errors.push("issue_approved");
  }
  return errors;
}

function labelsFromEnvironment(env) {
  const labels = JSON.parse(env.PR_LABELS_JSON ?? "[]");
  if (!Array.isArray(labels) || !labels.every((label) => typeof label === "string")) {
    throw new TypeError("PR_LABELS_JSON must be a JSON string array");
  }
  return labels;
}

async function fetchIssue(env, issueNumber, request) {
  if (!env.GITHUB_TOKEN) {
    throw new Error("GITHUB_TOKEN is required");
  }
  const endpoint = new URL(
    `/repos/${env.GITHUB_REPOSITORY}/issues/${issueNumber}`,
    env.GITHUB_API_URL ?? "https://api.github.com",
  );
  const response = await request(endpoint, {
    headers: {
      Accept: "application/vnd.github+json",
      Authorization: `Bearer ${env.GITHUB_TOKEN}`,
      "X-GitHub-Api-Version": "2022-11-28",
    },
  });
  if (!response.ok) {
    throw new Error(`issue lookup failed with HTTP ${response.status}`);
  }
  const issue = await response.json();
  return {
    isPullRequest: Boolean(issue.pull_request),
    state: issue.state,
    labels: issue.labels.map((label) => label.name),
  };
}

export async function runPolicy(env = process.env, request = fetch) {
  const issueNumbers = linkedIssueNumbers(env.PR_BODY ?? "");
  if (issueNumbers.length !== 1) {
    return ["issue_linkage"];
  }

  let labels;
  let issue;
  try {
    labels = labelsFromEnvironment(env);
    issue = await fetchIssue(env, issueNumbers[0], request);
  } catch {
    return ["policy_input_or_issue_lookup"];
  }
  return validatePolicy({
    body: env.PR_BODY ?? "",
    labels,
    branch: env.PR_HEAD_REF ?? "",
    issue,
  });
}

async function main() {
  const errors = await runPolicy();
  if (errors.length === 0) {
    process.stdout.write("PR policy is valid.\n");
    return;
  }
  for (const error of errors) {
    process.stderr.write(`PR policy failed: ${error}\n`);
  }
  process.exitCode = 1;
}

if (process.argv[1] && import.meta.url === pathToFileURL(process.argv[1]).href) {
  await main();
}
