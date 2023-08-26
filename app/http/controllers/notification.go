package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/62teknologi/62sailfish/62golib/utils"
	pubsub_adapter "github.com/62teknologi/62sailfish/app/adapters/pubsub"
	"github.com/62teknologi/62sailfish/app/interfaces"
	sailfishUtils "github.com/62teknologi/62sailfish/app/utils"
	"github.com/62teknologi/62sailfish/config"
	"github.com/gin-gonic/gin"
)

type Notification struct {
	SingularName  string
	PluralName    string
	SingularLabel string
	PluralLabel   string
	Table         string
}

func (ctrl *Notification) Init(ctx *gin.Context) {
	ctrl.SingularName = "notification"
	ctrl.PluralName = "notifications"
	ctrl.SingularLabel = ctrl.SingularName
	ctrl.PluralLabel = ctrl.PluralName
	ctrl.Table = ctrl.PluralName
}

type NotificationStatus string

const (
	SENT NotificationStatus = "sent"
	READ NotificationStatus = "read"
)

type HealthStatus struct {
	ServerStatus   string            `json:"server_status"`
	DatabaseStatus string            `json:"database_status"`
	DatabaseName   string            `json:"database_name"`
	DatabaseHost   string            `json:"database_host"`
	Dependencies   map[string]string `json:"dependencies"`
}

func (ctrl *Notification) Health(ctx *gin.Context) {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load loadedConfig: %w", err)
		return
	}

	// TODO add more dependencies info
	dependencies := map[string]string{
		"rabbitmq": "ok",
	}

	if sailfishUtils.PingRabbitMQ(loadedConfig.RabbitmqUrl, loadedConfig.RabbitmqTopic, "ping") != nil {
		dependencies["rabbitmq"] = "error"
	}

	dbConn, _ := utils.DB.DB()
	parsedDsn, _ := url.Parse(loadedConfig.DBSource)
	dbHost := parsedDsn.Host
	dbName := parsedDsn.Path

	if dbHost == "" {
		// Parse DSN server format
		pairs := strings.Split(dbName, " ")
		data := make(map[string]string)
		for _, pair := range pairs {
			parts := strings.Split(pair, "=")
			if len(parts) == 2 {
				data[parts[0]] = parts[1]
			}
		}
		dbHost = data["dbHost"] + ":" + data["port"]
		dbName = data["dbname"]
	}

	if err := dbConn.Ping(); err != nil {
		ctx.JSON(http.StatusOK, utils.ResponseData("success", "Server running well", &HealthStatus{
			ServerStatus:   "ok",
			DatabaseStatus: "error",
			DatabaseName:   dbName,
			DatabaseHost:   dbHost,
			Dependencies:   dependencies,
		}))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "Server running well", &HealthStatus{
		ServerStatus:   "ok",
		DatabaseStatus: "ok",
		DatabaseName:   dbName,
		DatabaseHost:   dbHost,
		Dependencies:   dependencies,
	}))
}
func (ctrl *Notification) Find(ctx *gin.Context) {
	ctrl.Init(ctx)

	value := map[string]any{}
	columns := []string{ctrl.Table + ".*"}

	transformer, err := utils.JsonFileParser("setting/transformers/response/" + ctrl.Table + "/find.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	query := utils.DB.Table(ctrl.Table)

	utils.SetBelongsTo(query, transformer, &columns)
	delete(transformer, "filterable")

	if err := query.Select(columns).Where(ctrl.Table+".id = ?", ctx.Param("id")).Take(&value).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.SingularLabel+" not found", nil))
		return
	}

	utils.MapValuesShifter(transformer, value)
	utils.AttachBelongsTo(transformer, value)

	delete(transformer, "searchable")
	delete(transformer, "sortable")

	if err := utils.DB.Table(ctrl.Table+"").Where("id = ?", ctx.Param("id")).Update("status", READ).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+ctrl.SingularLabel+" success", transformer))
}

func notificationOrder(ctx *gin.Context) string {
	switch ctx.Query("order_by") {
	case "newest":
		return "created_at desc"
	case "oldest":
		return "created_at asc"
	default:
		return ""
	}
}

func (ctrl *Notification) FindAll(ctx *gin.Context) {
	ctrl.Init(ctx)

	values := []map[string]any{}
	columns := []string{ctrl.Table + ".*"}

	transformer, err := utils.JsonFileParser("setting/transformers/response/" + ctrl.Table + "/find.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	query := utils.DB.Table(ctrl.Table + "")
	filter := utils.SetFilterByQuery(query, transformer, ctx)
	filter["search"] = utils.SetGlobalSearch(query, transformer, ctx)
	utils.SetBelongsTo(query, transformer, &columns)

	pagination := utils.SetPagination(query, ctx)

	orderBy := notificationOrder(ctx)
	var sorter string

	if transformer["sortable"] != nil && orderBy != "" {
		query.Order(orderBy)
		sorter = orderBy
	}

	delete(transformer, "sortable")

	if err := query.Select(columns).Find(&values).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.PluralLabel+" not found", nil))
		return
	}

	delete(transformer, "searchable")
	delete(transformer, "sortable")

	customResponses := utils.MultiMapValuesShifter(transformer, values)
	response := utils.ResponseDataPaginate("success", "find "+ctrl.PluralLabel+" success", customResponses, pagination, filter, nil, nil)
	response["sorter"] = sorter
	ctx.JSON(http.StatusOK, response)
}

