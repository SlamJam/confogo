package main

import (
	"errors"

	"gopkg.in/yaml.v3"
)

func NewYamlContainer(data []byte) (Container, error) {
	var doc yaml.Node

	err := yaml.Unmarshal(data, &doc)
	if err != nil {
		return &yamlDoc{}, err
	}

	return &yamlDoc{&doc}, nil
}

type yamlDoc struct {
	root *yaml.Node
}

type yamlNodeAccesser interface {
	Access(*yaml.Node) (*yaml.Node, error)
}

type yamlNodeNameAccesser struct {
	nodeName string
}

func (a *yamlNodeNameAccesser) Access(node *yaml.Node) (*yaml.Node, error) {
	if node.Kind != yaml.MappingNode {
		return nil, errors.New("node's kind is not MappingNode")
	}

	if len(node.Content)%2 == 1 {
		return nil, errors.New("len children must be even")
	}

	pairsCount := len(node.Content) / 2
	for i := 0; i < pairsCount; i++ {
		keyNode := node.Content[2*i]
		valueNode := node.Content[2*i+1]
		if keyNode.Value == a.nodeName {
			return valueNode, nil
		}
	}

	return nil, errors.New("node not found")
}

type yamlNodeIndexAccesser struct {
	nodeIndex int
}

func (a *yamlNodeIndexAccesser) Access(node *yaml.Node) (*yaml.Node, error) {
	if node.Kind != yaml.SequenceNode {
		return nil, errors.New("kind is not SequenceNode")
	}

	if a.nodeIndex >= len(node.Content) {
		return nil, errors.New("too low children")
	}

	return node.Content[a.nodeIndex], nil
}

type yamlPathAccesser struct {
	accessers []yamlNodeAccesser
}

func (p *yamlPathAccesser) Access(root *yaml.Node) (*yaml.Node, error) {
	node := root

	if node.Kind == yaml.DocumentNode {
		if len(node.Content) != 1 {
			return nil, errors.New("malformed DocumentNode")
		}

		node = node.Content[0]
	}

	for _, a := range p.accessers {
		n, err := a.Access(node)
		if err != nil {
			return nil, err
		}
		node = n
	}

	return node, nil
}

type yamlPathAccesserBuilder struct {
	accessers []yamlNodeAccesser
}

func (y *yamlPathAccesserBuilder) build() yamlPathAccesser {
	return yamlPathAccesser{
		accessers: y.accessers,
	}
}

func (y *yamlPathAccesserBuilder) handlePath(p *PathPartName) {
	y.accessers = append(y.accessers, &yamlNodeNameAccesser{nodeName: p.name})
}

func (y *yamlPathAccesserBuilder) handleIndex(p *PathPartIndex) {
	y.accessers = append(y.accessers, &yamlNodeIndexAccesser{nodeIndex: p.index})
}

func yamlPathAccesserFromPath(path Path) yamlPathAccesser {
	v := yamlPathAccesserBuilder{}
	path.Visit(&v)

	return v.build()
}

func (doc *yamlDoc) GetValueAtPath(path Path) (string, error) {
	accesser := yamlPathAccesserFromPath(path)
	node, err := accesser.Access(doc.root)
	if err != nil {
		return "", err
	}

	if node.Kind != yaml.ScalarNode {
		return "", errors.New("node is not ScalarNode")
	}

	return node.Value, nil
}
