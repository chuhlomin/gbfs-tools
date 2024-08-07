package radix

import (
	"bufio"
	"bytes"
	"container/list"
	"context"
	"fmt"
	"net"

	"github.com/mediocregopher/radix/v4/internal/proc"
	"github.com/mediocregopher/radix/v4/resp"
	"github.com/mediocregopher/radix/v4/resp/resp3"
)

type stubUnmarshaler struct {
	unmarshalInto interface{}
	errCh         chan error
}

type stubCmdUnmarshaler struct {
	ctx         context.Context
	cmds        [][]string
	unmarshaler *stubUnmarshaler
}

type stub struct {
	proc          *proc.Proc
	network, addr string
	fn            func(context.Context, []string) interface{}
	ch            chan stubCmdUnmarshaler
}

// NewStubConn returns a (fake) Conn which pretends it is a Conn to a real redis
// instance, but is instead using the given callback to service requests. It is
// primarily useful for writing tests.
//
// When EncodeDecode is called the value to be marshaled is converted into a
// []string and passed to the callback. The return from the callback is then
// marshaled into an internal buffer. The value to be decoded is unmarshaled
// into using the internal buffer. If the internal buffer is empty at
// this step then the call will block.
//
// remoteNetwork and remoteAddr can be empty, but if given will be used as the
// return from the Addr method.
//
func NewStubConn(remoteNetwork, remoteAddr string, fn func(context.Context, []string) interface{}) Conn {
	s := &stub{
		proc:    proc.New(),
		network: remoteNetwork, addr: remoteAddr,
		fn: fn,
		ch: make(chan stubCmdUnmarshaler, 128),
	}
	s.proc.Run(s.responder)
	return s
}

func (s *stub) responder(ctx context.Context) {
	doneCh := ctx.Done()
	opts := resp.NewOpts()

	retBuf := new(bytes.Buffer)
	retBr := bufio.NewReader(retBuf)
	errList, unmarshalerList := list.New(), list.New()
	popFront := func(l *list.List) interface{} {
		e := l.Front()
		l.Remove(e)
		return e.Value
	}

	asErr := func(i interface{}) error {
		err, ok := i.(error)
		if !ok {
			return nil
		} else if _, ok = err.(resp.Marshaler); ok {
			return nil
		}
		return err
	}

	for {
		select {
		case <-doneCh:
			return
		case cu := <-s.ch:
			for _, cmd := range cu.cmds {
				ret := s.fn(cu.ctx, cmd)
				if err := asErr(ret); err != nil {
					errList.PushBack(err)
				} else if err := resp3.Marshal(retBuf, ret, opts); err != nil {
					panic(fmt.Sprintf("return from stub callback could not be marshaled: %v", err))
				}
			}
			if cu.unmarshaler != nil {
				unmarshalerList.PushBack(cu.unmarshaler)
			}
			for {
				if (retBuf.Len() == 0 && retBr.Buffered() == 0 && errList.Len() == 0) ||
					unmarshalerList.Len() == 0 {
					break
				}

				unmarshaler := popFront(unmarshalerList).(*stubUnmarshaler)

				if errList.Len() > 0 {
					err := popFront(errList).(error)
					unmarshaler.errCh <- err
					continue
				}

				err := resp3.Unmarshal(retBr, unmarshaler.unmarshalInto, opts)
				unmarshaler.errCh <- err
			}
		}
	}
}

type stubEncDecCtxKey int

const stubEncDecCtxKeyQueuedCh stubEncDecCtxKey = 0

func (s *stub) EncodeDecode(ctx context.Context, m, u interface{}) error {
	opts := resp.NewOpts()
	cu := stubCmdUnmarshaler{ctx: ctx}

	if m != nil {
		buf := new(bytes.Buffer)
		br := bufio.NewReader(buf)
		if err := resp3.Marshal(buf, m, opts); err != nil {
			return err
		}

		for buf.Len() > 0 || br.Buffered() > 0 {
			var cmd []string
			if err := resp3.Unmarshal(br, &cmd, opts); err != nil {
				panic(fmt.Sprintf("could not convert resp.Marshaler to []string: %v", err))
			}
			cu.cmds = append(cu.cmds, cmd)
		}
	}

	var errCh chan error
	if u != nil {
		errCh = make(chan error, 1)
		cu.unmarshaler = &stubUnmarshaler{
			unmarshalInto: u,
			errCh:         errCh,
		}
	}

	closedCh := s.proc.ClosedCh()

	select {
	case s.ch <- cu:
	case <-ctx.Done():
		return ctx.Err()
	case <-closedCh:
		return proc.ErrClosed
	}

	// This is a hack, but it lets the Pool tests have deterministic behavior
	// when it comes to performing concurrent commands against a single conn.
	if ch, _ := ctx.Value(stubEncDecCtxKeyQueuedCh).(chan struct{}); ch != nil {
		close(ch)
	}

	if errCh == nil {
		return nil
	}

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-closedCh:
		return proc.ErrClosed
	}
}

func (s *stub) Do(ctx context.Context, a Action) error {
	return a.Perform(ctx, s)
}

func (s *stub) Addr() net.Addr {
	return rawAddr{network: s.network, addr: s.addr}
}

func (s *stub) Close() error {
	return s.proc.Close(nil)
}
