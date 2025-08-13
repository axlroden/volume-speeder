package icon

import (
    "bytes"
    "image"
    "image/color"
    "image/png"
)

// Data contains PNG bytes for the tray icon. Generated procedurally at init time.
var Data = generateIconPNG()

func generateIconPNG() []byte {
    w, h := 24, 24
    img := image.NewNRGBA(image.Rect(0, 0, w, h))

    white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}

    set := func(x, y int, c color.NRGBA) {
        if x >= 0 && x < w && y >= 0 && y < h {
            img.SetNRGBA(x, y, c)
        }
    }

    // Speaker body
    for y := 9; y <= 15; y++ {
        for x := 3; x <= 7; x++ {
            set(x, y, white)
        }
    }
    for y := 9; y <= 15; y++ {
        dx := y - 12
        for x := 8; x <= 8+abs(dx); x++ {
            set(x, y, white)
        }
    }

    // Boost arrow and waves
    drawLine(set, 14, 16, 20, 10, white)
    drawLine(set, 20, 10, 20, 12, white)
    drawLine(set, 20, 10, 18, 10, white)

    drawArc(set, 11, 12, 4, white)
    drawArc(set, 12, 12, 6, white)

    var buf bytes.Buffer
    _ = png.Encode(&buf, img)
    return buf.Bytes()
}

func drawLine(set func(int, int, color.NRGBA), x0, y0, x1, y1 int, c color.NRGBA) {
    dx := abs(x1 - x0)
    sx := -1
    if x0 < x1 { sx = 1 }
    dy := -abs(y1 - y0)
    sy := -1
    if y0 < y1 { sy = 1 }
    err := dx + dy
    for {
        set(x0, y0, c)
        if x0 == x1 && y0 == y1 { break }
        e2 := 2 * err
        if e2 >= dy { err += dy; x0 += sx }
        if e2 <= dx { err += dx; y0 += sy }
    }
}

func drawArc(set func(int, int, color.NRGBA), cx, cy, r int, c color.NRGBA) {
    x, y := r, 0
    err := 0
    for x >= y {
        set(cx+x, cy+y, c)
        set(cx+y, cy+x, c)
        set(cx+x, cy-y, c)
        set(cx+y, cy-x, c)
        y++
        if err <= 0 { err += 2*y + 1 }
        if err > 0 { x--; err -= 2*x + 1 }
    }
}

func abs(v int) int { if v < 0 { return -v }; return v }
