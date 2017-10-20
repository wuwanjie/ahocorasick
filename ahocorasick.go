package goahocorasick

import (
	"fmt"
)

import (
	"github.com/anknown/darts"
)

const FAIL_STATE = -1
const ROOT_STATE = 1

type Machine struct {
	Trie    *godarts.DoubleArrayTrie
	Failure []int
	Output  map[int]([][]rune)
}

type Term struct {
	Pos  int
	Word []rune
}

func (m *Machine) Build(keywords [][]rune) (err error) {
	if len(keywords) == 0 {
		return fmt.Errorf("empty keywords")
	}

	d := new(godarts.Darts)

	trie := new(godarts.LinkedListTrie)
	m.Trie, trie, err = d.Build(keywords)
	if err != nil {
		return err
	}

	m.Output = make(map[int]([][]rune), 0)
	for idx, val := range d.Output {
		m.Output[idx] = append(m.Output[idx], val)
	}

	queue := make([](*godarts.LinkedListTrieNode), 0)
	m.Failure = make([]int, len(m.Trie.Base))
	for _, c := range trie.Root.Children {
		m.Failure[c.Base] = godarts.ROOT_NODE_BASE
	}
	queue = append(queue, trie.Root.Children...)
	for {
		if len(queue) == 0 {
			break
		}

		node := queue[0]
		for _, n := range node.Children {
			if n.Base == godarts.END_NODE_BASE {
				continue
			}
			inState := m.f(node.Base)
		set_state:
			outState := m.g(inState, n.Code-godarts.ROOT_NODE_BASE)
			if outState == FAIL_STATE {
				inState = m.f(inState)
				goto set_state
			}
			if _, ok := m.Output[outState]; ok != false {
				m.Output[n.Base] = append(m.Output[outState], m.Output[n.Base]...)
			}
			m.setF(n.Base, outState)
		}
		queue = append(queue, node.Children...)
		queue = queue[1:]
	}

	return nil
}

func (m *Machine) PrintFailure() {
	fmt.Printf("+-----+-----+\n")
	fmt.Printf("|%5s|%5s|\n", "index", "value")
	fmt.Printf("+-----+-----+\n")
	for i, v := range m.Failure {
		fmt.Printf("|%5d|%5d|\n", i, v)
	}
	fmt.Printf("+-----+-----+\n")
}

func (m *Machine) PrintOutput() {
	fmt.Printf("+-----+----------+\n")
	fmt.Printf("|%5s|%10s|\n", "index", "value")
	fmt.Printf("+-----+----------+\n")
	for i, v := range m.Output {
		var val string
		for _, o := range v {
			val = val + " " + string(o)
		}
		fmt.Printf("|%5d|%10s|\n", i, val)
	}
	fmt.Printf("+-----+----------+\n")
}

func (m *Machine) g(inState int, input rune) (outState int) {
	if inState == FAIL_STATE {
		return ROOT_STATE
	}

	t := inState + int(input) + godarts.ROOT_NODE_BASE
	if t >= len(m.Trie.Base) {
		if inState == ROOT_STATE {
			return ROOT_STATE
		}
		return FAIL_STATE
	}
	if inState == m.Trie.Check[t] {
		return m.Trie.Base[t]
	}

	if inState == ROOT_STATE {
		return ROOT_STATE
	}

	return FAIL_STATE
}

func (m *Machine) f(index int) (state int) {
	return m.Failure[index]
}

func (m *Machine) setF(inState, outState int) {
	m.Failure[inState] = outState
}

func (m *Machine) MultiPatternSearch(content []rune, returnImmediately bool) [](*Term) {
	terms := make([](*Term), 0)

	state := ROOT_STATE
	for pos, c := range content {
	start:
		if m.g(state, c) == FAIL_STATE {
			state = m.f(state)
			// failure不应该有0的状态呢
			if state == 0 {
				break
			}
			goto start
		} else {
			state = m.g(state, c)
			if val, ok := m.Output[state]; ok != false {
				for _, word := range val {
					term := new(Term)
					term.Pos = pos - len(word) + 1
					term.Word = word
					terms = append(terms, term)
					if returnImmediately {
						return terms
					}
				}
			}
		}
	}

	return terms
}

func (m *Machine) ExactSearch(content []rune) [](*Term) {
	if m.Trie.ExactMatchSearch(content, 0) {
		t := new(Term)
		t.Word = content
		t.Pos = 0
		return [](*Term){t}
	}

	return nil
}
