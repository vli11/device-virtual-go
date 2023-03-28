// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/edgexfoundry/device-sdk-go/v3/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/common"
)

const (
	deviceName              = "Random-Value-Device"
	enableRandomizationTrue = true
	rounds                  = 10

	nameBool         = common.ValueTypeBool
	nameBoolArray    = common.ValueTypeBoolArray
	nameInt8         = common.ValueTypeInt8
	nameInt8Array    = common.ValueTypeInt8Array
	nameInt16        = common.ValueTypeInt16
	nameInt16Array   = common.ValueTypeInt16Array
	nameInt32        = common.ValueTypeInt32
	nameInt32Array   = common.ValueTypeInt32Array
	nameInt64        = common.ValueTypeInt64
	nameInt64Array   = common.ValueTypeInt64Array
	nameUint8        = common.ValueTypeUint8
	nameUint8Array   = common.ValueTypeUint8Array
	nameUint16       = common.ValueTypeUint16
	nameUint16Array  = common.ValueTypeUint16Array
	nameUint32       = common.ValueTypeUint32
	nameUint32Array  = common.ValueTypeUint32Array
	nameUint64       = common.ValueTypeUint64
	nameUint64Array  = common.ValueTypeUint64Array
	nameFloat32      = common.ValueTypeFloat32
	nameFloat32Array = common.ValueTypeFloat32Array
	nameFloat64      = common.ValueTypeFloat64
	nameFloat64Array = common.ValueTypeFloat64Array
	nameBinary       = common.ValueTypeBinary
)

type resourceDef struct {
	devName    string
	cmdName    string
	resName    string
	randEnable bool
	dataType   string
	initValue  string
}

