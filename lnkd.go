package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	jsn "encoding/json"
	"fmt"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ostafen/clover"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var cloverDb *clover.DB

// This key is overwritten on startup by an ENV var, or an auto generated random string.
var jwtSecretKey = []byte("my_secret_key")
var jwtValidator *validator.Validator

var siteUrl = "https://drws.dev/"

var disallowedRoutes DisallowedRoutes

func main() {
	initDb()
	setupEnv()
	setupJwtValidation()
	setupDisallowedRoutes()
	router := gin.Default()
	authorized := router.Group("")
	authorized.Use(jwtAuthMiddleware())

	// Serve the frontend.
	router.StaticFile("/", "./lnkd-front/build/200.html")
	router.Static("/_app", "./lnkd-front/build/_app")

	// If there is a saved link, forward the user to the loaded redirect url.
	// Otherwise, send them to the home page.
	router.GET("/:pathVar", func(c *gin.Context) {
		name := c.Param("pathVar")
		name = siteUrl + name
		loadedRecord := LinkRecord{}
		links, _ := cloverDb.Query("linkRedirects").Where(clover.Field("lnkd_url").Eq(name)).FindAll()
		// If no links found redirect to app homepage.
		if len(links) < 1 {
			c.Redirect(http.StatusFound, "/")
			return
		}
		// If there is an error send user to homepage.
		err := links[0].Unmarshal(&loadedRecord)
		if err != nil {
			c.Redirect(http.StatusFound, "/")
			return
		}
		updates := make(map[string]interface{})
		updates["hits"] = loadedRecord.HitCount + 1

		query := cloverDb.Query("linkRedirects").Where(clover.Field("lnkd_url").Eq(loadedRecord.LnkdUrl))
		err = query.Update(updates)
		if err != nil {
			fmt.Println("Internal Error: Failed to update record hit count", loadedRecord)
		}
		// Redirect user to saved url.
		c.Redirect(http.StatusFound, loadedRecord.RedirectUrl)
	})
	// Default behavior is to allow anonymous link shortening.
	// However, they can not update another users links.
	router.POST("/api/lnkd", func(c *gin.Context) {
		var json LinkRecord
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		lnkdUrl := siteUrl + json.LnkdUrl
		links, _ := cloverDb.Query("linkRedirects").Where(clover.Field("lnkd_url").Eq(lnkdUrl)).FindAll()
		// Check if there is an existing record. If so only allow an authorized user to update it.
		claims := jwtGetClaims(c, c.GetHeader("Authorization"))
		var customClaims *CustomClaims
		if claims == nil {
			json.Username = ""
		} else {
			customClaims = claims.(*CustomClaims)
			if customClaims.Username != "" {
				json.Username = customClaims.Username
			}
		}
		if len(links) > 0 {
			loadedRecord := LinkRecord{}
			err := links[0].Unmarshal(&loadedRecord)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if customClaims.Username == loadedRecord.Username {
				// Actually handle edit.

				updates := make(map[string]interface{})
				updates["redirect_url"] = json.RedirectUrl

				query := cloverDb.Query("linkRedirects").Where(clover.Field("lnkd_url").Eq(loadedRecord.LnkdUrl).And(clover.Field("username").Eq(loadedRecord.Username)))
				err := query.Update(updates)
				if err != nil {
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "Lnkd object updated:" + loadedRecord.LnkdUrl, "lnkd": lnkdUrl})
				return
			}
			fmt.Println(customClaims.Username)

			c.JSON(http.StatusBadRequest, gin.H{"error": "Attempt to edit and existing link without proper authentication"})
			return
		}
		// Check if link being saved has an associated user. If so only allow authorized user.
		newId := ""
		if json.Username != "" && json.LnkdUrl != "" {
			claims := jwtGetClaims(c, c.GetHeader("Authorization"))
			if claims == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Attempt create a link without proper authentication"})
				return
			}
			newId = json.LnkdUrl
		} else {
			newId = generateNewHash(8)
		}

		docs, _ := cloverDb.Query("linkRedirects").Where(clover.Field("lnkd_url").Eq(newId)).FindAll()
		// The odds of a new id being generated that already exists are very low.
		// However, this should handle creating a new ID in that event.
		for len(docs) > 0 {
			newId = generateNewHash(8)
			docs, _ = cloverDb.Query("linkRedirects").Where(clover.Field("lnkd_url").Eq(newId)).FindAll()
		}
		json.LnkdUrl = newId

		lnkdUrl, err := createNewLnkdRecord(json)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating lnkd object."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "New lnkd object created", "lnkd": lnkdUrl})
	})
	router.POST("/login", func(c *gin.Context) {
		var json LoginRequest
		loadedUser := User{}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		users, _ := cloverDb.Query("users").Where(clover.Field("Username").Eq(json.Username)).FindAll()
		// Ensure only one user was found.
		if len(users) != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "110: Error loading user."})
			return
		}
		// Load user, then check password matches expected.
		err := users[0].Unmarshal(&loadedUser)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "120: Error loading user."})
			return
		}
		pwdMatch := comparePasswords(loadedUser.Password, []byte(json.Password))
		if pwdMatch {
			usersJwt, success := createUserJwt(loadedUser)
			if success {
				c.JSON(http.StatusOK, gin.H{"user": usersJwt})
				return
			}
		}
		// Default out with a generic error in the event login failed for any reason.
		c.JSON(http.StatusBadRequest, gin.H{"status": "Log in failed, check username or password"})
	})
	authorized.POST("/api/user", func(c *gin.Context) {
		var json User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		claims := jwtGetClaims(c, c.GetHeader("Authorization"))
		if claims == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access Denied"})
			return
		}
		customClaims := claims.(*CustomClaims)
		if customClaims.Roles != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have the required permissions"})
		}
		if !isUsernameAllowed(json.Username) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "122: Username Taken."})
			return
		}
		userName, err := createUser(json)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "122: Error creating user."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": userName})
	})

	authorized.PATCH("/api/user", func(c *gin.Context) {
		var json User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		claims := jwtGetClaims(c, c.GetHeader("Authorization"))
		if claims == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access Denied"})
			return
		}
		customClaims := claims.(*CustomClaims)
		if customClaims.Username != json.Username && customClaims.Roles != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have the required permissions"})
		}
		userName, err := editUser(json)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "122: Error editing user."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": userName})
	})

	authorized.GET("/api/lnkd", func(c *gin.Context) {
		showAll := c.Query("showAll")
		claims := jwtGetClaims(c, c.GetHeader("Authorization"))
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}
		customClaims := claims.(*CustomClaims)
		docs, _ := cloverDb.Query("linkRedirects").Where(clover.Field("username").Eq(customClaims.Username)).FindAll()

		// If admin allow listing of all LNKDs.
		if strings.Contains(customClaims.Roles, "admin") && showAll == "true" {
			docs, _ = cloverDb.Query("linkRedirects").FindAll()
		}
		var list LinkListing
		for _, doc := range docs {
			loadedRecord := LinkRecord{}
			err := doc.Unmarshal(&loadedRecord)
			if err != nil {
				fmt.Println(err)
			}
			list.Links = append(list.Links, loadedRecord)
		}
		c.JSON(http.StatusOK, gin.H{"links": list.Links})
	})

	router.GET("/api/lnkd/:route", func(c *gin.Context) {
		name := c.Param("route")
		allowed := isRouteAllowed(name)
		c.JSON(http.StatusOK, gin.H{"lnkd_allowed": allowed})
	})

	authorized.GET("/api/user", func(c *gin.Context) {
		docs, _ := cloverDb.Query("users").FindAll()
		claims := jwtGetClaims(c, c.GetHeader("Authorization"))
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Attempt to edit and existing link without proper authentication"})
			return
		}
		customClaims := claims.(*CustomClaims)
		if customClaims.Roles != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You do not have the required permissions"})
		}
		var list UserListing
		for _, doc := range docs {
			loadedUser := User{}
			err := doc.Unmarshal(&loadedUser)
			if err != nil {
				fmt.Println(err)
			}
			loadedUser.Password = ""
			list.Users = append(list.Users, loadedUser)
		}
		c.JSON(http.StatusOK, gin.H{"users": list.Users})
	})

	authorized.DELETE("/api/user/:username", func(c *gin.Context) {
		userName := c.Param("username")
		claims := jwtGetClaims(c, c.GetHeader("Authorization"))
		if claims == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Attempt to edit and existing link without proper authentication"})
			return
		}
		customClaims := claims.(*CustomClaims)
		if customClaims.Roles != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have the required permissions"})
		}
		docs, _ := cloverDb.Query("users").Where(clover.Field("Username").Eq(userName)).FindAll()
		if len(docs) != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "User not found, or user does not have access to delete link."})
			return
		}
		loadedUser := User{}
		err := docs[0].Unmarshal(&loadedUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "User deletion failed."})
			return
		}
		err = deleteUser(loadedUser.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "User deletion failed."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_deleted": loadedUser.Username})
	})

	authorized.DELETE("/api/lnkd/:link", func(c *gin.Context) {
		linkId := c.Param("link")

		claims := jwtGetClaims(c, c.GetHeader("Authorization"))
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Attempt to edit and existing link without proper authentication"})
			return
		}
		customClaims := claims.(*CustomClaims)

		docs, _ := cloverDb.Query("linkRedirects").Where(clover.Field("_id").Eq(linkId).And(clover.Field("username").Eq(customClaims.Username))).FindAll()
		if strings.Contains(customClaims.Roles, "admin") {
			docs, _ = cloverDb.Query("linkRedirects").Where(clover.Field("_id").Eq(linkId)).FindAll()
		}
		if len(docs) != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Link not found, or user does not have access to delete link."})
			return
		}
		err := deleteLnkdRecord(linkId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Link deletion failed."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"links_deleted": linkId})
	})

	err := router.Run(":8070")
	if err != nil {
		log.Fatal(err)
	}
}

