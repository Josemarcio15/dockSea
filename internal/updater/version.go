package updater

import (
	"strconv"
	"strings"
)

type parsedVersion struct {
	nums     [3]int
	stage    int // alpha=1, beta=2, rc=3, stable=4
	stageNum int
}

func parseVersion(v string) parsedVersion {
	v = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(v, "v"), "V"))
	main, pre, hasPre := strings.Cut(v, "-")
	var p parsedVersion
	for i, part := range strings.SplitN(main, ".", 3) {
		p.nums[i], _ = strconv.Atoi(strings.TrimSpace(part))
	}
	p.stage = 4
	if hasPre {
		pre = strings.ToLower(pre)
		label := strings.TrimRight(pre, "0123456789.")
		p.stageNum, _ = strconv.Atoi(strings.Trim(pre[len(label):], "."))
		switch strings.Trim(label, ".") {
		case "alpha", "a":
			p.stage = 1
		case "beta", "b":
			p.stage = 2
		case "rc":
			p.stage = 3
		default:
			p.stage = 0
		}
	}
	return p
}

// compareVersions returns >0 if a is newer than b, <0 if older, 0 if equal.
// Understands semver-like tags with alpha/beta/rc suffixes (e.g. v0.0.11-alpha).
func compareVersions(a, b string) int {
	pa, pb := parseVersion(a), parseVersion(b)
	for i := 0; i < 3; i++ {
		if pa.nums[i] != pb.nums[i] {
			return pa.nums[i] - pb.nums[i]
		}
	}
	if pa.stage != pb.stage {
		return pa.stage - pb.stage
	}
	return pa.stageNum - pb.stageNum
}
