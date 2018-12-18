package memory

import (
	"sync"
	"time"

	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
)

////////////////////////////////////////////////////////////////////////////////
// registerHostRequestTable
type registerHostRequestTable struct {
	m *sync.RWMutex
	// primary key is the HostID represented as a string
	data map[string]proto_rexecd.RegisterHostRequest
}

type registerHostRequestRow struct {
	key string
	val proto_rexecd.RegisterHostRequest
}

func newRegisterHostRequestTable() *registerHostRequestTable {
	return &registerHostRequestTable{
		m:    &sync.RWMutex{},
		data: map[string]proto_rexecd.RegisterHostRequest{},
	}
}

func (r *registerHostRequestTable) get(hostID string) (proto_rexecd.RegisterHostRequest, bool) {
	r.m.RLock()
	defer r.m.RUnlock()
	registerUserRequest, found := r.data[hostID]
	return registerUserRequest, found
}

func (r *registerHostRequestTable) set(registerHostRequest proto_rexecd.RegisterHostRequest) {
	r.m.Lock()
	defer r.m.Unlock()
	r.data[registerHostRequest.GetHostId()] = registerHostRequest
}

func (r *registerHostRequestTable) each() chan registerHostRequestRow {
	ch := make(chan registerHostRequestRow, 1)
	go func() {
		for key, val := range r.data {
			r.m.RLock()
			ch <- registerHostRequestRow{key: key, val: val}
			r.m.RUnlock()
		}
		close(ch)
	}()
	return ch
}

////////////////////////////////////////////////////////////////////////////////
// registerUserRequestTable
type registerUserRequestTable struct {
	m    *sync.RWMutex
	data map[string]proto_rexecd.RegisterUserRequest
}

type registerUserRequestRow struct {
	key string
	val proto_rexecd.RegisterUserRequest
}

func newRegisterUserRequestTable() *registerUserRequestTable {
	return &registerUserRequestTable{
		m:    &sync.RWMutex{},
		data: map[string]proto_rexecd.RegisterUserRequest{},
	}
}

func (r *registerUserRequestTable) get(username string) (proto_rexecd.RegisterUserRequest, bool) {
	r.m.RLock()
	defer r.m.RUnlock()
	registerUserRequest, found := r.data[username]
	return registerUserRequest, found
}

func (r *registerUserRequestTable) set(registerUserReqeust proto_rexecd.RegisterUserRequest) {
	r.m.Lock()
	defer r.m.Unlock()
	r.data[registerUserReqeust.GetUsername()] = registerUserReqeust
}

func (r *registerUserRequestTable) each() chan registerUserRequestRow {
	ch := make(chan registerUserRequestRow, 1)
	go func() {
		for key, val := range r.data {
			r.m.RLock()
			ch <- registerUserRequestRow{key: key, val: val}
			r.m.RUnlock()
		}
		close(ch)
	}()
	return ch
}

////////////////////////////////////////////////////////////////////////////////
// commandStore
type commandStore struct {
	data             map[time.Time]*proto_rexecd.CommandRequest
	userCommandStore *userCommandStore
	hostCommandStore *hostCommandStore
	m                *sync.RWMutex
}

func newCommandStore() *commandStore {
	return &commandStore{
		m:                &sync.RWMutex{},
		data:             map[time.Time]*proto_rexecd.CommandRequest{},
		userCommandStore: newUserCommandStore(),
		hostCommandStore: newHostCommandStore(),
	}
}

func (c *commandStore) addCommand(commandRequest proto_rexecd.CommandRequest, t time.Time) {
	c.userCommandStore.addCommand(commandRequest, t)
	for _, hostConfig := range commandRequest.GetHostConfig() {
		hostID := hostConfig.GetHostId()
		c.hostCommandStore.addCommand(hostID, t, commandRequest)
	}
}

