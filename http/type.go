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
	// BadGateway is a built-in reusable problem.Type that may be used to represent an HTTP Bad Gateway error.
	BadGateway = problem.Type{
		LogLevel: problem.LogLevelError,
		Status:   http.StatusBadGateway,
		Title:    http.StatusText(http.StatusBadGateway),
		TitleKey: "problem.http.BadGateway.title",
	}

	// BadRequest is a built-in reusable problem.Type that may be used to represent an HTTP Bad Request error.
	BadRequest = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusBadRequest,
		Title:    http.StatusText(http.StatusBadRequest),
		TitleKey: "problem.http.BadRequest.title",
	}

	// Conflict is a built-in reusable problem.Type that may be used to represent an HTTP Conflict error.
	Conflict = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusConflict,
		Title:    http.StatusText(http.StatusConflict),
		TitleKey: "problem.http.Conflict.title",
	}

	// ExpectationFailed is a built-in reusable problem.Type that may be used to represent an HTTP Expectation Failed
	// error.
	ExpectationFailed = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusExpectationFailed,
		Title:    http.StatusText(http.StatusExpectationFailed),
		TitleKey: "problem.http.ExpectationFailed.title",
	}

	// FailedDependency is a built-in reusable problem.Type that may be used to represent an HTTP Failed Dependency
	// error.
	FailedDependency = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusFailedDependency,
		Title:    http.StatusText(http.StatusFailedDependency),
		TitleKey: "problem.http.FailedDependency.title",
	}

	// Forbidden is a built-in reusable problem.Type that may be used to represent an HTTP Forbidden error.
	Forbidden = problem.Type{
		LogLevel: problem.LogLevelWarn,
		Status:   http.StatusForbidden,
		Title:    http.StatusText(http.StatusForbidden),
		TitleKey: "problem.http.Forbidden.title",
	}

	// GatewayTimeout is a built-in reusable problem.Type that may be used to represent an HTTP Gateway Timeout error.
	GatewayTimeout = problem.Type{
		LogLevel: problem.LogLevelError,
		Status:   http.StatusGatewayTimeout,
		Title:    http.StatusText(http.StatusGatewayTimeout),
		TitleKey: "problem.http.GatewayTimeout.title",
	}

	// Gone is a built-in reusable problem.Type that may be used to represent an HTTP Gone error.
	Gone = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusGone,
		Title:    http.StatusText(http.StatusGone),
		TitleKey: "problem.http.Gone.title",
	}

	// HTTPVersionNotSupported is a built-in reusable problem.Type that may be used to represent an HTTP Version Not
	// Supported error.
	HTTPVersionNotSupported = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusHTTPVersionNotSupported,
		Title:    http.StatusText(http.StatusHTTPVersionNotSupported),
		TitleKey: "problem.http.HTTPVersionNotSupported.title",
	}

	// InsufficientStorage is a built-in reusable problem.Type that may be used to represent an HTTP Insufficient
	// Storage error.
	InsufficientStorage = problem.Type{
		LogLevel: problem.LogLevelError,
		Status:   http.StatusInsufficientStorage,
		Title:    http.StatusText(http.StatusInsufficientStorage),
		TitleKey: "problem.http.InsufficientStorage.title",
	}

	// InternalServer is a built-in reusable problem.Type that may be used to represent an HTTP Internal Server error.
	InternalServer = problem.Type{
		LogLevel: problem.LogLevelError,
		Status:   http.StatusInternalServerError,
		Title:    http.StatusText(http.StatusInternalServerError),
		TitleKey: "problem.http.InternalServer.title",
	}

	// LengthRequired is a built-in reusable problem.Type that may be used to represent an HTTP Length Required error.
	LengthRequired = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusLengthRequired,
		Title:    http.StatusText(http.StatusLengthRequired),
		TitleKey: "problem.http.LengthRequired.title",
	}

	// Locked is a built-in reusable problem.Type that may be used to represent an HTTP Locked error.
	Locked = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusLocked,
		Title:    http.StatusText(http.StatusLocked),
		TitleKey: "problem.http.Locked.title",
	}

	// LoopDetected is a built-in reusable problem.Type that may be used to represent an HTTP Loop Detected error.
	LoopDetected = problem.Type{
		LogLevel: problem.LogLevelError,
		Status:   http.StatusLoopDetected,
		Title:    http.StatusText(http.StatusLoopDetected),
		TitleKey: "problem.http.LoopDetected.title",
	}

	// MethodNotAllowed is a built-in reusable problem.Type that may be used to represent an HTTP Method Not Allowed
	// error.
	MethodNotAllowed = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusMethodNotAllowed,
		Title:    http.StatusText(http.StatusMethodNotAllowed),
		TitleKey: "problem.http.MethodNotAllowed.title",
	}

	// MisdirectedRequest is a built-in reusable problem.Type that may be used to represent an HTTP Misdirected Request
	// error.
	MisdirectedRequest = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusMisdirectedRequest,
		Title:    http.StatusText(http.StatusMisdirectedRequest),
		TitleKey: "problem.http.MisdirectedRequest.title",
	}

	// NetworkAuthenticationRequired is a built-in reusable problem.Type that may be used to represent an HTTP Network
	// Authentication Required error.
	NetworkAuthenticationRequired = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusNetworkAuthenticationRequired,
		Title:    http.StatusText(http.StatusNetworkAuthenticationRequired),
		TitleKey: "problem.http.NetworkAuthenticationRequired.title",
	}

	// NotAcceptable is a built-in reusable problem.Type that may be used to represent an HTTP Not Acceptable error.
	NotAcceptable = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusNotAcceptable,
		Title:    http.StatusText(http.StatusNotAcceptable),
		TitleKey: "problem.http.NotAcceptable.title",
	}

	// NotExtended is a built-in reusable problem.Type that may be used to represent an HTTP Not Extended error.
	NotExtended = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusNotExtended,
		Title:    http.StatusText(http.StatusNotExtended),
		TitleKey: "problem.http.NotExtended.title",
	}

	// NotFound is a built-in reusable problem.Type that may be used to represent an HTTP Not Found error.
	NotFound = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusNotFound,
		Title:    http.StatusText(http.StatusNotFound),
		TitleKey: "problem.http.NotFound.title",
	}

	// NotImplemented is a built-in reusable problem.Type that may be used to represent an HTTP Not Implemented error.
	NotImplemented = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusNotImplemented,
		Title:    http.StatusText(http.StatusNotImplemented),
		TitleKey: "problem.http.NotImplemented.title",
	}

	// PaymentRequired is a built-in reusable problem.Type that may be used to represent an HTTP Payment Required error.
	PaymentRequired = problem.Type{
		LogLevel: problem.LogLevelWarn,
		Status:   http.StatusPaymentRequired,
		Title:    http.StatusText(http.StatusPaymentRequired),
		TitleKey: "problem.http.PaymentRequired.title",
	}

	// PreconditionFailed is a built-in reusable problem.Type that may be used to represent an HTTP Precondition Failed
	// error.
	PreconditionFailed = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusPreconditionFailed,
		Title:    http.StatusText(http.StatusPreconditionFailed),
		TitleKey: "problem.http.PreconditionFailed.title",
	}

	// PreconditionRequired is a built-in reusable problem.Type that may be used to represent an HTTP Precondition
	// Required error.
	PreconditionRequired = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusPreconditionRequired,
		Title:    http.StatusText(http.StatusPreconditionRequired),
		TitleKey: "problem.http.PreconditionRequired.title",
	}

	// ProxyAuthRequired is a built-in reusable problem.Type that may be used to represent an HTTP Proxy Authentication
	// Required error.
	ProxyAuthRequired = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusProxyAuthRequired,
		Title:    http.StatusText(http.StatusProxyAuthRequired),
		TitleKey: "problem.http.ProxyAuthRequired.title",
	}

	// RequestEntityTooLarge is a built-in reusable problem.Type that may be used to represent an HTTP Request Entity
	// Too Large error.
	RequestEntityTooLarge = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusRequestEntityTooLarge,
		Title:    http.StatusText(http.StatusRequestEntityTooLarge),
		TitleKey: "problem.http.RequestEntityTooLarge.title",
	}

	// RequestHeaderFieldsTooLarge is a built-in reusable problem.Type that may be used to represent an HTTP Request
	// Header Fields Too Large error.
	RequestHeaderFieldsTooLarge = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusRequestHeaderFieldsTooLarge,
		Title:    http.StatusText(http.StatusRequestHeaderFieldsTooLarge),
		TitleKey: "problem.http.RequestHeaderFieldsTooLarge.title",
	}

	// RequestTimeout is a built-in reusable problem.Type that may be used to represent an HTTP Request Timeout error.
	RequestTimeout = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusRequestTimeout,
		Title:    http.StatusText(http.StatusRequestTimeout),
		TitleKey: "problem.http.RequestTimeout.title",
	}

	// RequestURITooLong is a built-in reusable problem.Type that may be used to represent an HTTP Request URI Too Long
	// error.
	RequestURITooLong = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusRequestURITooLong,
		Title:    http.StatusText(http.StatusRequestURITooLong),
		TitleKey: "problem.http.RequestURITooLong.title",
	}

	// RequestedRangeNotSatisfiable is a built-in reusable problem.Type that may be used to represent an HTTP Requested
	// Range Not Satisfiable error.
	RequestedRangeNotSatisfiable = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusRequestedRangeNotSatisfiable,
		Title:    http.StatusText(http.StatusRequestedRangeNotSatisfiable),
		TitleKey: "problem.http.RequestedRangeNotSatisfiable.title",
	}

	// ServiceUnavailable is a built-in reusable problem.Type that may be used to represent an HTTP Service Unavailable
	// error.
	ServiceUnavailable = problem.Type{
		LogLevel: problem.LogLevelError,
		Status:   http.StatusServiceUnavailable,
		Title:    http.StatusText(http.StatusServiceUnavailable),
		TitleKey: "problem.http.ServiceUnavailable.title",
	}

	// Teapot is a built-in reusable problem.Type that may be used to represent an HTTP I'm a teapot error.
	Teapot = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusTeapot,
		Title:    http.StatusText(http.StatusTeapot),
		TitleKey: "problem.http.Teapot.title",
	}

	// TooEarly is a built-in reusable problem.Type that may be used to represent an HTTP Too Early error.
	TooEarly = problem.Type{
		LogLevel: problem.LogLevelWarn,
		Status:   http.StatusTooEarly,
		Title:    http.StatusText(http.StatusTooEarly),
		TitleKey: "problem.http.TooEarly.title",
	}

	// TooManyRequests is a built-in reusable problem.Type that may be used to represent an HTTP Too Many Requests
	// error.
	TooManyRequests = problem.Type{
		LogLevel: problem.LogLevelWarn,
		Status:   http.StatusTooManyRequests,
		Title:    http.StatusText(http.StatusTooManyRequests),
		TitleKey: "problem.http.TooManyRequests.title",
	}

	// Unauthorized is a built-in reusable problem.Type that may be used to represent an HTTP Unauthorized error.
	Unauthorized = problem.Type{
		LogLevel: problem.LogLevelWarn,
		Status:   http.StatusUnauthorized,
		Title:    http.StatusText(http.StatusUnauthorized),
		TitleKey: "problem.http.Unauthorized.title",
	}

	// UnavailableForLegalReasons is a built-in reusable problem.Type that may be used to represent an HTTP Unavailable
	// For Legal Reasons error.
	UnavailableForLegalReasons = problem.Type{
		LogLevel: problem.LogLevelWarn,
		Status:   http.StatusUnavailableForLegalReasons,
		Title:    http.StatusText(http.StatusUnavailableForLegalReasons),
		TitleKey: "problem.http.UnavailableForLegalReasons.title",
	}

	// UnprocessableEntity is a built-in reusable problem.Type that may be used to represent an HTTP Unprocessable
	// Entity error.
	UnprocessableEntity = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusUnprocessableEntity,
		Title:    http.StatusText(http.StatusUnprocessableEntity),
		TitleKey: "problem.http.UnprocessableEntity.title",
	}

	// UnsupportedMediaType is a built-in reusable problem.Type that may be used to represent an HTTP Unsupported Media
	// Type error.
	UnsupportedMediaType = problem.Type{
		LogLevel: problem.LogLevelWarn,
		Status:   http.StatusUnsupportedMediaType,
		Title:    http.StatusText(http.StatusUnsupportedMediaType),
		TitleKey: "problem.http.UnsupportedMediaType.title",
	}

	// UpgradeRequired is a built-in reusable problem.Type that may be used to represent an HTTP Upgrade Required error.
	UpgradeRequired = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusUpgradeRequired,
		Title:    http.StatusText(http.StatusUpgradeRequired),
		TitleKey: "problem.http.UpgradeRequired.title",
	}

	// VariantAlsoNegotiates is a built-in reusable problem.Type that may be used to represent an HTTP Variant Also
	// Negotiates error.
	VariantAlsoNegotiates = problem.Type{
		LogLevel: problem.LogLevelError,
		Status:   http.StatusVariantAlsoNegotiates,
		Title:    http.StatusText(http.StatusVariantAlsoNegotiates),
		TitleKey: "problem.http.VariantAlsoNegotiates.title",
	}
)

