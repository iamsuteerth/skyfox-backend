package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/nfnt/resize"
	"github.com/rs/zerolog/log"
)

type S3Service interface {
	UploadProfileImage(ctx context.Context, imageBytes []byte, username string, sha string) (string, error)
	DeleteProfileImage(ctx context.Context, objectKey string) error
	ValidateSHA(imageBytes []byte, providedSHA string) bool
	GeneratePresignedURL(ctx context.Context, objectKey string, duration time.Duration) (string, error)
}

type s3Service struct {
	s3Client *s3.S3
	bucket   string
}

func NewS3Service() S3Service {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("S3_BUCKET")

	if accessKey == "" || secretKey == "" || region == "" || bucket == "" {
		log.Error().Msg("S3 configuration missing")
		return nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to create AWS session")
		return nil
	}

	return &s3Service{
		s3Client: s3.New(sess),
		bucket:   bucket,
	}
}

func (s *s3Service) UploadProfileImage(ctx context.Context, imageBytes []byte, username string, providedSHA string) (string, error) {
	if !s.ValidateSHA(imageBytes, providedSHA) {
		return "", utils.NewBadRequestError("INVALID_IMAGE_HASH", "The image hash does not match", nil)
	}

	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode image bytes")
		return "", utils.NewInternalServerError("IMAGE_DECODE_FAILED", "Unable to decode image from byte array", err)
	}

	scaledImg := resize.Resize(64, 64, img, resize.Lanczos3)

	var buf bytes.Buffer

	if err := jpeg.Encode(&buf, scaledImg, &jpeg.Options{Quality: 100}); err != nil {
		log.Error().Err(err).Msg("Failed to encode scaled image to JPEG")
		return "", utils.NewInternalServerError("JPEG_ENCODE_FAILED", "Failed to encode image to JPEG", err)
	}

	timestamp := time.Now().Unix()
	objectKey := fmt.Sprintf("profile-images/%s_%d.jpg", username, timestamp)

	_, err = s.s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String("image/jpeg"),
		ACL:         aws.String("private"),
	})

	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to upload profile image to S3")
		return "", utils.NewInternalServerError("S3_UPLOAD_FAILED", "Failed to upload profile image", err)
	}

	imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		s.bucket,
		*s.s3Client.Config.Region,
		objectKey)

	return imageURL, nil
}

func (s *s3Service) DeleteProfileImage(ctx context.Context, objectKey string) error {
	if strings.HasPrefix(objectKey, "http") {
		parts := strings.Split(objectKey, ".com/")
		if len(parts) < 2 {
			return utils.NewBadRequestError("INVALID_OBJECT_KEY", "Invalid object URL format", nil)
		}
		objectKey = parts[1]
	}

	_, err := s.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		log.Error().Err(err).Str("objectKey", objectKey).Msg("Failed to delete profile image from S3")
		return utils.NewInternalServerError("S3_DELETE_FAILED", "Failed to delete profile image", err)
	}

	return nil
}

func (s *s3Service) ValidateSHA(imageBytes []byte, providedSHA string) bool {
	hash := sha256.Sum256(imageBytes)
	calculatedSHA := hex.EncodeToString(hash[:])

	return calculatedSHA == providedSHA
}

func (s *s3Service) GeneratePresignedURL(ctx context.Context, objectKey string, duration time.Duration) (string, error) {
    if strings.HasPrefix(objectKey, "http") {
        parts := strings.Split(objectKey, ".com/")
        if len(parts) < 2 {
            return "", utils.NewBadRequestError("INVALID_OBJECT_KEY", "Invalid object URL format", nil)
        }
        objectKey = parts[1]
    }

    req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(objectKey),
    })

    presignedURL, err := req.Presign(duration)
    if err != nil {
        log.Error().Err(err).Str("objectKey", objectKey).Msg("Failed to generate presigned URL")
        return "", utils.NewInternalServerError("PRESIGNED_URL_GENERATION_FAILED", "Failed to generate presigned URL", err)
    }

    return presignedURL, nil
}
