package main

var Dialogs = make([]Dialog, 0, 4)

type Dialog interface {
	GetWidth() int
	GetHeight() int
	GetDialogStrings() []string
	PrintLine(lineIndex int)
}

type MenuDialog struct {
}

func (d MenuDialog) GetWidth() {

}

func (d MenuDialog) GetHeight() {

}

func (d MenuDialog) GetDialogStrings() []string {
	return []string{
		" +------------------------+ ",
		" |    Info            [I] | ",
		" |    Rename          [R] | ",
		" |    Delete          [D] | ",
		" |    Parent Folder   [P] | ",
		" |    Close           [C] | ",
		" +------------------------+ ",
	}
}

func (d MenuDialog) PrintLine(lineIndex int) {

}

func addDialogToArr(d Dialog) {
	Dialogs = append(Dialogs, d)
}

func removeLastDialog() {
	if len(Dialogs) > 0 {
		Dialogs = Dialogs[:len(Dialogs)-1]
	}
}

func printDialogs() {

}
