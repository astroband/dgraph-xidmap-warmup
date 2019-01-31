package main;

import (
  "fmt"
  "github.com/docopt/docopt-go"
)

var usage = `dgraph live cache warmup utility.

Usage:
  dgraph-xidmap-warmup PREDICATE... --dir <dir> --url <url>

Arguments:
  PREDICATE predicate(s) holding hash keys

Options:
  -h --help        Show this screen.
  --version        Show version information and exit.
  --dir <dir>      xidmap path.
  --url <url>      DGraph url to connect.
`;

var Config struct {
  Predicate []string
  Dir string
  Url string
};

func main() {
  config, _ := docopt.ParseArgs(usage, nil, "1.0")
  config.Bind(&Config)
}
