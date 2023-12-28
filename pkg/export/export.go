// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package main

// Examples of private and public constructs.
// Read more in "Efficient Go"; Example 2-2.

const privateConst = 1
const PublicConst = 2

var privateVar int
var PublicVar int

func privateFunc() {}
func PublicFunc()  {}

type privateStruct struct {
	privateField int
	PublicField  int
}

func (privateStruct) privateMethod() {}
func (privateStruct) PublicMethod()  {}

type PublicStruct struct {
	privateField int
	PublicField  int
}

func (PublicStruct) privateMethod() {}
func (PublicStruct) PublicMethod()  {}

type privateInterface interface {
	privateMethod()
	PublicMethod()
}

type PublicInterface interface {
	privateMethod()
	PublicMethod()
}
