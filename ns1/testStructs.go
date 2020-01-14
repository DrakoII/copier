package ns1

type X struct {
	Name  *string
	Ys *[]Y
}

type Y struct {
	Name *string
}



type AstructSameName struct {
	//A *string
	B [][]X
}

type Bstruct struct {
	B ***X
}