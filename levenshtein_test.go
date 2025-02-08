package levenshtein

import (
	agnivade "github.com/agnivade/levenshtein"
	arbovm "github.com/arbovm/levenshtein"
	dgryski "github.com/dgryski/trifles/leven"

	"testing"
)

type testCase struct {
	a         string
	b         string
	threshold int
	want      int
	name      string
}

func computeDistanceTestCases() []testCase {
	return []testCase{
		{
			a:         "",
			b:         "",
			threshold: -1,
			want:      0,
			name:      "Empty strings, no threshold",
		},
		{
			a:         "hello",
			b:         "hello",
			threshold: -1,
			want:      0,
			name:      "Same strings, no threshold",
		},
		// the classic test case: kitten -> sitting
		{
			a:         "kitten",
			b:         "sitting",
			threshold: -1,
			want:      3,
			name:      "Simple example, no threshold",
		},
		{
			a:         "",
			b:         "abc",
			threshold: -1,
			want:      3,
			name:      "Empty vs non-empty, no threshold",
		},
		// the distance is 3, and it's below threshold of 5
		{
			a:         "kitten",
			b:         "sitting",
			threshold: 5,
			want:      3,
			name:      "Simple with threshold, distance < threshold",
		},
		// the distance is 2 (replace 'k'->'s', + remove (replace) 'X')
		{
			a:         "kitten",
			b:         "sittenX",
			threshold: 2,
			want:      2,
			name:      "Simple with threshold, distance == threshold",
		},
		// distance = 3, threshold 2 => returns threshold+1 = 3
		{
			a:         "kitten",
			b:         "sitting",
			threshold: 2,
			want:      3,
			name:      "Simple with threshold, distance > threshold",
		},
		// the distance is > 2 => returns threshold+1
		{
			a:         "abcdef",
			b:         "uvwxyz",
			threshold: 2,
			want:      3,
			name:      "Bigger difference, early exit expected",
		},
		// XYZ vs WXY — the distance is 2
		{
			a:         "abcXYZ",
			b:         "abcWXY",
			threshold: -1,
			want:      2,
			name:      "Leading/trailing identical runes, no threshold",
		},
		// distance=2 > threshold=1 => returns 1+1=2
		{
			a:         "abcXYZ",
			b:         "abcWXY",
			threshold: 1,
			want:      2,
			name:      "Leading/trailing identical runes, with threshold",
		},
	}
}

func TestComputeDistance(t *testing.T) {
	cases := computeDistanceTestCases()

	for _, caseEntry := range cases {
		t.Run(caseEntry.name, func(t *testing.T) {
			got := ComputeDistance(caseEntry.a, caseEntry.b, []int{caseEntry.threshold}...)
			if got != caseEntry.want {
				t.Errorf("TestComputeDistance(%q, %q, %v) = %d, want %d",
					caseEntry.a, caseEntry.b, caseEntry.threshold, got, caseEntry.want)
			}
		})
	}
}

func computeDistanceUnicodeTestCases() []testCase {
	return []testCase{
		{
			a:         "resumé and café",
			b:         "resumés and cafés",
			threshold: -1,
			want:      2,
			name:      "UC #1",
		},
		{
			a:         "resume and cafe",
			b:         "resumé and café",
			threshold: -1,
			want:      2,
			name:      "UC #2",
		},
		{
			a:         "Hafþór Júlíus Björnsson",
			b:         "Hafþor Julius Bjornsson",
			threshold: -1,
			want:      4,
			name:      "UC #3",
		},
		{
			a:         "།་གམ་འས་པ་་མ།",
			b:         "།་གམའས་པ་་མ",
			threshold: -1,
			want:      2,
			name:      "UC #4",
		},
		{
			a:         "Я был на этой планете бесконечным множеством",
			b:         "Я был на этой паланете бесконечным моножеством",
			threshold: -1,
			want:      2,
			name:      "UC #5",
		},
	}
}

func TestComputeDistanceUnicode(t *testing.T) {

	cases := computeDistanceUnicodeTestCases()

	for _, caseEntry := range cases {
		t.Run(caseEntry.name, func(t *testing.T) {
			got := ComputeDistance(caseEntry.a, caseEntry.b, []int{caseEntry.threshold}...)
			if got != caseEntry.want {
				t.Errorf("TestComputeDistanceUnicode(%q, %q, %v) = %d, want %d",
					caseEntry.a, caseEntry.b, caseEntry.threshold, got, caseEntry.want)
			}
		})
	}
}

