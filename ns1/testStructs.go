package ns1

type X struct {
	Name  *string
	Ys *[]Y
}

type Y struct {
	Name *string
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