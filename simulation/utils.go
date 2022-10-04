package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"math/rand"
)

// HexSlice2ByteSlice returns a slice of hex string hashes as byte slice.
func HexSlice2ByteSlice(str []string) [][]byte {
	var slice [][]byte
	for _, s := range str {
		slice = append(slice, Hex2Bytes(s))
	}
	return slice
}

// Hex2Bytes returns a hex string hash as byte slice.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// Bytes2Hex returns a string from a byte slice
func Bytes2Hex(byt []byte) string {
	h := hex.EncodeToString(byt)
	return h
}

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// random strings of fixed length from
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
// grabbed this function from lab1 part4 merkle tree benchmark
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// sums a list of ints to compute the average
func sum(slice []int) int {
	result := 0
	for _, i := range slice {
		result += i
	}
	return result
}

// returns the average of a slice
func average(slice []int) int {
	return sum(slice) / len(slice)
}
