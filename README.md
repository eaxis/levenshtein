# Early-Exit Levenshtein ![Build Status](https://github.com/eaxis/levenshtein/actions/workflows/ci.yml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/eaxis/levenshtein)](https://goreportcard.com/report/github.com/eaxis/levenshtein) [![PkgGoDev](https://pkg.go.dev/badge/github.com/eaxis/levenshtein)](https://pkg.go.dev/github.com/eaxis/levenshtein)

A [Go](http://golang.org) package to calculate the [Levenshtein Distance](http://en.wikipedia.org/wiki/Levenshtein_distance) with optional early-exit optimization for comparing large texts for similarity.

This library is based on the [agnivade's implementation](https://github.com/agnivade/levenshtein), and without providing a threshold, it works exactly the same.
But with a threshold, it stops calculating the distance when it exceeds the threshold, and returns the threshold + 1. This optimization saves CPU time when comparing large strings.

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

### Comparisons with other libraries (short strings with threshold)

```
name                                                    time/op
BenchmarkCompetitorsWithThreshold/Nordic/eaxis-12       ~206ns
BenchmarkCompetitorsWithThreshold/Nordic/agniva-12      ~267ns
BenchmarkCompetitorsWithThreshold/Nordic/arbovm-12      ~708ns
BenchmarkCompetitorsWithThreshold/Nordic/dgryski-12     ~700ns
```

### Comparisons with other libraries (long strings with threshold)

```
name                                                    time/op
BenchmarkCompetitorsWithThreshold/Russian/eaxis-12      ~1860ns
BenchmarkCompetitorsWithThreshold/Russian/agniva-12     ~86593ns
BenchmarkCompetitorsWithThreshold/Russian/arbovm-12     ~64309ns
BenchmarkCompetitorsWithThreshold/Russian/dgryski-12    ~63692ns
```

### Comparisons with other libraries (short strings)

```
name                                                    time/op
BenchmarkCompetitors/Nordic/eaxis-12                    ~315ns
BenchmarkCompetitors/Nordic/agniva-12                   ~271ns
BenchmarkCompetitors/Nordic/arbovm-12                   ~714ns
BenchmarkCompetitors/Nordic/dgryski-12                  ~691ns
```

### Comparisons with other libraries (long strings)

```
name                                                    time/op
BenchmarkCompetitors/Russian/eaxis-12                   ~91058ns
BenchmarkCompetitors/Russian/agniva-12                  ~88403ns
BenchmarkCompetitors/Russian/arbovm-12                  ~68272ns
BenchmarkCompetitors/Russian/dgryski-12                 ~65855ns
```