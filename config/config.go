package config

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/cihub/seelog"
	"github.com/go-ini/ini"
)

type CheckTimers struct {
	Process     *time.Ticker
	Connections *time.Ticker
	RealTime    *time.Ticker
}

type AgentConfig struct {
	Enabled       bool
	APIKey        string
	HostName      string
	APIEndpoint   *url.URL
	LogLevel      string
	QueueSize     int
	Blacklist     []*regexp.Regexp
	MaxProcFDs    int
	ProcLimit     int
	AllowRealTime bool
	Concurrency   int
	Proxy         *url.URL
	Timers        *CheckTimers
	Logger        *LoggerConfig
}

const (
	defaultEndpoint = "https://process.datadoghq.com"
	maxProcLimit    = 100
)

// NewDefaultAgentConfig returns an AgentConfig with defaults initialized
func NewDefaultAgentConfig() *AgentConfig {
	hostname, err := getHostname()
	if err != nil {
		hostname = ""
	}
	u, err := url.Parse(defaultEndpoint)
	if err != nil {
		// This is a hardcoded URL so parsing it should not fail
		panic(err)
	}
	ac := &AgentConfig{
		Enabled:       false,
		HostName:      hostname,
		APIEndpoint:   u,
		LogLevel:      "info",
		QueueSize:     20,
		MaxProcFDs:    200,
		ProcLimit:     100,
		AllowRealTime: true,
		Concurrency:   4,
		Timers: &CheckTimers{
			Process:     time.NewTicker(10 * time.Second),
			Connections: time.NewTicker(3 * 60 * time.Minute),
			RealTime:    time.NewTicker(2 * time.Second),
		},
	}

	return ac
}

// NewAgentConfig returns an AgentConfig using a conf and legacy configuration.
// conf will be nil if there is no configuration available but legacyConf will
// give an error if nil.
func NewAgentConfig(agentConf, legacyConf *File) (*AgentConfig, error) {
	cfg := NewDefaultAgentConfig()

	var ns string
	var file *File
	var section *ini.Section
	if agentConf != nil {
		section, _ = agentConf.GetSection("Main")
	}

	// Prefer the dd-agent config file.
	if section != nil {
		file = agentConf
		ns = "process.config"
		a, err := agentConf.Get("Main", "api_key")
		if err != nil {
			return nil, err
		}
		ak := strings.Split(a, ",")
		cfg.APIKey = ak[0]
		cfg.LogLevel = strings.ToLower(agentConf.GetDefault("Main", "log_level", "INFO"))
		cfg.Proxy = getProxySettings(section)
		e := agentConf.GetDefault(ns, "endpoint", defaultEndpoint)
		u, err := url.Parse(e)
		if err != nil {
			return nil, fmt.Errorf("invalid endpoint URL: %s", err)
		}
		if v, _ := agentConf.Get("Main", "process_agent_enabled"); v == "true" {
			cfg.Enabled = true
		}
		cfg.APIEndpoint = u
	}

	// But legacy conf will override dd-agent.
	if legacyConf != nil {
		file = legacyConf
		ns = "dd-process-agent"
		cfg.LogLevel = strings.ToLower(legacyConf.GetDefault(ns, "log_level", cfg.LogLevel))

		s, err := legacyConf.Get(ns, "server_url")
		if err != nil {
			return nil, err
		}
		u, err := url.Parse(s)
		if err != nil {
			return nil, err
		}
		cfg.APIEndpoint = u

		a, err := legacyConf.Get(ns, "api_key")
		if err != nil {
			return nil, err
		}
		cfg.APIKey = a
		proxy := legacyConf.GetDefault(ns, "proxy", "")
		if proxy != "" {
			cfg.Proxy, err = url.Parse(proxy)
			if err != nil {
				log.Errorf("Could not parse proxy url from configuration: %s", err)
			}
		}
	}

	// We can have no configuration in ENV-only case.
	if file != nil {
		cfg.QueueSize = file.GetIntDefault(ns, "queue_size", cfg.QueueSize)
		cfg.MaxProcFDs = file.GetIntDefault(ns, "max_proc_fds", cfg.MaxProcFDs)
		cfg.AllowRealTime = file.GetBool(ns, "allow_real_time", cfg.AllowRealTime)

		blacklistPats := file.GetStrArrayDefault(ns, "blacklist", ",", []string{})
		blacklist := make([]*regexp.Regexp, 0, len(blacklistPats))
		for _, b := range blacklistPats {
			r, err := regexp.Compile(b)
			if err == nil {
				blacklist = append(blacklist, r)
			}
		}
		cfg.Blacklist = blacklist
		procLimit := file.GetIntDefault(ns, "proc_limit", cfg.ProcLimit)
		if procLimit <= maxProcLimit {
			cfg.ProcLimit = procLimit
		} else {
			log.Warn("Overriding the configured process limit because it exceeds maximum")
			cfg.ProcLimit = maxProcLimit
		}
		cfg.Concurrency = file.GetIntDefault(ns, "concurrency", cfg.Concurrency)
		t := cfg.Timers
		t.Process = time.NewTicker(file.GetDurationDefault(ns, "process_interval", time.Second, 10*time.Second))
		t.Connections = time.NewTicker(file.GetDurationDefault(ns, "connection_interval", time.Minute, 3*60*time.Minute))
		t.RealTime = time.NewTicker(file.GetDurationDefault(ns, "realtime_interval", time.Second, 2*time.Second))
	}

	cfg = mergeEnv(cfg)

	// (Re)configure the logging from our configuration
	NewLoggerLevel(cfg.LogLevel)

	return cfg, nil
}

