package api

import (
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fankserver/fankserver-cli/api/admin"
	"github.com/fankserver/fankserver-cli/config"
	"github.com/fankserver/fankserver-cli/connection"
	"github.com/fankserver/fankserver-cli/models"

	"gopkg.in/kataras/iris.v6"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	cli "gopkg.in/urfave/cli.v2"
	"gopkg.in/mgo.v2/bson"
)

func Listen(ctx *cli.Context) error {
	err := config.ReadConfigFile(ctx.String("config"))
	if err != nil {
		return err
	}

	// Setup mongodb
	Db := connection.MongoDB{}
	Db.Init()

	// Setup model indizes
	models.Init()

	app := iris.New()
	app.Adapt(iris.DevLogger())
	app.Adapt(
		httprouter.New(),
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		}),
	)

	auth := app.Party("/auth", func(ctx *iris.Context) {
		ctx.Next()
	})
	NewAuthAPI(auth)

	api := app.Party("/api", jwtMiddleware().Serve)
	{
		apiAdmin := api.Party("/admin", userHasTags("admin"))
		{
			admin.NewTag(apiAdmin.Party("/tag"))
			admin.NewUser(apiAdmin.Party("/user"))
		}
		api.Get("/", func(ctx *iris.Context) {
			user := ctx.Get("jwt").(*jwt.Token)

			ctx.Writef("This is an authenticated request\n")
			ctx.Writef("Claim content:\n")

			ctx.Writef("%s\n%s", user.Signature, user.Claims.(jwt.MapClaims))
		})
	}

	app.Listen(fmt.Sprintf("%s:%d", ctx.String("interface"), ctx.Uint("port")))
	return nil
}

func userHasTags(tag string) iris.HandlerFunc{
	return func(ctx *iris.Context) {
		jwtUser := ctx.Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)["user"]
		db := connection.MongoDB{}
		db.Init()
		defer db.Close()

		operations := []bson.M{
			models.MatchById(jwtUser.(map[string]interface{})["id"].(string)),
			models.UserLookupTags,
			bson.M{
				"$match": bson.M{
					"tags.name": tag,
				},
			},
		}

		var user *models.User
		err := db.C(models.UserCollection).Pipe(operations).One(&user)
		if err == mgo.ErrNotFound || user == nil {
			ctx.EmitError(iris.StatusForbidden)
			return
		} else if err != nil {
			log.Println(err)
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}

		ctx.Next()
	}
}
