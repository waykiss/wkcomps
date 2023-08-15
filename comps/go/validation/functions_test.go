package validation

import (
	"testing"
)

// TestIsValidEan tests the IsValidEan function
func TestIsValidEan(t *testing.T) {
	expectedResults := map[string]bool{
		"40700719670720": true,  //ean14
		"5012345678900":  true,  //ean13
		"012345678905":   true,  //ean12
		"78912342":       true,  //ean8
		"":               true,  //EMPTY ean
		"40700719670721": false, //INVALID ean14
		"5012345678901":  false, //INVALID ean13
		"012345678906":   false, //INVALID ean12
		"78912343":       false, //INVALID ean8
		"SEM GTIM":       false, //INVALID ean
		"1":              false, //INVALID ean
	}
	for ean, expected := range expectedResults {
		if IsValidEan(ean) != expected {
			t.Errorf("Fail test for EAN %s, expected : %t got: %t ", ean, !expected, expected)
		}
	}
}
