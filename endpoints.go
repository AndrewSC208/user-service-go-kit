package users

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	PostUserEndpoint   endpoint.Endpoint
	GetUserEndpoint    endpoint.Endpoint
	PutUserEndpoint    endpoint.Endpoint
	PatchUserEndpoint  endpoint.Endpoint
	DeleteUserEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a users server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostProfileEndpoint:   MakePostProfileEndpoint(s),
		GetProfileEndpoint:    MakeGetProfileEndpoint(s),
		PutProfileEndpoint:    MakePutProfileEndpoint(s),
		PatchProfileEndpoint:  MakePatchProfileEndpoint(s),
		DeleteProfileEndpoint: MakeDeleteProfileEndpoint(s),
	}
}

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in a profilesvc client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	// Note that the request encoders need to modify the request URL, changing
	// the path and method. That's fine: we simply need to provide specific
	// encoders for each endpoint.

	return Endpoints{
		PostUserEndpoint:   httptransport.NewClient("POST", tgt, encodePostUserRequest, decodePostUserResponse, options...).Endpoint(),
		GetUserEndpoint:    httptransport.NewClient("GET", tgt, encodeGetUserRequest, decodeGetUserResponse, options...).Endpoint(),
		PutUserEndpoint:    httptransport.NewClient("PUT", tgt, encodePutUserRequest, decodePutUserResponse, options...).Endpoint(),
		PatchUserEndpoint:  httptransport.NewClient("PATCH", tgt, encodePatchUserRequest, decodePatchUserResponse, options...).Endpoint(),
		DeleteUserEndpoint: httptransport.NewClient("DELETE", tgt, encodeDeleteUserRequest, decodeDeleteUserResponse, options...).Endpoint(),
	}, nil
}

/**
 * METHODS
 */
// PostUser implements Service. Primarily useful in a client.
func (e Endpoints) PostUser(ctx context.Context, u User) error {
	request := postUserRequest{User: u}
	response, err := e.PostUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postUserResponse)
	return resp.Err
}

// GetUser implements Service. Primarily useful in a client.
func (e Endpoints) GetUser(ctx context.Context, username string) (User, error) {
	request := getUserRequest{Username: username}
	response, err := e.GetUserEndpoint(ctx, request)
	if err != nil {
		return User{}, err
	}
	resp := response.(getUserResponse)
	return resp.User, resp.Err
}

// PutUser implements Service. Primarily useful in a client.
func (e Endpoints) PutUser(ctx context.Context, username string, u User) error {
	request := putUserRequest{Username: username, User: u}
	response, err := e.PutUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(putUserResponse)
	return resp.Err
}

// PatchUser implements Service. Primarily useful in a client.
func (e Endpoints) PatchUser(ctx context.Context, username string, u User) error {
	request := patchUserRequest{Username: username, User: u}
	response, err := e.PatchUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(patchUserResponse)
	return resp.Err
}

// DeleteUser implements Service. Primarily useful in a client.
func (e Endpoints) DeleteUser(ctx context.Context, username string) error {
	request := deleteUserRequest{Username: username}
	response, err := e.DeleteUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteUserResponse)
	return resp.Err
}
/**
 * ENDPOINT FACTORIES
 */
// MakePostUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postUserRequest)
		e := s.PostUser(ctx, req.User)
		return postUserResponse{Err: e}, nil
	}
}

// MakeGetUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUserRequest)
		u, e := s.GetUser(ctx, req.User)
		return getUserResponse{User: u, Err: e}, nil
	}
}

// MakePutUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePutUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putUserRequest)
		e := s.PutUser(ctx, req.Username, req.User)
		return putUserResponse{Err: e}, nil
	}
}

// MakePatchUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePatchUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(patchUserRequest)
		e := s.PatchUser(ctx, req.Username, req.User)
		return patchUserResponse{Err: e}, nil
	}
}

// MakeDeleteUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteUserRequest)
		e := s.DeleteUser(ctx, req.Username)
		return deleteUserResponse{Err: e}, nil
	}
}

// We have two options to return errors from the business logic.
//
// We could return the error via the endpoint itself. That makes certain things
// a little bit easier, like providing non-200 HTTP responses to the client. But
// Go kit assumes that endpoint errors are (or may be treated as)
// transport-domain errors. For example, an endpoint error will count against a
// circuit breaker error count.
//
// Therefore, it's often better to return service (business logic) errors in the
// response object. This means we have to do a bit more work in the HTTP
// response encoder to detect e.g. a not-found error and provide a proper HTTP
// status code. That work is done with the errorer interface, in transport.go.
// Response types that may contain business-logic errors implement that
// interface.

type postUserRequest struct {
	User User
}

type postUserResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postUserResponse) error() error { return r.Err }

type getUserRequest struct {
	Username string
}

type getUserResponse struct {
	User    User    `json:"user,omitempty"`
	Err     error   `json:"err,omitempty"`
}

func (r getUserResponse) error() error { return r.Err }

type putUserRequest struct {
	Username      string
	User 		  User
}

type putUserResponse struct {
	Err error `json:"err,omitempty"`
}

func (r putUserResponse) error() error { return nil }

type patchUserRequest struct {
	Username      string
	User 		  User
}

type patchUserResponse struct {
	Err error `json:"err,omitempty"`
}

func (r patchUserResponse) error() error { return r.Err }

type deleteUserRequest struct {
	Username string
}

type deleteUserResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteUserResponse) error() error { return r.Err }