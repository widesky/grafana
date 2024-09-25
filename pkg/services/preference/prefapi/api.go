// shared logic between httpserver and teamapi
package prefapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/kinds/preferences"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/org"
	pref "github.com/grafana/grafana/pkg/services/preference"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
)

func UpdatePreferencesFor(ctx context.Context,
	dashboardService dashboards.DashboardService, preferenceService pref.Service,
	orgID, userID, teamId int64, dtoCmd *dtos.UpdatePrefsCmd) response.Response {
	if dtoCmd.Theme != "" && !pref.IsValidThemeID(dtoCmd.Theme) {
		return response.Error(http.StatusBadRequest, "Invalid theme", nil)
	}

	dashboardID := dtoCmd.HomeDashboardID
	if dtoCmd.HomeDashboardUID != nil {
		query := dashboards.GetDashboardQuery{UID: *dtoCmd.HomeDashboardUID, OrgID: orgID}
		if query.UID == "" {
			// clear the value
			dashboardID = 0
		} else {
			queryResult, err := dashboardService.GetDashboard(ctx, &query)
			if err != nil {
				return response.Error(http.StatusNotFound, "Dashboard not found", err)
			}
			dashboardID = queryResult.ID
		}
	}
	dtoCmd.HomeDashboardID = dashboardID

	saveCmd := pref.SavePreferenceCommand{
		UserID:            userID,
		OrgID:             orgID,
		TeamID:            teamId,
		Theme:             dtoCmd.Theme,
		Language:          dtoCmd.Language,
		Timezone:          dtoCmd.Timezone,
		WeekStart:         dtoCmd.WeekStart,
		HomeDashboardID:   dtoCmd.HomeDashboardID,
		QueryHistory:      dtoCmd.QueryHistory,
		CookiePreferences: dtoCmd.Cookies,
	}

	if err := preferenceService.Save(ctx, &saveCmd); err != nil {
		return response.ErrOrFallback(http.StatusInternalServerError, "Failed to save preferences", err)
	}

	return response.Success("Preferences updated")
}

func UpdateUsersOrgByLanguage(ctx context.Context, userService user.Service, preferenceService pref.Service, cfg *setting.Cfg, orgsJson []org.UserOrgDTO, orgID, userID int64, dtoCmd *dtos.UpdatePrefsCmd) {
	saveCmd := pref.SavePreferenceCommand{
		UserID:            userID,
		OrgID:             orgID,
		TeamID:            0,
		Theme:             dtoCmd.Theme,
		Language:          dtoCmd.Language,
		Timezone:          dtoCmd.Timezone,
		WeekStart:         dtoCmd.WeekStart,
		HomeDashboardID:   dtoCmd.HomeDashboardID,
		QueryHistory:      dtoCmd.QueryHistory,
		CookiePreferences: dtoCmd.Cookies,
	}

	// remove suffix from current org name if present
	currentOrg := findOrgById(orgsJson, orgID)

	if index := indexOf(currentOrg.Name, "_"); index != -1 {
		currentOrg.Name = currentOrg.Name[:index]
	}

	// construct expected org name
	var orgName string
	if orgSuffix, hasLang := cfg.WideSkyProvisioner.LanguagePermissions.Languages[dtoCmd.Language]; !hasLang {
		// this language assignment doesn't exist
		// no organisation id changes needed
		cfg.Logger.Info("Language %s is not in the list of orgAssignments\n", dtoCmd.Language)
		cfg.Logger.Info("%#v\n", cfg.WideSkyProvisioner.LanguagePermissions.Languages)
		return
	} else if orgSuffix == "" {
		orgName = currentOrg.Name
	} else {
		orgName = fmt.Sprintf("%s_%s", currentOrg.Name, orgSuffix)
	}

	foundOrg := findOrgByName(orgsJson, orgName)

	if foundOrg.OrgID == -1 {
		// not found, no organisation id changes needed
		return
	}

	// Set user's organisation
	cmd := user.SetUsingOrgCommand{UserID: userID, OrgID: foundOrg.OrgID}
	if err := userService.SetUsingOrg(ctx, &cmd); err != nil {
		cfg.Logger.Error("Failed to change the organisation by language preference")
		return
	}

	// Update organisation id for this command as well
	saveCmd.OrgID = foundOrg.OrgID
	if err := preferenceService.Save(ctx, &saveCmd); err != nil {
		cfg.Logger.Error("Failed to save preferenes")
		return
	}
}

func GetPreferencesFor(ctx context.Context,
	dashboardService dashboards.DashboardService, preferenceService pref.Service,
	orgID, userID, teamID int64) response.Response {
	prefsQuery := pref.GetPreferenceQuery{UserID: userID, OrgID: orgID, TeamID: teamID}

	preference, err := preferenceService.Get(ctx, &prefsQuery)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to get preferences", err)
	}

	var dashboardUID string

	// when homedashboardID is 0, that means it is the default home dashboard, no UID would be returned in the response
	if preference.HomeDashboardID != 0 {
		query := dashboards.GetDashboardQuery{ID: preference.HomeDashboardID, OrgID: orgID}
		queryResult, err := dashboardService.GetDashboard(ctx, &query)
		if err == nil {
			dashboardUID = queryResult.UID
		}
	}

	dto := preferences.Spec{}

	if preference.WeekStart != nil && *preference.WeekStart != "" {
		dto.WeekStart = preference.WeekStart
	}
	if preference.Theme != "" {
		dto.Theme = &preference.Theme
	}
	if dashboardUID != "" {
		dto.HomeDashboardUID = &dashboardUID
	}
	if preference.Timezone != "" {
		dto.Timezone = &preference.Timezone
	}

	if preference.JSONData != nil {
		if preference.JSONData.Language != "" {
			dto.Language = &preference.JSONData.Language
		}

		if preference.JSONData.QueryHistory.HomeTab != "" {
			dto.QueryHistory = &preferences.QueryHistoryPreference{
				HomeTab: &preference.JSONData.QueryHistory.HomeTab,
			}
		}
	}

	return response.JSON(http.StatusOK, &dto)
}

func findOrgByName(orgs []org.UserOrgDTO, toFind string) org.UserOrgDTO {
	for _, org := range orgs {
		if org.Name == toFind {
			return org
		}
	}

	return org.UserOrgDTO{
		Name:  "",
		OrgID: -1,
		Role:  "",
	}
}

func findOrgById(orgs []org.UserOrgDTO, id int64) org.UserOrgDTO {
	for _, org := range orgs {
		if org.OrgID == id {
			return org
		}
	}

	return org.UserOrgDTO{
		Name:  "",
		OrgID: -1,
		Role:  "",
	}
}

func indexOf(str string, toMatch string) int {
	for i := 0; i < len(str); i++ {
		if subset := str[i : i+len(toMatch)]; subset == toMatch {
			return i
		}
	}

	return -1
}
