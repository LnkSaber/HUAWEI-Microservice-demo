package highway

import (
	"context"
	"sync"

	"github.com/go-chassis/go-chassis/core/client"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/invocation"
)

//const timeout
const (
	//Name is a variable of type string
	Name                  = "highway"
	DefaultConnectTimeOut = 60
	DefaultSendTimeOut    = 300
)

//highwayClient
//Deprecated
type highwayClient struct {
	once     sync.Once
	opts     client.Options
	reqMutex sync.Mutex // protects following
}

//NewHighwayClient is a function
//Deprecated
func NewHighwayClient(options client.Options) (client.ProtocolClient, error) {

	rc := &highwayClient{
		once: sync.Once{},
		opts: options,
	}

	c := client.ProtocolClient(rc)

	return c, nil
}

func (c *highwayClient) String() string {
	return "highway_client"
}
func (c *highwayClient) Close() error {
	return nil
}

// GetOptions method return opts
func (c *highwayClient) GetOptions() client.Options {
	return c.opts
}

// ReloadConfigs
func (c *highwayClient) ReloadConfigs(opts client.Options) {
	c.opts = opts
}
func invocation2Req(inv *invocation.Invocation) *Request {
	highwayReq := &Request{}
	highwayReq.MsgID = uint64(int(GenerateMsgID()))
	highwayReq.MethodName = inv.OperationID
	highwayReq.Schema = inv.SchemaID
	highwayReq.Arg = inv.Args
	highwayReq.SvcName = inv.MicroServiceName
	return highwayReq
}
func (c *highwayClient) Call(ctx context.Context, addr string, inv *invocation.Invocation, rsp interface{}) error {
	connParams := &ConnParams{}
	connParams.TLSConfig = c.opts.TLSConfig
	connParams.Addr = addr
	connParams.Timeout = DefaultConnectTimeOut
	baseClient, err := CachedClients.GetClient(connParams)
	if err != nil {
		return err
	}
	tmpRsp := &Response{0, Ok, "", 0, rsp, nil}
	highwayReq := invocation2Req(inv)
	//Current only twoway
	highwayReq.TwoWay = true
	highwayReq.Attachments = common.FromContext(ctx)

	err = baseClient.Send(highwayReq, tmpRsp, DefaultSendTimeOut)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	client.InstallPlugin(Name, NewHighwayClient)

}
