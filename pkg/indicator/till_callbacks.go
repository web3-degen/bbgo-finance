// Code generated by "callbackgen -type TILL"; DO NOT EDIT.

package indicator

import ()

func (inc *TILL) OnUpdate(cb func(value float64)) {
	inc.updateCallbacks = append(inc.updateCallbacks, cb)
}

func (inc *TILL) EmitUpdate(value float64) {
	for _, cb := range inc.updateCallbacks {
		cb(value)
	}
}