// Benchmarks

type benchmarkCase struct {
	a         string
	b         string
	threshold int
	name      string
}

func computeDistanceBenchmarkCases(threshold int) []benchmarkCase {
	return []benchmarkCase{
		{
			"levenshtein",
			"frankenstein",
			threshold,
			"ASCII",
		},
		// Testing acutes and umlauts
		{
			"resumé and café",
			"resumés and cafés",
			threshold,
			"French",
		},
		{
			"Hafþór Júlíus Björnsson",
			"Hafþor Julius Bjornsson",
			threshold,
			"Nordic",
		},
		{
			"разными ощущениями и разными стремлениями",
			"разные ощущения и разные стремления",
			threshold,
			"Russian",
		},
		// Only 2 characters are less in the 2nd string
		{
			"།་གམ་འས་པ་་མ།",
			"།་གམའས་པ་་མ", threshold,
			"Tibetan",
		},
		// Long strings
		{
			"a very long string that is meant to exceed",
			"another very long string that is meant to exceed",
			threshold,
			"Long lead",
		},
		{
			"a very long string with a word in the middle that is different",
			"a very long string with some text in the middle that is different",
			threshold,
			"Long middle",
		},
		{
			"a very long string with some text at the end that is not the same",
			"a very long string with some text at the end that is very different",
			threshold,
			"Long trail",
		},
		{
			"+a very long string with different leading and trailing characters+",
			"-a very long string with different leading and trailing characters-",
			threshold,
			"Long diff",
		},
	}
}

func BenchmarkComputeDistance(b *testing.B) {
	cases := computeDistanceBenchmarkCases(-1)

	for _, caseEntry := range cases {
		b.Run(caseEntry.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				ComputeDistance(caseEntry.a, caseEntry.b)
			}
		})
	}
}

func BenchmarkComputeDistanceWithThreshold(b *testing.B) {
	cases := computeDistanceBenchmarkCases(2)

	for _, caseEntry := range cases {
		b.Run(caseEntry.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				ComputeDistance(caseEntry.a, caseEntry.b, caseEntry.threshold)
			}
		})
	}
}

func competitorsBenchmarkCasesShort(threshold int) []benchmarkCase {
	return []benchmarkCase{
		// ASCII
		{
			"levenshtein",
			"frankenstein",
			threshold,
			"ASCII short",
		},
		// Testing acutes and umlauts
		{
			"resumé and café",
			"resumés and cafés",
			threshold,
			"French short",
		},
		{
			"Hafþór Júlíus Björnsson",
			"Hafþor Julius Bjornsson",
			threshold,
			"Nordic short",
		},
		// Only 2 characters are less in the 2nd string
		{
			"།་གམ་འས་པ་་མ།",
			"།་གམའས་པ་་མ",
			threshold,
			"Tibetan short",
		},
		{
			"разными ощущениями и разными стремлениями",
			"разные ощущения и разные стремления",
			threshold,
			"Russian short",
		},
	}
}

