package main

import "testing"

func TestCheckIsProxyHeaderSet(t *testing.T) {
	emptyProxyHeaderMock := map[string][]string{"": {""}}
	xRealIpProxyHeaderMock := map[string][]string{"X-Real-Ip": {"10.0.0.1"}}
	xForwardedForProxyHeaderMock := map[string][]string{"Accept": {"*/*"}, "X-Forwarded-For": {"10.0.0.1"}}
	multiXProxyHeaderMock := map[string][]string{"X-Forwarded-For": {"10.0.0.2"}, "X-Real-Ip": {"10.0.0.1"}}

	t.Run(`Return "" and error if proxy key not found`, func(t *testing.T) {
		got, err := isProxyHeaderSet(emptyProxyHeaderMock)

		if got != "" && err != nil {
			t.Errorf("Got header %q and err %v", got, err)
		}
	},
	)

	t.Run(`Return header key "X-Real-Ip" when found and nil error`, func(t *testing.T) {
		got, err := isProxyHeaderSet(xRealIpProxyHeaderMock)

		if got != "X-Real-Ip" || err != nil {
			t.Errorf("Got header %q and err %v", got, err)
		}
	},
	)

	t.Run(`Return header key "X-Forwarded" when found and nil error from multiple headers`, func(t *testing.T) {
		got, err := isProxyHeaderSet(xForwardedForProxyHeaderMock)

		if got != "X-Forwarded-For" || err != nil {
			t.Errorf("Got header %q and err %v", got, err)
		}
	},
	)

	t.Run(`Return "X-Real-Ip" header key when found and nil error if multiple proxy headers`, func(t *testing.T) {
		got, err := isProxyHeaderSet(multiXProxyHeaderMock)

		if got != "X-Real-Ip" || err != nil {
			t.Errorf("Got header %q and err %v", got, err)
		}
	},
	)
}
