package setting

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/grafana/pkg/infra/log"
)

type WideSkyProvisionerSettings struct {
	PrimaryOrg int64             `json:"primaryOrg"`
	Languages  map[string]string `json:"languages"`
}

func (cfg *Cfg) readWideSkyProvisioner() WideSkyProvisionerSettings {
	if !cfg.Raw.HasSection("widesky_app_provisioner") {
		return WideSkyProvisionerSettings{
			PrimaryOrg: -1,
			Languages:  make(map[string]string),
		}
	}

	wsProvisionerSec := cfg.Raw.Section("widesky_app_provisioner")

	return WideSkyProvisionerSettings{
		PrimaryOrg: wsProvisionerSec.Key("primary_org").MustInt64(2),
		Languages:  readLanguageAssignments(cfg),
	}
}

func count(match string, str string) int {
	count := 0
	for i := 0; i < len(str); i++ {
		subset := str[i : i+len(match)]
		if subset == match {
			count++
		}
	}

	return count
}

func includesStr(array []string, toFind string) bool {
	for _, element := range array {
		if element == toFind {
			return true
		}
	}

	return false
}

func isAlphaNumeric(str string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(str)
}

func addAssignment(item string, codes *map[string]string, suffixAssignment *map[string]string, logger log.Logger) bool {
	if item == "" || count(":", item) != 1 {
		logger.Error("Missing content or incorrect usage of ':'", "item", item)
		return false
	}

	item = strings.Trim(item, " ")

	langCodes := []string{"en-US", "fr-FR", "es-ES", "de-DE", "zh-Hans"}
	split := strings.Split(item, ":")
	langCode, orgSuffix := strings.Trim(split[0], " "), strings.Trim(split[1], " ")
	if dupLang := (*suffixAssignment)[orgSuffix]; dupLang != "" {
		// usage of organisation suffixs must be unique
		logger.Error(
			fmt.Sprintf("Language %s has already been assigned suffix %s. Duplicate suffixes are not allowed.",
				dupLang, orgSuffix))
		return false
	}

	if langCode == "" {
		// assign default value english
		langCode = "en-US"
	}

	if !includesStr(langCodes, langCode) {
		// Is it a valid language code from the set?
		logger.Error(
			fmt.Sprintf("Invalid language code %s. Must be one of %s", langCode, strings.Join(langCodes, ", ")))
		return false
	}

	if (*codes)[langCode] != "" {
		// usage of language codes must be unique
		logger.Error(fmt.Sprintf("Language code %s is already present.", langCode))
		return false
	}

	if !isAlphaNumeric(orgSuffix) {
		logger.Error("Organisation suffix must comprise of alphanumeric character")
		return false
	}

	// all checks completed
	(*suffixAssignment)[orgSuffix] = langCode
	(*codes)[langCode] = orgSuffix

	return true
}

func readLanguageAssignments(cfg *Cfg) map[string]string {
	if !cfg.Raw.HasSection("ws_internationalization") {
		return make(map[string]string)
	}

	wsInternationalizationSec := cfg.Raw.Section("ws_internationalization")
	assignments := wsInternationalizationSec.Key("org_language_assignments").MustString("en-US:english")

	languages := make(map[string]string)
	suffixAssignments := make(map[string]string)

	split := strings.Split(strings.Trim(assignments, " "), ",")

	if len(split) == 1 && split[0] == "" {
		return languages
	}

	for _, item := range split {
		if success := addAssignment(item, &languages, &suffixAssignments, cfg.Logger); !success {
			return languages
		}
	}

	return languages
}
