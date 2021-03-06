package binutils

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestMarshalNumbers(t *testing.T) {
	s := uint32(1604)
	result := Marshal(&s)
	expected := []byte{0, 0, 6, 68}
	if bytes.Compare(result, expected) != 0 {
		t.Errorf("Result does not match expected, expected '%v' now '%v'", expected, result)
	}
}

func TestUnsupportedTypes(t *testing.T) {
	type intType struct {
		FieldA int
		FieldB uint
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expecting panic on unsupported type")
		}
	}()

	var result intType
	b := []byte{0, 0, 0, 13, 0, 0, 0, 50}
	Unmarshal(b, &result)

}

func TestMarshalString(t *testing.T) {
	testStr := "HelloWorld"
	expected := []byte{0, 0, 0, 10, 72, 101, 108, 108, 111, 87, 111, 114, 108, 100}
	result := Marshal(&testStr)
	if !bytes.Equal(result, expected) {
		t.Errorf("Result does not match expected, expected '%v' now '%v'", expected, result)
	}
}

func TestMarshalStruct(t *testing.T) {
	type sshptyRequest struct {
		Term    string
		Width   uint32
		Height  uint32
		PWidth  uint32
		PHeight uint32
	}

	original := sshptyRequest{Term: "xterm", Width: 80, Height: 24}
	expected := []byte{0, 0, 0, 5, 120, 116, 101, 114, 109, 0, 0, 0, 80, 0, 0, 0, 24, 0, 0, 0, 0, 0, 0,
		0, 0}
	result := Marshal(&original)
	if !bytes.Equal(result, expected) {
		t.Errorf("Result does not match expected, expected:\n '%v'\n now:\n '%v'", expected, result)
	}
}

func TestMarshalComplexStruct(t *testing.T) {
	type innerStruct struct {
		A string
		B uint64
	}
	type structA struct {
		A uint8
		B uint32
		C string
		D [3]byte
		E []innerStruct
	}

	v := structA{A: 7,
		B: 12,
		C: "Good2",
		D: [3]byte{5, 7, 9},
		E: []innerStruct{
			{A: "Hello",
				B: 17},
			{A: "Goodbye",
				B: 27},
			{A: "Helloe",
				B: 95},
		}}
	expected := []byte{7, 0, 0, 0, 12, 0, 0, 0, 5, 71, 111, 111, 100, 50,
		5, 7, 9, 0, 0, 0, 3,
		0, 0, 0, 5, 72, 101, 108, 108, 111, 0, 0, 0, 0, 0, 0, 0, 17,
		0, 0, 0, 7, 71, 111, 111, 100, 98, 121, 101, 0, 0, 0, 0, 0, 0, 0, 27,
		0, 0, 0, 6, 72, 101, 108, 108, 111, 101, 0, 0, 0, 0, 0, 0, 0, 95}
	result := Marshal(&v)
	if !bytes.Equal(result, expected) {
		PrintCompareByteArray(expected, result, t)
		t.Errorf("Result does not match expected, expected:\n '%v'(%v)\n now:\n '%v'(%v)", expected, len(expected), result, len(result))
	}
}

func PrintCompareByteArray(a, b []byte, t *testing.T) {
	for i, v := range a {
		t.Logf("A:%v B:%v", v, b[i])
	}
}

func TestSSHPtyRequest(t *testing.T) {

	type sshptyRequest struct {
		Term     string
		Width    uint32
		Height   uint32
		PWidth   uint32
		PHeight  uint32
		TermMode []byte
	}

	var result sshptyRequest
	expected := sshptyRequest{Term: "xterm", Width: 80, Height: 24, TermMode: []byte{3, 0, 0, 0, 127, 42, 0, 0, 0, 1, 128, 0, 0, 150, 0, 129, 0, 0, 150, 0, 0}}
	b := []byte{0, 0, 0, 5, 120, 116, 101, 114, 109, 0, 0, 0, 80, 0, 0, 0, 24, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 21, 3, 0, 0, 0, 127, 42, 0, 0, 0, 1, 128, 0, 0, 150, 0, 129, 0, 0, 150, 0, 0}
	Unmarshal(b, &result)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Result does not match expected, expected '%v' now '%v'", expected, result)
	}
}

