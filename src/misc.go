package main

func isExtensionPDF(filepath string) bool {
	if len(filepath) > 4 {
		if filepath[len(filepath)-4:] == ".pdf" {
			return true
		}
	}
	return false
}

func isExtensionHTML(filepath string) bool {
	if len(filepath) > 5 {
		if filepath[len(filepath)-5:] == ".html" {
			return true
		}
	}
	return false
}
