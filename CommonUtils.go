package main

import (
    "strconv"
    "time"
)

//
//Extract the last 4 digit of a number
//
func LastFour(value int64) (trimValue string) {
    // Take substring of first word with runes.
    // ... This handles any kind of rune in the string.
    tempString := strconv.FormatInt(value, 10)
    length := len(tempString)
    //
    if length > 4 {
        runes := []rune(tempString)
        safeSubstring := string(runes[(length - 4):length])
        return safeSubstring
    }
    return tempString
}

//
// Returns the current time stamp
//
func Timestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}
