// Copyright (c) 2015, Jean Niklas L'orange
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
// may be used to endorse or promote products derived from this software without
// specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package errmonad

import (
	"fmt"
	"reflect"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// Bind takes a sequence of functions and pipes the input
func Bind(fns ...interface{}) interface{} {
	if len(fns) == 0 {
		panic("Needs at least 1 function to bind over")
	}
	rfns := make([]reflect.Value, len(fns), len(fns))
	rtyps := make([]reflect.Type, len(fns), len(fns))
	for i, fn := range fns {
		rfns[i] = reflect.ValueOf(fn)
		rtyps[i] = rfns[i].Type()
		if rtyps[i].Kind() != reflect.Func {
			panic("Argument is not function")
		}
		if rtyps[i].NumOut() == 0 {
			panic("Function must at least return one argument")
		}
		if rtyps[i].Out(rtyps[i].NumOut()-1) != errorType {
			panic("Last output argument must be of type error")
		}
	}

	// attempt to match types
	for i := 0; i < len(rfns)-1; i++ {
		if rtyps[i].NumOut()-1 != rtyps[i+1].NumIn() {
			panic("Argument count mismatch")
		}
		for j := 0; j < rtyps[i+1].NumIn(); j++ {
			if !rtyps[i].Out(j).AssignableTo(rtyps[i+1].In(j)) {
				panic(fmt.Sprintf("Cannot assign %s to %s", rtyps[i].Out(j), rtyps[i+1].In(j)))
			}
		}
	}

	// Alright, let's attempt this thing.
	bindFn := func(in []reflect.Value) []reflect.Value {
		for _, fn := range rfns {
			res := fn.Call(in)
			rerr := res[len(res)-1]
			// error: short-circuit
			if err, ok := rerr.Interface().(error); ok && err != nil {
				ret := []reflect.Value{}
				for i := 0; i < rtyps[len(rtyps)-1].NumOut()-1; i++ {
					ret = append(ret, reflect.Zero(rtyps[len(rtyps)-1].Out(i)))
				}
				return append(ret, rerr)
			}
			in = res[:len(res)-1]
		}
		return append(in, reflect.Zero(errorType))
	}
	bindIn := []reflect.Type{}
	for i := 0; i < rtyps[0].NumIn(); i++ {
		bindIn = append(bindIn, rtyps[0].In(i))
	}
	bindOut := []reflect.Type{}
	for i := 0; i < rtyps[len(rtyps)-1].NumOut(); i++ {
		bindOut = append(bindOut, rtyps[len(rtyps)-1].Out(i))
	}

	outrFn := reflect.FuncOf(bindIn, bindOut, false)
	return reflect.MakeFunc(outrFn, bindFn).Interface()
}
