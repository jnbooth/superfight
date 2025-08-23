package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func WriteInvalid(w http.ResponseWriter, message []byte) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Header().Add("Content-Type", "text/plain")
	w.Write(message)
}

func FormString(r *http.Request, key string) (value string, err []byte) {
	if vs := r.Form[key]; len(vs) > 0 {
		return vs[0], nil
	}
	return "", fmt.Append(nil, "Missing parameter: ", key)
}

func FormStrings(r *http.Request, key string, count int) (values []string, err []byte) {
	vs := r.Form[key]
	if count == 0 {
		return vs, nil
	}
	if lenVs := len(vs); lenVs != count {
		return nil, fmt.Append(nil, key, ": expected ", count, " values, got ", lenVs)
	}
	return vs, nil
}

func formIntN(r *http.Request, key string, bitSize int) (value int64, err []byte) {
	val, err := FormString(r, key)
	if err != nil {
		return 0, err
	}
	value, errN := strconv.ParseInt(val, 10, bitSize)
	if errN != nil {
		return 0, fmt.Append(nil, key, errN.Error())
	}
	return value, nil
}

func formUintN(r *http.Request, key string, bitSize int) (value uint64, err []byte) {
	val, err := FormString(r, key)
	if err != nil {
		return 0, err
	}
	value, errN := strconv.ParseUint(val, 10, bitSize)
	if errN != nil {
		return 0, fmt.Append(nil, key, errN.Error())
	}
	return value, nil
}

func FormInt(r *http.Request, key string, min int, max int) (value int, err []byte) {
	val, err := formIntN(r, key, 0)
	if err != nil {
		return 0, err
	}
	value = int(val)
	if value < min || value > max {
		return 0, fmt.Append(nil, key, ": should be in range ", min, "-", max)
	}
	return value, nil
}

func FormByte(r *http.Request, key string, min byte, max byte) (value byte, err []byte) {
	val, err := formUintN(r, key, 8)
	if err != nil {
		return 0, err
	}
	value = byte(val)
	if value < min || value > max {
		return 0, fmt.Append(nil, key, ": should be in range ", min, "-", max)
	}
	return value, nil
}

func FormBytes(r *http.Request, key string, min byte, max byte, count int) (values []byte, err []byte) {
	vals, err := FormStrings(r, key, count)
	if err != nil {
		return nil, err
	}
	values = make([]byte, 0, len(vals))
	for _, sval := range vals {
		val, err := strconv.ParseUint(sval, 10, 8)
		if err != nil {
			return nil, fmt.Append(nil, key, err.Error())
		}
		value := byte(val)
		if value < min || value > max {
			return nil, fmt.Append(nil, key, ": should be in range ", min, "-", max)
		}
		values = append(values, byte(value))
	}
	return values, nil
}
