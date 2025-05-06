package tcp

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/AdhityaRamadhanus/zdb"
	"github.com/AdhityaRamadhanus/zdb/commands"
	"github.com/AdhityaRamadhanus/zdb/miniresp3"
	"github.com/pkg/errors"
)

var serverInfo = map[string]interface{}{
	"server":  "zdb",
	"proto":   3,
	"version": "0.0.1",
}

type dataCmd struct {
	name string
	args commands.CmdArgs
}

type eventCmd struct {
	cmd    []dataCmd
	writer *miniresp3.Writer
}

type Server struct {
	avlab     zdb.ZDB
	proto     string
	addr      string
	clients   int
	eventChan chan *eventCmd
}

func NewServer(proto, addr string) *Server {
	return &Server{
		avlab:     *zdb.NewZDB(16),
		proto:     proto,
		addr:      addr,
		eventChan: make(chan *eventCmd, 1000),
	}
}

func (srv *Server) eventLoop() {
	for ev := range srv.eventChan {
		for _, evcmd := range ev.cmd {
			//TODO: Maybe change to function map if it doesn't affect performance too much
			switch evcmd.name {
			case "hello":
				ev.writer.AppendMap(serverInfo)
			case "echo":
				ev.writer.AppendBulkStr(evcmd.args[0])
			case "ping":
				ev.writer.AppendSimpleStr("OK")
			case "shards":
				lengths := srv.avlab.ShardStats()
				ev.writer.AppendArrInt(lengths)
			case "scan":
				cmd := &commands.ScanCmd{}
				if err := cmd.Build(evcmd.args); err != nil {
					ev.writer.AppendSimpleError(err.Error())
					continue
				}
				keys, nextCursor, err := srv.avlab.Scan(cmd)
				if err != nil {
					ev.writer.AppendSimpleError(err.Error())
					continue
				}
				ev.writer.AppendArrAny([]interface{}{nextCursor, keys})
			case "zadd":
				cmd := &commands.ZADDCmd{}
				if err := cmd.Build(evcmd.args); err != nil {
					ev.writer.AppendSimpleError(err.Error())
					continue
				}
				success := srv.avlab.ZAdd(cmd)
				ev.writer.AppendInt(success)
			case "zcard":
				cmd := &commands.ZCardCmd{}
				if err := cmd.Build(evcmd.args); err != nil {
					ev.writer.AppendSimpleError(err.Error())
					continue
				}
				success := srv.avlab.ZCard(cmd)
				ev.writer.AppendInt(success)
			case "zcount":
				cmd := &commands.ZCountCmd{}
				if err := cmd.Build(evcmd.args); err != nil {
					ev.writer.AppendSimpleError(err.Error())
					continue
				}
				count := srv.avlab.ZCount(cmd)
				ev.writer.AppendInt(count)
			case "zrange":
				cmd := &commands.ZRangeCmd{}
				if err := cmd.Build(evcmd.args); err != nil {
					ev.writer.AppendSimpleError(err.Error())
					continue
				}
				nodes := srv.avlab.ZRange(cmd)
				if cmd.WithScores {
					ev.writer.AppendArrStr(zdb.Reduce(nodes, func(acc []string, n zdb.Node) []string {
						acc = append(acc, n.Key())
						acc = append(acc, fmt.Sprintf("%.2f", n.Score()))
						return acc
					}, []string{}))
				} else {
					ev.writer.AppendArrStr(zdb.Map(nodes, func(n zdb.Node) string {
						return n.Key()
					}))
				}
			case "zrank":
				cmd := &commands.ZRankCmd{}
				if err := cmd.Build(evcmd.args); err != nil {
					ev.writer.AppendSimpleError(err.Error())
					continue
				}
				success := srv.avlab.ZRank(cmd)
				ev.writer.AppendInt(success)
			case "zrem":
				cmd := &commands.ZRemCmd{}
				if err := cmd.Build(evcmd.args); err != nil {
					ev.writer.AppendSimpleError(err.Error())
					continue
				}
				success := srv.avlab.ZRem(cmd)
				ev.writer.AppendInt(success)
			case "zscan":
				cmd := &commands.ZScanCmd{}
				if err := cmd.Build(evcmd.args); err != nil {
					ev.writer.AppendSimpleError(err.Error())
					continue
				}
				keys, nextCursor := srv.avlab.ZScan(cmd)
				ev.writer.AppendArrAny([]interface{}{nextCursor, keys})
			case "zscore":
				cmd := &commands.ZScoreCmd{}
				if err := cmd.Build(evcmd.args); err != nil {
					ev.writer.AppendSimpleError(err.Error())
					continue
				}

				score, err := srv.avlab.ZScore(cmd)
				if err != nil {
					ev.writer.AppendNil()
					continue
				}

				ev.writer.AppendFloat64(score)
			default:
				ev.writer.AppendSimpleStr("OK")
			}
		}

		ev.writer.Write()
	}
}

func (srv *Server) Run() {
	l, err := net.Listen(srv.proto, srv.addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	go srv.eventLoop()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("error in accepting conn", err)
			continue
		}

		go srv.handleClient(conn)
	}
}

func (srv *Server) handleClient(conn net.Conn) (err error) {
	defer func() {
		if err != nil && err.Error() != "EOF" {
			log.Println(err)
		} else {
			log.Println("Closing Connection")
		}
	}()

	errChan := make(chan error, 10)
	go srv.handleData(conn, errChan)
	for err := range errChan {
		return err
	}

	return nil
}

func (srv *Server) handleData(conn net.Conn, errChan chan<- error) (err error) {
	defer func() {
		if err != nil {
			errChan <- err
		}
	}()

	r := miniresp3.NewReader(conn)
	w := miniresp3.NewWriter(conn)
	cmds := []dataCmd{}
	for {
		arrayCount, err := r.ReadArrayHeader()
		if err != nil {
			return err
		}

		args := []string{}
		for i := 0; i <= int(arrayCount-1); i++ {
			bulkStr, err := r.ReadBulkString()
			if err != nil {
				return errors.Wrap(err, "failed to parse Bulk String")
			}
			args = append(args, bulkStr)
		}

		cmds = append(cmds, dataCmd{
			name: strings.ToLower(args[0]),
			args: args[1:],
		})

		if r.IsAllRead() {
			srv.eventChan <- &eventCmd{
				cmd:    cmds,
				writer: w,
			}
			cmds = []dataCmd{}
		}
	}
}
