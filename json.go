package main

import (
	"encoding/json"
	"errors"
)

func NewJsonContainer(data []byte) (Container, error) {
	var doc interface{}

	err := json.Unmarshal(data, &doc)
	if err != nil {
		return &jsonDoc{}, err
	}

	return &jsonDoc{doc}, nil
}

type jsonDoc struct {
	root interface{}
}

type jsonNode interface{}

type jsonNodeAccesser interface {
	Access(jsonNode) (jsonNode, error)
}

type jsonNodeNameAccesser struct {
	nodeName string
}

func (a *jsonNodeNameAccesser) Access(node jsonNode) (jsonNode, error) {
	nodeMap, ok := node.(map[string]any)
	if !ok {
		return nil, errors.New("node type is not map")
	}

	res, ok := nodeMap[a.nodeName]
	if !ok {
		return nil, errors.New("node not found")
	}

	return res, nil
}

type jsonNodeIndexAccesser struct {
	nodeIndex int
}

func (a *jsonNodeIndexAccesser) Access(node jsonNode) (jsonNode, error) {
	nodeArray, ok := node.([]any)
	if !ok {
		return nil, errors.New("node type is not array")
	}

	if a.nodeIndex >= len(nodeArray) {
		return nil, errors.New("array len too small")
	}

	return nodeArray[a.nodeIndex], nil
}

type jsonPathAccesser struct {
	accessers []jsonNodeAccesser
}

func (p *jsonPathAccesser) Access(root jsonNode) (jsonNode, error) {
	node := root

	if node == nil {
		return nil, errors.New("empty document")
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

type jsonPathAccesserBuilder struct {
	accessers []jsonNodeAccesser
}

func (y *jsonPathAccesserBuilder) build() jsonPathAccesser {
	return jsonPathAccesser{
		accessers: y.accessers,
	}
}

func (y *jsonPathAccesserBuilder) handlePath(p *PathPartName) {
	y.accessers = append(y.accessers, &jsonNodeNameAccesser{nodeName: p.name})
}

func (y *jsonPathAccesserBuilder) handleIndex(p *PathPartIndex) {
	y.accessers = append(y.accessers, &jsonNodeIndexAccesser{nodeIndex: p.index})
}

func jsonPathAccesserFromPath(path Path) jsonPathAccesser {
	v := jsonPathAccesserBuilder{}
	path.Visit(&v)

	return v.build()
}

func (doc *jsonDoc) GetValueAtPath(path Path) (string, error) {
	accesser := jsonPathAccesserFromPath(path)
	node, err := accesser.Access(doc.root)
	if err != nil {
		return "", err
	}

	switch node.(type) {
	case string:
		return node.(string), nil
	default:
		return "", errors.New("node is not string")
	}
}
