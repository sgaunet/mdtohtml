package main

func isExtensionHTML(filepath string) bool {
	if len(filepath) > 5 {
		if filepath[len(filepath)-5:] == ".html" {
			return true
		}
	}
	return false
}
