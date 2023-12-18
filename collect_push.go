package main

import (
	"context"

	"github.com/golang/snappy"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	dto "github.com/prometheus/client_model/go"

	"github.com/prometheus/prometheus/storage/remote"
	"github.com/prometheus/prometheus/util/fmtutil"

	"github.com/prometheus/client_golang/prometheus"
)

// Push metric to a remote write url
func Push(client *remote.Client, data []byte) error {
	// Encode the request body into snappy encoding.
	compressed := snappy.Encode(nil, data)
	err := client.Store(context.Background(), compressed, 0)
	if err != nil {
		return err
	}
	return nil
}

func CollectAndEncode(logger log.Logger, jobLabel string, labels map[string]string) []byte {
	var raw []byte
	up := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "up",
			Help: "the prometheus keepalive agent up metric",
		},
	)

	reg := prometheus.NewRegistry()
	reg.MustRegister(up)

	up.Add(1)

	mfs, err := reg.Gather()
	if err != nil {
		level.Error(logger).Log("msg", "Could not collect metric", "err", err) // #nosec G104
		return raw
	}

	// Convert []*dto.MetricFamily to a map to create a write request
	mfsMap := map[string]*dto.MetricFamily{"up": mfs[0]}

	// add job labels to labels
	labels["job"] = jobLabel

	// create the write request
	metricsData, err := fmtutil.MetricFamiliesToWriteRequest(mfsMap, labels)
	if err != nil {
		level.Error(logger).Log("msg", "Could not create remote write request", "err", err) // #nosec G104
		return raw
	}

	raw, err = metricsData.Marshal()
	if err != nil {
		level.Error(logger).Log("msg", "Could not encode metric data", "err", err) // #nosec G104
	}
	return raw
}
