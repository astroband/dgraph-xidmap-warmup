package main;

import (
  "os"
  "log"

  "google.golang.org/grpc"
  "github.com/dgraph-io/dgo"
  "github.com/dgraph-io/dgo/protos/api"

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

var Client *dgo.Dgraph;

func main() {
  var conn *grpc.ClientConn;

  config, _ := docopt.ParseArgs(usage, nil, "1.0")
  config.Bind(&Config)

  Client, conn = connect()
  defer conn.Close()

  os.MkdirAll(Config.Dir, os.ModePerm);
}

func connect() (*dgo.Dgraph, *grpc.ClientConn) {
  conn, err := grpc.Dial(Config.Url, grpc.WithInsecure())
  if err != nil {
    log.Fatal(err)
  }

  dc := api.NewDgraphClient(conn)
  return dgo.NewDgraphClient(dc), conn
}
