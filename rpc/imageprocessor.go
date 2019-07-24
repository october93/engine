package rpc

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	gif "image/gif"
	jpeg "image/jpeg"
	png "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/1l0/identicon"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/kit/log"
)

// Used to extract the file type in order to add it again to the filename after
// generating a unique random filename.
var fileTypeRegexp = regexp.MustCompile("^[a-zA-Z0-9.]*")

// ImageProcessor ties together methods for downloading and processing images.
type ImageProcessor struct {
	aws    *session.Session
	config *Config
	log    log.Logger
	elog   errorLogWriter
}

type errorLogWriter struct {
	log.Logger
}

func (lw errorLogWriter) Write(p []byte) (n int, err error) {
	lw.Logger.Error(errors.New(string(p)))
	return len(p), nil
}

// NewImageProcessor returns a new instance of ImageProcessor and an error if
// it was not able to ensure all the image directories have been created.
func NewImageProcessor(c *Config, l log.Logger) (*ImageProcessor, error) {
	if err := os.MkdirAll(filepath.Join(c.PublicImagesPath, c.ProfileImagesPath), 0700); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(c.PublicImagesPath, c.CardImagesPath), 0700); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(c.PublicImagesPath, c.CardContentImagesPath), 0700); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(c.PublicImagesPath, c.OriginalCardImagesPath), 0700); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(c.PublicImagesPath, c.CoverImagesPath), 0700); err != nil {
		return nil, err
	}
	ip := &ImageProcessor{config: c, log: l, elog: errorLogWriter{Logger: l}}
	if c.S3Upload {
		credentials := credentials.NewStaticCredentials(c.AccessKeyID, c.AccessKeySecret, "")
		s, err := session.NewSession(&aws.Config{Region: &c.S3Region, Credentials: credentials})
		if err != nil {
			return nil, err
		}
		ip.aws = s
	}
	return ip, nil
}

func (ip *ImageProcessor) SaveBase64CardImage(data string) (string, string, error) {
	return ip.saveBase64Image(data, ip.config.OriginalCardImagesPath)
}

func (ip *ImageProcessor) SaveBase64CardContentImage(data string) (string, string, error) {
	return ip.saveBase64Image(data, ip.config.CardContentImagesPath)
}

func (ip *ImageProcessor) SaveBase64ProfileImage(data string) (string, string, error) {
	return ip.saveBase64Image(data, ip.config.ProfileImagesPath)
}

func (ip *ImageProcessor) SaveBase64CoverImage(data string) (string, string, error) {
	return ip.saveBase64Image(data, ip.config.CoverImagesPath)
}

func (ip *ImageProcessor) DownloadCardImage(url string) (string, string, error) {
	return ip.downloadImage(url, ip.config.OriginalCardImagesPath)
}

func (ip *ImageProcessor) DownloadProfileImage(url string) (string, string, error) {
	return ip.downloadImage(url, ip.config.ProfileImagesPath)
}

func (ip *ImageProcessor) GenerateDefaultProfileImage() (string, string, error) {
	identiconGenerator := identicon.New()
	identiconGenerator.Margin = 150

	dest := ip.config.ProfileImagesPath

	filename := fmt.Sprintf("%v.%s", globalid.Next(), "png")
	f, err := os.Create(filepath.Join(ip.config.PublicImagesPath, dest, filename))
	if err != nil {
		return "", "", err
	}
	err = identiconGenerator.GeneratePNG(f)
	if err != nil {
		return "", "", err
	}
	url := fmt.Sprintf("%s/%s/%s/%s", ip.config.ImageHost, ip.config.PublicImagesPath, dest, filename)
	if ip.config.S3Upload {
		buf := new(bytes.Buffer)
		err = identiconGenerator.GeneratePNG(buf)
		if err != nil {
			return "", "", err
		}
		url, err = ip.uploadToS3(buf.Bytes(), dest, filename)
		if err != nil {
			return "", "", err
		}
	}

	return url, fmt.Sprintf("%s/%s/%s", ip.config.PublicImagesPath, dest, filename), nil
}

