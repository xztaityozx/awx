package cmd

type Summary struct {
	Files  []string
	Status bool
}

func (this Summary) ToString() string {
	var rt = "AWX Summary:\n\tStatus: "
	if this.Status {
		rt += "Success\n"
	} else {
		rt += "Failed\n"
	}
	rt += "\tFiles:\n"
	for _, value := range this.Files {
		rt += "\t\t" + value + "\n"
	}

	return rt
}
