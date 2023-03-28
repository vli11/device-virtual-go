// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020-2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func parseStrToInt(str string, bitSize int) (int64, error) {
	if i, err := strconv.ParseInt(str, 10, bitSize); err != nil {
		return i, err
	} else {
		return i, nil
	}
}

func parseIntMinimumMaximum(minimum, maximum *float64, dataType string) (int64, int64, error) {
	var err error
	if minimum == nil || maximum == nil {
		err = fmt.Errorf("minimum:%v maximum:%v are not in valid range, use default value", minimum, maximum)
		return 0, 0, err
	}

	var min, max int64
	min = int64(*minimum)
	var mmax float64
	mmax = *maximum
	max = int64(mmax)

	if max <= min {
		err = fmt.Errorf("minimum:%v maximum:%v are not in valid range, use default value", minimum, maximum)
		return 0, 0, err
	}

	// switch dataType {
	// case common.ValueTypeInt8, common.ValueTypeInt8Array:
	// 	min, err1 = parseStrToInt(minimum, 8)
	// 	max, err2 = parseStrToInt(maximum, 8)
	// 	if max <= min || err1 != nil || err2 != nil {
	// 		err = fmt.Errorf("minimum:%s maximum:%s not in valid range, use default value", minimum, maximum)
	// 	}
	// case common.ValueTypeInt16, common.ValueTypeInt16Array:
	// 	min, err1 = parseStrToInt(minimum, 16)
	// 	max, err2 = parseStrToInt(maximum, 16)
	// 	if max <= min || err1 != nil || err2 != nil {
	// 		err = fmt.Errorf("minimum:%s maximum:%s not in valid range, use default value", minimum, maximum)
	// 	}
	// case common.ValueTypeInt32, common.ValueTypeInt32Array:
	// 	min, err1 = parseStrToInt(minimum, 32)
	// 	max, err2 = parseStrToInt(maximum, 32)
	// 	if max <= min || err1 != nil || err2 != nil {
	// 		err = fmt.Errorf("minimum:%s maximum:%s not in valid range, use default value", minimum, maximum)
	// 	}
	// case common.ValueTypeInt64, common.ValueTypeInt64Array:
	// 	min, err1 = parseStrToInt(minimum, 64)
	// 	max, err2 = parseStrToInt(maximum, 64)
	// 	if max <= min || err1 != nil || err2 != nil {
	// 		err = fmt.Errorf("minimum:%s maximum:%s not in valid range, use default value", minimum, maximum)
	// 	}
	// }

	return min, max, err
}

func randomInt(min, max int64) int64 {
	if max > 0 && min < 0 {
		var negativePart int64
		var positivePart int64
		//min~0
		if min == int64(math.MinInt64) {
			negativePart = rand.Int63n(int64(math.MaxInt64)) + min - rand.Int63n(int64(1)) //nolint:gosec
		} else {
			negativePart = rand.Int63n(-min+int64(1)) + min //nolint:gosec
		}
		//0~max
		if max == int64(math.MaxInt64) {
			positivePart = rand.Int63n(max) + rand.Int63n(int64(1)) //nolint:gosec
		} else {
			positivePart = rand.Int63n(max + int64(1)) //nolint:gosec
		}
		return negativePart + positivePart
	} else {
		if max == int64(math.MaxInt64) && min == 0 {
			return rand.Int63n(max) + rand.Int63n(int64(1)) //nolint:gosec
		} else if min == int64(math.MinInt64) && max == 0 {
			return rand.Int63n(int64(math.MaxInt64)) + min - rand.Int63n(int64(1)) //nolint:gosec
		} else {
			return rand.Int63n(max-min+1) + min //nolint:gosec
		}
	}
}

func randomUint(min, max uint64) uint64 {
	//nolint // SA1019: rand.Seed has been deprecated
	rand.Seed(time.Now().UnixNano())
	if max-min < uint64(math.MaxInt64) {
		return uint64(rand.Int63n(int64(max-min+1))) + min //nolint:gosec
	}
	x := rand.Uint64() //nolint:gosec
	for x > max-min {
		x = rand.Uint64() //nolint:gosec
	}
	return x + min
}

func parseStrToUint(str string, bitSize int) (uint64, error) {
	if i, err := strconv.ParseUint(str, 10, bitSize); err != nil {
		return i, err
	} else {
		return i, nil
	}
}

func parseUintMinimumMaximum(minimum, maximum *float64, dataType string) (uint64, uint64, error) {
	var err error
	if minimum == nil || maximum == nil {
		err = fmt.Errorf("minimum:%v maximum:%v are not in valid range, use default value", minimum, maximum)
		return 0, 0, err
	}

	var min, max uint64
	min = uint64(*minimum)
	max = uint64(*maximum)

	return min, max, err
}

func randomFloat(min, max float64) float64 {
	//nolint // SA1019: rand.Seed has been deprecated
	rand.Seed(time.Now().UnixNano())
	if max > 0 && min < 0 {
		var negativePart float64
		var positivePart float64
		negativePart = rand.Float64() * min //nolint:gosec
		positivePart = rand.Float64() * max //nolint:gosec
		return negativePart + positivePart
	} else {
		return rand.Float64()*(max-min) + min //nolint:gosec
	}
}

func parseStrToFloat(str string, bitSize int) (float64, error) {
	if f, err := strconv.ParseFloat(str, bitSize); err != nil {
		return f, err
	} else {
		return f, nil
	}
}

func parseFloatMinimumMaximum(minimum, maximum *float64, dataType string) (float64, float64, error) {
	var err error
	if minimum == nil || maximum == nil {
		err = fmt.Errorf("minimum:%v maximum:%v are not in valid range, use default value", minimum, maximum)
		return 0, 0, err
	}
	return *minimum, *maximum, err
}
