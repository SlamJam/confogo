package main

/*
 Абстрактный "путь" до значения.
 Есть два типа аксесоров: по имени и по индексу.
 Собственно, в дереве мы можем выбрать узел или по имени или по индексу.
*/

type PathPartVisitor interface {
	handlePath(*PathPartName)
	handleIndex(*PathPartIndex)
}

type PathPartVisitable interface {
	Visit(v PathPartVisitor)
}

// Часть (part) пути, адресующая по имени
type PathPartName struct {
	name string
}

func (p *PathPartName) Visit(v PathPartVisitor) {
	v.handlePath(p)
}

func WithPath(part string) *PathPartName {
	return &PathPartName{name: part}
}

var _ PathPartVisitable = &PathPartName{}

// Часть (part) пути, адресующая по индексу
type PathPartIndex struct {
	index int
}

func (p *PathPartIndex) Visit(v PathPartVisitor) {
	v.handleIndex(p)
}

func WithIndex(index int) *PathPartIndex {
	return &PathPartIndex{index: index}
}

var _ PathPartVisitable = &PathPartIndex{}

// Полный путь, состоящий из "кусочков" (parts)
type Path struct {
	parts []PathPartVisitable
}

func (path *Path) Visit(v PathPartVisitor) {
	for _, part := range path.parts {
		part.Visit(v)
	}
}

// assert
var _ PathPartVisitable = &Path{}

// Публичный конструктор, собирающий путь из кусочков
func NewPath(parts ...PathPartVisitable) Path {
	return Path{parts: parts}
}
