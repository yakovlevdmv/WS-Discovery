package WS_Discovery

import (
	"fmt"
	"os"
	"github.com/beevik/etree"
	"github.com/satori/go.uuid"
	"strconv"
)

const (
	WS_namespace = "http://docs.oasis-open.org/ws-dd/ns/discovery/2009/01"
	multicast_ip = "239.255.255.250"
	multicast_port = 3702
)

var deviceTypes = map[DeviceType]string{
	NVD : "NetworkVideoDisplay",
	NVS : "NetworkVideoStorage",
	NVA : "NetworkVideoAnalytics",
	NVT : "NetworkVideoTransmitter",
}

type DeviceType int

const (
	NVD DeviceType = iota
	NVS
	NVA
	NVT
)

func (devType DeviceType) String() string {
	stringRepresentation := []string {
		"NVD : Network Video Display",
		"NVS : Network Video Storage",
		"NVA : Network Video Analytics",
		"NVT : Network Video Transmitter",
	}
	i := uint8(devType)
	switch {
	case i <= uint8(NVT):
		return stringRepresentation[i]
	default:
		return strconv.Itoa(int(i))
	}
}

func BuildProbeMessage(uuidV4 string, scopes, types []string) string {
	//Список namespace
	namespaces := make(map[string]string)
	namespaces["a"] = "http://www.w3.org/2005/08/addressing"
	namespaces["d"] = "http://schemas.xmlsoap.org/ws/2005/04/discovery"
	namespaces["dn"] = "http://www.onvif.org/ver10/network/wsdl"

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

	//Содержимое Body
	var bodyContent []*etree.Element

	probe := etree.NewElement("d:Probe")
	typesTag := etree.NewElement("d:Types")
	typesTag.SetText("dn:" + deviceTypes[NVT])
	probe.AddChild(typesTag)
	scopesTag := etree.NewElement("d:Scopes")
	probe.AddChild(scopesTag)
	bodyContent = append(bodyContent, probe)

	return BuildSoapMessage(headerContent, bodyContent, namespaces)
}

func SendProbe(interfaceName string) []string{
	// Creating UUID Version 4
	uuidV4 := uuid.Must(uuid.NewV4())
	fmt.Printf("UUIDv4: %s\n", uuidV4)

	types := []string{NVT.String()}

	probeSOAP := BuildProbeMessage(uuidV4.String(), nil, types)

	fmt.Println(probeSOAP)

	return sendUDPMulticast(probeSOAP, interfaceName)

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}