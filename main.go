package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/database"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"bwastartup/user"

	webHandler "bwastartup/web/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func main() {
	// // refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	// dsn := "root@tcp(localhost:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	// // dsn := "root:a@tcp(0.0.0.0:5555)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal(err.Error())
	// } else {
	// 	fmt.Println(db, "connection succes")
	// }
	db := database.ConnectToDB()
	database.AutoMigrate(db)

	campaignRepository := campaign.NewRepository(db)
	userRepository := user.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	paymentService := payment.NewPaymentService()
	transactionService := transaction.NewServiceTransaction(transactionRepository, campaignRepository, paymentService)
	UserService := user.NewSevice(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authservice := auth.Newjwtservice()

	campaign, _ := campaignService.FindCampaigns(1)
	fmt.Println(campaign)

	fmt.Println(campaign, "camm")

	// test crte transaction
	// user, _ := UserService.GetUserById(43)
	// input := transaction.CreateTransactionInput{
	// 	CampaignID: 12,
	// 	User:       user,
	// 	Amount:     1000000,
	// }
	// transactionService.CreateTransaction(input)
	// //

	token, err := authservice.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo1fQ.KZkpB7JjhcG-J7DbexpFgFPkT1pLlw28t7sndIG66sM")
	fmt.Println(token, "token err", err)
	if err != nil {
		fmt.Println("error")
		fmt.Println("error")
		fmt.Println("error")
	}
	if token.Valid {
		fmt.Println("valid")
		fmt.Println("valid")
		fmt.Println("valid")
	} else {
		fmt.Println("inva")
		fmt.Println("inva")
		fmt.Println("inva")
	}

	transactionHandler := handler.NewTransactionHandler(transactionService)
	userHandler := handler.NewHandler(UserService, authservice)
	campainHandler := handler.NewCampaignHandler(campaignService)

	userWebHandler := webHandler.NewUserHandler(UserService)
	campaignWebHandler := webHandler.NewCampaignHandler(campaignService, UserService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	SessionWebHandler := webHandler.NewSessionHandler(UserService)

	router := gin.Default()

	store := cookie.NewStore(auth.SECRET_KEY)
	router.Use(sessions.Sessions("mysession", store))
	router.Use(cors.Default())
	router.HTMLRender = loadTemplates("./web/templates/")
	router.Static("/images", "./images")
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/webfonts", "./web/assets/webfonts")

	api := router.Group("api/v1")

	// api.POST("/users", userHandler.RegisterUser)
	// api.POST("/login", userHandler.LoginUser)
	// api.POST("/checkemail", userHandler.HandlerEmailAvailability)
	// api.POST("/uploadavatar", authMidleware(authservice, UserService), userHandler.UploadAvatar)
	// api.POST("/user/fetch", authMidleware(authservice, UserService), userHandler.FetchUser)

	// api.GET("/campains", campainHandler.GetCampaigns)
	// api.POST("/campains", authMidleware(authservice, UserService), campainHandler.CreateCampaign)
	// api.PUT("/campains/:id", authMidleware(authservice, UserService), campainHandler.UpdateCampaign)
	// api.GET("/campain/:id", campainHandler.GetCampaign)
	// api.POST("/uploadcampainimages", authMidleware(authservice, UserService), campainHandler.UploadCampaignImage)

	// api.GET("/campaigns/:id/transaction", authMidleware(authservice, UserService), transactionHandler.GetCampaignTransactions)
	// api.GET("/transaction", authMidleware(authservice, UserService), transactionHandler.GetUserTransactions)
	// api.POST("/transaction", authMidleware(authservice, UserService), transactionHandler.CreateTransaction)
	// api.POST("/transaction/notification", transactionHandler.GetNotif)

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.LoginUser)
	api.POST("/email_checkers", userHandler.HandlerEmailAvailability)
	api.POST("/avatars", authMidleware(authservice, UserService), userHandler.UploadAvatar)
	api.POST("/users/fetch", authMidleware(authservice, UserService), userHandler.FetchUser)

	api.GET("/campaigns", campainHandler.GetCampaigns)
	api.POST("/campaigns", authMidleware(authservice, UserService), campainHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMidleware(authservice, UserService), campainHandler.UpdateCampaign)
	api.GET("/campaigns/:id", campainHandler.GetCampaign)
	api.POST("/campaign-images", authMidleware(authservice, UserService), campainHandler.UploadCampaignImage)

	api.GET("/campaigns/:id/transactions", authMidleware(authservice, UserService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMidleware(authservice, UserService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMidleware(authservice, UserService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotif)

	router.GET("/users", userWebHandler.Index)
	router.GET("/users/new", userWebHandler.CreateView)
	router.POST("/users", userWebHandler.CreateUser)
	router.GET("/users/edit/:id", userWebHandler.UpdateUser)
	router.POST("/users/update/:id", userWebHandler.PostUpdateUser)
	router.GET("/users/avatar/:id", userWebHandler.UploadAvatarView)
	router.POST("/users/avatar/:id", userWebHandler.UploadAvatarSubmit)

	router.GET("/campaigns", sessionhMidleware(), campaignWebHandler.Index)
	router.GET("/campaigns/new", campaignWebHandler.CreateView)
	router.POST("/campaigns", campaignWebHandler.CreateCampaign)
	// router.POST("/campaigns", campaignWebHandler.CreateCampaign)
	router.GET("/campaigns/image/:id", campaignWebHandler.UploadImageView)
	router.POST("/campaigns/image/:id", campaignWebHandler.UploadImage)
	router.GET("/campaigns/edit/:id", campaignWebHandler.EditCampaign)
	router.POST("/campaigns/update/:id", campaignWebHandler.PostUpdateCampaign)
	router.GET("/campaigns/show/:id", campaignWebHandler.ShowDetailCampaign)

	router.GET("/transactions", transactionWebHandler.GetAllTransaction)

	router.GET("/login", SessionWebHandler.NewSession)
	router.POST("/login", SessionWebHandler.PostSession)
	router.GET("/logout", SessionWebHandler.Destroy)

	router.Run(":8080")

}

func sessionhMidleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)
		userLogin := session.Get("userId")
		if userLogin == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

	}
}

func authMidleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unathorized", http.StatusUnauthorized, "unauthorized", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		fmt.Println("proses authHeader", authHeader)
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		fmt.Println("proses arraytoken", arrayToken, tokenString, len(arrayToken))
		fmt.Println("proses arraytoken", arrayToken[0], arrayToken[1])
		if len(arrayToken) == 2 {
			fmt.Println("proses arraytoken", arrayToken[0], arrayToken[1], "2")
			tokenString = arrayToken[1]
		}
		token, err := authService.ValidateToken(tokenString)
		fmt.Println("proses token", tokenString, token, err)
		if err != nil {
			fmt.Println("proses token", tokenString, token, err, "masuk sini kah")
			response := helper.APIResponse("Unathorized", http.StatusUnauthorized, "unauthorized", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)
		fmt.Println("proses claim", claim)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unathorized", http.StatusUnauthorized, "unauthorized", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userId := int(claim["user_id"].(float64))

		user, err := userService.GetUserById(userId)
		if err != nil {
			response := helper.APIResponse("Unathorized", http.StatusUnauthorized, "unauthorized", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("curresntuser", user)
	}
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