// mergeEnv applies overrides from environment variables to the trace agent configuration
func mergeEnv(c *AgentConfig) *AgentConfig {
	if v := os.Getenv("DD_PROCESS_AGENT_ENABLED"); v == "true" {
		c.Enabled = true
	} else if v == "false" {
		c.Enabled = false
	}

	if v := os.Getenv("DD_HOSTNAME"); v != "" {
		log.Info("overriding hostname from env DD_HOSTNAME value")
		c.HostName = v
	}

	// Support API_KEY and DD_API_KEY but prefer DD_API_KEY.
	var apiKey string
	if v := os.Getenv("API_KEY"); v != "" {
		apiKey = v
		log.Info("overriding API key from env API_KEY value")
	}
	if v := os.Getenv("DD_API_KEY"); v != "" {
		apiKey = v
		log.Info("overriding API key from env DD_API_KEY value")
	}
	if apiKey != "" {
		vals := strings.Split(apiKey, ",")
		for i := range vals {
			vals[i] = strings.TrimSpace(vals[i])
		}
		c.APIKey = vals[0]
	}

	// Support LOG_LEVEL and DD_LOG_LEVEL but prefer DD_LOG_LEVE
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		c.LogLevel = v
	}
	if v := os.Getenv("DD_LOG_LEVEL"); v != "" {
		c.LogLevel = v
	}

	c.Proxy = proxyFromEnv(c.Proxy)

	if v := os.Getenv("DD_PROCESS_AGENT_URL"); v != "" {
		u, err := url.Parse(v)
		if err != nil {
			log.Warnf("DD_PROCESS_AGENT_URL is invalid: %s", err)
		} else {
			log.Infof("overriding API endpoint from env")
			c.APIEndpoint = u
		}
	}

	return c
}

// IsBlacklisted returns a boolean indicating if the given command is blacklisted by our config.
func IsBlacklisted(cmdline []string, blacklist []*regexp.Regexp) bool {
	cmd := strings.Join(cmdline, " ")
	for _, b := range blacklist {
		if b.MatchString(cmd) {
			return true
		}
	}
	return false
}

