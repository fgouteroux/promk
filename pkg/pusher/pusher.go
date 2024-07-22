package pusher

import (
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kingpin/v2"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"

	"github.com/prometheus/prometheus/storage/remote"

	customlog "github.com/fgouteroux/promk/pkg/log"
)

type Pusher struct {
	logger   log.Logger
	client   *remote.Client
	interval time.Duration
	jobLabel string
	labels   map[string]string
}

func Setup() (p *Pusher) {
	app := kingpin.New(filepath.Base(os.Args[0]), "Prometheus Keepalive Agent.").UsageWriter(os.Stdout)
	baseURL := app.Flag("remote-write-url", "Prometheus remote-write url").Envar("PROMK_URL").Required().URL()
	username := app.Flag("basic-auth.username", "Prometheus remote-write username").Envar("PROMK_USERNAME").String()
	password := app.Flag("basic-auth.password", "Prometheus remote-write password").Envar("PROMK_PASSWORD").String()
	clientSkipTLSVerify := app.Flag("client-tls-skip-verify", "Prometheus remote-write skip TLS verify").Envar("PROMK_SKIP_TLS_VERIFY").Bool()
	clientTLSCertPath := app.Flag("client-tls-cert-path", "Prometheus remote-write client TLS certificate path").Envar("PROMK_CLIENT_TLS_CERT_PATH").String()
	clientTLSKeyPath := app.Flag("client-tls-key-path", "Prometheus remote-write client TLS key path").Envar("PROMK_CLIENT_TLS_KEY_PATH").String()
	clientTLSCaPath := app.Flag("client-tls-ca-path", "Prometheus remote-write client TLS ca path").Envar("PROMK_CLIENT_TLS_CA_PATH").String()
	jobLabel := app.Flag("job-label", "Job label to attach.").Default("promk").String()
	labels := app.Flag("labels", "Add labels to the metric").StringMap()
	pushInterval := app.Flag("push-interval", "Time Internal to push the metric").Default("30s").Duration()
	pushTimeout := app.Flag("push-timeout", "Timeout to push the metric").Default("15s").Duration()
	pushHeaders := app.Flag("push-headers", "Add headers for the push request").StringMap()

	promlogConfig := &promlog.Config{}
	flag.AddFlags(app, promlogConfig)
	app.Version(version.Print("promk"))
	app.HelpFlag.Short('h')

	kingpin.MustParse(app.Parse(os.Args[1:]))

	logger, err := customlog.InitLogger(promlogConfig)
	if err != nil {
		var logger log.Logger
		level.Error(logger).Log("Failed to init custom logger", err)
		os.Exit(1)
	}

	level.Info(logger).Log("msg", "Starting promk", "version", version.Info())              // #nosec G104
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext()) // #nosec G104

	roundtripper, err := initHTTPTransport(*clientTLSCaPath, *clientTLSKeyPath, *clientTLSCertPath, *clientSkipTLSVerify)
	if err != nil {
		level.Error(logger).Log("msg", "Could not init remote write client roundtripper", "err", err) // #nosec G104
		os.Exit(1)
	}
	remoteWriteClient, err := initRemoteWriteClient(*baseURL, *pushTimeout, roundtripper, *username, *password, *pushHeaders)
	if err != nil {
		level.Error(logger).Log("msg", "Could not init remote write client", "err", err) // #nosec G104
		os.Exit(1)
	}

	return &Pusher{
		logger:   logger,
		client:   remoteWriteClient,
		interval: *pushInterval,
		jobLabel: *jobLabel,
		labels:   *labels,
	}
}

func (p *Pusher) Run() {
	// create a new Ticker
	tk := time.NewTicker(p.interval)

	// start the ticker
	for range tk.C {
		metric := CollectAndEncode(p.logger, p.jobLabel, p.labels)
		err := Push(p.client, metric)
		if err != nil {
			level.Error(p.logger).Log("msg", "Could not push to the remote write", "err", err) // #nosec G104
		} else {
			level.Info(p.logger).Log("msg", "Successful push to the remote write") // #nosec G104
		}
	}
}
