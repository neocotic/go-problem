// Copyright (C) 2024 neocotic
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package http

import (
	"github.com/neocotic/go-problem"
	"net/http"
)

var (
	// BadGatewayDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Bad Gateway
	// error.
	BadGatewayDefinition = problem.Definition{
		DetailKey: "problem.http.BadGatewayDefinition.detail",
		Type:      BadGateway,
	}

	// BadRequestDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Bad Request
	// error.
	BadRequestDefinition = problem.Definition{
		DetailKey: "problem.http.BadRequestDefinition.detail",
		Type:      BadRequest,
	}

	// ConflictDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Conflict
	// error.
	ConflictDefinition = problem.Definition{
		DetailKey: "problem.http.ConflictDefinition.detail",
		Type:      Conflict,
	}

	// ExpectationFailedDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Expectation Failed error.
	ExpectationFailedDefinition = problem.Definition{
		DetailKey: "problem.http.ExpectationFailedDefinition.detail",
		Type:      ExpectationFailed,
	}

	// FailedDependencyDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Failed
	// Dependency error.
	FailedDependencyDefinition = problem.Definition{
		DetailKey: "problem.http.FailedDependencyDefinition.detail",
		Type:      FailedDependency,
	}

	// ForbiddenDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Forbidden
	// error.
	ForbiddenDefinition = problem.Definition{
		DetailKey: "problem.http.ForbiddenDefinition.detail",
		Type:      Forbidden,
	}

	// GatewayTimeoutDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Gateway
	// Timeout error.
	GatewayTimeoutDefinition = problem.Definition{
		DetailKey: "problem.http.GatewayTimeoutDefinition.detail",
		Type:      GatewayTimeout,
	}

	// GoneDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Gone error.
	GoneDefinition = problem.Definition{
		DetailKey: "problem.http.GoneDefinition.detail",
		Type:      Gone,
	}

	// HTTPVersionNotSupportedDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Version Not Supported error.
	HTTPVersionNotSupportedDefinition = problem.Definition{
		DetailKey: "problem.http.HTTPVersionNotSupportedDefinition.detail",
		Type:      HTTPVersionNotSupported,
	}

	// InsufficientStorageDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Insufficient Storage error.
	InsufficientStorageDefinition = problem.Definition{
		DetailKey: "problem.http.InsufficientStorageDefinition.detail",
		Type:      InsufficientStorage,
	}

	// InternalServerDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Internal
	// Server error.
	InternalServerDefinition = problem.Definition{
		DetailKey: "problem.http.InternalServerDefinition.detail",
		Type:      InternalServer,
	}

	// LengthRequiredDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Length
	// Required error.
	LengthRequiredDefinition = problem.Definition{
		DetailKey: "problem.http.LengthRequiredDefinition.detail",
		Type:      LengthRequired,
	}

	// LockedDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Locked error.
	LockedDefinition = problem.Definition{
		DetailKey: "problem.http.LockedDefinition.detail",
		Type:      Locked,
	}

	// LoopDetectedDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Loop
	// Detected error.
	LoopDetectedDefinition = problem.Definition{
		DetailKey: "problem.http.LoopDetectedDefinition.detail",
		Type:      LoopDetected,
	}

	// MethodNotAllowedDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Method
	// Not Allowed error.
	MethodNotAllowedDefinition = problem.Definition{
		DetailKey: "problem.http.MethodNotAllowedDefinition.detail",
		Type:      MethodNotAllowed,
	}

	// MisdirectedRequestDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Misdirected Request error.
	MisdirectedRequestDefinition = problem.Definition{
		DetailKey: "problem.http.MisdirectedRequestDefinition.detail",
		Type:      MisdirectedRequest,
	}

	// NetworkAuthenticationRequiredDefinition is a built-in reusable problem.Definition that may be used to represent
	// an HTTP Network Authentication Required error.
	NetworkAuthenticationRequiredDefinition = problem.Definition{
		DetailKey: "problem.http.NetworkAuthenticationRequiredDefinition.detail",
		Type:      NetworkAuthenticationRequired,
	}

	// NotAcceptableDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Not
	// Acceptable error.
	NotAcceptableDefinition = problem.Definition{
		DetailKey: "problem.http.NotAcceptableDefinition.detail",
		Type:      NotAcceptable,
	}

	// NotFoundDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Not Found
	// error.
	NotFoundDefinition = problem.Definition{
		DetailKey: "problem.http.NotFoundDefinition.detail",
		Type:      NotFound,
	}

	// NotExtendedDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Not
	// Extended error.
	NotExtendedDefinition = problem.Definition{
		DetailKey: "problem.http.NotExtendedDefinition.detail",
		Type:      NotExtended,
	}

	// NotImplementedDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Not
	// Implemented error.
	NotImplementedDefinition = problem.Definition{
		DetailKey: "problem.http.NotImplementedDefinition.detail",
		Type:      NotImplemented,
	}

	// PaymentRequiredDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Payment
	// Required error.
	PaymentRequiredDefinition = problem.Definition{
		DetailKey: "problem.http.PaymentRequiredDefinition.detail",
		Type:      PaymentRequired,
	}

	// PreconditionFailedDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Precondition Failed error.
	PreconditionFailedDefinition = problem.Definition{
		DetailKey: "problem.http.PreconditionFailedDefinition.detail",
		Type:      PreconditionFailed,
	}

	// PreconditionRequiredDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Precondition Required error.
	PreconditionRequiredDefinition = problem.Definition{
		DetailKey: "problem.http.PreconditionRequiredDefinition.detail",
		Type:      PreconditionRequired,
	}

	// ProxyAuthRequiredDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Proxy
	// Authentication Required error.
	ProxyAuthRequiredDefinition = problem.Definition{
		DetailKey: "problem.http.ProxyAuthRequiredDefinition.detail",
		Type:      ProxyAuthRequired,
	}

	// RequestEntityTooLargeDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Request Entity Too Large error.
	RequestEntityTooLargeDefinition = problem.Definition{
		DetailKey: "problem.http.RequestEntityTooLargeDefinition.detail",
		Type:      RequestEntityTooLarge,
	}

	// RequestHeaderFieldsTooLargeDefinition is a built-in reusable problem.Definition that may be used to represent an
	// HTTP Request Header Fields Too Large error.
	RequestHeaderFieldsTooLargeDefinition = problem.Definition{
		DetailKey: "problem.http.RequestHeaderFieldsTooLargeDefinition.detail",
		Type:      RequestHeaderFieldsTooLarge,
	}

	// RequestTimeoutDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Request
	// Timeout error.
	RequestTimeoutDefinition = problem.Definition{
		DetailKey: "problem.http.RequestTimeoutDefinition.detail",
		Type:      RequestTimeout,
	}

	// RequestURITooLongDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Request URI Too Longer error.
	RequestURITooLongDefinition = problem.Definition{
		DetailKey: "problem.http.RequestURITooLongDefinition.detail",
		Type:      RequestURITooLong,
	}

	// RequestedRangeNotSatisfiableDefinition is a built-in reusable problem.Definition that may be used to represent an
	// HTTP Requested Range Not Satisfiable error.
	RequestedRangeNotSatisfiableDefinition = problem.Definition{
		DetailKey: "problem.http.RequestedRangeNotSatisfiableDefinition.detail",
		Type:      RequestedRangeNotSatisfiable,
	}

	// ServiceUnavailableDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Service Unavailable error.
	ServiceUnavailableDefinition = problem.Definition{
		DetailKey: "problem.http.ServiceUnavailableDefinition.detail",
		Type:      ServiceUnavailable,
	}

	// TeapotDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP I'm a teapot
	// error.
	TeapotDefinition = problem.Definition{
		DetailKey: "problem.http.TeapotDefinition.detail",
		Type:      Teapot,
	}

	// TooEarlyDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Too Early
	// error.
	TooEarlyDefinition = problem.Definition{
		DetailKey: "problem.http.TooEarlyDefinition.detail",
		Type:      TooEarly,
	}

	// TooManyRequestsDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Too
	// Many Requests error.
	TooManyRequestsDefinition = problem.Definition{
		DetailKey: "problem.http.TooManyRequestsDefinition.detail",
		Type:      TooManyRequests,
	}

	// UnauthorizedDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Unauthorized error.
	UnauthorizedDefinition = problem.Definition{
		DetailKey: "problem.http.UnauthorizedDefinition.detail",
		Type:      Unauthorized,
	}

	// UnavailableForLegalReasonsDefinition is a built-in reusable problem.Definition that may be used to represent an
	// HTTP Unavailable For Legal Reasons error.
	UnavailableForLegalReasonsDefinition = problem.Definition{
		DetailKey: "problem.http.UnavailableForLegalReasonsDefinition.detail",
		Type:      UnavailableForLegalReasons,
	}

	// UnprocessableEntityDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Unprocessable Entity error.
	UnprocessableEntityDefinition = problem.Definition{
		DetailKey: "problem.http.UnprocessableEntityDefinition.detail",
		Type:      UnprocessableEntity,
	}

	// UnsupportedMediaTypeDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Unsupported Media Type error.
	UnsupportedMediaTypeDefinition = problem.Definition{
		DetailKey: "problem.http.UnsupportedMediaTypeDefinition.detail",
		Type:      UnsupportedMediaType,
	}

	// UpgradeRequiredDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP Upgrade
	// Required error.
	UpgradeRequiredDefinition = problem.Definition{
		DetailKey: "problem.http.UpgradeRequiredDefinition.detail",
		Type:      UpgradeRequired,
	}

	// VariantAlsoNegotiatesDefinition is a built-in reusable problem.Definition that may be used to represent an HTTP
	// Variant Also Negotiates error.
	VariantAlsoNegotiatesDefinition = problem.Definition{
		DetailKey: "problem.http.VariantAlsoNegotiatesDefinition.detail",
		Type:      VariantAlsoNegotiates,
	}
)