// SaveBase64Image creates a new file and decodes the Base64 encoded image to
// the file. SaveBase64Image is used to persist images uploaded by the clients.
func (ip *ImageProcessor) saveBase64Image(data, dest string) (string, string, error) {
	delimiterIndex := strings.IndexByte(data, ',')
	if delimiterIndex != -1 {
		data = data[strings.IndexByte(data, ',')+1:]
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	img, format, err := image.Decode(reader)
	if err != nil {
		return "", "", err
	}
	filename := fmt.Sprintf("%v.%s", globalid.Next(), format)
	f, err := os.Create(filepath.Join(ip.config.PublicImagesPath, dest, filename))
	if err != nil {
		return "", "", err
	}
	err = ip.encodeImage(f, img, data, format)
	if err != nil {
		return "", "", err
	}
	url := fmt.Sprintf("%s/%s/%s/%s", ip.config.ImageHost, ip.config.PublicImagesPath, dest, filename)
	if ip.config.S3Upload {
		buf := new(bytes.Buffer)
		err = ip.encodeImage(buf, img, data, format)
		if err != nil {
			return "", "", err
		}
		url, err = ip.uploadToS3(buf.Bytes(), dest, filename)
		if err != nil {
			return "", "", err
		}
	}

	return url, fmt.Sprintf("%s/%s/%s", ip.config.PublicImagesPath, dest, filename), nil
}

func (ip *ImageProcessor) encodeImage(w io.Writer, img image.Image, data, format string) error {
	var err error
	switch format {
	case "jpg", "jpeg":
		err = jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
	case "png":
		err = png.Encode(w, img)
	case "gif":
		err = encodeGIF(w, data)
	}
	return err
}

func (ip *ImageProcessor) BlendImage(src, gradient string) (string, error) {
	// desaturate image (grayscale), increase brightness, decrease contrast
	tmp := filepath.Join(ip.config.PublicImagesPath, ip.config.CardImagesPath, globalid.Next().String())
	cmd := exec.Command("convert", src, "-modulate", "100,0", "-brightness-contrast", "43x-35%", tmp) // #nosec
	cmd.Stderr = ip.elog
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	// increase shadows by 100% and incrase highlights by 50%
	cmd = exec.Command("shadowhighlight", "-sa", "100", "-ha", "50", tmp, tmp) // #nosec
	err = cmd.Run()
	cmd.Stderr = ip.elog
	if err != nil {
		return "", err
	}
	// use identify to read out image width and height
	out, err := exec.Command("identify", tmp).Output() // #nosec
	if err != nil {
		return "", err
	}
	split := strings.Split(string(out), " ")
	if len(split) < 3 {
		return "", fmt.Errorf("identify: unexpected output %s", string(out))
	}
	dimension := fmt.Sprintf("%s!", split[2])
	// resize gradient to match image
	gradientPath := filepath.Join(ip.config.AssetsPath, gradient) + ".png"
	resizedGradientPath := filepath.Join(ip.config.PublicImagesPath, ip.config.CardImagesPath, globalid.Next().String()) + ".png"
	cmd = exec.Command("convert", gradientPath, "-resize", dimension, resizedGradientPath) // #nosec
	cmd.Stderr = ip.elog
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	// color image by blending with gradient background
	filename := fmt.Sprintf("%v.png", globalid.Next())
	dest := filepath.Join(ip.config.PublicImagesPath, ip.config.CardImagesPath, filename)
	cmd = exec.Command("composite", "-compose", "Multiply", "-gravity", "center", resizedGradientPath, tmp, dest) // #nosec
	cmd.Stderr = ip.elog
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	// delete temporary files again
	cleanup, err := filepath.Glob(fmt.Sprintf("%s*", tmp))
	if err != nil {
		return "", err
	}
	cleanup = append(cleanup, resizedGradientPath)
	for _, f := range cleanup {
		err = os.Remove(f)
		if err != nil {
			return "", err
		}
	}
	url := fmt.Sprintf("%s/%s", ip.config.ImageHost, dest)
	if ip.config.S3Upload {
		var data []byte
		data, err = ioutil.ReadFile(dest)
		if err != nil {
			return "", err
		}
		url, err = ip.uploadToS3(data, ip.config.CardImagesPath, filename)
		if err != nil {
			return "", err
		}
	}
	return url, err
}

// GradientImage creates a background image just out of the gradient by
// converting the existing gradient to a jpeg. Returns an error on unknown
// gradients.
func (ip *ImageProcessor) GradientImage(gradient string) (string, string, error) {
	src := filepath.Join(ip.config.AssetsPath, gradient) + ".png"
	if _, err := os.Stat(src); err != nil {
		if os.IsNotExist(err) {
			return "", "", fmt.Errorf("unknown gradient: %s", gradient)
		} else {
			return "", "", err
		}
	}
	url := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s/%s.png", ip.config.S3Region, ip.config.S3Bucket, ip.config.BackgroundImagesPath, gradient)
	return "", url, nil
}

// encodeGIF encodes the image by using the dedicated GIF decoder in order to
// preserve all frames.
func encodeGIF(w io.Writer, data string) error {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	g, err := gif.DecodeAll(reader)
	if err != nil {
		return err
	}
	return gif.EncodeAll(w, g)
}

func (ip *ImageProcessor) downloadImage(url, dest string) (string, string, error) {
	// The comma is a reserved character and sometimes poorly implemented on
	// the server-side. Escaping it is neccesarry to make links work that would
	// otherwise break.
	//
	// https://tools.ietf.org/html/rfc3986#section-2.2
	url = strings.Replace(url, ",", "%2C", -1)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		ip.log.Error(errors.Wrap(errors.New(resp.Status), "downloadImage failed"))
		return "", "", nil
	}

	defer ip.close(resp.Body)
	fileType := fileTypeRegexp.FindString(path.Ext(url))
	filename := fmt.Sprintf("%s%s", globalid.Next().String(), fileType)
	f, err := os.Create(filepath.Join(ip.config.PublicImagesPath, dest, filename))
	defer ip.close(f)
	if err != nil {
		return "", "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	_, err = f.Write(data)
	if err != nil {
		return "", "", err
	}
	url = fmt.Sprintf("%s/%s/%s/%s", ip.config.ImageHost, ip.config.PublicImagesPath, dest, filename)
	if ip.config.S3Upload {
		url, err = ip.uploadToS3(data, dest, filename)
		if err != nil {
			return "", "", err
		}
	}
	return url, fmt.Sprintf("%s/%s/%s", ip.config.PublicImagesPath, dest, filename), nil
}

func (ip *ImageProcessor) close(c io.Closer) {
	err := c.Close()
	if err != nil {
		ip.log.Error(err)
	}
}

func (ip *ImageProcessor) uploadToS3(data []byte, prefix, key string) (string, error) {
	mime := http.DetectContentType(data)
	_, err := s3.New(ip.aws).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(ip.config.S3Bucket),
		Key:                aws.String(fmt.Sprintf("%s/%s", prefix, key)),
		ACL:                aws.String("public-read"),
		ContentType:        aws.String(mime),
		ContentLength:      aws.Int64(int64(len(data))),
		ContentDisposition: aws.String("inline"),
		Body:               bytes.NewReader(data),
	})
	url := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s/%s", ip.config.S3Region, ip.config.S3Bucket, prefix, key)
	return url, err
}
