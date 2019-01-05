# lexrankmmr

[![Build Status](https://travis-ci.org/ramenjuniti/lexrank-mmr.svg?branch=master)](https://travis-ci.org/ramenjuniti/lexrank-mmr)

[GoDoc](https://godoc.org/github.com/ramenjuniti/lexrank-mmr)

## Algorithm

[LexRank: Graph-based Lexical Centrality as Salience in Text Summarization](https://www.cs.cmu.edu/afs/cs/project/jair/pub/volume22/erkan04a-html/erkan04a.html)

[The Use of MMR, Diversity-Based Reranking for Reordering Documents and Producing Summaries](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.188.3982&rep=rep1&type=pdf)

## Dependency

- [github.com/dcadenas/pagerank](https://github.com/dcadenas/pagerank)
- [github.com/gaspiman/cosine_similarity](https://github.com/gaspiman/cosine_similarity)
- [github.com/ikawaha/kagome](https://github.com/ikawaha/kagome)

## install

```sh
go get github.com/ramenjuniti/lexrankmmr
```

## Usage

```go
package main

import github.com/ramenjuniti/lexrankmmr

func main() {
    text := "Please input the document you want to summarize here."
    summary, err := lexrank.New(
        lexrank.MaxLines(maxLines),            // option (default 0)
        lexrank.MaxCharacters(maxCharacters),  // option (default 0)
        lexrank.Threshold(threshold),          // option (default 0.1)
        lexrank.Tolerance(tolerance),          // option (default 0.0001)
        lexrank.Damping(damping),              // option (default 0.85)
        lexrank.Lambda(lambda),                // option (default 1.0)
    )
    err = summary.Summarize(text)
    if err != nil {
        log.Fatal(err)
    }
}
```

## License

This software is released under the MIT License, see LICENSE.
