/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package trace

import (
	"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type fakeSpanFactory struct{}

func (fakeSpanFactory) New(Span, string) Span                                     { return fakeSpan{} }
func (fakeSpanFactory) NewClientSpan(parent Span, serviceName, label string) Span { return fakeSpan{} }
func (fakeSpanFactory) FromContext(context.Context) (Span, bool)                  { return nil, false }
func (fakeSpanFactory) NewContext(parent context.Context, _ Span) context.Context { return parent }
func (fakeSpanFactory) AddGrpcServerOptions(addInterceptors func(s grpc.StreamServerInterceptor, u grpc.UnaryServerInterceptor)) {
}
func (fakeSpanFactory) AddGrpcClientOptions(addInterceptors func(s grpc.StreamClientInterceptor, u grpc.UnaryClientInterceptor)) {
}

// fakeSpan implements Span with no-op methods.
type fakeSpan struct{}

func (fakeSpan) Finish()                      {}
func (fakeSpan) Annotate(string, interface{}) {}

func init() {
	tracingBackendFactories["noop"] = func(_ string) (tracingService, io.Closer, error) {
		return fakeSpanFactory{}, &nilCloser{}, nil
	}
}
