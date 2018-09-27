package flexy

type Config struct {
	Groups []Group
	OCVars map[string]string
}

type Group struct {
	Name string
	Hosts []Host
}

type Host struct {
	ID string
	PublicDNS string
	OCNodeGroupName string
	OCSchedulable bool
	IPv4PublicIP string
}