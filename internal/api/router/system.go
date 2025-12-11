package router

import (
	"backend-go/internal/api/handler"
	"backend-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterSystemRoutes 注册系统管理模块路由
func RegisterSystemRoutes(engine *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	tenantHandler *handler.TenantHandler,
	dictHandler *handler.DictHandler,
	deptHandler *handler.DeptHandler,
	postHandler *handler.PostHandler,
	roleHandler *handler.RoleHandler,
	permissionHandler *handler.PermissionHandler,
	noticeHandler *handler.NoticeHandler,
	loginLogHandler *handler.LoginLogHandler,
	operateLogHandler *handler.OperateLogHandler,
	configHandler *handler.ConfigHandler,
	smsChannelHandler *handler.SmsChannelHandler,
	smsTemplateHandler *handler.SmsTemplateHandler,
	smsLogHandler *handler.SmsLogHandler,
	fileConfigHandler *handler.FileConfigHandler,
	fileHandler *handler.FileHandler,
	jobHandler *handler.JobHandler,
	jobLogHandler *handler.JobLogHandler,
	apiAccessLogHandler *handler.ApiAccessLogHandler,
	apiErrorLogHandler *handler.ApiErrorLogHandler,
	socialClientHandler *handler.SocialClientHandler,
	socialUserHandler *handler.SocialUserHandler,
) {
	api := engine.Group("/admin-api")
	{
		systemGroup := api.Group("/system")
		{
			// Auth
			authGroup := systemGroup.Group("/auth")
			{
				authGroup.POST("/login", authHandler.Login)
				authGroup.POST("/logout", authHandler.Logout)
				authGroup.POST("/refresh-token", authHandler.RefreshToken)
				authGroup.POST("/sms-login", authHandler.SmsLogin)
				authGroup.POST("/send-sms-code", authHandler.SendSmsCode)
				authGroup.POST("/register", authHandler.Register)
				authGroup.POST("/reset-password", authHandler.ResetPassword)
				authGroup.GET("/social-auth-redirect", authHandler.SocialAuthRedirect)
				authGroup.POST("/social-login", authHandler.SocialLogin)
				// 需要认证的接口
				authGroup.Use(middleware.Auth())
				authGroup.GET("/get-permission-info", authHandler.GetPermissionInfo)
			}

			// Tenant (No Auth needed for simple list? Usually Tenant mgmt needs auth. Checked Java: @PreAuthorize)
			// Assuming Auth middleware is applied to systemGroup, so we just add routes.
			tenantGroup := systemGroup.Group("/tenant")
			{
				tenantGroup.GET("/simple-list", tenantHandler.GetTenantSimpleList)
				tenantGroup.GET("/get-by-website", tenantHandler.GetTenantByWebsite)
				tenantGroup.GET("/get-id-by-name", tenantHandler.GetTenantIdByName)
				tenantGroup.POST("/create", tenantHandler.CreateTenant)
				tenantGroup.PUT("/update", tenantHandler.UpdateTenant)
				tenantGroup.DELETE("/delete", tenantHandler.DeleteTenant)
				tenantGroup.GET("/get", tenantHandler.GetTenant)
				tenantGroup.GET("/page", tenantHandler.GetTenantPage)
				tenantGroup.GET("/export-excel", tenantHandler.ExportTenantExcel)
			}

			// Dict Type
			dictTypeGroup := systemGroup.Group("/dict-type")
			{
				dictTypeGroup.GET("/simple-list", dictHandler.GetSimpleDictTypeList)
				dictTypeGroup.GET("/page", dictHandler.GetDictTypePage)
				dictTypeGroup.GET("/get", dictHandler.GetDictType)
				dictTypeGroup.POST("/create", dictHandler.CreateDictType)
				dictTypeGroup.PUT("/update", dictHandler.UpdateDictType)
				dictTypeGroup.DELETE("/delete", dictHandler.DeleteDictType)
				dictTypeGroup.GET("/export-excel", dictHandler.ExportDictTypeExcel)
			}

			// Dict Data
			dictDataGroup := systemGroup.Group("/dict-data")
			{
				dictDataGroup.GET("/simple-list", dictHandler.GetSimpleDictDataList)
				dictDataGroup.GET("/list-all-simple", dictHandler.GetSimpleDictDataList)
				dictDataGroup.GET("/page", dictHandler.GetDictDataPage)
				dictDataGroup.GET("/get", dictHandler.GetDictData)
				dictDataGroup.POST("/create", dictHandler.CreateDictData)
				dictDataGroup.PUT("/update", dictHandler.UpdateDictData)
				dictDataGroup.DELETE("/delete", dictHandler.DeleteDictData)
			}

			// Dept
			deptGroup := systemGroup.Group("/dept")
			{
				deptGroup.GET("/list", deptHandler.GetDeptList)
				deptGroup.GET("/list-all-simple", deptHandler.GetSimpleDeptList)
				deptGroup.GET("/simple-list", deptHandler.GetSimpleDeptList)
				deptGroup.GET("/get", deptHandler.GetDept)
				deptGroup.POST("/create", deptHandler.CreateDept)
				deptGroup.PUT("/update", deptHandler.UpdateDept)
				deptGroup.DELETE("/delete", deptHandler.DeleteDept)
			}

			// Post
			postGroup := systemGroup.Group("/post")
			{
				postGroup.GET("/page", postHandler.GetPostPage)
				postGroup.GET("/simple-list", postHandler.GetSimplePostList)
				postGroup.GET("/get", postHandler.GetPost)
				postGroup.POST("/create", postHandler.CreatePost)
				postGroup.PUT("/update", postHandler.UpdatePost)
				postGroup.DELETE("/delete", postHandler.DeletePost)
			}

			// User
			userGroup := systemGroup.Group("/user")
			{
				userGroup.GET("/page", userHandler.GetUserPage)
				userGroup.GET("/list-all-simple", userHandler.GetSimpleUserList)
				userGroup.GET("/simple-list", userHandler.GetSimpleUserList)
				userGroup.GET("/get", userHandler.GetUser)
				userGroup.POST("/create", userHandler.CreateUser)
				userGroup.PUT("/update", userHandler.UpdateUser)
				userGroup.DELETE("/delete", userHandler.DeleteUser)
				userGroup.PUT("/update-status", userHandler.UpdateUserStatus)
				userGroup.PUT("/update-password", userHandler.UpdateUserPassword)
				userGroup.GET("/export", userHandler.ExportUser)
				userGroup.GET("/get-import-template", userHandler.GetImportTemplate)
				userGroup.POST("/import", userHandler.ImportUser)
			}

			// Role
			roleGroup := systemGroup.Group("/role")
			{
				roleGroup.GET("/page", roleHandler.GetRolePage)
				roleGroup.GET("/list-all-simple", roleHandler.GetSimpleRoleList)
				roleGroup.GET("/simple-list", roleHandler.GetSimpleRoleList)
				roleGroup.GET("/get", roleHandler.GetRole)
				roleGroup.POST("/create", roleHandler.CreateRole)
				roleGroup.PUT("/update", roleHandler.UpdateRole)
				roleGroup.PUT("/update-status", roleHandler.UpdateRoleStatus)
				roleGroup.DELETE("/delete", roleHandler.DeleteRole)
			}

			// Permission
			permGroup := systemGroup.Group("/permission")
			{
				permGroup.GET("/list-role-menus", permissionHandler.GetRoleMenuList)
				permGroup.POST("/assign-role-menu", permissionHandler.AssignRoleMenu)
				permGroup.POST("/assign-role-data-scope", permissionHandler.AssignRoleDataScope)
				permGroup.GET("/list-user-roles", permissionHandler.GetUserRoleList)
				permGroup.POST("/assign-user-role", permissionHandler.AssignUserRole)
			}

			// Notice
			noticeGroup := systemGroup.Group("/notice")
			{
				noticeGroup.GET("/page", noticeHandler.GetNoticePage)
				noticeGroup.GET("/get", noticeHandler.GetNotice)
				noticeGroup.POST("/create", noticeHandler.CreateNotice)
				noticeGroup.PUT("/update", noticeHandler.UpdateNotice)
				noticeGroup.DELETE("/delete", noticeHandler.DeleteNotice)
				noticeGroup.POST("/push", noticeHandler.Push)
			}

			// Login Log
			loginLogGroup := systemGroup.Group("/login-log")
			{
				loginLogGroup.GET("/page", loginLogHandler.GetLoginLogPage)
			}

			// Operate Log
			operateLogGroup := systemGroup.Group("/operate-log")
			{
				operateLogGroup.GET("/page", operateLogHandler.GetOperateLogPage)
			}
		}

		// Sms Channel
		smsChannelGroup := api.Group("/system/sms-channel")
		{
			smsChannelGroup.POST("/create", smsChannelHandler.CreateSmsChannel)
			smsChannelGroup.PUT("/update", smsChannelHandler.UpdateSmsChannel)
			smsChannelGroup.DELETE("/delete", smsChannelHandler.DeleteSmsChannel)
			smsChannelGroup.GET("/get", smsChannelHandler.GetSmsChannel)
			smsChannelGroup.GET("/page", smsChannelHandler.GetSmsChannelPage)
			smsChannelGroup.GET("/simple-list", smsChannelHandler.GetSimpleSmsChannelList)
		}

		// Sms Template
		smsTemplateGroup := api.Group("/system/sms-template")
		{
			smsTemplateGroup.POST("/create", smsTemplateHandler.CreateSmsTemplate)
			smsTemplateGroup.PUT("/update", smsTemplateHandler.UpdateSmsTemplate)
			smsTemplateGroup.DELETE("/delete", smsTemplateHandler.DeleteSmsTemplate)
			smsTemplateGroup.GET("/get", smsTemplateHandler.GetSmsTemplate)
			smsTemplateGroup.GET("/page", smsTemplateHandler.GetSmsTemplatePage)
		}

		// Sms Log
		smsLogGroup := api.Group("/system/sms-log")
		{
			smsLogGroup.GET("/page", smsLogHandler.GetSmsLogPage)
		}

		// Config
		configGroup := api.Group("/infra/config")
		{
			configGroup.GET("/page", configHandler.GetConfigPage)
			configGroup.GET("/get", configHandler.GetConfig)
			configGroup.GET("/get-value-by-key", configHandler.GetConfigKey)
			configGroup.POST("/create", configHandler.CreateConfig)
			configGroup.PUT("/update", configHandler.UpdateConfig)
			configGroup.DELETE("/delete", configHandler.DeleteConfig)
		}

		// Infra
		infraGroup := api.Group("/infra")
		{
			// File Config
			fileConfigGroup := infraGroup.Group("/file-config")
			{
				fileConfigGroup.POST("/create", fileConfigHandler.CreateFileConfig)
				fileConfigGroup.PUT("/update", fileConfigHandler.UpdateFileConfig)
				fileConfigGroup.PUT("/update-master", fileConfigHandler.UpdateFileConfigMaster)
				fileConfigGroup.DELETE("/delete", fileConfigHandler.DeleteFileConfig)
				fileConfigGroup.GET("/page", fileConfigHandler.GetFileConfigPage)
				fileConfigGroup.GET("/get", fileConfigHandler.GetFileConfig)
			}

			// File
			fileGroup := infraGroup.Group("/file")
			{
				fileGroup.POST("/upload", fileHandler.UploadFile)
				fileGroup.DELETE("/delete", fileHandler.DeleteFile)
				fileGroup.GET("/page", fileHandler.GetFilePage)
			}

			// Job
			jobGroup := infraGroup.Group("/job")
			{
				jobGroup.POST("/create", jobHandler.CreateJob)
				jobGroup.PUT("/update", jobHandler.UpdateJob)
				jobGroup.PUT("/update-status", jobHandler.UpdateJobStatus)
				jobGroup.DELETE("/delete", jobHandler.DeleteJob)
				jobGroup.GET("/get", jobHandler.GetJob)
				jobGroup.GET("/page", jobHandler.GetJobPage)
				jobGroup.PUT("/trigger", jobHandler.TriggerJob)
			}

			// Job Log
			jobLogGroup := infraGroup.Group("/job-log")
			{
				jobLogGroup.GET("/get", jobLogHandler.GetJobLog)
				jobLogGroup.GET("/page", jobLogHandler.GetJobLogPage)
			}

			// API Access Log
			apiAccessLogGroup := infraGroup.Group("/api-access-log")
			{
				apiAccessLogGroup.GET("/page", apiAccessLogHandler.GetApiAccessLogPage)
			}

			// API Error Log
			apiErrorLogGroup := infraGroup.Group("/api-error-log")
			{
				apiErrorLogGroup.GET("/page", apiErrorLogHandler.GetApiErrorLogPage)
				apiErrorLogGroup.PUT("/update-status", apiErrorLogHandler.UpdateApiErrorLogProcess)
			}

			// Social Client
			socialClientGroup := infraGroup.Group("/social-client")
			{
				socialClientGroup.POST("/create", socialClientHandler.CreateSocialClient)
				socialClientGroup.PUT("/update", socialClientHandler.UpdateSocialClient)
				socialClientGroup.DELETE("/delete", socialClientHandler.DeleteSocialClient)
				socialClientGroup.GET("/get", socialClientHandler.GetSocialClient)
				socialClientGroup.GET("/page", socialClientHandler.GetSocialClientPage)
			}

			// Social User
			socialUserGroup := infraGroup.Group("/social-user")
			{
				socialUserGroup.POST("/bind", socialUserHandler.BindSocialUser)
				socialUserGroup.POST("/unbind", socialUserHandler.UnbindSocialUser)
				socialUserGroup.GET("/list", socialUserHandler.GetSocialUserList)
				socialUserGroup.GET("/get", socialUserHandler.GetSocialUser)
				socialUserGroup.GET("/page", socialUserHandler.GetSocialUserPage)
			}
		}
	}
}
