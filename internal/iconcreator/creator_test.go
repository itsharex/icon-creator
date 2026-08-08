package iconcreator

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestMinimumZoom(t *testing.T) {
	tests := []struct {
		width, height int
		want          float64
	}{
		{1024, 701, 701.0 / 1024.0},
		{701, 1024, 701.0 / 1024.0},
		{512, 512, 1},
		{0, 512, 1},
	}
	for _, test := range tests {
		if got := MinimumZoom(test.width, test.height); math.Abs(got-test.want) > 1e-9 {
			t.Errorf("MinimumZoom(%d, %d) = %f, want %f", test.width, test.height, got, test.want)
		}
	}
}

func TestRoundedIconExactFitPreservesWideSource(t *testing.T) {
	src := labeledSource(12, 6)
	got := RoundedIcon(src, 12, 0, MinimumZoom(12, 6), 0, 0)

	assertColorNear(t, got.NRGBAAt(0, 3), color.NRGBA{R: 255, A: 255})
	assertColorNear(t, got.NRGBAAt(11, 3), color.NRGBA{B: 255, A: 255})
	assertAlpha(t, got.NRGBAAt(6, 0), 0)
	assertAlpha(t, got.NRGBAAt(6, 11), 0)
}

func TestRoundedIconFillAndPanWideSource(t *testing.T) {
	src := labeledSource(12, 6)
	centered := RoundedIcon(src, 12, 0, 1, 0, 0)
	left := RoundedIcon(src, 12, 0, 1, 100, 0)
	right := RoundedIcon(src, 12, 0, 1, -100, 0)

	if centered.NRGBAAt(6, 6) == left.NRGBAAt(6, 6) || centered.NRGBAAt(6, 6) == right.NRGBAAt(6, 6) {
		t.Fatal("expected horizontal pan to change a wide source at fill zoom")
	}
	assertColorNear(t, left.NRGBAAt(0, 6), color.NRGBA{R: 255, A: 255})
	assertColorNear(t, right.NRGBAAt(11, 6), color.NRGBA{B: 255, A: 255})
}

func TestRoundedIconExactFitPreservesTallSource(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 6, 12))
	for y := 0; y < 12; y++ {
		c := color.NRGBA{R: uint8(20 * y), G: 70, A: 255}
		for x := 0; x < 6; x++ {
			src.SetNRGBA(x, y, c)
		}
	}
	got := RoundedIcon(src, 12, 0, MinimumZoom(6, 12), 0, 0)
	assertAlpha(t, got.NRGBAAt(0, 6), 0)
	assertAlpha(t, got.NRGBAAt(11, 6), 0)
	assertAlpha(t, got.NRGBAAt(6, 0), 255)
	assertAlpha(t, got.NRGBAAt(6, 11), 255)
}

func TestRoundedIconClampsBelowExactFit(t *testing.T) {
	src := labeledSource(12, 6)
	minimum := RoundedIcon(src, 12, 0, MinimumZoom(12, 6), 0, 0)
	below := RoundedIcon(src, 12, 0, 0.1, 0, 0)
	for y := 0; y < 12; y++ {
		for x := 0; x < 12; x++ {
			if minimum.NRGBAAt(x, y) != below.NRGBAAt(x, y) {
				t.Fatalf("below-fit zoom differs at %d,%d", x, y)
			}
		}
	}
}

func TestRoundedIconSquareSourceStillUsesOneAsMinimumZoom(t *testing.T) {
	src := labeledSource(12, 12)
	defaultResult := RoundedIcon(src, 12, 0, 1, 0, 0)
	belowOne := RoundedIcon(src, 12, 0, 0.5, 0, 0)
	for y := 0; y < 12; y++ {
		for x := 0; x < 12; x++ {
			if defaultResult.NRGBAAt(x, y) != belowOne.NRGBAAt(x, y) {
				t.Fatalf("square source below 100%% differs at %d,%d", x, y)
			}
		}
	}
}

func TestTransparentBackgroundBeforeLetterboxing(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 64, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			c := color.NRGBA{R: 248, G: 248, B: 246, A: 255}
			if x >= 24 && x < 40 && y >= 8 && y < 24 {
				c = color.NRGBA{R: 210, G: 54, B: 62, A: 255}
			}
			src.SetNRGBA(x, y, c)
		}
	}
	got := RoundedIconWithOptions(src, 64, 0, MinimumZoom(64, 32), 0, 0, true)
	assertAlpha(t, got.NRGBAAt(4, 20), 0)
	assertAlpha(t, got.NRGBAAt(32, 32), 255)
}

func labeledSource(width, height int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetNRGBA(x, y, color.NRGBA{R: uint8(255 * (width - 1 - x) / (width - 1)), B: uint8(255 * x / (width - 1)), A: 255})
		}
	}
	return img
}

func assertAlpha(t *testing.T, got color.NRGBA, want uint8) {
	t.Helper()
	if got.A != want {
		t.Fatalf("alpha = %d, want %d", got.A, want)
	}
}

func assertColorNear(t *testing.T, got, want color.NRGBA) {
	t.Helper()
	const tolerance = 30
	if absByte(got.R, want.R) > tolerance || absByte(got.G, want.G) > tolerance || absByte(got.B, want.B) > tolerance || got.A != want.A {
		t.Fatalf("color = %#v, want near %#v", got, want)
	}
}

func absByte(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}