func ExampleMarshal() {
	type sshptyRequest struct {
		Term     string
		Width    uint32
		Height   uint32
		PWidth   uint32
		PHeight  uint32
		TermMode []byte
	}

	var req sshptyRequest
	b := []byte{0, 0, 0, 5, 120, 116, 101, 114, 109, 0, 0, 0, 80, 0, 0, 0, 24, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 21, 3, 0, 0, 0, 127, 42, 0, 0, 0, 1, 128, 0, 0, 150, 0, 129, 0, 0, 150, 0, 0}
	Unmarshal(b, &req)
	fmt.Printf("Term:%v (%v x %v)", req.Term, req.Width, req.Height)
	// Output: Term:xterm (80 x 24)
}

func TestStructWithArrayFields(t *testing.T) {
	type arrayStruct struct {
		FieldA [3]byte
		FieldB [3]uint16
		FieldC [3]string
	}

	b := []byte{10, 17, 10, 0, 0, 0, 1, 1, 2, 0, 0, 0, 5, 72, 101, 108, 108, 111, 0, 0, 0, 4, 72, 101,
		108, 108, 0, 0, 0, 3, 72, 101, 108}
	var result arrayStruct
	expected := arrayStruct{
		FieldA: [3]byte{10, 17, 10},
		FieldB: [3]uint16{0, 1, 258},
		FieldC: [3]string{"Hello", "Hell", "Hel"},
	}
	Unmarshal(b, &result)
	if result != expected {
		t.Errorf("Result not expected: Expected: %v Actual %v", expected, result)
	}
}

func TestString(t *testing.T) {
	var result string
	b := []byte{0, 0, 0, 5, 72, 101, 108, 108, 111}

	Unmarshal(b, &result)
	if result != "Hello" {
		t.Errorf("Result not expected: Expected: %v Actual %v", "Hello", result)
	}
}

func TestAnonNestedStruct(t *testing.T) {
	type nestedStruct struct {
		FieldA uint16
		FieldB struct {
			FieldB1 uint16
			FieldB2 uint16
		}
		FieldC uint16
	}

	var result nestedStruct
	expected := nestedStruct{
		FieldA: 17,
		FieldB: struct {
			FieldB1 uint16
			FieldB2 uint16
		}{
			FieldB1: 13,
			FieldB2: 15},
		FieldC: 21,
	}
	b := []byte{0, 17, 0, 13, 0, 15, 0, 21}
	Unmarshal(b, &result)
	if result != expected {
		t.Errorf("Result not expected: Expected: %v Actual %v", expected, result)
	}
}

func TestArrayOfStruct(t *testing.T) {
	type innerStruct struct {
		FieldA1 int16
		FieldA2 bool
		FieldA3 int32
	}
	type arrayOfStruct struct {
		FieldA [3]innerStruct
	}

	expected := arrayOfStruct{FieldA: [3]innerStruct{
		innerStruct{FieldA1: 258, FieldA2: false, FieldA3: 50595078},
		innerStruct{FieldA1: 2828, FieldA2: true, FieldA3: 219025168},
		innerStruct{FieldA1: 5398, FieldA2: false, FieldA3: 387455258}}}
	b := []byte{1, 2, 0, 3, 4, 5, 6, 11, 12, 1, 13, 14, 15, 16, 21, 22, 0, 23, 24, 25, 26}
	var result arrayOfStruct
	Unmarshal(b, &result)
	if result != expected {
		t.Errorf("Result not expected: Expected: %v Actual %v", expected, result)
	}

}

func TestNestedStruct(t *testing.T) {
	type innerStruct struct {
		FieldA string
		FieldB uint16
	}
	type nestedStruct struct {
		FieldA uint16
		FieldB innerStruct
		FieldC uint16
	}

	var result nestedStruct
	expected := nestedStruct{FieldA: 17, FieldB: innerStruct{FieldA: "Hello", FieldB: 19}, FieldC: 50}
	b := []byte{0, 17, 0, 0, 0, 5, 72, 101, 108, 108, 111, 0, 19, 0, 50}
	Unmarshal(b, &result)
	if result != expected {
		t.Errorf("Result not expected: Expected: %v Actual %v", expected, result)
	}
}

func TestBasicArray(t *testing.T) {
	var result []uint16
	b := []byte{0, 0, 0, 1, 0, 2}
	expected := []uint16{uint16(0), uint16(1), uint16(2)}
	Unmarshal(b, &result)
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Result not expected: Expected: %v Actual %v", expected, result)
		}
	}
}
