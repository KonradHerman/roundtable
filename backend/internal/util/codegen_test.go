package util

import (
	"strings"
	"testing"
)

func TestGenerateRoomCode_Length(t *testing.T) {
	t.Parallel()

	// Generate multiple codes and verify length
	for i := 0; i < 100; i++ {
		code := GenerateRoomCode()

		if len(code) != 6 {
			t.Errorf("code length = %d, want 6 (code: %s)", len(code), code)
		}
	}
}

func TestGenerateRoomCode_Charset(t *testing.T) {
	t.Parallel()

	// Valid characters (excluding 0, O, 1, I, L for readability)
	validChars := "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	// Generate many codes and verify all characters are valid
	for i := 0; i < 100; i++ {
		code := GenerateRoomCode()

		for _, char := range code {
			if !strings.ContainsRune(validChars, char) {
				t.Errorf("code contains invalid character %c (code: %s)", char, code)
			}
		}
	}
}

func TestGenerateRoomCode_NoAmbiguousCharacters(t *testing.T) {
	t.Parallel()

	// Characters that should NOT appear (ambiguous) - based on actual charset
	// Note: The charset is "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" which excludes 0, O, 1, I
	ambiguousChars := "01IO"

	// Generate many codes and verify no ambiguous characters
	for i := 0; i < 100; i++ {
		code := GenerateRoomCode()

		for _, char := range code {
			if strings.ContainsRune(ambiguousChars, char) {
				t.Errorf("code contains ambiguous character %c (code: %s)", char, code)
			}
		}
	}
}

func TestGenerateRoomCode_Uniqueness(t *testing.T) {
	t.Parallel()

	// Generate many codes and check for uniqueness
	numCodes := 1000
	codes := make(map[string]bool)
	duplicates := 0

	for i := 0; i < numCodes; i++ {
		code := GenerateRoomCode()

		if codes[code] {
			duplicates++
		}
		codes[code] = true
	}

	// With a 6-character code and 32 possible characters (32^6 = ~1 billion combinations),
	// we should see very few duplicates in 1000 attempts
	duplicateRate := float64(duplicates) / float64(numCodes)

	if duplicateRate > 0.01 { // Allow up to 1% duplicate rate (very generous)
		t.Errorf("duplicate rate too high: %.2f%% (%d duplicates in %d codes)",
			duplicateRate*100, duplicates, numCodes)
	}
}

func TestGenerateRoomCode_Randomness(t *testing.T) {
	t.Parallel()

	t.Run("codes are not identical", func(t *testing.T) {
		t.Parallel()

		code1 := GenerateRoomCode()
		code2 := GenerateRoomCode()
		code3 := GenerateRoomCode()

		// It's extremely unlikely that all three codes are the same
		if code1 == code2 && code2 == code3 {
			t.Errorf("all three codes are identical: %s", code1)
		}
	})

	t.Run("character distribution is reasonable", func(t *testing.T) {
		t.Parallel()

		// Generate many codes and count character frequencies
		numCodes := 1000
		charCounts := make(map[rune]int)

		for i := 0; i < numCodes; i++ {
			code := GenerateRoomCode()
			for _, char := range code {
				charCounts[char]++
			}
		}

		// Total characters generated
		totalChars := numCodes * 6

		// With 32 possible characters, expected frequency per char is ~3.125%
		// We'll check that no character appears more than 10% or less than 0.5%
		for char, count := range charCounts {
			frequency := float64(count) / float64(totalChars)

			if frequency > 0.10 {
				t.Errorf("character %c appears too frequently: %.2f%%", char, frequency*100)
			}

			if frequency < 0.005 {
				t.Errorf("character %c appears too rarely: %.2f%%", char, frequency*100)
			}
		}
	})

	t.Run("positional independence", func(t *testing.T) {
		t.Parallel()

		// Generate codes and check that each position has variety
		numCodes := 100
		positions := make([]map[rune]bool, 6)
		for i := 0; i < 6; i++ {
			positions[i] = make(map[rune]bool)
		}

		for i := 0; i < numCodes; i++ {
			code := GenerateRoomCode()
			for pos, char := range code {
				positions[pos][char] = true
			}
		}

		// Each position should have at least a few different characters
		for pos, chars := range positions {
			if len(chars) < 5 {
				t.Errorf("position %d has too few unique characters: %d", pos, len(chars))
			}
		}
	})
}

