// Package levenshtein is a Go implementation to calculate Levenshtein Distance with optional threshold.
//
// Original implementation taken from
// https://github.com/agnivade/levenshtein
// which is based on
// https://gist.github.com/andrei-m/982927#gistcomment-1931258
package levenshtein

import "unicode/utf8"

// minLengthThreshold is the length of the string beyond which
// an allocation will be made. Strings smaller than this will be
// zero alloc.
const minLengthThreshold = 32

// ComputeDistance computes the Levenshtein distance between the two
// strings passed as arguments.
//
//   - If you provide a threshold (maxDist) as the third parameter, the function
//     will attempt an early exit. That is, if the function detects during
//     computation that the distance must exceed `maxDist`, it will return
//     `maxDist + 1` immediately instead of computing the full distance.
//
//   - If you call it without a threshold or set a negative threshold, the
//     function will compute the distance in full, just like the original
//     implementation.
//
// Example 1: No threshold => always compute the full distance
// dist := ComputeDistance("kitten", "sitting")
// dist is 3, default behavior
//
// Example 2: Threshold = 2 => early exit if distance exceeds 2
// dist := ComputeDistance("kittenA", "sittingB", 2)
// dist is 3, the distance exceeds given threshold, and we do not compute the full distance
// and return threshold+1
//
// The function works on runes (Unicode code points) but does not normalize
// the input strings. See https://blog.golang.org/normalization
// and the golang.org/x/text/unicode/norm package if you need normalization.
func ComputeDistance(a, b string, threshold ...int) int {
	var maxDist int
	if len(threshold) > 0 {
		maxDist = threshold[0]
	} else {
		maxDist = -1
	}

	if len(a) == 0 {
		dist := utf8.RuneCountInString(b)
		if maxDist >= 0 && dist > maxDist {
			return maxDist + 1
		}
		return dist
	}

	if len(b) == 0 {
		dist := utf8.RuneCountInString(a)
		if maxDist >= 0 && dist > maxDist {
			return maxDist + 1
		}
		return dist
	}

	if a == b {
		return 0
	}

	s1 := []rune(a)
	s2 := []rune(b)

	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}

	for i := 0; i < len(s1); i++ {
		if s1[len(s1)-1-i] != s2[len(s2)-1-i] {
			s1 = s1[:len(s1)-i]
			s2 = s2[:len(s2)-i]
			break
		}
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			s1 = s1[i:]
			s2 = s2[i:]
			break
		}
	}

	lenS1 := len(s1)
	lenS2 := len(s2)

	if maxDist >= 0 && absInt(lenS1-lenS2) > maxDist {
		return maxDist + 1
	}

	var x []uint16
	if lenS1+1 > minLengthThreshold {
		x = make([]uint16, lenS1+1)
	} else {
		x = make([]uint16, minLengthThreshold)
		x = x[:lenS1+1]
	}

	for i := 1; i < len(x); i++ {
		x[i] = uint16(i)
	}

	for i := 1; i <= lenS2; i++ {
		prev := uint16(i)
		minInRow := prev

		for j := 1; j <= lenS1; j++ {
			current := x[j-1]

			if s2[i-1] != s1[j-1] {
				current = min(x[j-1]+1, prev+1, x[j]+1)
			}

			x[j-1] = prev
			prev = current

			if current < minInRow {
				minInRow = current
			}
		}

		x[lenS1] = prev

		if maxDist >= 0 && minInRow > uint16(maxDist) {
			return maxDist + 1
		}
	}

	return int(x[lenS1])
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
