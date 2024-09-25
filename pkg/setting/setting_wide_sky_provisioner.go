package setting

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/grafana/pkg/infra/log"
)

type TeamPermissions struct {
	Enabled    bool                           `json:"enabled"`
	Plugins    map[string]map[string][]string `json:"plugins"`
	Dashboards map[string][]string            `json:"dashboards"`
}

type LanguagePermissions struct {
	Enabled   bool              `json:"enabled"`
	Languages map[string]string `json:"languages"`
}

type WideSkyProvisionerSettings struct {
	TeamPermissions     TeamPermissions     `json:"teamPermissions"`
	LanguagePermissions LanguagePermissions `json:"languagePermissions"`
}

func (cfg *Cfg) readWideSkyProvisioner() {
	teamPermissions := readWideSkyTeamPermissions(cfg)
	cfg.Logger.Info("WideSky team permissions", "Enabled", teamPermissions.Enabled)

	languagePermissions := readWideSkyLanguagePermissions(cfg)
	cfg.Logger.Info("WideSky language permission", "Enabled", languagePermissions.Enabled)
	if languagePermissions.Enabled {
		cfg.Logger.Info(fmt.Sprintf("WideSky organisation has %d language assignments", len(languagePermissions.Languages)))
	}

	cfg.WideSkyProvisioner = &WideSkyProvisionerSettings{
		TeamPermissions:     teamPermissions,
		LanguagePermissions: languagePermissions,
	}
}

func readWideSkyTeamPermissions(cfg *Cfg) TeamPermissions {
	dashboards := make(map[string][]string)
	plugins := make(map[string]map[string][]string)

	if !cfg.Raw.HasSection("ws_security") {
		return TeamPermissions{
			Dashboards: dashboards,
			Plugins:    plugins,
			Enabled:    false,
		}
	}

	teamPermissionsSec := cfg.Raw.Section("ws_security")

	return TeamPermissions{
		Dashboards: dashboards,
		Plugins:    plugins,
		Enabled:    teamPermissionsSec.Key("enable_team_app_plugin_permission").MustBool(true),
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

func readWideSkyLanguagePermissions(cfg *Cfg) LanguagePermissions {
	emptyPermissions := LanguagePermissions{
		Languages: make(map[string]string),
		Enabled:   false,
	}

	if !cfg.Raw.HasSection("ws_internationalization") {
		return emptyPermissions
	}

	languagePermissionsSec := cfg.Raw.Section("ws_internationalization")
	enabled := languagePermissionsSec.Key("enabled").MustBool(false)
	assignmentStr := languagePermissionsSec.Key("org_language_assignments").String()

	if !enabled {
		return emptyPermissions
	}

	languages := make(map[string]string)
	suffixAssignments := make(map[string]string)

	split := strings.Split(strings.Trim(assignmentStr, " "), ",")

	if len(split) == 1 && split[0] == "" {
		return emptyPermissions
	}

	for _, item := range split {
		if item == "" || count(":", item) != 1 {
			cfg.Logger.Error("Missing content or incorrect usage of ':'", "item", item)
			return emptyPermissions
		}

		if success := addAssignment(item, &languages, &suffixAssignments, cfg.Logger); !success {
			return emptyPermissions
		}
	}

	return LanguagePermissions{
		Languages: languages,
		Enabled:   len(languages) > 0,
	}
}
