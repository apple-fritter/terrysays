/**  
*  Author: Ian Shaw (bitwitch)
*
*  Markov chain algorithm taken from https://golang.org/doc/codewalk/markov/
* 
*  That algorithm written by The Go Authors thus I must include:
*  Copyright 2011 The Go Authors. All rights reserved.
*  Use of this source code is governed by a BSD-style
*  license that can be found in the LICENSE file.
*
*  Based on the program presented in the "Design and Implementation" chapter
*  of The Practice of Programming (Kernighan and Pike, Addison-Wesley 1999).
*  See also Computer Recreations, Scientific American 260, 122 - 125 (1989).
*
*  A Markov chain algorithm generates text by creating a statistical model of
*  potential textual suffixes for a given prefix.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
	"log"
	"os"
)

// Prefix is a Markov chain prefix of one or more words
type Prefix []string

// String returns the Prefix as a string (for use as a map key)
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map of prefixes to a list of suffixes
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

// Generate returns a string of at most n words generated from Chain
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

func learnVocabulary(prefixLen int) *Chain {
	// initialize a new Chain
	c := NewChain(prefixLen) 
	
	// Read video transcriptions to build markov chain 
	dirname := "./vocabulary/"
    fs, _ := ioutil.ReadDir(dirname)
    for _, f := range fs { 
        if strings.HasSuffix(f.Name(), ".txt") {
			transcript, err := os.Open(dirname + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			c.Build(transcript)
		}
	}
    return c	
}

func main() {
	// register command-line flags
	numWords := flag.Int("words", 100, "maximum number of words to print")
	prefixLen := flag.Int("prefix", 2, "prefix length in words")
	flag.Parse()

    // seed the random number generator
	rand.Seed(time.Now().UnixNano())
	
	// build markov chain
	c := learnVocabulary(*prefixLen)

	// hey terry, what cha got to say <3
	soanyway := c.Generate(*numWords)
	fmt.Println("Terry says...\n" + soanyway)
}

