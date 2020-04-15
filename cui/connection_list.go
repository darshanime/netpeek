package cui

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/jroimartin/gocui"
)

func AddConnection(netflow, tcpflow gopacket.Flow) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("conns")
		if err != nil {
			return err
		}
		if _, ok := connMap[netflow]; !ok {
			connMap[netflow] = 1
		}
		addLineToViewConns(v,
			netflow.Src().String()+":"+tcpflow.Src().String(),
			netflow.Dst().String()+":"+tcpflow.Dst().String())
		return nil
	})
}

func getConnectionListView(tcp *layers.TCP) string {
	return tcp.TransportFlow().String()
}
