module.exports = {
  platform: "github",
  repositories: [
    process.env.RENOVATE_REPO
  ],
  extends: [
    ":ignoreModulesAndTests",
    ":enablePreCommit"
  ],
  pruneStaleBranches: false,
  branchPrefix: "chore/",
  branchNameStrict: true,
  semanticCommits: "enabled",
  semanticCommitScope: "dependencies",
  commitMessageLowerCase: "never",
  onboarding: true,
  onboardingBranch: "chore/onboard-renovate",
  onboardingCommitMessage: "🚀 Onboard Renovate",
  rangeStrategy: "pin",
  labels: ["renovate"],
  username: "renovate-bot",
  gitAuthor: "Renovate Bot <bot@renovateapp.com>",
  cachePrivatePackages: true,
  packageRules: [
    {
      matchPackagePatterns: [
        "*"
      ],
      semanticCommitType: "chore"
    },
    {
      matchDepTypes: [
        "dependencies",
        "require"
      ],
      semanticCommitType: "fix"
    },
    {
      matchManagers: [
        "github-actions"
      ],
      semanticCommitType: "ci",
    },
    {
      matchDatasources: [
        "docker"
      ],
      pinDigests: false
    },
  ],
  customManagers: [
    {
      customType: "regex",
      datasourceTemplate: "docker",
      managerFilePatterns: [
        "/(^|/)Chart\\.yaml$/"
      ],
      matchStrings: [
        "#\\s*renovate: image=(?<depName>.*?)\\s+appVersion:\\s*[\"']?(?<currentValue>[\\w+\\.\\-]*)"
      ]
    }
  ]
};
