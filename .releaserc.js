module.exports = {
  branches: [
    {
      name: "main"
    },
    {
      name: "feature/*",
      prerelease: "rc-${name.split('/').join('-').toLowerCase()}"
    }
  ],
  tagFormat: "${version}",
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/changelog",
    ["@semantic-release/exec", {
      "publishCmd": "goreleaser release --release-notes CHANGELOG.md --clean"
    }],
  ]
};
