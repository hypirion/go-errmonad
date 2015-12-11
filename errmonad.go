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
func Bind(fn1, fn2 interface{}) interface{} {
	rfn1 := reflect.ValueOf(fn1)
	rtyp1 := rfn1.Type()
	if rtyp1.Kind() != reflect.Func {
		panic("Argument is not function")
	}
	if rtyp1.NumOut() == 0 {
		panic("Function must at least return one argument")
	}
	if rtyp1.Out(rtyp1.NumOut()-1) != errorType {
		panic("Last output argument must be of type error")
	}

	rfn2 := reflect.ValueOf(fn2)
	rtyp2 := rfn2.Type()
	if rtyp2.Kind() != reflect.Func {
		panic("Argument is not function")
	}
	if rtyp2.NumOut() == 0 {
		panic("Function must at least return one argument")
	}
	if rtyp2.Out(rtyp2.NumOut()-1) != errorType {
		panic("Last output argument must be of type error")
	}

	// Attempt to match types
	if rtyp1.NumOut()-1 != rtyp2.NumIn() {
		panic("Argument count mismatch")
	}

	for i := 0; i < rtyp2.NumIn(); i++ {
		if !rtyp1.Out(i).AssignableTo(rtyp2.In(i)) {
			panic(fmt.Sprintf("Cannot assign %s to %s", rtyp1.Out(i), rtyp2.In(i)))
		}
	}

	// Alright, let's attempt this thing.
	bindFn := func(in []reflect.Value) []reflect.Value {
		res := rfn1.Call(in)
		rerr := res[len(res)-1]
		if err, ok := rerr.Interface().(error); ok && err != nil {
			ret := []reflect.Value{}
			for i := 0; i < rtyp2.NumOut()-1; i++ {
				ret = append(ret, reflect.Zero(rtyp2.Out(i)))
			}
			return append(ret, rerr)
		}
		return rfn2.Call(res[:len(res)-1])
	}
	rtyp1In := []reflect.Type{}
	for i := 0; i < rtyp1.NumIn(); i++ {
		rtyp1In = append(rtyp1In, rtyp1.In(i))
	}
	rtyp2Out := []reflect.Type{}
	for i := 0; i < rtyp2.NumOut(); i++ {
		rtyp2Out = append(rtyp2Out, rtyp2.Out(i))
	}

	outrFn := reflect.FuncOf(rtyp1In, rtyp2Out, false)
	return reflect.MakeFunc(outrFn, bindFn).Interface()
}