// getHostname shells out to obtain the hostname used by the infra agent
// falling back to os.Hostname() if it is unavailable
func getHostname() (string, error) {
	ddAgentPy := "/opt/datadog-agent/embedded/bin/python"
	getHostnameCmd := "from utils.hostname import get_hostname; print get_hostname()"

	cmd := exec.Command(ddAgentPy, "-c", getHostnameCmd)
	dockerEnv := os.Getenv("DOCKER_DD_AGENT")
	cmd.Env = []string{
		"PYTHONPATH=/opt/datadog-agent/agent",
		fmt.Sprintf("DOCKER_DD_AGENT=%s", dockerEnv),
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Infof("error retrieving dd-agent hostname, falling back to os.Hostname(): %v", err)
		return os.Hostname()
	}

	hostname := strings.TrimSpace(stdout.String())

	if hostname == "" {
		log.Infof("error retrieving dd-agent hostname, falling back to os.Hostname(): %s", stderr.String())
		return os.Hostname()
	}

	return hostname, err
}

// getProxySettings returns a url.Url for the proxy configuration from datadog.conf, if available.
// In the case of invalid settings an error is logged and nil is returned. If settings are missing,
// meaning we don't want a proxy, then nil is returned with no error.
func getProxySettings(m *ini.Section) *url.URL {
	var host, scheme string
	if v := m.Key("proxy_host").MustString(""); v != "" {
		// accept either http://myproxy.com or myproxy.com
		if i := strings.Index(v, "://"); i != -1 {
			// when available, parse the scheme from the url
			scheme = v[0:i]
			host = v[i+3:]
		} else {
			host = v
		}
	}

	if host == "" {
		return nil
	}

	var port int
	if v := m.Key("proxy_port").MustInt(-1); v != -1 {
		port = v
	}
	var user, password string
	if v := m.Key("proxy_user").MustString(""); v != "" {
		user = v
	}
	if v := m.Key("proxy_password").MustString(""); v != "" {
		password = v
	}
	return constructProxy(host, scheme, port, user, password)
}

// proxyFromEnv parses out the proxy configuration from the ENV variables in a
// similar way to getProxySettings and, if enough values are available, returns
// a new proxy URL value. If the environment is not set for this then the
// `defaultVal` is returned.
func proxyFromEnv(defaultVal *url.URL) *url.URL {
	var host, scheme string
	if v := os.Getenv("PROXY_HOST"); v != "" {
		// accept either http://myproxy.com or myproxy.com
		if i := strings.Index(v, "://"); i != -1 {
			// when available, parse the scheme from the url
			scheme = v[0:i]
			host = v[i+3:]
		} else {
			host = v
		}
	}

	if host == "" {
		return defaultVal
	}

	var port int
	if v := os.Getenv("PROXY_PORT"); v != "" {
		port, _ = strconv.Atoi(v)
	}
	var user, password string
	if v := os.Getenv("PROXY_USER"); v != "" {
		user = v
	}
	if v := os.Getenv("PROXY_PASSWORD"); v != "" {
		password = v
	}

	return constructProxy(host, scheme, port, user, password)
}

// constructProxy constructs a *url.Url for a proxy given the parts of a
// Note that we assume we have at least a non-empty host for this call but
// all other values can be their defaults (empty string or 0).
func constructProxy(host, scheme string, port int, user, password string) *url.URL {
	var userpass *url.Userinfo
	if user != "" {
		if password != "" {
			userpass = url.UserPassword(user, password)
		} else {
			userpass = url.User(user)
		}
	}

	var path string
	if userpass != nil {
		path = fmt.Sprintf("%s://%s@%s:%v", scheme, userpass.String(), host, port)
	} else {
		path = fmt.Sprintf("%s://%s:%v", scheme, host, port)
	}

	u, err := url.Parse(path)
	if err != nil {
		log.Errorf("error parsing proxy settings, not using a proxy: %s", err)
		return nil
	}
	return u
}
