package tparser

import (
	"reflect"
	"testing"

	"github.com/saturn4er/go-parse-types/test_examples"
	. "github.com/smartystreets/goconvey/convey"
	"fmt"
	"github.com/saturn4er/go-parse-types/test_examples/sub_pack"
)

func TestParse(t *testing.T) {
	var parser *packageParse
	var err error
	Convey("Сheck text_examples", t, func() {
		parser, err = New("./test_examples")
		So(err, ShouldBeNil)
		So(err, ShouldBeNil)
		Convey("Check struct `A`", func() {
			aType := reflect.TypeOf(test_examples.A{})
			pAType := reflect.TypeOf(&test_examples.A{})
			ty := parser.getTypeByName("A")
			Convey("Check kind", func() {
				So(ty.Kind, ShouldEqual, aType.Kind())
			})
			Convey("Check fields count", func() {
				So(ty.Fields, ShouldHaveLength, aType.NumField())
			})
			Convey("Check methods count", func() {
				So(ty.Methods, ShouldHaveLength, aType.NumMethod()+pAType.NumMethod())
			})
		})
		Convey("Check func TestFunc", func() {
			fType := reflect.TypeOf(test_examples.TestFunc)
			ty := parser.getTypeByName("TestFunc")
			Convey("Check kind", func() {
				So(ty.Kind, ShouldEqual, fType.Kind())
			})
			Convey("Check in count", func() {
				So(ty.In, ShouldHaveLength, fType.NumIn())
			})
			Convey("Check out count", func() {
				So(ty.Out, ShouldHaveLength, fType.NumOut())
			})
		})
	})

	for _, value := range parser.types {
		fmt.Println("----------")
		fmt.Println(value)
	}
	Convey("Сheck sub_pack", t, func() {
		parser, err = New("./test_examples/sub_pack")
		So(err, ShouldBeNil)
		err = parser.parse()
		So(err, ShouldBeNil)
		Convey("Check struct `B`", func() {
			aType := reflect.TypeOf(sub_pack.B{})
			pAType := reflect.TypeOf(&sub_pack.B{})
			ty := parser.getTypeByName("B")
			Convey("Check kind", func() {
				So(ty.Kind, ShouldEqual, aType.Kind())
			})
			Convey("Check fields count", func() {
				So(ty.Fields, ShouldHaveLength, aType.NumField())
			})
			Convey("Check methods count", func() {
				So(ty.Methods, ShouldHaveLength, aType.NumMethod()+pAType.NumMethod())
			})
		})

	})

}
