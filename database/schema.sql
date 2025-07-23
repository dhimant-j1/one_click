-- Create ENUM types first
CREATE TYPE user_role AS ENUM ('MasterAdmin', 'AgencyAdmin', 'LocationAdmin', 'Agent', 'Customer');
CREATE TYPE quote_status AS ENUM ('Draft', 'Presented', 'Bound');
CREATE TYPE policy_status AS ENUM ('Active', 'Expired', 'Cancelled');

CREATE TABLE Agencies (
    AgencyID SERIAL PRIMARY KEY,
    AgencyName VARCHAR(255) NOT NULL,
    AgentCode VARCHAR(255),
    AdditionalData JSONB
);

CREATE TABLE AgencyLocations (
    LocationID SERIAL PRIMARY KEY,
    AgencyID INT NOT NULL,
    Address TEXT NOT NULL,
    AdditionalData JSONB,
    FOREIGN KEY (AgencyID) REFERENCES Agencies(AgencyID)
);

CREATE TABLE Users (
    UserID SERIAL PRIMARY KEY,
    FirstName VARCHAR(255),
    LastName VARCHAR(255),
    Email VARCHAR(255) UNIQUE NOT NULL,
    PasswordHash VARCHAR(255) NOT NULL,
    Role user_role NOT NULL,
    AgencyID INT,
    LocationID INT,
    Address TEXT,
    LicenseNumber VARCHAR(255),
    LicenseState VARCHAR(50),
    DateOfBirth DATE,
    MobileNumber VARCHAR(50),
    HasPhysicalImpairment BOOLEAN,
    NeedsFinancialFiling BOOLEAN,
    AdditionalData JSONB,
    FOREIGN KEY (AgencyID) REFERENCES Agencies(AgencyID) ON DELETE SET NULL,
    FOREIGN KEY (LocationID) REFERENCES AgencyLocations(LocationID) ON DELETE SET NULL
);

CREATE TABLE UserSessions (
    SessionID SERIAL PRIMARY KEY,
    UserID INT NOT NULL,
    SessionToken VARCHAR(255) UNIQUE NOT NULL,
    ExpiresAt TIMESTAMP NOT NULL,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    IPAddress VARCHAR(45),
    UserAgent TEXT,
    FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE CASCADE
);

CREATE TABLE PasswordResetTokens (
    TokenID SERIAL PRIMARY KEY,
    UserID INT NOT NULL,
    ResetToken VARCHAR(255) UNIQUE NOT NULL,
    ExpiresAt TIMESTAMP NOT NULL,
    IsUsed BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE CASCADE
);

CREATE TABLE InsuranceProviders (
    ProviderID SERIAL PRIMARY KEY,
    ProviderName VARCHAR(255) NOT NULL,
    ContactInfo TEXT,
    AdditionalData JSONB
);

CREATE TABLE AgencyProviderAccess (
    AgencyID INT NOT NULL,
    ProviderID INT NOT NULL,
    PRIMARY KEY (AgencyID, ProviderID),
    FOREIGN KEY (AgencyID) REFERENCES Agencies(AgencyID) ON DELETE CASCADE,
    FOREIGN KEY (ProviderID) REFERENCES InsuranceProviders(ProviderID) ON DELETE CASCADE
);

CREATE TABLE Vehicles (
    VehicleID SERIAL PRIMARY KEY,
    CustomerID INT NOT NULL,
    VIN VARCHAR(255) NOT NULL,
    Make VARCHAR(100),
    Model VARCHAR(100),
    Year INT,
    Type VARCHAR(100),
    PlateNumber VARCHAR(50),
    PlateType VARCHAR(50),
    BodyStyle VARCHAR(100),
    VehicleUse VARCHAR(100),
    VehicleHistory TEXT,
    AnnualMileage INT,
    CurrentOdometerReading INT,
    OdometerReadingDate DATE,
    PurchasedInLast90Days BOOLEAN,
    LengthOfOwnership VARCHAR(100),
    OwnershipStatus VARCHAR(100),
    HasRacingEquipment BOOLEAN,
    HasExistingDamage BOOLEAN,
    CostNew DECIMAL(10, 2),
    AdditionalData JSONB,
    FOREIGN KEY (CustomerID) REFERENCES Users(UserID) ON DELETE CASCADE
);

