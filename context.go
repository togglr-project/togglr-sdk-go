package togglr

// RequestContext represents the context for feature evaluation
type RequestContext map[string]any

// NewContext creates a new empty request context
func NewContext() RequestContext {
	return make(RequestContext)
}

// Chainable helper methods for building context

// WithUserID sets the user ID
func (r RequestContext) WithUserID(id string) RequestContext {
	r[AttrUserID] = id

	return r
}

// WithUserEmail sets the user email
func (r RequestContext) WithUserEmail(email string) RequestContext {
	r[AttrUserEmail] = email

	return r
}

// WithAnonymous sets whether the user is anonymous
func (r RequestContext) WithAnonymous(flag bool) RequestContext {
	r[AttrUserAnonymous] = flag

	return r
}

// WithCountry sets the country code
func (r RequestContext) WithCountry(code string) RequestContext {
	r[AttrCountryCode] = code

	return r
}

// WithRegion sets the region
func (r RequestContext) WithRegion(region string) RequestContext {
	r[AttrRegion] = region

	return r
}

// WithCity sets the city
func (r RequestContext) WithCity(city string) RequestContext {
	r[AttrCity] = city

	return r
}

// WithManufacturer sets the device manufacturer
func (r RequestContext) WithManufacturer(m string) RequestContext {
	r[AttrManufacturer] = m

	return r
}

// WithDeviceType sets the device type
func (r RequestContext) WithDeviceType(t string) RequestContext {
	r[AttrDeviceType] = t

	return r
}

// WithOS sets the operating system
func (r RequestContext) WithOS(os string) RequestContext {
	r[AttrOS] = os

	return r
}

// WithOSVersion sets the operating system version
func (r RequestContext) WithOSVersion(v string) RequestContext {
	r[AttrOSVersion] = v

	return r
}

// WithBrowser sets the browser
func (r RequestContext) WithBrowser(b string) RequestContext {
	r[AttrBrowser] = b

	return r
}

// WithBrowserVersion sets the browser version
func (r RequestContext) WithBrowserVersion(v string) RequestContext {
	r[AttrBrowserVersion] = v
	return r
}

// WithLanguage sets the language
func (r RequestContext) WithLanguage(lang string) RequestContext {
	r[AttrLanguage] = lang

	return r
}

// WithConnectionType sets the connection type
func (r RequestContext) WithConnectionType(ct string) RequestContext {
	r[AttrConnectionType] = ct

	return r
}

// WithAge sets the user age
func (r RequestContext) WithAge(age int) RequestContext {
	r[AttrAge] = age

	return r
}

// WithGender sets the user gender
func (r RequestContext) WithGender(g string) RequestContext {
	r[AttrGender] = g

	return r
}

// WithIP sets the IP address
func (r RequestContext) WithIP(ip string) RequestContext {
	r[AttrIP] = ip

	return r
}

// WithAppVersion sets the application version
func (r RequestContext) WithAppVersion(ver string) RequestContext {
	r[AttrAppVersion] = ver

	return r
}

// WithPlatform sets the platform
func (r RequestContext) WithPlatform(p string) RequestContext {
	r[AttrPlatform] = p

	return r
}

// Set sets an arbitrary key-value pair
func (r RequestContext) Set(key string, value any) RequestContext {
	r[key] = value

	return r
}
