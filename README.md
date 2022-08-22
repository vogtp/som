# Service Oriented Monitoring (SOM)
<!--
[![GoDoc](https://pkg.go.dev/badge/github.com/vogtp/som?status.svg)](https://pkg.go.dev/github.com/vogtp/som?tab=doc) &nbsp;
[![Go Report Card](https://goreportcard.com/badge/github.com/vogtp/som)](https://goreportcard.com/report/github.com/vogtp/som) &nbsp;
[![Release](https://img.shields.io/github/release/vogtp/som.svg?style=flat-square)](https://github.com/vogtp/som/releases)&nbsp;
[![Go](https://github.com/vogtp/som/actions/workflows/go.yml/badge.svg)](https://github.com/vogtp/som/actions/workflows/go.yml)&nbsp;
[![codecov](https://codecov.io/gh/vogtp/som/branch/main/graph/badge.svg?token=DV0IDZ2FXE)](https://codecov.io/gh/vogtp/som)&nbsp;
-->

The Service Oriented Monitoring (SOM) monitors web applications from the
perspective of an user.

Szenarios are automatically correlated over users and regions.

Monitoring is done by instrumenting a chrome browser with the 
<a href="https://chromedevtools.github.io/devtools-protocol/" target="_blank">Chrome DevTools Protocol</a>.

## Conecpts

SOM uses a monitoring as code approach.

### Szenario

A szenario is the central concept of SOM.

A szenarion is a bit of go code that monitors an application.

All szenarios with the same name are correlated.

### User

A szenario is run in the context of an user. 
The user provides a szenarion with password (and possibly more information).

Most users will have more than one szenario.

#### User Type

A user type (e.g. staf, client) assisiates users with a list of 
szenarios.

### Region

A region is the place (e.g. internal network, internet) 
where the monitoring originates from.

Regions are used to correlate the szenarios.

## Correlation

A szenarion is correlated first by user (by szenario runs of the user) 
then by region (by the users of the region).

A correlation group (i.e. user, region, szenario) has a level:
Level   | Num Value
--------|----------
Unknown | 0
OK      | 1
Issues  | 2
Warning | 3
Down    | 4

The level of a group is determined by summing the level of its children
and dividing it by the number of children and then rounded.

The following restriction apply:

1. If the last szenario run of a user was OK the level of the user is OK
1. Issues is only reached if the level is at least Issues (2)
1. Down is only reached if the level is at leat Down - 0.2 (3.8)
1. If there are no children the level is Unknown

A correlation tree looks like this:

    Szenario OWA: OK
        Region development: OK
            User som-user-dev-1: OK
                08.08.2022 19:22:49:  2.16s OK
                08.08.2022 19:17:22:  2.35s OK
                08.08.2022 19:13:41:  4.43s OK
        Region default: OK
            User som-user-dev-1: OK
                08.08.2022 13:16:04:  3.36s OK
                08.08.2022 13:15:51:  3.19s OK
    Szenario CourseSearch: OK
        Region development: OK
            User som-user-dev-1: OK
                08.08.2022 19:23:51: 22.46s OK
                08.08.2022 19:18:25: 22.37s OK
                08.08.2022 19:14:45: "CourseSearch" step "Loading" failed: page load error net::ERR_INTERNET_DISCONNECTED
            User som-world-dev-1: OK
                09.08.2022 04:06:02: 22.22s OK
                09.08.2022 04:03:37: 24.32s OK
                09.08.2022 04:01:59: "CourseSearch" step "Loading" failed: page load error net::ERR_INTERNET_DISCONNECTED
            User som-world-dev-2: Issues
                09.08.2022 04:07:50: "CourseSearch" step "Loading" failed: page load error net::ERR_INTERNET_DISCONNECTED
                09.08.2022 04:05:24: 22.26s OK
                09.08.2022 04:03:00: 22.95s OK
    Szenario PersSearch: OK
        Region development: OK
            User som-world-dev-1: Issues
                09.08.2022 04:07:24: "PersSearch" step "Loading" failed: page load error net::ERR_INTERNET_DISCONNECTED
                09.08.2022 04:05:01:  0.79s OK
                09.08.2022 04:02:36:  0.84s OK
            User som-world-dev-2: OK
                09.08.2022 04:08:27:  0.94s OK
                09.08.2022 04:06:46:  0.79s OK
                09.08.2022 04:04:23: "PersSearch" step "Loading" failed: page load error net::ERR_INTERNET_DISCONNECTED
            User som-user-dev-1: OK
                08.08.2022 19:21:48:  0.75s OK
                08.08.2022 19:16:22:  0.86s OK
                08.08.2022 19:12:40:  0.97s OK
        Region default: OK
            User som-user-dev-1: OK
                08.08.2022 13:15:36:  1.34s OK
                08.08.2022 13:15:11:  1.81s OK
                08.08.2022 12:48:14:  1.40s OK
    Szenario Intranet: Down
        Region development: Down
            User som-user-dev-1: Down
                08.08.2022 19:25:13: "Intranet" step "Loading" failed: timeout 1m0s
                08.08.2022 19:19:47: "Intranet" step "Loading" failed: timeout 1m0s
                08.08.2022 19:15:45: "Intranet" step "Loading" failed: page load error net::ERR_INTERNET_DISCONNECTED


## Writing Szenarios

To create a szenario you have to inherit `*szenario.Base`:

    type customSzenario struct {
        *szenario.Base
        Search string
    }

and implement its `Execute` method:

    func (s customSzenario) Execute(engine szenario.Engine) (err error) {
        engine.Step("Loading",
            chromedp.Navigate("https://google.ch/"),
            chromedp.WaitVisible(`#tophf`, chromedp.ByID),
        )
        engine.Step("Check",
            engine.Body(engine.Contains("Google Search"), engine.Contains("About"), engine.Bigger(1000)),
        )

        return nil
    }

`engine.Step(stepName string, actions ...chromedp.Action)` executes a step with a name (stepName) and one or more actions.
All actions runnable by https://pkg.go.dev/github.com/chromedp/chromedp@v0.8.4#Run can be used.

`chromedp.Navigate("https://google.ch/")` opens https://google.ch/

`chromedp.WaitVisible(`#tophf`, chromedp.ByID)` wait for an element with the ID #tophf to apear.

    engine.Body(engine.Contains("Google Search"), engine.Contains("About"), engine.Bigger(1000)),
    

checks if the Body conatins "Google Search" and "About" and it's size is bigger than 1000 characters.

A input can be done like this:

        engine.Step("searching",
			chromedp.SendKeys(
				`document.querySelector("body > div.L3eUgb > div.o3j99.ikrT4e.om7nvf > form > div:nth-child(1) > div.A8SBwf > div.RNNXgb > div > div.a4bIc > input")`,
				fmt.Sprintf("%s\r", s.Search),
				chromedp.ByJSPath, // copy JSPath from chrom developer tools
			),
		)

Finally the szenario has to be added to the szenarion config in `loader.go`:

For a working google example see szenario/google.go

### Sniplets

#### Click Button

The following clicks the button with the ID buttonId

    engine.Step("Click Accept", chromedp.Click("#buttonId", chromedp.ByID))

#### Check if something is visible

The following checks if a button with id buttonId is present and clicks it.

    if ok := engine.IsPresent("#buttonId", chromedp.ByID); ok {
        engine.Step("Click Accept", chromedp.Click("#buttonId", chromedp.ByID))
    }

#### Login / Enter values

		engine.Step("Login",
			chromedp.WaitVisible(`#username`, chromedp.ByID),
			chromedp.SendKeys(`#username`, s.User().Name()+"\r", chromedp.ByID),
			chromedp.WaitReady(`#password`, chromedp.ByID),
			chromedp.SendKeys(`#password`, s.User().Password()+"\r", chromedp.ByID),
		)

## Engine Methods

The engine passed to the Execute method has the following methods:

Method Signature                                              | Explaination
--------------------------------------------------------------|----------------------------------------------------------------------------- 
Step(name string, actions ...chromedp.Action)                 | Step executes the actions given and records how long it takes
IsPresent(sel interface{}, opts ...chromedp.QueryOption) bool | IsPresent checks if something is present
SetStatus(key, val string)                                    | SetStatus sets a status of the event
AddErr(err error)                                             | AddErr adds a error to the event
Body(checks ...CheckFunc) chromedp.Action                     | Body is used to check the content of the page
Contains(s string) CheckFunc                                  | Contains looks for a string in the body
NotContains(s string) CheckFunc                               | NotContains looks for a string in the body and errs if found
Bigger(i int) CheckFunc                                       | Bigger checks if the size of the body (in bytes) in bigger than i
Strings(html *string) CheckFunc                               | Strings gets the body as plaintext
Headless() bool                                               | Headless indicates if the browser is headless (i.e. does not show on screen)
WaitForEver()                                                 | WaitForEver blocks until the timeout is reached
Dump() CheckFunc                                              | Dump prints the body and its size to log

