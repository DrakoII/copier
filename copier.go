package copier

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

// Copy copy things
func Copy(toValue interface{}, fromValue interface{}) (err error) {
	var (
		isSlice   bool
		amount    = 1
		from      = indirect(reflect.ValueOf(fromValue))
		to        = indirect(reflect.ValueOf(toValue))
	)

	if !to.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return
	}

	fromType := indirectType(from.Type())
	toType := indirectType(to.Type())

	// Just set it if possible to assign
	// And need to do copy anyway if the type is struct
	//This is the first easy kind of assignments - could a slice be assigned here?
	if fromType.Kind() != reflect.Struct && from.Type().AssignableTo(to.Type()) {
		to.Set(from)
		return
	}

	isSlice2 := fromType.Kind() != reflect.Slice
	fmt.Println(isSlice2)
	//if both are slices, dont return here. Otherwise we expect both to be struct (otherwise return)
	if (fromType.Kind() != reflect.Slice || toType.Kind() != reflect.Slice) && (fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct) {
		return
	}

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i)) //this gives value on i-th index of slice - this could be another slice in 2d slices
				//here, source is the inner slice
			} else {
				source = indirect(from)
			}
			// dest
			dest = indirect(reflect.New(toType).Elem()) //here, dest is a single elem of the underlying type
		} else {
			source = indirect(from)
			dest = indirect(to)
		}

		// check source
		if source.IsValid() {
			fromTypeFields := deepFields(fromType) //this gets all fields in the struct //the fromType here  perhaps not expected to be a slice, but it can be in 2d slices
			//fmt.Println( fromTypeFields)
			// Copy from field to field or method
			for _, field := range fromTypeFields {
				name := field.Name

				fromField := source.FieldByName(name)
				if fromField.IsValid() { //this line is failing - name = "Name", what is source? source must be struct, but it is probably the inner array
					// has field
					if toField := dest.FieldByName(name); toField.IsValid() { //why does toField have to be valid(no-zero)?
						if toField.CanSet() {
							if !set(toField, fromField) {
								if err := Copy(toField.Addr().Interface(), fromField.Interface()); err != nil { //Not sure why the toField has to do .Addr() , but formField does not
									return err
								}
							}
						}
					} else {
						// try to set to method
						var toMethod reflect.Value
						if dest.CanAddr() {
							toMethod = dest.Addr().MethodByName(name)
						} else {
							toMethod = dest.MethodByName(name)
						}

						if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
							toMethod.Call([]reflect.Value{fromField})
						}
					}
				}
			}

			// Copy from method to field
			for _, field := range deepFields(toType) {
				name := field.Name

				var fromMethod reflect.Value
				if source.CanAddr() {
					fromMethod = source.Addr().MethodByName(name)
				} else {
					fromMethod = source.MethodByName(name)
				}

				if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 {
					if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
						values := fromMethod.Call([]reflect.Value{})
						if len(values) >= 1 {
							set(toField, values[0])
						}
					}
				}
			}
		}
		if isSlice { //for copying from struct (or slice) to slice
			t := dest.Type()
			tt := to.Type()
			fmt.Println(t, tt)
			if dest.Kind() == reflect.Slice{
//
			err :=Copy(dest.Addr().Interface(), source.Interface() ) //here we want to copy the contents of source slice to dest slice. Maybe the syntax needs to be a bit different
			if err != nil{
				fmt.Println("unexpected error when copying nested slice")
			}
			//err := Copy(toField.Addr().Interface(), fromField.Interface()); err != nil  //Not sure why the toField has to do .Addr() , but formField does not
				//
			}

				if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}
		}
	}
	return
}

//return dimensions of slice. 0 if type is not slice
func SliceDims(reflectType reflect.Type) int {
	dims := 0

	for reflectType.Kind() == reflect.Slice {
		dims += 1
		reflectType = reflectType.Elem()
	}
	return dims
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = dereferencedType(reflectType); reflectType.Kind() == reflect.Struct { //CHANGED
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) reflect.Type {
	if reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}

	for reflectType.Kind() == reflect.Ptr /*|| reflectType.Kind() == reflect.Slice */ { //CHANGED
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func dereferencedType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

//this does the actual assigning if it can be done. If not false is returned and copy is recursivelly called
func set(to, from reflect.Value) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			//set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			err := scanner.Scan(from.Interface())
			if err != nil {
				return false
			}
		} else if from.Kind() == reflect.Ptr {
			return set(to, from.Elem())
		} else {
			return false //Here if the Set was not successful, we return false and we have not set anything
		}
	}
	return true
}
