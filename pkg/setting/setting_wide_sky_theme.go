package setting

type ColorGroup struct {
	Main  string `json:"main"`
	Shade string `json:"shade"`
	Text  string `json:"text"`
}

type BrandGradient struct {
	ColorA string `json:"colorA"`
	ColorB string `json:"colorB"`
	All    string `json:"all"`
}

type WideSkyThemeSettings struct {
	Name        string `json:"name"`
	BaseColor   string `json:"baseColor"`
	BorderColor string `json:"borderColor"`

	FontHeader string `json:"fontHeader"`
	FontBody   string `json:"fontBody"`

	Horizontal BrandGradient `json:"horizontal"`
	Vertical   BrandGradient `json:"vertical"`

	Primary   ColorGroup `json:"primary"`
	Secondary ColorGroup `json:"secondary"`
	Info      ColorGroup `json:"info"`
	Success   ColorGroup `json:"success"`
	Warning   ColorGroup `json:"warning"`
	Error     ColorGroup `json:"error"`

	TextPrimary   string `json:"textPrimary"`
	TextSecondary string `json:"textSecondary"`
	TextDisabled  string `json:"textDisabled"`
	TextLink      string `json:"textLink"`

	BackgroundCanvas    string `json:"backgroundCanvas"`
	BackgroundPrimary   string `json:"backgroundPrimary"`
	BackgroundSecondary string `json:"backgroundSecondary"`
	BackgroundNavBar    string `json:"backgroundNavBar"`

	BorderWeak   string `json:"borderWeak"`
	BorderMedium string `json:"borderMedium"`
	BorderStrong string `json:"borderStrong"`

	ActionHover              string `json:"actionHover"`
	ActionFocus              string `json:"actionFocus"`
	ActionSelected           string `json:"actionSelected"`
	ActionDisabledBackground string `json:"actionDisabledBackground"`
}

func (cfg *Cfg) readWideSkyTheme() {
	hasThemeSec := cfg.Raw.HasSection("ws_theme")

	if !hasThemeSec {
		cfg.Logger.Info("A custom WideSky theme is not configured")
		return
	}

	cfg.Logger.Info("Creating a WideSky theme configuration")
	themeSec := cfg.Raw.Section("ws_theme")
	NOT_SET := "NOT_SET"

	baseColor := themeSec.Key("base_colour").MustString("255,255,255")

	cfg.WideSkyTheme = &WideSkyThemeSettings{
		Name:        themeSec.Key("name").MustString("WideSky"),
		BaseColor:   baseColor,
		BorderColor: themeSec.Key("border_colour").MustString(baseColor),

		FontHeader: themeSec.Key("font_header").MustString("MuseoSansRounded"),
		FontBody:   themeSec.Key("font_body").MustString("MuseoSansRounded"),

		Horizontal: BrandGradient{
			ColorA: themeSec.Key("brand_gradient_horizontal_colour_a").MustString("#D1D3D4"),
			ColorB: themeSec.Key("brand_gradient_horizontal_colour_b").MustString("#A7A9AC"),
			All:    themeSec.Key("brand_gradient_horizontal_all").MustString(NOT_SET),
		},
		Vertical: BrandGradient{
			ColorA: themeSec.Key("brand_gradient_vertical_colour_a").MustString("#D1D3D4"),
			ColorB: themeSec.Key("brand_gradient_vertical_colour_b").MustString("#A7A9AC"),
			All:    themeSec.Key("brand_gradient_vertical_all").MustString(NOT_SET),
		},

		Primary: ColorGroup{
			Main:  themeSec.Key("primary_main").MustString("#0073BC"),
			Shade: themeSec.Key("primary_hover").MustString("rgb(38, 136, 198)"),
			Text:  themeSec.Key("primary_text").MustString("#FFFFFF"),
		},
		Secondary: ColorGroup{
			Main:  themeSec.Key("secondary_main").MustString("rgba(204, 204, 220, 0.16)"),
			Shade: themeSec.Key("secondary_hover").MustString("rgba(204, 204, 220, 0.20)"),
			Text:  themeSec.Key("secondary_text").MustString("rgb(204, 204, 220)"),
		},
		Info: ColorGroup{
			Main:  themeSec.Key("info_main").MustString("#0073BC"),
			Shade: themeSec.Key("info_hover").MustString("rgb(38, 136, 198)"),
			Text:  themeSec.Key("info_text").MustString("#FFFFFF"),
		},
		Success: ColorGroup{
			Main:  themeSec.Key("success_main").MustString("#1A7F4B"),
			Shade: themeSec.Key("success_hover").MustString("rgb(60, 146, 102)"),
			Text:  themeSec.Key("success_text").MustString("#FFFFFF"),
		},
		Warning: ColorGroup{
			Main:  themeSec.Key("warning_main").MustString("#F5B73D"),
			Shade: themeSec.Key("warning_hover").MustString("rgb(246, 193, 90)"),
			Text:  themeSec.Key("warning_text").MustString("#000000"),
		},
		Error: ColorGroup{
			Main:  themeSec.Key("error_main").MustString("#D10E5C"),
			Shade: themeSec.Key("error_hover").MustString("rgb(215, 50, 116)"),
			Text:  themeSec.Key("error_text").MustString("#FFFFFF"),
		},

		TextPrimary:   themeSec.Key("text_primary").MustString("rgb(204, 204, 220)"),
		TextSecondary: themeSec.Key("text_secondary").MustString("rgba(204, 204, 220, 0.65)"),
		TextDisabled:  themeSec.Key("text_disabled").MustString("rgba(204, 204, 220, 0.6)"),
		TextLink:      themeSec.Key("text_link").MustString(NOT_SET),

		BackgroundCanvas:    themeSec.Key("background_canvas").MustString("#0D294B"),
		BackgroundPrimary:   themeSec.Key("background_primary").MustString("#0073bc"),
		BackgroundSecondary: themeSec.Key("background_secondary").MustString("#4e9ed7"),
		BackgroundNavBar:    themeSec.Key("background_nav_bar").MustString(NOT_SET),

		BorderWeak:   themeSec.Key("border_weak").MustString("rgba(204, 204, 220, 0.15)"),
		BorderMedium: themeSec.Key("border_medium").MustString("rgba(204, 204, 220, 0.25)"),
		BorderStrong: themeSec.Key("border_strong").MustString("rgba(204, 204, 220, 0.07)"),

		ActionHover:              themeSec.Key("action_hover").MustString("rgba(204, 204, 220, 0.16)"),
		ActionFocus:              themeSec.Key("action_focus").MustString("rgba(204, 204, 220, 0.16)"),
		ActionSelected:           themeSec.Key("action_selected").MustString("rgba(204, 204, 220, 0.12)"),
		ActionDisabledBackground: themeSec.Key("action_disabled_background").MustString("rgba(204, 204, 220, 0.04)"),
	}
}