func prepareDB() *db {
	db := getDb()

	if err := db.init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ds := []resourceDef{
		{deviceName, nameBool, nameBool, enableRandomizationTrue, nameBool, "true"},
		{deviceName, nameBoolArray, nameBoolArray, enableRandomizationTrue, nameBoolArray, "[true]"},
		{deviceName, nameInt8, nameInt8, enableRandomizationTrue, nameInt8, "0"},
		{deviceName, nameInt8Array, nameInt8Array, enableRandomizationTrue, nameInt8Array, "[0]"},
		{deviceName, nameInt16, nameInt16, enableRandomizationTrue, nameInt16, "0"},
		{deviceName, nameInt16Array, nameInt16Array, enableRandomizationTrue, nameInt16Array, "[0]"},
		{deviceName, nameInt32, nameInt32, enableRandomizationTrue, nameInt32, "0"},
		{deviceName, nameInt32Array, nameInt32Array, enableRandomizationTrue, nameInt32Array, "[0]"},
		{deviceName, nameInt64, nameInt64, enableRandomizationTrue, nameInt64, "0"},
		{deviceName, nameInt64Array, nameInt64Array, enableRandomizationTrue, nameInt64Array, "[0]"},
		{deviceName, nameUint8, nameUint8, enableRandomizationTrue, nameUint8, "0"},
		{deviceName, nameUint8Array, nameUint8Array, enableRandomizationTrue, nameUint8Array, "[0]"},
		{deviceName, nameUint16, nameUint16, enableRandomizationTrue, nameUint16, "0"},
		{deviceName, nameUint16Array, nameUint16Array, enableRandomizationTrue, nameUint16Array, "[0]"},
		{deviceName, nameUint32, nameUint32, enableRandomizationTrue, nameUint32, "0"},
		{deviceName, nameUint32Array, nameUint32Array, enableRandomizationTrue, nameUint32Array, "[0]"},
		{deviceName, nameUint64, nameUint64, enableRandomizationTrue, nameUint64, "0"},
		{deviceName, nameUint64Array, nameUint64Array, enableRandomizationTrue, nameUint64Array, "[0]"},
		{deviceName, nameFloat32, nameFloat32, enableRandomizationTrue, nameFloat32, "0"},
		{deviceName, nameFloat32Array, nameFloat32Array, enableRandomizationTrue, nameFloat32Array, "[0]"},
		{deviceName, nameFloat64, nameFloat64, enableRandomizationTrue, nameFloat64, "0"},
		{deviceName, nameFloat64Array, nameFloat64Array, enableRandomizationTrue, nameFloat64Array, "[0]"},
	}
	for _, d := range ds {
		if err := db.addResource(d.devName, d.cmdName, d.resName, d.randEnable, d.dataType, d.initValue); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return db
}

func TestValueBool(t *testing.T) {
	db := prepareDB()
	defer func() {
		if err := db.closeDb(); err != nil {
			t.Fatal(err)
		}
	}()

	vd := newVirtualDevice()
	v1, err := vd.read(deviceName, nameBool, nameBool, nil, nil, db)
	if err != nil {
		t.Fatal(err)
	}

	//the return string must be convertible to boolean
	b1, err := v1.BoolValue()
	if err != nil {
		t.Fatal(err)
	}

	//EnableRandomization = true
	for x := 1; x <= rounds; x++ {
		v2, _ := vd.read(deviceName, nameBool, nameBool, nil, nil, db)
		b2, _ := v2.BoolValue()
		if b1 != b2 {
			break
		}
		if x == rounds {
			t.Fatalf("EnableRandomization is true, but got same read in %d rounds", rounds)
		}
	}

	//EnableRandomization = false
	if err := db.updateResourceRandomization(false, deviceName, nameBool); err != nil {
		t.Fatal(err)
	}

	v1, _ = vd.read(deviceName, nameBool, nameBool, nil, nil, db)
	b1, _ = v1.BoolValue()
	for x := 0; x <= rounds; x++ {
		v2, _ := vd.read(deviceName, nameBool, nameBool, nil, nil, db)
		b2, _ := v2.BoolValue()
		if b1 != b2 {
			t.Fatalf("EnableRandomization is false, but got different read")
		}
	}
}

func TestValueBoolArray(t *testing.T) {
	db := prepareDB()
	defer func() {
		if err := db.closeDb(); err != nil {
			t.Fatal(err)
		}
	}()

	vd := newVirtualDevice()
	v1, err := vd.read(deviceName, nameBoolArray, nameBoolArray, nil, nil, db)
	if err != nil {
		t.Fatal(err)
	}

	// the returned string must be convertible to boolean array
	b1, err := v1.BoolArrayValue()
	if err != nil {
		t.Fatal(err)
	}

	// EnableRandomization = true
	for x := 1; x <= rounds; x++ {
		v2, _ := vd.read(deviceName, nameBoolArray, nameBoolArray, nil, nil, db)
		b2, _ := v2.BoolArrayValue()
		if !reflect.DeepEqual(b1, b2) {
			break
		}
		if x == rounds {
			t.Fatalf("EnableRandomization is true, but got the same reading in %d rounds", rounds)
		}
	}

	// EnableRandomization = false
	if err := db.updateResourceRandomization(false, deviceName, nameBoolArray); err != nil {
		t.Fatal(err)
	}

	v1, _ = vd.read(deviceName, nameBoolArray, nameBoolArray, nil, nil, db)
	b1, _ = v1.BoolArrayValue()
	for x := 0; x <= rounds; x++ {
		v2, _ := vd.read(deviceName, nameBoolArray, nameBoolArray, nil, nil, db)
		b2, _ := v2.BoolArrayValue()
		if !reflect.DeepEqual(b1, b2) {
			t.Fatalf("EnableRandomization is false, but got different reading")
		}
	}
}

func TestValueIntx(t *testing.T) {
	var min, max float64
	min = -128
	max = 127
	ValueIntx(t, nameInt8, nameInt8, &min, &max)
	ValueIntx(t, nameInt8, nameInt8, nil, nil)
	min = -32768
	max = 32767
	ValueIntx(t, nameInt16, nameInt16, &min, &max)
	ValueIntx(t, nameInt16, nameInt16, nil, nil)
	min = -2147483648
	max = 2147483647
	ValueIntx(t, nameInt32, nameInt32, &min, &max)
	ValueIntx(t, nameInt32, nameInt32, nil, nil)
	min = -9223372036854775808
	max = math.MaxInt64 //9223372036854775807

	ValueIntx(t, nameInt64, nameInt64, &min, &max)
	ValueIntx(t, nameInt64, nameInt64, nil, nil)
}

func TestValueIntxArray(t *testing.T) {
	var min, max float64
	min = -128
	max = 127
	ValueIntxArray(t, nameInt8Array, nameInt8Array, &min, &max)
	ValueIntxArray(t, nameInt8Array, nameInt8Array, nil, nil)
	min = -32768
	max = 32767
	ValueIntxArray(t, nameInt16Array, nameInt16Array, &min, &max)
	ValueIntxArray(t, nameInt16Array, nameInt16Array, nil, nil)
	min = -2147483648
	max = 2147483647
	ValueIntxArray(t, nameInt32Array, nameInt32Array, &min, &max)
	ValueIntxArray(t, nameInt32Array, nameInt32Array, nil, nil)
	min = -9223372036854775808
	max = 9223372036854775807
	ValueIntxArray(t, nameInt64Array, nameInt64Array, &min, &max)
	ValueIntxArray(t, nameInt64Array, nameInt64Array, nil, nil)
}

func TestValueUintx(t *testing.T) {
	var min, max float64
	min = 0
	max = 255
	ValueUintx(t, nameUint8, nameUint8, &min, &max)
	ValueUintx(t, nameUint8, nameUint8, nil, nil)
	min = 0
	max = 65535
	ValueUintx(t, nameUint16, nameUint16, &min, &max)
	ValueUintx(t, nameUint16, nameUint16, nil, nil)
	min = 0
	max = 4294967295
	ValueUintx(t, nameUint32, nameUint32, &min, &max)
	ValueUintx(t, nameUint32, nameUint32, nil, nil)
	min = 0
	max = 18446744073709551615
	ValueUintx(t, nameUint64, nameUint64, &min, &max)
	ValueUintx(t, nameUint64, nameUint64, nil, nil)
}

func TestValueUintxArray(t *testing.T) {
	var min, max float64
	min = 0
	max = 255
	ValueUintxArray(t, nameUint8Array, nameUint8Array, &min, &max)
	ValueUintxArray(t, nameUint8Array, nameUint8Array, nil, nil)
	min = 0
	max = 65535
	ValueUintxArray(t, nameUint16Array, nameUint16Array, &min, &max)
	ValueUintxArray(t, nameUint16Array, nameUint16Array, nil, nil)
	min = 0
	max = 4294967295
	ValueUintxArray(t, nameUint32Array, nameUint32Array, &min, &max)
	ValueUintxArray(t, nameUint32Array, nameUint32Array, nil, nil)
	min = 0
	max = 18446744073709551615
	ValueUintxArray(t, nameUint64Array, nameUint64Array, &min, &max)
	ValueUintxArray(t, nameUint64Array, nameUint64Array, nil, nil)
}

func TestValueFloatx(t *testing.T) {
	var min, max float64
	min = -3.40282346638528859811704183484516925440e+38
	max = 3.40282346638528859811704183484516925440e+38
	ValueFloatx(t, nameFloat32, nameFloat32, &min, &max)
	ValueFloatx(t, nameFloat32, nameFloat32, nil, nil)
	min = -1.797693134862315708145274237317043567981e+308
	max = 1.797693134862315708145274237317043567981e+308
	ValueFloatx(t, nameFloat64, nameFloat64, &min, &max)
	ValueFloatx(t, nameFloat64, nameFloat64, nil, nil)
}

func TestValueFloatxArray(t *testing.T) {
	var min, max float64
	min = -3.40282346638528859811704183484516925440e+38
	max = 3.40282346638528859811704183484516925440e+38
	ValueFloatxArray(t, nameFloat32Array, nameFloat32Array, &min, &max)
	ValueFloatxArray(t, nameFloat32Array, nameFloat32Array, nil, nil)
	min = -1.797693134862315708145274237317043567981e+308
	max = 1.797693134862315708145274237317043567981e+308
	ValueFloatxArray(t, nameFloat64Array, nameFloat64Array, &min, &max)
	ValueFloatxArray(t, nameFloat64Array, nameFloat64Array, nil, nil)
}

func TestValueBinary(t *testing.T) {
	vd := newVirtualDevice()
	v1, err := vd.read(deviceName, nameBinary, nameBinary, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	//the return string must be convertible to binary
	_, err = v1.BinaryValue()
	if err != nil {
		t.Fatal(err)
	}
}

func ValueIntx(t *testing.T, dr, typeName string, min, max *float64) {
	db := prepareDB()
	defer func() {
		if err := db.closeDb(); err != nil {
			t.Fatal(err)
		}
	}()

	//EnableRandomization = true
	if err := db.updateResourceRandomization(true, deviceName, dr); err != nil {
		t.Fatal(err)
	}

	vd := newVirtualDevice()

	var i1 int64
	for x := 1; x <= rounds; x++ {
		vn, err := vd.read(deviceName, dr, typeName, min, max, db)
		if err != nil {
			t.Fatal(err)
		}
		in := getIntValue(vn)

		if x == 1 {
			i1 = in
		}
		if i1 != in {
			break
		}
		if x == rounds {
			t.Fatalf("EnableRandomization is true, but got same read in %d rounds", rounds)
		}
	}

	for x := 1; x <= rounds; x++ {
		v, err := vd.read(deviceName, dr, typeName, min, max, db)

		if err != nil {
			t.Fatal(err)
		}
		i := getIntValue(v)
		if err != nil {
			t.Fatal(err)
		}
		if min != nil && max != nil {
			var floati float64
			floati = float64(i)
			if floati < *min || floati > *max {
				t.Fatalf("the random reading: %d is out of range: %v ~ %v", i, *min, *max)
			}
		}
	}

	//EnableRandomization = false
	if err := db.updateResourceRandomization(false, deviceName, dr); err != nil {
		t.Fatal(err)
	}

	v1, _ := vd.read(deviceName, dr, typeName, min, max, db)
	i1 = getIntValue(v1)
	for x := 1; x <= rounds; x++ {
		v2, _ := vd.read(deviceName, dr, typeName, min, max, db)
		i2 := getIntValue(v2)
		if i1 != i2 {
			t.Fatalf("EnableRandomization is false, but got different read")
		}
	}
}

func ValueIntxArray(t *testing.T, dr, typeName string, min, max *float64) {
	db := prepareDB()
	defer func() {
		if err := db.closeDb(); err != nil {
			t.Fatal(err)
		}
	}()

	// EnableRandomization = true
	if err := db.updateResourceRandomization(true, deviceName, dr); err != nil {
		t.Fatal(err)
	}

	vd := newVirtualDevice()

	var intArr1 []int64
	for x := 1; x <= rounds; x++ {
		vn, err := vd.read(deviceName, dr, typeName, min, max, db)
		if err != nil {
			t.Fatal(err)
		}
		intArrN := getIntArrayValue(vn)

		if x == 1 {
			intArr1 = intArrN
		}
		if !reflect.DeepEqual(intArr1, intArrN) {
			break
		}
		if x == rounds {
			t.Fatalf("EnableRandomization is true, but got the same reading in %d rounds", rounds)
		}
	}

	for x := 1; x <= rounds; x++ {
		v, err := vd.read(deviceName, dr, typeName, min, max, db)

		if err != nil {
			t.Fatal(err)
		}
		intArr := getIntArrayValue(v)
		if err != nil {
			t.Fatal(err)
		}
		if min != nil && max != nil {
			for _, i := range intArr {
				var floati float64
				floati = float64(i)
				if floati < *min || floati > *max {
					t.Fatalf("the random reading: %d is out of range: %v ~ %v", i, *min, *max)
				}
			}
		}
	}

	// EnableRandomization = false
	if err := db.updateResourceRandomization(false, deviceName, dr); err != nil {
		t.Fatal(err)
	}

	v1, _ := vd.read(deviceName, dr, typeName, min, max, db)
	intArr1 = getIntArrayValue(v1)
	for x := 1; x <= rounds; x++ {
		v2, _ := vd.read(deviceName, dr, typeName, min, max, db)
		intArr2 := getIntArrayValue(v2)
		if !reflect.DeepEqual(intArr1, intArr2) {
			t.Fatalf("EnableRandomization is false, but got different reading")
		}
	}
}

func ValueUintx(t *testing.T, dr, typeName string, min, max *float64) {
	db := prepareDB()
	defer func() {
		if err := db.closeDb(); err != nil {
			t.Fatal(err)
		}
	}()

	// EnableRandomization = true
	if err := db.updateResourceRandomization(true, deviceName, dr); err != nil {
		t.Fatal(err)
	}

	vd := newVirtualDevice()

	var i1 uint64
	for x := 1; x <= rounds; x++ {
		vn, _ := vd.read(deviceName, dr, typeName, min, max, db)
		in := getUintValue(vn)

		if x == 1 {
			i1 = in
		}
		if i1 != in {
			break
		}
		if x == rounds {
			t.Fatalf("EnableRandomization is true, but got same reading in %d rounds", rounds)
		}
	}

	for x := 1; x <= rounds; x++ {
		v, err := vd.read(deviceName, dr, typeName, min, max, db)
		if err != nil {
			t.Fatal(err)
		}
		i := getUintValue(v)
		if err != nil {
			t.Fatal(err)
		}
		if min != nil && max != nil {
			var floati float64
			floati = float64(i)
			if floati < *min || floati > *max {
				t.Fatalf("the random reading: %d is out of range: %v ~ %v", i, *min, *max)
			}
		}
	}

	// EnableRandomization = false
	if err := db.updateResourceRandomization(false, deviceName, dr); err != nil {
		t.Fatal(err)
	}

	v1, _ := vd.read(deviceName, dr, typeName, min, max, db)
	i1 = getUintValue(v1)
	for x := 1; x <= rounds; x++ {
		v2, _ := vd.read(deviceName, dr, typeName, min, max, db)
		i2 := getUintValue(v2)
		if i1 != i2 {
			t.Fatalf("EnableRandomization is false, but got different reading")
		}
	}
}

func ValueUintxArray(t *testing.T, dr, typeName string, min, max *float64) {
	db := prepareDB()
	defer func() {
		if err := db.closeDb(); err != nil {
			t.Fatal(err)
		}
	}()

	// EnableRandomization = true
	if err := db.updateResourceRandomization(true, deviceName, dr); err != nil {
		t.Fatal(err)
	}

	vd := newVirtualDevice()

	var uintArr1 []uint64
	for x := 1; x <= rounds; x++ {
		vn, _ := vd.read(deviceName, dr, typeName, min, max, db)
		uintArrN := getUintArrayValue(vn)

		if x == 1 {
			uintArr1 = uintArrN
		}
		if !reflect.DeepEqual(uintArr1, uintArrN) {
			break
		}
		if x == rounds {
			t.Fatalf("EnableRandomization is true, but got the same reading in %d rounds", rounds)
		}
	}

	for x := 1; x <= rounds; x++ {
		v, err := vd.read(deviceName, dr, typeName, min, max, db)
		if err != nil {
			t.Fatal(err)
		}
		uintArr := getUintArrayValue(v)
		if err != nil {
			t.Fatal(err)
		}
		if min != nil && max != nil {
			for _, i := range uintArr {
				var floati float64
				floati = float64(i)
				if floati < *min || floati > *max {
					t.Fatalf("the random reading: %d is out of range: %v ~ %v", i, *min, *max)
				}
			}
		}
	}

	// EnableRandomization = false
	if err := db.updateResourceRandomization(false, deviceName, dr); err != nil {
		t.Fatal(err)
	}

	v1, _ := vd.read(deviceName, dr, typeName, min, max, db)
	uintArr1 = getUintArrayValue(v1)
	for x := 1; x <= rounds; x++ {
		v2, _ := vd.read(deviceName, dr, typeName, min, max, db)
		uintArr2 := getUintArrayValue(v2)
		if !reflect.DeepEqual(uintArr1, uintArr2) {
			t.Fatalf("EnableRandomization is false, but got different reading")
		}
	}
}

func ValueFloatx(t *testing.T, dr, typeName string, min, max *float64) {
	db := prepareDB()
	defer func() {
		if err := db.closeDb(); err != nil {
			t.Fatal(err)
		}
	}()

	// EnableRandomization = true
	if err := db.updateResourceRandomization(true, deviceName, dr); err != nil {
		t.Fatal(err)
	}

	vd := newVirtualDevice()

	var f1 float64
	for x := 1; x <= rounds; x++ {
		vn, _ := vd.read(deviceName, dr, typeName, min, max, db)
		fn := getFloatValue(vn)
		if x == 1 {
			f1 = fn
		}
		if f1 != fn {
			break
		}
		if x == rounds {
			t.Fatalf("EnableRandomization is true, but got the same reading in %d rounds", rounds)
		}
	}

	for x := 1; x <= rounds; x++ {
		v, err := vd.read(deviceName, dr, typeName, min, max, db)
		if err != nil {
			t.Fatal(err)
		}
		f := getFloatValue(v)
		if err != nil {
			t.Fatal(err)
		}
		if min != nil && max != nil {
			if f < *min || f > *max {
				t.Fatalf("the random reading: %f is out of range: %v ~ %v", f, *min, *max)
			}
		}
	}

	// EnableRandomization = false
	if err := db.updateResourceRandomization(false, deviceName, dr); err != nil {
		t.Fatal(err)
	}

	v1, _ := vd.read(deviceName, dr, typeName, min, max, db)
	f1 = getFloatValue(v1)
	for x := 1; x <= rounds; x++ {
		v2, _ := vd.read(deviceName, dr, typeName, min, max, db)
		f2 := getFloatValue(v2)
		if f1 != f2 {
			t.Fatalf("EnableRandomization is false, but got different reading")
		}
	}
}

func ValueFloatxArray(t *testing.T, dr, typeName string, min, max *float64) {
	db := prepareDB()
	defer func() {
		if err := db.closeDb(); err != nil {
			t.Fatal(err)
		}
	}()

	// EnableRandomization = true
	if err := db.updateResourceRandomization(true, deviceName, dr); err != nil {
		t.Fatal(err)
	}

	vd := newVirtualDevice()

	var floatArr1 []float64
	for x := 1; x <= rounds; x++ {
		vn, _ := vd.read(deviceName, dr, typeName, min, max, db)
		floatArrN := getFloatArrayValue(vn)

		if x == 1 {
			floatArr1 = floatArrN
		}
		if !reflect.DeepEqual(floatArr1, floatArrN) {
			break
		}
		if x == rounds {
			t.Fatalf("EnableRandomization is true, but got the same reading in %d rounds", rounds)
		}
	}

	for x := 1; x <= rounds; x++ {
		v, err := vd.read(deviceName, dr, typeName, min, max, db)
		if err != nil {
			t.Fatal(err)
		}
		floatArr := getFloatArrayValue(v)
		if err != nil {
			t.Fatal(err)
		}
		if min != nil && max != nil {
			for _, f := range floatArr {
				if f < *min || f > *max {
					t.Fatalf("the random reading: %f is out of range: %v ~ %v", f, *min, *max)
				}
			}
		}
	}

	// EnableRandomization = false
	if err := db.updateResourceRandomization(false, deviceName, dr); err != nil {
		t.Fatal(err)
	}

	v1, _ := vd.read(deviceName, dr, typeName, min, max, db)
	floatArr1 = getFloatArrayValue(v1)
	for x := 1; x <= rounds; x++ {
		v2, _ := vd.read(deviceName, dr, typeName, min, max, db)
		floatArr2 := getFloatArrayValue(v2)
		if !reflect.DeepEqual(floatArr1, floatArr2) {
			t.Fatalf("EnableRandomization is false, but got different reading")
		}
	}
}

func getIntValue(cv *models.CommandValue) int64 {
	switch cv.Type {
	case common.ValueTypeInt8:
		v, _ := cv.Int8Value()
		return int64(v)
	case common.ValueTypeInt16:
		v, _ := cv.Int16Value()
		return int64(v)
	case common.ValueTypeInt32:
		v, _ := cv.Int32Value()
		return int64(v)
	case common.ValueTypeInt64:
		v, _ := cv.Int64Value()
		return v
	default:
		return 0
	}
}

func getIntArrayValue(cv *models.CommandValue) []int64 {
	var value []int64
	switch cv.Type {
	case common.ValueTypeInt8Array:
		int8Arr, _ := cv.Int8ArrayValue()
		for _, i := range int8Arr {
			value = append(value, int64(i))
		}
	case common.ValueTypeInt16Array:
		int16Arr, _ := cv.Int16ArrayValue()
		for _, i := range int16Arr {
			value = append(value, int64(i))
		}
	case common.ValueTypeInt32Array:
		int32Arr, _ := cv.Int32ArrayValue()
		for _, i := range int32Arr {
			value = append(value, int64(i))
		}
	case common.ValueTypeInt64Array:
		value, _ = cv.Int64ArrayValue()
	default:
		value = []int64{0}
	}
	return value
}

func getUintValue(cv *models.CommandValue) uint64 {
	switch cv.Type {
	case common.ValueTypeUint8:
		v, _ := cv.Uint8Value()
		return uint64(v)
	case common.ValueTypeUint16:
		v, _ := cv.Uint16Value()
		return uint64(v)
	case common.ValueTypeUint32:
		v, _ := cv.Uint32Value()
		return uint64(v)
	case common.ValueTypeUint64:
		v, _ := cv.Uint64Value()
		return v
	default:
		return 0
	}
}

func getUintArrayValue(cv *models.CommandValue) []uint64 {
	var value []uint64
	switch cv.Type {
	case common.ValueTypeUint8Array:
		uint8Arr, _ := cv.Uint8ArrayValue()
		for _, i := range uint8Arr {
			value = append(value, uint64(i))
		}
	case common.ValueTypeUint16Array:
		uint16Arr, _ := cv.Uint16ArrayValue()
		for _, i := range uint16Arr {
			value = append(value, uint64(i))
		}
	case common.ValueTypeUint32Array:
		uint32Arr, _ := cv.Uint32ArrayValue()
		for _, i := range uint32Arr {
			value = append(value, uint64(i))
		}
	case common.ValueTypeUint64Array:
		value, _ = cv.Uint64ArrayValue()
	default:
		value = []uint64{0}
	}
	return value
}

func getFloatValue(cv *models.CommandValue) float64 {
	switch cv.Type {
	case common.ValueTypeFloat32:
		v, _ := cv.Float32Value()
		return float64(v)
	case common.ValueTypeFloat64:
		v, _ := cv.Float64Value()
		return v
	default:
		return 0
	}
}

func getFloatArrayValue(cv *models.CommandValue) []float64 {
	var value []float64
	switch cv.Type {
	case common.ValueTypeFloat32Array:
		float32Arr, _ := cv.Float32ArrayValue()
		for _, f := range float32Arr {
			value = append(value, float64(f))
		}
	case common.ValueTypeFloat64Array:
		value, _ = cv.Float64ArrayValue()
	default:
		value = []float64{0}
	}
	return value
}
