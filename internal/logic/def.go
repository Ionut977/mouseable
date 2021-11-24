package logic

import (
	"math"

	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/hook"
)

type logicState struct {
	fixedSpeed      float64
	speedX, speedY  float64
	steppingMap     map[*logicDef]struct{}
	wasCursorMoving bool
}

type logicDef struct {
	function *def.Function
	onStart  func(state *logicState)
	onStep   func(state *logicState)
	onStop   func(state *logicState)
}

var logicDefs = []*logicDef{
	{
		function: def.MoveRight,
		onStep: func(s *logicState) {
			s.speedX += dataMap[def.Acceleration]
		},
	},
	{
		function: def.MoveUp,
		onStep: func(s *logicState) {
			s.speedY -= dataMap[def.Acceleration]
		},
	},
	{
		function: def.MoveLeft,
		onStep: func(s *logicState) {
			s.speedX -= dataMap[def.Acceleration]
		},
	},
	{
		function: def.MoveDown,
		onStep: func(s *logicState) {
			s.speedY += dataMap[def.Acceleration]
		},
	},
	{
		function: def.ClickLeft,
		onStart: func(_ *logicState) {
			DI.MouseDown(0)
		},
		onStop: func(_ *logicState) {
			DI.MouseUp(0)
		},
	},
	{
		function: def.ClickRight,
		onStart: func(_ *logicState) {
			DI.MouseDown(1)
		},
		onStop: func(_ *logicState) {
			DI.MouseUp(1)
		},
	},
	{
		function: def.ClickMiddle,
		onStart: func(_ *logicState) {
			DI.MouseDown(2)
		},
		onStop: func(_ *logicState) {
			DI.MouseUp(2)
		},
	},
	{
		function: def.WheelUp,
		onStep: func(_ *logicState) {
			DI.Wheel(int(dataMap[def.WheelAmount]), false)
		},
	},
	{
		function: def.WheelDown,
		onStep: func(_ *logicState) {
			DI.Wheel(-int(dataMap[def.WheelAmount]), false)
		},
	},
	{
		function: def.SniperMode,
		onStart: func(s *logicState) {
			s.fixedSpeed = dataMap[def.SniperModeSpeed]
		},
		onStop: func(s *logicState) {
			s.fixedSpeed = 0
		},
	},
	{
		function: def.Flash,
		onStart: func(s *logicState) {
			if math.Abs(s.speedX) < 0.5 && math.Abs(s.speedY) < 0.5 {
				return
			}
			distance := dataMap[def.FlashDistance]
			var dx int32
			var dy int32
			angle := math.Atan2(s.speedX, s.speedY)
			dx = int32(distance * math.Sin(angle))
			dy = int32(distance * math.Cos(angle))
			hook.AddCursorPos(dx, dy)
		},
	},
}
