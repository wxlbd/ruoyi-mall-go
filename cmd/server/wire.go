//go:build wireinject
// +build wireinject

package main

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler" // Statistics Handlers
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/router"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/permission"
	ws "github.com/wxlbd/ruoyi-mall-go/internal/pkg/websocket"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo" // Pay Repo
	iotRepo "github.com/wxlbd/ruoyi-mall-go/internal/repo/iot"
	payRepo "github.com/wxlbd/ruoyi-mall-go/internal/repo/pay"
	tradeRepo "github.com/wxlbd/ruoyi-mall-go/internal/repo/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/infra"
	infraJob "github.com/wxlbd/ruoyi-mall-go/internal/service/infra/job"
	iotSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	promotionJob "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion/job"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay/job"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"

	productRepo "github.com/wxlbd/ruoyi-mall-go/internal/repo/product" // Product Statistics Repo
	product "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/product"
	promotionSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	tradeSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade"
	tradeBrokerageSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade/calculators"
	deliveryClient "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade/delivery/client"
	tradeJob "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade/job"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	paySvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay/client"
	_ "github.com/wxlbd/ruoyi-mall-go/internal/service/pay/client/alipay"
	_ "github.com/wxlbd/ruoyi-mall-go/internal/service/pay/client/weixin"
	payWalletSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay/wallet"

	"github.com/wxlbd/ruoyi-mall-go/pkg/cache"
	"github.com/wxlbd/ruoyi-mall-go/pkg/database"
	"github.com/wxlbd/ruoyi-mall-go/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitApp() (*gin.Engine, error) {
	wire.Build(
		database.InitDB,
		cache.InitRedis,
		logger.NewLogger,
		// Repo (GORM Gen)
		repo.NewQuery,
		iotRepo.ProviderSet,
		// Service
		system.NewOAuth2TokenService,
		system.NewAuthService,
		system.NewMenuService,
		system.NewRoleService,
		system.NewPermissionService,
		system.NewTenantService,
		system.NewTenantPackageService,
		system.NewUserService,
		system.NewDictService,
		system.NewDeptService,
		system.NewPostService,
		system.NewNoticeService,
		system.NewConfigService,
		system.NewSmsClientFactory,             // Added SmsClientFactory
		system.NewSmsChannelService,            // Added SmsChannelService
		system.NewSmsTemplateService,           // Added SmsTemplateService
		system.NewSmsLogService,                // Added SmsLogService
		system.NewSmsSendService,               // Added SmsSendService
		infra.NewFileConfigService,             // Added FileConfigService
		infra.NewFileService,                   // Added FileService
		system.NewSmsCodeService,               // Added SmsCodeService
		system.NewLoginLogService,              // Added LoginLogService
		system.NewOperateLogService,            // Added OperateLogService
		infra.NewScheduler,                     // Added Scheduler
		infra.NewJobService,                    // Added JobService
		infra.NewJobLogService,                 // Added JobLogService
		infra.NewApiAccessLogService,           // Added ApiAccessLogService
		infra.NewApiErrorLogService,            // Added ApiErrorLogService
		system.NewSocialClientService,          // Added SocialClientService
		system.NewSocialUserService,            // Added SocialUserService
		system.NewMailService,                  // Added MailService
		system.NewNotifyService,                // Added NotifyService
		system.NewOAuth2ClientService,          // Added OAuth2ClientService
		system.NewCaptchaService,               // Added CaptchaService
		iotSvc.ProviderSet,                     // Added IotService ProviderSet
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

		// Handlers
		handler.ProviderSet,
		iot.ProviderSet,

		// Trade
		tradeSvc.NewCartService,
		tradeSvc.NewTradeOrderQueryService,
		tradeSvc.NewTradePriceService,
		tradeSvc.NewPriceCalculatorHelper,
		calculators.ProviderSet,
		ProvidePriceCalculators,
		tradeSvc.NewTradeOrderUpdateService,
		tradeSvc.NewTradeAfterSaleService,
		tradeSvc.NewAfterSaleLogService,  // Added
		tradeSvc.NewTradeConfigService,   // Added Config
		tradeSvc.NewTradeOrderLogService, // Added Log

		// Delivery
		tradeSvc.NewDeliveryExpressService,
		tradeSvc.NewDeliveryPickUpStoreService,
		tradeSvc.NewDeliveryExpressTemplateService,

		tradeBrokerageSvc.NewBrokerageUserService,
		// Brokerage
		tradeBrokerageSvc.NewBrokerageRecordService,
		tradeBrokerageSvc.NewBrokerageWithdrawService, // Added

		// Wallet Services
		payWalletSvc.NewPayWalletService,
		payWalletSvc.NewPayWalletRechargeService,
		payWalletSvc.NewPayWalletRechargePackageService,
		payWalletSvc.NewPayWalletTransactionService,

		// Pay Repositories
		payRepo.NewPayTransferRepository,
		payRepo.NewPayNoRedisDAO,
		// Trade Repositories
		tradeRepo.NewTradeNoRedisDAO,
		// System Repositories
		repo.NewNotifyTemplateRepository,
		repo.NewNotifyMessageRepository,
		// Statistics
		repo.NewTradeStatisticsRepository,
		repo.NewTradeOrderStatisticsRepository,
		repo.NewTradeOrderLogRepository, // Added Log Repo
		repo.NewAfterSaleLogRepository,  // Added AfterSale Log Repo
		repo.NewAfterSaleStatisticsRepository,
		repo.NewBrokerageStatisticsRepository,
		repo.NewMemberStatisticsRepository,
		repo.NewApiAccessLogStatisticsRepository,
		repo.NewPayWalletStatisticsRepository,
		productRepo.NewProductStatisticsRepository, // Product
		product.NewProductStatisticsService,        // Product
		tradeSvc.NewTradeStatisticsService,         // Trade
		tradeSvc.NewTradeOrderStatisticsServiceV2,
		tradeSvc.NewAfterSaleStatisticsService,
		tradeSvc.NewBrokerageStatisticsService,
		memberSvc.NewMemberStatisticsService,
		infra.NewApiAccessLogStatisticsService,
		paySvc.NewPayWalletStatisticsService,
		job.NewPayTransferSyncJob, // Added PayTransferSyncJob
		job.NewPayNotifyJob,       // Added PayNotifyJob
		job.NewPayOrderSyncJob,    // Added PayOrderSyncJob
		job.NewPayOrderExpireJob,  // Added PayOrderExpireJob
		job.NewPayRefundSyncJob,   // Added PayRefundSyncJob
		promotionJob.NewCouponExpireJob,
		infraJob.NewJobLogCleanJob,
		infraJob.NewErrorLogCleanJob,
		infraJob.NewAccessLogCleanJob,
		tradeJob.NewTradeOrderAutoCancelJob,
		tradeJob.NewTradeOrderAutoReceiveJob,
		tradeJob.NewTradeOrderAutoCommentJob,
		tradeJob.NewBrokerageRecordUnfreezeJob,

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
		promotionSvc.NewPointActivityService, promotionSvc.NewKefuService,

		// WebSocket
		ws.ProviderSet,

		// Casbin
		permission.InitEnforcer,
		middleware.NewCasbinMiddleware,

		// Pay
		paySvc.NewPayAppService,
		paySvc.NewPayChannelService,
		paySvc.NewPayOrderService,
		paySvc.NewPayRefundService,
		paySvc.NewPayNotifyService,
		paySvc.NewPayTransferService,
		client.NewPayClientFactory,

		deliveryClient.NewExpressClientFactory, // Added ExpressClientFactory
		wire.Bind(new(deliveryClient.ExpressClientFactory), new(*deliveryClient.ExpressClientFactoryImpl)),

		// Statistics Repository 接口绑定
		wire.Bind(new(tradeSvc.TradeStatisticsRepository), new(*repo.TradeStatisticsRepositoryImpl)),
		wire.Bind(new(tradeSvc.TradeOrderStatisticsRepository), new(*repo.TradeOrderStatisticsRepositoryImpl)),
		wire.Bind(new(tradeSvc.AfterSaleStatisticsRepository), new(*repo.AfterSaleStatisticsRepositoryImpl)),
		wire.Bind(new(tradeSvc.BrokerageStatisticsRepository), new(*repo.BrokerageStatisticsRepositoryImpl)),
		wire.Bind(new(memberSvc.MemberStatisticsRepository), new(*repo.MemberStatisticsRepositoryImpl)),
		wire.Bind(new(infra.ApiAccessLogStatisticsRepository), new(*repo.ApiAccessLogStatisticsRepositoryImpl)),
		wire.Bind(new(paySvc.PayWalletStatisticsRepository), new(*repo.PayWalletStatisticsRepositoryImpl)),
		wire.Bind(new(product.ProductStatisticsRepository), new(*productRepo.ProductStatisticsRepositoryImpl)),
		wire.Bind(new(system.NotifyTemplateRepository), new(*repo.NotifyTemplateRepositoryImpl)),
		wire.Bind(new(system.NotifyMessageRepository), new(*repo.NotifyMessageRepositoryImpl)),

		wire.Bind(new(promotionSvc.CombinationTradeOrderService), new(*tradeSvc.TradeOrderUpdateService)),
		wire.Bind(new(promotionSvc.CombinationSocialClientService), new(*system.SocialClientService)),
		tradeSvc.NewDefaultPromotionPriceCalculator,
		wire.Bind(new(tradeSvc.PromotionPriceCalculator), new(*tradeSvc.DefaultPromotionPriceCalculator)),

		// 接口绑定: TradeOrderUpdateService 依赖
		wire.Bind(new(tradeSvc.PayOrderServiceAPI), new(*paySvc.PayOrderService)),
		wire.Bind(new(tradeSvc.PayRefundServiceAPI), new(*paySvc.PayRefundService)),
		wire.Bind(new(tradeSvc.PayAppServiceAPI), new(*paySvc.PayAppService)),
		wire.Bind(new(tradeSvc.TradeConfigServiceAPI), new(*tradeSvc.TradeConfigService)),
		wire.Bind(new(tradeSvc.ProductSkuServiceAPI), new(*product.ProductSkuService)),
		wire.Bind(new(tradeSvc.ProductCommentServiceAPI), new(*product.ProductCommentService)),
		wire.Bind(new(tradeSvc.CouponUserServiceAPI), new(*promotionSvc.CouponUserService)),
		wire.Bind(new(tradeSvc.MemberUserServiceAPI), new(*memberSvc.MemberUserService)),
		wire.Bind(new(tradeSvc.TradeNoRedisDAOAPI), new(*tradeRepo.TradeNoRedisDAO)),

		// 接口绑定: SkuPromotionCalculator 供 CalculateProductPrice 复用
		wire.Bind(new(tradeSvc.SkuPromotionCalculator), new(*calculators.DiscountActivityPriceCalculator)),

		// Router
		router.InitRouter,

		// Job Handlers (Slice Injection for Scheduler)
		ProvideJobHandlers,
	)
	return &gin.Engine{}, nil
}

