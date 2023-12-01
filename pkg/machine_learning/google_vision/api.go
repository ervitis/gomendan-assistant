package google_vision

import (
	"context"
	"io"

	apivisionv1 "cloud.google.com/go/vision/apiv1"
	apivisionv2 "cloud.google.com/go/vision/v2/apiv1"
	visionv2 "cloud.google.com/go/vision/v2/apiv1/visionpb"

	"github.com/ervitis/gomendan-assistant/pkg/core"
)

type client struct {
	gc *apivisionv2.ImageAnnotatorClient

	reqs []*visionv2.AnnotateImageRequest
}

func (c *client) FaceEmotion(ctx context.Context, bf io.Reader) (*core.Emotion, error) {
	img, err := apivisionv1.NewImageFromReader(bf)
	if err != nil {
		return nil, err
	}

	if len(img.GetContent()) == 0 {
		return nil, nil
	}

	c.reqs[0].Image = img
	annotations, err := c.gc.BatchAnnotateImages(ctx, &visionv2.BatchAnnotateImagesRequest{
		Requests: c.reqs,
	})
	if err != nil {
		return nil, err
	}
	respAnnotations := annotations.GetResponses()
	if len(respAnnotations) == 0 {
		return nil, nil
	}
	if len(respAnnotations[0].FaceAnnotations) == 0 {
		return nil, nil
	}
	faceAnnotation := respAnnotations[0].FaceAnnotations[0]

	emotion := FaceAnnotation(faceAnnotation)
	// print the result
	// fmt.Println(emotion.String())

	return &emotion, nil
}

func (c *client) Close() error {
	return c.gc.Close()
}

func Into(likelihood visionv2.Likelihood) core.EmotionType {
	switch likelihood {
	case visionv2.Likelihood_LIKELY:
		return core.LIKELY
	case visionv2.Likelihood_POSSIBLE:
		return core.POSSIBLE
	case visionv2.Likelihood_VERY_LIKELY:
		return core.VERY_LIKELY
	default:
		return core.UNKNOWN
	}
}

func FaceAnnotation(annotation *visionv2.FaceAnnotation) core.Emotion {
	return core.Emotion{
		Anger:    Into(annotation.GetAngerLikelihood()),
		Joy:      Into(annotation.GetJoyLikelihood()),
		Surprise: Into(annotation.GetSurpriseLikelihood()),
		Sorrow:   Into(annotation.GetSorrowLikelihood()),
	}
}

func NewClient(ctx context.Context) (core.Detector, error) {
	c, err := apivisionv2.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	features := make([]*visionv2.Feature, 1)
	features[0] = &visionv2.Feature{
		MaxResults: 1,
		Type:       visionv2.Feature_FACE_DETECTION,
		Model:      "builtin/stable",
	}
	reqs := make([]*visionv2.AnnotateImageRequest, 1)
	reqs[0] = &visionv2.AnnotateImageRequest{Features: features}
	return &client{
		gc:   c,
		reqs: reqs,
	}, nil
}
