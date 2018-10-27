package Error

// if error do a panic
func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}


