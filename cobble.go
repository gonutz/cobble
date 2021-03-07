package main

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"

	"github.com/gonutz/prototype/draw"
)

func main() {
	const windowW, windowH = 800, 600

	r := startRectangles()
	if loaded, err := loadRectangles(); err == nil {
		r = loaded
	}
	defer func() {
		saveRectangles(r)
	}()

	dragIndex := -1
	var dragDx, dragDy int
	wasLeftMouseDown := false
	draw.RunWindow("Cobble", windowW, windowH, func(window draw.Window) {
		if window.WasKeyPressed(draw.KeyEscape) {
			window.Close()
		}
		window.DrawText("Drag with left mouse.\nClick right while dragging to rotate.", 0, 0, draw.White)

		mouseX, mouseY := window.MousePosition()
		isLeftMouseDown := window.IsMouseDown(draw.LeftButton)
		defer func() { wasLeftMouseDown = isLeftMouseDown }()

		underMouse := -1
		for i, r := range r {
			window.FillRect(r.x, r.y, r.w, r.h, draw.White)
			window.DrawRect(r.x, r.y, r.w, r.h, draw.Gray)
			if r.contains(mouseX, mouseY) {
				underMouse = i
			}
		}

		if !wasLeftMouseDown && isLeftMouseDown {
			dragIndex = underMouse
			if dragIndex != -1 {
				r := r[dragIndex]
				dragDx = mouseX - r.x
				dragDy = mouseY - r.y
			}
		}
		if !isLeftMouseDown {
			dragIndex = -1
		}

		if dragIndex != -1 {
			r := &r[dragIndex]
			r.x, r.y = roundTo(mouseX-dragDx, 5), roundTo(mouseY-dragDy, 5)

			for _, c := range window.Clicks() {
				if c.Button == draw.RightButton {
					r.w, r.h = r.h, r.w
				}
			}
		}

		highlighted := dragIndex
		if highlighted == -1 {
			highlighted = underMouse
		}
		if highlighted != -1 {
			r := r[highlighted]
			window.FillRect(r.x, r.y, r.w, r.h, draw.DarkGreen)
		}
	})
}

type rectangle struct {
	x, y, w, h int
}

func rect(x, y, w, h int) rectangle {
	return rectangle{x: x, y: y, w: w, h: h}
}

func (r rectangle) contains(x, y int) bool {
	return r.x <= x && x < r.x+r.w &&
		r.y <= y && y < r.y+r.h
}

func roundTo(x, to int) int {
	return (x + to/2) / to * to
}

func startRectangles() []rectangle {
	var r []rectangle
	for i := 0; i < 20; i++ {
		r = append(r, rectangle{
			x: 150 + i%5*50,
			y: 150 + i/5*50,
			w: 50,
			h: 50,
		})
	}
	for i := 0; i < 4; i++ {
		r = append(r, rectangle{
			x: 400,
			y: 150 + i*50,
			w: 25,
			h: 50,
		})

		r = append(r, rectangle{
			x: 150 + i*50,
			y: 400,
			w: 50,
			h: 25,
		})
	}
	for i := 0; i < 2; i++ {
		r = append(r, rectangle{
			x: 150 + i*75,
			y: 350,
			w: 75,
			h: 50,
		})
	}
	for i := 0; i < 3; i++ {
		r = append(r, rectangle{
			x: 550,
			y: 150 + i*45,
			w: 40,
			h: 40,
		})
	}
	return r
}

const saveFile = "cobbles.rects"

var enc = binary.LittleEndian

func loadRectangles() ([]rectangle, error) {
	data, err := ioutil.ReadFile(saveFile)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(data)
	var rects []rectangle
	err = nil
	for err == nil {
		var x, y, w, h int32
		binary.Read(r, enc, &x)
		binary.Read(r, enc, &y)
		binary.Read(r, enc, &w)
		err = binary.Read(r, enc, &h)
		rects = append(rects, rect(int(x), int(y), int(w), int(h)))
	}
	return rects, nil
}

func saveRectangles(r []rectangle) {
	var buf bytes.Buffer
	w := &buf
	for _, r := range r {
		binary.Write(w, enc, int32(r.x))
		binary.Write(w, enc, int32(r.y))
		binary.Write(w, enc, int32(r.w))
		binary.Write(w, enc, int32(r.h))
	}
	ioutil.WriteFile(saveFile, buf.Bytes(), 0666)
}
