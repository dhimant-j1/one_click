package models

import "time"

// Enum type definitions
type UserRole string
type QuoteStatus string
type PolicyStatus string

// UserRole enum constants
const (
	RoleMasterAdmin   UserRole = "MasterAdmin"
	RoleAgencyAdmin   UserRole = "AgencyAdmin"
	RoleLocationAdmin UserRole = "LocationAdmin"
	RoleAgent         UserRole = "Agent"
	RoleCustomer      UserRole = "Customer"
)

// QuoteStatus enum constants
const (
	QuoteStatusDraft     QuoteStatus = "Draft"
	QuoteStatusPresented QuoteStatus = "Presented"
	QuoteStatusBound     QuoteStatus = "Bound"
)

// PolicyStatus enum constants
const (
	PolicyStatusActive    PolicyStatus = "Active"
	PolicyStatusExpired   PolicyStatus = "Expired"
	PolicyStatusCancelled PolicyStatus = "Cancelled"
)

// User represents the Users table

type User struct {
	UserID                int        `json:"user_id"`
	FirstName             string     `json:"first_name"`
	LastName              string     `json:"last_name"`
	Email                 string     `json:"email"`
	PasswordHash          string     `json:"password_hash"`
	Role                  UserRole   `json:"role"`
	AgencyID              *int       `json:"agency_id,omitempty"`
	LocationID            *int       `json:"location_id,omitempty"`
	Address               *string    `json:"address,omitempty"`
	LicenseNumber         *string    `json:"license_number,omitempty"`
	LicenseState          *string    `json:"license_state,omitempty"`
	DateOfBirth           *time.Time `json:"date_of_birth,omitempty"`
	MobileNumber          *string    `json:"mobile_number,omitempty"`
	HasPhysicalImpairment *bool      `json:"has_physical_impairment,omitempty"`
	NeedsFinancialFiling  *bool      `json:"needs_financial_filing,omitempty"`
	AdditionalData        *string    `json:"additional_data,omitempty"`
}

// Agency represents the Agencies table

type Agency struct {
	AgencyID       int     `json:"agency_id"`
	AgencyName     string  `json:"agency_name"`
	AgentCode      *string `json:"agent_code,omitempty"`
	AdditionalData *string `json:"additional_data,omitempty"`
}

// InsuranceProvider represents the InsuranceProviders table

type InsuranceProvider struct {
	ProviderID     int     `json:"provider_id"`
	ProviderName   string  `json:"provider_name"`
	ContactInfo    string  `json:"contact_info"`
	AdditionalData *string `json:"additional_data,omitempty"`
}

// AgencyProviderAccess represents the AgencyProviderAccess table

type AgencyProviderAccess struct {
	AgencyID   int `json:"agency_id"`
	ProviderID int `json:"provider_id"`
}

// Quote represents the Quotes table

type Quote struct {
	QuoteID        int         `json:"quote_id"`
	AgentID        int         `json:"agent_id"`
	CustomerID     int         `json:"customer_id"`
	VehicleID      *int        `json:"vehicle_id,omitempty"`
	QuoteDate      time.Time   `json:"quote_date"`
	Status         QuoteStatus `json:"status"`
	AdditionalData *string     `json:"additional_data,omitempty"`
}

// Policy represents the Policies table

type Policy struct {
	PolicyID       int          `json:"policy_id"`
	CustomerID     int          `json:"customer_id"`
	ProviderID     int          `json:"provider_id"`
	QuoteID        *int         `json:"quote_id,omitempty"`
	PolicyNumber   string       `json:"policy_number"`
	EffectiveDate  time.Time    `json:"effective_date"`
	ExpirationDate time.Time    `json:"expiration_date"`
	Status         PolicyStatus `json:"status"`
	AdditionalData *string      `json:"additional_data,omitempty"`
}

// AgencyLocation represents the AgencyLocations table
type AgencyLocation struct {
	LocationID     int     `json:"location_id"`
	AgencyID       int     `json:"agency_id"`
	Address        string  `json:"address"`
	AdditionalData *string `json:"additional_data,omitempty"`
}

