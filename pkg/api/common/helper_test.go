package common

import (
	"fmt"
	"testing"
)

func TestValidateDNSPrefix(t *testing.T) {
	cases := []struct {
		dnsPrefix   string
		expectedErr error
	}{
		{
			"validDnsPrefix",
			nil,
		},
		{
			"",
			fmt.Errorf("DNSPrefix '' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 0)"),
		},
		{
			"a",
			fmt.Errorf("DNSPrefix 'a' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 1)"),
		},
		{
			"1234",
			fmt.Errorf("DNSPrefix '1234' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 4)"),
		},
		{
			"verylongdnsprefixthatismorethan45characterslong",
			fmt.Errorf("DNSPrefix 'verylongdnsprefixthatismorethan45characterslong' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 47)"),
		},
		{
			"dnswith_special?char",
			fmt.Errorf("DNSPrefix 'dnswith_special?char' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 20)"),
		},
		{
			"myDNS-1234",
			nil,
		},
	}

	for _, c := range cases {
		err := ValidateDNSPrefix(c.dnsPrefix)
		if err != nil && c.expectedErr != nil {
			if err.Error() != c.expectedErr.Error() {
				t.Fatalf("expected validateDNSPrefix to return error %s, but instead got %s", c.expectedErr.Error(), err.Error())
			}
		} else {
			if c.expectedErr != nil {
				t.Fatalf("expected validateDNSPrefix to return error %s, but instead got no error", c.expectedErr.Error())
			} else if err != nil {
				t.Fatalf("expected validateDNSPrefix to return no error, but instead got %s", err.Error())
			}
		}
	}
}