/*
	Handle any DB Errors.
*/
func handleDbErrors(errorStr string) {
	fmt.Println(errorStr)
	//log.Fatal(errorStr)
}

func initDb() {
	dbPath := "clover-db"
	if _, err := os.Stat("/data/lnkd"); !os.IsNotExist(err) {
		dbPath = "/data/lnkd"
	}
	cloverDb, _ = clover.Open(dbPath)

	hasCollection, err := cloverDb.HasCollection("linkRedirects")
	if err != nil {
		handleDbErrors("Error checking for linkRedirects collection.")
	}
	if !hasCollection {
		err := cloverDb.CreateCollection("linkRedirects")
		if err != nil {
			handleDbErrors("Error creating linkRedirects collection")
		}
	}
	hasUserCollection, err := cloverDb.HasCollection("users")
	if err != nil {
		handleDbErrors("Error checking for users collection.")
	}
	if !hasUserCollection {
		err := cloverDb.CreateCollection("users")
		createInitialUser()
		if err != nil {
			handleDbErrors("Error creating users collection")
		}
	}
}

func (c *CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func setupJwtValidation() {
	// Creates a new random key.
	// This is more secure than using the same key, however users will be logged out if the server reboots.
	// Note this can also cause issues if the service is load balanced and don't get the same backend everytime.
	keyFunc := func(ctx context.Context) (interface{}, error) {
		// Our token must be signed using this data.
		return jwtSecretKey, nil
	}

	customClaims := func() validator.CustomClaims {
		return &CustomClaims{}
	}

	jValidator, err := validator.New(
		keyFunc,
		validator.HS256,
		siteUrl,
		[]string{"lnkd"},
		validator.WithCustomClaims(customClaims),
	)
	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}
	jwtValidator = jValidator
}

