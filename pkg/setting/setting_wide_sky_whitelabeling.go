package setting

import (
	"sort"
	"strconv"
	"strings"
)

type FooterItem struct {
	Text  string `json:"text"`
	Url   string `json:"url"`
	Icon  string `json:"icon"`
	Hide  bool   `json:"hide"`
	Color string `json:"color"`
}

type HelpOption struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}

type LoginBox struct {
	Logo                        string `json:"logo"`
	LogoMaxWidth                int    `json:"logoMaxWidth"`
	LogoMaxWidthMediaBreakpoint int    `json:"logoMaxWidthMediaBreakpoint"`

	Color                   string `json:"color"`
	TextColor               string `json:"textColor"`
	TextboxBorderColor      string `json:"textboxBorderColor"`
	TextboxPlaceholderColor string `json:"textboxPlaceholderColor"`
	TextboxBackgroundColor  string `json:"textboxBackgroundColor"`

	ButtonBgColor        string `json:"buttonBgColor"`
	ButtonHoverBgColor   string `json:"buttonHoverBgColor"`
	ButtonTextColor      string `json:"buttonTextColor"`
	ButtonHoverTextColor string `json:"buttonHoverTextColor"`

	LinkButtonBgColor        string `json:"linkButtonBgColor"`
	LinkButtonHoverBgColor   string `json:"linkButtonHoverBgColor"`
	LinkButtonTextColor      string `json:"linkButtonTextColor"`
	LinkButtonHoverTextColor string `json:"linkButtonHoverTextColor"`
}

type LoginBackground struct {
	Image      string `json:"image"`
	Background string `json:"background"`
	Position   string `json:"position"`
	Color      string `json:"color"`
	MinHeight  string `json:"minHeight"`
}

type WideSkyWhitelabelingSettings struct {
	AppSidebarLogo            string `json:"appSidebarLogo"`
	ApplicationName           string `json:"applicationName"`
	ApplicationLogo           string `json:"applicationLogo"`
	LoadingLogo               string `json:"loadingLogo"`
	BrowserTabFavicon         string
	BrowserTabTitle           string `json:"browserTabTitle"`
	BrowserTabSubtitle        string `json:"browserTabSubtitle"`
	ShowEnterpriseUpgradeInfo bool   `json:"showEnterpriseUpgradeInfo"`
	AdminSettingsMessage      string `json:"adminSettingsMessage"`

	FooterItems     []FooterItem `json:"footerItems"`
	FooterPipeColor string       `json:"footerPipeColor"`

	HelpOptions []HelpOption `json:"helpOptions"`

	LoginBox        LoginBox        `json:"loginBox"`
	LoginBackground LoginBackground `json:"loginBackground"`

	EntityNotFoundText       string `json:"entityNotFoundText"`
	EntityNotFoundLink       string `json:"entityNotFoundLink"`
	EntityNotFoundLinkText   string `json:"entityNotFoundLinkText"`
	TemplateVariableHelpLink string `json:"templateVariableHelpLink"`
}

