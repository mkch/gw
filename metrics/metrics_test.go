package metrics_test

import (
	"testing"

	"github.com/mkch/gw/metrics"
)

func TestMetrics(t *testing.T) {
	dip := metrics.Dip(100)
	px := dip.Px(192)

	if px.Value() != 200 {
		t.Errorf("Expected 200px, got %vpx", px.Value())
	}

	px = metrics.Px(100)
	dip = px.Dip(192)

	if dip.Value() != 50 {
		t.Errorf("Expected 50dip, got %vdip", dip.Value())
	}
}

func TestPackageHelper(t *testing.T) {
	px := metrics.ToPx(metrics.Dip(100), 192)
	if px.Value() != 200 {
		t.Errorf("Expected 200px, got %vpx", px.Value())
	}

	dip := metrics.ToDip(metrics.Px(100), 192)
	if dip.Value() != 50 {
		t.Errorf("Expected 50dip, got %vdip", dip.Value())
	}

	// nil Dimension should be treated as 0.
	px = metrics.ToPx(metrics.Dimension(nil), 192)
	if px != 0 {
		t.Errorf("Expected 0px, got %vpx", px.Value())
	}

	dip = metrics.ToDip(metrics.Dimension(nil), 192)
	if dip != 0 {
		t.Errorf("Expected 0dip, got %vdip", dip.Value())
	}
}
