package util

func CheckErr(err error) {
	if err == nil {
		return
	}
	handleErr(err)
}

func handleErr(e error) {

}

