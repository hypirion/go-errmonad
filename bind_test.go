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

package errmonad_test

import (
	"encoding/json"
	"fmt"
	monad "gopkg.in/hyPiRion/go-errmonad.v1"
)

const MaxBananaCount = 90 // Max amount of bananas in a box

type BananaBox struct {
	Type    string
	Bananas int
}

// Double takes a banana crate and returns a new banana crate where the amount
// of bananas in it is doubled. If the new banana count will be larger than the
// maximum banana count, this method will error.
func (bc BananaBox) Double() (BananaBox, error) {
	return bc.AddBananas(bc.Bananas)
}

// AddBananas add n bananas to a banana box. If the new banana count will be
// larger than the maximum banana count a single crate can contain, this method
// will error.
func (bc BananaBox) AddBananas(n int) (BananaBox, error) {
	if n+bc.Bananas > MaxBananaCount {
		return BananaBox{}, fmt.Errorf("Tried to add %d bananas to a box with %d bananas already inside it, will go over the limit", n, bc.Bananas)
	}
	bc.Bananas += n
	return bc, nil
}

func jsonBananaBox(bs []byte) (bb BananaBox, err error) {
	err = json.Unmarshal(bs, &bb)
	return
}

var doubleBananaBox = monad.Bind(
	jsonBananaBox,
	(BananaBox).Double,
	json.Marshal,
).(func([]byte) ([]byte, error))

var quadrupleBananaBox = monad.Bind(
	jsonBananaBox,
	(BananaBox).Double,
	(BananaBox).Double,
	json.Marshal,
).(func([]byte) ([]byte, error))

func Example() {
	examples := []string{
		`[]`,
		`{"Bananas": "0"}`,
		`{"Type": "Dwarf Cavendish", "Bananas": 41}`,
		`{"Type": "Grand Nain", "Bananas": 16}`,
	}
	conversions := []func([]byte) ([]byte, error){
		doubleBananaBox,
		quadrupleBananaBox,
	}
	for _, example := range examples {
		for _, convert := range conversions {
			bs, err := convert([]byte(example))
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(bs))
			}
		}
	}
	// Output:
	// json: cannot unmarshal array into Go value of type errmonad_test.BananaBox
	// json: cannot unmarshal array into Go value of type errmonad_test.BananaBox
	// json: cannot unmarshal string into Go value of type int
	// json: cannot unmarshal string into Go value of type int
	// {"Type":"Dwarf Cavendish","Bananas":82}
	// Tried to add 82 bananas to a box with 82 bananas already inside it, will go over the limit
	// {"Type":"Grand Nain","Bananas":32}
	// {"Type":"Grand Nain","Bananas":64}
}
