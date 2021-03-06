package ui

import (
	"fmt"

	"github.com/gravitational/robotest/e2e/framework"
	"github.com/gravitational/robotest/e2e/model/ui/defaults"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	web "github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

func DeleteSite(page *web.Page, domainName string) {
	deploymentAvailable := func() bool {
		return getDeploymentIndex(page, domainName) >= 0
	}

	Eventually(deploymentAvailable, defaults.DeleteTimeout).Should(BeTrue(),
		fmt.Sprintf("deployment %q should be in the deployment list", domainName))

	deploymentIndex := getDeploymentIndex(page, domainName)
	Expect(deploymentIndex).To(BeNumerically(">=", 0), "expected to find a valid deployment index")

	By("selecting a delete site item")
	SetDropdownValue2(page, fmt.Sprintf(".grv-portal-sites tr:nth-child(%v)", deploymentIndex+1), "button", "Delete...")

	By("entering AWS credentials")
	elems := page.FindByName("aws_access_key")
	count, _ := elems.Count()
	if count > 0 {
		Expect(elems).To(BeFound(), "expected to find an input field for AWS access key")
		Expect(elems.SendKeys(framework.TestContext.AWS.AccessKey)).To(Succeed(),
			"expected to input AWS access key")

		elems = page.FindByName("aws_secret_key")
		Expect(elems).To(BeFound(), "expected to find an input field for AWS secret key")
		Expect(elems.SendKeys(framework.TestContext.AWS.SecretKey)).To(Succeed(),
			"expected to input AWS secret key")
	}

	By("entering domain name")
	elems = page.FindByName("deploymentName")
	Expect(elems).To(BeFound())
	Expect(elems.SendKeys(domainName)).To(Succeed())

	By("confirming the action")
	Expect(page.Find(".grv-dialog .btn-danger").Click()).To(Succeed())

	Eventually(deploymentAvailable, defaults.DeleteTimeout).ShouldNot(BeTrue(),
		fmt.Sprintf("deployment %q should disappear from the deployment list", domainName))
}

func getDeploymentIndex(page *web.Page, domainName string) int {
	const scriptTemplate = `
            var rows = Array.prototype.slice.call(document.querySelectorAll(".grv-portal-sites .grv-table .grv-portal-sites-tag"));
            return rows.findIndex( (tag) => {
		    return (tag.innerText == "Name:%v");
            });
        `

	script := fmt.Sprintf(scriptTemplate, domainName)
	var deploymentIndex int

	Expect(page.RunScript(script, nil, &deploymentIndex)).To(Succeed())

	return deploymentIndex
}