func TestGenerateRoomCode_UppercaseOnly(t *testing.T) {
	t.Parallel()

	// Generate codes and verify all are uppercase
	for i := 0; i < 100; i++ {
		code := GenerateRoomCode()

		if strings.ToUpper(code) != code {
			t.Errorf("code is not all uppercase: %s", code)
		}

		// Also verify no lowercase letters
		for _, char := range code {
			if char >= 'a' && char <= 'z' {
				t.Errorf("code contains lowercase character %c (code: %s)", char, code)
			}
		}
	}
}

func TestGenerateRoomCode_NoSpaces(t *testing.T) {
	t.Parallel()

	// Generate codes and verify no spaces
	for i := 0; i < 100; i++ {
		code := GenerateRoomCode()

		if strings.Contains(code, " ") {
			t.Errorf("code contains space: %s", code)
		}

		// Also check for other whitespace
		if strings.TrimSpace(code) != code {
			t.Errorf("code contains whitespace: %s", code)
		}
	}
}

func TestGenerateRoomCode_Consistency(t *testing.T) {
	t.Parallel()

	t.Run("always returns a string", func(t *testing.T) {
		t.Parallel()

		for i := 0; i < 10; i++ {
			code := GenerateRoomCode()

			if code == "" {
				t.Error("GenerateRoomCode returned empty string")
			}
		}
	})

	t.Run("no special characters", func(t *testing.T) {
		t.Parallel()

		specialChars := "!@#$%^&*()_+-=[]{}\\|;:'\",.<>?/`~"

		for i := 0; i < 100; i++ {
			code := GenerateRoomCode()

			for _, char := range code {
				if strings.ContainsRune(specialChars, char) {
					t.Errorf("code contains special character %c (code: %s)", char, code)
				}
			}
		}
	})
}

func TestGenerateRoomCode_Performance(t *testing.T) {
	t.Parallel()

	// This test just verifies it doesn't hang or panic
	// Generate a reasonable number of codes quickly
	for i := 0; i < 10000; i++ {
		_ = GenerateRoomCode()
	}
}

func BenchmarkGenerateRoomCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateRoomCode()
	}
}

func TestGenerateRoomCode_StatisticalProperties(t *testing.T) {
	t.Parallel()

	t.Run("no repeating patterns", func(t *testing.T) {
		t.Parallel()

		// Generate codes and check for suspicious patterns
		numCodes := 100

		for i := 0; i < numCodes; i++ {
			code := GenerateRoomCode()

			// Check for all same character (AAAAAA)
			allSame := true
			firstChar := rune(code[0])
			for _, char := range code {
				if char != firstChar {
					allSame = false
					break
				}
			}

			if allSame {
				t.Errorf("code has all same characters: %s", code)
			}

			// Check for simple sequential patterns (ABCDEF, 123456)
			// This is a very basic check; real sequential detection would be more complex
			isSequential := true
			for i := 1; i < len(code); i++ {
				if code[i] != code[i-1]+1 {
					isSequential = false
					break
				}
			}

			if isSequential {
				t.Errorf("code appears sequential: %s", code)
			}
		}
	})
}

func TestGenerateRoomCode_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("concurrent generation", func(t *testing.T) {
		t.Parallel()

		// Generate codes concurrently to ensure thread safety
		numGoroutines := 10
		codesPerGoroutine := 100

		results := make(chan string, numGoroutines*codesPerGoroutine)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				for j := 0; j < codesPerGoroutine; j++ {
					code := GenerateRoomCode()
					results <- code
				}
			}()
		}

		// Collect all codes
		codes := make([]string, 0, numGoroutines*codesPerGoroutine)
		for i := 0; i < numGoroutines*codesPerGoroutine; i++ {
			code := <-results
			codes = append(codes, code)

			// Verify each code is valid
			if len(code) != 6 {
				t.Errorf("concurrent generation produced invalid code length: %d", len(code))
			}
		}
	})
}