func competitorsBenchmarkCasesLong(threshold int) []benchmarkCase {
	return []benchmarkCase{
		{
			"I am so deeply fulfilled in my understanding that it feels as though I have lived for hundreds of trillions of billions of years on trillions upon trillions of planets just like this Earth. This world is completely clear to me, and here I seek only one thing—peace, tranquility, and harmony, from merging with the infinitely eternal, from contemplating the grand fractal resemblance, and from this marvelous unity of existence. I am so deeply fulfilled in my understanding that it feels as though I have lived for hundreds of trillions of billions of years on trillions upon trillions of planets just like this Earth. This world is completely clear to me, and here I seek only one thing—peace, tranquility, and harmony, from merging with the infinitely eternal, from contemplating the grand fractal resemblance, and from this marvelous unity of existence.",
			"I am so deeply fulfilled in my understanding that it feels as though I have lived for hundreds of trillions of billions of years on trillions upon trillions of planets just like this Earth. And here I seek only one thing—peace, tranquility, and harmony, from merging with the infinitely eternal, from contemplating the grand fractal resemblance, and from this marvelous unity of existence, infinitely eternal. I am so deeply fulfilled in my understanding that it feels as though I have lived for hundreds of trillions of billions of years on trillions upon trillions of planets just like this Earth. And here I seek only one thing—peace, tranquility, and harmony, from merging with the infinitely eternal, from contemplating the grand fractal resemblance, and from this marvelous unity of existence, infinitely eternal.",
			threshold,
			"ASCII long",
		},
		{
			"Я в своем познании настолько преисполнился, что я как будто бы уже сто триллионов миллиардов лет проживаю на триллионах и триллионах таких же планет, как эта Земля, мне этот мир абсолютно понятен, и я здесь ищу только одного - покоя, умиротворения и вот этой гармонии, от слияния с бесконечно вечным, от созерцания великого фрактального подобия и от вот этого замечательного всеединства существа. Я в своем познании настолько преисполнился, что я как будто бы уже сто триллионов миллиардов лет проживаю на триллионах и триллионах таких же планет, как эта Земля, мне этот мир абсолютно понятен, и я здесь ищу только одного - покоя, умиротворения и вот этой гармонии, от слияния с бесконечно вечным, от созерцания великого фрактального подобия и от вот этого замечательного всеединства существа.",
			"Я в своем познании настолько преисполнился, что я как будто бы уже сто триллионов миллиардов лет проживаю на триллионах и триллионах таких же планет, как эта Земля и я здесь ищу только одного - покоя, умиротворения и вот этой гармонии, от слияния с бесконечно вечным, от созерцания великого фрактального подобия и от вот этого замечательного всеединства существа, бесконечно вечного. Я в своем познании настолько преисполнился, что я как будто бы уже сто триллионов миллиардов лет проживаю на триллионах и триллионах таких же планет, как эта Земля и я здесь ищу только одного - покоя, умиротворения и вот этой гармонии, от слияния с бесконечно вечным, от созерцания великого фрактального подобия и от вот этого замечательного всеединства существа, бесконечно вечного.",
			threshold,
			"Russian long",
		},
	}
}

func benchmarkCompetitors(b *testing.B, cases []benchmarkCase) {
	for _, caseEntry := range cases {
		b.Run(caseEntry.name, func(b *testing.B) {
			b.Run("eaxis", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					ComputeDistance(caseEntry.a, caseEntry.b, caseEntry.threshold)
				}
			})
			b.Run("agniva", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					agnivade.ComputeDistance(caseEntry.a, caseEntry.b)
				}
			})
			b.Run("arbovm", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					arbovm.Distance(caseEntry.a, caseEntry.b)
				}
			})
			b.Run("dgryski", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					dgryski.Levenshtein([]rune(caseEntry.a), []rune(caseEntry.b))
				}
			})
		})
	}
}

func BenchmarkCompetitors(b *testing.B) {
	cases := append(
		competitorsBenchmarkCasesShort(-1),
		competitorsBenchmarkCasesLong(-1)...,
	)

	benchmarkCompetitors(b, cases)
}

func BenchmarkCompetitorsWithThreshold(b *testing.B) {
	cases := append(
		competitorsBenchmarkCasesShort(2),
		competitorsBenchmarkCasesLong(10)...,
	)

	benchmarkCompetitors(b, cases)
}

// Fuzzing

func FuzzComputeDistanceDifferent(f *testing.F) {
	testcases := []struct{ a, b string }{
		{"levenshtein", "frankenstein"},
		{"resumé and café", "resumés and cafés"},
		{"Hafþór Júlíus Björnsson", "Hafþor Julius Bjornsson"},
		{"།་གམ་འས་པ་་མ།", "།་གམའས་པ་་མ"},
		{`_p~𕍞`, `b잖PwN`},
		{`7ȪJR`, `6L)wӝ`},
		{`_p~𕍞`, `Y>q8օ݌`},
	}
	for _, tc := range testcases {
		f.Add(tc.a, tc.b)
	}
	f.Fuzz(func(t *testing.T, a, b string) {
		n := ComputeDistance(a, b)
		if n < 0 {
			t.Errorf("Distance can not be negative: %d, a: %q, b: %q", n, a, b)
		}
		if n > len(a)+len(b) {
			t.Errorf("Distance can not be greater than sum of lengths of a and b: %d, a: %q, b: %q", n, a, b)
		}
	})
}

func FuzzComputeDistanceEqual(f *testing.F) {
	testcases := []string{
		"levenshtein", "frankenstein",
		"resumé and café", "resumés and cafés",
		"Hafþór Júlíus Björnsson", "Hafþor Julius Bjornsson",
		"།་གམ་འས་པ་་མ།", "།་གམའས་པ་་མ",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, a string) {
		n := ComputeDistance(a, a)
		if n != 0 {
			t.Errorf("Distance must be zero: %d, a: %q", n, a)
		}
	})
}
