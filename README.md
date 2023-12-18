# promk

## Prometheus Keepalive Agent

An agent to send `up` prometheus metric with a given interval to a remote write.

When using the prometheus agent or grafana agent, we could received an keepalive alert for a node.
At this stage we don't know if it's the node or the prometheus/grafana agent which is down.
There is many reasons of agent failures like config, memory errors.

This tool is complementary with any other agent and aims to reduce false positive caused by an agent down.

### Usage

```
usage: promk --remote-write-url=REMOTE-WRITE-URL [<flags>]

Prometheus Keepalive Agent.


Flags:
  -h, --[no-]help          Show context-sensitive help (also try --help-long and --help-man).
      --remote-write-url=REMOTE-WRITE-URL  
                           Prometheus remote-write url ($PROMK_URL)
      --basic-auth.username=BASIC-AUTH.USERNAME  
                           Prometheus remote-write username ($PROMK_USERNAME)
      --basic-auth.password=BASIC-AUTH.PASSWORD  
                           Prometheus remote-write password ($PROMK_PASSWORD)
      --[no-]client-tls-skip-verify  
                           Prometheus remote-write skip TLS verify ($PROMK_SKIP_TLS_VERIFY)
      --client-tls-cert-path=CLIENT-TLS-CERT-PATH  
                           Prometheus remote-write client TLS certificate path ($PROMK_CLIENT_TLS_CERT_PATH)
      --client-tls-key-path=CLIENT-TLS-KEY-PATH  
                           Prometheus remote-write client TLS key path ($PROMK_CLIENT_TLS_KEY_PATH)
      --client-tls-ca-path=CLIENT-TLS-CA-PATH  
                           Prometheus remote-write client TLS ca path ($PROMK_CLIENT_TLS_CA_PATH)
      --job-label="promk"  Job label to attach.
      --labels=LABELS ...  Add labels to the metric
      --push-interval=30s  Time Internal to push the metric
      --push-timeout=15s   Timeout to push the metric
      --push-headers=PUSH-HEADERS ...  
                           Add headers for the push request
      --log.level=info     Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt  Output format of log messages. One of: [logfmt, json]
      --[no-]version       Show application version.
```

## Sources

- [promtool](https://prometheus.io/docs/prometheus/latest/command-line/promtool/)
- [monitoring grafana agent](https://grafana.com/blog/2020/11/18/best-practices-for-meta-monitoring-the-grafana-agent/)