package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/envoyproxy/envoy/contrib/golang/filters/http/source/go/pkg/api"
)

type request struct {
	host   string
	method string
	path   string
}

type filter struct {
	callbacks api.FilterCallbackHandler
	config    *config

	request *request
}

func (f *filter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	host, _ := header.Get(":authority")
	method, _ := header.Get(":method")
	path, _ := header.Get(":path")

	f.request = &request{
		host:   host,
		method: method,
		path:   path,
	}

	return api.Continue
}

func (f *filter) DecodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	return api.Continue
}

func (f *filter) DecodeTrailers(trailers api.RequestTrailerMap) api.StatusType {
	return api.Continue
}

func (f *filter) EncodeHeaders(header api.ResponseHeaderMap, endStream bool) api.StatusType {
	status, _ := header.Get(":status")

	input := &cloudwatch.PutMetricDataInput{
		Namespace: aws.String(f.config.metricNamespace),
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String(f.config.metricName),
				Timestamp:  aws.Time(time.Now().UTC()),
				Unit:       aws.String("Count"),
				Value:      aws.Float64(1),
				Dimensions: []*cloudwatch.Dimension{
					{Name: aws.String("host"), Value: aws.String(f.request.host)},
					{Name: aws.String("method"), Value: aws.String(f.request.method)},
					{Name: aws.String("path"), Value: aws.String(f.request.path)},
					{Name: aws.String("status"), Value: aws.String(status)},
				},
			},
		},
	}

	f.callbacks.Log(api.Debug, fmt.Sprintf("input: %v", input))

	sess, err := session.NewSession()
	if err != nil {
		f.callbacks.Log(api.Error, fmt.Sprintf("failed to create a new session: %v", err))
	}
	cw := cloudwatch.New(sess, aws.NewConfig().WithRegion(f.config.region))

	output, err := cw.PutMetricData(input)
	if err != nil {
		f.callbacks.Log(api.Error, fmt.Sprintf("failed to call PutMetricData: %s", err))
		return api.Continue
	}
	f.callbacks.Log(api.Info, fmt.Sprintf("output: %s", output.GoString()))

	return api.Continue
}

func (f *filter) EncodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	return api.Continue
}

func (f *filter) EncodeTrailers(trailers api.ResponseTrailerMap) api.StatusType {
	return api.Continue
}

func (f *filter) OnDestroy(reason api.DestroyReason) {
}
