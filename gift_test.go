package gift

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

type testFilter struct {
	z int
}

func (p *testFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	dstBounds = image.Rect(0, 0, srcBounds.Dx()+p.z, srcBounds.Dy()+p.z*2)
	return
}

func (p *testFilter) Draw(dst draw.Image, src image.Image, options *Options) {
	dst.Set(dst.Bounds().Min.X, dst.Bounds().Min.Y, color.Gray{123})
	return
}

func TestGIFT(t *testing.T) {
	g := New()
	if g.Parallelization() != defaultOptions.Parallelization {
		t.Error("unexpected parallelization property")
	}
	g.SetParallelization(true)
	if g.Parallelization() != true {
		t.Error("unexpected parallelization property")
	}
	g.SetParallelization(false)
	if g.Parallelization() != false {
		t.Error("unexpected parallelization property")
	}

	g = New(
		&testFilter{1},
		&testFilter{2},
		&testFilter{3},
	)
	if len(g.Filters) != 3 {
		t.Error("unexpected filters count")
	}

	g.Add(
		&testFilter{4},
		&testFilter{5},
		&testFilter{6},
	)
	if len(g.Filters) != 6 {
		t.Error("unexpected filters count")
	}
	b := g.Bounds(image.Rect(0, 0, 1, 2))
	if !b.Eq(image.Rect(0, 0, 22, 44)) {
		t.Error("unexpected gift bounds")
	}

	g.Empty()
	if len(g.Filters) != 0 {
		t.Error("unexpected filters count")
	}
	b = g.Bounds(image.Rect(0, 0, 1, 2))
	if !b.Eq(image.Rect(0, 0, 1, 2)) {
		t.Error("unexpected gift bounds")
	}

	g = &GIFT{}
	src := image.NewGray(image.Rect(-1, -1, 1, 1))
	src.Pix = []uint8{
		1, 2,
		3, 4,
	}
	dst := image.NewGray(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	if !dst.Bounds().Size().Eq(src.Bounds().Size()) {
		t.Error("unexpected dst bounds")
	}
	for i := range dst.Pix {
		if dst.Pix[i] != src.Pix[i] {
			t.Error("unexpected dst pix")
		}
	}

	g.Add(&testFilter{1})
	g.Add(&testFilter{2})
	dst = image.NewGray(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	if dst.Bounds().Dx() != src.Bounds().Dx()+3 || dst.Bounds().Dy() != src.Bounds().Dy()+6 {
		t.Error("unexpected dst bounds")
	}
	if dst.Pix[0] != 123 {
		t.Error("unexpected dst pix")
	}
}

func TestDrawAt(t *testing.T) {
	testDataGray := []struct {
		desc                     string
		filters                  []Filter
		pt                       image.Point
		op                       Operator
		srcb, dstb               image.Rectangle
		srcPix, dstPix0, dstPix1 []uint8
	}{
		{
			"draw at (Gray, [], -2, -2, copy)",
			[]Filter{},
			image.Pt(-2, -2),
			CopyOperator,
			image.Rect(-1, -1, 2, 2),
			image.Rect(-2, -2, 2, 2),
			[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]uint8{1, 2, 3, 0, 4, 5, 6, 0, 7, 8, 9, 0, 0, 0, 0, 0},
		},
		{
			"draw at (Gray, [], -1, -1, copy)",
			[]Filter{},
			image.Pt(-1, -1),
			CopyOperator,
			image.Rect(-1, -1, 2, 2),
			image.Rect(-2, -2, 2, 2),
			[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]uint8{0, 0, 0, 0, 0, 1, 2, 3, 0, 4, 5, 6, 0, 7, 8, 9},
		},
		{
			"draw at (Gray, [], 0, 0, copy)",
			[]Filter{},
			image.Pt(0, 0),
			CopyOperator,
			image.Rect(-1, -1, 2, 2),
			image.Rect(-2, -2, 2, 2),
			[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 0, 0, 4, 5},
		},
		{
			"draw at (Gray, [], 2, 2, copy)",
			[]Filter{},
			image.Pt(2, 2),
			CopyOperator,
			image.Rect(-1, -1, 2, 2),
			image.Rect(-2, -2, 2, 2),
			[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			"draw at (Gray, [], 0, -10, copy)",
			[]Filter{},
			image.Pt(0, -10),
			CopyOperator,
			image.Rect(-1, -1, 2, 2),
			image.Rect(-2, -2, 2, 2),
			[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			"draw at (Gray, [], -3, -3, copy)",
			[]Filter{},
			image.Pt(-3, -3),
			CopyOperator,
			image.Rect(-1, -1, 2, 2),
			image.Rect(-2, -2, 2, 2),
			[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]uint8{5, 6, 0, 0, 8, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			"draw at (Gray, [], -3, -3, over)",
			[]Filter{},
			image.Pt(-3, -3),
			OverOperator,
			image.Rect(-1, -1, 2, 2),
			image.Rect(-2, -2, 2, 2),
			[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]uint8{5, 6, 0, 0, 8, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			"draw at (Gray, [Resize], -2, -2, copy)",
			[]Filter{Resize(6, 6, NearestNeighborResampling)},
			image.Pt(-2, -2),
			CopyOperator,
			image.Rect(-1, -1, 2, 2),
			image.Rect(-2, -2, 2, 2),
			[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]uint8{1, 1, 2, 2, 1, 1, 2, 2, 4, 4, 5, 5, 4, 4, 5, 5},
		},
		{
			"draw at (Gray, [Resize], -3, -3, copy)",
			[]Filter{Resize(6, 6, NearestNeighborResampling)},
			image.Pt(-3, -3),
			CopyOperator,
			image.Rect(-1, -1, 2, 2),
			image.Rect(-2, -2, 2, 2),
			[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]uint8{1, 2, 2, 3, 4, 5, 5, 6, 4, 5, 5, 6, 7, 8, 8, 9},
		},
		{
			"draw at (Gray, [Resize], -1, -1, copy)",
			[]Filter{Resize(6, 6, NearestNeighborResampling)},
			image.Pt(-1, -1),
			CopyOperator,
			image.Rect(-1, -1, 2, 2),
			image.Rect(-2, -2, 2, 2),
			[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]uint8{0, 0, 0, 0, 0, 1, 1, 2, 0, 1, 1, 2, 0, 4, 4, 5},
		},
		{
			"draw at (Gray, [Resize], -1, -1, copy, empty)",
			[]Filter{Resize(6, 6, NearestNeighborResampling)},
			image.Pt(-1, -1),
			CopyOperator,
			image.Rect(0, 0, 0, 0),
			image.Rect(0, 0, 0, 0),
			[]uint8{},
			[]uint8{},
			[]uint8{},
		},
	}

	for _, d := range testDataGray {
		src := image.NewGray(d.srcb)
		src.Pix = d.srcPix

		g := New(d.filters...)

		dst := image.NewGray(d.dstb)
		dst.Pix = d.dstPix0

		g.DrawAt(dst, src, d.pt, d.op)

		if !checkBoundsAndPix(dst.Bounds(), d.dstb, dst.Pix, d.dstPix1) {
			t.Errorf("test [%s] failed: %#v, %#v", d.desc, dst.Bounds(), dst.Pix)
		}
	}

	testDataNRGBA := []struct {
		desc                     string
		filters                  []Filter
		pt                       image.Point
		op                       Operator
		srcb, dstb               image.Rectangle
		srcPix, dstPix0, dstPix1 []uint8
	}{
		{
			"draw at (NRGBA, [], 1, 1, over, 0% 100% alpha)",
			[]Filter{},
			image.Pt(1, 1),
			OverOperator,
			image.Rect(0, 0, 2, 2),
			image.Rect(0, 0, 3, 3),
			[]uint8{
				10, 20, 30, 255, 40, 50, 60, 255,
				100, 200, 0, 255, 0, 250, 200, 255,
			},
			[]uint8{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			[]uint8{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 10, 20, 30, 255, 40, 50, 60, 255,
				0, 0, 0, 0, 100, 200, 0, 255, 0, 250, 200, 255,
			},
		},
		{
			"draw at (NRGBA, [], 1, 1, over, 0% 50% alpha)",
			[]Filter{},
			image.Pt(1, 1),
			OverOperator,
			image.Rect(0, 0, 2, 2),
			image.Rect(0, 0, 3, 3),
			[]uint8{
				10, 20, 30, 127, 40, 50, 60, 127,
				100, 200, 0, 127, 0, 250, 200, 127,
			},
			[]uint8{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			[]uint8{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 10, 20, 30, 127, 40, 50, 60, 127,
				0, 0, 0, 0, 100, 200, 0, 127, 0, 250, 200, 127,
			},
		},
		{
			"draw at (NRGBA, [], 1, 1, over, 100% 50% alpha)",
			[]Filter{},
			image.Pt(1, 1),
			OverOperator,
			image.Rect(0, 0, 2, 2),
			image.Rect(0, 0, 3, 3),
			[]uint8{
				10, 20, 30, 128, 40, 50, 60, 128,
				100, 200, 0, 128, 0, 250, 200, 128,
			},
			[]uint8{
				0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255,
				0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255,
				0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255,
			},
			[]uint8{
				0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 0, 255,
				0, 0, 0, 255, 5, 10, 15, 255, 20, 25, 30, 255,
				0, 0, 0, 255, 50, 100, 0, 255, 0, 125, 100, 255,
			},
		},
		{
			"draw at (NRGBA, [], 1, 1, over, 100% 25% alpha)",
			[]Filter{},
			image.Pt(1, 1),
			OverOperator,
			image.Rect(0, 0, 2, 2),
			image.Rect(0, 0, 3, 3),
			[]uint8{
				20, 40, 80, 64, 40, 80, 120, 64,
				100, 200, 0, 64, 0, 100, 200, 64,
			},
			[]uint8{
				0, 0, 0, 255, 1, 2, 3, 255, 0, 0, 0, 255,
				0, 0, 0, 255, 40, 80, 120, 255, 40, 40, 40, 255,
				0, 0, 0, 255, 200, 200, 12, 255, 0, 0, 0, 255,
			},
			[]uint8{
				0, 0, 0, 255, 1, 2, 3, 255, 0, 0, 0, 255,
				0, 0, 0, 255, 35, 70, 110, 255, 40, 50, 60, 255,
				0, 0, 0, 255, 175, 200, 9, 255, 0, 25, 50, 255,
			},
		},
		{
			"draw at (NRGBA, [], 1, 1, over, shape)",
			[]Filter{},
			image.Pt(1, 1),
			OverOperator,
			image.Rect(0, 0, 2, 2),
			image.Rect(0, 0, 3, 3),
			[]uint8{
				100, 100, 100, 255, 100, 100, 100, 255,
				100, 100, 100, 255, 100, 100, 100, 0,
			},
			[]uint8{
				10, 10, 10, 255, 10, 10, 10, 255, 10, 10, 10, 255,
				10, 10, 10, 255, 10, 10, 10, 255, 10, 10, 10, 255,
				10, 10, 10, 255, 10, 10, 10, 255, 10, 10, 10, 255,
			},
			[]uint8{
				10, 10, 10, 255, 10, 10, 10, 255, 10, 10, 10, 255,
				10, 10, 10, 255, 100, 100, 100, 255, 100, 100, 100, 255,
				10, 10, 10, 255, 100, 100, 100, 255, 10, 10, 10, 255,
			},
		},
	}

	for _, d := range testDataNRGBA {
		src := image.NewNRGBA(d.srcb)
		src.Pix = d.srcPix

		g := New(d.filters...)

		dst := image.NewNRGBA(d.dstb)
		dst.Pix = d.dstPix0

		g.DrawAt(dst, src, d.pt, d.op)

		if !checkBoundsAndPix(dst.Bounds(), d.dstb, dst.Pix, d.dstPix1) {
			t.Errorf("test [%s] failed: %#v, %#v", d.desc, dst.Bounds(), dst.Pix)
		}
	}

}

type fakeDrawImage struct {
	r image.Rectangle
}

func (p fakeDrawImage) Bounds() image.Rectangle     { return p.r }
func (p fakeDrawImage) At(x, y int) color.Color     { return color.NRGBA{0, 0, 0, 0} }
func (p fakeDrawImage) ColorModel() color.Model     { return color.NRGBAModel }
func (p fakeDrawImage) Set(x, y int, c color.Color) {}

func TestSubImage(t *testing.T) {
	testData := []struct {
		desc string
		img  draw.Image
		ok   bool
	}{
		{
			"sub image (Gray)",
			image.NewGray(image.Rect(0, 0, 10, 10)),
			true,
		},
		{
			"sub image (Gray16)",
			image.NewGray16(image.Rect(0, 0, 10, 10)),
			true,
		},
		{
			"sub image (RGBA)",
			image.NewRGBA(image.Rect(0, 0, 10, 10)),
			true,
		},
		{
			"sub image (RGBA64)",
			image.NewRGBA64(image.Rect(0, 0, 10, 10)),
			true,
		},
		{
			"sub image (NRGBA)",
			image.NewNRGBA(image.Rect(0, 0, 10, 10)),
			true,
		},
		{
			"sub image (NRGBA64)",
			image.NewNRGBA64(image.Rect(0, 0, 10, 10)),
			true,
		},
		{
			"sub image (fake)",
			fakeDrawImage{image.Rect(0, 0, 10, 10)},
			false,
		},
	}

	for _, d := range testData {
		simg, ok := getSubImage(d.img, image.Pt(3, 3))
		if ok != d.ok {
			t.Errorf("test [%s] failed: expected %#v, got %#v", d.desc, d.ok, ok)
		} else if ok {
			simg.Set(5, 5, color.NRGBA{255, 255, 255, 255})
			r, g, b, a := d.img.At(5, 5).RGBA()
			if r != 0xffff || g != 0xffff || b != 0xffff || a != 0xffff {
				t.Errorf("test [%s] failed: expected (0xffff, 0xffff, 0xffff, 0xffff), got (%d, %d, %d, %d)", d.desc, r, g, b, a)
			}
		}
	}

}

func TestDraw(t *testing.T) {
	filters := [][]Filter{
		[]Filter{},
		[]Filter{Resize(2, 2, NearestNeighborResampling), Crop(image.Rect(0, 0, 1, 1))},
		[]Filter{Resize(2, 2, NearestNeighborResampling), CropToSize(1, 1, CenterAnchor)},
		[]Filter{FlipHorizontal()},
		[]Filter{FlipVertical()},
		[]Filter{Resize(2, 2, NearestNeighborResampling), Resize(1, 1, NearestNeighborResampling)},
		[]Filter{Resize(2, 2, NearestNeighborResampling), ResizeToFit(1, 1, NearestNeighborResampling)},
		[]Filter{Resize(2, 2, NearestNeighborResampling), ResizeToFill(1, 1, NearestNeighborResampling, CenterAnchor)},
		[]Filter{Rotate(45, color.NRGBA{0, 0, 0, 0}, NearestNeighborInterpolation)},
		[]Filter{Rotate90()},
		[]Filter{Rotate180()},
		[]Filter{Rotate270()},
		[]Filter{Transpose()},
		[]Filter{Transverse()},
		[]Filter{Brightness(10)},
		[]Filter{ColorBalance(10, 10, 10)},
		[]Filter{ColorFunc(func(r0, g0, b0, a0 float32) (r, g, b, a float32) { return 1, 1, 1, 1 })},
		[]Filter{Colorize(240, 50, 100)},
		[]Filter{ColorspaceLinearToSRGB()},
		[]Filter{ColorspaceSRGBToLinear()},
		[]Filter{Contrast(10)},
		[]Filter{Convolution([]float32{-1, -1, 0, -1, 1, 1, 0, 1, 1}, false, false, false, 0.0)},
		[]Filter{Gamma(1.1)},
		[]Filter{GaussianBlur(3)},
		[]Filter{Grayscale()},
		[]Filter{Hue(90)},
		[]Filter{Invert()},
		[]Filter{Maximum(3, true)},
		[]Filter{Minimum(3, true)},
		[]Filter{Mean(3, true)},
		[]Filter{Median(3, true)},
		[]Filter{Pixelate(3)},
		[]Filter{Saturation(10)},
		[]Filter{Sepia(10)},
		[]Filter{Sigmoid(0.5, 5)},
		[]Filter{Sobel()},
		[]Filter{UnsharpMask(1.0, 1.5, 0.001)},
	}

	for i, f := range filters {
		src := image.NewNRGBA(image.Rect(1, 1, 2, 2))
		src.Pix = []uint8{255, 255, 255, 255}
		g := New(f...)
		dst := image.NewNRGBA(image.Rect(-100, -100, -95, -95))
		g.Draw(dst, src)
		for x := dst.Bounds().Min.X; x < dst.Bounds().Max.X; x++ {
			for y := dst.Bounds().Min.Y; y < dst.Bounds().Max.Y; y++ {
				failed := false
				if x == -100 && y == -100 {
					if (dst.NRGBAAt(x, y) == color.NRGBA{0, 0, 0, 0}) {
						failed = true
					}
				} else {
					if (dst.NRGBAAt(x, y) != color.NRGBA{0, 0, 0, 0}) {
						failed = true
					}
				}
				if failed {
					t.Errorf("test draw pos failed: %d %#v %#v", i, f, dst.Pix)
				}
			}
		}
	}
}
