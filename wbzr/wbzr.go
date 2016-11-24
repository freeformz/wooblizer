// Package wbzr provides tool to create a wooble in a given language and to package
// (wrap) some woobles.
package wbzr

import (
	"bytes"
	"errors"
	"io/ioutil"
	"text/template"

	"github.com/woobleio/wooblizer/wbzr/engine"
)

// ScriptLang are constants for implemented script languages.
type ScriptLang int

const (
	JSES5 ScriptLang = iota
)

type wbzr struct {
	DomainsSec []string
	Scripts    []engine.Script

	lang     ScriptLang
	skeleton string
}

// New takes a script language which is used to inject and output a file.
func New(sl ScriptLang) *wbzr {
	var skeleton string
	switch sl {
	case JSES5:
		skeleton = wbJses5
	default:
		panic("Language not supported")
	}

	return &wbzr{
		nil,
		make([]engine.Script, 0),
		sl,
		skeleton,
	}
}

// Get returns an injected source.
func (wb *wbzr) Get(name string) (engine.Script, error) {
	for _, sc := range wb.Scripts {
		if sc.GetName() == name {
			return sc, nil
		}
	}
	return nil, errors.New("Wooble " + name + " not found")
}

// Inject injects a source code to be wooblized. It takes a name which must be
// unique. Src can be empty, it'll create a default object
func (wb *wbzr) Inject(src string, name string) (engine.Script, error) {
	if _, err := wb.Get(name); err == nil {
		return nil, errors.New("Wooble " + name + " already exists. A name must be unique")
	}
	var sc engine.Script
	var err error

	switch wb.lang {
	case JSES5:
		sc, err = engine.NewJSES5(src, name)
	}

	if err != nil {
		return nil, err
	}

	wb.Scripts = append(wb.Scripts, sc)
	return sc, nil
}

// InjectFile injects a source from a file.
func (wb *wbzr) InjectFile(path string, name string) (engine.Script, error) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return wb.Inject(string(c[:]), name)
}

// Secure set some domains to protect the script and make it works only for specific domains
func (wb *wbzr) Secure(domains ...string) {
	wb.DomainsSec = domains
}

func (wb *wbzr) SecureAndWrap(domains ...string) (*bytes.Buffer, error) {
	wb.Secure(domains...)
	return wb.Wrap()
}

// Wrap packages some woobles (all the woobles injected in the wbzr)
// and build a file which contains the wooble library.
func (wb *wbzr) Wrap() (*bytes.Buffer, error) {
	for _, sc := range wb.Scripts {
		if _, err := sc.Build(); err != nil {
			return nil, err
		}
	}

	tmpl := template.Must(template.New("WbJSES5").Parse(wbJses5))

	var out bytes.Buffer
	if err := tmpl.Execute(&out, wb); err != nil {
		return nil, err
	}

	return &out, nil
}

var wbJses5 = `
function Wb(id) {
	{{if .DomainsSec}}
	var ah = [{{range $i, $o := .DomainsSec}}"{{$o}}"{{if not $i}},{{end}}{{end}}];
  var xx = ah.indexOf(window.location.hostname);
  if(ah.indexOf(window.location.hostname) == -1) {
  	console.log("Wooble error : domain restricted");
    return;
  }
	{{end}}

	if(window === this) {
  	return new Wb(id);
  }

  var cs = {
  	{{range $i, $o := .Scripts}}"{{$o.GetName}}":{{"{"}}{{$o.GetSource}}{{"}"}}{{if not $i}},{{end}}{{end}}
  }

  var c = cs[id];
  if(typeof c == 'undefined') {
  	console.log("Wooble error : creation", id, "not found");
    return undefined;
  }

  this.init = function (target) {
    if(document.querySelector(target) == null) {
    	console.log("Wooble error : Element", target, "not found in the document");
      return;
    }

    if("_buildDoc" in c) c._buildDoc(target);
    if("_buildStyle" in c) c._buildStyle();
    if("_init" in c) c._init();
  }

  this.get = function() {
  	return c;
  }

  return this;
}
`