func setupEnv() {
	viper.SetEnvPrefix("lnkd")
	err := viper.BindEnv("url")
	if err != nil {
		log.Println(err)
	}
	err = viper.BindEnv("jwtkey")
	if err != nil {
		log.Println(err)
	}
	if viper.GetString("jwtkey") != "" {
		jwtSecretKey = []byte(viper.GetString("jwtkey"))
	} else {
		generateJwtSecret()
	}
	fmt.Println(viper.Get("url"))
	if viper.GetString("url") != "" {
		siteUrl = viper.GetString("url")
	}
}

func createUserJwt(loadedUser User) (string, bool) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loadedUser.Username,
		"iss":      siteUrl,
		"aud":      "lnkd",
		"roles":    loadedUser.Roles,
	})

	// Sign and get the complete encoded token as a string using the secret.
	tokenString, err := token.SignedString(jwtSecretKey)

	if err != nil {
		return "", false
	}
	return tokenString, true
}

func jwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtClaims, err := jwtValidator.ValidateToken(c, c.GetHeader("Authorization"))
		if err != nil {
			respondWithError(401, "Invalid JWT", c)
			fmt.Println(err)
			return
		}
		claims := jwtClaims.(*validator.ValidatedClaims)
		if claims.RegisteredClaims.Audience[0] != "lnkd" {
			respondWithError(401, "Invalid JWT claims", c)
			return
		}
	}
}