func (ctrl *Notification) Create(ctx *gin.Context) {
	ctrl.Init(ctx)

	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load loadedConfig: %w", err)
		return
	}

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
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "validation error", validation.Errors))
		return
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	if transformer["template"] != "" {
		templateFile, err := os.Open("public/" + transformer["template"].(string))
		if err != nil {
			fmt.Errorf("error while load template: %w", err)
			return
		}
		defer templateFile.Close()

		templateBytes, err := io.ReadAll(templateFile)
		if err != nil {
			fmt.Errorf("error while convert template to string: %w", err)
			return
		}
		templateString := string(templateBytes)

		t, err := template.New("webpage").Parse(templateString)
		if err != nil {
			fmt.Errorf("error while parse template string: %w", err)
			return
		}

		if transformer["template_params"] != nil {
			transformer["template_params"].(map[string]any)["HTTP_SERVER_ADDRESS"] = loadedConfig.BaseUrl
		}

		var buf bytes.Buffer
		err = t.Execute(&buf, transformer["template_params"])
		if err != nil {
			fmt.Errorf("error execute template: %w", err)
			return
		}

		html := buf.String()

		html = strings.ReplaceAll(html, "\n", "")
		html = strings.ReplaceAll(html, "\"", "'")

		transformer["template"] = html
	}

	delete(transformer, "template_params")
	transformer["status"] = SENT
	transformer["created_at"] = time.Now()

	if err := utils.DB.Table(ctrl.Table + "").Create(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl *Notification) Push(ctx *gin.Context) {
	ctrl.Init(ctx)

	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	transformer, err := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/push.json")
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

	if input["type"] == "email" {
		sailfishUtils.EmailSender(input["template"].(string), input["template_params"], []sailfishUtils.EmailReceiver{
			{
				Subject: input["title"].(string),
				Address: input["recipient"].(string),
				Name:    input["recipient_name"].(string),
			},
		})
	}

	if input["type"] == "rabbitmq" || input["type"] == "google_pubsub" {
		var pubsubAdapter interfaces.PubSub
		var pubsubTopic string
		switch input["type"] {
		case "rabbitmq":
			pubsubAdapter, err = pubsub_adapter.NewRabbitMQ(loadedConfig.RabbitmqUrl)
			pubsubTopic = loadedConfig.RabbitmqTopic
		case "google_pubsub":
			pubsubAdapter, err = pubsub_adapter.NewGoogleCloudPubSub(loadedConfig.GooglePubSubProjectId)
			pubsubTopic = loadedConfig.GooglePubSubTopic
		}

		jsonTransformer, _ := json.Marshal(transformer)
		err = pubsubAdapter.Publish(ctx, pubsubTopic, jsonTransformer)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
			return
		}
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "success push notification", nil))
}

func (ctrl *Notification) Consume(ctx *gin.Context) {
	ctrl.Init(ctx)

	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load loadedConfig: %w", err)
		return
	}

	// Set the response headers to enable streaming
	ctx.Header("Content-Type", "text/plain")
	ctx.Header("Transfer-Encoding", "chunked")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	// Subscribe to the RabbitMQ topic
	rabbitmq, err := pubsub_adapter.NewRabbitMQ(loadedConfig.RabbitmqUrl)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}
	messages, err := rabbitmq.Subscribe(ctx, loadedConfig.RabbitmqTopic)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("Failed to subscribe: %v", err.Error(), nil))
		return
	}

	// Start streaming the messages to the client
	flusher := ctx.Writer.(http.Flusher)
	for msg := range messages {
		ctx.String(http.StatusOK, "%s\n", msg)
		flusher.Flush()
		time.Sleep(100 * time.Millisecond) // Optional: introduce a small delay between messages
	}
}

func (ctrl *Notification) Update(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, err := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/update.json")
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

	if err := utils.DB.Table(ctrl.Table+"").Where("id = ?", ctx.Param("id")).Updates(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "update "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl *Notification) Delete(ctx *gin.Context) {
	ctrl.Init(ctx)

	if err := utils.DB.Table(ctrl.Table+"").Where("id = ?", ctx.Param("id")).Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+ctrl.SingularLabel+" success", nil))
}
