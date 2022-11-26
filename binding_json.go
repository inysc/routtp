package routtp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/inysc/facade"
)

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// interface{} as a Number instead of as a float64.
var EnableDecoderUseNumber = false

// EnableDecoderDisallowUnknownFields is used to call the DisallowUnknownFields method
// on the JSON Decoder instance. DisallowUnknownFields causes the Decoder to
// return an error when the destination is a struct and the input contains object
// keys which do not match any non-ignored, exported fields in the destination.
var EnableDecoderDisallowUnknownFields = false

type jsonBinding struct{}

var jsonbinding = &jsonBinding{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = req.Body.Close()
	if err != nil {
		facade.Errorf("failed to close request body<%s>", err)
		return err
	}
	facade.Debugf("request body: %s", body)
	return jsonbinding.BindBody(body, obj)
}

func (jsonBinding) BindBody(body []byte, obj any) error {
	return json.Unmarshal(body, obj)
}
