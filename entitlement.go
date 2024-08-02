package discordgo

// EntitlementType is the type of entitlement (see EntitlementType* consts)
// https://discord.com/developers/docs/monetization/entitlements#entitlement-object-entitlement-types
type EntitlementType int

// Valid EntitlementType values
const (
	EntitlementTypePurchase = 1
	EntitlementTypePremiumSubscription = 2
	EntitlementTypeDeveloperGift = 3
	EntitlementTypeTestModePurchase = 4
	EntitlementTypeFreePurchase = 5
	EntitlementTypeUserGift = 6
	EntitlementTypePremiumPurchase = 7
	EntitlementTypeApplicationSubscription = 8
)

// Entitlements represent that a user or guild has access to a premium offering
// in your application.
type Entitlement struct {
	// The ID of the entitlement
	ID string `json:"id"`

	// The ID of the SKU
	SKUID string `json:"sku_id"`

	// The ID of the parent application
	ApplicationID string `json:"application_id"`

	// The ID of the user that is granted access to the entitlement's sku
	UserID string `json:"user_id"`

	// The type of entitlement
	Type EntitlementType `json:"type"`

	// The entitlement was deleted
	Deleted bool `json:"deleted"`

	// The start date at which the entitlement is valid. 
	// Not present when using test entitlements.
	StartsAt string `json:"starts_at"`

	// The date at which the entitlement is no longer valid. 
	// Not present when using test entitlements.
	EndsAt string `json:"ends_at"`

	// The ID of the guild that is granted access to the entitlement's sku.
	GuildID string `json:"guild_id"`

	// Whether or not the entitlement has been consumed.
	// Only available for consumable items.
	Consumed bool `json:"consumed"`
}