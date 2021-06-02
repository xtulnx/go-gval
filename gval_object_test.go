package gval

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
)

func showMethod(rv1 reflect.Type) {
	nm := rv1.NumMethod()
	fmt.Printf("[%s] NumMethod: %d\n", rv1.Name(), nm)
	for i := nm; i > 0; i-- {
		fmt.Printf("%v\n", rv1.Method(i-1))
	}
	fmt.Println("-------")
}

func TestObject(t *testing.T) {
	f := Foo1{}
	//showMethod(reflect.TypeOf(f))
	//showMethod(reflect.TypeOf(&f))

	fmt.Println(f.F1(2).F2(4).F3(6).F4())

	testEvaluate(
		[]evaluationTest{
			{
				name:       "测试001",
				expression: `f("温度高于30℃，当前值:%g, 请及时检查！", 76.41234)`,
				parameter: map[string]interface{}{
					"f": func(f string, args ...interface{}) string {
						return fmt.Sprintf(f, args...)
					},
				},
				want: "温度高于30℃，当前值:76.41234, 请及时检查！",
			},
			{
				name:       "函数返回对象继续引用",
				expression: `f("Hello world", 0).F1(2).F2(4).F3(6).F4()`,
				parameter: map[string]interface{}{
					"f": func(f string, c float64) *Foo1 {
						log.Printf("fff:%s\n", f)
						return &Foo1{}
					},
				},
				want: "Foo4->Foo3->Foo2->Foo1:151",
			},
			{
				name:       "函数返回对象继续引用",
				expression: "f.F1(2).F2(4).F3(6).F4()",
				parameter: map[string]interface{}{
					"f": &f,
				},
				want: "Foo4->Foo3->Foo2->Foo1:151",
			},
		},
		t,
	)
}

type FooX interface {
	Super() FooX
	Name() string
}
type Foo1 struct {
}
type Foo2 struct {
	FooX
	Msg int
}
type Foo3 struct {
	FooX
	Msg int
}
type Foo4 struct {
	FooX
	Msg int
}

func (Foo1) Super() FooX {
	return nil
}
func (Foo1) Name() string {
	return "Foo1"
}
func (f *Foo1) F1(v float64) *Foo2 {
	return &Foo2{
		FooX: f,
		Msg:  3 + int(v),
	}
}

func (f *Foo2) Super() FooX {
	return f.FooX
}
func (Foo2) Name() string {
	return "Foo2"
}
func (f *Foo2) F2(v float64) *Foo3 {
	return &Foo3{
		FooX: f,
		Msg:  f.Msg*5 + int(v),
	}
}

func (Foo3) Name() string {
	return "Foo3"
}
func (f *Foo3) Super() FooX {
	return f.FooX
}
func (f *Foo3) F3(v float64) *Foo4 {
	return &Foo4{
		FooX: f,
		Msg:  f.Msg*5 + int(v),
	}
}

func (f *Foo4) Super() FooX {
	return f.FooX
}
func (f *Foo4) Name() string {
	return "Foo4"
}
func (f *Foo4) F4() string {
	s1 := []string{}
	var f1 FooX = f
	for f1 != nil {
		s1 = append(s1, f1.Name())
		f1 = f1.Super()
	}
	return fmt.Sprintf("%s:%d", strings.Join(s1, "->"), f.Msg)
}
