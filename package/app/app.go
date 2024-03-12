package app

import (
	"encoding/json"

	db "272-backend/package/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/swagger"
)

/* type FormT struct {
	ID       string   `json:"id,omitempty" bson:"_id,omitempty"`
	Username string   `json:"username" bson:"username"`
	Password string   `json:"password" bson:"password"`
	Roles    []string `json:"roles" bson:"roles"`
	Token    string   `json:"token" bson:"token"`
} */

var (
	App          *fiber.App
	SessionStore *session.Store
)

func init() {
	SessionStore = session.New()
	App = fiber.New()
	App.Use(
		cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept",
		}),
		logger.New(),
		/* getSession, */
	)
	App.Get("/swagger/*", swagger.HandlerDefault) // default

	/* App.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://127.0.0.1:5000/swagger/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	})) */

}

/*
func getSession(c *fiber.Ctx) error {
	sess, err := SessionStore.Get(c)
	if err != nil {
		log.Println(err)
	}
	var userDoc FormT
	token := sess.Get("token")
	if token == nil {
		return c.Next()
	}
	userDocx := sess.Get("userDoc")
	if userDocx == nil {
		docString, _ := db.Redis.Get(token.(string))
		if len(docString) == 0 {
			if _doc, err := GetAuthDocumentByToken(token.(string)); err != nil {
				log.Println("parseError: | GetAuthDocumentByToken" + err.Error())
			} else {
				sess.Set("userDoc", _doc)
				if err := sess.Save(); err != nil {
					log.Println("parseError: | session" + err.Error())
					log.Println(err)
				}
				out, _ := json.Marshal(_doc)
				if err := db.Redis.Set(_doc.Token, string(out)); err != nil {
					log.Println("parseError: | Redis" + err.Error())
					log.Println(err)
				}
			}
		} else {
			if err := json.Unmarshal([]byte(docString), &userDoc); err != nil {
				log.Println("parseError: " + err.Error())
			}
			sess.Set("userDoc", docString)
			if err := sess.Save(); err != nil {
				log.Println("error: " + err.Error())
			}
		}
		return c.Next()
	}
	if err := json.Unmarshal([]byte(userDocx.(string)), &userDoc); err != nil {
		log.Println("parseError: " + err.Error())
	}
	return c.Next()
}
*/
/* func isAuthorized(c *fiber.Ctx) error {
	sess, err := SessionStore.Get(c)
	if err != nil {
		log.Println(err)
	}
	token := sess.Get("token")
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	return c.Next()
}
*/
/*
func GetAuthDocumentByToken(token string) (FormT, error) {
	var auth FormT
	log.Println("App | GetAuthDocumentByToken | token: " + token + "")
	if result, err := db.Redis.Get(auth.Token); err == nil {
		if err := json.Unmarshal([]byte(result), &auth); err != nil {
			log.Println("App | GetAuthDocumentByToken | Unmarshal: " + err.Error())
			return auth, err
		}
		return auth, nil
	}
	if err := db.Auths.FindOne(context.TODO(), bson.D{
		{Key: "token", Value: token},
	}).Decode(&auth); err != nil {
		log.Println("App | GetAuthDocumentByToken: " + err.Error())
		return auth, err
	}
	return auth, nil
}
*/

func GetRoles(c *fiber.Ctx) []string {
	sess, err := SessionStore.Get(c)
	if err != nil {
		return []string{}
	}
	token := sess.Get("token")
	rolesString, err := db.Redis.Get(token.(string))
	if err != nil {
		return []string{}
	}
	var roles []string
	err = json.Unmarshal([]byte(rolesString), &roles)
	if err != nil {
		return []string{}
	}
	return roles
}
