package handle

import "testing"

func TestExtractPercentage(t *testing.T) {
	p, e := extractPercentage("[download]  29.5% of   74.90MiB at  152.04KiB/s ETA 05:55")
	if !e {
		t.Fatal(e)
	} else if p != 29.5 {
		t.Fatal(p)
	}
}
