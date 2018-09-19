package main

import "os"

// WriteStringToFile writes a string to a file
func WriteStringToFile(contents string) {
	// For more granular writes, open a file for writing.
	f, err := os.Create("/tmp/dat2")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Write([]byte(contents))
	if err != nil {
		panic(err)
	}

}
