package updater

import "testing"

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int // sign
	}{
		{"v0.0.12-alpha", "0.0.11-alpha", 1},
		{"0.0.11-alpha", "0.0.11-alpha", 0},
		{"0.0.11", "0.0.11-alpha", 1},
		{"0.0.11-beta", "0.0.11-alpha", 1},
		{"0.0.11-rc.2", "0.0.11-rc.1", 1},
		{"0.1.0-alpha", "0.0.99", 1},
		{"0.0.10", "0.0.11-alpha", -1},
	}
	for _, c := range cases {
		got := compareVersions(c.a, c.b)
		if (got > 0) != (c.want > 0) || (got < 0) != (c.want < 0) {
			t.Errorf("compareVersions(%q,%q)=%d want sign %d", c.a, c.b, got, c.want)
		}
	}
}

func TestFindDebAsset(t *testing.T) {
	assets := []ghAsset{
		{Name: "DockSea_1.0.0_linux_arm64.deb"},
		{Name: "DockSea_1.0.0_linux_amd64.deb"},
		{Name: "DockSea_1.0.0_linux_amd64.deb.sha256"},
	}
	a := findDebAsset(assets)
	if a == nil {
		t.Fatal("expected asset")
	}
}
