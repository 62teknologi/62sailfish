package controllers

import (
	"net/http"

	"github.com/62teknologi/62sailfish/62golib/utils"
	"github.com/62teknologi/62sailfish/app/fcm"
	"github.com/62teknologi/62sailfish/app/interfaces"
	"github.com/gin-gonic/gin"
)

type Firebase struct {
	SingularName  string
	PluralName    string
	SingularLabel string
	PluralLabel   string
	Table         string
}

func (ctrl *Firebase) Init(ctx *gin.Context) {
	ctrl.SingularName = "fcm_token"
	ctrl.PluralName = "fcm_tokens"
	ctrl.SingularLabel = ctrl.SingularName
	ctrl.PluralLabel = ctrl.PluralName
	ctrl.Table = ctrl.PluralName
}

type Response struct {
	TokenError []string `json:"token_error"`
	Count      int      `json:"count"`
}

type ResponseTopic struct {
	Id string `json:"id"`
}

func (ctrl *Firebase) PostToken(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, err := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/create.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	var input map[string]any

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusOK, utils.ResponseData("failed", "validation", validation.Errors))
		return
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	result := utils.DB.Table(ctrl.Table+"").Where("user_id = ?", input["user_id"]).Updates(&transformer)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", result.Error.Error(), nil))
		return
	}
	if result.RowsAffected == 0 && result.Error == nil {
		if err := utils.DB.Table(ctrl.Table + "").Create(&transformer).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
			return
		}
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl *Firebase) PushNotification(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, err := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/send.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	var input map[string]any

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusOK, utils.ResponseData("failed", "validation", validation.Errors))
		return
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	if _, ok := input["image"]; !ok {
		input["image"] = ""
	}

	val, errTokens := input["tokens"]
	token, errToken := input["token"]

	var firebaseAdapter interfaces.FirebasePushNotification
	firebaseAdapter, err = fcm.NewFirebaseNotificationAdapter()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Something went wrong on connect to firebase", nil))
		return
	}

	if errTokens {
		delete(input, "tokens")
		delete(transformer, "tokens")
		tokensInterface, ok := val.([]interface{})
		if !ok {
			ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Token must be an array", nil))
			return
		}

		tokens := make([]string, len(tokensInterface))
		for i, token := range tokensInterface {
			tokenString, ok := token.(string)
			if !ok {
				ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Token must be an array", nil))
				return
			}
			tokens[i] = tokenString
		}
		if total := len(tokens); total > 0 {

			//Return token error to send notification
			te, _ := firebaseAdapter.PostNotifications(tokens, input)

			ctx.JSON(http.StatusOK, utils.ResponseData("success", "message sent", &Response{TokenError: te, Count: len(te)}))
			return
		}
	}
	if errToken {
		err := firebaseAdapter.PostNotification(token.(string), input)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "message not sent", err))
			return
		}

		ctx.JSON(http.StatusOK, utils.ResponseData("success", "message sent", nil))
		return
	}
	if !errTokens && !errToken {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "You must input tokens = [] or tokens = '' ", nil))
		return
	}
}

func (ctrl *Firebase) SubscribeTopic(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, err := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/subscribe.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	var input map[string]any

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("failed", "validation", validation.Errors))
		return
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	val, errTokens := input["tokens"]

	var firebaseAdapter interfaces.FirebasePushNotification
	firebaseAdapter, err = fcm.NewFirebaseNotificationAdapter()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Something went wrong on connect to firebase", nil))
	}

	if errTokens {
		delete(input, "tokens")
		delete(transformer, "tokens")
		tokensInterface, ok := val.([]interface{})
		if !ok {
			ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Token must be an array", nil))
			return
		}

		tokens := make([]string, len(tokensInterface))
		for i, token := range tokensInterface {
			tokenString, ok := token.(string)
			if !ok {
				ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Token must be an array", nil))
				return
			}
			tokens[i] = tokenString
		}
		if total := len(tokens); total > 0 {

			res, errTopic := firebaseAdapter.SubscribeTopic(tokens, input["topic"].(string))

			if errTopic != nil {
				ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", errTopic.Error(), nil))
				return
			}

			ctx.JSON(http.StatusOK, utils.ResponseData("success", "Topic subscribed", res))
			return
		}
	}
}

func (ctrl *Firebase) UnsubscribeTopic(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, err := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/subscribe.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	var input map[string]any

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusOK, utils.ResponseData("failed", "validation", validation.Errors))
		return
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	val, errTokens := input["tokens"]

	var firebaseAdapter interfaces.FirebasePushNotification
	firebaseAdapter, err = fcm.NewFirebaseNotificationAdapter()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Something went wrong on connect to firebase", nil))
	}

	if errTokens {
		delete(input, "tokens")
		delete(transformer, "tokens")
		tokensInterface, ok := val.([]interface{})
		if !ok {
			ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Token must be an array", nil))
			return
		}

		tokens := make([]string, len(tokensInterface))
		for i, token := range tokensInterface {
			tokenString, ok := token.(string)
			if !ok {
				ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Token must be an array", nil))
				return
			}
			tokens[i] = tokenString
		}
		if total := len(tokens); total > 0 {

			res, errTopic := firebaseAdapter.UnsubscribeTopic(tokens, input["topic"].(string))

			if errTopic != nil {
				ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", errTopic.Error(), nil))
				return
			}

			ctx.JSON(http.StatusOK, utils.ResponseData("success", "Topic unsubscribed", res))
			return
		}
	}
}

func (ctrl *Firebase) PushTopicNotification(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, err := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/post_topic.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	var input map[string]any

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusOK, utils.ResponseData("failed", "validation", validation.Errors))
		return
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	var firebaseAdapter interfaces.FirebasePushNotification
	firebaseAdapter, err = fcm.NewFirebaseNotificationAdapter()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Something went wrong on connect to firebase", nil))
	}

	topic, errInput := input["topic"]

	if _, ok := input["image"]; !ok {
		input["image"] = ""
	}

	if !errInput {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Topic are required", nil))
		return
	}

	res, errPost := firebaseAdapter.PostTopic(topic.(string), input)
	if errPost != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Error", errPost))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "Topic notification sended", &ResponseTopic{Id: res}))
	return
}
