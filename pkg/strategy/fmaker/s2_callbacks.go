// Code generated by "callbackgen -type S2"; DO NOT EDIT.

package fmaker

import ()

func (inc *S2) OnUpdate(cb func(value float64)) {
	inc.UpdateCallbacks = append(inc.UpdateCallbacks, cb)
}

func (inc *S2) EmitUpdate(value float64) {
	for _, cb := range inc.UpdateCallbacks {
		cb(value)
	}
}
