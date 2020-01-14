package ns2

type X struct {
	Name  string
	AnotherField *string
	Ys *[]Y
}

type Y struct {
	NameChanged *string
}



type Astruct struct {
	A *string
	B [][]X
}

type Bstruct struct {
	B ***X
}

type Cstruct struct {
	B [][][]X
}

type Dstruct struct {
	B ***[][][]X
}

type Fstruct struct {
	B [][]*X
}