// StatusType returns a problem.Type for the given HTTP status code or an empty/zero problem.Type if code is unknown.
//
// For example;
//
//	StatusType(400)  // BadRequest{}
//	StatusType(404)  // NotFound{}
//	StatusType(999)  // problem.Type{}
func StatusType(code int) problem.Type {
	return StatusTypeOrElse(code, problem.Type{})
}

// StatusTypeOrElse returns a problem.Type for the given HTTP status code or defaultType if code is unknown.
//
// For example;
//
//	defaultType := InternalServer{}
//	StatusTypeOrElse(400, defaultType)  // BadRequest{}
//	StatusTypeOrElse(404, defaultType)  // NotFound{}
//	StatusTypeOrElse(999, defaultType)  // InternalServer{}
func StatusTypeOrElse(code int, defaultType problem.Type) problem.Type {
	switch code {
	case http.StatusBadRequest:
		return BadRequest
	case http.StatusUnauthorized:
		return Unauthorized
	case http.StatusPaymentRequired:
		return PaymentRequired
	case http.StatusForbidden:
		return Forbidden
	case http.StatusNotFound:
		return NotFound
	case http.StatusMethodNotAllowed:
		return MethodNotAllowed
	case http.StatusNotAcceptable:
		return NotAcceptable
	case http.StatusProxyAuthRequired:
		return ProxyAuthRequired
	case http.StatusRequestTimeout:
		return RequestTimeout
	case http.StatusConflict:
		return Conflict
	case http.StatusGone:
		return Gone
	case http.StatusLengthRequired:
		return LengthRequired
	case http.StatusPreconditionFailed:
		return PreconditionFailed
	case http.StatusRequestEntityTooLarge:
		return RequestEntityTooLarge
	case http.StatusRequestURITooLong:
		return RequestURITooLong
	case http.StatusUnsupportedMediaType:
		return UnsupportedMediaType
	case http.StatusRequestedRangeNotSatisfiable:
		return RequestedRangeNotSatisfiable
	case http.StatusExpectationFailed:
		return ExpectationFailed
	case http.StatusTeapot:
		return Teapot
	case http.StatusMisdirectedRequest:
		return MisdirectedRequest
	case http.StatusUnprocessableEntity:
		return UnprocessableEntity
	case http.StatusLocked:
		return Locked
	case http.StatusFailedDependency:
		return FailedDependency
	case http.StatusTooEarly:
		return TooEarly
	case http.StatusUpgradeRequired:
		return UpgradeRequired
	case http.StatusPreconditionRequired:
		return PreconditionRequired
	case http.StatusTooManyRequests:
		return TooManyRequests
	case http.StatusRequestHeaderFieldsTooLarge:
		return RequestHeaderFieldsTooLarge
	case http.StatusUnavailableForLegalReasons:
		return UnavailableForLegalReasons
	case http.StatusInternalServerError:
		return InternalServer
	case http.StatusNotImplemented:
		return NotImplemented
	case http.StatusBadGateway:
		return BadGateway
	case http.StatusServiceUnavailable:
		return ServiceUnavailable
	case http.StatusGatewayTimeout:
		return GatewayTimeout
	case http.StatusHTTPVersionNotSupported:
		return HTTPVersionNotSupported
	case http.StatusVariantAlsoNegotiates:
		return VariantAlsoNegotiates
	case http.StatusInsufficientStorage:
		return InsufficientStorage
	case http.StatusLoopDetected:
		return LoopDetected
	case http.StatusNotExtended:
		return NotExtended
	case http.StatusNetworkAuthenticationRequired:
		return NetworkAuthenticationRequired
	default:
		return defaultType
	}
}
