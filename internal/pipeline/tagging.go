package pipeline

import (
	"strings"

	"github.com/iamlucif3r/trikal/internal/models"
)

var keywords = []string{
	"exploit", "zero day", "vulnerability", "rce", "breach", "malware",
	"phishing", "ransomware", "apt", "critical", "cve", "patch", "cyber attack",
	"rootkit", "leak", "ddos", "security advisory", "mitre", "cisa", "cert",
	"backdoor", "trojan", "worm", "botnet", "sql injection", "xss", "csrf",
	"privilege escalation", "data breach", "infostealer", "spyware", "adware",
	"social engineering", "threat actor", "ioc", "tactics", "techniques", "procedures",
	"zero trust", "firewall", "endpoint", "siem", "soc", "incident response",
	"threat intelligence", "forensics", "encryption", "dos", "supply chain attack",
	"insider threat", "password", "credential stuffing", "brute force", "mitm",
	"drive-by download", "watering hole", "sandbox", "obfuscation", "payload",
	"command and control", "c2", "red team", "blue team", "purple team",
	"security update", "exploit kit", "vulnerability management", "bug bounty",
	"security flaw", "security hole", "security patch", "security risk",
	"security incident", "security breach", "security vulnerability",
}

func Tag(in []models.NewsItem) []models.NewsItem {
	for i := range in {
		text := strings.ToLower(in[i].Title + " " + in[i].Summary)
		seen := map[string]struct{}{}
		for _, k := range keywords {
			if strings.Contains(text, strings.ToLower(k)) {
				seen[strings.ToLower(k)] = struct{}{}
			}
		}
		in[i].Tags = in[i].Tags[:0]
		for k := range seen {
			in[i].Tags = append(in[i].Tags, k)
		}
	}
	return in
}
