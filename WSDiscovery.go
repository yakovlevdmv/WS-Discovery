package WS_Discovery

import (
	"github.com/beevik/etree"
	"github.com/yakovlevdmv/gosoap"
	"fmt"
	"strings"
)

func BuildProbeMessage(uuidV4 string, scopes, types []string, nmsp map[string]string) gosoap.SoapMessage {
	//Список namespace
	namespaces := make(map[string]string)
	namespaces["a"] = "http://www.w3.org/2005/08/addressing"
	namespaces["d"] = "http://schemas.xmlsoap.org/ws/2005/04/discovery"

	probeMessage := gosoap.NewEmptySOAP()
	
	probeMessage.AddRootNamespaces(namespaces)
	if len(nmsp) != 0 {
		probeMessage.AddRootNamespaces(nmsp)
	}

	fmt.Println(probeMessage.String())

	//Содержимое Head
	var headerContent []*etree.Element

	action := etree.NewElement("a:Action")
	action.SetText("http://docs.oasis-open.org/ws-dd/ns/discovery/2009/01/Probe")
	action.CreateAttr("mustUnderstand", "1")

	msgID := etree.NewElement("a:MessageID")
	msgID.SetText("uuid:" + uuidV4)

	replyTo := etree.NewElement("a:ReplyTo")
	replyTo.CreateElement("a:Address").SetText("http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous")

	to := etree.NewElement("a:To")
	to.SetText("urn:schemas-xmlsoap-org:ws:2005:04:discovery")
	to.CreateAttr("mustUnderstand", "1")

	headerContent = append(headerContent, action, msgID, replyTo, to)
	probeMessage.AddHeaderContents(headerContent)

	//Содержимое Body
	probe := etree.NewElement("d:Probe")

	typesTag := etree.NewElement("d:Types")
	var typesString string
	for _, j := range types {
		typesString += j
		typesString += " "
	}

	typesTag.SetText(strings.TrimSpace(typesString))

	scopesTag := etree.NewElement("d:Scopes")
	var scopesString string
	for _, j := range scopes {
		scopesString += j
		scopesString += " "
	}
	scopesTag.SetText(strings.TrimSpace(scopesString))

	probe.AddChild(typesTag)
	probe.AddChild(scopesTag)

	probeMessage.AddBodyContent(probe)

	return probeMessage
}