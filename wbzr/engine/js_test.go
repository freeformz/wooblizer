package engine_test

import (
	"strings"
	"testing"

	"github.com/woobleio/wooblizer/wbzr/engine"
)

func TestIncludeHtml(t *testing.T) {
	src := `var Woobly=function Woobly(){_classCallCheck(this,Woobly);this.document=document};`

	s, errs := engine.NewJS("objForTest", src)

	t.Log(errs)

	if len(errs) > 0 {
		t.Error("The JS class is invalid")
	}

	s.IncludeHTMLCSS("<div class='heelo' id='hello'>test</div>", "div { color: red }")

	expected := `var _sr_ = document.querySelector(_t_).attachShadow({mode:'open'});var __b = document.createElement('div');__b.setAttribute('class', 'heelo');__b.setAttribute('id', 'hello');_sr_.appendChild(__b);var __c = document.createTextNode('test');__b.appendChild(__c);this.document = _sr_;var __s = document.createElement('style');__s.innerHTML = 'div { color: red }';this.document.appendChild(__s);`

	if !strings.Contains(s.Src, expected) {
		t.Error("Includes good HTML and good CSS : Unexpected source")
	}

	s, errs = engine.NewJS("objForTest", src)

	s.IncludeHTMLCSS("", "div { color: red; }")

	expected = `var __s = document.createElement('style');__s.innerHTML = 'div { color: red; }';this.document=document;this.document.head.appendChild(__s);`

	if !strings.Contains(s.Src, expected) {
		t.Error("Includes only HTML : Unexpected source")
	}

	s.Src = ""

	err := s.IncludeHTMLCSS("<div></div>", "")

	if err == nil {
		t.Error("Includes when no doc init is present : It should returns an error")
	}
}
