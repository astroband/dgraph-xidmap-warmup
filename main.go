package main;

import (
  "os"
  "log"
  "context"
  "encoding/json"
  "encoding/binary"

  "github.com/docopt/docopt-go"

  "google.golang.org/grpc"
  "github.com/dgraph-io/dgo"
  "github.com/dgraph-io/dgo/protos/api"

  "github.com/dgraph-io/badger"

  "github.com/bitherhq/go-bither/common/hexutil"
)

var usage = `dgraph live cache warmup utility.

Usage:
  dgraph-xidmap-warmup PREDICATE --dir <dir> --url <url>

Arguments:
  PREDICATE predicate holding hash keys

Options:
  -h --help        Show this screen.
  --version        Show version information and exit.
  --dir <dir>      xidmap path.
  --url <url>      DGraph url to connect.
`;

// TODO: multiple predicates for low/high
var Config struct {
  Predicate string
  Dir string
  Url string
};

var Client *dgo.Dgraph;
var DB *badger.DB;

func main() {
  var conn *grpc.ClientConn;

  log.SetPrefix("xidmap ")

  config, _ := docopt.ParseArgs(usage, nil, "1.0")
  config.Bind(&Config)

  Client, conn = connect()
  defer conn.Close()

  DB = open_db()
  defer DB.Close()

  os.MkdirAll(Config.Dir, os.ModePerm);

  result := query()
  store(result)
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

func query() ([]interface{}) {
  log.Println("Quering DGraph...")

  q := `
    query {
      all(func: has(` + Config.Predicate + `)) @cascade {
        uid
        ` + Config.Predicate + `
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

  return data["all"].([]interface{});
}

func store(r []interface{}) {
  log.Printf("Storing %d records...", len(r))

  for _, item := range r {
    i := item.(map[string]interface{})


    err := DB.Update(func(txn *badger.Txn) error {
      s := i["uid"].(string)

      uid, err := hexutil.DecodeUint64(s)
      if (err != nil) {
        return err;
      }

      buf := make([]byte, binary.MaxVarintLen64)
      n := binary.PutUvarint(buf, uid)

      err = txn.Set([]byte(i[Config.Predicate].(string)), buf[:n])
      return err;
    })

    if (err != nil) {
      log.Fatal(err)
    }
  }
}