func jwtGetClaims(ctx context.Context, authHeader string) validator.CustomClaims {
	jwtClaims, err := jwtValidator.ValidateToken(ctx, authHeader)
	if err != nil {
		return nil
	}
	claims := jwtClaims.(*validator.ValidatedClaims)
	return claims.CustomClaims
}

func createInitialUser() {
	defaultPassword := randomString(20)
	h := sha256.New()
	h.Write([]byte(defaultPassword))
	defaultPassword256 := hex.EncodeToString(h.Sum(nil))
	fmt.Println("Initial user created. Email:admin Password:" + defaultPassword)
	fmt.Println("Change this password")

	newUser := User{}
	newUser.Password = hashAndSalt([]byte(defaultPassword256))
	newUser.Status = 1
	newUser.Email = "admin@example.com"
	newUser.Username = "admin"
	newUser.Roles = "admin"
	doc := clover.NewDocumentOf(newUser)
	newId, err := cloverDb.InsertOne("users", doc)
	fmt.Println(newId)
	fmt.Println(err)
}

func createUser(userData User) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	userData.Username = reg.ReplaceAllString(userData.Username, "")

	newUser := User{}
	newUser.Password = hashAndSalt([]byte(userData.Password))
	newUser.Status = 1
	newUser.Email = userData.Email
	newUser.Username = userData.Username
	newUser.Roles = ""
	doc := clover.NewDocumentOf(newUser)
	_, err = cloverDb.InsertOne("users", doc)
	if err != nil {
		return "", err
	}
	return userData.Username, err
}

func editUser(userData User) (string, error) {

	userData.Password = hashAndSalt([]byte(userData.Password))
	updates := make(map[string]interface{})
	updates["password"] = userData.Password
	updates["email"] = userData.Email

	query := cloverDb.Query("users").Where(clover.Field("Username").Eq(userData.Username))
	err := query.Update(updates)

	return userData.Username, err
}

func createNewLnkdRecord(lnkd LinkRecord) (string, error) {
	link := clover.NewDocument()
	// Remove most special chars.
	reg, err := regexp.Compile("[^a-zA-Z0-9-_]+")
	if err != nil {
		return "", err
	}
	processedString := reg.ReplaceAllString(lnkd.LnkdUrl, "")

	lnkdUrl := siteUrl + processedString
	disallowedRoutes.Routes = append(disallowedRoutes.Routes, processedString)
	link.Set("lnkd_url", lnkdUrl)
	link.Set("redirect_url", lnkd.RedirectUrl)
	link.Set("hits", lnkd.HitCount)
	link.Set("username", lnkd.Username)

	_, err = cloverDb.InsertOne("linkRedirects", link)

	return lnkdUrl, err
}

func deleteLnkdRecord(lnkdId string) error {
	err := cloverDb.Query("linkRedirects").Where(clover.Field("_id").Eq(lnkdId)).Delete()
	return err
}

func deleteUser(lnkdId string) error {
	err := cloverDb.Query("users").Where(clover.Field("_id").Eq(lnkdId)).Delete()
	return err
}

func setupDisallowedRoutes() {
	routeFile, err := os.Open("./disallowed_routes.json")
	if err != nil {
		fmt.Println(err)
	}
	defer routeFile.Close()

	byteValue, _ := ioutil.ReadAll(routeFile)
	err = jsn.Unmarshal(byteValue, &disallowedRoutes)
	if err != nil {
		return
	}

	links, _ := cloverDb.Query("linkRedirects").FindAll()
	for _, link := range links {
		loadedRecord := LinkRecord{}
		err := link.Unmarshal(&loadedRecord)
		if err != nil {
			fmt.Println(err)
		}
		route := loadedRecord.LnkdUrl
		strings.ReplaceAll(route, siteUrl, "")
		disallowedRoutes.Routes = append(disallowedRoutes.Routes, route)
	}
}

func isRouteAllowed(lnkdUrl string) bool {
	found := false
	for _, element := range disallowedRoutes.Routes {
		i := strings.Index(lnkdUrl, element)
		if i == -1 || i >= 2 {
			found = false
		} else {
			return false
		}
	}
	return !found
}

func isUsernameAllowed(userName string) bool {
	users, _ := cloverDb.Query("users").Where(clover.Field("Username").Eq(userName)).FindAll()
	// Ensure only one user was found.
	if len(users) != 0 {
		return false
	}
	return true
}
