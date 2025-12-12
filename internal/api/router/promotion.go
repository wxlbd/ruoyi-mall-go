package router

import (
	promotionAdmin "backend-go/internal/api/handler/admin/promotion"
	"backend-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterPromotionRoutes 注册营销活动模块路由
func RegisterPromotionRoutes(engine *gin.Engine,
	couponHandler *promotionAdmin.CouponHandler,
	bannerHandler *promotionAdmin.BannerHandler,
	rewardActivityHandler *promotionAdmin.RewardActivityHandler,
	seckillConfigHandler *promotionAdmin.SeckillConfigHandler,
	seckillActivityHandler *promotionAdmin.SeckillActivityHandler,
	bargainActivityHandler *promotionAdmin.BargainActivityHandler,
	combinationActivityHandler *promotionAdmin.CombinationActivityHandler,
	discountActivityHandler *promotionAdmin.DiscountActivityHandler,
	articleCategoryHandler *promotionAdmin.ArticleCategoryHandler,
	articleHandler *promotionAdmin.ArticleHandler,
	diyTemplateHandler *promotionAdmin.DiyTemplateHandler,
	diyPageHandler *promotionAdmin.DiyPageHandler,
	kefuHandler *promotionAdmin.KefuHandler,
	pointActivityHandler *promotionAdmin.PointActivityHandler,
	bargainRecordHandler *promotionAdmin.BargainRecordHandler,
	combinationRecordHandler *promotionAdmin.CombinationRecordHandler,
	bargainHelpHandler *promotionAdmin.BargainHelpHandler,
) {
	promotionGroup := engine.Group("/admin-api/promotion")
	promotionGroup.Use(middleware.Auth())
	{
		// Coupon Template
		templateGroup := promotionGroup.Group("/coupon-template")
		{
			templateGroup.POST("/create", couponHandler.CreateCouponTemplate)
			templateGroup.PUT("/update", couponHandler.UpdateCouponTemplate)
			templateGroup.PUT("/update-status", couponHandler.UpdateCouponTemplateStatus)
			templateGroup.DELETE("/delete", couponHandler.DeleteCouponTemplate)
			templateGroup.GET("/get", couponHandler.GetCouponTemplate)
			templateGroup.GET("/page", couponHandler.GetCouponTemplatePage)
			templateGroup.GET("/list", couponHandler.GetCouponTemplateList)
		}

		// Coupon
		couponGroup := promotionGroup.Group("/coupon")
		{
			couponGroup.DELETE("/delete", couponHandler.DeleteCoupon)
			couponGroup.GET("/page", couponHandler.GetCouponPage)
			couponGroup.POST("/send", couponHandler.SendCoupon)
		}

		// Banner
		bannerGroup := promotionGroup.Group("/banner")
		{
			bannerGroup.POST("/create", bannerHandler.CreateBanner)
			bannerGroup.PUT("/update", bannerHandler.UpdateBanner)
			bannerGroup.DELETE("/delete", bannerHandler.DeleteBanner)
			bannerGroup.GET("/get", bannerHandler.GetBanner)
			bannerGroup.GET("/page", bannerHandler.GetBannerPage)
		}

		// Reward Activity
		rewardGroup := promotionGroup.Group("/reward-activity")
		{
			rewardGroup.POST("/create", rewardActivityHandler.CreateRewardActivity)
			rewardGroup.PUT("/update", rewardActivityHandler.UpdateRewardActivity)
			rewardGroup.DELETE("/delete", rewardActivityHandler.DeleteRewardActivity)
			rewardGroup.GET("/get", rewardActivityHandler.GetRewardActivity)
			rewardGroup.GET("/page", rewardActivityHandler.GetRewardActivityPage)
		}

		// Seckill Config
		seckillConfigGroup := promotionGroup.Group("/seckill-config")
		{
			seckillConfigGroup.POST("/create", seckillConfigHandler.CreateSeckillConfig)
			seckillConfigGroup.PUT("/update", seckillConfigHandler.UpdateSeckillConfig)
			seckillConfigGroup.PUT("/update-status", seckillConfigHandler.UpdateSeckillConfigStatus)
			seckillConfigGroup.DELETE("/delete", seckillConfigHandler.DeleteSeckillConfig)
			seckillConfigGroup.GET("/get", seckillConfigHandler.GetSeckillConfig)
			seckillConfigGroup.GET("/page", seckillConfigHandler.GetSeckillConfigPage)
			seckillConfigGroup.GET("/list", seckillConfigHandler.GetSeckillConfigList)
			seckillConfigGroup.GET("/simple-list", seckillConfigHandler.GetSeckillConfigSimpleList)
		}

		// Seckill Activity
		seckillActivityGroup := promotionGroup.Group("/seckill-activity")
		{
			seckillActivityGroup.POST("/create", seckillActivityHandler.CreateSeckillActivity)
			seckillActivityGroup.PUT("/update", seckillActivityHandler.UpdateSeckillActivity)
			seckillActivityGroup.DELETE("/delete", seckillActivityHandler.DeleteSeckillActivity)
			seckillActivityGroup.PUT("/close", seckillActivityHandler.CloseSeckillActivity)
			seckillActivityGroup.GET("/get", seckillActivityHandler.GetSeckillActivity)
			seckillActivityGroup.GET("/page", seckillActivityHandler.GetSeckillActivityPage)
		}

		// Bargain Activity
		bargainActivityGroup := promotionGroup.Group("/bargain-activity")
		{
			bargainActivityGroup.POST("/create", bargainActivityHandler.CreateBargainActivity)
			bargainActivityGroup.PUT("/update", bargainActivityHandler.UpdateBargainActivity)
			bargainActivityGroup.DELETE("/delete", bargainActivityHandler.DeleteBargainActivity)
			bargainActivityGroup.PUT("/close", bargainActivityHandler.CloseBargainActivity)
			bargainActivityGroup.GET("/get", bargainActivityHandler.GetBargainActivity)
			bargainActivityGroup.GET("/page", bargainActivityHandler.GetBargainActivityPage)
		}

		// Combination Activity
		combinationActivityGroup := promotionGroup.Group("/combination-activity")
		{
			combinationActivityGroup.POST("/create", combinationActivityHandler.CreateCombinationActivity)
			combinationActivityGroup.PUT("/update", combinationActivityHandler.UpdateCombinationActivity)
			combinationActivityGroup.DELETE("/delete", combinationActivityHandler.DeleteCombinationActivity)
			combinationActivityGroup.GET("/get", combinationActivityHandler.GetCombinationActivity)
			combinationActivityGroup.GET("/page", combinationActivityHandler.GetCombinationActivityPage)
		}

		// Discount Activity
		discountActivityGroup := promotionGroup.Group("/discount-activity")
		{
			discountActivityGroup.POST("/create", discountActivityHandler.CreateDiscountActivity)
			discountActivityGroup.PUT("/update", discountActivityHandler.UpdateDiscountActivity)
			discountActivityGroup.POST("/close", discountActivityHandler.CloseDiscountActivity)
			discountActivityGroup.DELETE("/delete", discountActivityHandler.DeleteDiscountActivity)
			discountActivityGroup.GET("/get", discountActivityHandler.GetDiscountActivity)
			discountActivityGroup.GET("/page", discountActivityHandler.GetDiscountActivityPage)
		}

		// Article Category
		articleCategoryGroup := promotionGroup.Group("/article-category")
		{
			articleCategoryGroup.POST("/create", articleCategoryHandler.CreateArticleCategory)
			articleCategoryGroup.PUT("/update", articleCategoryHandler.UpdateArticleCategory)
			articleCategoryGroup.DELETE("/delete", articleCategoryHandler.DeleteArticleCategory)
			articleCategoryGroup.GET("/get", articleCategoryHandler.GetArticleCategory)
			articleCategoryGroup.GET("/list", articleCategoryHandler.GetArticleCategoryList)
			articleCategoryGroup.GET("/simple-list", articleCategoryHandler.GetSimpleList)
		}

		// Article
		articleGroup := promotionGroup.Group("/article")
		{
			articleGroup.POST("/create", articleHandler.CreateArticle)
			articleGroup.PUT("/update", articleHandler.UpdateArticle)
			articleGroup.DELETE("/delete", articleHandler.DeleteArticle)
			articleGroup.GET("/get", articleHandler.GetArticle)
			articleGroup.GET("/page", articleHandler.GetArticlePage)
		}

		// DIY Template
		diyTemplateGroup := promotionGroup.Group("/diy-template")
		{
			diyTemplateGroup.POST("/create", diyTemplateHandler.CreateDiyTemplate)
			diyTemplateGroup.PUT("/update", diyTemplateHandler.UpdateDiyTemplate)
			diyTemplateGroup.DELETE("/delete", diyTemplateHandler.DeleteDiyTemplate)
			diyTemplateGroup.GET("/get", diyTemplateHandler.GetDiyTemplate)
			diyTemplateGroup.GET("/page", diyTemplateHandler.GetDiyTemplatePage)
			diyTemplateGroup.GET("/get-property", diyTemplateHandler.GetDiyTemplateProperty)
		}

		// DIY Page
		diyPageGroup := promotionGroup.Group("/diy-page")
		{
			diyPageGroup.POST("/create", diyPageHandler.CreateDiyPage)
			diyPageGroup.PUT("/update", diyPageHandler.UpdateDiyPage)
			diyPageGroup.DELETE("/delete", diyPageHandler.DeleteDiyPage)
			diyPageGroup.GET("/get", diyPageHandler.GetDiyPage)
			diyPageGroup.GET("/page", diyPageHandler.GetDiyPagePage)
			diyPageGroup.GET("/get-property", diyPageHandler.GetDiyPageProperty)
		}

		// Kefu Conversation (Admin)
		kefuConversationGroup := promotionGroup.Group("/kefu-conversation")
		{
			kefuConversationGroup.GET("/page", kefuHandler.GetConversationPage)
			kefuConversationGroup.DELETE("/delete", kefuHandler.DeleteConversation)
		}

		// Kefu Message (Admin)
		kefuMessageGroup := promotionGroup.Group("/kefu-message")
		{
			kefuMessageGroup.POST("/send", kefuHandler.SendMessage)
			kefuMessageGroup.GET("/page", kefuHandler.GetMessagePage)
		}

		// Point Activity
		pointActivityGroup := promotionGroup.Group("/point-activity")
		{
			pointActivityGroup.POST("/create", pointActivityHandler.CreatePointActivity)
			pointActivityGroup.PUT("/update", pointActivityHandler.UpdatePointActivity)
			pointActivityGroup.PUT("/close", pointActivityHandler.ClosePointActivity)
			pointActivityGroup.DELETE("/delete", pointActivityHandler.DeletePointActivity)
			pointActivityGroup.GET("/get", pointActivityHandler.GetPointActivity)
			pointActivityGroup.GET("/page", pointActivityHandler.GetPointActivityPage)
			pointActivityGroup.GET("/list-by-ids", pointActivityHandler.GetPointActivityListByIds)
		}

		// Bargain Record
		bargainRecordGroup := promotionGroup.Group("/bargain-record")
		{
			bargainRecordGroup.GET("/page", bargainRecordHandler.GetBargainRecordPage)
		}

		// Combination Record
		combinationRecordGroup := promotionGroup.Group("/combination-record")
		{
			combinationRecordGroup.GET("/page", combinationRecordHandler.GetCombinationRecordPage)
		}

		// Bargain Help
		bargainHelpGroup := promotionGroup.Group("/bargain-help")
		{
			bargainHelpGroup.GET("/page", bargainHelpHandler.GetBargainHelpPage)
		}
	}
}
