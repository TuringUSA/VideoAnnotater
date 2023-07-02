package controllers

import (
	"backend/Configs"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	storage "cloud.google.com/go/storage"
	video "cloud.google.com/go/videointelligence/apiv1p3beta1"
	videopb "cloud.google.com/go/videointelligence/apiv1p3beta1/videointelligencepb"
	"github.com/gin-gonic/gin"
	pag "github.com/gobeam/mongo-go-pagination"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Timestamp struct {
	Start int64 `bson:"start"`
	Stop  int64 `bson:"stop"`
}

type Annotation struct {
	Description string      `bson:"description"`
	Timestamps  []Timestamp `bson:"timestamps"`
}

type AnnotationBody struct {
	Keyword string `form:"keyword"`
	Page    int64  `form:"page"`
	PerPage int64  `form:"perPage"`
}

type paginatedResults struct {
	Videos     []Video
	Pagination pag.PaginatedData
}

type Video struct {
	ID          primitive.ObjectID `bson:"_id"`
	Filename    string             `bson:"filename"`
	Size        int                `bson:"size"`
	UUID        string             `bson:"uuid"`
	Annotations []Annotation       `bson:"features"`
	SignedUrl   string
}

func StoreVideo(c *gin.Context) {
	ctx := context.Background()
	client, err := video.NewClient(ctx)
	if err != nil {
		log.Fatal("Failed to create clent: %v", err)
	}

	defer func(client *video.Client) {
		err := client.Close()
		if err != nil {
		}
	}(client)

	fileUuid := uuid.New()

	f, uploadedFile, _ := c.Request.FormFile("video")

	if filepath.Ext(uploadedFile.Filename) != ".mp4" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Uploaded file must be mp4 video",
		})
		return
	}

	//Access cloud storage  bucket
	sc := Configs.StoarageBucket.Object(fileUuid.String()).NewWriter(ctx)

	_, err = io.Copy(sc, f)
	err = sc.Close()
	if err != nil {
		return
	}

	mongo_db := Configs.MongodbClient.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("videos")
	video_file := bson.D{{"filename", uploadedFile.Filename}, {"size", uploadedFile.Size}, {"uuid", fileUuid.String()}}
	mongoRecord, err := mongo_db.InsertOne(context.TODO(), video_file)

	if err != nil {
		panic(err)
	}

	annotationOp, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		InputUri: "gs://" + os.Getenv("GCP_STORAGE_BUCKET") + "/" + fileUuid.String(),
		Annotations: []videopb.Annotation{
			videopb.Feature_LABEL_DETECTION,
		},
	})

	if err != nil {
		log.Fatal("Failed to start annotation operation: %v", err)
	}

	resp, err := annotationOp.Wait(ctx)
	if err != nil {
		log.Fatal("Video Annotation Failed: %v", err)
	}
	result := resp.GetAnnotationResults()[0]

	var annotsList []Annotation
	for _, annotation := range result.SegmentLabelAnnotations {

		var timestamps []Timestamp
		for _, segment := range annotation.Segments {
			start := segment.Segment.StartTimeOffset.AsDuration()
			stop := segment.Segment.EndTimeOffset.AsDuration()
			timestamps = append(timestamps, Timestamp{
				Start: start.Milliseconds(),
				Stop:  stop.Milliseconds(),
			})
		}

		f := Annotation{Description: annotation.Entity.Description, Timestamps: timestamps}
		annotsList = append(annotsList, f)
	}

	filter := bson.D{{"_id", mongoRecord.InsertID}}
	update := bson.D{{"$set", bson.D{{"features", annotsList}}}}
	_, err = mongo_db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"features": annotsList,
	})
}

func getAnnotatedVideo(c *gin.Context) {
	var annotBody AnnotationBody
	err := c.BindQuery(&annotBody)
	if err != nil {
		println(err.Error())
		err = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	collection := Configs.MomgoClient.Database(os.Getenv("MONGO_DB")).Collection("videos")

	indexStage := bson.M{"$search": bson.D{{"index", os.Getenv("MONGO_DB_SEARCH_INDEX")}, {"text", bson.D{{"wildcard", "*"}}}, {"query", annotBody.Keyword}}}

	paginatedData, err := pag.New(collection).Context(context.TODO()).Limit(annotBody.PerPage).Page(annotBody.Page).Aggregate(indexStage)
	if err != nil {
		panic(err)
	}

	storageClient, _ := storage.NewClient(context.Background())

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}

	var output []Video
	for _, raw := range paginatedData.Data {
		var v *Video
		if marshallErr := bson.Unmarshal(raw, &v); marshallErr == nil {
			v.SignedUrl, _ = storageClient.Bucket(os.Getenv("GCP_STORAGE_BUCKET")).SignedURL(v.UUID, opts)
			output = append(output, *v)
		}
	}
	c.JSON(http.StatusOK, paginatedResults{Videos: output, Pagination: paginatedData.Pagination})
}
