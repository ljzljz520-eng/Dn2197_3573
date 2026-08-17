# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
ok  	gradebook/audit	0.002s
?   	gradebook/cmd/gradebook	[no test files]
?   	gradebook/fixtures	[no test files]
ok  	gradebook/domain	0.001s
ok  	gradebook/grading	0.001s
ok  	gradebook/ranking	0.002s
--- FAIL: TestWorkflowConcurrentConfirmationKeepsBothOperators (0.00s)
    workflow_confirm_test.go:42: confirmation count = 1, want 2
    workflow_confirm_test.go:49: confirmed operators = []string{"alice"}, want two
FAIL
FAIL	gradebook/service	0.029s
ok  	gradebook/storage	0.013s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/gradebook): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/gradebook): exit `0`
