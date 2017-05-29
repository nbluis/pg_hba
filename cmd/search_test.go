package cmd

import "testing"

// goals: must return a empty []hbaRule and return the ouput error
func TestProcessFileWithUnexistingFile(t *testing.T) {
	var output, err = processFile("./file_non_exists.conf_hba_mor")

	if len(output) > 0 {
		t.Error("Should returns an empty array")
	}

	if err == nil {
		t.Error("Should returns an error")
	}
}

// goals: create a empty file, check if is empty and return a empty array with no error
func TestProcessFileWithAEmptyOne(t *testing.T) {

	var output, err = processFile("./testfiles/empty_file.conf")

	if len(output) > 0 {
		t.Error("Should returns an empty array")
	}

	if err == nil {
		t.Error("Should returns an error")
	}

}

func TestProcessFileWithWrongConfiguration(t *testing.T) {

	var output, err = processFile("./testfiles/syntax_error.conf")

	if len(output) > 0 {
		t.Error("Should returns an empty array")
	}

	if err != nil {
		t.Error("Should returns an error")
	}
}
