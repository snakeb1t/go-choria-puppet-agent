package agent

import (
	"context"
	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-choria/plugin"
	"github.com/choria-io/go-choria/server"
	"github.com/choria-io/go-choria/server/agents"
	"github.com/choria-io/mcorpc-agent-provider/mcorpc"
)

type ResourceRequest struct {

}

type ResourceResponse struct {

}

type DisableRequest struct {

}

type DisableResponse struct {

}

type EnableRequest struct {

}

type EnableResponse struct {

}

type LastRunRequest struct {

}

type LastRunResponse struct {

}

type StatusRequest struct {

}

type StatusResponse struct {

}

type RunonceRequest struct {
	Force bool `json:"force"`
	Server string `json:"server"`
	Tags   string  `json:"tags"`
	Noop   bool    `json:"noop"`
	Splay  bool    `json:"splay"`
	SplayLimit int `json:"splaylimit"`
	Environment string `json:"environment"`
	UseCachedCatalog bool `json:"use_cached_catalog"`
}

type RunonceResponse struct {
	Summary string `json:"summary"`
	InitiatedAt string `json:"initiated_at"`
}

var metadata = &agents.Metadata{
	Name:        "puppet",
	Description: "Choria Puppet Agent",
	Author:      "Palantir Technologies <palantir.com>",
	Version:     "0.0.1",
	License:     "Apache-2",
	Timeout:     3600,
}

func New(mgr server.AgentManager) (agents.Agent, error) {
	agent := mcorpc.New("puppet", metadata, mgr.Choria(), mgr.Logger())

	agent.MustRegisterAction("runonce", runOnceAction)

	return agents.Agent(agent), nil
}

func runOnceAction(ctx context.Context, req *mcorpc.Request, reply *mcorpc.Reply, agent *mcorpc.Agent, conn choria.ConnectorInfo) {
	i := &RunonceRequest{}
	if !mcorpc.ParseRequestData(i, req, reply) {
		return
	}
	resp, err := runAgentOnce(i)
	if err != nil {
		reply.Statuscode = mcorpc.Aborted
		reply.Statusmsg = err.Error()
		return
	}
	reply.Data = resp
}

// ChoriaPlugin produces the Choria pluggable plugin it uses the metadata
// to dynamically answer questions of name and version
func ChoriaPlugin() plugin.Pluggable {
	return mcorpc.NewChoriaAgentPlugin(metadata, New)
}