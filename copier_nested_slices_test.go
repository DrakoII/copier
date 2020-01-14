package copier_test

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/copier/ns1"
	"github.com/jinzhu/copier/ns2"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func TestCopyNestedSlices(t *testing.T) {
	str := "test"
	ys := []ns1.Y{{Name: &str}, {Name: &str}}
	obj1 := ns1.Astruct{
		A: &str,
		B: [][]ns1.X{
			{{Name: &str, Ys: &ys}, {Name: &str, Ys: &ys}},
			{{Name: &str, Ys: &ys}, {Name: &str, Ys: &ys}, {Name: &str, Ys: &ys}},
		},
	}
	obj2 := &ns2.Astruct{}
	err := copier.Copy(obj2, &obj1)
	log.Println(err)
	assert.NoError(t, err, "unexpected error when copying")

	objC1 := ns1.Cstruct{
		//A: &str,
		B: [][][]ns1.X{
			{
				{{Name: &str, Ys: &ys}, {Name: &str, Ys: &ys}},
				{{Name: &str, Ys: &ys}},
			}, {
				{{Name: &str, Ys: &ys}, {Name: &str, Ys: &ys}},
				{{Name: &str, Ys: &ys}},
			},
		},
	}
	objC2 := &ns2.Cstruct{}
	err = copier.Copy(objC2, &objC1)
	log.Println(err)
	assert.NoError(t, err, "unexpected error when copying")

	x := [][][]ns1.X{
		{
			{{Name: &str, Ys: &ys}, {Name: &str, Ys: &ys}},
			{{Name: &str, Ys: &ys}},
		}, {
			{{Name: &str, Ys: &ys}, {Name: &str, Ys: &ys}},
			{{Name: &str, Ys: &ys}},
		},
	}
	atx := &x
	atatx := &atx
	objD1 := ns1.Dstruct{
		B: &(atatx),
	}
	objD2 := &ns2.Dstruct{}
	err = copier.Copy(objD2, &objD1)
	log.Println(err)
	assert.NoError(t, err, "unexpected error when copying")

	objF1 := ns1.Fstruct{
		//A: &str,
		B: [][]*ns1.X{
			{&ns1.X{Name: &str, Ys: &ys}, &ns1.X{Name: &str, Ys: &ys}},
			{&ns1.X{Name: &str, Ys: &ys}},
		},
	}
	objF2 := &ns2.Fstruct{}
	err = copier.Copy(objF2, &objF1)
	log.Println(err)
	assert.NoError(t, err, "unexpected error when copying")
}

func TestNestedPointersCopy(t *testing.T) {
	str := "ee"
	ys := []ns1.Y{{Name: &str}, {Name: &str}}
	x := ns1.X{Name: &str, Ys: &ys}
	atX := &x
	atatX := &atX
	obj1 := ns1.Bstruct{
		//A: &str,
		B: &(atatX),
	}
	obj2 := &ns2.Bstruct{}
	err := copier.Copy(obj2, &obj1)
	log.Println(err)
	assert.NoError(t, err, "unexpected error when copying")
}

func TestIndirectType(t *testing.T) {
	str := "ee"
	ys := []ns1.Y{{Name: &str}, {Name: &str}}
	obj1 := ns1.Astruct{
		A: &str,
		B: [][]ns1.X{
			{{Name: &str, Ys: &ys}},
			{{Name: &str, Ys: &ys}, {Name: &str, Ys: &ys}},
		},
	}

	iType := indirectType(reflect.TypeOf(obj1.B))
	fmt.Println(iType.Name())

	jagged := getJaggedSlice()
	assert.Equal(t, (reflect.TypeOf(jagged)).Elem(), indirectType(reflect.TypeOf(jagged)))

}

func
TestSliceDims(t *testing.T) {
	slice1 := []int{1, 2, 2}
	slice2 := [][]int{{1, 2}, {1, 2}}
	slice3 := [][][]int{{{1, 2}, {1, 2}}, {{1}, {2}}, {{}}}

	assert.Equal(t, 0, copier.SliceDims(reflect.TypeOf(2)))
	assert.Equal(t, 1, copier.SliceDims(reflect.TypeOf(slice1)))
	assert.Equal(t, 2, copier.SliceDims(reflect.TypeOf(slice2)))
	assert.Equal(t, 3, copier.SliceDims(reflect.TypeOf(slice3)))
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func getJaggedSlice() [][]ns1.X {
	str := "ee"
	ys := []ns1.Y{{Name: &str}, {Name: &str}}
	pic := make([][]ns1.X, 2) /* type declaration */
	for i := range pic {
		pic[i] = make([]ns1.X, 3) /* again the type? */
		for j := range pic[i] {
			pic[i][j] = ns1.X{Name: &str, Ys: &ys}
		}
	}
	return pic
}
