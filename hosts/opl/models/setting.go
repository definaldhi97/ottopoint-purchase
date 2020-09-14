package models

type SettingOPL struct {
	Settings struct {
		PointsDaysActiveCount                  int      `json:"pointsDaysActiveCount"`
		ExpirePointsNotificationDays           int      `json:"expirePointsNotificationDays"`
		ExpireCouponsNotificationDays          int      `json:"expireCouponsNotificationDays"`
		ExpireLevelsNotificationDays           int      `json:"expireLevelsNotificationDays"`
		Returns                                bool     `json:"returns"`
		AllowCustomersProfileEdits             bool     `json:"allowCustomersProfileEdits"`
		AllTimeNotLocked                       bool     `json:"allTimeNotLocked"`
		LevelResetPointsOnDowngrade            bool     `json:"levelResetPointsOnDowngrade"`
		Webhooks                               bool     `json:"webhooks"`
		ExcludeDeliveryCostsFromTierAssignment bool     `json:"excludeDeliveryCostsFromTierAssignment"`
		ExcludedLevelCategories                []string `json:"excludedLevelCategories"`
		CustomerStatusesEarning                []string `json:"customerStatusesEarning"`
		CustomerStatusesSpending               []string `json:"customerStatusesSpending"`
		CustomersIdentificationPriority        []struct {
			Priority int    `json:"priority"`
			Field    string `json:"field"`
		} `json:"customersIdentificationPriority"`
		ExcludedDeliverySKUs []interface{} `json:"excludedDeliverySKUs"`
		ExcludedLevelSKUs    []interface{} `json:"excludedLevelSKUs"`
		Logo                 struct {
			Path  string        `json:"path"`
			Mime  string        `json:"mime"`
			Sizes []interface{} `json:"sizes"`
		} `json:"logo"`
		SmallLogo struct {
			Path  string        `json:"path"`
			Mime  string        `json:"mime"`
			Sizes []interface{} `json:"sizes"`
		} `json:"small-logo"`
		HeroImage struct {
			Path  string        `json:"path"`
			Mime  string        `json:"mime"`
			Sizes []interface{} `json:"sizes"`
		} `json:"hero-image"`
		AdminCockpitLogo struct {
			Path  string        `json:"path"`
			Mime  string        `json:"mime"`
			Sizes []interface{} `json:"sizes"`
		} `json:"admin-cockpit-logo"`
		ClientCockpitLogoBig struct {
			Path  string        `json:"path"`
			Mime  string        `json:"mime"`
			Sizes []interface{} `json:"sizes"`
		} `json:"client-cockpit-logo-big"`
		ClientCockpitLogoSmall struct {
			Path  string        `json:"path"`
			Mime  string        `json:"mime"`
			Sizes []interface{} `json:"sizes"`
		} `json:"client-cockpit-logo-small"`
		ClientCockpitHeroImage struct {
			Path  string        `json:"path"`
			Mime  string        `json:"mime"`
			Sizes []interface{} `json:"sizes"`
		} `json:"client-cockpit-hero-image"`
		Currency                       string `json:"currency"`
		Timezone                       string `json:"timezone"`
		ProgramName                    string `json:"programName"`
		ProgramPointsSingular          string `json:"programPointsSingular"`
		ProgramPointsPlural            string `json:"programPointsPlural"`
		PointsDaysExpiryAfter          string `json:"pointsDaysExpiryAfter"`
		TierAssignType                 string `json:"tierAssignType"`
		LevelDowngradeMode             string `json:"levelDowngradeMode"`
		LevelDowngradeBase             string `json:"levelDowngradeBase"`
		AccountActivationMethod        string `json:"accountActivationMethod"`
		MarketingVendorsValue          string `json:"marketingVendorsValue"`
		PushySecretKey                 string `json:"pushySecretKey"`
		MaxPointsRedeemed              string `json:"maxPointsRedeemed"`
		TransactionTriggeredSmsContent string `json:"transactionTriggeredSmsContent"`
		ProgramConditionsURL           string `json:"programConditionsUrl"`
		ProgramFaqURL                  string `json:"programFaqUrl"`
		ProgramURL                     string `json:"programUrl"`
		HelpEmailAddress               string `json:"helpEmailAddress"`
		URIWebhooks                    string `json:"uriWebhooks"`
		WebhookHeaderName              string `json:"webhookHeaderName"`
		WebhookHeaderValue             string `json:"webhookHeaderValue"`
		AccentColor                    string `json:"accentColor"`
		CSSTemplate                    string `json:"cssTemplate"`
	} `json:"settings"`
}
