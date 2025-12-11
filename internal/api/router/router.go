package router

import (
	"backend-go/internal/api/handler"
	adminHandler "backend-go/internal/api/handler/admin" // Statistics Handlers
	memberAdmin "backend-go/internal/api/handler/admin/member"
	payAdmin "backend-go/internal/api/handler/admin/pay"
	productHandler "backend-go/internal/api/handler/admin/product"
	promotionAdmin "backend-go/internal/api/handler/admin/promotion"
	tradeAdmin "backend-go/internal/api/handler/admin/trade"
	tradeBrokerageAdmin "backend-go/internal/api/handler/admin/trade/brokerage"
	memberHandler "backend-go/internal/api/handler/app/member"
	productApp "backend-go/internal/api/handler/app/product"
	promotionApp "backend-go/internal/api/handler/app/promotion"
	tradeApp "backend-go/internal/api/handler/app/trade"
	appBrokerage "backend-go/internal/api/handler/app/trade/brokerage"

	"backend-go/internal/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, rdb *redis.Client,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	tenantHandler *handler.TenantHandler,
	dictHandler *handler.DictHandler,
	deptHandler *handler.DeptHandler,
	postHandler *handler.PostHandler,
	roleHandler *handler.RoleHandler,
	menuHandler *handler.MenuHandler,
	permissionHandler *handler.PermissionHandler,
	noticeHandler *handler.NoticeHandler,
	configHandler *handler.ConfigHandler,
	smsChannelHandler *handler.SmsChannelHandler,
	smsTemplateHandler *handler.SmsTemplateHandler,
	smsLogHandler *handler.SmsLogHandler,
	fileConfigHandler *handler.FileConfigHandler,
	fileHandler *handler.FileHandler,
	appAuthHandler *memberHandler.AppAuthHandler,
	appMemberUserHandler *memberHandler.AppMemberUserHandler,
	appMemberAddressHandler *memberHandler.AppMemberAddressHandler,
	productCategoryHandler *productHandler.ProductCategoryHandler,
	productPropertyHandler *productHandler.ProductPropertyHandler,
	productBrandHandler *productHandler.ProductBrandHandler,
	productSpuHandler *productHandler.ProductSpuHandler,
	productCommentHandler *productHandler.ProductCommentHandler,
	productFavoriteHandler *productHandler.ProductFavoriteHandler,
	productBrowseHistoryHandler *productHandler.ProductBrowseHistoryHandler,
	appProductFavoriteHandler *productApp.AppProductFavoriteHandler,
	appProductBrowseHistoryHandler *productApp.AppProductBrowseHistoryHandler,
	appProductSpuHandler *productApp.AppProductSpuHandler,
	appProductCommentHandler *productApp.AppProductCommentHandler,
	appCartHandler *tradeApp.AppCartHandler,
	tradeOrderHandler *tradeAdmin.TradeOrderHandler,
	appTradeOrderHandler *tradeApp.AppTradeOrderHandler,
	tradeAfterSaleHandler *tradeAdmin.TradeAfterSaleHandler,
	appTradeAfterSaleHandler *tradeApp.AppTradeAfterSaleHandler,
	// Promotion
	couponHandler *promotionAdmin.CouponHandler,
	combinationActivityHandler *promotionAdmin.CombinationActivityHandler,
	discountActivityHandler *promotionAdmin.DiscountActivityHandler,
	appCombinationActivityHandler *promotionApp.AppCombinationActivityHandler,
	appCombinationRecordHandler *promotionApp.AppCombinationRecordHandler,
	appCouponHandler *promotionApp.AppCouponHandler,
	deliveryExpressHandler *tradeAdmin.DeliveryExpressHandler,
	deliveryPickUpStoreHandler *tradeAdmin.DeliveryPickUpStoreHandler,
	deliveryFreightTemplateHandler *tradeAdmin.DeliveryFreightTemplateHandler,
	bannerHandler *promotionAdmin.BannerHandler,
	rewardActivityHandler *promotionAdmin.RewardActivityHandler,
	seckillConfigHandler *promotionAdmin.SeckillConfigHandler,
	seckillActivityHandler *promotionAdmin.SeckillActivityHandler,
	bargainActivityHandler *promotionAdmin.BargainActivityHandler,
	appBannerHandler *promotionApp.AppBannerHandler,
	memberLevelHandler *memberAdmin.MemberLevelHandler,
	memberGroupHandler *memberAdmin.MemberGroupHandler,
	memberTagHandler *memberAdmin.MemberTagHandler,
	memberConfigHandler *memberAdmin.MemberConfigHandler,
	memberPointRecordHandler *memberAdmin.MemberPointRecordHandler,
	appMemberPointRecordHandler *memberHandler.AppMemberPointRecordHandler,
	memberSignInConfigHandler *memberAdmin.MemberSignInConfigHandler,
	memberSignInRecordHandler *memberAdmin.MemberSignInRecordHandler,
	appMemberSignInRecordHandler *memberHandler.AppMemberSignInRecordHandler,
	memberUserHandler *memberAdmin.MemberUserHandler,
	payAppHandler *payAdmin.PayAppHandler,
	payChannelHandler *payAdmin.PayChannelHandler,
	payOrderHandler *payAdmin.PayOrderHandler,
	payRefundHandler *payAdmin.PayRefundHandler,
	payNotifyHandler *payAdmin.PayNotifyHandler,
	loginLogHandler *handler.LoginLogHandler,
	operateLogHandler *handler.OperateLogHandler,
	jobHandler *handler.JobHandler,
	jobLogHandler *handler.JobLogHandler,
	apiAccessLogHandler *handler.ApiAccessLogHandler,
	apiErrorLogHandler *handler.ApiErrorLogHandler,
	socialClientHandler *handler.SocialClientHandler,
	socialUserHandler *handler.SocialUserHandler,
	sensitiveWordHandler *handler.SensitiveWordHandler,
	mailHandler *handler.MailHandler,
	notifyHandler *handler.NotifyHandler,
	oauth2ClientHandler *handler.OAuth2ClientHandler, // Added OAuth2ClientHandler
	appBargainActivityHandler *promotionApp.AppBargainActivityHandler,
	appBargainRecordHandler *promotionApp.AppBargainRecordHandler,
	appBargainHelpHandler *promotionApp.AppBargainHelpHandler,
	// Article
	articleCategoryHandler *promotionAdmin.ArticleCategoryHandler,
	articleHandler *promotionAdmin.ArticleHandler,
	appArticleHandler *promotionApp.AppArticleHandler,
	// DIY
	diyTemplateHandler *promotionAdmin.DiyTemplateHandler,
	diyPageHandler *promotionAdmin.DiyPageHandler,
	appDiyPageHandler *promotionApp.AppDiyPageHandler,
	// Kefu
	kefuHandler *promotionAdmin.KefuHandler,
	appKefuHandler *promotionApp.AppKefuHandler,
	// Trade Config
	tradeConfigHandler *tradeAdmin.TradeConfigHandler,
	appTradeConfigHandler *tradeApp.AppTradeConfigHandler,
	brokerageUserHandler *tradeBrokerageAdmin.BrokerageUserHandler,
	brokerageRecordHandler *tradeBrokerageAdmin.BrokerageRecordHandler,
	brokerageWithdrawHandler *tradeBrokerageAdmin.BrokerageWithdrawHandler,
	// Statistics
	tradeStatisticsHandler *adminHandler.TradeStatisticsHandler,
	productStatisticsHandler *adminHandler.ProductStatisticsHandler,
	memberStatisticsHandler *adminHandler.MemberStatisticsHandler,
	payStatisticsHandler *adminHandler.PayStatisticsHandler,
	appBrokerageUserHandler *appBrokerage.AppBrokerageUserHandler,
	appBrokerageRecordHandler *appBrokerage.AppBrokerageRecordHandler,
	appBrokerageWithdrawHandler *appBrokerage.AppBrokerageWithdrawHandler,
) *gin.Engine {
	// Debug log to confirm router init
	fmt.Println("Initializing Router...")
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())
	r.Use(gin.Logger())

	// 基础路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// ========== 模块化路由注册 ==========

	// System 模块 (Auth, Tenant, Dict, Dept, Post, User, Role, Permission, Logs, SMS, File, Infra)
	// System 模块 (Auth, Tenant, Dict, Dept, Post, User, Role, Permission, Logs, SMS, File, Infra)
	RegisterSystemRoutes(r,
		authHandler, userHandler, tenantHandler, dictHandler, deptHandler,
		postHandler, roleHandler, menuHandler, permissionHandler, noticeHandler,
		loginLogHandler, operateLogHandler, configHandler,
		smsChannelHandler, smsTemplateHandler, smsLogHandler,
		fileConfigHandler, fileHandler,
		jobHandler, jobLogHandler, apiAccessLogHandler, apiErrorLogHandler,
		socialClientHandler, socialUserHandler, sensitiveWordHandler, mailHandler, notifyHandler, oauth2ClientHandler,
	)

	// Product 模块
	RegisterProductRoutes(r,
		productCategoryHandler, productBrandHandler, productPropertyHandler,
		productSpuHandler, productCommentHandler, productFavoriteHandler,
		productBrowseHistoryHandler,
	)

	// Promotion 模块
	RegisterPromotionRoutes(r,
		couponHandler, bannerHandler, rewardActivityHandler,
		seckillConfigHandler, seckillActivityHandler, bargainActivityHandler,
		combinationActivityHandler, discountActivityHandler,
		articleCategoryHandler, articleHandler,
		diyTemplateHandler, diyPageHandler, kefuHandler,
	)

	// Trade 模块
	RegisterTradeRoutes(r,
		tradeOrderHandler, tradeAfterSaleHandler,
		deliveryExpressHandler, deliveryPickUpStoreHandler, deliveryFreightTemplateHandler,
		tradeConfigHandler,
		brokerageUserHandler,
		brokerageRecordHandler,
		brokerageWithdrawHandler,
	)

	// Member 模块 (Admin)
	RegisterMemberRoutes(r,
		memberSignInConfigHandler, memberSignInRecordHandler,
		memberPointRecordHandler,
		memberConfigHandler, memberGroupHandler, memberLevelHandler, memberTagHandler,
		memberUserHandler,
	)

	// Pay 模块
	RegisterPayRoutes(r,
		payAppHandler, payChannelHandler, payOrderHandler, payRefundHandler, payNotifyHandler,
	)

	// App 模块 (移动端)
	RegisterAppRoutes(r,
		// Member
		appAuthHandler, appMemberUserHandler, appMemberAddressHandler,
		appMemberPointRecordHandler, appMemberSignInRecordHandler,
		// Product
		appProductFavoriteHandler, appProductBrowseHistoryHandler,
		appProductSpuHandler, appProductCommentHandler,
		// Trade
		appCartHandler, appTradeOrderHandler, appTradeAfterSaleHandler, appTradeConfigHandler,
		// Promotion
		appCouponHandler, appBannerHandler, appArticleHandler, appDiyPageHandler, appKefuHandler,
		appCombinationActivityHandler, appCombinationRecordHandler,
		appBargainActivityHandler, appBargainRecordHandler, appBargainHelpHandler,
		appBrokerageUserHandler,
		appBrokerageRecordHandler,
		appBrokerageWithdrawHandler,
	)

	// Statistics 模块
	RegisterStatisticsRoutes(r,
		tradeStatisticsHandler, productStatisticsHandler,
		memberStatisticsHandler, payStatisticsHandler,
	)

	return r
}
