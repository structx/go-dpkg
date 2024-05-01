// Package chain defined service interfaces
package chain

// interfaces in this package are used to define
// required service functionality to be used by platforms
// above the decentralized network
//
// will define as empty to begin

// AccessControl access control required functionality
//
//go:generate mockery --name AccessControl
type AccessControl interface{}

// DDNS decentralized dns required functionality
//
//go:generate mockery --name DDNS
type DDNS interface{}

// LightNode light node required functionality
//
//go:generate mockery --name LightNode
type LightNode interface{}

// TxManager transaction management required functionality
//
//go:generate mockery --name TxManager
type TxManager interface{}