CREATE TABLE VehicleDrivers (
    VehicleID INT NOT NULL,
    UserID INT NOT NULL,
    DriverType VARCHAR(100),
    RelationshipToInsured VARCHAR(100),
    Gender VARCHAR(10),
    MaritalStatus VARCHAR(50),
    LicensingException TEXT,
    NeedsFinancialResponsibilityFiling BOOLEAN,
    HasUncompensatedImpairment BOOLEAN,
    MobileNumber VARCHAR(50),
    AdditionalData JSONB,
    PRIMARY KEY (VehicleID, UserID),
    FOREIGN KEY (VehicleID) REFERENCES Vehicles(VehicleID) ON DELETE CASCADE,
    FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE CASCADE
);

CREATE TABLE DrivingHistory (
    HistoryID SERIAL PRIMARY KEY,
    UserID INT NOT NULL,
    IncidentType VARCHAR(255),
    IncidentDate DATE,
    ConvictionDate DATE,
    Amount DECIMAL(10, 2),
    Description TEXT,
    AdditionalData JSONB,
    FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE CASCADE
);

CREATE TABLE InsuranceHistory (
    HistoryID SERIAL PRIMARY KEY,
    UserID INT NOT NULL,
    InsuranceStatus VARCHAR(100),
    CurrentCarrier VARCHAR(255),
    CurrentBodilyInjuryLimits VARCHAR(100),
    LengthWithCurrentCompany VARCHAR(100),
    ContinuousInsurance BOOLEAN,
    VehicleRegisteredToOther BOOLEAN,
    DriverWithoutLicense BOOLEAN,
    LicenseSuspendedRevoked BOOLEAN,
    DeclinedCancelledNonRenewed BOOLEAN,
    MilitaryDeployment BOOLEAN,
    PrimaryResidence TEXT,
    AdditionalData JSONB,
    FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE CASCADE
);

CREATE TABLE Coverages (
    CoverageID SERIAL PRIMARY KEY,
    CoverageName VARCHAR(255) NOT NULL,
    Description TEXT,
    AdditionalData JSONB
);

CREATE TABLE Quotes (
    QuoteID SERIAL PRIMARY KEY,
    AgentID INT NOT NULL,
    CustomerID INT NOT NULL,
    VehicleID INT,
    QuoteDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Status quote_status,
    AdditionalData JSONB,
    FOREIGN KEY (AgentID) REFERENCES Users(UserID) ON DELETE CASCADE,
    FOREIGN KEY (CustomerID) REFERENCES Users(UserID) ON DELETE CASCADE,
    FOREIGN KEY (VehicleID) REFERENCES Vehicles(VehicleID) ON DELETE SET NULL
);

CREATE TABLE QuoteLineItems (
    LineItemID SERIAL PRIMARY KEY,
    QuoteID INT NOT NULL,
    ProviderID INT NOT NULL,
    CoverageID INT NOT NULL,
    Price DECIMAL(10, 2) NOT NULL,
    LimitAmount VARCHAR(100),
    DeductibleAmount VARCHAR(100),
    AdditionalData JSONB,
    FOREIGN KEY (QuoteID) REFERENCES Quotes(QuoteID) ON DELETE CASCADE,
    FOREIGN KEY (ProviderID) REFERENCES InsuranceProviders(ProviderID) ON DELETE CASCADE,
    FOREIGN KEY (CoverageID) REFERENCES Coverages(CoverageID) ON DELETE CASCADE
);

CREATE TABLE Policies (
    PolicyID SERIAL PRIMARY KEY,
    CustomerID INT NOT NULL,
    ProviderID INT NOT NULL,
    QuoteID INT,
    PolicyNumber VARCHAR(255) UNIQUE NOT NULL,
    EffectiveDate DATE NOT NULL,
    ExpirationDate DATE NOT NULL,
    Status policy_status,
    AdditionalData JSONB,
    FOREIGN KEY (CustomerID) REFERENCES Users(UserID) ON DELETE CASCADE,
    FOREIGN KEY (ProviderID) REFERENCES InsuranceProviders(ProviderID) ON DELETE CASCADE,
    FOREIGN KEY (QuoteID) REFERENCES Quotes(QuoteID) ON DELETE SET NULL
);
