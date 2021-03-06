package agent

import (
	"fmt"

	utils "github.com/gravitational/robotest/e2e/model/ui"
	"github.com/gravitational/robotest/infra"

	. "github.com/onsi/gomega"
	web "github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

type AgentServer struct {
	Hostname string
	index    int
	page     *web.Page
}

func CreateAgentServer(page *web.Page, index int) AgentServer {
	cssSelector := getServerCssSelector(index)
	elem := page.Find(cssSelector)
	Expect(elem).To(BeFound())

	hostname, _ := elem.Find(".grv-provision-req-server-hostname span").Text()
	Expect(hostname).NotTo(BeEmpty(), "should have a hostname")

	return AgentServer{page: page, Hostname: hostname, index: index}
}

func (a *AgentServer) SetIP(value string) {
	cssSelector := fmt.Sprintf("%v .grv-provision-req-server-interface", getServerCssSelector(a.index))

	utils.SetDropdownValue2(a.page, cssSelector, "", value)
}

func (a *AgentServer) SetDockerDevice(value string) {
	cssSelector := fmt.Sprintf(`%v input[placeholder="loopback"]`, getServerCssSelector(a.index))
	Expect(a.page.Find(cssSelector).Fill(value)).To(
		Succeed(),
		"should set a docker device value")
}

func (a *AgentServer) GetIPs() []string {
	const scriptTemplate = `
            var result = [];
            var cssSelector = "%v .grv-provision-req-server-interface li a";
            var children = document.querySelectorAll(cssSelector);
            children.forEach( z => result.push(z.text) );
            return result; `
	var result []string

	script := fmt.Sprintf(scriptTemplate, getServerCssSelector(a.index))
	a.page.RunScript(script, nil, &result)
	return result
}

func (a *AgentServer) SetIPByInfra(provisioner infra.Provisioner) {
	ips := a.GetIPs()
	if len(ips) < 2 {
		return
	}
	var node infra.Node
	for _, ip := range ips {
		node, _ = provisioner.NodePool().Node(ip)
		if node != nil {
			break
		}
	}

	descriptionText := fmt.Sprintf("cannot find node matching any of %v IPs", a.Hostname)
	Expect(node).NotTo(BeNil(), descriptionText)
	a.SetIP(node.Addr())
}

func getServerCssSelector(index int) string {
	return fmt.Sprintf(".grv-provision-req-server:nth-child(%v)", index+1)
}
