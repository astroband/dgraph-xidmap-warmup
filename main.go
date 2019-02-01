package main;

import (
  "os"
  "log"
  "context"
  // "encoding/hex"
  "encoding/json"

  "github.com/docopt/docopt-go"

  "google.golang.org/grpc"
  "github.com/dgraph-io/dgo"
  "github.com/dgraph-io/dgo/protos/api"

  "github.com/dgraph-io/badger"
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
var DB *badger.DB;

func main() {
  var conn *grpc.ClientConn;

  config, _ := docopt.ParseArgs(usage, nil, "1.0")
  config.Bind(&Config)

  Client, conn = connect()
  defer conn.Close()

  DB = open_db()
  defer DB.Close()

  os.MkdirAll(Config.Dir, os.ModePerm);

  query();
}

func connect() (*dgo.Dgraph, *grpc.ClientConn) {
  conn, err := grpc.Dial(Config.Url, grpc.WithInsecure())
  if err != nil {
    log.Println("Failed to connect to DGraph")
    log.Fatal(err)
  }

  dc := api.NewDgraphClient(conn)
  return dgo.NewDgraphClient(dc), conn
}

func open_db() (*badger.DB) {
  opts := badger.DefaultOptions
  opts.Dir = Config.Dir
  opts.ValueDir = Config.Dir

  db, err := badger.Open(opts)
  if err != nil {
    log.Println("Failed to open badger database")
	  log.Fatal(err)
  }
  return db
}

func query() {
  q := `
    query {
      all(func: has(key)) @cascade {
        uid
        key
      }
    }
  `;

  ctx := context.Background()

  txn := Client.NewTxn()
  defer txn.Discard(ctx)

  resp, err := txn.Query(ctx, q)
  if (err != nil) {
    log.Fatal(err)
  }

  var data map[string]interface{}

  err = json.Unmarshal(resp.Json, &data)
  if err != nil {
    log.Fatal(err)
  }

  //log.Println(data["all"])

  for _, item := range data["all"].([]interface{}) {
    x := item.(map[string]interface{})
    log.Println(x["key"], x["uid"])
  }
}
