package ghtml_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestParse(t *testing.T) {
	s := `<p>Links:</p>
	<ul>
		<li>
		<a href="foo">Foo</a>
		<li>
		<a href="/bar/baz">BarBaz</a>
	</ul>
	`
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			msg := ""
			msg += "Data:" + n.Data
			for _, at := range n.Attr {
				msg += "[" + at.Key + ", " + at.Namespace + ", " + at.Val + "]"
			}
			fmt.Println(msg)
		} else {
			switch n.Type {
			case html.TextNode:
				fmt.Println("type:TextNode")
			case html.DocumentNode:
				fmt.Println("type:DocumentNode")
			case html.CommentNode:
				fmt.Println("type:CommentNode")
			case html.DoctypeNode:
				fmt.Println("type:DoctypeNode")
			default:
				fmt.Println("type:ErrorNode")
			}
		}

		/*
			if n.Type == html.ElementNode && n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "href" {
						fmt.Println(a.Val)
						break
					}
				}
			}
			//*/
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fmt.Println("------")
			f(c)
		}
	}
	f(doc)
}

func TestLoop(t *testing.T) {
	i := 0
	j := 0

END:
	fmt.Println(">> START")
LOOP:
	for {
		fmt.Println("i :", i)
		j = 0
		i++
		for {
			fmt.Println(" j :", j)

			if i == 9 {
				goto END
			}
			if i == 10 {
				break LOOP
			}

			j++
			if j == 3 {
				break
			}
		}
	}

	fmt.Println(">> end loop")
}
