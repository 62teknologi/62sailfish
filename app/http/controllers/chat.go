package controllers

import (
	"fmt"
	"net/http"

	"github.com/62teknologi/62sailfish/62golib/utils"
	chat_adapter "github.com/62teknologi/62sailfish/app/adapters/chat"
	"github.com/62teknologi/62sailfish/app/interfaces"
	"github.com/62teknologi/62sailfish/config"
	"github.com/gin-gonic/gin"
)

type Chat struct {
	SingularName  string
	PluralName    string
	SingularLabel string
	PluralLabel   string
	Table         string
}

func (ctrl *Chat) Init(ctx *gin.Context) {
	ctrl.SingularName = "chat"
	ctrl.PluralName = "chats"
	ctrl.SingularLabel = ctrl.SingularName
	ctrl.PluralLabel = ctrl.PluralName
	ctrl.Table = ctrl.PluralName
}

func (ctrl *Chat) Send(ctx *gin.Context) {
	ctrl.Init(ctx)

	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load loadedConfig: %e", err)
		return
	}

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
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "validation error", validation.Errors))
		return
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	var sender string
	if transformer["sender"] == nil {
		sender = ""
	} else {
		sender = transformer["sender"].(string)
	}

	var chatAdapter interfaces.Chat
	switch input["type"] {
	case "whatsapp":
		chatAdapter = chat_adapter.NewVonageWhatsapp(loadedConfig)
	case "sms":
		chatAdapter = chat_adapter.NewVonageSms(loadedConfig)
	default:
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "invalid chat type", nil))
		return
	}

	err = chatAdapter.SendMessage(sender, transformer["recipient"].(string), transformer["message"].(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "success send chat", nil))
}
