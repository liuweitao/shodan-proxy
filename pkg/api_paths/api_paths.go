package api_paths

import "regexp"

var AllowedPaths = []*regexp.Regexp{
    // Search Methods
    regexp.MustCompile(`^/shodan/host/\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`),
    regexp.MustCompile(`^/shodan/host/count$`),
    regexp.MustCompile(`^/shodan/host/search$`),
    regexp.MustCompile(`^/shodan/host/search/facets$`),
    regexp.MustCompile(`^/shodan/host/search/filters$`),
    regexp.MustCompile(`^/shodan/host/search/tokens$`),

    // On-Demand Scanning
    regexp.MustCompile(`^/shodan/ports$`),
    regexp.MustCompile(`^/shodan/protocols$`),
    regexp.MustCompile(`^/shodan/scan$`),
    regexp.MustCompile(`^/shodan/scan/internet$`),
    regexp.MustCompile(`^/shodan/scans$`),
    regexp.MustCompile(`^/shodan/scan/[a-zA-Z0-9]+$`),

    // Network Alerts
    regexp.MustCompile(`^/shodan/alert$`),
    regexp.MustCompile(`^/shodan/alert/[a-zA-Z0-9]+/info$`),
    regexp.MustCompile(`^/shodan/alert/[a-zA-Z0-9]+$`),
    regexp.MustCompile(`^/shodan/alert/info$`),
    regexp.MustCompile(`^/shodan/alert/triggers$`),
    regexp.MustCompile(`^/shodan/alert/[a-zA-Z0-9]+/trigger/[a-zA-Z0-9]+$`),
    regexp.MustCompile(`^/shodan/alert/[a-zA-Z0-9]+/trigger/[a-zA-Z0-9]+/ignore/[a-zA-Z0-9]+$`),
    regexp.MustCompile(`^/shodan/alert/[a-zA-Z0-9]+/notifier/[a-zA-Z0-9]+$`),

    // Notifiers
    regexp.MustCompile(`^/notifier$`),
    regexp.MustCompile(`^/notifier/provider$`),
    regexp.MustCompile(`^/notifier/[a-zA-Z0-9]+$`),

    // Directory Methods
    regexp.MustCompile(`^/shodan/query$`),
    regexp.MustCompile(`^/shodan/query/search$`),
    regexp.MustCompile(`^/shodan/query/tags$`),

    // Bulk Data (Enterprise)
    regexp.MustCompile(`^/shodan/data$`),
    regexp.MustCompile(`^/shodan/data/[a-zA-Z0-9]+$`),

    // Manage Organization (Enterprise)
    regexp.MustCompile(`^/org$`),
    regexp.MustCompile(`^/org/member/[a-zA-Z0-9]+$`),

    // Account Methods
    regexp.MustCompile(`^/account/profile$`),

    // DNS Methods
    regexp.MustCompile(`^/dns/domain/[a-zA-Z0-9.-]+$`),
    regexp.MustCompile(`^/dns/resolve$`),
    regexp.MustCompile(`^/dns/reverse$`),

    // Utility Methods
    regexp.MustCompile(`^/tools/httpheaders$`),
    regexp.MustCompile(`^/tools/myip$`),

    // API Status Methods
    regexp.MustCompile(`^/api-info$`),
}
