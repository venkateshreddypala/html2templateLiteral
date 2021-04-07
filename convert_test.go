package html2templateLiteral_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/domano/html2jsx"
)

type ConvertTest struct {
	In  string
	Out string
	Err error
}

func Test_Convert(t *testing.T) {
	var tests = []ConvertTest{
		{
			In:  `<div class="bg-white"></div>`,
			Out: `<div className="bg-white"></div>`,
			Err: nil,
		},
		{
			In:  `<div class="bg-white">abc</div>`,
			Out: `<div className="bg-white">abc</div>`,
			Err: nil,
		},
		{
			In:  `<div custom-attribute="foo"></div>`,
			Out: `<div customAttribute="foo"></div>`,
			Err: nil,
		},
		{
			In:  `<div custom-attribute="foo">abc</div>`,
			Out: `<div customAttribute="foo">abc</div>`,
			Err: nil,
		},
		{
			In:  `<div aria-hidden="true" class="bg-white"></div>`,
			Out: `<div aria-hidden="true" className="bg-white"></div>`,
			Err: nil,
		},
		{
			In:  `<div aria-hidden="true" class="bg-white">abc</div>`,
			Out: `<div aria-hidden="true" className="bg-white">abc</div>`,
			Err: nil,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			rIn := strings.NewReader(test.In)
			wOut := &bytes.Buffer{}

			err := html2jsx.Convert(rIn, wOut)

			if errors.Unwrap(err) != test.Err {
				t.Errorf("Errors did not match: expected %s but got %s", test.Err, err)
			}
			if out := wOut.String(); out != test.Out {
				t.Errorf("Outputs did not match: expected \n%s but got \n%s", test.Out, out)
			}

		})
	}
}

func RunTest(ct ConvertTest) {

}
