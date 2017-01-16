package oauthgoogle

import (
	"bytes"
	"encoding/json"
	"github.com/mattermost/platform/einterfaces"
	"github.com/mattermost/platform/model"
	"io"
	"io/ioutil"

	l4g "github.com/alecthomas/log4go"
)

type GoogleProvider struct {
}

type GoogleUser struct {
	id             string
	email          string
	verified_email bool
	name           string
	given_name     string
	family_name    string
	picture        string
	locale         string
	hd             string
}

func init() {
	l4g.Debug("oauthgoogle: Initializing GoogleProvider")
	provider := &GoogleProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_GOOGLE, provider)
	l4g.Debug("oauthgoogle: GoogleProvider initialized @ %s", provider)
}

func (m *GoogleProvider) userFromGoogleUser(gu *GoogleUser) *model.User {

	var authData string
	user := &model.User{}
	if gu.IsValid() != true {
		return user
	}

	username := buildUsernameFromIdAndDomain(gu.given_name, gu.family_name)
	user.Username = model.CleanUsername(username)

	user.FirstName = gu.given_name
	user.LastName = gu.family_name
	user.Email = gu.email
	authData = gu.getAuthData()
	user.AuthData = &authData
	user.AuthService = m.GetIdentifier()

	return user
}

func (gu *GoogleUser) getAuthData() string {
	l4g.Debug("Entering getAuthData()")

	var authData string

	if len(gu.id) != 21 {
		authData = "false"
	} else {
		authData = gu.id
	}
	l4g.Debug("getAuthData: %s", authData)
	return authData
}

func (m *GoogleProvider) GetIdentifier() string {
	l4g.Debug("Entering getIdentifier()")
	l4g.Debug("GetIdentifier: %s", model.USER_AUTH_SERVICE_GOOGLE)
	return model.USER_AUTH_SERVICE_GOOGLE
}

func (m *GoogleProvider) GetUserFromJson(data io.Reader) *model.User {
	l4g.Debug("Entering GetUserFromJson()")
	l4g.Debug("GetUserFromJson: %s", data)
	gu := jsonDataDecode(data)
	if gu.IsValid() == true {
		return m.userFromGoogleUser(gu)
	}
	return &model.User{}

}

/* t, err := decode.Token()

	if err != nil {
		l4g.Debug("\t --> googleUserFromJson()->decode: ERROR->%s (%v)", err, t)
		return &gAuth
	}
	err := decode.Decode(&gAuth)
	if err != nil {
		l4g.Debug("\t --> googleUserFromJson()->decode: ERROR->%s", err)
		return &gAuth
	} else {
		l4g.Debug("\t --> googleUserFromJson()->decode: SUCCESS -> %v", gAuth)
	}

	var more string

	for {
		t, err := decode.Token()
		if err == io.EOF {
			l4g.Debug("\t --> googleUserFromJson()->decode: io.EOF")
			break
		}
		if err != nil {
			l4g.Debug("\t --> googleUserFromJson()->decode: ERROR->%s", err)
			return nil
		}
		if decode.More() {
			more = "(cont)"
		} else {
			more = "done!"
		}

		l4g.Debug("\t --> googleUserFromJson()->decode: %T: %v %s", t, t, more)
	}
	aliagnAuthData(&gAuth)
	return &gAuth
}
*/

func (gu *GoogleUser) IsValid() bool {
	l4g.Debug("Entering IsValid()")
	if gu == nil {
		l4g.Debug("(gu is NULL) IsValid() false")
		return false
	}

	if len(gu.email) == 0 {
		l4g.Debug("len(gu.email) == 0 <-> IsValid() false")
		return false
	}

	if len(gu.id) == 0 {
		l4g.Debug("len(gu.id) == 0 <-> IsValid() false")
		return false
	}

	return true
}

func buildUsernameFromIdAndDomain(left string, right string) string {
	var buffer bytes.Buffer
	buffer.WriteString(left)
	buffer.WriteString("_")
	buffer.WriteString(right)
	return buffer.String()

}

func __getEmailFromJason(
func jsonDataDecode(data io.Reader) *GoogleUser {
	l4g.Debug("entering googleAuthDecode()")

	var gu GoogleUser
	var err error
	var myJson []byte
	var parsed map[string]interface{}

	if myJson, err = ioutil.ReadAll(data); err != nil {
		l4g.Error("\t--> Unable to read oauth json data for parsing... (%s)", err)
		return nil
	}

	if err = json.Unmarshal(myJson, &parsed); err != nil {
		l4g.Error("\t--> Unable to decode json for google Authentication: (%s)", err)
		return nil
	}
	//err := json.NewDecoder(data).Decode(&gu)

	if (parsed["id"] == nil) || (len(parsed["id"].(string)) == 0) {
		l4g.Error("\t--> Unable to parse ID from OAuth json!")
		return nil
	}

	if (parsed["email"] == nil) || (len(parsed["email"].(string)) == 0) {
		l4g.Error("\t--> Unable to parse email from OAuth json!")
		return nil
	}

	gu.id = parsed["id"].(string)
	gu.email = parsed["email"].(string)
	//gu.username = buildUsernameFromIdAndDomain(parsed["id"].(string), parsed["hd"].(string))
	gu.verified_email = parsed["verified_email"].(bool)
	gu.given_name = parsed["given_name"].(string)
	gu.family_name = parsed["family_name"].(string)
	gu.locale = parsed["locale"].(string)
	gu.hd = parsed["hd"].(string)
	l4g.Debug("Good: jsonDataDecode(data, gu): %s <-> %s", gu.id, gu.email)
	return &gu
}

func (gu *GoogleUser) checkParsed(parsed interface{}) bool {
	if parsed == nil {
		return false
	}

	if len(parsed.(string)) == 0 {
		return false
	}
	return true
}

func (m *GoogleProvider) GetAuthDataFromJson(data io.Reader) string {
	l4g.Debug("Entering GetAuthDataFromJson()")
	l4g.Debug("GetAuthDataFromJson %s", data)

	gu := jsonDataDecode(data)
	if gu == nil {
		l4g.Debug("\t --> decoder failed with error: ")
		return ""
	}

	if len(gu.email) == 0 {
		l4g.Debug("GetAuthDataFromJson: len(gu.email) == 0 <-> IsValid() false")
		return ""
	}

	if gu.IsValid() {
		return gu.getAuthData()
	}

	return "000000000000000000001"
}
