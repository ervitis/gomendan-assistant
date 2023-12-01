package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"

	"github.com/ervitis/gomendan-assistant/pkg/core"
	"github.com/ervitis/gomendan-assistant/pkg/machine_learning/google_vision"
)

func streamCapture(ctx context.Context, capt *gocv.VideoCapture, stream *mjpeg.Stream, classifier gocv.CascadeClassifier, img gocv.Mat, mlClient core.Detector) {
	var (
		buf     *gocv.NativeByteBuffer
		emotion *core.Emotion
		err     error
	)

	for {
		if ok := capt.Read(&img); !ok {
			fmt.Println("device closed")
			return
		}
		if img.Empty() {
			continue
		}

		// get the scanned image and buffer it
		buf, _ = gocv.IMEncode(".jpg", img)

		if capt.IsOpened() && buf != nil && !img.Empty() {
			b := bytes.NewBuffer(buf.GetBytes())
			emotion, err = mlClient.FaceEmotion(ctx, b)
			if err != nil {
				fmt.Printf("error face emotion: %v\n", err)
				time.Sleep(3 * time.Second)
			}
		}

		rects := classifier.DetectMultiScale(img)
		for _, r := range rects {
			gocv.Rectangle(&img, r, color.RGBA{R: 255}, 2)
			size := gocv.GetTextSize(emotion.String(), gocv.FontHersheyPlain, 1.4, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-size.X/2, r.Min.Y-2)
			gocv.PutText(&img, emotion.String(), pt, gocv.FontHersheyPlain, 1.4, color.RGBA{B: 255}, 2)
		}

		// update the image with the rects
		buf, _ = gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf.GetBytes())
		buf.Close()
	}
}

func cleanStream(capt *gocv.VideoCapture, img gocv.Mat, classifier gocv.CascadeClassifier) {
	_ = capt.Close()
	_ = img.Close()
	_ = classifier.Close()
	fmt.Println("capture terminated")
}

func main() {
	done := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signal.Notify(done, os.Kill, os.Interrupt, syscall.SIGTERM)

	mlClient, err := google_vision.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		fmt.Printf("Error capturing data from device: %v\n", err)
		return
	}

	img := gocv.NewMat()
	classifier := gocv.NewCascadeClassifier()
	if !classifier.Load("./data/haarcascade_frontalface_default.xml") {
		fmt.Printf("Error reading file")
		return
	}

	stream := mjpeg.NewStream()
	stream.FrameInterval = 500 * time.Millisecond

	go func() {
		streamCapture(ctx, webcam, stream, classifier, img, mlClient)
	}()

	mux := http.NewServeMux()
	mux.Handle("/mendan", stream)

	// this is a test, the stream will stop after 1 minute
	// a better approach would be using websockets
	srv := &http.Server{
		Addr:         ":8880",
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-done
	_ = mlClient.Close()
	cleanStream(webcam, img, classifier)
	ctxEnd, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxEnd); err != nil {
		fmt.Printf("shutdown server error: %v\n", err)
	}

	fmt.Println("server terminated")
}
