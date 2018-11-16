package config

// Prefixes
const procPrefix = "process_config."
const networkPrefix = "network_tracer_config."
const intervalsPrefix = procPrefix + "intervals."
const windowsPrefix = procPrefix + "windows."

// List of the keys used in the yaml configuration and environment
// kXXX -> YAML
// envXXX -> Environment
const (
	// Top level
	kDDURL   = "dd_url"
	kAPIKey  = "api_key"
	envDDURL = "DD_PROCESS_AGENT_URL"

	// All keys in the Yaml file
	// A string indicate the enabled state of the Agent.
	// If "false" (the default) we will only collect containers.
	// If "true" we will collect containers and processes.
	// If "disabled" the agent will be disabled altogether and won't start.
	kEnabled   = procPrefix + "enabled"
	envEnabled = "DD_PROCESS_AGENT_ENABLED"

	// The full path to the file where process-agent logs will be written.
	kLogFile = procPrefix + "log_file"

	// The interval, in seconds, at which we will run each check. If you want consistent
	// behavior between real-time you may set the Container/ProcessRT intervals to 10.
	kIntervalsContainer   = intervalsPrefix + "container"
	kIntervalsContainerRT = intervalsPrefix + "container_realtime"
	kIntervalsProcess     = intervalsPrefix + "process"
	kIntervalsProcessRT   = intervalsPrefix + "process_realtime"
	kIntervalsConnections = intervalsPrefix + "connections"

	// A list of regex patterns that will exclude a process if matched.
	kBlacklistPatterns = procPrefix + "blacklist_patterns"
	// Enable/Disable the DataScrubber to obfuscate process args
	kScrubArgs   = procPrefix + "scrub_args"
	envScrubArgs = "DD_SCRUB_ARGS"

	// A custom word list to enhance the default one used by the DataScrubber
	kCustomSensitiveWords   = procPrefix + "custom_sensitive_words"
	envCustomSensitiveWords = "DD_CUSTOM_SENSITIVE_WORDS"

	// Strips all process arguments
	kStripProcessArguments   = procPrefix + "strip_proc_arguments"
	envStripProcessArguments = "DD_STRIP_PROCESS_ARGS"

	// How many check results to buffer in memory when POST fails. The default is usually fine.
	kQueueSize = procPrefix + "queue_size"
	// The maximum number of file descriptors to open when collecting net connections.
	// Only change if you are running out of file descriptors from the Agent.
	kMaxProcFDs = procPrefix + "max_proc_fds"
	// The maximum number of processes, connections or containers per message.
	// Only change if the defaults are causing issues.
	kMaxPerMessage = procPrefix + "max_per_message"
	// Overrides the path to the Agent bin used for getting the hostname. The default is usually fine.
	kDDAgentBin   = procPrefix + "dd_agent_bin"
	envDDAgentBin = "DD_AGENT_BIN"

	// Overrides of the environment we pass to fetch the hostname. The default is usually fine.
	kDDAgentEnv   = procPrefix + "dd_agent_env"
	envDDAgentEnv = "DD_AGENT_ENV"

	// Optional additional pairs of endpoint_url => []apiKeys to submit to other locations.
	kAdditionalEndpoints = procPrefix + "additional_endpoints"

	// Windows config

	// Sets windows process table refresh rate (in number of check runs)
	kWinArgsRefreshInterval = windowsPrefix + "args_refresh_interval"
	// Controls getting process arguments immediately when a new process is discovered
	kWinAddNewArgs = windowsPrefix + "add_new_args"

	// Network tracer config
	// A string indicating the enabled state of the network tracer.
	kNetworkTracingEnabled   = networkPrefix + "enabled"
	envNetworkTracingEnabled = "DD_NETWORK_TRACING_ENABLED"

	// The full path to the location of the unix socket where network traces will be accessed
	kNetworkUnixSocketPath   = networkPrefix + "nettracer_socket"
	envNetworkUnixSocketPath = "DD_NETTRACER_SOCKET"

	// The full path to the file where network-tracer logs will be written.
	kNetworkLogFile = networkPrefix + "log_file"
)
