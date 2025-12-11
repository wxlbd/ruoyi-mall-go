//go:build wireinject
// +build wireinject

package main

import (
	"backend-go/internal/api/handler"
	adminHandler "backend-go/internal/api/handler/admin" // Statistics Handlers
	memberAdmin "backend-go/internal/api/handler/admin/member"
	payAdmin "backend-go/internal/api/handler/admin/pay"
	productHandler "backend-go/internal/api/handler/admin/product"
	promotionAdmin "backend-go/internal/api/handler/admin/promotion"
	tradeAdmin "backend-go/internal/api/handler/admin/trade"
	brokerage "backend-go/internal/api/handler/admin/trade/brokerage"
	memberHandler "backend-go/internal/api/handler/app/member"
	productApp "backend-go/internal/api/handler/app/product"
	promotionApp "backend-go/internal/api/handler/app/promotion"
	tradeApp "backend-go/internal/api/handler/app/trade"
	appBrokerage "backend-go/internal/api/handler/app/trade/brokerage"
	"backend-go/internal/api/router"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo"
	productRepo "backend-go/internal/repo/product" // Product Statistics Repo
	"backend-go/internal/service"
	memberSvc "backend-go/internal/service/member"
	paySvc "backend-go/internal/service/pay"
	"backend-go/internal/service/pay/client"
	_ "backend-go/internal/service/pay/client/alipay"
	_ "backend-go/internal/service/pay/client/weixin"
	product "backend-go/internal/service/product"
	promotionSvc "backend-go/internal/service/promotion"
	tradeSvc "backend-go/internal/service/trade"
	tradeBrokerageSvc "backend-go/internal/service/trade/brokerage"
	deliveryClient "backend-go/internal/service/trade/delivery/client"

	"backend-go/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitApp() (*gin.Engine, error) {
	wire.Build(
		core.InitDB,
		core.InitRedis,
		logger.NewLogger,
		// Repo (GORM Gen)
		repo.NewQuery,
		// Service
		service.NewOAuth2TokenService,
		service.NewAuthService,
		service.NewMenuService,
		service.NewRoleService,
		service.NewPermissionService,
		service.NewTenantService,
		service.NewUserService,
		service.NewDictService,
		service.NewDeptService,
		service.NewPostService,
		service.NewNoticeService,
		service.NewConfigService,
		service.NewSmsClientFactory,            // Added SmsClientFactory
		service.NewSmsChannelService,           // Added SmsChannelService
		service.NewSmsTemplateService,          // Added SmsTemplateService
		service.NewSmsLogService,               // Added SmsLogService
		service.NewFileConfigService,           // Added FileConfigService
		service.NewFileService,                 // Added FileService
		service.NewSmsCodeService,              // Added SmsCodeService
		service.NewLoginLogService,             // Added LoginLogService
		service.NewOperateLogService,           // Added OperateLogService
		service.NewScheduler,                   // Added Scheduler
		service.NewJobService,                  // Added JobService
		service.NewJobLogService,               // Added JobLogService
		service.NewApiAccessLogService,         // Added ApiAccessLogService
		service.NewApiErrorLogService,          // Added ApiErrorLogService
		service.NewSocialClientService,         // Added SocialClientService
		service.NewSocialUserService,           // Added SocialUserService
		service.NewSensitiveWordService,        // Added SensitiveWordService
		service.NewMailService,                 // Added MailService
		service.NewNotifyService,               // Added NotifyService
		service.NewOAuth2ClientService,         // Added OAuth2ClientService
		memberSvc.NewMemberAuthService,         // Added MemberAuthService
		memberSvc.NewMemberUserService,         // Added MemberUserService
		memberSvc.NewMemberAddressService,      // Added MemberAddressService
		memberSvc.NewMemberLevelService,        // Added MemberLevelService
		memberSvc.NewMemberGroupService,        // Added MemberGroupService
		memberSvc.NewMemberTagService,          // Added MemberTagService
		memberSvc.NewMemberConfigService,       // Added MemberConfigService
		memberSvc.NewMemberPointRecordService,  // Added MemberPointRecordService
		product.NewProductCategoryService,      // Added ProductCategoryService
		product.NewProductPropertyService,      // Added ProductPropertyService
		product.NewProductPropertyValueService, // Added ProductPropertyValueService
		product.NewProductBrandService,         // Added ProductBrandService
		product.NewProductSkuService,           // Added ProductSkuService
		product.NewProductSpuService,           // Added ProductSpuService
		product.NewProductCommentService,       // Added ProductCommentService
		product.NewProductFavoriteService,      // Added ProductFavoriteService
		product.NewProductBrowseHistoryService, // Added ProductBrowseHistoryService
		// Member Sign-in
		memberSvc.NewMemberSignInConfigService,
		memberSvc.NewMemberSignInRecordService,

		// Handler
		handler.NewAuthHandler,
		handler.NewUserHandler,
		handler.NewTenantHandler,
		handler.NewDictHandler,
		handler.NewDeptHandler,
		handler.NewPostHandler,
		handler.NewRoleHandler,
		handler.NewMenuHandler, // Added MenuHandler
		handler.NewPermissionHandler,
		handler.NewNoticeHandler,
		handler.NewConfigHandler,
		handler.NewSmsChannelHandler,    // Added SmsChannelHandler
		handler.NewSmsTemplateHandler,   // Added SmsTemplateHandler
		handler.NewSmsLogHandler,        // Added SmsLogHandler
		handler.NewFileConfigHandler,    // Added FileConfigHandler
		handler.NewFileHandler,          // Added FileHandler
		memberHandler.NewAppAuthHandler, // Added AppAuthHandler
		handler.NewLoginLogHandler,      // Added LoginLogHandler
		handler.NewOperateLogHandler,    // Added OperateLogHandler
		handler.NewJobHandler,           // Added JobHandler
		handler.NewJobLogHandler,        // Added JobLogHandler
		handler.NewApiAccessLogHandler,  // Added ApiAccessLogHandler
		handler.NewApiErrorLogHandler,   // Added ApiErrorLogHandler
		handler.NewSocialClientHandler,  // Added SocialClientHandler
		handler.NewSocialUserHandler,    // Added SocialUserHandler
		handler.NewSensitiveWordHandler, // Added SensitiveWordHandler
		handler.NewMailHandler,          // Added MailHandler
		handler.NewNotifyHandler,        // Added NotifyHandler
		handler.NewOAuth2ClientHandler,  // Added OAuth2ClientHandler
		// Member
		memberAdmin.NewMemberLevelHandler,             // Added MemberLevelHandler for admin
		memberAdmin.NewMemberGroupHandler,             // Added MemberGroupHandler for admin
		memberAdmin.NewMemberTagHandler,               // Added MemberTagHandler for admin
		memberAdmin.NewMemberConfigHandler,            // Added MemberConfigHandler for admin
		memberAdmin.NewMemberPointRecordHandler,       // Added MemberPointRecordHandler for admin
		memberAdmin.NewMemberSignInConfigHandler,      // Added MemberSignInConfigHandler
		memberAdmin.NewMemberSignInRecordHandler,      // Added MemberSignInRecordHandler
		memberAdmin.NewMemberUserHandler,              // Added MemberUserHandler
		memberHandler.NewAppMemberUserHandler,         // Added AppMemberUserHandler
		memberHandler.NewAppMemberAddressHandler,      // Added AppMemberAddressHandler
		memberHandler.NewAppMemberPointRecordHandler,  // Added AppMemberPointRecordHandler
		memberHandler.NewAppMemberSignInRecordHandler, // Added AppMemberSignInRecordHandler
		productHandler.NewProductCategoryHandler,      // Added ProductCategoryHandler
		productHandler.NewProductPropertyHandler,      // Added ProductPropertyHandler
		productHandler.NewProductBrandHandler,         // Added ProductBrandHandler
		productHandler.NewProductSpuHandler,           // Added ProductSpuHandler
		productHandler.NewProductCommentHandler,
		productHandler.NewProductFavoriteHandler,
		productHandler.NewProductBrowseHistoryHandler,

		// App handlers
		productApp.NewAppProductFavoriteHandler,
		productApp.NewAppProductBrowseHistoryHandler,
		productApp.NewAppProductSpuHandler,
		productApp.NewAppProductCommentHandler,
		// Trade
		tradeSvc.NewCartService,
		tradeSvc.NewTradeOrderQueryService,
		tradeSvc.NewTradePriceService,
		tradeSvc.NewTradeOrderUpdateService,
		tradeSvc.NewTradeAfterSaleService,
		tradeSvc.NewTradeConfigService,   // Added Config
		tradeSvc.NewTradeOrderLogService, // Added Log
		tradeApp.NewAppCartHandler,
		tradeApp.NewAppTradeOrderHandler,
		tradeApp.NewAppTradeAfterSaleHandler,
		tradeApp.NewAppTradeConfigHandler, // Added Config
		tradeAdmin.NewTradeOrderHandler,
		tradeAdmin.NewTradeAfterSaleHandler,
		tradeAdmin.NewTradeConfigHandler, // Added Config
		// Delivery
		tradeSvc.NewDeliveryExpressService,
		tradeSvc.NewDeliveryPickUpStoreService,
		tradeSvc.NewDeliveryFreightTemplateService,
		tradeAdmin.NewDeliveryExpressHandler,
		tradeAdmin.NewDeliveryPickUpStoreHandler,
		tradeAdmin.NewDeliveryFreightTemplateHandler,
		tradeBrokerageSvc.NewBrokerageUserService,
		// Brokerage
		tradeBrokerageSvc.NewBrokerageRecordService,
		tradeBrokerageSvc.NewBrokerageWithdrawService, // Added
		paySvc.NewPayTransferService,                  // Placeholder
		paySvc.NewPayWalletService,                    // Placeholder
		brokerage.NewBrokerageUserHandler,
		brokerage.NewBrokerageRecordHandler,
		brokerage.NewBrokerageWithdrawHandler,
		appBrokerage.NewAppBrokerageUserHandler, // Added
		appBrokerage.NewAppBrokerageRecordHandler,
		appBrokerage.NewAppBrokerageWithdrawHandler,

		// Statistics
		repo.NewTradeStatisticsRepository,
		repo.NewTradeOrderStatisticsRepository,
		repo.NewTradeOrderLogRepository, // Added Log Repo
		repo.NewAfterSaleStatisticsRepository,
		repo.NewBrokerageStatisticsRepository,
		repo.NewMemberStatisticsRepository,
		repo.NewApiAccessLogStatisticsRepository,
		repo.NewPayWalletStatisticsRepository,
		productRepo.NewProductStatisticsRepository, // Product
		service.NewProductStatisticsService,        // Product
		service.NewTradeStatisticsService,          // Trade
		service.NewTradeOrderStatisticsServiceV2,
		service.NewAfterSaleStatisticsService,
		service.NewBrokerageStatisticsService,
		service.NewMemberStatisticsService,
		service.NewApiAccessLogStatisticsService,
		service.NewPayWalletStatisticsService,
		adminHandler.NewTradeStatisticsHandler,
		adminHandler.NewProductStatisticsHandler,
		adminHandler.NewMemberStatisticsHandler,
		adminHandler.NewPayStatisticsHandler,

		// Promotion
		promotionSvc.NewCouponService,
		promotionSvc.NewCouponUserService,
		promotionSvc.NewPromotionBannerService,     // Added Banner
		promotionSvc.NewRewardActivityService,      // Added Activity
		promotionSvc.NewSeckillConfigService,       // Added Seckill Config
		promotionSvc.NewSeckillActivityService,     // Added Seckill Activity
		promotionSvc.NewBargainActivityService,     // Added Bargain Activity
		promotionSvc.NewBargainRecordService,       // Added Bargain Record
		promotionSvc.NewBargainHelpService,         // Added Bargain Help
		promotionSvc.NewCombinationActivityService, // Added Combination Activity
		promotionSvc.NewCombinationRecordService,   // Added Combination Record
		promotionSvc.NewDiscountActivityService,    // Added Discount Activity
		promotionSvc.NewArticleCategoryService,     // Added Article Category
		promotionSvc.NewArticleService,             // Added Article
		promotionSvc.NewDiyTemplateService,         // Added Diy Template
		promotionSvc.NewDiyPageService,             // Added Diy Page
		promotionSvc.NewKefuService,
		promotionAdmin.NewCouponHandler,
		promotionAdmin.NewBannerHandler,              // Added Banner
		promotionAdmin.NewRewardActivityHandler,      // Added Activity
		promotionAdmin.NewSeckillConfigHandler,       // Added Seckill Config
		promotionAdmin.NewSeckillActivityHandler,     // Added Seckill Activity
		promotionAdmin.NewBargainActivityHandler,     // Added Bargain Activity
		promotionAdmin.NewCombinationActivityHandler, // Added Combination Activity
		promotionAdmin.NewDiscountActivityHandler,    // Added Discount Activity
		promotionAdmin.NewArticleCategoryHandler,     // Added Article Category
		promotionAdmin.NewArticleHandler,             // Added Article
		promotionAdmin.NewDiyTemplateHandler,         // Added Diy Template
		promotionAdmin.NewDiyPageHandler,             // Added Diy Page
		promotionAdmin.NewKefuHandler,
		promotionApp.NewAppKefuHandler,
		promotionApp.NewAppCouponHandler,
		promotionApp.NewAppBannerHandler, // Added Banner
		promotionApp.NewAppBargainActivityHandler,
		promotionApp.NewAppBargainRecordHandler,
		promotionApp.NewAppBargainHelpHandler,
		promotionApp.NewAppCombinationActivityHandler, // Added Combination Activity
		promotionApp.NewAppCombinationRecordHandler,   // Added Combination Record
		promotionApp.NewAppArticleHandler,             // Added Article
		promotionApp.NewAppDiyPageHandler,             // Added Diy Page
		// Pay
		paySvc.NewPayAppService,
		paySvc.NewPayChannelService,
		paySvc.NewPayOrderService,
		paySvc.NewPayRefundService,
		paySvc.NewPayNotifyService,
		client.NewPayClientFactory,

		deliveryClient.NewExpressClientFactory, // Added ExpressClientFactory
		wire.Bind(new(deliveryClient.ExpressClientFactory), new(*deliveryClient.ExpressClientFactoryImpl)),

		payAdmin.NewPayAppHandler,
		payAdmin.NewPayChannelHandler,
		payAdmin.NewPayOrderHandler,
		payAdmin.NewPayRefundHandler,
		payAdmin.NewPayNotifyHandler,

		// Router
		router.InitRouter,
	)
	return &gin.Engine{}, nil
}
