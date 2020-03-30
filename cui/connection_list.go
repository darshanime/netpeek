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
		viewConnsAddLine(v,
			netflow.Src().String()+":"+tcpflow.Src().String(),
			netflow.Dst().String()+":"+tcpflow.Dst().String())
		return nil
	})
}

func getConnectionListView(tcp *layers.TCP) string {
	return tcp.TransportFlow().String()
}

func actionViewConnectionsUp(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorUp(g, v, 2)
	return nil
}

func actionViewConnectionsDown(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorDown(g, v, false)
	return nil
}
