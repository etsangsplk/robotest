package installer

import (
	"fmt"

	. "github.com/onsi/gomega"
	web "github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

type OnPremServer struct {
	Hostname string
	index    int
	page     *web.Page
}

func (p *OnPremServer) GetIPs() []string {
	const scriptTemplate = `
		var result = [];
		var cssSelector = "%v .grv-provision-req-server-interface li a";
		var children = document.querySelectorAll(cssSelector);
		children.forEach( z => result.push(z.text) );
		return result; `
	var result []string

	script := fmt.Sprintf(scriptTemplate, getServerCssSelector(p.index))
	p.page.RunScript(script, nil, &result)
	return result
}

func createServer(page *web.Page, index int) OnPremServer {
	cssSelector := getServerCssSelector(index)
	elem := page.Find(cssSelector)
	Expect(elem).To(BeFound())

	hostname, _ := elem.Find(".grv-provision-req-server-hostname span").Text()
	Expect(hostname).NotTo(BeEmpty())

	return OnPremServer{page: page, Hostname: hostname, index: index}

}

func getServerCssSelector(index int) string {
	return fmt.Sprintf(".grv-provision-req-server:nth-child(%v)", index+1)
}
