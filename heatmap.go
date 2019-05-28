package onitamago

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
)

type BitboardHeatmap [5][5]uint // y, x

func (b *BitboardHeatmap) AddCard(card Card) {
	board := card.Bitboard() >> (CardOffset - BitboardCenterPos)
	b.add(board)
}

func (b *BitboardHeatmap) AddBoard(board Bitboard) {
	b.add(board)
}

func (b *BitboardHeatmap) clearMiddle() {
	b[2][2] = 0
}

func (b *BitboardHeatmap) add(board Bitboard) {
	horizontalMasks := [...]Bitboard{ // top -> bottom
		R5Mask, R4Mask, R3Mask, R2Mask, R1Mask,
	}
	verticalMasks := [...]Bitboard{ // left -> right
		AMask, BMask, CMask, DMask, EMask,
	}

	// TODO: speed up using popcount, this is just naive approach
	for y := range b {
		for x := range b[y] {
			r := horizontalMasks[y] | verticalMasks[x]
			if r == 0 {
			}
			if board&(horizontalMasks[y]&verticalMasks[x]) > 0 {
				b[y][x]++
			}
		}
	}
	b.clearMiddle()
}

func (b *BitboardHeatmap) Merge(m *BitboardHeatmap) {
	for y := range b {
		for x := range b[y] {
			b[y][x] += m[y][x]
		}
	}
	b.clearMiddle()
}

func (b *BitboardHeatmap) render(opacity float64) (img *image.NRGBA) {
	const H = 25 * 8
	const W = H
	const DotSize = H / len(b)

	img = imaging.New(W, H, color.NRGBA{255, 255, 255, 255})

	heat := imaging.New(DotSize, DotSize, color.NRGBA{10, 30, 205, 255})
	for y := range b {
		for x := range b[y] {
			img = imaging.Overlay(img, heat, image.Pt(x*DotSize, y*DotSize), opacity*float64(b[y][x]))
		}
	}

	piece := imaging.New(DotSize, DotSize, color.NRGBA{255, 255, 10, 255})
	img = imaging.Overlay(img, piece, image.Pt(2*DotSize, 2*DotSize), 1)

	lineColor := color.NRGBA{50, 50, 50, 255}
	verticalLine := imaging.New(1, H, lineColor)
	horizontalLine := imaging.New(W, 1, lineColor)
	for i := 1; i < 5; i++ {
		img = imaging.Overlay(img, verticalLine, image.Pt(i*DotSize, 0), 1)
		img = imaging.Overlay(img, horizontalLine, image.Pt(0, i*DotSize), 1)
	}

	return img
}

func (b *BitboardHeatmap) Render() (img *image.NRGBA) {
	return b.render(0.2)
}

func (b *BitboardHeatmap) RenderOneCard() (img *image.NRGBA) {
	return b.render(1)
}
