# Contributing

Thank you for your interest in this project. This document describes how to contribute to the specification and reference implementation.

## Project Phase

This repository is currently in **specification phase**. The full reference implementation is scheduled to land under NLnet Restack fund support (application planned for the fund's first open call). See [STATUS.md](./STATUS.md) and [ROADMAP.md](./ROADMAP.md) for the current state.

Contributions during the spec phase are particularly welcome for:

- Errors, ambiguities, or under-specified edge cases in [SPEC.md](./SPEC.md)
- Threat-model coverage gaps
- Comparisons with related work
- Test vector additions for the implemented portion (envelope-field encryption)
- Documentation improvements

We are not yet accepting large code contributions to the reference implementation, since the implementation roadmap is funded and managed against milestones. After milestone 3 is complete, this policy will relax.

## How to Propose a Specification Change

1. **Open an issue first.** Describe the problem, the proposed change, and the rationale. This avoids wasted effort on changes that turn out to be out of scope or incompatible with planned work.
2. **Wait for triage.** Maintainers will tag the issue and discuss within ~7 days.
3. **If accepted, open a pull request** with the proposed wording change. Quote section numbers and provide before/after text.
4. **Cryptographic changes require additional scrutiny** and may take longer to merge. Expect at least one round of review focused on correctness and threat-model implications.

## How to Propose a Code Change

For now (pre-milestone-3), code changes are limited to:

- Bug fixes to currently-implemented code (`EncryptField`, `DecryptField`, `NewSessionKey`)
- Improvements to test coverage, lint cleanups, build hygiene
- Documentation in `go/` (godoc comments, examples)

For these:

1. Fork the repository
2. Create a feature branch
3. Make the change
4. Ensure `go test ./...` passes
5. Ensure `go vet ./...` passes
6. Open a pull request describing the change and any test additions

## Code Style

- Standard Go formatting (`gofmt`)
- Standard Go vet hygiene (`go vet`)
- Comments on exported functions follow godoc conventions
- Avoid third-party dependencies unless necessary; the standard library is usually sufficient
- Cryptographic code MUST use vetted primitives from the standard library or audited libraries; ad-hoc cryptographic constructions will not be accepted

## Sign-off

By contributing, you agree that your contributions are licensed under the project's [Apache 2.0 license](./LICENSE) and that you have the right to make the contribution.

We use the [Developer Certificate of Origin](https://developercertificate.org/) (DCO) sign-off process. Please add a `Signed-off-by` line to each commit:

```
Signed-off-by: Your Name <your.email@example.com>
```

Use `git commit -s` to add this automatically.

## Communication

- Specification questions: GitHub issues on this repository
- Security reports: see [SECURITY.md](./SECURITY.md)
- General discussion: GitHub discussions (once enabled)

## Code of Conduct

Be respectful. Discussions should focus on technical merit. Personal attacks, harassment, or sustained disruption will result in moderation action.

---

Last updated: 2026-05-02