func ProvidePriceCalculators(
	bargain *calculators.BargainActivityPriceCalculator,
	combination *calculators.CombinationActivityPriceCalculator,
	coupon *calculators.CouponPriceCalculator,
	delivery *calculators.DeliveryPriceCalculator,
	discount *calculators.DiscountActivityPriceCalculator,
	pointActivity *calculators.PointActivityPriceCalculator,
	pointGive *calculators.PointGivePriceCalculator,
	pointUse *calculators.PointUsePriceCalculator,
	reward *calculators.RewardActivityPriceCalculator,
	seckill *calculators.SeckillActivityPriceCalculator,
) []tradeSvc.PriceCalculator {
	// 对齐 Java TradePriceCalculator.ORDER_* 常量定义的顺序
	// ORDER_SECKILL_ACTIVITY = 8
	// ORDER_BARGAIN_ACTIVITY = 8
	// ORDER_COMBINATION_ACTIVITY = 8
	// ORDER_POINT_ACTIVITY = 8
	// ORDER_DISCOUNT_ACTIVITY = 10  ← 折扣活动必须在优惠券之前！
	// ORDER_REWARD_ACTIVITY = 20
	// ORDER_COUPON = 30
	// ORDER_POINT_USE = 40
	// ORDER_DELIVERY = 50
	// ORDER_POINT_GIVE = 999
	return []tradeSvc.PriceCalculator{
		seckill,       // 8
		bargain,       // 8
		combination,   // 8
		pointActivity, // 8
		discount,      // 10 ← 关键：discount必须在coupon之前
		reward,        // 20
		coupon,        // 30
		pointUse,      // 40
		delivery,      // 50
		pointGive,     // 999
	}
}

// ProvideJobHandlers 聚合所有定时任务处理器，供 Wire 使用
func ProvideJobHandlers(
	h1 *job.PayTransferSyncJob,
	h2 *job.PayNotifyJob,
	h3 *job.PayOrderSyncJob,
	h4 *job.PayOrderExpireJob,
	h5 *job.PayRefundSyncJob,
	h6 *promotionJob.CouponExpireJob,
	h7 *infraJob.JobLogCleanJob,
	h8 *infraJob.ErrorLogCleanJob,
	h9 *infraJob.AccessLogCleanJob,
	h10 *tradeJob.TradeOrderAutoCancelJob,
	h11 *tradeJob.TradeOrderAutoReceiveJob,
	h12 *tradeJob.TradeOrderAutoCommentJob,
	h13 *tradeJob.BrokerageRecordUnfreezeJob,
) []infra.JobHandler {
	return []infra.JobHandler{h1, h2, h3, h4, h5, h6, h7, h8, h9, h10, h11, h12, h13}
}
