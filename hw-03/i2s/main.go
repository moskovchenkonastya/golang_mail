package main

import(
	"reflect"
	"fmt"
)

func i2s(data interface{}, out interface{}) error {
	
	strData := reflect.ValueOf(data)
	strOut := reflect.ValueOf(out)

	var err error = nil
	
	switch reflect.TypeOf(data).Kind() {
		case reflect.Invalid:
			return fmt.Errorf("No value")

		case reflect.Map:

			if strOut.Type().Kind() == reflect.Struct {
				return fmt.Errorf("Out does not have ptr type!")	
			}
			err = parseStruct(strData, strOut)

		case reflect.Slice:

			if strData.Type().Kind() != reflect.Indirect(strOut).Type().Kind() {
				return fmt.Errorf("Data and Out should have slice type.")
			}
			value := reflect.Indirect(strData)
			sliceValue := reflect.MakeSlice(strOut.Type().Elem(), value.Len(), value.Cap())

			for i := 0; i < value.Len(); i++ {

				err = parseStruct(value.Index(i).Elem(), sliceValue.Index(i))
			}
			strOut.Elem().Set(sliceValue)
			
		default:
			return fmt.Errorf("Not support this type")
		
	}

	return err
}	

func parseStruct(data, out reflect.Value) error {

	var err error = nil
	value := reflect.Indirect(out)

	for i := 0; i < value.NumField(); i++ { 
		
		typeField := value.Type().Field(i)
		elem := data.MapIndex(reflect.ValueOf(typeField.Name)).Elem()

		switch typeField.Type.Kind(){

			case reflect.Int:	
				
				switch elem.Type().Kind(){
					case reflect.Float64:			
						value.Field(i).SetInt(int64(elem.Float()))
					case reflect.Int:			
						value.Field(i).SetInt(elem.Int())
					case reflect.String:
						return fmt.Errorf("Sring вместо Int!")
				}

			case reflect.String:
				if elem.Type().Kind() != reflect.String{
					return fmt.Errorf("Int вместо String!")
				}
				value.Field(i).SetString(elem.String())
			
			case reflect.Bool:
				if elem.Type().Kind() == reflect.String {
					return fmt.Errorf("String вместо Bool!")
				}	
				value.Field(i).SetBool(elem.Bool())
				
			case reflect.Struct:

				if elem.Type().Kind() == reflect.Bool {
					return fmt.Errorf("Bool вместо Struct!")
				}
				strValue := reflect.New(typeField.Type)
				parseStruct(elem, strValue)
				value.Field(i).Set(strValue.Elem())

			case reflect.Slice:
				if elem.Type().Kind() != reflect.Slice {
					return fmt.Errorf("Struct вместо Slice!")
				}	
				sliceValue := reflect.MakeSlice(typeField.Type, elem.Len(), elem.Cap())
				for i := 0; i < elem.Len(); i++ {
					parseStruct(elem.Index(i).Elem(), sliceValue.Index(i))
				}
	
				value.Field(i).Set(sliceValue)
				
			default:
				return fmt.Errorf("Can not parse this type")
		}
	}	
	return err
}