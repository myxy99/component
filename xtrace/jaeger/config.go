package jaeger

import (
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xdefer"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-server/xapp"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jCfg "github.com/uber/jaeger-client-go/config"
	"os"
	"time"
)

type Config struct {
	ServiceName      string                `mapStructure:"service_name"`
	Sampler          *jCfg.SamplerConfig   `mapStructure:"sampler"`
	Reporter         *jCfg.ReporterConfig  `mapStructure:"reporter"`
	Headers          *jaeger.HeadersConfig `mapStructure:"headers"`
	EnableRPCMetrics bool                  `mapStructure:"enable_rpc_metrics"`
	Tags             []opentracing.Tag     `mapStructure:"tags"`
	Options          []jCfg.Option         `mapStructure:"options"`
	PanicOnError     bool                  `mapStructure:"panic_on_error"`
}

func DefaultConfig() *Config {
	agentAddr := "127.0.0.1:6831"
	headerName := "x-trace-id"
	if addr := os.Getenv("JAEGER_AGENT_ADDR"); addr != "" {
		agentAddr = addr
	}
	return &Config{
		ServiceName: xapp.Name(),
		Sampler: &jCfg.SamplerConfig{
			Type:  "const",
			Param: 0.001,
		},
		Reporter: &jCfg.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentAddr,
		},
		EnableRPCMetrics: true,
		Headers: &jaeger.HeadersConfig{
			TraceBaggageHeaderPrefix: "ctx-",
			TraceContextHeaderName:   headerName,
		},
		Tags: []opentracing.Tag{
			{Key: "host_name", Value: xapp.HostName()},
			{Key: "host_ip", Value: xapp.HostIP()},
			{Key: "app_id", Value: xapp.AppId()},
			{Key: "app_name", Value: xapp.Name()},
		},
		PanicOnError: true,
	}
}

func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := xcfg.UnmarshalKey(key, config); err != nil {
		xlog.Panic("Application Starting",
			xlog.FieldComponentName("XTrace"),
			xlog.FieldMethod("XTrace.Jaeger.RawConfig"),
			xlog.FieldDescription("UnmarshalKey Error"),
			xlog.FieldErr(err),
		)
	}
	return config
}

func (config *Config) WithTag(tags ...opentracing.Tag) *Config {
	if config.Tags == nil {
		config.Tags = make([]opentracing.Tag, 0)
	}
	config.Tags = append(config.Tags, tags...)
	return config
}

func (config *Config) WithOption(options ...jCfg.Option) *Config {
	if config.Options == nil {
		config.Options = make([]jCfg.Option, 0)
	}
	config.Options = append(config.Options, options...)
	return config
}

func (config *Config) Build() opentracing.Tracer {
	var configuration = jCfg.Configuration{
		ServiceName: config.ServiceName,
		Sampler:     config.Sampler,
		Reporter:    config.Reporter,
		RPCMetrics:  config.EnableRPCMetrics,
		Headers:     config.Headers,
		Tags:        config.Tags,
	}
	tracer, closer, err := configuration.NewTracer(config.Options...)
	if err != nil {
		if config.PanicOnError {
			xlog.Panic("Application Starting",
				xlog.FieldComponentName("XTrace"),
				xlog.FieldMethod("XTrace.Jaeger.Build"),
				xlog.FieldDescription("Build Jaeger Error"),
				xlog.FieldErr(err),
			)
		} else {
			xlog.Error("Application Starting",
				xlog.FieldComponentName("XTrace"),
				xlog.FieldMethod("XTrace.Jaeger.Build"),
				xlog.FieldDescription("Build Jaeger Error"),
				xlog.FieldErr(err),
			)
		}
	}
	xdefer.Register(func() error {
		xlog.Info("Application Stopping",
			xlog.FieldComponentName("XTrace"),
			xlog.FieldMethod("Jaeger.Build.Register"),
			xlog.FieldDescription("XTrace jaeger server shutdown"),
		)
		return closer.Close()
	})
	return tracer
}
