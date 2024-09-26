package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Progress struct {
	status int
}

func startProgress(title string) *Progress {
	fmt.Print(title + ": [" + strings.Repeat("-", 40) + "]" + strings.Repeat(string(rune(8)), 41))
	return &Progress{0}
}

func (p *Progress) progress(total, current int) {
	quantity := (100 * current) / total
	x := int(float64(quantity * 40 / 100))
	fmt.Print(strings.Repeat("#", x-p.status))
	p.status = x
}

func (p *Progress) end() {
	fmt.Println(strings.Repeat("#", (40-p.status)) + "]")
	p.status = 0
}

type IndefiniteLoadingBar struct {
	size      int
	direction int
	position  int
	finish    bool
	mu        sync.Mutex
}

func NewIndefiniteLoadingBar() *IndefiniteLoadingBar {
	return &IndefiniteLoadingBar{
		size:      40,
		direction: 1,
		position:  0,
		finish:    false,
	}
}

func (i *IndefiniteLoadingBar) progress() {
	fmt.Print("\r")
	spaces := strings.Repeat("-", (i.size - i.position))
	traces := strings.Repeat("-", i.position)
	fmt.Print("[" + traces + "#" + spaces + "]")
	i.position += i.direction
	if i.position >= i.size || i.position <= 0 {
		i.direction *= -1
	}
}

func (i *IndefiniteLoadingBar) end() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.finish = true
	fmt.Println()
}

func (i *IndefiniteLoadingBar) start() {
	for {
		i.mu.Lock()
		finish := i.finish
		i.mu.Unlock()
		if finish {
			break
		}
		i.progress()
		time.Sleep(10 * time.Millisecond)
	}
}