// Vehicle represents the Vehicles table
type Vehicle struct {
	VehicleID              int        `json:"vehicle_id"`
	CustomerID             int        `json:"customer_id"`
	VIN                    string     `json:"vin"`
	Make                   *string    `json:"make,omitempty"`
	Model                  *string    `json:"model,omitempty"`
	Year                   *int       `json:"year,omitempty"`
	Type                   *string    `json:"type,omitempty"`
	PlateNumber            *string    `json:"plate_number,omitempty"`
	PlateType              *string    `json:"plate_type,omitempty"`
	BodyStyle              *string    `json:"body_style,omitempty"`
	VehicleUse             *string    `json:"vehicle_use,omitempty"`
	VehicleHistory         *string    `json:"vehicle_history,omitempty"`
	AnnualMileage          *int       `json:"annual_mileage,omitempty"`
	CurrentOdometerReading *int       `json:"current_odometer_reading,omitempty"`
	OdometerReadingDate    *time.Time `json:"odometer_reading_date,omitempty"`
	PurchasedInLast90Days  *bool      `json:"purchased_in_last_90_days,omitempty"`
	LengthOfOwnership      *string    `json:"length_of_ownership,omitempty"`
	OwnershipStatus        *string    `json:"ownership_status,omitempty"`
	HasRacingEquipment     *bool      `json:"has_racing_equipment,omitempty"`
	HasExistingDamage      *bool      `json:"has_existing_damage,omitempty"`
	CostNew                *float64   `json:"cost_new,omitempty"`
	AdditionalData         *string    `json:"additional_data,omitempty"`
}

// VehicleDriver represents the VehicleDrivers table
type VehicleDriver struct {
	VehicleID                          int     `json:"vehicle_id"`
	UserID                             int     `json:"user_id"`
	DriverType                         *string `json:"driver_type,omitempty"`
	RelationshipToInsured              *string `json:"relationship_to_insured,omitempty"`
	Gender                             *string `json:"gender,omitempty"`
	MaritalStatus                      *string `json:"marital_status,omitempty"`
	LicensingException                 *string `json:"licensing_exception,omitempty"`
	NeedsFinancialResponsibilityFiling *bool   `json:"needs_financial_responsibility_filing,omitempty"`
	HasUncompensatedImpairment         *bool   `json:"has_uncompensated_impairment,omitempty"`
	MobileNumber                       *string `json:"mobile_number,omitempty"`
	AdditionalData                     *string `json:"additional_data,omitempty"`
}

// DrivingHistory represents the DrivingHistory table
type DrivingHistory struct {
	HistoryID      int        `json:"history_id"`
	UserID         int        `json:"user_id"`
	IncidentType   *string    `json:"incident_type,omitempty"`
	IncidentDate   *time.Time `json:"incident_date,omitempty"`
	ConvictionDate *time.Time `json:"conviction_date,omitempty"`
	Amount         *float64   `json:"amount,omitempty"`
	Description    *string    `json:"description,omitempty"`
	AdditionalData *string    `json:"additional_data,omitempty"`
}

// InsuranceHistory represents the InsuranceHistory table
type InsuranceHistory struct {
	HistoryID                   int     `json:"history_id"`
	UserID                      int     `json:"user_id"`
	InsuranceStatus             *string `json:"insurance_status,omitempty"`
	CurrentCarrier              *string `json:"current_carrier,omitempty"`
	CurrentBodilyInjuryLimits   *string `json:"current_bodily_injury_limits,omitempty"`
	LengthWithCurrentCompany    *string `json:"length_with_current_company,omitempty"`
	ContinuousInsurance         *bool   `json:"continuous_insurance,omitempty"`
	VehicleRegisteredToOther    *bool   `json:"vehicle_registered_to_other,omitempty"`
	DriverWithoutLicense        *bool   `json:"driver_without_license,omitempty"`
	LicenseSuspendedRevoked     *bool   `json:"license_suspended_revoked,omitempty"`
	DeclinedCancelledNonRenewed *bool   `json:"declined_cancelled_non_renewed,omitempty"`
	MilitaryDeployment          *bool   `json:"military_deployment,omitempty"`
	PrimaryResidence            *string `json:"primary_residence,omitempty"`
	AdditionalData              *string `json:"additional_data,omitempty"`
}

// Coverage represents the Coverages table
type Coverage struct {
	CoverageID     int     `json:"coverage_id"`
	CoverageName   string  `json:"coverage_name"`
	Description    *string `json:"description,omitempty"`
	AdditionalData *string `json:"additional_data,omitempty"`
}

// QuoteLineItem represents the QuoteLineItems table
type QuoteLineItem struct {
	LineItemID       int     `json:"line_item_id"`
	QuoteID          int     `json:"quote_id"`
	ProviderID       int     `json:"provider_id"`
	CoverageID       int     `json:"coverage_id"`
	Price            float64 `json:"price"`
	LimitAmount      *string `json:"limit_amount,omitempty"`
	DeductibleAmount *string `json:"deductible_amount,omitempty"`
	AdditionalData   *string `json:"additional_data,omitempty"`
}
