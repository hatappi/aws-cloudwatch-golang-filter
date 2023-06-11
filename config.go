package main

import (
	"errors"
	"fmt"

	xds "github.com/cncf/xds/go/xds/type/v3"
	"github.com/envoyproxy/envoy/contrib/golang/filters/http/source/go/pkg/api"
	"github.com/envoyproxy/envoy/contrib/golang/filters/http/source/go/pkg/http"
	"google.golang.org/protobuf/types/known/anypb"
)

const Name = "aws-cloudwatch"

func init() {
	http.RegisterHttpFilterConfigFactory(Name, ConfigFactory)
	http.RegisterHttpFilterConfigParser(&parser{})
}

type config struct {
	region          string
	metricNamespace string
	metricName      string
}

type parser struct {
}

func (p *parser) Parse(any *anypb.Any) (interface{}, error) {
	configStruct := &xds.TypedStruct{}
	if err := any.UnmarshalTo(configStruct); err != nil {
		return nil, err
	}

	v := configStruct.Value
	conf := &config{}

	region, ok := v.AsMap()["region"]
	if !ok {
		return nil, errors.New("missing `region`")
	}
	if str, ok := region.(string); ok {
		conf.region = str
	} else {
		return nil, fmt.Errorf("region: expect string while got %T", region)
	}

	metricNamespace, ok := v.AsMap()["metric_namespace"]
	if !ok {
		return nil, errors.New("missing `metric_namespace`")
	}
	if str, ok := metricNamespace.(string); ok {
		conf.metricNamespace = str
	} else {
		return nil, fmt.Errorf("metric_namespace: expect string while got %T", metricNamespace)
	}

	metricName, ok := v.AsMap()["metric_name"]
	if !ok {
		return nil, errors.New("missing `metric_name`")
	}
	if str, ok := metricName.(string); ok {
		conf.metricName = str
	} else {
		return nil, fmt.Errorf("metric_name: expect string while got %T", metricName)
	}

	return conf, nil
}

func (p *parser) Merge(parent interface{}, child interface{}) interface{} {
	parentConfig := parent.(*config)
	childConfig := child.(*config)

	// copy one, do not update parentConfig directly.
	newConfig := *parentConfig
	if childConfig.region != "" {
		newConfig.region = childConfig.region
	}
	if childConfig.metricNamespace != "" {
		newConfig.metricNamespace = childConfig.metricNamespace
	}
	if childConfig.metricName != "" {
		newConfig.metricName = childConfig.metricName
	}

	return &newConfig
}

func ConfigFactory(c interface{}) api.StreamFilterFactory {
	return func(callbacks api.FilterCallbackHandler) api.StreamFilter {
		return &filter{
			callbacks: callbacks,
			config:    c.(*config),
		}
	}
}

func main() {}
