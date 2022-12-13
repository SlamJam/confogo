package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func NewEnvContainer(prefix string) (Container, error) {
	c := envDoc{}
	c.root = make(map[string]string, 0)

	envs := os.Environ()
	for _, e := range envs {
		kv := strings.SplitN(e, "=", 2)

		if len(kv) != 2 {
			return &c, errors.New("fail format")
		}
		key := kv[0]
		value := kv[1]

		if strings.HasPrefix(key, prefix) {
			if _, ok := c.root[key]; ok {
				return &c, errors.New("key already exists")
			}
			c.root[strings.TrimPrefix(key, prefix)] = value
		}
	}

	return &c, nil
}

type envDoc struct {
	root map[string]string
}

type envPathBuilder struct {
	parts []string
}

func (y *envPathBuilder) build() string {
	return strings.Join(y.parts, ".")
}

func (y *envPathBuilder) handlePath(p *PathPartName) {
	y.parts = append(y.parts, p.name)
}

func (y *envPathBuilder) handleIndex(p *PathPartIndex) {
	y.parts = append(y.parts, fmt.Sprintf("[%d]", p.index))
}

func envPathFromPath(path Path) string {
	v := envPathBuilder{}
	path.Visit(&v)

	return v.build()
}

func (doc *envDoc) GetValueAtPath(path Path) (string, error) {
	envPath := envPathFromPath(path)

	value, ok := doc.root[envPath]
	if ok {
		return value, nil
	}

	return "", errors.New("path not found")
}
