package httpserver

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"prometheus-test/infrastructure/drivers"
	"prometheus-test/infrastructure/http_client/trace_http"
	"prometheus-test/model"
	"strings"
)

const (
	MethodGET = iota
	MethodPOST
	MethodAll
)

type action struct {
	Path     string
	Method   int
	Handlers []gin.HandlerFunc
}

var actionMaps = make([]*action, 0)

func registerGinHttpAction(path string, method int, handlers ...gin.HandlerFunc) {
	if len(handlers) == 0 {
		panic("action no handlers!")
	}
	actionMaps = append(actionMaps, &action{
		Path:     path,
		Method:   method,
		Handlers: handlers,
	})
}

func init() {
	registerGinHttpAction(
		"/metrics",
		MethodGET,
		func(ctx *gin.Context) {
			promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
		})

	registerGinHttpAction(
		"/GenerateImagesUsingText",
		MethodPOST,
		GenerateImagesUsingText,
	)
	registerGinHttpAction(
		"/getImages",
		MethodGET,
		getImages,
	)
}

func getImages(ctx *gin.Context) {
	// 创建一个 200x200 的红色矩形图像
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	red := color.RGBA{R: 255, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: red}, image.Point{}, draw.Src)

	// 将图像编码为 PNG 格式的二进制数据
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, img)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 发送二进制数据到客户端
	data := buffer.Bytes()
	ctx.Data(http.StatusOK, "image/png", data)
}

func GenerateImagesUsingText(ctx *gin.Context) {
	//简单实现 没有做限流 缓存 并发等场景逻辑处理
	keywords := ctx.Request.PostFormValue("keywords")

	if len(keywords) > 1024 {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	if len(keywords) < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	db, err := drivers.GetMysqlEngineTest(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	//获取关键词列表
	modelKeywords := getKeywords(db, keywords)
	//假设url每次都是唯一的
	modelImage := model.Image{Url: "http://127.0.0.1:9000/getImages"}
	db.Create(&modelImage)
	var modelImageMappings []model.ImageMapping
	//关键词和图片映射绑定 后期利用关键词和映射表查找所有相关的图片
	for _, keyword := range modelKeywords {
		modelImageMappings = append(
			modelImageMappings,
			model.ImageMapping{
				KeywordID: keyword.ID,
				ImageID:   modelImage.ID,
			},
		)
	}
	db.Create(&modelImageMappings)

	Client := trace_http.FetchDefaultTraceClient().GET(ctx, modelImage.Url)
	if Client.Err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	data := Client.Result
	ctx.Data(http.StatusOK, "image/png", data)
}

func getKeywords(db *gorm.DB, text string) map[string]*model.Keyword {

	textSplit := strings.Split(text, " ")
	var md5HashList []string
	var findKeywords []*model.Keyword
	var modelKeywords []*model.Keyword
	var KeywordsMap = make(map[string]*model.Keyword, len(textSplit))
	for _, s := range textSplit {
		md5Hash := md5.Sum([]byte(s))
		md5str := hex.EncodeToString(md5Hash[:])
		if _, ok := KeywordsMap[s]; !ok {
			KeywordsMap[s] = &model.Keyword{
				Content: s,
				MD5:     md5str,
			}
			md5HashList = append(md5HashList, md5str)
		}
	}
	db.Where("md5 IN ?", md5HashList).Find(&findKeywords)
	for _, keyword := range findKeywords {
		if _, ok := KeywordsMap[keyword.Content]; ok {
			if KeywordsMap[keyword.Content].MD5 == keyword.MD5 {
				KeywordsMap[keyword.Content].ID = keyword.ID
			}
		}
	}
	for s := range KeywordsMap {
		if KeywordsMap[s].ID == 0 {
			modelKeywords = append(modelKeywords, KeywordsMap[s])
		}
	}
	if len(modelKeywords) > 0 {
		db.Create(&modelKeywords)
	}
	return KeywordsMap
}