// StatusDefinition returns a problem.Definition for the given HTTP status code or an empty/zero problem.Definition if
// code is unknown.
//
// For example;
//
//	StatusDefinition(400)  // BadRequestDefinition{}
//	StatusDefinition(404)  // NotFoundDefinition{}
//	StatusDefinition(999)  // problem.Definition{}
func StatusDefinition(code int) problem.Definition {
	return StatusDefinitionOrElse(code, problem.Definition{})
}

// StatusDefinitionOrElse returns a problem.Definition for the given HTTP status code or defaultDefinition if code is
// unknown.
//
// For example;
//
//	defaultDef := InternalServerDefinition{}
//	StatusDefinitionOrElse(400, defaultDef)  // BadRequestDefinition{}
//	StatusDefinitionOrElse(404, defaultDef)  // NotFoundDefinition{}
//	StatusDefinitionOrElse(999, defaultDef)  // InternalServerDefinition{}
func StatusDefinitionOrElse(code int, defaultDefinition problem.Definition) problem.Definition {
	switch code {
	case http.StatusBadRequest:
		return BadRequestDefinition
	case http.StatusUnauthorized:
		return UnauthorizedDefinition
	case http.StatusPaymentRequired:
		return PaymentRequiredDefinition
	case http.StatusForbidden:
		return ForbiddenDefinition
	case http.StatusNotFound:
		return NotFoundDefinition
	case http.StatusMethodNotAllowed:
		return MethodNotAllowedDefinition
	case http.StatusNotAcceptable:
		return NotAcceptableDefinition
	case http.StatusProxyAuthRequired:
		return ProxyAuthRequiredDefinition
	case http.StatusRequestTimeout:
		return RequestTimeoutDefinition
	case http.StatusConflict:
		return ConflictDefinition
	case http.StatusGone:
		return GoneDefinition
	case http.StatusLengthRequired:
		return LengthRequiredDefinition
	case http.StatusPreconditionFailed:
		return PreconditionFailedDefinition
	case http.StatusRequestEntityTooLarge:
		return RequestEntityTooLargeDefinition
	case http.StatusRequestURITooLong:
		return RequestURITooLongDefinition
	case http.StatusUnsupportedMediaType:
		return UnsupportedMediaTypeDefinition
	case http.StatusRequestedRangeNotSatisfiable:
		return RequestedRangeNotSatisfiableDefinition
	case http.StatusExpectationFailed:
		return ExpectationFailedDefinition
	case http.StatusTeapot:
		return TeapotDefinition
	case http.StatusMisdirectedRequest:
		return MisdirectedRequestDefinition
	case http.StatusUnprocessableEntity:
		return UnprocessableEntityDefinition
	case http.StatusLocked:
		return LockedDefinition
	case http.StatusFailedDependency:
		return FailedDependencyDefinition
	case http.StatusTooEarly:
		return TooEarlyDefinition
	case http.StatusUpgradeRequired:
		return UpgradeRequiredDefinition
	case http.StatusPreconditionRequired:
		return PreconditionRequiredDefinition
	case http.StatusTooManyRequests:
		return TooManyRequestsDefinition
	case http.StatusRequestHeaderFieldsTooLarge:
		return RequestHeaderFieldsTooLargeDefinition
	case http.StatusUnavailableForLegalReasons:
		return UnavailableForLegalReasonsDefinition
	case http.StatusInternalServerError:
		return InternalServerDefinition
	case http.StatusNotImplemented:
		return NotImplementedDefinition
	case http.StatusBadGateway:
		return BadGatewayDefinition
	case http.StatusServiceUnavailable:
		return ServiceUnavailableDefinition
	case http.StatusGatewayTimeout:
		return GatewayTimeoutDefinition
	case http.StatusHTTPVersionNotSupported:
		return HTTPVersionNotSupportedDefinition
	case http.StatusVariantAlsoNegotiates:
		return VariantAlsoNegotiatesDefinition
	case http.StatusInsufficientStorage:
		return InsufficientStorageDefinition
	case http.StatusLoopDetected:
		return LoopDetectedDefinition
	case http.StatusNotExtended:
		return NotExtendedDefinition
	case http.StatusNetworkAuthenticationRequired:
		return NetworkAuthenticationRequiredDefinition
	default:
		return defaultDefinition
	}
}
