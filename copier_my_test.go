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

func TestSameNameCopyAcrossNamespaces(t *testing.T) {
	str := "ee"
	ys := []ns1.Y{{Name: &str}, {Name: &str}}
	obj1 := ns1.AstructSameName{
		//A: &str,
		B: [][]ns1.X{
			{{Name: &str, Ys: &ys}, {Name: &str, Ys: &ys}},
			{{Name: &str, Ys: &ys}, {Name: &str, Ys: &ys}, {Name: &str, Ys: &ys}},
		},
	}
	obj2 := &ns2.AstructSameName{}
	err := copier.Copy(obj2, &obj1)
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
	obj1 := ns1.AstructSameName{
		//A: &str,
		B: [][]ns1.X{
			{{Name: &str, Ys: &ys}},
			{{Name: &str, Ys: &ys}, {Name: &str, Ys: &ys}},
		},
	}

	//simpleObj := []ns1.X{{Name: &str, Ys: &ys}}

	iType := indirectType(reflect.TypeOf(obj1.B))
	fmt.Println(iType.Name())

	jagged := getJaggedSlice()

	//assert.Equal(t, iType, indirectType(reflect.TypeOf(simpleObj)))
	//assert.Equal(t, (reflect.TypeOf(jagged)).Elem(), indirectType(reflect.TypeOf(simpleObj)))
	assert.Equal(t, (reflect.TypeOf(jagged)).Elem(), indirectType(reflect.TypeOf(jagged)))

}

func TestSliceDims(t *testing.T){
	slice1 := []int {1,2,2}
	slice2:= [][]int {{1,2},{1,2}}
	slice3:= [][][]int {{{1,2},{1,2}},{{1},{2}},{{}}}

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

/*type x1 struct{ Name *string }
type x2 struct {
	NameChanged  *string
	anotherField string
}

type astructSameName1 struct {
	A string
	B *[]x1
}

type astructSameName2 struct {
	A string
	B *[]x2
}

type bstructSameName1 struct {
	A string
	B int64
	C time.Time
}

type bstructSameName2 struct {
	A string
	B time.Time
	C int64
}

func TestSameNameCopy(t *testing.T) {
	str := "ee"
	obj1 := astructSameName1{A: "123", B: &[]x1{{Name: &str}}}
	obj2 := &astructSameName2{}
	err := copier.Copy(obj2, &obj1)
	log.Println(err)
	assert.NoError(t, err, "unexpected error when copying")
}*/

/*
func TestCopyFieldsWithSameNameButDifferentTypes2(t *testing.T) {
	obj1 := bstructSameName1{A: "123", B: 2, C: time.Now()}
	obj2 := &bstructSameName2{}
	err := copier.Copy(obj2, &obj1)
	if err != nil {
		t.Error("Should not raise error")
	}

	if obj2.A != obj1.A {
		t.Errorf("Field A should be copied")
	}
}

/*func TestSameNameSliceCopy(t *testing.T) {
	str := "ee"
	slice1 := []astructSameName1{{A: "123", B: &[]x1{{Name: &str}}}, {A: "123", B: &[]x1{{Name: &str}}}}
	slice2 := []astructSameName2{}
	err := copier.Copy(slice2, &slice1)
	log.Println(err)
	assert.NoError(t, err, "unexpected error when copying slices")
}*/

/*************************************

/
*/
/*
type TypeStruct1 struct {
	Field1 string
	Field2 string
	Field3 TypeStruct2
	Field4 *TypeStruct2
	Field5 []*TypeStruct2
	Field6 []TypeStruct2
	Field7 []*TypeStruct2
	Field8 []TypeStruct2
}

type TypeStruct2 struct {
	Field1 int
	Field2 string
}

type TypeStruct3 struct {
	Field1 interface{}
	Field2 string
	Field3 TypeStruct4
	Field4 *TypeStruct4
	Field5 []*TypeStruct4
	Field6 []*TypeStruct4
	Field7 []TypeStruct4
	Field8 []TypeStruct4
}

type TypeStruct4 struct {
	field1 int
	Field2 string
}

func (t *TypeStruct4) Field1(i int) {
	t.field1 = i
}

func TestCopyDifferentFieldType(t *testing.T) {
	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
	}
	ts2 := &TypeStruct2{}

	copier.Copy(ts2, ts)

	if ts2.Field2 != ts.Field2 || ts2.Field1 != 0 {
		t.Errorf("Should be able to copy from ts to ts2")
	}
}

func TestCopyDifferentTypeMethod(t *testing.T) {
	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
	}
	ts4 := &TypeStruct4{}

	copier.Copy(ts4, ts)

	if ts4.Field2 != ts.Field2 || ts4.field1 != 0 {
		t.Errorf("Should be able to copy from ts to ts4")
	}
}

func TestAssignableType(t *testing.T) {
	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
		Field3: TypeStruct2{
			Field1: 666,
			Field2: "str2",
		},
		Field4: &TypeStruct2{
			Field1: 666,
			Field2: "str2",
		},
		Field5: []*TypeStruct2{
			{
				Field1: 666,
				Field2: "str2",
			},
		},
		Field6: []TypeStruct2{
			{
				Field1: 666,
				Field2: "str2",
			},
		},
		Field7: []*TypeStruct2{
			{
				Field1: 666,
				Field2: "str2",
			},
		},
	}

	ts3 := &TypeStruct3{}

	copier.Copy(&ts3, &ts)

	if v, ok := ts3.Field1.(string); !ok {
		t.Error("Assign to interface{} type did not succeed")
	} else if v != "str1" {
		t.Error("String haven't been copied correctly")
	}

	if ts3.Field2 != ts.Field2 {
		t.Errorf("Field2 should be copied")
	}

	checkType2WithType4(ts.Field3, ts3.Field3, t, "Field3")
	checkType2WithType4(*ts.Field4, *ts3.Field4, t, "Field4")

	for idx, f := range ts.Field5 {
		checkType2WithType4(*f, *(ts3.Field5[idx]), t, "Field5")
	}

	for idx, f := range ts.Field6 {
		checkType2WithType4(f, *(ts3.Field6[idx]), t, "Field6")
	}

	for idx, f := range ts.Field7 {
		checkType2WithType4(*f, ts3.Field7[idx], t, "Field7")
	}

	for idx, f := range ts.Field8 {
		checkType2WithType4(f, ts3.Field8[idx], t, "Field8")
	}
}

func checkType2WithType4(t2 TypeStruct2, t4 TypeStruct4, t *testing.T, testCase string) {
	if t2.Field1 != t4.field1 || t2.Field2 != t4.Field2 {
		t.Errorf("%v: type struct 4 and type struct 2 is not equal", testCase)
	}
}*/
