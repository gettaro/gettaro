-- Backfill PR prefixes based on team pr_prefix values
-- This updates pull_requests where the title starts with a team's pr_prefix followed by a hyphen
-- Example: If team prefix is "WL", it will match PRs like "WL-123: Fix bug"
-- The matching is case-insensitive and matches the logic in the GitHub provider

UPDATE pull_requests pr
SET prefix = t.pr_prefix
FROM source_control_accounts sca
INNER JOIN teams t ON t.organization_id = sca.organization_id
WHERE pr.source_control_account_id = sca.id
  AND t.pr_prefix IS NOT NULL
  AND t.pr_prefix != ''
  AND pr.prefix IS NULL
  AND UPPER(pr.title) LIKE UPPER(t.pr_prefix) || '-%';

-- To backfill for a specific organization only, uncomment and use:
-- AND sca.organization_id = 'YOUR_ORGANIZATION_ID_HERE';

-- Note: If a PR title matches multiple team prefixes, only one will be set (non-deterministic which one)
-- Run multiple times if you need to catch PRs that were added after teams were created

