package di

import (
	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/io"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/logic"
	"github.com/wirekang/mouseable/internal/overlay"
	"github.com/wirekang/mouseable/internal/view"
	"github.com/wirekang/mouseable/internal/winapi"
)

func Init() {
	logic.DI.SetCursorPos = func(x, y int) {
		winapi.SetCursorPos(x, y)
		lg.Logf("logic.DI.SetCursorPos(%d, %d)", x, y)
	}
	logic.DI.GetCursorPos = func() (x, y int) {
		x, y = winapi.GetCursorPos()
		lg.Logf("logic.DI.GetCursorPos() (%d, %d)", x, y)
		return
	}
	logic.DI.AddCursorPos = func(dx, dy int) {
		winapi.AddCursorPos(dx, dy)
		lg.Logf("logic.DI.AddCursorPos(%d, %d)", dx, dy)
	}
	logic.DI.MouseDown = func(button int) {
		winapi.MouseDown(button)
		lg.Logf("logic.DI.MouseDown(%d)", button)
	}
	logic.DI.MouseUp = func(button int) {
		winapi.MouseUp(button)
		lg.Logf("logic.DI.MouseUp(%d)", button)
	}
	logic.DI.Wheel = func(amount int, isHorizontal bool) {
		winapi.Wheel(amount, isHorizontal)
		lg.Logf("logic.DI.Wheel(%d, %v)", amount, isHorizontal)
	}
	logic.DI.OnActivated = func() {
		overlay.OnActivated()
		lg.Logf("logic.DI.OnActivated()")
	}
	logic.DI.OnDeactivated = func() {
		overlay.OnDeactivated()
		lg.Logf("logic.DI.OnDeactivated()")
	}
	winapi.DI.OnKey = func(keyCode uint32, isDown bool) (preventDefault bool) {
		preventDefault = logic.OnKey(keyCode, isDown)
		lg.Logf("winapi.DI.OnKey(%d(0x%X), %v) %v", keyCode, keyCode, isDown, preventDefault)
		return
	}
	winapi.DI.OnCursorMove = func(x, y int) {
		overlay.OnCursorMove(x, y)
		lg.Logf("winapi.DI.OnCursorMove(%d, %d)", x, y)
	}
	view.DI.LoadConfig = func() (config def.Config, err error) {
		config, err = io.LoadConfig()
		lg.Logf("view.DI.LoadConfig() %+v %+v", config, err)
		return
	}
	view.DI.SaveConfigJSON = func(json string) (err error) {
		err = io.SaveConfigJSON(json)
		lg.Logf("view.DI.SaveConfigJSON(%s) %+v", json, err)
		return
	}
	io.DI.SetConfig = func(config def.Config) {
		logic.SetConfig(config)
		overlay.SetConfig(config)
		lg.Logf("io.DI.SetConfig(%+v)", config)
	}
	view.DI.GetKeyText = func(keyCode uint32) (string, bool) {
		txt, ok := winapi.GetKeyText(keyCode)
		lg.Logf("view.DI.GetKeyText(%s, %v)", txt, ok)
		return txt, ok
	}
}