func (cfg *Cfg) readWideSkyWhitelabeling() {
	hasWhitelabelSec := cfg.Raw.HasSection("ws_whitelabel")

	if !hasWhitelabelSec {
		return
	}

	whitelabelSec := cfg.Raw.Section("ws_whitelabel")

	linkMap := make(map[string]*FooterItem)
	for _, key := range whitelabelSec.Keys() {
		linkParts := strings.Split(key.Name(), "_")
		if len(linkParts) < 4 || linkParts[0] != "footer" || linkParts[1] != "element" {
			continue
		}

		index := linkParts[2]
		if _, exists := linkMap[index]; !exists {
			linkMap[index] = &FooterItem{}
		}

		switch linkParts[3] {
		case "text":
			linkMap[index].Text = key.String()
		case "link":
			linkMap[index].Url = key.String()
		case "icon":
			linkMap[index].Icon = key.String()
		case "hide":
			linkMap[index].Hide = key.MustBool(false)
		case "color":
			linkMap[index].Color = key.String()
		}
	}

	var footerItems []FooterItem
	{
		keys := make([]string, 0, len(linkMap))
		for k := range linkMap {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			iNum, _ := strconv.Atoi(keys[i])
			jNum, _ := strconv.Atoi(keys[j])
			return iNum < jNum
		})

		for _, k := range keys {
			footerItems = append(footerItems, *linkMap[k])
		}
	}

	helpMap := make(map[string]*HelpOption)
	for _, key := range whitelabelSec.Keys() {
		helpParts := strings.Split(key.Name(), "_")
		if len(helpParts) < 4 || helpParts[0] != "help" || helpParts[1] != "option" {
			continue
		}

		index := helpParts[2]
		if _, exists := helpMap[index]; !exists {
			helpMap[index] = &HelpOption{}
		}

		switch helpParts[3] {
		case "label":
			helpMap[index].Label = key.String()
		case "href":
			helpMap[index].Href = key.String()
		}
	}

	var helpOptions []HelpOption
	{
		keys := make([]string, 0, len(helpMap))
		for k := range helpMap {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			iNum, _ := strconv.Atoi(keys[i])
			jNum, _ := strconv.Atoi(keys[j])
			return iNum < jNum
		})

		for _, k := range keys {
			helpOptions = append(helpOptions, *helpMap[k])
		}
	}

	cfg.WideSkyWhitelabeling = &WideSkyWhitelabelingSettings{
		AppSidebarLogo:            whitelabelSec.Key("app_sidebar_logo").String(),
		ApplicationName:           whitelabelSec.Key("application_name").String(),
		ApplicationLogo:           whitelabelSec.Key("application_logo").String(),
		LoadingLogo:               whitelabelSec.Key("loading_logo").String(),
		BrowserTabFavicon:         whitelabelSec.Key("browser_tab_favicon").String(),
		BrowserTabTitle:           whitelabelSec.Key("browser_tab_title").String(),
		BrowserTabSubtitle:        whitelabelSec.Key("browser_tab_sub_title").String(),
		ShowEnterpriseUpgradeInfo: whitelabelSec.Key("show_enterprise_upgrade_info").MustBool(false),
		AdminSettingsMessage:      whitelabelSec.Key("admin_settings_message").String(),

		FooterItems:     footerItems,
		FooterPipeColor: whitelabelSec.Key("footer_pipe_color").String(),

		HelpOptions: helpOptions,

		LoginBox: LoginBox{
			Logo:                        whitelabelSec.Key("login_box_logo").String(),
			LogoMaxWidth:                whitelabelSec.Key("login_box_logo_max_width").MustInt(-1),
			LogoMaxWidthMediaBreakpoint: whitelabelSec.Key("login_box_logo_max_width_media_breakpoint").MustInt(-1),

			Color:                   whitelabelSec.Key("login_box_colour").String(),
			TextColor:               whitelabelSec.Key("login_box_text_colour").String(),
			TextboxBorderColor:      whitelabelSec.Key("login_box_textbox_border_colour").String(),
			TextboxPlaceholderColor: whitelabelSec.Key("login_box_textbox_placeholder_colour").String(),
			TextboxBackgroundColor:  whitelabelSec.Key("login_box_textbox_background_colour").String(),

			ButtonBgColor:        whitelabelSec.Key("login_box_button_bg_colour").String(),
			ButtonHoverBgColor:   whitelabelSec.Key("login_box_button_hover_bg_colour").String(),
			ButtonTextColor:      whitelabelSec.Key("login_box_button_text_colour").String(),
			ButtonHoverTextColor: whitelabelSec.Key("login_box_button_hover_text_colour").String(),

			LinkButtonBgColor:        whitelabelSec.Key("login_box_link_button_bg_colour").String(),
			LinkButtonHoverBgColor:   whitelabelSec.Key("login_box_link_button_hover_bg_colour").String(),
			LinkButtonTextColor:      whitelabelSec.Key("login_box_link_button_text_colour").String(),
			LinkButtonHoverTextColor: whitelabelSec.Key("login_box_link_button_hover_text_colour").String(),
		},
		LoginBackground: LoginBackground{
			Image:      whitelabelSec.Key("login_background_image").String(),
			Background: whitelabelSec.Key("login_background_background").String(),
			Position:   whitelabelSec.Key("login_background_position").String(),
			Color:      whitelabelSec.Key("login_background_color").String(),
			MinHeight:  whitelabelSec.Key("login_background_min_height").String(),
		},

		EntityNotFoundText:       whitelabelSec.Key("entity_not_found_text").String(),
		EntityNotFoundLink:       whitelabelSec.Key("entity_not_found_link").String(),
		EntityNotFoundLinkText:   whitelabelSec.Key("entity_not_found_link_text").String(),
		TemplateVariableHelpLink: whitelabelSec.Key("template_variable_help_link").String(),
	}
}