func (c *commandStore) execData(hostID string, t time.Time) (*execData, bool) {
	_, found := c.hostCommandStore.data[hostID]
	if !found {
		return newExecData(hostID, t, proto_rexecd.CommandRequest{}), false
	}
	execData, found := c.hostCommandStore.data[hostID][t]
	return execData, found
}

////////////////////////////////////////////////////////////////////////////////
// userCommandStore
type userCommandStore struct {
	m    *sync.RWMutex
	data map[string]map[time.Time]proto_rexecd.CommandRequest
}

func newUserCommandStore() *userCommandStore {
	return &userCommandStore{
		m:    &sync.RWMutex{},
		data: map[string]map[time.Time]proto_rexecd.CommandRequest{},
	}
}

func (u *userCommandStore) addCommand(commandRequest proto_rexecd.CommandRequest, t time.Time) {
	u.m.Lock()
	defer u.m.Unlock()
	username := commandRequest.GetUsername()
	_, found := u.data[username]
	if !found {
		u.data[username] = map[time.Time]proto_rexecd.CommandRequest{}
	}
	u.data[username][t] = commandRequest
}

////////////////////////////////////////////////////////////////////////////////
// hostCommandStore
type hostCommandStore struct {
	m    *sync.RWMutex
	data map[string]map[time.Time]*execData
}

func newHostCommandStore() *hostCommandStore {
	return &hostCommandStore{
		m:    &sync.RWMutex{},
		data: map[string]map[time.Time]*execData{},
	}
}

func (h *hostCommandStore) addCommand(hostID string, t time.Time, commandRequest proto_rexecd.CommandRequest) {
	h.m.Lock()
	defer h.m.Unlock()
	for _, hostConfig := range commandRequest.GetHostConfig() {
		_, found := h.data[hostID]
		if !found {
			h.data[hostID] = map[time.Time]*execData{}
		}
		h.data[hostID][t] = newExecData(hostConfig.GetHostId(), t, commandRequest)
	}
}

////////////////////////////////////////////////////////////////////////////////
// execData
type execData struct {
	m              *sync.RWMutex
	hostID         string
	time           time.Time
	commandRequest proto_rexecd.CommandRequest
	stdoutLines    []*proto_rexecd.Line
	stderrLines    []*proto_rexecd.Line
	done           bool
}

func newExecData(hostID string, t time.Time, commandRequest proto_rexecd.CommandRequest) *execData {
	return &execData{
		m:              &sync.RWMutex{},
		hostID:         hostID,
		time:           t,
		commandRequest: commandRequest,
		stdoutLines:    []*proto_rexecd.Line{},
		stderrLines:    []*proto_rexecd.Line{},
	}
}

func (e *execData) addStdoutLine(b []byte) []byte {
	e.m.Lock()
	defer e.m.Unlock()
	e.stdoutLines = append(e.stdoutLines, &proto_rexecd.Line{
		Line: b,
		Time: int64(time.Now().Nanosecond()),
	})
	return b
}

func (e *execData) addStderrLine(b []byte) []byte {
	e.m.Lock()
	defer e.m.Unlock()
	e.stderrLines = append(e.stderrLines, &proto_rexecd.Line{
		Line: b,
		Time: int64(time.Now().Nanosecond()),
	})
	return b
}

func (e *execData) setDone() {
	e.m.Lock()
	defer e.m.Unlock()
	e.done = true
}

func (e *execData) eachStdoutLine() chan proto_rexecd.Line {
	ch := make(chan proto_rexecd.Line, 1)
	go func() {
		for _, line := range e.stdoutLines {
			e.m.RLock()
			ch <- *line
			e.m.RUnlock()
		}
		close(ch)
	}()
	return ch
}

func (e *execData) eachStderrLine() chan proto_rexecd.Line {
	ch := make(chan proto_rexecd.Line, 1)
	go func() {
		for _, line := range e.stderrLines {
			e.m.RLock()
			ch <- *line
			e.m.RUnlock()
		}
		close(ch)
	}()
	return ch
}
