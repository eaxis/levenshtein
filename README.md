# Early-Exit Levenshtein ![Build Status](https://github.com/eaxis/levenshtein/actions/workflows/ci.yml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/eaxis/levenshtein)](https://goreportcard.com/report/github.com/eaxis/levenshtein) [![PkgGoDev](https://pkg.go.dev/badge/github.com/eaxis/levenshtein)](https://pkg.go.dev/github.com/eaxis/levenshtein)

A [Go](http://golang.org) package to calculate the [Levenshtein Distance](http://en.wikipedia.org/wiki/Levenshtein_distance) with optional early-exit optimization for comparing large texts for similarity.

This library is based on [agnivade's implementation](https://github.com/agnivade/levenshtein), and works the same way if no threshold is provided.
When a threshold is set, the library stops calculating the distance as soon as the distance exceeds the threshold and returns `threshold + 1` instead of calculating the remaining distance, saving significant CPU time when comparing long strings.

The library also nests the following features/limitations:
- The library is fully capable of working with non-ascii strings. But the strings are not normalized.
- As a performance optimization, the library can handle strings only up to 65536 characters (runes).

## Motivation

I created this library because I needed to compare thousands of posts for duplicates in my side project.
The process was disappointingly slow and only got worse as more posts were added.
By using an early-exit approach, I’ve achieved over 100x speedup in my certain case.
You can find detailed benchmarks below.

## Install
```
go get github.com/eaxis/levenshtein
```

## A simple example

```go
package main

import (
	"fmt"
	"github.com/eaxis/levenshtein"
)

func main() {
	s1 := "kitten"
	s2 := "sitting"
	distance := levenshtein.ComputeDistance(s1, s2)
	fmt.Printf("The distance is %d.\n", distance) // The distance is 3.
}

```

## An example with early-exit optimization

```go
package main

import (
	"fmt"
	"github.com/eaxis/levenshtein"
)

func main() {
	similarityThreshold := 10

	// The Levenstein distance between these strings is 47.
	// Since the similarityThreshold is 10, the function will stop calculating the distance at 10 and return 11.
	// Which means the distance is greater than the similarityThreshold.
	s1 := "these strings are completely different and have nothing in common"
	s2 := "calculating the full distance is just a waste of time"

	distance := levenshtein.ComputeDistance(s1, s2, similarityThreshold)
	fmt.Printf("The distance is at least %d.\n", distance) // The distance is at least 11.

	if distance <= similarityThreshold {
		fmt.Println("The strings are similar.")
	} else {
		fmt.Println("The strings are not similar.") // this will be printed
	}
}

```

## Benchmarks

### Comparisons with other libraries (short strings with threshold = 2)

```
BenchmarkCompetitorsWithThreshold/ASCII_short/eaxis-12 	        ~61 ns/op
BenchmarkCompetitorsWithThreshold/ASCII_short/agniva-12         ~90 ns/op
BenchmarkCompetitorsWithThreshold/ASCII_short/arbovm-12        ~221 ns/op
BenchmarkCompetitorsWithThreshold/ASCII_short/dgryski-12       ~220 ns/op

```

### Comparisons with other libraries (long strings with threshold = 10)

```
BenchmarkCompetitorsWithThreshold/ASCII_long/eaxis-12         ~1995 ns/op
BenchmarkCompetitorsWithThreshold/ASCII_long/agniva-12      ~634766 ns/op
BenchmarkCompetitorsWithThreshold/ASCII_long/arbovm-12      ~917835 ns/op
BenchmarkCompetitorsWithThreshold/ASCII_long/dgryski-12     ~921690 ns/op
```

### Comparisons with other libraries (short strings)

```
BenchmarkCompetitors/ASCII_short/eaxis-12              	       ~111 ns/op
BenchmarkCompetitors/ASCII_short/agniva-12             	        ~91 ns/op
BenchmarkCompetitors/ASCII_short/arbovm-12             	       ~219 ns/op
BenchmarkCompetitors/ASCII_short/dgryski-12            	       ~223 ns/op
```

### Comparisons with other libraries (long strings)

```
BenchmarkCompetitors/ASCII_long/eaxis-12               	    ~726370 ns/op
BenchmarkCompetitors/ASCII_long/agniva-12              	    ~633286 ns/op
BenchmarkCompetitors/ASCII_long/arbovm-12              	    ~900986 ns/op
BenchmarkCompetitors/ASCII_long/dgryski-12             	    ~912527 ns/op
```