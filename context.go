package togglr

type RequestContext map[string]any

func NewContext() RequestContext {
	return make(RequestContext)
}

func (r RequestContext) WithUserID(id string) RequestContext {
	r[AttrUserID] = id

	return r
}

func (r RequestContext) WithUserEmail(email string) RequestContext {
	r[AttrUserEmail] = email

	return r
}

func (r RequestContext) WithAnonymous(flag bool) RequestContext {
	r[AttrUserAnonymous] = flag

	return r
}

func (r RequestContext) WithCountry(code string) RequestContext {
	r[AttrCountryCode] = code

	return r
}

func (r RequestContext) WithRegion(region string) RequestContext {
	r[AttrRegion] = region

	return r
}

func (r RequestContext) WithCity(city string) RequestContext {
	r[AttrCity] = city

	return r
}

func (r RequestContext) WithManufacturer(m string) RequestContext {
	r[AttrManufacturer] = m

	return r
}

func (r RequestContext) WithDeviceType(t string) RequestContext {
	r[AttrDeviceType] = t

	return r
}

func (r RequestContext) WithOS(os string) RequestContext {
	r[AttrOS] = os

	return r
}

func (r RequestContext) WithOSVersion(v string) RequestContext {
	r[AttrOSVersion] = v

	return r
}

func (r RequestContext) WithBrowser(b string) RequestContext {
	r[AttrBrowser] = b

	return r
}

func (r RequestContext) WithBrowserVersion(v string) RequestContext {
	r[AttrBrowserVersion] = v
	return r
}

func (r RequestContext) WithLanguage(lang string) RequestContext {
	r[AttrLanguage] = lang

	return r
}

func (r RequestContext) WithConnectionType(ct string) RequestContext {
	r[AttrConnectionType] = ct

	return r
}

func (r RequestContext) WithAge(age int) RequestContext {
	r[AttrAge] = age

	return r
}

func (r RequestContext) WithGender(g string) RequestContext {
	r[AttrGender] = g

	return r
}

func (r RequestContext) WithIP(ip string) RequestContext {
	r[AttrIP] = ip

	return r
}

func (r RequestContext) WithAppVersion(ver string) RequestContext {
	r[AttrAppVersion] = ver

	return r
}

func (r RequestContext) WithPlatform(p string) RequestContext {
	r[AttrPlatform] = p

	return r
}

func (r RequestContext) Set(key string, value any) RequestContext {
	r[key] = value

	return r
}
