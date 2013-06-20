package main

//------------------------------------------------------------------------------
//
// Commands
//
//------------------------------------------------------------------------------

import (
	"encoding/json"
	"github.com/benbjohnson/go-raft"
	"github.com/xiangli-cmu/raft-etcd/store"
	"time"
)

// A command represents an action to be taken on the replicated state machine.
type Command interface {
	CommandName() string
	Apply(server *raft.Server) ([]byte, error)
}

// Set command
type SetCommand struct {
	Key        string    `json:"key"`
	Value      string    `json:"value"`
	ExpireTime time.Time `json:"expireTime"`
}

// The name of the command in the log
func (c *SetCommand) CommandName() string {
	return "set"
}

// Set the value of key to value
func (c *SetCommand) Apply(server *raft.Server) ([]byte, error) {
	return store.Set(c.Key, c.Value, c.ExpireTime)
}

// Get the path for http request
func (c *SetCommand) GeneratePath() string {
	return "set/" + c.Key
}

// Get command
type GetCommand struct {
	Key string `json:"key"`
}

// The name of the command in the log
func (c *GetCommand) CommandName() string {
	return "get"
}

// Set the value of key to value
func (c *GetCommand) Apply(server *raft.Server) ([]byte, error) {
	res := store.Get(c.Key)
	return json.Marshal(res)
}

func (c *GetCommand) GeneratePath() string {
	return "get/" + c.Key
}

// Delete command
type DeleteCommand struct {
	Key string `json:"key"`
}

// The name of the command in the log
func (c *DeleteCommand) CommandName() string {
	return "delete"
}

// Delete the key
func (c *DeleteCommand) Apply(server *raft.Server) ([]byte, error) {
	return store.Delete(c.Key)
}

// Watch command
type WatchCommand struct {
	Key string `json:"key"`
}

//The name of the command in the log
func (c *WatchCommand) CommandName() string {
	return "watch"
}

func (c *WatchCommand) Apply(server *raft.Server) ([]byte, error) {
	ch := make(chan store.Response)

	// add to the watchers list
	store.AddWatcher(c.Key, ch)

	// wait for the notification for any changing
	res := <-ch

	return json.Marshal(res)
}

// JoinCommand
type JoinCommand struct {
	Name string `json:"name"`
}

func (c *JoinCommand) CommandName() string {
	return "join"
}

func (c *JoinCommand) Apply(server *raft.Server) ([]byte, error) {
	err := server.AddPeer(c.Name)
	// no result will be returned
	return nil, err
}