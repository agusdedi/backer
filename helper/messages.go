package helper

// Common messages
const (
	MsgInvalidInput = "Invalid input"
)

// User messages
const (
	MsgEmailAlreadyRegistered         = "Email already registered"
	MsgRegistrationFailed             = "Registration failed"
	MsgFailedToGenerateToken          = "Failed to generate authentication token"
	MsgAccountRegisteredSuccessfully  = "Account registered successfully"
	MsgInvalidEmailOrPassword         = "Invalid email or password"
	MsgSuccessfullyLoggedIn           = "Successfully logged in"
	MsgEmailValidationFailed          = "Email validation failed"
	MsgEmailCheckFailed               = "Email check failed"
	MsgUnableToCheckEmailAvailability = "Unable to check email availability"
	MsgEmailAlreadyRegisteredCheck    = "Email is already registered"
	MsgEmailAvailable                 = "Email is available"
	MsgNoAvatarFileProvided           = "No avatar file provided"
	MsgFailedToSaveAvatarFile         = "Failed to save avatar file"
	MsgFailedToUpdateAvatar           = "Failed to update avatar"
	MsgAvatarUploadedSuccessfully     = "Avatar uploaded successfully"
)

// Campaign messages
const (
	MsgFailedToGetCampaigns              = "Failed to get campaigns"
	MsgListOfCampaignsRetrieved          = "List of campaigns retrieved successfully"
	MsgInvalidCampaignID                 = "Invalid campaign ID"
	MsgCampaignNotFound                  = "Campaign not found"
	MsgCampaignDetailRetrieved           = "Campaign detail retrieved successfully"
	MsgFailedToCreateCampaign            = "Failed to create campaign"
	MsgCampaignCreatedSuccessfully       = "Campaign created successfully"
	MsgNotAuthorizedToUpdateCampaign     = "You are not authorized to update this campaign"
	MsgFailedToUpdateCampaign            = "Failed to update campaign"
	MsgCampaignUpdatedSuccessfully       = "Campaign updated successfully"
	MsgNoFileUploaded                    = "No file uploaded"
	MsgInvalidFileType                   = "Invalid file type. Only JPG, JPEG, and PNG are allowed"
	MsgFileTooLarge                      = "File size too large. Maximum 2MB"
	MsgFailedToSaveFileToServer          = "Failed to save file to server"
	MsgNotAuthorizedToUploadImage        = "You are not authorized to upload image for this campaign"
	MsgFailedToSaveImageToDatabase       = "Failed to save image to database"
	MsgCampaignImageUploadedSuccessfully = "Campaign image uploaded successfully"
)
