package env

import "os"

var API_KEY_GOOGLE = ""

func Load() {

	API_KEY_GOOGLE = os.Getenv("API_KEY_GOOGLE")

	if API_KEY_GOOGLE == "" {
		println("API_KEY_GOOGLE not found")
	}
}
