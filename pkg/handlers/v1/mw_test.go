package v1

import "testing"

func TestResponse(t *testing.T) {
	r := Response{
		Data: "any type could be used here",
	}
	if r.Error != nil {
		t.Error("Response should not contain error")
	}
}
