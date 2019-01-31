package main;

import (
  "fmt"
  "github.com/docopt/docopt-go"
)

var usage = `dgraph live cache warmup utility.

Usage:
  dgraph-xidmap-warmup PREDICATE... --dir <path> --connect <url>

Arguments:
  PREDICATE predicate(s) holding hash keys

Options:
  -h --help        Show this screen.
  --version        Show version information and exit.
  --dir <dir>      xidmap path.
  --connect <url>  DGraph url to connect.
`;

func main() {
  arguments, _ := docopt.ParseArgs(usage, nil, "1.0")
  fmt.Println(arguments["PREDICATE"])
}
