package apierror

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

func TestExtractCodeFromARMHttpResponse(t *testing.T) {
	RegisterTestingT(t)

	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(`{"error":{"code":"ResourceGroupNotFound","message":"Resource group 'jiren-fakegroup' could not be found."}}`)),
	}

	code := ExtractCodeFromARMHttpResponse(resp, "")
	Expect(code).To(Equal(ErrorCode("ResourceGroupNotFound")))
}
