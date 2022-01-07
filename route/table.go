package route

type Route struct {
	Cidr    string
	Name    string
	Address string
}

type Table struct {
	Name   string
	Subnet []string
	Routes map[string]Route
}

// NewTable 创建新的路由表格
func NewTable(name string, subnet []string) Table {
	return Table{
		Name:   name,
		Subnet: subnet,
		Routes: make(map[string]Route),
	}
}

func (table Table) Find(address string) *Route {
	route := table.Routes[address]
	return &route
}
