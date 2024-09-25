package setting

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
)

func Test_readWideSkyWhitelabeling(t *testing.T) {
	testCases := []struct {
		desc              string
		whitelabelOptions map[string]string
		verifyCfg         func(*testing.T, Cfg)
	}{
		{
			desc: "should apply the correct general settings",
			whitelabelOptions: map[string]string{
				"app_sidebar_logo":             "public/img/widesky_icon.svg",
				"application_name":             "WideSky",
				"application_logo":             "public/img/widesky_icon.svg",
				"loading_logo":                 "public/img/widesky_icon.svg",
				"browser_tab_favicon":          "public/img/widesky_icon.svg",
				"browser_tab_title":            "WideSky",
				"browser_tab_sub_title":        "",
				"show_enterprise_upgrade_info": "false",
				"admin_settings_message":       "These system settings are managed by the <a className=\"external-link\" href=\"https://widesky.cloud/contact-us/\" rel=\"noreferrer\" target=\"_blank\">WideSky</a> team. Contact us to request changes.",
			},
			verifyCfg: func(t *testing.T, cfg Cfg) {
				require.Equal(t, "public/img/widesky_icon.svg", cfg.WideSkyWhitelabeling.AppSidebarLogo)
				require.Equal(t, "WideSky", cfg.WideSkyWhitelabeling.ApplicationName)
				require.Equal(t, "public/img/widesky_icon.svg", cfg.WideSkyWhitelabeling.ApplicationLogo)
				require.Equal(t, "public/img/widesky_icon.svg", cfg.WideSkyWhitelabeling.LoadingLogo)
				require.Equal(t, "public/img/widesky_icon.svg", cfg.WideSkyWhitelabeling.BrowserTabFavicon)
				require.Equal(t, "WideSky", cfg.WideSkyWhitelabeling.BrowserTabTitle)
				require.Equal(t, "", cfg.WideSkyWhitelabeling.BrowserTabSubtitle)
				require.Equal(t, false, cfg.WideSkyWhitelabeling.ShowEnterpriseUpgradeInfo)
				require.Equal(t, "These system settings are managed by the <a className=\"external-link\" href=\"https://widesky.cloud/contact-us/\" rel=\"noreferrer\" target=\"_blank\">WideSky</a> team. Contact us to request changes.", cfg.WideSkyWhitelabeling.AdminSettingsMessage)
			},
		},
		{
			desc: "should apply the correct help settings",
			whitelabelOptions: map[string]string{
				"help_option_1_label": "Documentation",
				"help_option_1_href":  "https://widesky.cloud/docs/",
			},
			verifyCfg: func(t *testing.T, cfg Cfg) {
				require.Equal(t, HelpOption{Label: "Documentation", Href: "https://widesky.cloud/docs/"}, cfg.WideSkyWhitelabeling.HelpOptions[0])
			},
		},
		{
			desc: "should apply the correct footer settings",
			whitelabelOptions: map[string]string{
				"footer_element_1_hide":  "false",
				"footer_element_2_hide":  "true",
				"footer_element_3_hide":  "false",
				"footer_element_1_icon":  "document-info",
				"footer_element_2_icon":  "envelope",
				"footer_element_3_icon":  "",
				"footer_element_1_text":  "Documentation",
				"footer_element_2_text":  "Support",
				"footer_element_3_text":  "Open Source",
				"footer_element_1_link":  "https://widesky.cloud/docs/",
				"footer_element_2_link":  "https://widesky.cloud/contact-us/",
				"footer_element_3_link":  "https://github.com/widesky/grafana",
				"footer_element_1_color": "#000000",
				"footer_element_2_color": "#111111",
				"footer_element_3_color": "#222222",
				"footer_pipe_color":      "#ccccdc",
			},
			verifyCfg: func(t *testing.T, cfg Cfg) {
				require.Equal(t, FooterItem{Hide: false, Icon: "document-info", Text: "Documentation", Url: "https://widesky.cloud/docs/", Color: "#000000"}, cfg.WideSkyWhitelabeling.FooterItems[0])
				require.Equal(t, FooterItem{Hide: true, Icon: "envelope", Text: "Support", Url: "https://widesky.cloud/contact-us/", Color: "#111111"}, cfg.WideSkyWhitelabeling.FooterItems[1])
				require.Equal(t, FooterItem{Hide: false, Icon: "", Text: "Open Source", Url: "https://github.com/widesky/grafana", Color: "#222222"}, cfg.WideSkyWhitelabeling.FooterItems[2])
				require.Equal(t, "#ccccdc", cfg.WideSkyWhitelabeling.FooterPipeColor)
			},
		},
		{
			desc: "should apply the correct login box settings",
			whitelabelOptions: map[string]string{
				"login_box_logo":                            "public/img/WideSky_logo_login.svg",
				"login_box_logo_max_width":                  "60",
				"login_box_logo_max_width_media_breakpoint": "200",
				"login_box_colour":                          "#fff",
				"login_box_text_colour":                     "#04275f",
				"login_box_textbox_placeholder_colour":      "",
				"login_box_textbox_border_colour":           "",
				"login_box_textbox_background_colour":           "",
				"login_box_button_bg_colour":                "#2688c6",
				"login_box_button_hover_bg_colour":          "#2688c6",
				"login_box_button_text_colour":              "#fff",
				"login_box_button_hover_text_colour":        "#fff",
				"login_box_link_button_bg_colour":           "0000",
				"login_box_link_button_hover_bg_colour":     "#2688c614",
				"login_box_link_button_text_colour":         "#04275f",
				"login_box_link_button_hover_text_colour":   "#04275f",
			},
			verifyCfg: func(t *testing.T, cfg Cfg) {
				require.Equal(t, "public/img/WideSky_logo_login.svg", cfg.WideSkyWhitelabeling.LoginBox.Logo)
				require.Equal(t, 60, cfg.WideSkyWhitelabeling.LoginBox.LogoMaxWidth)
				require.Equal(t, 200, cfg.WideSkyWhitelabeling.LoginBox.LogoMaxWidthMediaBreakpoint)
				require.Equal(t, "#fff", cfg.WideSkyWhitelabeling.LoginBox.Color)
				require.Equal(t, "#04275f", cfg.WideSkyWhitelabeling.LoginBox.TextColor)
				require.Equal(t, "", cfg.WideSkyWhitelabeling.LoginBox.TextboxPlaceholderColor)
				require.Equal(t, "", cfg.WideSkyWhitelabeling.LoginBox.TextboxBorderColor)
				require.Equal(t, "", cfg.WideSkyWhitelabeling.LoginBox.TextboxBackgroundColor)
				require.Equal(t, "#2688c6", cfg.WideSkyWhitelabeling.LoginBox.ButtonBgColor)
				require.Equal(t, "#2688c6", cfg.WideSkyWhitelabeling.LoginBox.ButtonHoverBgColor)
				require.Equal(t, "#fff", cfg.WideSkyWhitelabeling.LoginBox.ButtonTextColor)
				require.Equal(t, "#fff", cfg.WideSkyWhitelabeling.LoginBox.ButtonHoverTextColor)
				require.Equal(t, "0000", cfg.WideSkyWhitelabeling.LoginBox.LinkButtonBgColor)
				require.Equal(t, "#2688c614", cfg.WideSkyWhitelabeling.LoginBox.LinkButtonHoverBgColor)
				require.Equal(t, "#04275f", cfg.WideSkyWhitelabeling.LoginBox.LinkButtonTextColor)
				require.Equal(t, "#04275f", cfg.WideSkyWhitelabeling.LoginBox.LinkButtonHoverTextColor)
			},
		},
		{
			desc: "should apply the correct login background settings",
			whitelabelOptions: map[string]string{
				"login_background_image":      "public/img/WideskyHeroBackground.svg",
				"login_background_background": "transparent",
				"login_background_position":   "bottom",
				"login_background_color":      "transparent radial-gradient(closest-side at 50\\% 50\\%,#0073bc 0%,#003052 100%) 0\\% 0\\% no-repeat padding-box",
				"login_background_min_height": "100vh",
			},
			verifyCfg: func(t *testing.T, cfg Cfg) {
				require.Equal(t, "public/img/WideskyHeroBackground.svg", cfg.WideSkyWhitelabeling.LoginBackground.Image)
				require.Equal(t, "transparent", cfg.WideSkyWhitelabeling.LoginBackground.Background)
				require.Equal(t, "bottom", cfg.WideSkyWhitelabeling.LoginBackground.Position)
				require.Equal(t, "transparent radial-gradient(closest-side at 50\\% 50\\%,#0073bc 0%,#003052 100%) 0\\% 0\\% no-repeat padding-box", cfg.WideSkyWhitelabeling.LoginBackground.Color)
				require.Equal(t, "100vh", cfg.WideSkyWhitelabeling.LoginBackground.MinHeight)
			},
		},
		{
			desc: "should apply the correct entity not found & template help link settings",
			whitelabelOptions: map[string]string{
				"entity_not_found_text":       "seeking help from your WideSky administrator",
				"entity_not_found_link":       "mailto: support@widesky.cloud",
				"entity_not_found_link_text":  "Contact Support.",
				"template_variable_help_link": "https://widesky.cloud/docs/products/cloud/vis/template-variables/",
			},
			verifyCfg: func(t *testing.T, cfg Cfg) {
				require.Equal(t, "seeking help from your WideSky administrator", cfg.WideSkyWhitelabeling.EntityNotFoundText)
				require.Equal(t, "mailto: support@widesky.cloud", cfg.WideSkyWhitelabeling.EntityNotFoundLink)
				require.Equal(t, "Contact Support.", cfg.WideSkyWhitelabeling.EntityNotFoundLinkText)
				require.Equal(t, "https://widesky.cloud/docs/products/cloud/vis/template-variables/", cfg.WideSkyWhitelabeling.TemplateVariableHelpLink)
			},
		},
	}

	t.Run("Should not set cfg when setion is not set", func(t *testing.T) {
		f := ini.Empty()
		cfg := NewCfg()
		cfg.Raw = f
		cfg.readWideSkyWhitelabeling()
		require.Equal(t, (*WideSkyWhitelabelingSettings)(nil), cfg.WideSkyWhitelabeling)
	})

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			f := ini.Empty()
			cfg := NewCfg()
			wideSkyWhitelabelingSec, err := f.NewSection("ws_whitelabel")
			require.NoError(t, err)
			for k, v := range tc.whitelabelOptions {
				_, err = wideSkyWhitelabelingSec.NewKey(k, v)
				require.NoError(t, err)
			}
			cfg.Raw = f
			cfg.readWideSkyWhitelabeling()
			tc.verifyCfg(t, *cfg)
		})
	}
}
