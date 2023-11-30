package core

import (
	"context"
	"fmt"
	"io"
	"strings"
)

type (
	Detector interface {
		FaceEmotion(context.Context, io.Reader) (*Emotion, error)
		Close() error
	}
)

type EmotionType string

const (
	UNKNOWN     EmotionType = "unknown"
	POSSIBLE    EmotionType = "possible"
	LIKELY      EmotionType = "likely"
	VERY_LIKELY EmotionType = "very likely"
)

func (et EmotionType) Number() float32 {
	if et == POSSIBLE {
		return 0.4
	}
	if et == LIKELY {
		return 0.7
	}
	if et == VERY_LIKELY {
		return 0.8
	}
	return 0
}

func (et EmotionType) Likelihood() string {
	return string(et)
}

type Emotion struct {
	Anger    EmotionType
	Joy      EmotionType
	Surprise EmotionType
	Sorrow   EmotionType
}

func (em *Emotion) String() string {
	if em == nil {
		return ""
	}
	emotions := []string{
		fmt.Sprintf("%.2f anger", em.Anger.Number()),
		fmt.Sprintf("%.2f joy", em.Joy.Number()),
		fmt.Sprintf("%.2f surprised", em.Surprise.Number()),
		fmt.Sprintf("%.2f sorrowed", em.Surprise.Number()),
	}
	return fmt.Sprintf("the subject emotion is %s", strings.Join(emotions, ", "))
}
