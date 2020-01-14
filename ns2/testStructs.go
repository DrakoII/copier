package ns2

type X struct {
	NameChange  *string
	AnotherField *string
	Ys *[]Y
}

type Y struct {
	NameChanged *string
}



type AstructSameName struct {
	//A *string
	B [][]X
}

type Bstruct struct {
	B ***X
}
